// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// maplist_test.go --- List mapping tests.
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

func TestMapListString(t *testing.T) {
	list := []string{"yes", "ja", "da", "si", "oui"}

	t.Run("Produces list of matches", func(t *testing.T) {
		res := MapList(list, func(elt string) bool {
			return elt == "da"
		})

		if len(res) != 1 {
			t.Errorf("Unexpected result: %#v", res)
		}
	})

	t.Run("Produces empty list for no matches", func(t *testing.T) {
		res := MapList(list, func(elt string) bool {
			return elt == "yup"
		})

		if len(res) != 0 {
			t.Errorf("Unexpected result: %#v", res)
		}
	})
}

func TestMapListInt(t *testing.T) {
	list := []int{2, 3, 4, 5, 6, 7, 8}

	t.Run("Produces list of matches", func(t *testing.T) {
		res := MapList(list, func(elt int) bool {
			return elt%2 == 0
		})

		if len(res) != 4 {
			t.Errorf("Unexpected result: %#v", res)
		}
	})

	t.Run("Produces empty list for no matches", func(t *testing.T) {
		res := MapList(list, func(elt int) bool {
			return elt == 42
		})

		if len(res) != 0 {
			t.Errorf("Unexpected result: %#v", res)
		}
	})
}

func TestMapListFloat(t *testing.T) {
	list := []float64{2, 3, 4, 5, 6, 7, 8}

	t.Run("Produces list of matches", func(t *testing.T) {
		res := MapList(list, func(elt float64) bool {
			return math.Mod(elt, 2) == 0
		})

		if len(res) != 4 {
			t.Errorf("Unexpected result: %#v", res)
		}
	})

	t.Run("Produces empty list for no matches", func(t *testing.T) {
		res := MapList(list, func(elt float64) bool {
			return elt > 42.0
		})

		if len(res) != 0 {
			t.Errorf("Unexpected result: %#v", res)
		}
	})
}

// maplist_test.go ends here.
