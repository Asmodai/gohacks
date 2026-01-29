// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// astmodifier.go --- AST `modifier' node.
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
	"github.com/Asmodai/gohacks/debug"
)

// * Constants:

const (
	// This modifier will mark the predicate to which it is attached as
	// having a required value.  If the value is not present, then the
	// predicate will fail.
	ModRequire ModifierKind = iota

	// This modifier will mark the predicate to which it is attached as
	// having a prohibited value.  If the value is present, then the
	// predicate will fail.
	ModProhibit
)

// * Variables:

var (
	// Map of `ModKind -> string` for pretty-printing.
	//
	//nolint:gochecknoglobals
	modKindStrings = map[ModifierKind]string{
		ModRequire:  "Require",
		ModProhibit: "Prohibit",
	}
)

// * Code:

// ** Type:

// Modifier kind type.
type ModifierKind int

// ** Structure:

// An AST node for a `modifier' to an operation..
type ASTModifier struct {
	Kid  ASTNode      // Node to which the modifier applies.
	span *Span        // Source code span.
	Kind ModifierKind // Modifier kind.
}

// ** Methods:

// Return the span for the AST node.
func (n ASTModifier) Span() *Span {
	return n.span
}

// Display debugging information.
func (n ASTModifier) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("AST 'Modifier' Node")

	dbg.Init(params...)
	dbg.Printf("Kind: %v", ModifierKindToString(n.Kind))
	dbg.Printf("Span: %s", n.span.String())
	dbg.Printf("Child:")

	if n.Kid != nil {
		n.Kid.Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// ** Functions:

// Return the string representation of a modifier.
func ModifierKindToString(kind ModifierKind) string {
	if str, found := modKindStrings[kind]; found {
		return str
	}

	return invalidStr
}

// * astmodifier.go ends here.
