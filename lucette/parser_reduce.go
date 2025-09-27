// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// parser_reduce.go --- Parser reduction methods.
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

// * Code:

// ** Methods:

// Reduce an operator from the stack.
//
//nolint:forcetypeassert
func (p *parser) reduceOne() {
	operand, ok := p.popOp()
	if !ok {
		return
	}

	switch operand.kind {
	case exprOr:
		right, left := p.popNode(), p.popNode()
		p.pushNode(reduceOr(left, right))

	case exprAnd, exprImplicitAnd:
		right, left := p.popNode(), p.popNode()
		p.pushNode(reduceAnd(left, right))

	case exprNot:
		kid := p.popNode()
		p.pushNode(&ASTNot{
			Kid:  kid,
			span: NewSpan(operand.span.start, kid.Span().end)})

	case exprRequire:
		kid := p.popNode()
		p.pushNode(&ASTModifier{
			Kind: ModRequire,
			Kid:  kid,
			span: NewSpan(operand.span.start, kid.Span().end)})

	case exprProhibit:
		kid := p.popNode()
		p.pushNode(&ASTModifier{
			Kind: ModProhibit,
			Kid:  kid,
			span: NewSpan(operand.span.start, kid.Span().end)})

	case exprFieldApply:
		fld := operand.aux.(string)
		kid := p.popNode()

		p.pushNode(applyField(
			fld,
			kid,
			NewSpan(operand.span.start, kid.Span().end)))

	case exprLParen: // Barrier.
	}
}

// Reduce operators from the top of the stack until the given minimum
// binding precedence is met.
func (p *parser) reduceWhile(minPrec bindPrec) {
	for {
		top, topOk := p.topOp()
		if !topOk {
			return
		}

		if top.kind == exprLParen {
			return
		}

		prec := precedence(top.kind)
		if prec < minPrec {
			return
		}

		p.reduceOne()
	}
}

// ** Functions:

// Return the precedence of the operator kind.
func precedence(kind exprKind) bindPrec {
	prec, found := precedences[kind]

	if !found {
		return 0
	}

	return prec
}

// Reduce an AND expression.
func reduceAnd(lhs, rhs ASTNode) ASTNode {
	splhs := lhs.Span()
	sprhs := rhs.Span()

	if andlhs, ok := lhs.(*ASTAnd); ok {
		if andrhs, ok := rhs.(*ASTAnd); ok {
			return &ASTAnd{
				Kids: append(andlhs.Kids, andrhs.Kids...),
				span: NewSpan(splhs.start, sprhs.end)}
		}

		return &ASTAnd{
			Kids: append(andlhs.Kids, rhs),
			span: NewSpan(splhs.start, sprhs.end)}
	}

	if andrhs, ok := rhs.(*ASTAnd); ok {
		return &ASTAnd{
			Kids: append([]ASTNode{lhs}, andrhs.Kids...),
			span: NewSpan(splhs.start, sprhs.end)}
	}

	return &ASTAnd{
		Kids: []ASTNode{lhs, rhs},
		span: NewSpan(splhs.start, sprhs.end)}
}

// Reduce an OR expression.
func reduceOr(lhs, rhs ASTNode) ASTNode {
	slhs := lhs.Span()
	srhs := rhs.Span()

	if orlhs, ok := lhs.(*ASTOr); ok {
		if orrhs, ok := rhs.(*ASTOr); ok {
			return &ASTOr{
				Kids: append(orlhs.Kids, orrhs.Kids...),
				span: NewSpan(slhs.start, srhs.end)}
		}

		return &ASTOr{
			Kids: append(orlhs.Kids, rhs),
			span: NewSpan(slhs.start, srhs.end)}
	}

	if orrhs, ok := rhs.(*ASTOr); ok {
		return &ASTOr{
			Kids: append([]ASTNode{lhs}, orrhs.Kids...),
			span: NewSpan(slhs.start, srhs.end)}
	}

	return &ASTOr{
		Kids: []ASTNode{lhs, rhs},
		span: NewSpan(slhs.start, srhs.end)}
}

// Apply the given field to the given node.
//
// Applies the field to the child nodes where applicable.
func applyField(field string, node ASTNode, span *Span) ASTNode {
	switch val := node.(type) {
	case *ASTPredicate:
		res := *val
		res.Field = field
		res.span = span

		return &res

	case *ASTAnd:
		kids := make([]ASTNode, len(val.Kids))

		for idx, elt := range val.Kids {
			kids[idx] = applyField(
				field,
				elt,
				NewSpan(val.span.start, elt.Span().end))
		}

		return &ASTAnd{Kids: kids, span: span}

	case *ASTOr:
		kids := make([]ASTNode, len(val.Kids))

		for idx, elt := range val.Kids {
			kids[idx] = applyField(
				field,
				elt,
				NewSpan(val.span.start, elt.Span().end))
		}

		return &ASTOr{Kids: kids, span: span}

	case *ASTNot:
		return &ASTNot{
			Kid:  applyField(field, val.Kid, val.span),
			span: span}

	case *ASTModifier:
		return &ASTModifier{
			Kind: val.Kind,
			Kid:  applyField(field, val.Kid, val.span),
			span: span}

	default:
		return node
	}
}

// * parser_reduce.go ends here.
