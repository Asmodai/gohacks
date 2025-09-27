// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// astliteral.go --- AST `Literal' node.
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
	"github.com/Asmodai/gohacks/debug"
)

// * Constants:

const (
	LString    LiteralKind = iota // Literal is a string.
	LNumber                       // Literal is a number.
	LUnbounded                    // Literal has no value bound to it.
)

// * Variables:

var (
	// Map of `LiteralKind -> string` for pretty-printing.
	//
	//nolint:gochecknoglobals
	litKindStrings = map[LiteralKind]string{
		LString:    "String",
		LNumber:    "Number",
		LUnbounded: "Unbounded",
	}
)

// * Code:

// ** Types:

// Literal kind type.
type LiteralKind int

// ** Structure:

// An AST node for a `literal' of some kind.
type ASTLiteral struct {
	String string      // String value.
	span   *Span       // Source code span.
	Kind   LiteralKind // Kind of the literal.
	Number float64     // Numeric value.
}

// ** Methods:

// Return the span for the AST node.
func (n ASTLiteral) Span() *Span {
	return n.span
}

// Display debugging information.
func (n ASTLiteral) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("AST 'Literal' Node")

	dbg.Init(params...)
	dbg.Printf("Kind:   %v", LiteralKindToString(n.Kind))
	dbg.Printf("Span:   %s", n.span.String())
	dbg.Printf("String: %s", n.String)
	dbg.Printf("Number: %g", n.Number)

	dbg.End()
	dbg.Print()

	return dbg
}

// ** Functions:

// Return the string representation of a literal type.
func LiteralKindToString(lit LiteralKind) string {
	if str, found := litKindStrings[lit]; found {
		return str
	}

	return invalidStr
}

// * astliteral.go ends here.
