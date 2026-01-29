// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// position_test.go --- Position tests.
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

package lucette

// * Imports:

import "testing"

// * Code:

func TestPosition(t *testing.T) {
	var testPos Position

	line := 12
	column := 42
	wantstr := "12:42"

	t.Run("NewPosition", func(t *testing.T) {
		testPos = NewPosition(line, column)

		if testPos.Line != line {
			t.Errorf("Line mismatch: %d != %d",
				testPos.Line,
				line)
		}

		if testPos.Column != column {
			t.Errorf("Column mismatch: %d != %d",
				testPos.Column,
				column)
		}
	})

	t.Run("String", func(t *testing.T) {
		got := testPos.String()

		if got != wantstr {
			t.Errorf("String mismatch: %q != %q", got, wantstr)
		}
	})
}

// * position_test.go ends here.
