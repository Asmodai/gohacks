// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// bool_test.go --- Boolean conversion tests.
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

import (
	"testing"
)

// * Code:

// ** Types:

type expectedBool struct {
	name string
	val  any
	want bool
}

// ** Tests:

func TestToBool(t *testing.T) {
	want := []expectedBool{
		expectedBool{"bool", true, true},
		expectedBool{"float64", float64(1), true},
		expectedBool{"float32", float32(0), false},
		expectedBool{"int", int(1), true},
		expectedBool{"int8", int8(0), false},
		expectedBool{"int16", int16(1), true},
		expectedBool{"int32", int32(0), false},
		expectedBool{"int64", int64(1), true},
		expectedBool{"uint", uint(0), false},
		expectedBool{"uint8", uint8(1), true},
		expectedBool{"uint16", uint16(0), false},
		expectedBool{"uint32", uint32(1), true},
		expectedBool{"uint64", uint64(0), false},
		expectedBool{"string", "true", true},
		expectedBool{"string", "no", false},
	}

	t.Run("Numeric types", func(t *testing.T) {
		for idx := range want {
			res, ok := ToBool(want[idx].val)

			if !ok {
				t.Errorf("'%s' value not ok", want[idx].name)
			}

			if res != want[idx].want {
				t.Fatalf(
					"Value mismatch: %#v != %#v",
					res,
					want[idx].want,
				)
			}
		}
	})

	t.Run("Non-numeric types", func(t *testing.T) {
		res, ok := ToFloat64("Huh?")

		if ok {
			t.Error("Unexpectedly ok!")
		}

		if res != 0 {
			t.Errorf("Value mismatch: %#v != 0", res)
		}
	})
}

// ** Benchmarks:

func BenchmarkToBool(b *testing.B) {
	b.ReportAllocs()

	for val := range b.N {
		ToBool(val)
	}
}

// * bool_test.go ends here.
