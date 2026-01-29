// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_iir_test.go --- IIR tests.
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

package dag

// * Imports:

import (
	"context"
	"testing"
)

// * Code:

func TestIIRPredicate(t *testing.T) {
	var (
		goodHaystack []int = []int{4, 8}
		badHaystack  []int = []int{100, 200}
		good         int   = 4
	)

	input := NewDataInputFromMap(map[string]any{"Range": good})
	builder := &IIRBuilder{}
	pred1, _ := builder.Build("Range", goodHaystack, nil, false)
	pred2, _ := builder.Build("Range", badHaystack, nil, false)

	t.Run(pred1.String(), func(t *testing.T) {
		if !pred1.Eval(context.TODO(), input) {
			t.Errorf("%s - failed.  ! %v < %v < %v",
				pred1.String(),
				goodHaystack[0],
				good,
				goodHaystack[1])
		}
	})

	t.Run(pred2.String(), func(t *testing.T) {
		if pred2.Eval(context.TODO(), input) {
			t.Errorf("%s - failed.  %v < %v < %v",
				pred1.String(),
				badHaystack[0],
				good,
				badHaystack[1])
		}
	})
}

// * predicate_iir_test.go ends here.
