// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// pop_test.go --- Test of array pop function.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

//
//
//

// * Package:

package utils

// * Imports:

import "testing"

// * Code:

// ** Tests:

func TestPop(t *testing.T) {
	var arr []string = []string{"one", "two", "three", "four"}
	var elt string
	var last string
	var rest []string

	for len(arr) > 0 {
		elt = arr[len(arr)-1]
		last, rest = Pop(arr)

		if last != elt {
			t.Errorf("Unexpected '%v'", last)
			return
		}

		if len(rest) != len(arr)-1 {
			t.Errorf(
				"Length of rest does not match: %d != %d",
				len(rest),
				len(arr)-1,
			)
		}

		arr = arr[:len(arr)-1]
	}

	// nolint:ineffassign,staticcheck
	// Now test zero-length array.
	last, rest = Pop([]string{})
	if last != "" {
		t.Errorf("Unexpected '%v'", last)
	}
}

// * pop_test.go ends here.
