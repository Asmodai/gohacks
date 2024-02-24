/*
 * pair_test.go --- Pair type tests.
 *
 * Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package types

import (
	"testing"
)

func TestPairType(t *testing.T) {
	var fail bool = false

	car := "The Car"
	cdr := 42

	p1 := NewEmptyPair()
	p2 := NewPair(car, nil)
	p3 := NewPair(car, cdr)

	if p1.First == nil && p1.Second == nil {
		t.Log("`NewEmptyPair` works.")
	} else {
		t.Error("`NewEmptyPair` does not work.")
		fail = true
	}

	if p2.First == car && p2.Second == nil {
		t.Log("`NewPair(value, nil)` works.")
	} else {
		t.Error("`NewPair(value, nil)` does not work.")
		fail = true
	}

	if p3.First == car && p3.Second == cdr {
		t.Log("`NewPair(value, value)` works.")
	} else {
		t.Error("`NewPair(value, value)` does not work.")
		fail = true
	}

	if fail {
		t.Error("Test(s) failed.")
	}
}

/* pair_test.go ends here. */
