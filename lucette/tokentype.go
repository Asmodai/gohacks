// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// tokentype.go --- Token types.
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

// * Constants:

const (
	ttEOF TokenType = iota // End of file.

	ttNumber // Numeric value.
	ttPhrase // String value.
	ttField  // Field name.
	ttRegexp // Regular expression.

	ttPlus     // '+'
	ttMinus    // '-'
	ttStar     // '*'
	ttQuestion // '?'
	ttLParen   // '('
	ttLBracket // '['
	ttLCurly   // '{'
	ttRParen   // ')'
	ttRBracket // ']'
	ttRCurly   // '}'
	ttColon    // ':'
	ttTilde    // '~'
	ttCaret    // '^'

	ttTo  // 'TO' lexeme.
	ttAnd // 'AND'/'&&' lexeme.
	ttOr  // 'OR'/'||' lexeme.
	ttNot // 'NOT'/'!' lexeme.
	ttLT  // '<' lexeme.
	ttLTE // '<=' lexeme.
	ttGT  // '>' lexeme.
	ttGTE // '>=' lexeme.

	ttIllegal // Illegal token.
	ttUnknown // Unknown token.
)

// ** Variables:

var (
	//nolint:gochecknoglobals
	tokenTypeNames = map[TokenType]string{
		ttEOF:      "ttEOF",
		ttNumber:   "ttNumber",
		ttPhrase:   "ttPhrase",
		ttField:    "ttField",
		ttRegexp:   "ttRegexp",
		ttPlus:     "ttPlus",
		ttMinus:    "ttMinus",
		ttStar:     "ttStar",
		ttQuestion: "ttQuestion",
		ttLParen:   "ttLParen",
		ttLBracket: "ttLBracket",
		ttLCurly:   "ttLCurly",
		ttRParen:   "ttRParen",
		ttRBracket: "ttRBracket",
		ttRCurly:   "ttRCurly",
		ttColon:    "ttColon",
		ttTilde:    "ttTilde",
		ttCaret:    "ttCaret",
		ttTo:       "ttTo",
		ttAnd:      "ttAnd",
		ttOr:       "ttOr",
		ttNot:      "ttNot",
		ttLT:       "ttLT",
		ttLTE:      "ttLTE",
		ttGT:       "ttGT",
		ttGTE:      "ttGTE",
		ttIllegal:  "ttIllegal",
		ttUnknown:  "ttUnknown",
	}

	//nolint:gochecknoglobals
	tokenTypeKeyword = map[string]TokenType{
		"NOT": ttNot,
		"OR":  ttOr,
		"AND": ttAnd,
		"TO":  ttTo,
	}

	//nolint:gochecknoglobals
	tokenTypeLiterals = map[TokenType]string{
		ttEOF:      "<EOF>",
		ttPlus:     "+",
		ttMinus:    "-",
		ttStar:     "*",
		ttQuestion: "?",
		ttLParen:   "(",
		ttLBracket: "[",
		ttLCurly:   "{",
		ttRParen:   ")",
		ttRBracket: "]",
		ttRCurly:   "}",
		ttColon:    ":",
		ttTilde:    "~",
		ttCaret:    "^",
		ttIllegal:  "<illegal>",
		ttUnknown:  "<unknown>",
	}

	//nolint:gochecknoglobals
	tokenPunct = map[rune]TokenType{
		'+': ttPlus,
		'-': ttMinus,
		'*': ttStar,
		'?': ttQuestion,
		'(': ttLParen,
		'{': ttLCurly,
		'[': ttLBracket,
		')': ttRParen,
		'}': ttRCurly,
		']': ttRBracket,
		':': ttColon,
		'~': ttTilde,
		'^': ttCaret,
		'!': ttNot,
	}
)

// * Code:

// ** Type:

// Lexer token.
type TokenType uint

// ** Methods:

// Stringer method for lexer tokens.
func (t TokenType) String() string {
	if t >= ttUnknown {
		return tokenTypeNames[ttUnknown]
	}

	return tokenTypeNames[t]
}

func (t TokenType) Literal() string {
	if t < ttPlus || t >= ttCaret {
		// Special case for the `!` operator.
		if t == ttNot {
			return "!"
		}

		return tokenTypeLiterals[ttIllegal]
	}

	return tokenTypeLiterals[t]
}

// * tokentype.go ends here.
