// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// parser.go --- The parser.
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

//nolint:unused
package lucette

// * Imports:

import ()

// * Constants:

const (
	stateExpectPrimary parserState = iota // Expecting primary expr.
	stateAfterPrimary                     // Not a primary expr.

	exprOr          exprKind = iota // OR expression.
	exprAnd                         // AND expression.
	exprImplicitAnd                 // Implicit AND expression.
	exprNot                         // NOT expression.
	exprRequire                     // Required term expression.
	exprProhibit                    // Prohibited term expresssion.
	exprFieldApply                  // Field load expression.
	exprLParen                      // Left parenthesis.

	precOR         bindPrec = 10 // OR precedence.
	precAND        bindPrec = 20 // AND precedence.
	precNOT        bindPrec = 30 // NOT precedence.
	precFieldApply bindPrec = 85 // Field Apply precedence.
)

// * Variables:

var (
	//nolint:gochecknoglobals
	primaryStarters = map[Token]struct{}{
		TokenPhrase:   struct{}{},
		TokenNumber:   struct{}{},
		TokenRegex:    struct{}{},
		TokenField:    struct{}{},
		TokenNot:      struct{}{},
		TokenLParen:   struct{}{},
		TokenLBracket: struct{}{},
		TokenLCurly:   struct{}{},
		TokenLT:       struct{}{},
		TokenLTE:      struct{}{},
		TokenGT:       struct{}{},
		TokenGTE:      struct{}{},
	}

	//nolint:gochecknoglobals,unused
	primaryEnders = map[Token]struct{}{
		TokenPhrase:   struct{}{},
		TokenNumber:   struct{}{},
		TokenRegex:    struct{}{},
		TokenRParen:   struct{}{},
		TokenRBracket: struct{}{},
		TokenRCurly:   struct{}{},
	}

	//nolint:gochecknoglobals
	precedences = map[exprKind]bindPrec{
		exprOr:          precOR,
		exprAnd:         precAND,
		exprImplicitAnd: precAND,
		exprNot:         precNOT,
		exprRequire:     precNOT,
		exprProhibit:    precNOT,
		exprFieldApply:  precFieldApply,
	}

	//nolint:gochecknoglobals
	comparators = map[Token]ComparatorKind{
		TokenLT:  ComparatorLT,
		TokenLTE: ComparatorLTE,
		TokenGT:  ComparatorGT,
		TokenGTE: ComparatorGTE,
	}
)

// * Code:

// ** Types:

// Parser state.
type parserState int

// Expression type.
type exprKind int

// Binding precedence.
type bindPrec int

// Operation stack entry.
type opEntry struct {
	kind exprKind // Expression type.
	aux  any      // Auxiliary data.
	span *Span    // Source code span.
}

// ** Structure:

type parser struct {
	tokens []LexedToken // Token list from the lexer.
	index  int          // Current position index.
	ast    []ASTNode    // Abstract syntax tree nodes.
	stack  []opEntry    // Operation stack.
	state  parserState  // Parser state.
	diags  []Diagnostic // Diagnostic messages.
}

// ** Methods:

// Reset the parser state.
func (p *parser) Reset() {
	p.tokens = make([]LexedToken, 0)
	p.ast = make([]ASTNode, 0)
	p.stack = make([]opEntry, 0, 1)
	p.diags = make([]Diagnostic, 0, 1)
	p.index = 0
	p.state = stateExpectPrimary
}

// Return a list of diagnostic messages.
func (p *parser) Diagnostics() []Diagnostic {
	return p.diags
}

// Parse the list of lexed tokens and generate an AST.
//
//nolint:cyclop
func (p *parser) Parse(lexed []LexedToken) (ASTNode, []Diagnostic) {
	p.Reset()

	if len(lexed) < 1 {
		p.errorf(NewEmptySpan(), "no tokens provided")

		return nil, p.diags
	}

	p.tokens = lexed

	for {
		tok := p.next()

		// If we are at EOF, then go to the drain.
		if tok.Token == TokenEOF {
			goto DRAIN
		}

		switch p.state {
		case stateExpectPrimary:
			p.parsePrimary(tok)

		case stateAfterPrimary:
			p.parseAfterPrimary(tok)
		}
	}

DRAIN:
	for {
		if len(p.stack) == 0 {
			break
		}

		top, _ := p.topOp()
		if top.kind == exprLParen {
			p.errorf(top.span, "unmatched '('")

			_, _ = p.popOp()

			continue
		}

		p.reduceOne()
	}

	if len(p.ast) == 0 {
		return &ASTPredicate{Kind: PredicateANY, span: NewEmptySpan()}, p.diags
	}

	if len(p.ast) != 1 {
		p.errorf(p.stack[0].span, "internal: multiple roots after reduce")
	}

	return p.ast[0], p.diags
}

// ** Functions:

// Generate s span for the given token.
func spanTok(tok LexedToken) *Span {
	return NewSpan(tok.Start, tok.End)
}

// Does the given token start a primary expression?
func startsPrimary(token Token) bool {
	_, found := primaryStarters[token]

	return found
}

// Does the given token end a primary expression?
func endsPrimary(token Token) bool {
	_, found := primaryEnders[token]

	return found
}

// Attach a boost to the node.
func attachBoost(node ASTNode, val float64) ASTNode {
	if pred, ok := node.(*ASTPredicate); ok {
		res := *pred

		res.Boost = &val

		return &res
	}

	return node
}

// Attach a fuzz to the node.
func attachFuzz(node ASTNode, val *float64) ASTNode {
	if pred, ok := node.(*ASTPredicate); ok {
		res := *pred

		res.Fuzz = val

		return &res
	}

	return node
}

// Create a new parser instance.
func NewParser() Parser {
	return &parser{
		state: stateExpectPrimary,
	}
}

// * parser.go ends here.
