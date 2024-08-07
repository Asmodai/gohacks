// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// any_test.go --- `Any` tests.
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

package generics

import (
	"math"
	"testing"
)

func TestAnyString(t *testing.T) {
	goodArr := []string{"yes", "ja", "yes", "yes"}
	badArr := []string{"yes", "nein", "yes", "yes"}

	t.Run("Returns true for good array", func(t *testing.T) {
		res := Any(goodArr, func(elt string) bool {
			return elt == "ja"
		})

		if !res {
			t.Error("Unexpected result!")
		}
	})

	t.Run("Returns false for bad array", func(t *testing.T) {
		res := Any(badArr, func(elt string) bool {
			return elt == "ja"
		})

		if res {
			t.Error("Unexpected result!")
		}
	})
}

func TestAnyInt(t *testing.T) {
	goodArr := []int{1, 2, 5, 7, 9}
	badArr := []int{1, 3, 5, 7, 9}

	t.Run("Returns true for good array", func(t *testing.T) {
		res := Any(goodArr, func(elt int) bool {
			return elt%2 == 0
		})

		if !res {
			t.Error("Unexpected result!")
		}
	})

	t.Run("Returns false for bad array", func(t *testing.T) {
		res := Any(badArr, func(elt int) bool {
			return elt%2 == 0
		})

		if res {
			t.Error("Unexpected result!")
		}
	})
}

func TestAnyFloat(t *testing.T) {
	goodArr := []float64{1.0, 2.0, 5.0, 7.0, 9.0}
	badArr := []float64{1.0, 3.0, 5.0, 7.0, 9.0}

	t.Run("Returns true for good array", func(t *testing.T) {
		res := Any(goodArr, func(elt float64) bool {
			return math.Mod(elt, 2) == 0
		})

		if !res {
			t.Error("Unexpected result!")
		}
	})

	t.Run("Returns false for bad array", func(t *testing.T) {
		res := Any(badArr, func(elt float64) bool {
			return math.Mod(elt, 2) == 0
		})

		if res {
			t.Error("Unexpected result!")
		}
	})
}

// any_test.go ends here.
