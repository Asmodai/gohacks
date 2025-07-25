// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// memoise_test.go --- Memoisation tests.
//
// Copyright (c) 2024 Paul Ward <asmodai@gmail.com>
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

package memoise

import (
	"gitlab.com/tozd/go/errors"

	"testing"
)

var (
	memo      Memoise
	testError error = errors.Base("this is an error")
)

func TestMemoise(t *testing.T) {
	t.Run("Works as expected", func(t *testing.T) {
		memo = NewMemoise()

		r1, _ := memo.Check(
			"Test1",
			func() (any, error) { return 42, nil },
		)

		// r2 is sufficiently different from r1 so we can detect
		// whether there's a miss or not.
		r2, _ := memo.Check(
			"Test1",
			func() (any, error) { return 84, nil },
		)

		if r1 != r2 {
			t.Errorf("Result mismatch: %v!=%v", r1, r2)
		}
	})

	t.Run("Handles errors", func(t *testing.T) {
		memo = NewMemoise()

		_, err := memo.Check(
			"Test2",
			func() (any, error) { return nil, testError },
		)

		if !errors.Is(err, testError) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})
}

// memoise_test.go ends here.
