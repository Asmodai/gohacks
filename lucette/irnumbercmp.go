// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// irnumbercmp.go --- IR numerical compare node.
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

// This node will resolve dynamically to any of the numeric operations.

// * Package:

//nolint:dupl
package lucette

// * Imports:

import (
	"fmt"

	"github.com/Asmodai/gohacks/debug"
)

// * Code:

// ** Structure:

type IRNumberCmp struct {
	Field string
	Op    ComparatorKind
	Value float64
}

// ** Methods:

// Generate the key.
func (n IRNumberCmp) Key() string {
	return fmt.Sprintf("cn|%s|%d|%g", n.Field, n.Op, n.Value)
}

// Display debugging information.
func (n IRNumberCmp) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Number Comparator")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.Field)
	dbg.Printf("Op:     %s", ComparatorKindToString(n.Op))
	dbg.Printf("Value:  %f", n.Value)

	dbg.End()
	dbg.Print()

	return dbg
}

// Generate opcode.
func (n IRNumberCmp) Emit(program *Program, trueLabel, falseLabel LabelID) {
	operator := GetNumberComparator(n.Op, OpNumberEQ)
	fidx := program.AddFieldConstant(n.Field)
	nidx := program.AddNumberConstant(n.Value)

	program.AppendIsn(OpLoadField, fidx)    // LDFLD fIdx
	program.AppendIsn(operator, nidx)       // <op> nIdx
	program.AppendJump(OpJumpNZ, trueLabel) // JNZ true continuation
	program.AppendJump(OpJump, falseLabel)  // JMP false continuation
}

// * irnumbercmp.go ends here.
