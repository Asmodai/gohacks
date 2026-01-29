// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// astpredicate.go --- AST `Predicate` node.
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

// * Comments:

// * Package:

package lucette

// * Imports:

import (
	"regexp"

	"github.com/Asmodai/gohacks/debug"
)

// * Constants:

const (
	PredicateCMP    PredicateKind = iota // Predicate is a comparator.
	PredicateEQS                         // Predicate is `EQ.S'.
	PredicateANY                         // Predicate is `ANY'.
	PredicateGLOB                        // Predicate is `GLOB'.
	PredicateNEQS                        // Predicate is `NEQ.S'.
	PredicatePHRASE                      // Predicate is `PHRASE'.
	PredicatePREFIX                      // Predicate is `PREFIX'.
	PredicateRANGE                       // Predicate is `RANGE'.
	PredicateREGEX                       // Predicate is `REGEX'.
)

// * Variables:

var (
	// Map of `PredicateKind -> string` for pretty-printing.
	//
	//nolint:gochecknoglobals
	predKindStrings = map[PredicateKind]string{
		PredicateEQS:    "EQ.S",
		PredicateNEQS:   "NEQ.S",
		PredicatePREFIX: "PREFIX",
		PredicateGLOB:   "GLOB",
		PredicateREGEX:  "REGEX",
		PredicatePHRASE: "PHRASE",
		PredicateANY:    "ANY",
		PredicateRANGE:  "RANGE",
		PredicateCMP:    "<comparator>",
	}
)

// * Code:

// ** Types:

// Predicate kind type.
type PredicateKind int

// ** Structure:

// An AST node for predicates.
type ASTPredicate struct {
	Range      *ASTRange      // Target range value
	compiled   *regexp.Regexp // Compiled regex pattern.
	Comparator *ASTComparator // Comparator to use.
	Fuzz       *float64       // Levenshtein Distance.
	Boost      *float64       // Boost value.
	Field      string         // Target field.
	String     string         // Target string value.
	Regex      string         // Target regex pattern.
	span       *Span          // Source code span.
	Kind       PredicateKind  // Predicate kind.
	Number     float64        // Target numeric value.
	Proximity  int            // String promity.
}

// ** Methods:

// Return the span for the AST node.
func (n ASTPredicate) Span() *Span {
	return n.span
}

// Display debugging information.
func (n ASTPredicate) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("AST 'Predicate' Node")

	dbg.Init(params...)
	dbg.Printf("Kind:       %v", PredicateKindToString(n.Kind))
	dbg.Printf("Span:       %s", n.span.String())
	dbg.Printf("Field:      %s", n.Field)
	dbg.Printf("String:     %q", n.String)
	dbg.Printf("Number:     %f", n.Number)
	dbg.Printf("Regexp:     %q", n.Regex)
	dbg.Printf("Proximity:  %d", n.Proximity)

	if n.Fuzz != nil {
		dbg.Printf("Fuzz:       %g", *n.Fuzz)
	}

	if n.Boost != nil {
		dbg.Printf("Boost:      %g", *n.Boost)
	}

	if n.Comparator != nil {
		n.Comparator.Debug(&dbg, "Comparator")
	}

	if n.Range != nil {
		n.Range.Debug(&dbg, "Range")
	}

	if n.compiled != nil {
		dbg.Printf("")
		dbg.Printf("Regex is compiled")
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** Functions:

// Return the string representation for a predicate kind.
func PredicateKindToString(kind PredicateKind) string {
	if str, found := predKindStrings[kind]; found {
		return str
	}

	return invalidStr
}

// * astpredicate.go ends here.
