// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// literal_test.go --- Literal tests.
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

import (
	"testing"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

// * Variables:

// * Code:

func TestLiteral(t *testing.T) {
	strval := "foo"
	errval := errors.Base("test error")
	errstr := "<error: " + errval.Error() + ">"
	vallit := literal{value: strval}
	errlit := literal{err: errval}

	t.Run("String", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			got := vallit.String()

			if got != strval {
				t.Errorf("Value mismatch: %q != %q",
					got,
					strval)
			}
		})

		t.Run("Error", func(t *testing.T) {
			got := errlit.String()

			if got != errstr {
				t.Errorf("Error mismatch: %q != %q",
					got,
					errval)
			}
		})
	})

	t.Run("IsError", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			if vallit.IsError() {
				t.Error("Value thinks its an error")
			}
		})

		t.Run("Error", func(t *testing.T) {
			if !errlit.IsError() {
				t.Error("Error thinks its not an error")
			}
		})
	})

	t.Run("Error", func(t *testing.T) {
		t.Run("Value", func(t *testing.T) {
			if len(vallit.Error()) > 0 {
				t.Errorf("Value has error: %q",
					vallit.Error())
			}
		})

		t.Run("Error", func(t *testing.T) {
			if len(errlit.Error()) == 0 {
				t.Errorf("Error has no error: %q",
					errlit.Error())
			}
		})
	})
}

// * literal_test.go ends here.
