// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// irnumberrange.go --- IR IP address `Range' node.
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
	"net/netip"

	"github.com/Asmodai/gohacks/debug"
)

// * Code:

// ** Structure:

type IRIPRange struct {
	Field string
	Lo    netip.Addr
	Hi    netip.Addr
	IncL  bool
	IncH  bool
}

// ** Methods:

// Generate the key.
func (n IRIPRange) Key() string {
	return fmt.Sprintf("rip|%s|%s|%s|%t|%t",
		n.Field,
		n.Lo.String(),
		n.Hi.String(),
		n.IncL,
		n.IncH)
}

// Display debugging information.
func (n IRIPRange) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("IP Address Range")

	dbg.Init(params...)
	dbg.Printf("Field:          %s", n.Field)
	dbg.Printf("Low:            %s", n.Lo.String())
	dbg.Printf("High:           %s", n.Hi.String())
	dbg.Printf("Increment Low:  %v", n.IncL)
	dbg.Printf("Increment High: %v", n.IncH)

	dbg.End()
	dbg.Print()

	return dbg
}

// Generate opcode.
func (n IRIPRange) Emit(program *Program, trueLabel, falseLabel LabelID) {
	fidx := program.AddFieldConstant(n.Field)
	lidx := program.AddIPConstant(n.Lo)
	hidx := program.AddIPConstant(n.Hi)

	program.AppendIsn(OpLoadField, fidx) // LDFLD fIdx

	program.AppendIsn(OpIPRange, // RNG.IP lIdx hIdx IncL IncH
		lidx,
		hidx,
		n.IncL,
		n.IncH)

	program.AppendJump(OpJumpNZ, trueLabel) // JNZ true continuation
	program.AppendJump(OpJump, falseLabel)  // JMP false continuation
}

// * irnumberrange.go ends here.
