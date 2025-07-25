// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// pair_test.go --- Pair type tests.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

package types

import (
	"testing"
)

func TestPairType(t *testing.T) {
	car := "The Car"
	cdr := 42

	t.Run("NewEmptyPair", func(t *testing.T) {
		pair := NewEmptyPair()

		if pair.First != nil || pair.Second != nil {
			t.Errorf("Unexpected pair contents: %v", pair)
		}
	})

	t.Run("NewPair(value, nil)", func(t *testing.T) {
		pair := NewPair(car, nil)

		if pair.First != car || pair.Second != nil {
			t.Errorf("Unexpected pair contents: %v", pair)
		}
	})

	t.Run("NewPair(value, value)", func(t *testing.T) {
		pair := NewPair(car, cdr)
		if pair.First != car || pair.Second != cdr {
			t.Errorf("Unexpected pair contents: %v", pair)
		}
	})

	t.Run("String", func(t *testing.T) {
		pair := NewPair(car, cdr)

		if pair.String() != `("The Car" 42)` {
			t.Errorf("Unexpected string representation: %s", pair.String())
		}
	})
}

// pair_test.go ends here.
