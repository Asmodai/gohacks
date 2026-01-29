// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// parser_stack.go --- Stack methods.
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

// * Constants:

// * Variables:

// * Code:

// Push an operator to the stack.
func (p *parser) pushOp(kind exprKind, aux any, span *Span) {
	p.stack = append(
		p.stack,
		opEntry{kind: kind, aux: aux, span: span})
}

// Pop an operator from the stack.
func (p *parser) popOp() (opEntry, bool) {
	if len(p.stack) == 0 {
		return opEntry{}, false
	}

	elt := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]

	return elt, true
}

// Return the operator on the top of the stack.
func (p *parser) topOp() (opEntry, bool) {
	if len(p.stack) == 0 {
		return opEntry{}, false
	}

	return p.stack[len(p.stack)-1], true
}

// Push a node to the AST.
func (p *parser) pushNode(node ASTNode) {
	p.ast = append(p.ast, node)
}

// Pop a node from the AST.
func (p *parser) popNode() ASTNode {
	if len(p.ast) == 0 {
		p.errorf(NewEmptySpan(), "stack underflow")

		return &ASTPredicate{Kind: PredicateANY}
	}

	end := len(p.ast) - 1
	node := p.ast[end]
	p.ast = p.ast[:end]

	return node
}

// * parser_stack.go ends here.
