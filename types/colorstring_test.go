// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// colorstring_test.go --- Colour string tests.
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

// * Comments:

//
//
//

// * Package:

package types

// * Imports:

import (
	"testing"
)

// * Code:

// ** Tests:

func TestColorString(t *testing.T) {
	onlyFG := "\x1b[0;31;49mTest\x1b[0m"
	onlyBG := "\x1b[0;39;41mTest\x1b[0m"
	onlyAttr := "\x1b[1;39;49mTest\x1b[0m"

	t.Run("Sets FG properly", func(t *testing.T) {
		cl := NewColorString()

		cl.SetString("Test")
		cl.SetFG(Red)

		if cl.String() != onlyFG {
			t.Errorf("Unexpected: %#v != %#v",
				cl.String(),
				onlyFG)
		}
	})

	t.Run("Sets BG properly", func(t *testing.T) {
		cl := NewColorString()

		cl.SetString("Test")
		cl.SetBG(Red)

		if cl.String() != onlyBG {
			t.Errorf("Unexpected: %#v != %#v",
				cl.String(),
				onlyBG)
		}
	})

	t.Run("Sets attribute properly", func(t *testing.T) {
		cl := NewColorString()

		cl.SetString("Test")
		cl.SetBold()

		if cl.String() != onlyAttr {
			t.Errorf("Unexpected '%#v'", cl.String())
		}
	})
}

// * colorstring_test.go ends here.
