// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// typer.go --- Generate Typed IR.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
//mock:yes

// * Comments:

// * Package:

package lucette

// * Imports:

import (
	"fmt"
	"net/netip"
	"strconv"
	"time"

	"github.com/Asmodai/gohacks/conversion"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	//nolint:gochecknoglobals
	zeroIP = netip.Addr{}
)

// * Code:

// ** Interface:

type Typer interface {
	// Generate typed IR from the given AST root node.
	Type(ASTNode) (IRNode, []Diagnostic)

	// Return the diagnostic messages generated during IR generation.
	Diagnostics() []Diagnostic
}

// ** Structure:

type typer struct {
	Sch   Schema
	Diags []Diagnostic
}

// ** Utility methods:

func (t *typer) addDiag(span *Span, msg string, args ...any) {
	t.Diags = append(t.Diags,
		NewDiagnostic(fmt.Sprintf(msg, args...), span))
}

func (t *typer) Diagnostics() []Diagnostic {
	return t.Diags
}

// ** IR methods:

// Generate an IR `And' node.
func (t *typer) makeIRAnd(node *ASTAnd) IRNode {
	kids := make([]IRNode, 0, len(node.Kids))

	for _, kid := range node.Kids {
		kids = append(kids, t.typeIR(kid))
	}

	return &IRAnd{Kids: kids}
}

// Generate an IR `Or' node.
func (t *typer) makeIROr(node *ASTOr) IRNode {
	kids := make([]IRNode, 0, len(node.Kids))

	for _, kid := range node.Kids {
		kids = append(kids, t.typeIR(kid))
	}

	return &IROr{Kids: kids}
}

func (t *typer) makeIRModifier(node *ASTModifier) IRNode {
	// -x => NOT X.
	if node.Kind == ModProhibit {
		return &IRNot{Kid: t.typeIR(node.Kid)}
	}

	// +X => X.
	return t.typeIR(node.Kid)
}

func (t *typer) makeIRPredicate(node *ASTPredicate) IRNode {
	spec, ok := t.Sch[node.Field]
	if !ok {
		t.addDiag(node.span, "unknown field %q", node.Field)

		// Treat as a keyword so we can at least proceed.
		spec = FieldSpec{Name: node.Field, FType: FTKeyword}
	}

	return t.typeLeaf(node, spec)
}

func (t *typer) typeCompare(pred *ASTPredicate, spec FieldSpec) IRNode {
	//nolint:exhaustive
	switch spec.FType {
	case FTNumeric:
		val := pickNumber(pred.Comparator.Atom, t, pred.span)

		return &IRNumberCmp{
			Field: spec.Name,
			Op:    pred.Comparator.Op,
			Value: val}

	case FTDateTime:
		str, _, _ := stringOrNumber(pred.Comparator.Atom)

		val, err := toEpoch(str, spec.Layouts)
		if err != nil {
			t.addDiag(pred.span, "%v", err)

			val = 0
		}

		return &IRTimeCmp{
			Field: spec.Name,
			Op:    pred.Comparator.Op,
			Value: val}

	case FTIP:
		str, _, _ := stringOrNumber(pred.Comparator.Atom)

		val, err := toIP(str)
		if err != nil {
			t.addDiag(pred.span, "%v", err)
		}

		return &IRIPCmp{
			Field: spec.Name,
			Op:    pred.Comparator.Op,
			Value: val}

	default:
		t.addDiag(pred.span,
			"%q comparators not supported on field %q",
			FieldTypeToString(spec.FType),
			spec.Name)

		return &IRStringEQ{
			Field: spec.Name,
			Value: fmt.Sprintf("%v", pred.Comparator.Atom)}
	}
}

func (t *typer) typeStringEQ(pred *ASTPredicate, spec FieldSpec) IRNode {
	if spec.FType != FTText && spec.FType != FTKeyword {
		return t.typeCompare(pred, spec)
	}

	return &IRStringEQ{
		Field: spec.Name,
		Value: pred.String}
}

func (t *typer) typeStringNEQ(pred *ASTPredicate, spec FieldSpec) IRNode {
	if spec.FType != FTText && spec.FType != FTKeyword {
		return t.typeCompare(pred, spec)
	}

	return &IRStringNEQ{
		Field: spec.Name,
		Value: pred.String}
}

