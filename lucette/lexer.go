// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// lexer.go --- Lexer.
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

//nolint:dupword
/*

   Query = Disj , { WS Disj } ;
   Disj  = Conj , { WS ( "OR" | "||" ) WS Conj } ;
   Conj  = ModClause , { WS ( "AND" | "&&" ) WS ModClause } ;

   ModClause = [ ( "+" | "-" | "!" | "NOT" ) WS? ] Clause ;
   Clause    = Field WS? Value [ WS? Boost ] ;

   Field       = ( BareField | QuotedField ) WS? ":" ;
   BareField   = IDENT ;
   QuotedField = "'" QUOTED_CONTENT "'" ;

   Value = Phrase [ Fuzzy ]
	 | Number
	 | Regex
	 | Range
	 ;

   Phrase = '"' STRING_CONTENT '"' ;
   Fuzzy  = "~" [ Number ] ;
   Regex  = "/" REGEX_CONTENT "/" ;

   Range      = Interval | Comparator ;
   Interval   = ( "[" | "{" ) WS? RangeBound WS "TO" WS RangeBound WS? ( "]" | "}" ) ;
   RangeBound = Phrase | Number | "*" ;

   Comparator = CompOp WS? RangeAtom ;
   CompOp     = ">" | "<" | ">=" | "<=" ;
   RangeAtom  = Phrase | Number ;

   Boost = "^" Number ;

   Lexical classes:

   IDENT      = IDENT_CHAR , { IDENT_CHAR } ;
   IDENT_CHAR = Letter | Digit | "_" | "." | "+" | "-" ;

   Number = INT [ FRACT ] [ EXP ]
   INT    = DIGIT { [ "_" ] DIGIT };
   FRACT  = "." DIGIT { [ "_" ] DIGIT } ;
   EXP    = ( "e" | "E" ) [ "+" | "-" ] DIGIT {DIGIT } ;

   STRING_CONTENT = { ESC | ~["\\] } ;
   QUOTED_CONTENT = { QESC | ~['\\] } ;
   REGEX_CONTENT  = { RESC | ~[/\\] } ;

   ESC  = "\\" ( "\"" | "\\" | "/" | "b" | "f" | "n" | "r" | "t"
		      | "u" HEX HEX HEX HEX )
		      ;

   QESC = "\\" ( "'"  | "\\" | "/" | "b" | "f" | "n" | "r" | "t"
		      | "u" HEX HEX HEX HEX )
		      ;

   RESC = "\\" ( "/"  | "\\" | "." | "*" | "+" | "?" | "|" | "(" | ")"
		      | "[" | "]" | "{" | "}" | "d" | "s" | "w" | "b" | "B"
		      | "u" HEX HEX HEX HEX )
		      ;

   HEX    = DIGIT | "A"…"F" | "a"…"f" ;
   DIGIT  = "0"…"9" ;
   Letter = UnicodeLetter ;
   Digit  = UnicodeDigit ;

   WS      = WS_CHAR { WS_CHAR } ;
   WS_CHAR = " " | "\t" | "\r" | "\n" | "\f" | U+3000 ;

*/

// * Package:

//nolint:unused
package lucette

import (
	"bufio"
	"io"
	"strconv"
	"strings"
	"unicode"

	"gitlab.com/tozd/go/errors"
)

// * Imports:

// * Constants:

const (
	unicodeSize          = 4
	initialTokenCapacity = 4
)

// * Variables:

