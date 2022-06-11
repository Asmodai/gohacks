/*
 * params_test.go --- Query params tests.
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

package apiclient

import "testing"

func TestAccessors(t *testing.T) {
	p := NewParams()

	t.Run("SetUseBasic works", func(t *testing.T) {
		p.SetUseBasic(true)

		if !p.UseBasic {
			t.Error("No, UseBasic is not set!")
		}
	})

	t.Run("SetUseToken works", func(t *testing.T) {
		p.SetUseToken(true)

		if !p.UseToken {
			t.Error("No, UseToken is not set!")
		}
	})

	t.Run("AddQueryParam works", func(t *testing.T) {
		p.AddQueryParam("drink", "yes please")

		if len(p.Queries) != 1 {
			t.Errorf("No, length of queries is %v", len(p.Queries))
		}

		if p.Queries[0].Name != "drink" {
			t.Error("Query name does not match.")
		}

		if p.Queries[0].Content != "yes please" {
			t.Error("Query content does not match.")
		}
	})

	t.Run("ClearQueryParams works", func(t *testing.T) {
		before := len(p.Queries)
		p.ClearQueryParams()
		after := len(p.Queries)

		if before == 0 {
			t.Error("Array was empty before the test ran!")
			return
		}

		if before == after {
			t.Error("Nothing was removed!")
		}

		if after != 0 {
			t.Error("Array still contains elements.")
		}
	})
}

/* params_test.go ends here. */
