// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// token.go --- Token type.
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

import (
	"strings"

	"github.com/Asmodai/gohacks/utils"
)

// * Code:

// ** Type:

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   literal
	start     Position
	end       Position
}

// ** Methods:

func (t *Token) String() string {
	const (
		padType   = 10
		padLexeme = 4
		padPos    = 8
	)

	var sbld strings.Builder

	sbld.WriteString("Token -- start:")
	sbld.WriteString(utils.Pad(t.start.String(), padPos))

	sbld.WriteString(" end:")
	sbld.WriteString(utils.Pad(t.end.String(), padPos))

	sbld.WriteString(" start:")
	sbld.WriteString(utils.Pad(t.tokenType.String(), padType))

	sbld.WriteString(" lexeme:")
	sbld.WriteString(utils.Pad(t.lexeme, padLexeme))

	sbld.WriteString(" literal:\"")
	sbld.WriteString(t.literal.String())
	sbld.WriteRune('"')

	return sbld.String()
}

// ** Functions:

func newToken(tokType TokenType, lexeme string, start, end Position) Token {
	return newTokenWithLiteral(tokType, lexeme, "", start, end)
}

func newTokenWithLiteral(tokType TokenType, lexeme string, lit any, start, end Position) Token {
	return Token{
		tokenType: tokType,
		lexeme:    lexeme,
		literal:   literal{value: lit},
		start:     start,
		end:       end,
	}
}

func newTokenWithError(tokType TokenType, lexeme string, err error, start, end Position) Token {
	return Token{
		tokenType: tokType,
		lexeme:    lexeme,
		literal:   literal{err: err},
		start:     start,
		end:       end,
	}
}

// * token.go ends here.
