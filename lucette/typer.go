// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// typer.go --- Semantic typer.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
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

// * Comments:

// * Package:

package lucette

import (
	"fmt"
	"net/netip"
	"strconv"
	"time"

	"github.com/Asmodai/gohacks/conversion"
	"gitlab.com/tozd/go/errors"
)

// * Imports:

// * Constants:

const (
	FTKeyword FieldType = iota
	FTText
	FTNumeric
	FTDateTime
	FTIP
)

// * Variables:

var (
	ErrBadDateTime    = errors.Base("bad datetime")
	ErrUnknownLiteral = errors.Base("unknown literal")

	//nolint:gochecknoglobals
	zeroIP = netip.Addr{}
)

// * Code:

// ** Types:

type FieldType int

type Schema map[string]FieldSpec

// ** Field specification:

type FieldSpec struct {
	Name     string
	FType    FieldType
	Analyser string
	Layouts  []string
}

func ftypeToString(ftype FieldType) string {
	var ftypeString = map[FieldType]string{
		FTKeyword:  "Keyword",
		FTText:     "Text",
		FTNumeric:  "Numeric",
		FTDateTime: "Datetime",
		FTIP:       "IP",
	}

	if str, found := ftypeString[ftype]; found {
		return str
	}

	return invalidStr
}

// ** Typer:

type Typer struct {
	Sch   Schema
	Diags []Diagnostic
}

// ** Methods:

func (t *Typer) diag(span Span, msg string, args ...any) {
	t.Diags = append(
		t.Diags,
		Diagnostic{
			Msg: fmt.Sprintf(msg, args...),
			At:  span})
}

func (t *Typer) Type(node Node) (TypedNode, []Diagnostic) {
	out := t.typeNode(node)

	return out, t.Diags
}

//nolint:cyclop
func (t *Typer) typeNode(node Node) TypedNode {
	switch val := node.(type) {
	case *NodeAnd:
		kids := make([]TypedNode, 0, len(val.kids))

		for _, kid := range val.kids {
			kids = append(kids, t.typeNode(kid))
		}

		return &TypedNodeAnd{kids: kids}

	case *NodeOr:
		kids := make([]TypedNode, 0, len(val.kids))

		for _, kid := range val.kids {
			kids = append(kids, t.typeNode(kid))
		}

		return &TypedNodeOr{kids: kids}

	case *NodeNot:
		return &TypedNodeNot{kid: t.typeNode(val.kid)}

	case *NodeMod:
		// -X => NOT X
		if val.kind == ModProhibit {
			return &TypedNodeNot{kid: t.typeNode(val.kid)}
		}

		// +X => X.
		return t.typeNode(val.kid)

	case *NodePred:
		spec, ok := t.Sch[val.field]
		if !ok && len(val.field) > 0 {
			t.diag(val.span, "unknown field %q", val.field)

			// Treat unknown as a keyword so we can proceed.
			spec = FieldSpec{Name: val.field, FType: FTKeyword}
		}

		return t.typeLeaf(val, spec)

	default:
		t.diag(Span{}, "internal: unknown node")

		return &TypedNodeAnd{kids: nil}
	}
}

