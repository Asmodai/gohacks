/*
 * all_test.go --- `All` tests.
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

func TestIntAll(t *testing.T) {
	goodArr := []int{2, 4, 6, 8, 10}
	badArr := []int{2, 3, 6, 8, 10}

	t.Run("Returns true for good array", func(t *testing.T) {
		res := IntAll(goodArr, func(elt int) bool {
			return elt%2 == 0
		})

		if !res {
			t.Error("Unexpected result!")
		}
	})

	t.Run("Returns false for bad array", func(t *testing.T) {
		res := IntAll(badArr, func(elt int) bool {
			return elt%2 == 0
		})

		if res {
			t.Error("Unexpected result!")
		}
	})
}

/* all_test.go ends here. */
