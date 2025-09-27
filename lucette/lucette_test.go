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
		tokens []LexedToken
		err    error
	}{
		{
			input:  "field:string",
			tokens: []LexedToken{},
			err:    ErrUnexpectedBareword,
		}, {
			input: `level:42 && message:"whee"`,
			tokens: []LexedToken{
				makeLiteralToken(TokenField, "", "level"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenNumber, "", float64(42)),
				makeLiteralToken(TokenAnd, "&&", ""),
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "whee"),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `'the level':42 && 'the message':"whee"`,
			tokens: []LexedToken{
				makeLiteralToken(TokenField, "", "the level"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenNumber, "", float64(42)),
				makeLiteralToken(TokenAnd, "&&", ""),
				makeLiteralToken(TokenField, "", "the message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "whee"),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `message:"value"~5 && level:95`,
			tokens: []LexedToken{
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "value"),
				makeLiteralToken(TokenTilde, "~", ""),
				makeLiteralToken(TokenNumber, "", float64(5)),
				makeLiteralToken(TokenAnd, "&&", ""),
				makeLiteralToken(TokenField, "", "level"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenNumber, "", float64(95)),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `percent:3_123_456 || percent:3.1e+02`,
			tokens: []LexedToken{
				makeLiteralToken(TokenField, "", "percent"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenNumber, "", float64(3_123_456)),
				makeLiteralToken(TokenOr, "||", ""),
				makeLiteralToken(TokenField, "", "percent"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenNumber, "", float64(3.1e+02)),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `(message:"yes" or message:"maybe") and !message:"no"`,
			tokens: []LexedToken{
				makeLiteralToken(TokenLParen, "(", ""),
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "yes"),
				makeLiteralToken(TokenOr, "or", ""),
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "maybe"),
				makeLiteralToken(TokenRParen, ")", ""),
				makeLiteralToken(TokenAnd, "and", ""),
				makeLiteralToken(TokenNot, "!", ""),
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "no"),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `(message:"ye*" or message:"n?") and (not message:"chuckles" and !+message:"ass")`,
			tokens: []LexedToken{
				makeLiteralToken(TokenLParen, "(", ""),
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "ye*"),
				makeLiteralToken(TokenOr, "or", ""),
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "n?"),
				makeLiteralToken(TokenRParen, ")", ""),
				makeLiteralToken(TokenAnd, "and", ""),
				makeLiteralToken(TokenLParen, "(", ""),
				makeLiteralToken(TokenNot, "not", ""),
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "chuckles"),
				makeLiteralToken(TokenAnd, "and", ""),
				makeLiteralToken(TokenNot, "!", ""),
				makeLiteralToken(TokenPlus, "+", ""),
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "ass"),
				makeLiteralToken(TokenRParen, ")", ""),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `(level:[1 to *})`,
			tokens: []LexedToken{
				makeLiteralToken(TokenLParen, "(", ""),
				makeLiteralToken(TokenField, "", "level"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenLBracket, "[", ""),
				makeLiteralToken(TokenNumber, "", float64(1)),
				makeLiteralToken(TokenTo, "to", ""),
				makeLiteralToken(TokenStar, "*", ""),
				makeLiteralToken(TokenRCurly, "}", ""),
				makeLiteralToken(TokenRParen, ")", ""),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `(source:["192.168.1.1" to "192.168.1.254"])`,
			tokens: []LexedToken{
				makeLiteralToken(TokenLParen, "(", ""),
				makeLiteralToken(TokenField, "", "source"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenLBracket, "[", ""),
				makeLiteralToken(TokenPhrase, "", "192.168.1.1"),
				makeLiteralToken(TokenTo, "to", ""),
				makeLiteralToken(TokenPhrase, "", "192.168.1.254"),
				makeLiteralToken(TokenRBracket, "]", ""),
				makeLiteralToken(TokenRParen, ")", ""),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `!source:"8.8.8.8"`,
			tokens: []LexedToken{
				makeLiteralToken(TokenNot, "!", ""),
				makeLiteralToken(TokenField, "", "source"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenPhrase, "", "8.8.8.8"),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `source:>="8.8.8.8"`,
			tokens: []LexedToken{
				makeLiteralToken(TokenField, "", "source"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenGTE, ">=", ""),
				makeLiteralToken(TokenPhrase, "", "8.8.8.8"),
				makeLiteralToken(TokenEOF, "", ""),
			},
		}, {
			input: `message:/testing/i`,
			tokens: []LexedToken{
				makeLiteralToken(TokenField, "", "message"),
				makeLiteralToken(TokenColon, ":", ""),
				makeLiteralToken(TokenRegex, "", "(?i)testing"),
				makeLiteralToken(TokenEOF, "", ""),
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

func dumpTokens(toks []LexedToken) string {
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

func makeLiteralToken(token Token, lexeme string, literal any) LexedToken {
	return NewLexedTokenWithLiteral(token,
		lexeme,
		literal,
		NewPosition(0, 0),
		NewPosition(0, 0))
}

// ** Tests:

func TestLucette(t *testing.T) {
	for testno, test := range tests {
		var (
			lx     Lexer
			pr     Parser
			tr     Typer
			lexed  []LexedToken
			parsed ASTNode
			typed  IRNode
			nnf    IRNode
			simple IRNode
			code   *Program
			diags  []Diagnostic
			err    error
		)

		lx = NewLexer()
		if lx == nil {
			t.Fatal("Didn't get a working lexer")
		}

		// Lexer test.
		t.Run(fmt.Sprintf("lexer test %04d", testno), func(t *testing.T) {
			lexed, err = lx.Lex(strings.NewReader(test.input))
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
				if lexed[idx].Token != tok.Token {
					t.Errorf("%04d: Token type mismatch: %q != %q",
						idx,
						lexed[idx].Token.String(),
						tok.Token.String())
				}

				if lexed[idx].Lexeme != tok.Lexeme {
					t.Errorf("%04d: Lexeme mismatch: %q != %q",
						idx,
						lexed[idx].Lexeme,
						tok.Lexeme)
				}

				if lexed[idx].Literal != tok.Literal {
					t.Errorf("%04d: Literal mismatch: %q != %q",
						idx,
						lexed[idx].Literal.Value,
						tok.Literal.Value)
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
			pr = NewParser()

			parsed, diags = pr.Parse(lexed)

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
			nf := NewNNF()
			nnf = nf.NNF(typed)

			//nnf.Debug()
		})

		// Simplifier.
		t.Run(fmt.Sprintf("Simplifier test %04d", testno), func(t *testing.T) {
			si := NewSimplifier()
			simple = si.Simplify(nnf)

			simple.Debug()
		})

		// Code emiTokener.
		t.Run(fmt.Sprintf("Code emit test %04d", testno), func(t *testing.T) {
			code = NewProgram()
			code.Emit(simple)

			ds := NewDefaultDisassembler()
			ds.SetProgram(code)
			ds.Dissassemble(os.Stdout)
			//code.PrettyPrint(os.Stdout, NewDefaultPreTokenyPrinterOptions())
		})
	}
}

// ** Benchmarks:

func BenchmarkLucette(b *testing.B) {
	var (
		lexed  []LexedToken
		parsed ASTNode
		typed  IRNode
		nnf    IRNode
		simple IRNode
		diags  []Diagnostic
		err    error
	)

	input := tests[6].input
	schema := MakeStructSchema()

	b.Run("Lexer", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			lx := NewLexer()

			if lexed, err = lx.Lex(strings.NewReader(input)); err != nil {
				b.Fatalf("Unexpected error: %#v", err)
			}
		}
	})

	b.Run("Parser", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			pr := NewParser()

			if parsed, diags = pr.Parse(lexed); len(diags) > 0 {
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
			nf := NewNNF()

			nnf = nf.NNF(typed)
		}
	})

	b.Run("Simplifier", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			si := NewSimplifier()

			simple = si.Simplify(nnf)
		}
	})

	b.Run("Compiler", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			code := NewProgram()

			code.Emit(simple)
		}
	})
}

// * lucette_test.go ends here.
