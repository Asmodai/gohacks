// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// irand.go --- IR `And' node.
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

//nolint:dupl
package lucette

// * Imports:

import (
	"sort"
	"strings"

	"github.com/Asmodai/gohacks/debug"
)

// * Code:

// ** Structure:

type IRAnd struct {
	Kids []IRNode
}

// ** Methods:

// Generate key.
func (n IRAnd) Key() string {
	keys := make([]string, 0, len(n.Kids))

	for _, kid := range n.Kids {
		keys = append(keys, kid.Key())
	}

	sort.Strings(keys)

	return "and|" + strings.Join(keys, "|")
}

// Display debugging information.
func (n IRAnd) Debug(params ...any) *debug.Debug {
	dbg := debug.NewDebug("Typed AND Node")

	dbg.Init(params...)
	dbg.Printf("Children:")

	for idx := range n.Kids {
		n.Kids[idx].Debug(&dbg)
	}

	dbg.End()
	dbg.Print()

	return dbg
}

// Emit opcode.
func (n IRAnd) Emit(program *Program, tLabel, fLabel LabelID) {
	kidlen := len(n.Kids)

	if kidlen == 0 {
		program.AppendJump(OpJump, tLabel)

		return
	}

	for idx := range kidlen - 1 {
		continuation := program.NewLabel()

		n.Kids[idx].Emit(program, continuation, fLabel)
		program.BindLabel(continuation)
	}

	n.Kids[kidlen-1].Emit(program, tLabel, fLabel)
}

// * irand.go ends here.
