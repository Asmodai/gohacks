// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// irnot.go --- IR `Not' node.
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

import "github.com/Asmodai/gohacks/debug"

// * Code:

// ** Structure:

type IRNot struct {
	Kid IRNode
}

// ** Methods:

// Generate key.
func (n IRNot) Key() string {
	return "not|" + n.Kid.Key()
}

// Display debugging information.
func (n IRNot) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("NOT Node")

	dbg.Init(params...)
	dbg.Printf("Child:")

	if n.Kid != nil {
		n.Kid.Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// Emit opcode.
func (n IRNot) Emit(program *Program, trueLabel, falseLabel LabelID) {
	if n.Kid == nil {
		program.AppendJump(OpJump, falseLabel)

		return
	}

	cont := program.NewLabel()

	n.Kid.Emit(program, trueLabel, cont)

	program.AppendIsn(OpNot)
	program.BindLabel(cont)
}

// * irnot.go ends here.
