/*
 * pop_test.go --- Test of array pop function.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package utils

import "testing"

func TestPop(t *testing.T) {
	var arr []string = []string{"one", "two", "three", "four"}
	var elt string
	var last string
	var rest []string

	for len(arr) > 0 {
		elt = arr[len(arr)-1]
		last, rest = Pop(arr)

		if last != elt {
			t.Errorf("Unexpected '%v'", last)
			return
		}

		if len(rest) != len(arr)-1 {
			t.Errorf(
				"Length of rest does not match: %d != %d",
				len(rest),
				len(arr)-1,
			)
		}

		arr = arr[:len(arr)-1]
	}

	// Now test zero-length array.
	last, rest = Pop([]string{})
	if last != "" {
		t.Errorf("Unexpected '%v'", last)
	}
}

/* pop_test.go ends here. */
