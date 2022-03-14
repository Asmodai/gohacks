/*
 * elide_test.go --- `elide` tests.
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

func ElideTest(t *testing.T) {
	length := 8
	one := "1234"
	two := "123456789ABCD"

	t.Run("Does not elide when string is within limit", func(t *testing.T) {
		res := Elidable(one).Elide(length)

		if res != "1234" {
			t.Errorf("Unexpected result '%s'", res)
		}
	})

	t.Run("Elide when string is longer than limit", func(t *testing.T) {
		res := Elidable(two).Elide(length)

		if res != "12345..." {
			t.Errorf("Unexpected result '%s'", res)
		}
	})
}

/* elide_test.go ends here. */
