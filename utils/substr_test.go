/*
 * substr_test.go --- Substring tests.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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

package utils

import "testing"

func TestSubstr(t *testing.T) {
	var input string = "One two three"

	t.Run("Works with sane arguments", func(t *testing.T) {
		if output := Substr(input, 4, 3); output != "two" {
			t.Errorf("Unexpected '%v'", output)
		}
	})

	t.Run("Works when length is longer than input", func(t *testing.T) {
		if output := Substr(input, 4, 100); output != "two three" {
			t.Errorf("Unexpected '%v'", output)
		}
	})

	t.Run("Returns empty string with insane args", func(t *testing.T) {
		if output := Substr(input, 100, 24); output != "" {
			t.Errorf("Unexpected '%v'", output)
		}
	})
}

/* substr_test.go ends here. */
