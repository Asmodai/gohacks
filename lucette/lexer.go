// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// lexer.go --- The lexer.
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
	"bufio"
	"io"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

// * Variables:

var (
	// Map of `rune -> rune` for conversion of escapes within strings.
	//nolint:gochecknoglobals
	escapePhrase = map[rune]rune{
		'"':  '"',
		'\\': '\\',
		'/':  '/',
		'b':  '\b',
		'f':  '\f',
		't':  '\t',
	}

	// Map of `rune -> rune` for conversion of escapes within quoted
	// strings.
	//
	//nolint:gochecknoglobals
	escapeQuoted = map[rune]rune{
		'\'': '\'',
		'\\': '\\',
		'/':  '/',
		'b':  '\b',
		'f':  '\f',
		't':  '\t',
	}
)

// * Code:

// ** Type:

// The lexer.
type lexer struct {
	reader   *bufio.Reader // IO reader.
	tokens   []LexedToken  // Token list.
	currPos  Position      // Current source position.
	startPos Position      // Starting position.
	lastPos  Position      // Last read position.
	lastRead bool          // Can unread?
}

// ** Methods:

// Reset the lexer's internal state.
func (l *lexer) Reset() {
	l.currPos = NewPosition(1, 0)
	l.startPos = l.currPos
	l.lastPos = NewEmptyPosition()
	l.lastRead = false
	l.tokens = make([]LexedToken, 0, initialTokenCapacity)
	l.reader = nil
}

// Return the list of lexed tokens.
func (l *lexer) Tokens() []LexedToken {
	return l.tokens
}

// Invoke the lexer on the specified reader.
func (l *lexer) Lex(reader io.Reader) ([]LexedToken, error) {
	l.Reset()

	if reader == nil {
		return []LexedToken{}, errors.WithStack(ErrInvalidReader)
	}

	l.reader = bufio.NewReader(reader)

	for {
		l.startPos = l.currPos

		err := l.lexToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return []LexedToken{}, errors.WithStack(err)
		}
	}

	return l.tokens, nil
}

// ** Functions:

// Create a new lexer.
func NewLexer() Lexer {
	return &lexer{
		currPos: NewPosition(1, 0),
		tokens:  make([]LexedToken, 0, initialTokenCapacity)}
}

// * lexer.go ends here.