//nolint:funlen
func (t *Typer) typeEqS(pred *NodePred, spec FieldSpec) TypedNode {
	switch spec.FType {
	case FTKeyword, FTText:
		return &TypedNodeEqS{
			field: spec.Name,
			value: pred.strval}

	case FTNumeric:
		val, err := strconv.ParseFloat(pred.strval, 64)
		if err == nil {
			return &TypedNodeCmpN{
				field: spec.Name,
				op:    CmpEQ,
				value: val}
		}

		t.diag(pred.span,
			"field %q is numeric; %q is not a number",
			spec.Name,
			pred.strval)

		return &TypedNodeEqS{
			field: spec.Name,
			value: pred.strval}

	case FTDateTime:
		val, err := toEpoch(pred.strval, spec.Layouts)
		if err == nil {
			return &TypedNodeCmpT{
				field: spec.Name,
				op:    CmpEQ,
				value: val}
		}

		t.diag(pred.span,
			"datetime parse failed for %q",
			pred.strval)

		return &TypedNodeEqS{
			field: spec.Name,
			value: pred.strval}

	case FTIP:
		ipaddr, err := toIP(pred.strval)
		if err != nil {
			t.diag(pred.span,
				"ip parse failed: %v",
				err)

			return &TypedNodeEqS{
				field: spec.Name,
				value: pred.strval}
		}

		return &TypedNodeCmpIP{
			field: spec.Name,
			op:    CmpEQ,
			value: ipaddr}

	default:
		t.diag(pred.span,
			"unknown string equality comparison")

		return &TypedNodeEqS{
			field: spec.Name,
			value: pred.strval}
	}
}

//nolint:exhaustive
func (t *Typer) typeCmp(pred *NodePred, spec FieldSpec) TypedNode {
	switch spec.FType {
	case FTNumeric:
		val := pickNumber(pred.cmp.Atom, t, pred.span)

		return &TypedNodeCmpN{
			field: spec.Name,
			op:    pred.cmp.Op,
			value: val}

	case FTDateTime:
		str, _, _ := stringOrNumber(pred.cmp.Atom)

		val, err := toEpoch(str, spec.Layouts)
		if err != nil {
			t.diag(pred.span, "%v", err)

			val = 0
		}

		return &TypedNodeCmpT{
			field: spec.Name,
			op:    pred.cmp.Op,
			value: val}

	case FTIP:
		str, _, _ := stringOrNumber(pred.cmp.Atom)

		val, err := toIP(str)
		if err != nil {
			t.diag(pred.span, "%v", err)
		}

		return &TypedNodeCmpIP{
			field: spec.Name,
			op:    pred.cmp.Op,
			value: val}

	default:
		t.diag(pred.span,
			"%q comparators not supported on field %q",
			ftypeToString(spec.FType),
			spec.Name)

		return &TypedNodeEqS{
			field: spec.Name,
			value: fmt.Sprintf("%v", pred.cmp.Atom)}
	}
}

//nolint:exhaustive
func (t *Typer) typeRange(pred *NodePred, spec FieldSpec) TypedNode {
	switch spec.FType {
	case FTNumeric:
		low := pickNumberPtr(pred.rnge.Low, t, pred.span)
		high := pickNumberPtr(pred.rnge.High, t, pred.span)

		return &TypedNodeRangeN{
			field: spec.Name,
			low:   low,
			high:  high,
			incl:  pred.rnge.IncL,
			inch:  pred.rnge.IncH}

	case FTDateTime:
		low := pickTimePtr(pred.rnge.Low, t, pred.span, spec.Layouts)
		high := pickTimePtr(pred.rnge.High, t, pred.span, spec.Layouts)

		return &TypedNodeRangeT{
			field: spec.Name,
			low:   low,
			high:  high,
			incl:  pred.rnge.IncL,
			inch:  pred.rnge.IncH}

	case FTIP:
		low := pickIPPtr(pred.rnge.Low, t, pred.span)
		high := pickIPPtr(pred.rnge.High, t, pred.span)

		return &TypedNodeRangeIP{
			field: spec.Name,
			low:   low,
			high:  high,
			incl:  pred.rnge.IncL,
			inch:  pred.rnge.IncH}

	default:
		t.diag(pred.span,
			"ranges not supported on field %q",
			spec.Name)

		return &TypedNodeEqS{
			field: spec.Name,
			value: ""}
	}
}

func (t *Typer) typeRegex(pred *NodePred, spec FieldSpec) TypedNode {
	if spec.FType != FTKeyword && spec.FType != FTText {
		t.diag(pred.span,
			"regex not supported on field %q",
			spec.Name)
	}

	return &TypedNodeRegex{
		field:    spec.Name,
		pattern:  pred.reval,
		compiled: pred.repat}
}