func (t *typer) typeAny(_ *ASTPredicate, spec FieldSpec) IRNode {
	return &IRAny{Field: spec.Name}
}

func (t *typer) typeGlob(pred *ASTPredicate, spec FieldSpec) IRNode {
	if spec.FType != FTText {
		t.addDiag(pred.span,
			"glob only supported for text fields")

		return &IRAny{Field: spec.Name}
	}

	return &IRGlob{
		Field: spec.Name,
		Glob:  pred.String}
}

func (t *typer) typePhrase(pred *ASTPredicate, spec FieldSpec) IRNode {
	if spec.FType != FTText && pred.Proximity != 0 {
		t.addDiag(pred.span,
			"proximity only supported for text fields")
	}

	if spec.FType == FTIP {
		ipaddr, err := toIP(pred.String)
		if err != nil {
			t.addDiag(pred.span,
				"ip parse failed: %v",
				err)

			goto proceedAsString
		}

		return &IRIPCmp{
			Field: spec.Name,
			Op:    ComparatorEQ,
			Value: ipaddr}
	}

proceedAsString:
	return &IRPhrase{
		Field:     spec.Name,
		Phrase:    pred.String,
		Proximity: pred.Proximity,
		Fuzz:      pred.Fuzz,
		Boost:     pred.Boost}
}

func (t *typer) typePrefix(pred *ASTPredicate, spec FieldSpec) IRNode {
	if spec.FType != FTKeyword && spec.FType != FTText {
		t.addDiag(pred.span,
			"prefix not supported on field %q",
			spec.Name)

		return &IRAny{Field: spec.Name}
	}

	return &IRPrefix{
		Field:  spec.Name,
		Prefix: pred.String}
}

func (t *typer) typeRange(pred *ASTPredicate, spec FieldSpec) IRNode {
	//nolint:exhaustive
	switch spec.FType {
	case FTNumeric:
		low := pickNumberPtr(pred.Range.Lo, t, pred.span)
		high := pickNumberPtr(pred.Range.Hi, t, pred.span)

		return &IRNumberRange{
			Field: spec.Name,
			Lo:    low,
			Hi:    high,
			IncL:  pred.Range.IncL,
			IncH:  pred.Range.IncH}

	case FTDateTime:
		low := pickTimePtr(pred.Range.Lo, t, pred.span, spec.Layouts)
		high := pickTimePtr(pred.Range.Hi, t, pred.span, spec.Layouts)

		return &IRTimeRange{
			Field: spec.Name,
			Lo:    low,
			Hi:    high,
			IncL:  pred.Range.IncL,
			IncH:  pred.Range.IncH}

	case FTIP:
		low := pickIPPtr(pred.Range.Lo, t, pred.span)
		high := pickIPPtr(pred.Range.Hi, t, pred.span)

		return &IRIPRange{
			Field: spec.Name,
			Lo:    low,
			Hi:    high,
			IncL:  pred.Range.IncL,
			IncH:  pred.Range.IncH}

	default:
		t.addDiag(pred.span,
			"ranges not supported on field %q",
			spec.Name)

		return &IRStringEQ{Field: spec.Name, Value: ""}
	}
}

func (t *typer) typeRegex(pred *ASTPredicate, spec FieldSpec) IRNode {
	if spec.FType != FTKeyword && spec.FType != FTText {
		t.addDiag(pred.span,
			"regex not supported on field %q",
			spec.Name)

		return &IRAny{Field: spec.Name}
	}

	return &IRRegex{
		Field:    spec.Name,
		Pattern:  pred.Regex,
		Compiled: pred.compiled}
}

func (t *typer) typeLeaf(pred *ASTPredicate, spec FieldSpec) IRNode {
	switch pred.Kind {
	case PredicateCMP:
		return t.typeCompare(pred, spec)

	case PredicateEQS:
		return t.typeStringEQ(pred, spec)

	case PredicateNEQS:
		return t.typeStringNEQ(pred, spec)

	case PredicateANY:
		return t.typeAny(pred, spec)

	case PredicateGLOB:
		return t.typeGlob(pred, spec)

	case PredicatePHRASE:
		return t.typePhrase(pred, spec)

	case PredicatePREFIX:
		return t.typePrefix(pred, spec)

	case PredicateRANGE:
		return t.typeRange(pred, spec)

	case PredicateREGEX:
		return t.typeRegex(pred, spec)
	}

	t.addDiag(pred.span,
		"unhandled predicate kind %q",
		PredicateKindToString(pred.Kind))

	return &IRAny{Field: spec.Name}
}

