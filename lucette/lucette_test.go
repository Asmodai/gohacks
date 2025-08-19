// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// lucette_test.go --- Full parser test.
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
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

// * Variables:

var (
	tests = []struct {
		input  string
		tokens []Token
		err    error
	}{
		{
			input:  "field:string",
			tokens: []Token{},
			err:    ErrUnexpectedBareword,
		}, {
			input: `level:42 && message:"whee"`,
			tokens: []Token{
				makeLiteralToken(ttField, "", "level"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttNumber, "", float64(42)),
				makeLiteralToken(ttAnd, "&&", ""),
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "whee"),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `'the level':42 && 'the message':"whee"`,
			tokens: []Token{
				makeLiteralToken(ttField, "", "the level"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttNumber, "", float64(42)),
				makeLiteralToken(ttAnd, "&&", ""),
				makeLiteralToken(ttField, "", "the message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "whee"),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `message:"value"~5 && level:95`,
			tokens: []Token{
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "value"),
				makeLiteralToken(ttTilde, "~", ""),
				makeLiteralToken(ttNumber, "", float64(5)),
				makeLiteralToken(ttAnd, "&&", ""),
				makeLiteralToken(ttField, "", "level"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttNumber, "", float64(95)),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `percent:3_123_456 || percent:3.1e+02`,
			tokens: []Token{
				makeLiteralToken(ttField, "", "percent"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttNumber, "", float64(3_123_456)),
				makeLiteralToken(ttOr, "||", ""),
				makeLiteralToken(ttField, "", "percent"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttNumber, "", float64(3.1e+02)),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `(message:"yes" or message:"maybe") and !message:"no"`,
			tokens: []Token{
				makeLiteralToken(ttLParen, "(", ""),
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "yes"),
				makeLiteralToken(ttOr, "or", ""),
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "maybe"),
				makeLiteralToken(ttRParen, ")", ""),
				makeLiteralToken(ttAnd, "and", ""),
				makeLiteralToken(ttNot, "!", ""),
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "no"),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `(message:"ye*" or message:"n?") and (not message:"chuckles" and !+message:"ass")`,
			tokens: []Token{
				makeLiteralToken(ttLParen, "(", ""),
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "ye*"),
				makeLiteralToken(ttOr, "or", ""),
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "n?"),
				makeLiteralToken(ttRParen, ")", ""),
				makeLiteralToken(ttAnd, "and", ""),
				makeLiteralToken(ttLParen, "(", ""),
				makeLiteralToken(ttNot, "not", ""),
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "chuckles"),
				makeLiteralToken(ttAnd, "and", ""),
				makeLiteralToken(ttNot, "!", ""),
				makeLiteralToken(ttPlus, "+", ""),
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "ass"),
				makeLiteralToken(ttRParen, ")", ""),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `(level:[1 to *})`,
			tokens: []Token{
				makeLiteralToken(ttLParen, "(", ""),
				makeLiteralToken(ttField, "", "level"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttLBracket, "[", ""),
				makeLiteralToken(ttNumber, "", float64(1)),
				makeLiteralToken(ttTo, "to", ""),
				makeLiteralToken(ttStar, "*", ""),
				makeLiteralToken(ttRCurly, "}", ""),
				makeLiteralToken(ttRParen, ")", ""),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `(source:["192.168.1.1" to "192.168.1.254"])`,
			tokens: []Token{
				makeLiteralToken(ttLParen, "(", ""),
				makeLiteralToken(ttField, "", "source"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttLBracket, "[", ""),
				makeLiteralToken(ttPhrase, "", "192.168.1.1"),
				makeLiteralToken(ttTo, "to", ""),
				makeLiteralToken(ttPhrase, "", "192.168.1.254"),
				makeLiteralToken(ttRBracket, "]", ""),
				makeLiteralToken(ttRParen, ")", ""),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `!source:"8.8.8.8"`,
			tokens: []Token{
				makeLiteralToken(ttNot, "!", ""),
				makeLiteralToken(ttField, "", "source"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttPhrase, "", "8.8.8.8"),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `source:>="8.8.8.8"`,
			tokens: []Token{
				makeLiteralToken(ttField, "", "source"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttGTE, ">=", ""),
				makeLiteralToken(ttPhrase, "", "8.8.8.8"),
				makeLiteralToken(ttEOF, "", ""),
			},
		}, {
			input: `message:/testing/i`,
			tokens: []Token{
				makeLiteralToken(ttField, "", "message"),
				makeLiteralToken(ttColon, ":", ""),
				makeLiteralToken(ttRegexp, "", "(?i)testing"),
				makeLiteralToken(ttEOF, "", ""),
			},
		},
	}
)

// * Code:

// ** Test structure:

type TestStruct struct {
	Timestamp time.Time
	Facility  string
	Level     int
	Percent   float64
	Source    string
	Message   string
}

func MakeStructSchema() Schema {
	ret := Schema{}

	ret["timestamp"] = FieldSpec{
		Name:    "Timestamp",
		FType:   FTDateTime,
		Layouts: []string{time.RFC3339},
	}

	ret["facility"] = FieldSpec{Name: "Facility", FType: FTText}
	ret["level"] = FieldSpec{Name: "Level", FType: FTNumeric}
	ret["percent"] = FieldSpec{Name: "Percent", FType: FTNumeric}
	ret["source"] = FieldSpec{Name: "Source", FType: FTIP}
	ret["message"] = FieldSpec{Name: "Message", FType: FTText}

	// These exist for lexer tests.
	ret["the level"] = FieldSpec{Name: "Level", FType: FTNumeric}
	ret["the message"] = FieldSpec{Name: "Message", FType: FTText}

	return ret
}

// ** Utilities:

func dumpTokens(toks []Token) string {
	var sbld strings.Builder

	length := len(toks)

	for idx, tok := range toks {
		sbld.WriteString(tok.String())

		if idx < length {
			sbld.WriteString("\n")
		}
	}

	return sbld.String()
}

func makeLiteralToken(ttype TokenType, lexeme string, literal any) Token {
	return newTokenWithLiteral(ttype,
		lexeme,
		literal,
		newPosition(0, 0),
		newPosition(0, 0))
}

// ** Tests:

func TestLucette(t *testing.T) {
	for testno, test := range tests {
		var (
			lx     *Lexer
			pr     *Parser
			tr     *Typer
			lexed  []Token
			parsed Node
			typed  TypedNode
			nnf    TypedNode
			simple TypedNode
			code   *Program
			diags  []Diagnostic
			err    error
		)

		lx = NewLexer(strings.NewReader(test.input))
		if lx == nil {
			t.Fatal("Didn't get a working lexer")
		}

		// Lexer test.
		t.Run(fmt.Sprintf("lexer test %04d", testno), func(t *testing.T) {
			lexed, err = lx.lex()
			if err != nil && test.err == nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if test.err != nil {
				if !errors.Is(err, test.err) {
					t.Fatalf("Error mismatch: %#v != %#v",
						err,
						test.err)
				}

				t.Logf("Query: %s\nError: %s",
					test.input,
					err.Error())

				// Consider the test a pass.
				return
			}

			t.Logf("Query: %s\nResult:\n%s",
				test.input,
				dumpTokens(lexed))

			if len(lexed) != len(test.tokens) {
				t.Fatalf("Token count mismatch: %d != %d",
					len(lexed),
					len(test.tokens))
			}

			for idx, tok := range test.tokens {
				if lexed[idx].tokenType != tok.tokenType {
					t.Errorf("%04d: Token type mismatch: %q != %q",
						idx,
						lexed[idx].tokenType.String(),
						tok.tokenType.String())
				}

				if lexed[idx].lexeme != tok.lexeme {
					t.Errorf("%04d: Lexeme mismatch: %q != %q",
						idx,
						lexed[idx].lexeme,
						tok.lexeme)
				}

				if lexed[idx].literal != tok.literal {
					t.Errorf("%04d: Literal mismatch: %q != %q",
						idx,
						lexed[idx].literal.value,
						tok.literal.value)
				}
			}
		})

		// If the lexer test is one that's testing error generation,
		// we shouldn't continue.
		if test.err != nil {
			continue
		}

		// Parser test.
		t.Run(fmt.Sprintf("parser test %04d", testno), func(t *testing.T) {
			pr = NewParser(lexed)

			parsed, diags = pr.Parse()

			parsed.Debug()

			if len(diags) > 0 {
				for idx := range diags {
					t.Logf("Diagnostic %02d: %s",
						idx,
						diags[idx].String())
				}
				t.Fatalf("Parser errors.")
			}
		})

		// Typer tests.
		t.Run(fmt.Sprintf("typer test %04d", testno), func(t *testing.T) {
			tr = NewTyper(MakeStructSchema())

			typed, diags = tr.Type(parsed)

			typed.Debug()

			if len(diags) > 0 {
				for idx := range diags {
					t.Logf("Diagnostic %02d: %s",
						idx,
						diags[idx].String())
				}
				t.Fatalf("Typer errors.")
			}
		})

		// NNF tests.
		t.Run(fmt.Sprintf("NNF test %04d", testno), func(t *testing.T) {
			nnf = ToNNF(typed)

			//nnf.Debug()
		})

		// Simplifier.
		t.Run(fmt.Sprintf("Simplifier test %04d", testno), func(t *testing.T) {
			simple = Simplify(nnf)

			simple.Debug()
		})

		// Code emitter.
		t.Run(fmt.Sprintf("Code emit test %04d", testno), func(t *testing.T) {
			code = NewProgram()
			code.Emit(simple)
			code.PrettyPrint(os.Stdout, NewDefaultPrettyPrinterOptions())
		})
	}
}

// ** Benchmarks:

func BenchmarkLucette(b *testing.B) {
	var (
		lexed  []Token
		parsed Node
		typed  TypedNode
		nnf    TypedNode
		//simple TypedNode
		diags []Diagnostic
		err   error
	)

	input := tests[6].input
	schema := MakeStructSchema()

	b.Run("Lexer", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			lx := NewLexer(strings.NewReader(input))

			if lexed, err = lx.lex(); err != nil {
				b.Fatalf("Unexpected error: %#v", err)
			}
		}
	})

	b.Run("Parser", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			pr := NewParser(lexed)

			if parsed, diags = pr.Parse(); len(diags) > 0 {
				b.Fatalf("Unexpected error: %#v", diags)
			}
		}
	})

	b.Run("Typer", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			tr := NewTyper(schema)

			if typed, diags = tr.Type(parsed); len(diags) > 0 {
				b.Fatalf("Unexpected error: %#v", diags)
			}
		}
	})

	b.Run("NNF", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			nnf = ToNNF(typed)
		}
	})

	b.Run("Simplifier", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			_ = Simplify(nnf)
		}
	})

}

// * lucette_test.go ends here.
