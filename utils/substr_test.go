// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// substr_test.go --- Substring tests.
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

package utils

import "testing"

func TestSubstr(t *testing.T) {
	var input string = "One two three"

	t.Run("Works with sane arguments", func(t *testing.T) {
		if output := Substr(input, 4, 3); output != "two" {
			t.Errorf("Unexpected '%v'", output)
		}
	})

	t.Run("Works when length is longer than input", func(t *testing.T) {
		if output := Substr(input, 4, 100); output != "two three" {
			t.Errorf("Unexpected '%v'", output)
		}
	})

	t.Run("Returns empty string with insane args", func(t *testing.T) {
		if output := Substr(input, 100, 24); output != "" {
			t.Errorf("Unexpected '%v'", output)
		}
	})
}

// substr_test.go ends here.
