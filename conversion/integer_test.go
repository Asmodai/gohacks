// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// integer_test.go --- Integer conversion tests.
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

package conversion

// * Imports:

import "testing"

// * Code:

// ** Tests:

func TestToInt64(t *testing.T) {
	want := []struct {
		name   string
		val    any
		want   int64
		result bool
	}{
		{"string", "42", int64(0), false},
		{"int8", int8(42), int64(42), true},
		{"float32", float32(42.4), int64(0), false},
		{"bool", true, int64(1), true},
	}

	t.Run("Base types", func(t *testing.T) {
		for idx := range want {
			res, ok := ToInt64(want[idx].val)

			if ok != want[idx].result {
				t.Fatalf("%q: ok != %#v",
					want[idx].name,
					want[idx].result,
				)
			}

			if res != want[idx].want {
				t.Fatalf("%q: %#v != %#v",
					want[idx].name,
					want[idx].result,
					res,
				)
			}
		}
	})
}

func TestToUint64(t *testing.T) {
	want := []struct {
		name   string
		val    any
		want   uint64
		result bool
	}{
		{"string", "42", uint64(0), false},
		{"int8", int8(42), uint64(42), true},
		{"float32", float32(42.4), uint64(0), false},
		{"bool", true, uint64(1), true},
	}

	t.Run("Base types", func(t *testing.T) {
		for idx := range want {
			res, ok := ToUint64(want[idx].val)

			if ok != want[idx].result {
				t.Fatalf("%q: ok != %#v",
					want[idx].name,
					want[idx].result,
				)
			}

			if res != want[idx].want {
				t.Fatalf("%q: %#v != %#v",
					want[idx].name,
					want[idx].result,
					res,
				)
			}
		}
	})
}

// * integer_test.go ends here.
