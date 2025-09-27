// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// irglob.go --- IR `Glob' node.
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

//nolint:dupl
package lucette

// * Imports:

import "github.com/Asmodai/gohacks/debug"

// * Code:

// ** Structure:

type IRGlob struct {
	Field string
	Glob  string
}

// ** Methods:

// Generate the key.
func (n IRGlob) Key() string {
	return "glob|" + n.Field + "|" + n.Glob
}

// Display debugging information.
func (n IRGlob) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Glob")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.Field)
	dbg.Printf("Glob:   %q", n.Glob)

	dbg.End()
	dbg.Print()

	return dbg
}

// Emit opcode.
func (n IRGlob) Emit(program *Program, trueLabel, falseLabel LabelID) {
	fidx := program.AddFieldConstant(n.Field)
	sidx := program.AddStringConstant(n.Glob)

	program.AppendIsn(OpLoadField, fidx)    // LDFLD fIdx
	program.AppendIsn(OpGlob, sidx)         // GLB.S sIdx
	program.AppendJump(OpJumpNZ, trueLabel) // JNZ true continuation
	program.AppendJump(OpJump, falseLabel)  // JMP false continuation
}

// * irglob.go ends here.
