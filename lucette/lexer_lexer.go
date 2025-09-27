// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// lexer_lexer.go --- Lexer methods.
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
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/Asmodai/gohacks/stringy"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

// * Variables:

// * Code:

// Add a lexed regular expression token.
//
// The lexeme is checked for presence of POSIX/PCRE2 regex flags.  If any
// are found then the lexeme is modified to include the relevant flag for
// the Go regex implementation.
func (l *lexer) addRegex(lexeme string) error {
	var flags strings.Builder

	for {
		peek, err := l.peekRune()
		if err != nil {
			break
		}

		if !unicode.IsLetter(peek) {
			break
		}

		if _, err := l.readRune(); err != nil {
			return err
		}

		flags.WriteRune(peek)
	}

	fstr := flags.String()
	if len(fstr) > 0 {
		for _, flag := range fstr {
			switch flag {
			case 'i', 'm', 's', 'U':
			// Flag is handled, do nothing.

			default:
				// Flag is not something we handle.
				return errors.WithMessagef(
					ErrRegexFlags,
					"%q at %s",
					flag,
					l.currPos.String())
			}
		}

		lexeme = "(?" + fstr + ")" + lexeme
	}

	l.addTokenWithLiteral(TokenRegex, "", lexeme)

	return nil
}

// Scan for a regular expression.
func (l *lexer) lexRegex() error {
	var sbld strings.Builder

	for {
		next, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.WithMessagef(
					ErrUnterminatedRegex,
					l.startPos.String())
			}

			return err
		}

		switch next {
		case '\r', '\n':
			return errors.WithMessagef(ErrNewlineInRegex, l.currPos.String())

		case '\\':
			next, err := l.readRune()
			if err != nil {
				return err
			}

			sbld.WriteRune('\\')
			sbld.WriteRune(next)

		case '/':
			return l.addRegex(sbld.String())

		default:
			sbld.WriteRune(next)
		}
	}
}

// Parse and add a numeric value lexed token.
func (l *lexer) addNumber(input string) error {
	strval := strings.ReplaceAll(input, "_", "")

	fval, err := strconv.ParseFloat(strval, 64)
	if err != nil {
		return errors.WithMessagef(
			err,
			"bad number %q at %s",
			input,
			l.currPos.String())
	}

	l.addTokenWithLiteral(TokenNumber, "", fval)

	return nil
}

// Scan for a number.
//
//nolint:cyclop
func (l *lexer) lexNumber() error {
	if err := l.unreadRune(); err != nil {
		return err
	}

	var (
		sbld    strings.Builder
		sawDot  bool
		sawExp  bool
		atStart = true
		prev    = rune(0)
	)

	for {
		accept := false

		next, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		switch {
		// `[0-9]`?
		case unicode.IsDigit(next):
			accept = true

			// `[0-9]_`?
		case next == '_' && prev != 0 && unicode.IsDigit(prev):
			accept = true

			// Ensure we just have one decimal point.
			// Also must come before the exponent.
		case next == '.' && !sawDot && !sawExp:
			sawDot, accept = true, true

			// Ensure we just have one exponent.
		case (next == 'e' || next == 'E') && !sawExp:
			sawExp, accept = true, true

			// Check for a +/- at the start of the number and
			// in the exponent.
		case (next == '+' || next == '-') && (atStart || prev == 'e' || prev == 'E'):
			accept = true
		}

		if !accept {
			if err := l.unreadRune(); err != nil {
				return err
			}

			break
		}

		sbld.WriteRune(next)
		prev = next
		atStart = false
	}

	return l.addNumber(sbld.String())
}

// Scan a Unicode escape (\u0000) to a rune.
func (l *lexer) lexUnicodeEscape() (rune, error) {
	var hexRunes [unicodeEscapeSize]rune

	for idx := range unicodeEscapeSize {
		hex, err := l.readRune()
		if err != nil {
			return 0, err
		}

		if !stringy.IsHexadecimal(hex) {
			return 0, errors.WithMessagef(
				ErrUnexpectedRune,
				"invalid \\u escape at %s",
				l.currPos.String())
		}

		hexRunes[idx] = hex
	}

	val, _ := strconv.ParseInt(string(hexRunes[:]), 16, 32)

	return rune(val), nil
}

