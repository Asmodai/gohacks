// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// tokentype.go --- Token types.
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

// * Constants:

const (
	TokenEOF Token = iota // End of file.

	TokenNumber // Numeric value.
	TokenPhrase // String phrase value.
	TokenField  // Field name.
	TokenRegex  // Regular expression.

	TokenPlus     // '+'
	TokenMinus    // '-'
	TokenStar     // '*'
	TokenQuestion // '?'
	TokenLParen   // '('
	TokenLBracket // '['
	TokenLCurly   // '{'
	TokenRParen   // ')'
	TokenRBracket // ']'
	TokenRCurly   // '}'
	TokenColon    // ':'
	TokenTilde    // '~'
	TokenCaret    // '^'

	TokenTo  // 'TO'.
	TokenAnd // 'AND'/'&&'.
	TokenOr  // 'OR'/'||'.
	TokenNot // 'NOT'/'!'.
	TokenLT  // '<'.
	TokenLTE // '<='.
	TokenGT  // '>'
	TokenGTE // '>='

	TokenIllegal // Illegal token.
	TokenUnknown // Unknown token.
)

// * Variables:

var (
	// Map of `tokens -> strings` for use with pretty-printing.
	//
	//nolint:gochecknoglobals
	tokenStrings = map[Token]string{
		TokenEOF:      "TokenEOF",
		TokenNumber:   "TokenNumber",
		TokenPhrase:   "TokenPhrase",
		TokenField:    "TokenField",
		TokenRegex:    "TokenRegex",
		TokenPlus:     "TokenPlus",
		TokenMinus:    "TokenMinus",
		TokenStar:     "TokenStar",
		TokenQuestion: "TokenQuestion",
		TokenLParen:   "TokenLParen",
		TokenLBracket: "TokenLBracket",
		TokenLCurly:   "TokenLCurly",
		TokenRParen:   "TokenRParen",
		TokenRBracket: "TokenRBracket",
		TokenRCurly:   "TokenRCurly",
		TokenColon:    "TokenColon",
		TokenTilde:    "TokenTilde",
		TokenCaret:    "TokenCaret",
		TokenTo:       "TokenTo",
		TokenAnd:      "TokenAnd",
		TokenOr:       "TokenOr",
		TokenNot:      "TokenNot",
		TokenLT:       "TokenLT",
		TokenLTE:      "TokenLTE",
		TokenGT:       "TokenGT",
		TokenGTE:      "TokenGTE",
		TokenIllegal:  "TokenIllegal",
		TokenUnknown:  "TokenUnknown",
	}

	// Map of `keyword strings -> tokens` for use with identifier parsing.
	//
	//nolint:gochecknoglobals
	tokenKeywords = map[string]Token{
		"NOT": TokenNot,
		"OR":  TokenOr,
		"AND": TokenAnd,
		"TO":  TokenTo,
	}

	// Map of `tokens -> literal strings` for use with building lexed
	// tokens.
	//
	//nolint:gochecknoglobals
	tokenLiterals = map[Token]string{
		TokenEOF:      "EOF",
		TokenPlus:     "+",
		TokenMinus:    "-",
		TokenStar:     "*",
		TokenQuestion: "?",
		TokenLParen:   "(",
		TokenLBracket: "[",
		TokenLCurly:   "{",
		TokenRParen:   ")",
		TokenRBracket: "]",
		TokenRCurly:   "}",
		TokenColon:    ":",
		TokenTilde:    "~",
		TokenCaret:    "^",
		TokenNot:      "!",
		TokenIllegal:  "Illegal",
		TokenUnknown:  "Unknown",
	}

	// Map of `runes -> tokens` for use with parsing of runes to tokens.
	//
	//nolint:gochecknoglobals
	tokenPunct = map[rune]Token{
		'+': TokenPlus,
		'-': TokenMinus,
		'*': TokenStar,
		'?': TokenQuestion,
		'(': TokenLParen,
		'[': TokenLBracket,
		'{': TokenLCurly,
		')': TokenRParen,
		']': TokenRBracket,
		'}': TokenRCurly,
		':': TokenColon,
		'~': TokenTilde,
		'^': TokenCaret,
		'!': TokenNot,
	}
)

// * Code:

// ** Types:

// Lucette token type.
type Token int

// ** Methods:

// Return the string representation of a token.
func (t Token) String() string {
	if t >= TokenUnknown {
		return tokenStrings[TokenUnknown]
	}

	return tokenStrings[t]
}

// Return the literal string representation of a token if it has one.
func (t Token) Literal() string {
	if (t < TokenPlus || t > TokenCaret) && t != TokenNot {
		return tokenLiterals[TokenIllegal]
	}

	return tokenLiterals[t]
}

// * tokentype.go ends here.
