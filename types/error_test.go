/*
 * error_test.go --- Error tests.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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
