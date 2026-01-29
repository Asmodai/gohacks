// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// errors.go --- Errors.
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

import "gitlab.com/tozd/go/errors"

// * Code:

var (
	// Returned when the typer detects an invalid datetime.
	ErrBadDateTime = errors.Base("bad datetime value")

	// Returned should an attempt be made to `unread` after a rune has
	// already been put back into the reader.
	ErrDoubleUnread = errors.Base("double unread")

	// Returned if the lexer is invoked without a valid reader.
	ErrInvalidReader = errors.Base("reader is not valid")

	// Returned when the code generator detects a label without a target.
	ErrJumpMissingArg = errors.Base("jump missing target arg")

	// Returned when the code generator detects a label with an invalid
	// target.
	ErrJumpNotLabelID = errors.Base("jump target arg not LabelID")

	// Returned when the code generator detects a label with a bad ID.
	ErrLabelBadIDType = errors.Base("LABEL has bad id type")

	// Returned when the code generator detects a label that lacks an ID.
	ErrLabelMissingID = errors.Base("LABEL missing id")

	// Returned when the lexer detects an embedded newline in a field
	// name.
	ErrNewlineInField = errors.Base("embedded newline in field")

	// Returned when the lexer detects an embedded newline in a phrase.
	ErrNewlineInPhrase = errors.Base("embedded newline in phrase")

	// Returned when the lexer detects a newline in a regular expression.
	ErrNewlineInRegex = errors.Base("embedded newline in regular expression")

	// Returned if no tokens were provided.
	ErrNoTokens = errors.Base("no tokens")

	// Returned when the lexer detects unsupported flags in a regular
	// expression.
	ErrRegexFlags = errors.Base("regex flags not supported")

	// Returned when the code generator detects a label that has not been
	// bound to a target.
	ErrUnboundLabel = errors.Base("unbound label")

	// Returned when the lexer detects an unexpected bareword in the
	// source code.
	ErrUnexpectedBareword = errors.Base("unexpected bareword (missing quotes or field?)")

	// Returned when the lexer detects an unexpected character.
	ErrUnexpectedRune = errors.Base("unexpected rune")

	// Returned when the lexer detects an unexpected token in the source
	// code.
	ErrUnexpectedToken = errors.Base("unexpected token")

	// Returned when the typer detects an unknown literal.
	ErrUnknownLiteral = errors.Base("unknown literal")

	// Returned when the lexer detects that a quoted field name is
	// unterminated.
	ErrUnterminatedField = errors.Base("unterminated quoted field")

	// Returned when the lexer detects an unterminated regular expression.
	ErrUnterminatedRegex = errors.Base("unterminated regular expression")

	// Returned when the lexer detects an unterminated quoted string.
	ErrUnterminatedString = errors.Base("unterminated string")
)

// * errors.go ends here.