func (t *Typer) typePhrase(pred *NodePred, spec FieldSpec) TypedNode {
	if spec.FType != FTText && pred.prox != 0 {
		t.diag(pred.span,
			"proximity only supported for text fields")
	}

	return &TypedNodePhrase{
		field:  spec.Name,
		phrase: pred.strval,
		prox:   pred.prox,
		fuzz:   pred.fuzz,
		boost:  pred.boost}
}

func (t *Typer) typePrefix(pred *NodePred, spec FieldSpec) TypedNode {
	if spec.FType != FTKeyword && spec.FType != FTText {
		t.diag(pred.span, "prefix unsupported on field %q", spec.Name)
	}

	return &TypedNodePrefix{
		field:  spec.Name,
		prefix: pred.strval}
}

//nolint:exhaustive
func (t *Typer) typeLeaf(pred *NodePred, spec FieldSpec) TypedNode {
	switch pred.kind {
	case PK_EXISTS:
		return &TypedNodeExists{field: spec.Name}

	case PK_EQ_S:
		return t.typeEqS(pred, spec)

	case PK_CMP:
		return t.typeCmp(pred, spec)

	case PK_RANGE:
		return t.typeRange(pred, spec)

	case PK_REGEX:
		return t.typeRegex(pred, spec)

	case PK_PHRASE:
		return t.typePhrase(pred, spec)

	case PK_PREFIX:
		return t.typePrefix(pred, spec)
	}

	t.diag(pred.span,
		"unhandled predicate kind %q",
		predKindToString(pred.kind))

	return &TypedNodeExists{field: spec.Name}
}

// ** Functions:

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

// NOTE: Replace this... should have net.Addr in the typed node.
func toIP(val string) (netip.Addr, error) {
	addr, err := netip.ParseAddr(val)
	if err != nil {
		return zeroIP, errors.WithStack(err)
	}

	return addr, nil
}

//nolint:unparam
func stringOrNumber(lit NodeLit) (string, *float64, error) {
	switch lit.kind {
	case LString:
		return lit.strval, nil, nil

	case LNumber:
		val := lit.numval
		sval, _ := conversion.ToString(val)

		return sval, &val, nil

	case LUnbounded:
		return "", nil, nil

	default:
		return "", nil, errors.WithStack(ErrUnknownLiteral)
	}
}

//nolint:exhaustive
func pickNumber(lit NodeLit, inst *Typer, span Span) float64 {
	switch lit.kind {
	case LNumber:
		return lit.numval

	case LString:
		val, err := strconv.ParseFloat(lit.strval, 64)
		if err != nil {
			inst.diag(span, "not a number: %q", lit.strval)

			return 0
		}

		return val

	default:
		inst.diag(span, "unbounded not allowed here")

		return 0
	}
}

func pickNumberPtr(lit *NodeLit, inst *Typer, span Span) *float64 {
	if lit == nil || lit.kind == LUnbounded {
		return nil
	}

	val := pickNumber(*lit, inst, span)

	return &val
}

func pickTimePtr(lit *NodeLit, inst *Typer, span Span, layouts []string) *int64 {
	if lit == nil || lit.kind == LUnbounded {
		return nil
	}

	str, _, _ := stringOrNumber(*lit)

	val, err := toEpoch(str, layouts)
	if err != nil {
		inst.diag(span, "%v", err)

		return nil
	}

	return &val
}

func pickIPPtr(lit *NodeLit, inst *Typer, span Span) netip.Addr {
	if lit == nil || lit.kind == LUnbounded {
		return zeroIP
	}

	str, _, _ := stringOrNumber(*lit)

	val, err := toIP(str)
	if err != nil {
		inst.diag(span, "%v", err)

		return zeroIP
	}

	return val
}

func NewTyper(sch Schema) *Typer {
	return &Typer{Sch: sch}
}

// * typer.go ends here.
