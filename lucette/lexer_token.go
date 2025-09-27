// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// lexer_token.go --- Lexer token methods.
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

// Add a new lexed token with the given type and lexeme to the token list.
func (l *lexer) addToken(token Token, lexeme string) {
	l.addTokenWithLiteral(token, lexeme, "")
}

// Add a new lexed token with the given token type, lexeme, and literal
// value to the token list.
func (l *lexer) addTokenWithLiteral(token Token, lexeme string, literal any) {
	l.tokens = append(
		l.tokens,
		NewLexedTokenWithLiteral(token, lexeme, literal, l.startPos, l.currPos))
}

// If the next rune in the lexer matches `match`, then a new lexed token
// is created using `tokThen` as the token and `lexThen` as the lexeme.
//
// If the match fails, then a new lexed token is created using `tokElse` as
// the token and `lexElse` as the lexeme instead.
func (l *lexer) addTokenIf(match rune, tokThen Token, lexThen string, tokElse Token, lexElse string) error {
	found, err := l.matchRune(match)
	if err != nil {
		return err
	}

	if found {
		l.addToken(tokThen, lexThen)

		return nil
	}

	l.addToken(tokElse, lexElse)

	return nil
}

// If the next rune in the lexer matches `match`, then a new lexed token
// using `token` and `lexeme` is created.
func (l *lexer) addTokenWhen(match rune, token Token, lexeme string) error {
	found, err := l.matchRune(match)
	if err != nil {
		return err
	}

	if found {
		l.addToken(token, lexeme)
	}

	return nil
}

// * lexer_token.go ends here.
