/*
 * pair_test.go --- Pair type tests.
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
