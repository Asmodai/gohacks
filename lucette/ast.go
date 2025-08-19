// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// ast.go --- Abstract syntax tree for Lucette.
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

// * Imports:

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Asmodai/gohacks/debug"
)

// * Constants{

const (
	invalidStr = "INVALID"
)

// * Code:

// ** Node interface:

// AST node.
type Node interface {
	// Return the span for this node.
	//
	// Spans can be used in diagnostics to show where in the source file
	// an issue exists.
	Span() Span

	// Print debugging information for the given node.
	Debug(...any) *debug.Debug
}

// ** Span type:

// Source code span.
type Span struct {
	start Position
	end   Position
}

// Return the string representation of a span.
func (s *Span) String() string {
	var sbld strings.Builder

	sbld.WriteString("From: ")
	sbld.WriteString(s.start.String())
	sbld.WriteString("  To: ")
	sbld.WriteString(s.end.String())

	return sbld.String()
}

// ** `And' type:

// An AST node for the `AND' logical operator.
type NodeAnd struct {
	kids []Node
	span Span
}

// Return the span for the AND within the source code.
func (n *NodeAnd) Span() Span {
	return n.span
}

// Display debugging information.
func (n *NodeAnd) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("AND Node")

	dbg.Init(params...)
	dbg.Printf("Span: %s", n.span.String())
	dbg.Printf("Children:")

	for idx := range n.kids {
		n.kids[idx].Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `Or' type:

// An AST node for the `OR' logical operator.
type NodeOr struct {
	kids []Node
	span Span
}

// Return the span for the OR within the source code.
func (n *NodeOr) Span() Span {
	return n.span
}

// Display debugging information.
func (n *NodeOr) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("OR Node")

	dbg.Init(params...)
	dbg.Printf("Span: %s", n.span.String())
	dbg.Printf("Children:")

	for idx := range n.kids {
		n.kids[idx].Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `Not' type:

// An AST node for the `NOT' logical operator.
type NodeNot struct {
	kid  Node
	span Span
}

// Return the span for the NOT within the source code.
func (n *NodeNot) Span() Span {
	return n.span
}

// Display debugging information.
func (n *NodeNot) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("NOT Node")

	dbg.Init(params...)
	dbg.Printf("Span: %s", n.span.String())
	dbg.Printf("Child:")

	if n.kid != nil {
		n.kid.Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** Modifiers:

// Modifier kind type.
type ModKind int

const (
	ModRequire ModKind = iota
	ModProhibit
)

// Return the string representation of a modifier.
func modKindToString(kind ModKind) string {
	var modKindString = map[ModKind]string{
		ModRequire:  "Require",
		ModProhibit: "Prohibit",
	}

	if str, found := modKindString[kind]; found {
		return str
	}

	return invalidStr
}

// An AST node representing a modifier.
type NodeMod struct {
	kind ModKind
	kid  Node
	span Span
}

// Return the span for the modifier.
func (n *NodeMod) Span() Span {
	return n.span
}

// Display debugging information.
func (n *NodeMod) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Modifier Node")

	dbg.Init(params...)
	dbg.Printf("Kind: %v", modKindToString(n.kind))
	dbg.Printf("Span: %s", n.span.String())
	dbg.Printf("Child:")

	if n.kid != nil {
		n.kid.Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** Literal for ranges.

// Literal kind type.
type LitKind int

const (
	LString LitKind = iota
	LNumber
	LUnbounded
)

// Return the string representation of a literal type.
func litKindToString(lit LitKind) string {
	var litKindStrings = map[LitKind]string{
		LString:    "String",
		LNumber:    "Number",
		LUnbounded: "Unbounded",
	}

	if str, found := litKindStrings[lit]; found {
		return str
	}

	return invalidStr
}

// An AST node representing a literal.
type NodeLit struct {
	kind   LitKind
	strval string
	numval float64
	span   Span
}

// Return the span for the node.
func (n *NodeLit) Span() Span {
	return n.span
}

// Display debugging information.
func (n *NodeLit) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Literal Node")

	dbg.Init(params...)
	dbg.Printf("Kind:   %v", litKindToString(n.kind))
	dbg.Printf("Span:   %s", n.span.String())
	dbg.Printf("String: %s", n.strval)
	dbg.Printf("Number: %f", n.numval)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** Leaf predicates:

// Predicate kind type.
type PredKind int

// Comparator kind type.
type CmpKind int

//nolint:revive,stylecheck
const (
	PK_EQ_S PredKind = iota
	PK_NEQ_S
	PK_PREFIX
	PK_GLOB
	PK_REGEX
	PK_PHRASE
	PK_EXISTS
	PK_CMP
	PK_RANGE

	CmpLT CmpKind = iota
	CmpLTE
	CmpGT
	CmpGTE
	CmpEQ
	CmpNEQ
)

// Return the string representation for a predicate kind.
func predKindToString(kind PredKind) string {
	var predKindStrings = map[PredKind]string{
		PK_EQ_S:   "EQ.S",
		PK_NEQ_S:  "NEQ.S",
		PK_PREFIX: "PREFIX",
		PK_GLOB:   "GLOB",
		PK_REGEX:  "REGEX",
		PK_PHRASE: "PHRASE",
		PK_EXISTS: "EXISTS",
		PK_CMP:    "<comparator>",
		PK_RANGE:  "<range>",
	}

	if str, found := predKindStrings[kind]; found {
		return str
	}

	return invalidStr
}

// Return the string representation for a comparator kind.
func cmpKindToString(kind CmpKind) string {
	var cmpKindStrings = map[CmpKind]string{
		CmpLT:  "LT",
		CmpLTE: "LTE",
		CmpGT:  "GT",
		CmpGTE: "GTE",
		CmpEQ:  "EQ",
		CmpNEQ: "NEQ",
	}

	if str, found := cmpKindStrings[kind]; found {
		return str
	}

	return fmt.Sprintf("%s [%d]", invalidStr, kind)
}

// *** Comparator:

// Comparator structure.
type Comparator struct {
	Op   CmpKind // Comparison operator.
	Atom NodeLit // Atom to compare.
}

// *** Range:

// Range structure.
type Range struct {
	Low  *NodeLit // Low literal.
	High *NodeLit // High literal.
	IncL bool     // Low is inclusive?
	IncH bool     // High is inclusive?
}

// *** Predicate:

// An AST node representing a predicate.
type NodePred struct {
	kind   PredKind
	field  string
	strval string
	numval float64
	reval  string
	repat  *regexp.Regexp
	cmp    *Comparator
	rnge   *Range
	prox   int
	fuzz   *float64
	boost  *float64
	span   Span
}

// Return the span.
func (n *NodePred) Span() Span {
	return n.span
}

// Display debugging information.
func (n *NodePred) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Predicate Node")

	dbg.Init(params...)
	dbg.Printf("Kind:       %v", predKindToString(n.kind))
	dbg.Printf("Span:       %s", n.span.String())
	dbg.Printf("Field:      %s", n.field)
	dbg.Printf("String:     %q", n.strval)
	dbg.Printf("Number:     %f", n.numval)
	dbg.Printf("Regexp:     %q", n.reval)

	dbg.Printf("Proximity:  %d", n.prox)

	if n.fuzz != nil {
		dbg.Printf("Fuzz:       %g", *n.fuzz)
	}

	if n.boost != nil {
		dbg.Printf("Boost:      %g", *n.boost)
	}

	if n.cmp != nil {
		dbg.Printf("")
		dbg.Printf("Comparator: %s", cmpKindToString(n.cmp.Op))
		n.cmp.Atom.Debug(&dbg, "Atom")
	}

	if n.rnge != nil {
		dbg.Printf("")
		dbg.Printf("Range:")
		n.rnge.Low.Debug(&dbg, "Range Low")
		n.rnge.High.Debug(&dbg, "Range High")
		dbg.Printf("Increment Low:  %v", n.rnge.IncL)
		dbg.Printf("Increment High: %v", n.rnge.IncH)
	}

	if n.repat != nil {
		dbg.Printf("")
		dbg.Printf("Regex is compiled")
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// * ast.go ends here.
