// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// span_test.go --- Span tests.
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

package lucette

// * Imports:

import "testing"

// * Code:

func TestSpan(t *testing.T) {
	var span *Span

	start := NewPosition(1, 10)
	end := NewPosition(1, 42)
	str := "1:10-1:42"

	t.Run("NewSpan", func(t *testing.T) {
		span = NewSpan(start, end)

		if span.start != start || span.end != end {
			t.Fatal("Did not create span")
		}
	})

	t.Run("Start", func(t *testing.T) {
		if span.Start() != start {
			t.Errorf("Mismatch: %#v != %#v", span.Start(), start)
		}
	})

	t.Run("End", func(t *testing.T) {
		if span.End() != end {
			t.Errorf("Mismatch: %#v != %#v", span.End(), end)
		}
	})

	t.Run("String", func(t *testing.T) {
		if span.String() != str {
			t.Errorf("Mismatch: %#v != %q", span.String(), str)
		}
	})
}

// * span_test.go ends here.
