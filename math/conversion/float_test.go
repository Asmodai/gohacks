// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// float_test.go --- Floating-point conversion tests.
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
//
//

// * Package:

package conversion

// * Imports:

import (
	"testing"
)

// * Code:

// ** Types:

type expected struct {
	name string
	val  any
	want float64
}

// ** Tests:

func TestToFloat64(t *testing.T) {
	want := []expected{
		expected{"float64", float64(42.0), float64(42.0)},
		expected{"float32", float32(42.0), float64(42.0)},
		expected{"int", int(42), float64(42.0)},
		expected{"int8", int8(42), float64(42.0)},
		expected{"int16", int16(42), float64(42.0)},
		expected{"int32", int32(42), float64(42.0)},
		expected{"int64", int64(42), float64(42.0)},
		expected{"uint", uint(42), float64(42.0)},
		expected{"uint8", uint8(42), float64(42.0)},
		expected{"uint16", uint16(42), float64(42.0)},
		expected{"uint32", uint32(42), float64(42.0)},
		expected{"uint64", uint64(42), float64(42.0)},
	}

	t.Run("Numeric types", func(t *testing.T) {
		for idx := range want {
			res, ok := ToFloat64(want[idx].val)

			if !ok {
				t.Errorf("'%s' value not ok", want[idx].name)
				return
			}

			if res != want[idx].want {
				t.Errorf(
					"Value mismatch: %#v != %#v",
					res,
					want[idx].want,
				)
				return
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

// * float_test.go ends here.
