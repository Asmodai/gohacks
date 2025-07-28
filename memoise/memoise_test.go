// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// memoise_test.go --- Memoisation tests.
//
// Copyright (c) 2024-2025 Paul Ward <paul@lisphacker.uk>
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

package memoise

// * Imports:

import (
	"sync"
	"testing"
	"time"

	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	testError error = errors.Base("this is an error")
)

// * Code:

// ** Tests:

func TestMemoise(t *testing.T) {
	t.Run("Works as expected", func(t *testing.T) {
		memo := NewMemoise()

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
			t.Errorf("Result mismatch: %v != %v", r1, r2)
		}
	})

	t.Run("Handles errors", func(t *testing.T) {
		memo := NewMemoise()

		_, err := memo.Check(
			"Test2",
			func() (any, error) { return nil, testError },
		)

		if !errors.Is(err, testError) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("Errors are not memoised", func(t *testing.T) {
		memo := NewMemoise()
		callCount := 0

		callback := func() (any, error) {
			callCount++
			return nil, testError
		}

		for i := 0; i < 3; i++ {
			_, _ = memo.Check("faily", callback)
		}

		if callCount != 3 {
			t.Errorf(
				"Expected callback to run 3 times, ran %d times",
				callCount,
			)
		}
	})

	t.Run("Resets", func(t *testing.T) {
		memo := NewMemoise()

		r1, _ := memo.Check(
			"Test1",
			func() (any, error) { return 100, nil },
		)

		// `Reset` should nuke the existing `Test` result.
		memo.Reset()

		r2, _ := memo.Check(
			"Test1",
			func() (any, error) { return 200, nil },
		)

		if r1 == r2 {
			t.Errorf("Result mismatch: %#v != %#v", r1, r2)
		}
	})
}

func TestMemoise_Concurrency(t *testing.T) {
	memo := NewMemoise()

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	counter := 0
	callback := func() (any, error) {
		time.Sleep(10 * time.Millisecond) // Simulate real work
		counter++
		return "hello", nil
	}

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := memo.Check("key", callback)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	if counter != 1 {
		t.Errorf("Expected callback to be called once, got %d", counter)
	}
}

// * memoise_test.go ends here.
