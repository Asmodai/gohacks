// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// astcomparator.go --- AST comparator type.
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

import "github.com/Asmodai/gohacks/debug"

// * Constants:

const (
	ComparatorLT  ComparatorKind = iota // Comparator is `LT'.
	ComparatorLTE                       // Comparator is `LTE'.
	ComparatorGT                        // Comparator is `GT'.
	ComparatorGTE                       // Comparator is `GTE'.
	ComparatorEQ                        // Comparator is `EQ'.
	ComparatorNEQ                       // Comparator is `NEQ'.
)

// * Variables:

//nolint:gochecknoglobals
var (
	// Map of `ComparatorKind -> string` for pretty-printing.
	cmpKindStrings = map[ComparatorKind]string{
		ComparatorLT:  "LT",
		ComparatorLTE: "LTE",
		ComparatorGT:  "GT",
		ComparatorGTE: "GTE",
		ComparatorEQ:  "EQ",
		ComparatorNEQ: "NEQ",
	}

	// Map of `ComparatorKind -> ComparatorKind` for inversion.
	cmpInverse = map[ComparatorKind]ComparatorKind{
		ComparatorLT:  ComparatorGTE,
		ComparatorLTE: ComparatorGT,
		ComparatorGT:  ComparatorLTE,
		ComparatorGTE: ComparatorLT,
		ComparatorEQ:  ComparatorNEQ,
		ComparatorNEQ: ComparatorEQ,
	}
)

// * Code:

// ** Types:

// Comparator kind type.
type ComparatorKind int

// ** Structure:

// Comparator structure.
type ASTComparator struct {
	Atom ASTLiteral     // Atom on which to operate.
	Op   ComparatorKind // Comparator operator.
}

// ** Methods:

// Display debugging information.
func (c ASTComparator) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Range")

	dbg.Init(params...)
	dbg.Printf("Operator: %s", ComparatorKindToString(c.Op))

	c.Atom.Debug(&dbg, "Atom")

	dbg.End()
	dbg.Print()

	return dbg
}

// ** Functions:

// Return the string representation for a comparator kind.
func ComparatorKindToString(kind ComparatorKind) string {
	if str, found := cmpKindStrings[kind]; found {
		return str
	}

	return invalidStr
}

func InvertComparator(kind ComparatorKind) ComparatorKind {
	return cmpInverse[kind]
}

// * astcomparator.go ends here.
