// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// mapmap_test.go --- Map mapping tests.
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
	"testing"
)

func TestMapMapString(t *testing.T) {
	data := map[string]string{
		"one":   "Learn Go",
		"two":   "???",
		"three": "Profit!",
	}

	res := MapMap(data, func(key string, val string) bool {
		return len(key) > 3 && len(val) > 3
	})

	if len(res) > 1 {
		t.Errorf("Unexpected result: %#v", res)
	}
}

func TestMapMapInt(t *testing.T) {
	data := map[int]string{
		1: "Hello",
		2: "There",
		3: "World",
	}

	res := MapMap(data, func(key int, _ string) bool {
		return key%2 != 0
	})

	if len(res) > 2 {
		t.Errorf("Unexpected result: %#v", res)
	}
}

func TestMapMapFloat(t *testing.T) {
	data := map[string]float64{
		"mon": 1.4,
		"tue": 0.5,
		"wed": 2.1,
		"thu": 0.1,
		"fri": 1.9,
	}

	res := MapMap(data, func(_ string, val float64) bool {
		return val > 2.0
	})

	if len(res) > 1 {
		t.Errorf("Unexpected result: %#v", res)
	}
}

// mapmap_test.go ends here.
