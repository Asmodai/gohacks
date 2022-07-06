/*
 * error_test.go --- Error tests.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
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
	"errors"
	"fmt"
	"testing"
)

func TestNewError(t *testing.T) {
	var fails bool = false

	e1 := NewError("TEST", "test %d", 1)
	e2 := NewErrorAndLog("TEST", "test %d", 2)

	if e1.Error() == "TEST: test 1" {
		t.Log("`NewError` works.")
	} else {
		t.Errorf("`NewError` returned %s", e1.Error())
		fails = true
	}

	if e2.Error() == "TEST: test 2" {
		t.Log("`NewErrorAndLog` works.")
	} else {
		t.Errorf("`NewErrorAndLog` returned %s", e2.Error())
		fails = true
	}

	if fails {
		t.Error("Error(s) occurred.")
		return
	}
}

func TestUnwrapping(t *testing.T) {
	inner := fmt.Errorf("inner error")
	outer := NewError("TEST", "%s", inner.Error())

	unwrapped := errors.Unwrap(outer)
	if unwrapped.Error() == "inner error" {
		t.Log("`Unwrap` works.")
		return
	}

	t.Errorf("`Unwrap` returned %s", unwrapped.Error())
}

func TestStringConvert(t *testing.T) {
	err := NewError("TEST", "no")

	if err.String() == "TEST: no" {
		t.Log("`String` works.")
		return
	}

	t.Error("`String` does not work.")
}

func TestJSONMarshal(t *testing.T) {
	t.Log("Can we marshal errors to JSON?")

	e1 := NewError("TEST", "no")
	json, err := e1.MarshalJSON()

	if err != nil {
		t.Errorf("JSON error: %s", err.Error())
		return
	}

	if string(json) == "\"TEST: no\"" {
		t.Log("Yes.")
		return
	}

	t.Errorf("Unexpected JSON: %s", string(json))
}

/* error_test.go ends here. */