var (
	ErrUnterminatedRegex  = errors.Base("unterminated regular expression")
	ErrNewlineInRegex     = errors.Base("embedded newline in regular expression")
	ErrRegexFlags         = errors.Base("regex flags not supported")
	ErrUnexpectedToken    = errors.Base("unexpected token")
	ErrNewlineInPhrase    = errors.Base("embedded newline in phrase")
	ErrUnterminatedField  = errors.Base("unterminated quoted field")
	ErrNewlineInField     = errors.Base("embedded newline in field")
	ErrUnterminatedString = errors.Base("unterminated string")
	ErrUnexpectedRune     = errors.Base("unexpected rune")
	ErrDoubleUnread       = errors.Base("double unread")
	ErrUnexpectedBareword = errors.Base("unexpected bareword (missing quotes or field?)")

	//nolint:gochecknoglobals
	escapePhrase = map[rune]rune{
		'"':  '"',
		'\\': '\\',
		'/':  '/',
		'b':  '\b',
		'f':  '\f',
		't':  '\t',
	}

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
//
// You are not meant to use this directly.
type Lexer struct {
	reader   *bufio.Reader // IO reader.
	pos      Position      // Current position in code.
	startPos Position      // Starting position in current tokenising.
	lastPos  Position      // Last read position.
	lastRead bool          // Can unread a character.
	tokens   []Token       // Token list.
}

// ** Methods:

// *** Reader methods:

// Read a rune and advance the reader.
func (l *Lexer) readRune() (rune, error) {
	l.lastPos = l.pos

	read, _, err := l.reader.ReadRune()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	l.lastRead = true

	switch {
	case read == '\n':
		l.pos.line++
		l.pos.column = 0

	default:
		l.pos.column++
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
func (l *Lexer) unreadRune() error {
	if !l.lastRead {
		return errors.WithMessagef(
			ErrDoubleUnread,
			"line %d, col %d",
			l.pos.line,
			l.pos.column)
	}

	if err := l.reader.UnreadRune(); err != nil {
		return errors.WithStack(err)
	}

	l.pos = l.lastPos
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
func (l *Lexer) peekRune() (rune, error) {
	// Read next rune.
	next, err := l.readRune()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	// Unread rune.
	if err := l.unreadRune(); err != nil {
		return 0, errors.WithStack(err)
	}

	return next, nil
}

// Check if the next rune is equal to `match`.
//
// Returns `false` if there is no match.
//
// Returns `false` and an error if there is a reader error.
func (l *Lexer) matchRune(match rune) (bool, error) {
	next, err := l.readRune()
	if err != nil {
		return false, errors.WithStack(err)
	}

	if next != match {
		// No match, so unread rune.
		return false, errors.WithStack(l.unreadRune())
	}

	return true, nil
}

// Consume runes from the buffer for as long as `cond` evaluates to `true`.
//
// Should `cond` evaluate to `false`, then the last-read rune is placed
// back into the buffer.
func (l *Lexer) readWhile(cond func(rune) bool) (string, error) {
	var sbld strings.Builder

	for {
		read, err := l.readRune()
		if err != nil {
			return sbld.String(), errors.WithStack(err)
		}

		if !cond(read) {
			if err := l.unreadRune(); err != nil {
				return sbld.String(), errors.WithStack(err)
			}

			break
		}

		sbld.WriteRune(read)
	}

	return sbld.String(), nil
}

// *** Token methods:

// Add the given token type and lexeme to the list of tokens.
func (l *Lexer) addToken(ttype TokenType, lexeme string) {
	l.addTokenWithLiteral(ttype, lexeme, "")
}

// Add the given token type, lexeme, and literal value to the list of tokens.
func (l *Lexer) addTokenWithLiteral(ttype TokenType, lexeme string, litValue any) {
	l.tokens = append(
		l.tokens,
		newTokenWithLiteral(
			ttype,
			lexeme,
			litValue,
			l.startPos,
			l.pos))
}

// Add the token `ttThen` with the lexeme `lexThen` if `match` is found.
//
// If `match` is not found, add the token `ttElse` with the lexeme `lexElse`.
func (l *Lexer) addIf(match rune, ttThen, ttElse TokenType, lexThen, lexElse string) error {
	found, err := l.matchRune(match)
	if err != nil {
		return errors.WithStack(err)
	}

	if found {
		l.addToken(ttThen, lexThen)

		return nil
	}

	l.addToken(ttElse, lexElse)

	return nil
}

// Add the token `ttThen` with the lexeme `lexThen` if `match` is found.
func (l *Lexer) addWhen(match rune, ttThen TokenType, lexThen string) error {
	found, err := l.matchRune(match)
	if err != nil {
		return errors.WithStack(err)
	}

	if found {
		l.addToken(ttThen, lexThen)
	}

	return nil
}

// *** Lexer methods:

// Possibly lex a regex with flags.  It not, a plain regex token is returned.
func (l *Lexer) lexRegexFlags(body string) error {
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
			return errors.WithStack(err)
		}

		flags.WriteRune(peek)
	}

	fstr := flags.String()
	if len(fstr) > 0 {
		for _, flag := range fstr {
			switch flag {
			case 'i', 'm', 's', 'U':

			default:
				return errors.WithMessagef(
					ErrRegexFlags,
					"unsupported flag %q at %s",
					flag,
					l.pos.String())
			}
		}

		body = "(?" + fstr + ")" + body
	}

	l.addTokenWithLiteral(ttRegexp, "", body)

	return nil
}

// Lex a regex.
func (l *Lexer) lexRegex() error {
	var sbld strings.Builder

	for {
		next, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.WithMessagef(
					ErrUnterminatedRegex,
					l.startPos.String())
			}

			return errors.WithStack(err)
		}

		switch next {
		case '\n', '\r':
			return errors.WithMessagef(
				ErrNewlineInRegex,
				l.pos.String())

		case '\\':
			next, err := l.readRune()
			if err != nil {
				return errors.WithStack(err)
			}

			sbld.WriteRune('\\')
			sbld.WriteRune(next)

		case '/':
			return l.lexRegexFlags(sbld.String())

		default:
			sbld.WriteRune(next)
		}
	}
}

