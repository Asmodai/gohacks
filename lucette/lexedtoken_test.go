// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// token_test.go --- Token tests.
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
	"testing"

	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Tests:

func TestToken(t *testing.T) {
	tokType := ttAnd
	tokLexeme := "AND"
	tokLiteral := "and"
	tokErr := errors.Base("Chungus")
	tokStart := newPosition(1, 2)
	tokEnd := newPosition(3, 4)

	t.Run("newToken", func(t *testing.T) {
		tok := newToken(tokType, tokLexeme, tokStart, tokEnd)

		if tok.tokenType != tokType {
			t.Errorf("Token type: %v != %v", tok.tokenType, tokType)
		}

		if tok.lexeme != tokLexeme {
			t.Errorf("Lexeme: %v != %v", tok.lexeme, tokLexeme)
		}

		if tok.start != tokStart {
			t.Errorf("Start: %#v != %#v", tok.start, tokStart)
		}

		if tok.end != tokEnd {
			t.Errorf("End: %#v != %#v", tok.end, tokEnd)
		}
	})

	t.Run("newTokenWithLiteral", func(t *testing.T) {
		tok := newTokenWithLiteral(
			tokType,
			tokLexeme,
			tokLiteral,
			tokStart,
			tokEnd)

		if tok.tokenType != tokType {
			t.Errorf("Token type: %v != %v", tok.tokenType, tokType)
		}

		if tok.lexeme != tokLexeme {
			t.Errorf("Lexeme: %v != %v", tok.lexeme, tokLexeme)
		}

		if tok.literal.value != tokLiteral {
			t.Errorf("Literal: %q != %q", tok.literal, tokLiteral)
		}

		if tok.start != tokStart {
			t.Errorf("Start: %#v != %#v", tok.start, tokStart)
		}

		if tok.end != tokEnd {
			t.Errorf("End: %#v != %#v", tok.end, tokEnd)
		}
	})

	t.Run("newTokenWithError", func(t *testing.T) {
		tok := newTokenWithError(
			tokType,
			tokLexeme,
			tokErr,
			tokStart,
			tokEnd)

		if tok.tokenType != tokType {
			t.Errorf("Token type: %v != %v", tok.tokenType, tokType)
		}

		if tok.lexeme != tokLexeme {
			t.Errorf("Lexeme: %v != %v", tok.lexeme, tokLexeme)
		}

		if tok.literal.err != tokErr {
			t.Errorf("Literal: %q != %q", tok.literal, tokLiteral)
		}

		if tok.start != tokStart {
			t.Errorf("Start: %#v != %#v", tok.start, tokStart)
		}

		if tok.end != tokEnd {
			t.Errorf("End: %#v != %#v", tok.end, tokEnd)
		}
	})
}

// * token_test.go ends here.