// Scan an escape sequence.
func (l *lexer) lexEscape(table map[rune]rune) (rune, error) {
	next, err := l.readRune()
	if err != nil {
		return 0, err
	}

	if replaced, found := table[next]; found {
		return replaced, nil
	}

	if next == 'u' {
		return l.lexUnicodeEscape()
	}

	return 0, errors.WithMessagef(
		ErrUnexpectedRune,
		"bad escape '\\%c' at %s",
		next,
		l.currPos.String())
}

// Scan for a quoted field.
func (l *lexer) lexQuotedField() error {
	var sbld strings.Builder

	for {
		next, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.WithMessagef(
					ErrUnterminatedField,
					l.currPos.String())
			}

			return err
		}

		switch next {
		case '\r', '\n':
			return errors.WithMessagef(
				ErrNewlineInField,
				l.currPos.String())

		case '\\':
			esc, err := l.lexEscape(escapeQuoted)
			if err != nil {
				return err
			}

			sbld.WriteRune(esc)

		case '\'':
			l.addTokenWithLiteral(TokenField, "", sbld.String())

			return nil

		default:
			sbld.WriteRune(next)
		}
	}
}

// Scan for an identifier.
func (l *lexer) lexIdentifier() (string, error) {
	res, err := l.readWhile(func(read rune) bool {
		if unicode.IsLetter(read) || unicode.IsDigit(read) {
			return true
		}

		switch read {
		case '.', '-', '+', '_':
			return true

		default:
			return false
		}
	})

	return res, err
}

// Scan for a field.
func (l *lexer) lexField() error {
	if err := l.unreadRune(); err != nil {
		return err
	}

	ident, err := l.lexIdentifier()
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	if next, err := l.peekRune(); err != nil || next != ':' {
		if kwd, ok := tokenKeywords[strings.ToUpper(ident)]; ok {
			l.addToken(kwd, ident)

			return nil
		}

		return errors.WithMessagef(
			ErrUnexpectedBareword,
			"%q at %s",
			ident,
			l.currPos.String())
	}

	l.addTokenWithLiteral(TokenField, "", ident)

	return nil
}

// Scan for a phrase.
func (l *lexer) lexPhrase() error {
	var sbld strings.Builder

	for {
		next, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.WithMessagef(
					ErrUnterminatedString,
					l.startPos.String())
			}

			return err
		}

		switch next {
		case '\\':
			escaped, err := l.lexEscape(escapePhrase)
			if err != nil {
				return err
			}

			next = escaped

		case '\r', '\n':
			return errors.WithMessagef(
				ErrNewlineInPhrase,
				"%v at %s",
				next,
				l.currPos.String())

		case '"':
			l.addTokenWithLiteral(
				TokenPhrase,
				"",
				sbld.String())

			return nil
		}

		sbld.WriteRune(next)
	}
}

// Scan for a literal value.
func (l *lexer) lexLiteral() error {
	if err := l.unreadRune(); err != nil {
		return err
	}

	next, err := l.readRune()
	if err != nil {
		return err
	}

	switch {
	case unicode.IsDigit(next):
		return l.lexNumber()

	case unicode.IsLetter(next):
		return l.lexField()

	default:
		return errors.WithMessagef(
			ErrUnexpectedToken,
			"'%c' [%d] at %s",
			next,
			next,
			l.currPos.String())
	}
}

// Scan for a token.
//
//nolint:cyclop
func (l *lexer) lexToken() error {
	next, err := l.readRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			l.addToken(TokenEOF, "")
		}

		return err
	}

	// Whitespace? Don't continue.
	if unicode.IsSpace(next) {
		return nil
	}

	// Simple punctuation? Create a token and return.
	if simple, found := tokenPunct[next]; found {
		l.addToken(simple, simple.Literal())

		return nil
	}

	switch next {
	case '<':
		return l.addTokenIf('=', TokenLTE, "<=", TokenLT, "<")

	case '>':
		return l.addTokenIf('=', TokenGTE, ">=", TokenGT, ">")

	case '|':
		return l.addTokenWhen('|', TokenOr, "||")

	case '&':
		return l.addTokenWhen('&', TokenAnd, "&&")

	case '"':
		return l.lexPhrase()

	case '/':
		return l.lexRegex()

	case '\'':
		return l.lexQuotedField()

	default:
		return l.lexLiteral()
	}
}

// * lexer_lexer.go ends here.
