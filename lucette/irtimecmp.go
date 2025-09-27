// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// irnumbercmp.go --- IR date/time compare node.
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

// This node will resolve dynamically to any of the date/time operations.

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

type IRTimeCmp struct {
	Field string
	Op    ComparatorKind
	Value int64
}

// ** Methods:

// Generate the key.
func (n IRTimeCmp) Key() string {
	return fmt.Sprintf("ct|%s|%d|%d", n.Field, n.Op, n.Value)
}

// Display debugging information.
func (n IRTimeCmp) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Date/Time Comparator")

	dbg.Init(params...)
	dbg.Printf("Field:  %s", n.Field)
	dbg.Printf("Op:     %s", ComparatorKindToString(n.Op))
	dbg.Printf("Value:  %d", n.Value)

	dbg.End()
	dbg.Print()

	return dbg
}

// Generate opcode.
func (n IRTimeCmp) Emit(program *Program, trueLabel, falseLabel LabelID) {
	operator := GetNumberComparator(n.Op, OpTimeEQ)
	fidx := program.AddFieldConstant(n.Field)
	tidx := program.AddTimeConstant(n.Value)

	program.AppendIsn(OpLoadField, fidx)    // LDFLD fIdx
	program.AppendIsn(operator, tidx)       // <op> tIdx
	program.AppendJump(OpJumpNZ, trueLabel) // JNZ true continuation
	program.AppendJump(OpJump, falseLabel)  // JMP false continuation
}

// * irnumbercmp.go ends here.
