// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_eq_test.go --- EQ instruction test.
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

package dag

// * Imports:

import (
	"context"
	"testing"
)

// * Code:

func TestEQPredicate(t *testing.T) {
	const (
		good int = 42
		bad  int = 76
	)

	input := NewDataInputFromMap(map[string]any{"Numeric": float64(good)})
	builder := &EQBuilder{}
	pred1, _ := builder.Build("Numeric", good, nil, false)
	pred2, _ := builder.Build("Numeric", bad, nil, false)

	t.Run(pred1.String(), func(t *testing.T) {
		if !pred1.Eval(context.TODO(), input) {
			t.Errorf("%s - failed.  %v != %v",
				pred1.String(),
				good,
				good)
		}
	})

	t.Run(pred2.String(), func(t *testing.T) {
		if pred2.Eval(context.TODO(), input) {
			t.Errorf("%s - failed.  %v == %v",
				pred2.String(),
				good,
				bad)
		}
	})
}

// * predicate_eq_test.go ends here.
