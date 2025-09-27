// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// lexedtoken.go --- Lexed token type.
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

	"github.com/Asmodai/gohacks/debug"
	"github.com/Asmodai/gohacks/utils"
)

// * Code:

// ** Type:

// A lexed token.
type LexedToken struct {
	Literal Literal  // Literal value for the token.
	Lexeme  string   // Lexeme for the token.
	Start   Position // Start position within source code.
	End     Position // End position within source code.
	Token   Token    // The token.
}

// ** Methods:

func (lt *LexedToken) String() string {
	const (
		padType   = 10
		padLexeme = 4
		padPos    = 8
	)

	var sbld strings.Builder

	sbld.WriteString("Token -- start:")
	sbld.WriteString(utils.Pad(lt.Start.String(), padPos))

	sbld.WriteString(" end:")
	sbld.WriteString(utils.Pad(lt.End.String(), padPos))

	sbld.WriteString(" start:")
	sbld.WriteString(utils.Pad(lt.Token.String(), padType))

	sbld.WriteString(" lexeme:")
	sbld.WriteString(utils.Pad(lt.Lexeme, padLexeme))

	sbld.WriteString(" literal:\"")
	sbld.WriteString(lt.Literal.String())
	sbld.WriteRune('"')

	return sbld.String()
}

// Display debugging information.
func (lt *LexedToken) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Lexed Token")

	dbg.Init(params...)

	dbg.Printf("Token:   %s", lt.Token.String())
	dbg.Printf("Lexeme:  %q", lt.Lexeme)
	dbg.Printf("Literal: %q", lt.Literal.String())
	dbg.Printf("Start:   %s", lt.Start.String())
	dbg.Printf("End:     %s", lt.End.String())

	dbg.End()
	dbg.Print()

	return dbg
}

// ** Functions:

// Return a new lexed token with the given lexeme.
func NewLexedToken(token Token, lexeme string, start, end Position) LexedToken {
	return NewLexedTokenWithLiteral(token, lexeme, "", start, end)
}

// Return a new lexed token with the given lexeme and literal.
func NewLexedTokenWithLiteral(token Token, lexeme string, lit any, start, end Position) LexedToken {
	return LexedToken{
		Token:   token,
		Lexeme:  lexeme,
		Literal: NewLiteral(lit),
		Start:   start,
		End:     end,
	}
}

// Return a new lexed token with the given lexeme and error message.
func NewLexedTokenWithError(token Token, lexeme string, err error, start, end Position) LexedToken {
	return LexedToken{
		Token:   token,
		Lexeme:  lexeme,
		Literal: NewErrorLiteral(err),
		Start:   start,
		End:     end,
	}
}

// * lexedtoken.go ends here.