func (l *Lexer) lexToFloat(input string) error {
	strval := strings.ReplaceAll(input, "_", "")

	fval, err := strconv.ParseFloat(strval, 64)
	if err != nil {
		return errors.WithMessagef(
			err,
			"bad number %q at %s",
			input,
			l.pos.String())
	}

	l.addTokenWithLiteral(ttNumber, "", fval)

	return nil
}

// Lex a number.
//
//nolint:cyclop
func (l *Lexer) lexNumber() error {
	if err := l.unreadRune(); err != nil {
		return errors.WithStack(err)
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
		if err != nil && !errors.Is(err, io.EOF) {
			return errors.WithStack(err)
		}

		switch {
		case unicode.IsDigit(next):
			accept = true

		case next == '_' && prev != 0 && unicode.IsDigit(prev):
			accept = true

		case next == '.' && !sawDot && !sawExp:
			sawDot, accept = true, true

		case (next == 'e' || next == 'E') && !sawExp:
			sawExp, accept = true, true

		case (next == '+' || next == '-') && (atStart || prev == 'e' || prev == 'E'):
			accept = true
		}

		if !accept {
			if err == nil {
				if err := l.unreadRune(); err != nil {
					return errors.WithStack(err)
				}
			}

			break
		}

		sbld.WriteRune(next)
		prev = next
		atStart = false

		if errors.Is(err, io.EOF) {
			break
		}
	}

	return l.lexToFloat(sbld.String())
}

func (l *Lexer) lexEscape(table map[rune]rune) (rune, error) {
	next, err := l.readRune()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if replaced, found := table[next]; found {
		return replaced, nil
	}

	if next == 'u' {
		var hexrune [unicodeSize]rune

		for idx := range unicodeSize {
			hex, err := l.readRune()
			if err != nil {
				return 0, errors.WithStack(err)
			}

			if !isHex(hex) {
				return 0, errors.WithMessagef(
					ErrUnexpectedRune,
					"invalid \\u escape at %s",
					l.pos.String())
			}

			hexrune[idx] = hex
		}

		val, _ := strconv.ParseInt(string(hexrune[:]), 16, 32)

		return rune(val), nil
	}

	return 0, errors.WithMessagef(
		ErrUnexpectedRune,
		"bad escape '\\%c' at %s",
		next,
		l.pos.String())
}

