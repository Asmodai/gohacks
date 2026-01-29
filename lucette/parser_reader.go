// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// parser_reader.go --- Parser reader methods.
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

import "fmt"

// * Code:

// Generate a diagnostic message.
func (p *parser) errorf(span *Span, format string, args ...any) {
	p.diags = append(p.diags, NewDiagnostic(fmt.Sprintf(format, args...), span))
}

// Return the current token without advancing the token reader.
func (p *parser) peek() LexedToken {
	return p.tokens[p.index]
}

// Get the next token.
func (p *parser) next() LexedToken {
	tok := p.tokens[p.index]
	p.index++

	return tok
}

// Unread a token back to the reader.
func (p *parser) unread() {
	if p.index > 0 {
		p.index--
	}
}

// Advance the token reader if the current token is of the given type.
func (p *parser) accept(tType Token) bool {
	if p.peek().Token == tType {
		p.index++

		return true
	}

	return false
}

// Expect the current token to be of the given type.
//
// If it is not, then generate a diagnostic message.
func (p *parser) expect(tType Token, msg string) LexedToken {
	tok := p.next()

	if tok.Token != tType {
		p.errorf(spanTok(tok), "expected %s", msg)

		return NewLexedToken(tType, "", tok.Start, tok.End)
	}

	return tok
}

// Return the previous token type.
//
//nolint:unused
func (p *parser) prevType() Token {
	if p.index == 0 {
		return TokenIllegal
	}

	return p.tokens[p.index-1].Token
}

// * parser_reader.go ends here.
