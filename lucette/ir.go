// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// ir.go --- Typed IR.
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
	"net/netip"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Asmodai/gohacks/debug"
)

// * Constants:

// * Variables:

// * Code:

// ** Interface:

type TypedNode interface {
	Key() string
	Debug(...any) *debug.Debug
}

// ** `True` node:

type TypedNodeTrue struct {
}

func (n *TypedNodeTrue) Key() string {
	return "true"
}

// Display debugging information.
func (n *TypedNodeTrue) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Typed TRUE node")

	dbg.Init(params...)
	dbg.End()
	dbg.Print()

	return dbg
}

// ** `False' node:

type TypedNodeFalse struct {
}

func (n *TypedNodeFalse) Key() string {
	return "false"
}

// Display debugging information.
func (n *TypedNodeFalse) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Typed FALSE node")

	dbg.Init(params...)
	dbg.End()
	dbg.Print()

	return dbg
}

// ** `And' node:

type TypedNodeAnd struct {
	kids []TypedNode
}

func (n *TypedNodeAnd) Key() string {
	keys := make([]string, 0, len(n.kids))

	for _, kid := range n.kids {
		keys = append(keys, kid.Key())
	}

	sort.Strings(keys)

	return "and|" + strings.Join(keys, "|")
}

// Display debugging information.
func (n *TypedNodeAnd) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Typed AND Node")

	dbg.Init(params...)
	dbg.Printf("Children:")

	for idx := range n.kids {
		n.kids[idx].Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `Or' node:

type TypedNodeOr struct {
	kids []TypedNode
}

func (n *TypedNodeOr) Key() string {
	keys := make([]string, 0, len(n.kids))

	for _, kid := range n.kids {
		keys = append(keys, kid.Key())
	}

	sort.Strings(keys)

	return "or|" + strings.Join(keys, "|")
}

// Display debugging information.
func (n *TypedNodeOr) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("OR Node")

	dbg.Init(params...)
	dbg.Printf("Children:")

	for idx := range n.kids {
		n.kids[idx].Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `Not' node:

type TypedNodeNot struct {
	kid TypedNode
}

func (n *TypedNodeNot) Key() string {
	return "not|" + n.kid.Key()
}

// Display debugging information.
func (n *TypedNodeNot) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("NOT Node")

	dbg.Init(params...)
	dbg.Printf("Child:")

	if n.kid != nil {
		n.kid.Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `EQ.S' node:

type TypedNodeEqS struct {
	field string
	value string
}

func (n *TypedNodeEqS) Key() string {
	return "eqs|" + n.field + "|" + n.value
}

// Display debugging information.
func (n *TypedNodeEqS) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("EQ.S")

	dbg.Init(params...)
	dbg.Printf("Field: %s", n.field)
	dbg.Printf("Value: %q", n.value)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `NEQ.S' node:

type TypedNodeNeqS struct {
	field string
	value string
}

func (n *TypedNodeNeqS) Key() string {
	return "neqs|" + n.field + "|" + n.value
}

// Display debugging information.
func (n *TypedNodeNeqS) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("NEQ.S")

	dbg.Init(params...)
	dbg.Printf("Field: %s", n.field)
	dbg.Printf("Value: %q", n.value)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `PREFIX' node:

type TypedNodePrefix struct {
	field  string
	prefix string
}

func (n *TypedNodePrefix) Key() string {
	return "prefix|" + n.field + "|" + n.prefix
}

// Display debugging information.
func (n *TypedNodePrefix) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Prefix")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.field)
	dbg.Printf("Prefix: %q", n.prefix)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `REGEX' node:

type TypedNodeRegex struct {
	field    string
	pattern  string
	compiled *regexp.Regexp
}

func (n *TypedNodeRegex) Key() string {
	return "re|" + n.field + "|" + n.pattern
}

// Display debugging information.
func (n *TypedNodeRegex) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Regex")

	dbg.Init(params...)
	dbg.Printf("Field:   %s", n.field)
	dbg.Printf("Pattern: %q", n.pattern)

	if n.compiled != nil {
		dbg.Printf("")
		dbg.Printf("Regex is compiled")
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `GLOB' node:

type TypedNodeGlob struct {
	field string
	glob  string
}

func (n *TypedNodeGlob) Key() string {
	return "glob|" + n.field + "|" + n.glob
}

// Display debugging information.
func (n *TypedNodeGlob) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Glob")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.field)
	dbg.Printf("Glob:   %q", n.glob)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `PHRASE' node:

type TypedNodePhrase struct {
	field  string
	phrase string
	prox   int
	fuzz   *float64
	boost  *float64
}

func (n *TypedNodePhrase) Key() string {
	return "phrase|" + n.field + "|" + n.phrase
}

// Display debugging information.
func (n *TypedNodePhrase) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Phrase")

	dbg.Init(params...)
	dbg.Printf("Field:     %s", n.field)
	dbg.Printf("Phrase:    %q", n.phrase)
	dbg.Printf("Proximity: %d", n.prox)

	if n.fuzz != nil {
		dbg.Printf("Fuzziness: %g", *n.fuzz)
	}

	if n.boost != nil {
		dbg.Printf("Boost:     %g", *n.boost)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `EXISTS' node:

type TypedNodeExists struct {
	field string
}

func (n *TypedNodeExists) Key() string {
	return "ex|" + n.field
}

// Display debugging information.
func (n *TypedNodeExists) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Exists")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.field)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `CmpN' node:

type TypedNodeCmpN struct {
	field string
	op    CmpKind
	value float64
}

func (n *TypedNodeCmpN) Key() string {
	return fmt.Sprintf("cn|%s|%d|%g", n.field, n.op, n.value)
}

// Display debugging information.
func (n *TypedNodeCmpN) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Cmp.N")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.field)
	dbg.Printf("Op:     %s", cmpKindToString(n.op))
	dbg.Printf("Value:  %f", n.value)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `RangeN' node:

type TypedNodeRangeN struct {
	field string
	low   *float64
	high  *float64
	incl  bool
	inch  bool
}

func (n *TypedNodeRangeN) Key() string {
	var low, high string

	if n.low != nil {
		low = fmt.Sprintf("%g", *n.low)
	}

	if n.high != nil {
		high = fmt.Sprintf("%g", *n.high)
	}

	return fmt.Sprintf("rn|%s|%s|%s|%t|%t",
		n.field,
		low,
		high,
		n.incl,
		n.inch)
}

// Display debugging information.
func (n *TypedNodeRangeN) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Range.N")

	dbg.Init(params...)
	dbg.Printf("Field:          %s", n.field)

	if n.low != nil {
		dbg.Printf("Low:            %f", *n.low)
	}

	if n.high != nil {
		dbg.Printf("High:           %f", *n.high)
	}

	dbg.Printf("Increment Low:  %v", n.incl)
	dbg.Printf("Increment High: %v", n.inch)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `CmpT' node:

type TypedNodeCmpT struct {
	field string
	op    CmpKind
	value int64
}

func (n *TypedNodeCmpT) Key() string {
	return fmt.Sprintf("ct|%s|%d|%d", n.field, n.op, n.value)
}

// Display debugging information.
func (n *TypedNodeCmpT) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Cmp.T")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.field)
	dbg.Printf("Op:     %s", cmpKindToString(n.op))
	dbg.Printf("Value:  %d", n.value)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `RangeT' node:

type TypedNodeRangeT struct {
	field string
	low   *int64
	high  *int64
	incl  bool
	inch  bool
}

func (n *TypedNodeRangeT) Key() string {
	var low, high string

	if n.low != nil {
		low = strconv.FormatInt(*n.low, 10)
	}

	if n.high != nil {
		high = strconv.FormatInt(*n.high, 10)
	}

	return fmt.Sprintf("rt|%s|%s|%s|%t|%t",
		n.field,
		low,
		high,
		n.incl,
		n.inch)
}

// Display debugging information.
func (n *TypedNodeRangeT) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Range.T")

	dbg.Init(params...)
	dbg.Printf("Field:          %s", n.field)

	if n.low != nil {
		dbg.Printf("Low:            %d", *n.low)
	}

	if n.high != nil {
		dbg.Printf("High:           %d", *n.high)
	}

	dbg.Printf("Increment Low:  %v", n.incl)
	dbg.Printf("Increment High: %v", n.inch)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `CmpIP' node:

type TypedNodeCmpIP struct {
	field string
	op    CmpKind
	value netip.Addr
}

func (n *TypedNodeCmpIP) Key() string {
	return fmt.Sprintf("cip|%s|%d|%s", n.field, n.op, n.value.String())
}

// Display debugging information.
func (n *TypedNodeCmpIP) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Cmp.IP")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.field)
	dbg.Printf("Op:     %s", cmpKindToString(n.op))
	dbg.Printf("Value:  %s", n.value.String())

	dbg.End()
	dbg.Print()

	return dbg
}

// ** `RangeIP' node:

type TypedNodeRangeIP struct {
	field string
	low   netip.Addr
	high  netip.Addr
	incl  bool
	inch  bool
}

func (n *TypedNodeRangeIP) Key() string {
	return fmt.Sprintf("rip|%s|%s|%s|%t|%t",
		n.field,
		n.low.String(),
		n.high.String(),
		n.incl,
		n.inch)
}

// Display debugging information.
func (n *TypedNodeRangeIP) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Range.IP")

	dbg.Init(params...)
	dbg.Printf("Field:          %s", n.field)
	dbg.Printf("Low:            %s", n.low.String())
	dbg.Printf("High:           %s", n.high.String())
	dbg.Printf("Increment Low:  %v", n.incl)
	dbg.Printf("Increment High: %v", n.inch)

	dbg.End()
	dbg.Print()

	return dbg
}

// * ir.go ends here.
