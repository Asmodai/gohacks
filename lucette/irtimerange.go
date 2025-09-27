// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// irnumberrange.go --- IR date/time `Range' node.
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
	"strconv"

	"github.com/Asmodai/gohacks/debug"
)

// * Code:

// ** Structure:

type IRTimeRange struct {
	Field string
	Lo    *int64
	Hi    *int64
	IncL  bool
	IncH  bool
}

// ** Methods:

// Generate the key.
func (n IRTimeRange) Key() string {
	var low, high string

	if n.Lo != nil {
		low = strconv.FormatInt(*n.Lo, 10)
	}

	if n.Hi != nil {
		high = strconv.FormatInt(*n.Hi, 10)
	}

	return fmt.Sprintf("rt|%s|%s|%s|%t|%t",
		n.Field,
		low,
		high,
		n.IncL,
		n.IncH)
}

// Display debugging information.
func (n IRTimeRange) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Date/Time Range")

	dbg.Init(params...)
	dbg.Printf("Field:          %s", n.Field)

	if n.Lo != nil {
		dbg.Printf("Low:            %d", *n.Lo)
	}

	if n.Hi != nil {
		dbg.Printf("High:           %d", *n.Hi)
	}

	dbg.Printf("Inclusive Low:  %v", n.IncL)
	dbg.Printf("Inclusive High: %v", n.IncH)

	dbg.End()
	dbg.Print()

	return dbg
}

// Generate opcode.
func (n IRTimeRange) Emit(program *Program, trueLabel, falseLabel LabelID) {
	fidx := program.AddFieldConstant(n.Field)
	lidx := 0
	hidx := 0

	if n.Lo != nil {
		lidx = program.AddTimeConstant(*n.Lo)
	}

	if n.Hi != nil {
		hidx = program.AddTimeConstant(*n.Hi)
	}

	program.AppendIsn(OpLoadField, fidx) // LDFLD fIdx

	program.AppendIsn(OpTimeRange, // RNG.T lIdx hIdx IncL IncH
		lidx,
		hidx,
		n.IncL,
		n.IncH)

	program.AppendJump(OpJumpNZ, trueLabel) // JNZ true continuation
	program.AppendJump(OpJump, falseLabel)  // JMP false continuation
}

// * irnumberrange.go ends here.