func (t *typer) typeIR(node ASTNode) IRNode {
	switch val := node.(type) {
	case *ASTAnd:
		return t.makeIRAnd(val)

	case *ASTOr:
		return t.makeIROr(val)

	case *ASTNot:
		return &IRNot{Kid: t.typeIR(val.Kid)}

	case *ASTModifier:
		return t.makeIRModifier(val)

	case *ASTPredicate:
		return t.makeIRPredicate(val)

	default:
		t.addDiag(ZeroSpan, "internal: unknown node type")

		return &IRAnd{Kids: nil}
	}
}

// ** Methods:

func (t *typer) Type(node ASTNode) (IRNode, []Diagnostic) {
	out := t.typeIR(node)

	return out, t.Diags
}

// ** Functions:

// Convert a string to a Unix epoch.
func toEpoch(val string, layouts []string) (int64, error) {
	for _, layout := range layouts {
		if tval, err := time.Parse(layout, val); err == nil {
			return tval.UnixNano(), nil
		}
	}

	// Try RFC3339 as a fallback.
	if tval, err := time.Parse(time.RFC3339, val); err == nil {
		return tval.UnixNano(), nil
	}

	// Try a Unix timestamp.
	if ival, err := strconv.ParseInt(val, 10, 64); err == nil {
		return time.Unix(ival, 0).UnixNano(), nil
	}

	return 0, errors.WithMessagef(
		ErrBadDateTime,
		"%q",
		val)
}

// Convert a string to an IP address.
func toIP(val string) (netip.Addr, error) {
	addr, err := netip.ParseAddr(val)
	if err != nil {
		return zeroIP, errors.WithStack(err)
	}

	return addr, nil
}

// Extract either a string or a number from a literal.
//
//nolint:unparam
func stringOrNumber(lit ASTLiteral) (string, *float64, error) {
	switch lit.Kind {
	case LString:
		return lit.String, nil, nil

	case LNumber:
		val := lit.Number
		sval, _ := conversion.ToString(val)

		return sval, &val, nil

	case LUnbounded:
		return "", nil, nil

	default:
		return "", nil, errors.WithStack(ErrUnknownLiteral)
	}
}

// Extract a number from a literal.
//
//nolint:exhaustive
func pickNumber(lit ASTLiteral, inst *typer, span *Span) float64 {
	switch lit.Kind {
	case LNumber:
		return lit.Number

	case LString:
		val, err := strconv.ParseFloat(lit.String, 64)
		if err != nil {
			inst.addDiag(span, "not a number: %q", lit.String)

			return 0
		}

		return val

	default:
		inst.addDiag(span, "unbounded not allowed here")

		return 0
	}
}

func pickNumberPtr(lit *ASTLiteral, inst *typer, span *Span) *float64 {
	if lit == nil || lit.Kind == LUnbounded {
		return nil
	}

	val := pickNumber(*lit, inst, span)

	return &val
}

func pickTimePtr(lit *ASTLiteral, inst *typer, span *Span, layouts []string) *int64 {
	if lit == nil || lit.Kind == LUnbounded {
		return nil
	}

	str, _, _ := stringOrNumber(*lit)

	val, err := toEpoch(str, layouts)
	if err != nil {
		inst.addDiag(span, "%v", err)

		return nil
	}

	return &val
}

func pickIPPtr(lit *ASTLiteral, inst *typer, span *Span) netip.Addr {
	if lit == nil || lit.Kind == LUnbounded {
		return zeroIP
	}

	str, _, _ := stringOrNumber(*lit)

	val, err := toIP(str)
	if err != nil {
		inst.addDiag(span, "%v", err)

		return zeroIP
	}

	return val
}

func NewTyper(sch Schema) Typer {
	return &typer{Sch: sch}
}

// * typer.go ends here.
