// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// token_test.go --- Token tests.
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
	"testing"

	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Tests:

func TestToken(t *testing.T) {
	tokType := TokenAnd
	tokLexeme := "AND"
	tokLiteral := "and"
	tokErr := errors.Base("Chungus")
	tokStart := NewPosition(1, 2)
	tokEnd := NewPosition(3, 4)

	t.Run("NewLexedToken", func(t *testing.T) {
		tok := NewLexedToken(tokType, tokLexeme, tokStart, tokEnd)

		if tok.Token != tokType {
			t.Errorf("Token type: %v != %v", tok.Token, tokType)
		}

		if tok.Lexeme != tokLexeme {
			t.Errorf("Lexeme: %v != %v", tok.Lexeme, tokLexeme)
		}

		if tok.Start != tokStart {
			t.Errorf("Start: %#v != %#v", tok.Start, tokStart)
		}

		if tok.End != tokEnd {
			t.Errorf("End: %#v != %#v", tok.End, tokEnd)
		}
	})

	t.Run("NewLexedTokenWithLiteral", func(t *testing.T) {
		tok := NewLexedTokenWithLiteral(
			tokType,
			tokLexeme,
			tokLiteral,
			tokStart,
			tokEnd)

		if tok.Token != tokType {
			t.Errorf("Token type: %v != %v", tok.Token, tokType)
		}

		if tok.Lexeme != tokLexeme {
			t.Errorf("Lexeme: %v != %v", tok.Lexeme, tokLexeme)
		}

		if tok.Literal.Value != tokLiteral {
			t.Errorf("Literal: %q != %q", tok.Literal, tokLiteral)
		}

		if tok.Start != tokStart {
			t.Errorf("Start: %#v != %#v", tok.Start, tokStart)
		}

		if tok.End != tokEnd {
			t.Errorf("End: %#v != %#v", tok.End, tokEnd)
		}
	})

	t.Run("NewLexedTokenWithError", func(t *testing.T) {
		tok := NewLexedTokenWithError(
			tokType,
			tokLexeme,
			tokErr,
			tokStart,
			tokEnd)

		if tok.Token != tokType {
			t.Errorf("Token type: %v != %v", tok.Token, tokType)
		}

		if tok.Lexeme != tokLexeme {
			t.Errorf("Lexeme: %v != %v", tok.Lexeme, tokLexeme)
		}

		if tok.Literal.Err != tokErr {
			t.Errorf("Literal: %q != %q", tok.Literal, tokLiteral)
		}

		if tok.Start != tokStart {
			t.Errorf("Start: %#v != %#v", tok.Start, tokStart)
		}

		if tok.End != tokEnd {
			t.Errorf("End: %#v != %#v", tok.End, tokEnd)
		}
	})
}

// * token_test.go ends here.
