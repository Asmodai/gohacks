// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// lexer_reader.go --- Lexer reader functions.
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
	"strings"

	"github.com/pkg/errors"
)

// * Code:

// ** Methods:

// Read a rune and advance the reader.
func (l *lexer) readRune() (rune, error) {
	l.lastPos = l.currPos

	read, _, err := l.reader.ReadRune()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	l.lastRead = true

	switch {
	case read == '\n':
		l.currPos.Line++
		l.currPos.Column = 0

	default:
		l.currPos.Column++
	}

	return read, nil
}

// Unread the previously read rune back to the buffer.
//
// If an attempt is made to unread without having a read, then
// `ErrDoubleUnread` is returned.
//
// This error could potentially be misleading.  It will trigger if an
// unread is attempted without first having done a read, or if an unread is
// attempted immediately after another unread.
func (l *lexer) unreadRune() error {
	if !l.lastRead {
		return errors.WithMessagef(
			ErrDoubleUnread,
			"%s",
			l.currPos.String())
	}

	if err := l.reader.UnreadRune(); err != nil {
		return errors.WithStack(err)
	}

	l.currPos = l.lastPos
	l.lastRead = false

	return nil
}

// Peek the next rune without advancing the reader.
//
// This is a compromise.  `bufio` doesn't provide a way of peeking the next
// rune, only bytes.  We do not want to jump through hoops here.  And even
// if we did, `bufio`'s `Peek` will prevent `UnreadRune` until a subsequent
// `ReadRune` is performed.
//
// For this implementation we read the next rune and then immediately unread
// it to simulate a `PeekRune`.
func (l *lexer) peekRune() (rune, error) {
	next, err := l.readRune()
	if err != nil {
		return 0, err
	}

	// Attempt to unread.
	if err := l.unreadRune(); err != nil {
		return 0, err
	}

	return next, nil
}

// Check if the next rune is equal to `match`.
//
// Returns `false` if there is no match.
//
// Returns `false` and an error if there is a reader error.
func (l *lexer) matchRune(match rune) (bool, error) {
	next, err := l.readRune()
	if err != nil {
		return false, err
	}

	if next != match {
		// No match, so unread the rune.
		return false, l.unreadRune()
	}

	return true, nil
}

// Consume runes from the buffer for as long as `cond` evaluates to `true`.
//
// Should `cond` evaluate to `false`, then the last-read rune is placed
// back into the buffer.
func (l *lexer) readWhile(cond func(rune) bool) (string, error) {
	var sbld strings.Builder

	for {
		read, err := l.readRune()
		if err != nil {
			return sbld.String(), err
		}

		// Call user's condition func.
		if !cond(read) {
			if err := l.unreadRune(); err != nil {
				return sbld.String(), err
			}

			break
		}

		sbld.WriteRune(read)
	}

	return sbld.String(), nil
}

// * lexer_reader.go ends here.
