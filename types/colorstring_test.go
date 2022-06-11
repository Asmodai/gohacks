/*
 * colorstring_test.go --- Colour string tests.
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

package types

import "testing"

func TestColorString(t *testing.T) {
	onlyFG := "\x1b[0;31;49mTest\x1b[0m"
	onlyBG := "\x1b[0;39;41mTest\x1b[0m"
	onlyAttr := "\x1b[1;39;49mTest\x1b[0m"

	t.Run("Sets FG properly", func(t *testing.T) {
		cl := MakeColorString()

		cl.SetString("Test")
		cl.SetFG(RED)

		if cl.String() != onlyFG {
			t.Errorf("Unexpected '%s'", cl.String())
		}
	})

	t.Run("Sets BG properly", func(t *testing.T) {
		cl := MakeColorString()

		cl.SetString("Test")
		cl.SetBG(RED)

		if cl.String() != onlyBG {
			t.Errorf("Unexpected '%s'", cl.String())
		}
	})

	t.Run("Sets attribute properly", func(t *testing.T) {
		cl := MakeColorString()

		cl.SetString("Test")
		cl.SetAttr(BOLD)

		if cl.String() != onlyAttr {
			t.Errorf("Unexpected '%s'", cl.String())
		}
	})
}

/* colorstring_test.go ends here. */
