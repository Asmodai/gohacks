// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// parser_primary.go --- Primary builders for the parser.
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
	"regexp"
)

// * Code:

// ** Methods:

// Generate a regex AST node.
func (p *parser) makeRegex(tok LexedToken) ASTNode {
	span := spanTok(tok)

	pattern, ok := tok.Literal.Value.(string)
	if !ok {
		p.errorf(span,
			"invalid regular expression value: %v",
			tok.Literal.Value)

		return &ASTPredicate{Kind: PredicateANY, span: span}
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		p.errorf(span,
			"regular expression compile failed: %q",
			err.Error())

		return &ASTPredicate{Kind: PredicateANY, span: span}
	}

	return &ASTPredicate{
		Kind:     PredicateREGEX,
		Regex:    pattern,
		compiled: compiled,
		span:     span}
}

// Make a predicate form from the given token.
//
//nolint:forcetypeassert
func (p *parser) makePredForm(tok LexedToken) ASTNode {
	span := spanTok(tok)

	//nolint:exhaustive
	switch tok.Token {
	case TokenPhrase:
		return &ASTPredicate{
			Kind:   PredicatePHRASE,
			String: tok.Literal.Value.(string),
			span:   span}

	case TokenNumber:
		return &ASTPredicate{
			Kind: PredicateCMP,
			span: span,
			Comparator: &ASTComparator{
				Op: ComparatorEQ,
				Atom: ASTLiteral{
					Kind:   LNumber,
					Number: tok.Literal.Value.(float64),
					span:   span}}}

	case TokenRegex:
		return p.makeRegex(tok)

	default:
		p.errorf(span, "unexpected primary")

		return &ASTPredicate{Kind: PredicateANY, span: span}
	}
}

// Expect that the next token in the reader is one that can be used within
// a range.
//
//nolint:forcetypeassert
func (p *parser) expectRangeAtom() ASTLiteral {
	tok := p.next()
	span := spanTok(tok)

	//nolint:exhaustive
	switch tok.Token {
	case TokenPhrase:
		return ASTLiteral{
			Kind:   LString,
			String: tok.Literal.Value.(string),
			span:   span}

	case TokenNumber:
		return ASTLiteral{
			Kind:   LNumber,
			Number: tok.Literal.Value.(float64),
			span:   span}

	default:
		p.errorf(span, "expected range atom (string or number)")

		return ASTLiteral{
			Kind:   LString,
			String: "",
			span:   span}
	}
}

// Read in a range bound from the token stream.
func (p *parser) readRangeBound() *ASTLiteral {
	if p.accept(TokenStar) {
		tok := p.peek()
		span := spanTok(tok)

		return &ASTLiteral{
			Kind: LUnbounded,
			span: span,
		}
	}

	atom := p.expectRangeAtom()

	return &atom
}

// Parse a RANGE or a comparator.
func (p *parser) parseRangeOrCmp(first LexedToken) ASTNode {
	span := spanTok(first)

	//nolint:exhaustive
	switch first.Token {
	case TokenLT, TokenLTE, TokenGT, TokenGTE:
		atom := p.expectRangeAtom()
		op := comparators[first.Token]

		return &ASTPredicate{
			Kind: PredicateCMP,
			Comparator: &ASTComparator{
				Op:   op,
				Atom: atom},
			span: NewSpan(span.start, atom.Span().end)}
	}

	low := p.readRangeBound()
	p.expect(TokenTo, "TO")
	high := p.readRangeBound()
	end := p.next()

	if end.Token != TokenRBracket && end.Token != TokenRCurly {
		p.errorf(spanTok(end), "unterminated interval")
	}

	incL := (first.Token == TokenLBracket)
	incH := (end.Token == TokenRBracket)

	return &ASTPredicate{
		Kind: PredicateRANGE,
		Range: &ASTRange{
			Lo:   low,
			Hi:   high,
			IncL: incL,
			IncH: incH},
		span: NewSpan(span.start, spanTok(end).end)}
}

// * parser_primary.go ends here.
