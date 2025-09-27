// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// parser_parse.go --- Main parser methods.
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

// Handle any required postfix translation.
//
//nolint:forcetypeassert,exhaustive
func (p *parser) handlePostFix() {
	for {
		switch p.peek().Token {
		case TokenCaret:
			p.next()

			numTok := p.expect(TokenNumber, "number after '^'")
			top := p.popNode()

			p.pushNode(attachBoost(top, numTok.Literal.Value.(float64)))

		case TokenTilde:
			p.next()

			var flt *float64

			if p.peek().Token == TokenNumber {
				val := p.next().Literal.Value.(float64)
				flt = &val
			}

			top := p.popNode()
			p.pushNode(attachFuzz(top, flt))

		default:
			return
		}
	}
}

// Parse a primary expression.
//
//nolint:forcetypeassert,exhaustive
func (p *parser) parsePrimary(tok LexedToken) {
	switch tok.Token {
	case TokenNot:
		p.reduceWhile(precNOT)
		p.pushOp(exprNot, nil, spanTok(tok))

	case TokenPlus:
		p.reduceWhile(precNOT)
		p.pushOp(exprRequire, nil, spanTok(tok))

	case TokenMinus:
		p.reduceWhile(precNOT)
		p.pushOp(exprProhibit, nil, spanTok(tok))

	case TokenField:
		p.expect(TokenColon, "':' after field")
		p.reduceWhile(precFieldApply)

		fld := tok.Literal.Value.(string)
		p.pushOp(exprFieldApply, fld, spanTok(tok))

	case TokenLParen:
		p.pushOp(exprLParen, nil, spanTok(tok))

	case TokenPhrase, TokenNumber, TokenRegex:
		p.pushNode(p.makePredForm(tok))
		p.handlePostFix()

		p.state = stateAfterPrimary

	case TokenLBracket, TokenLCurly, TokenLT, TokenLTE, TokenGT, TokenGTE:
		p.pushNode(p.parseRangeOrCmp(tok))
		p.handlePostFix()

		p.state = stateAfterPrimary

	default:
		p.errorf(spanTok(tok), "expected primary")
	}
}

// Parse an expression that comes after a primary.
//
//nolint:cyclop,forcetypeassert,exhaustive
func (p *parser) parseAfterPrimary(tok LexedToken) {
	switch tok.Token {
	case TokenCaret:
		numTok := p.expect(TokenNumber, "number after '^'")
		top := p.popNode()

		p.pushNode(attachBoost(top, numTok.Literal.Value.(float64)))

	case TokenTilde:
		var flt *float64

		if p.peek().Token == TokenNumber {
			val := p.next().Literal.Value.(float64)
			flt = &val
		}

		top := p.popNode()
		p.pushNode(attachFuzz(top, flt))

	case TokenAnd:
		p.reduceWhile(precAND)
		p.pushOp(exprAnd, nil, spanTok(tok))
		p.state = stateExpectPrimary

	case TokenOr:
		p.reduceWhile(precOR)
		p.pushOp(exprOr, nil, spanTok(tok))
		p.state = stateExpectPrimary

	case TokenRParen:
		for {
			top, ok := p.topOp()
			if !ok {
				p.errorf(spanTok(tok), "unmatched ')'")

				break
			}

			if top.kind == exprLParen {
				break
			}

			p.reduceOne()
		}

		if top, ok := p.popOp(); !ok || top.kind != exprLParen {
			p.errorf(spanTok(tok), "unmatched ')'")
		}

	default:
		if startsPrimary(tok.Token) {
			p.unread()
			p.reduceWhile(precAND)
			p.pushOp(exprImplicitAnd, nil, ZeroSpan)
			p.state = stateExpectPrimary
		} else {
			p.errorf(spanTok(tok), "expected operator")
		}
	}
}

// * parser_parse.go ends here.
