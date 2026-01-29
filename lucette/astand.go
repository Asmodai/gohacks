// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// astand.go --- AST `AND' node.
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

// * Variables:

// * Code:

// ** Type:

// An AST node for the `AND' logical operator.
type ASTAnd struct {
	Kids []ASTNode // Child nodes.
	span *Span     // Source code span.
}

// ** Methods:

// Return the span for the AST node.
func (n ASTAnd) Span() *Span {
	return n.span
}

// Display debugging information.
func (n ASTAnd) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("AST 'AND' Node")

	dbg.Init(params...)

	dbg.Printf("Span: %s", n.span.String())
	dbg.Printf("Children:")

	for idx := range n.Kids {
		n.Kids[idx].Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// * astand.go ends here.