func (l *Lexer) lexQuotedField() error {
	var sbld strings.Builder

	for {
		next, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.WithMessagef(
					ErrUnterminatedField,
					l.startPos.String())
			}

			return errors.WithStack(err)
		}

		switch next {
		case '\n', '\r':
			return errors.WithMessagef(ErrNewlineInField, l.pos.String())

		case '\\':
			esc, err := l.lexEscape(escapeQuoted)
			if err != nil {
				return errors.WithStack(err)
			}

			sbld.WriteRune(esc)

		case '\'':
			l.addTokenWithLiteral(ttField, "", sbld.String())

			return nil

		default:
			sbld.WriteRune(next)
		}
	}
}

// Lex a Lucette field.
//
//nolint:cyclop
func (l *Lexer) lexField() error {
	if err := l.unreadRune(); err != nil {
		return errors.WithStack(err)
	}

	result, err := l.readWhile(func(read rune) bool {
		if unicode.IsLetter(read) || unicode.IsDigit(read) {
			return true
		}

		switch read {
		case '.':
			fallthrough
		case '-':
			fallthrough
		case '+':
			fallthrough
		case '_':
			return true
		default:
			return false
		}
	})

	if err != nil && !errors.Is(err, io.EOF) {
		return errors.WithStack(err)
	}

	if next, err := l.peekRune(); err != nil || next != ':' {
		if kwd, ok := tokenTypeKeyword[strings.ToUpper(result)]; ok {
			l.addToken(kwd, result)

			return nil
		}

		return errors.WithMessagef(
			ErrUnexpectedBareword,
			"%q at %s",
			result,
			&l.startPos)
	}

	l.addTokenWithLiteral(ttField, "", result)

	return nil
}

func (l *Lexer) lexPhrase() error {
	var sbld strings.Builder

	for {
		next, err := l.readRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return errors.WithMessagef(
					ErrUnterminatedString,
					l.startPos.String())
			}

			return errors.WithStack(err)
		}

		switch next {
		case '\\':
			escaped, err := l.lexEscape(escapePhrase)
			if err != nil {
				return errors.WithStack(err)
			}

			next = escaped

		case '\n', '\r':
			return errors.WithMessagef(
				ErrNewlineInPhrase,
				"%s: %v",
				l.pos.String(),
				next)

		case '"':
			l.addTokenWithLiteral(
				ttPhrase,
				"",
				sbld.String())

			return nil
		}

		sbld.WriteRune(next)
	}
}

//nolint:cyclop
func (l *Lexer) lexToken() error {
	next, err := l.readRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			l.addToken(ttEOF, "")
		}

		return errors.WithStack(err)
	}

	// Don't even bother if the rune is whitespace.
	if unicode.IsSpace(next) {
		return nil
	}

	// Save time and look up simple punctuation tokens in a map.
	if simple, found := tokenPunct[next]; found {
		l.addToken(simple, simple.Literal())

		return nil
	}

	// Now lookup the other tokens.
	switch next {
	case '<':
		return l.addIf('=', ttLTE, ttLT, "<=", "<")

	case '>':
		return l.addIf('=', ttGTE, ttGT, ">=", ">")

	case '|':
		return l.addWhen('|', ttOr, "||")

	case '&':
		return l.addWhen('&', ttAnd, "&&")

	case '"':
		return l.lexPhrase()

	case '/':
		return l.lexRegex()

	case '\'':
		return l.lexQuotedField()

	default:
		switch {
		case unicode.IsDigit(next):
			return l.lexNumber()
		case unicode.IsLetter(next):
			return l.lexField()
		default:
			return errors.WithMessagef(
				ErrUnexpectedToken,
				"'%c' [%d]",
				next, next)
		}
	}
}

func (l *Lexer) lex() ([]Token, error) {
	for {
		l.startPos = l.pos

		err := l.lexToken()

		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return l.tokens, nil
}

// ** Functions:

// XXX This could live elsewhere.
func isHex(thing rune) bool {
	return (('0' <= thing && thing <= '9') ||
		('a' <= thing && thing <= 'f') ||
		('A' <= thing && thing <= 'F'))
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		pos:    newPosition(1, 0),
		tokens: make([]Token, 0, initialTokenCapacity),
	}
}

// * lexer.go ends here.
