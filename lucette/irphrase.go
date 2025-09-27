// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// irphrase.go --- IR `Phrase' node.
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
)

// * Code:

// ** Structure:

type IRPhrase struct {
	Field     string
	Phrase    string
	Proximity int
	Fuzz      *float64
	Boost     *float64
}

// ** Methods:

func (n IRPhrase) HasWildcard() bool {
	return strings.ContainsAny(n.Phrase, "?*")
}

// Generate the key.
func (n IRPhrase) Key() string {
	return "phrase|" + n.Field + "|" + n.Phrase
}

// Display debugging information.
func (n IRPhrase) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Phrase")

	dbg.Init(params...)
	dbg.Printf("Field:     %s", n.Field)
	dbg.Printf("Phrase:    %q", n.Phrase)
	dbg.Printf("Proximity: %d", n.Proximity)

	if n.Fuzz != nil {
		dbg.Printf("Fuzziness: %g", *n.Fuzz)
	}

	if n.Boost != nil {
		dbg.Printf("Boost:     %g", *n.Boost)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// Generate opcode.
func (n IRPhrase) Emit(program *Program, trueLabel, falseLabel LabelID) {
	fidx := program.AddFieldConstant(n.Field)
	sidx := program.AddStringConstant(n.Phrase)

	program.AppendIsn(OpLoadField, fidx) // LDFLD fIdx

	if n.Fuzz != nil {
		program.AppendIsn(OpLoadFuzzy, *n.Fuzz) // LDFZY fuzz
	}

	if n.Boost != nil {
		program.AppendIsn(OpLoadBoost, *n.Boost) // LDBST boost
	}

	program.AppendIsn(OpPhrase, sidx, n.Proximity) // PHR.S sIdx prox
	program.AppendIsn(OpJumpNZ, trueLabel)         // JNZ true cont.
	program.AppendIsn(OpJump, falseLabel)          // JMP false cont.
}

// * irphrase.go ends here.
