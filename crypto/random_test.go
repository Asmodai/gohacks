/*
 * random_test.go --- RNG tests.
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

package crypto

import (
	"bytes"
	"testing"
)

func TestRandomBytes(t *testing.T) {
	var r1 []byte = []byte("")
	var r2 []byte = []byte("")
	var err error = nil

	t.Run("Generates 16 bytes", func(t *testing.T) {
		r1, err = GenerateRandomBytes(16)
		if err != nil {
			t.Errorf("No, '%v'", err.Error())
			return
		}

		if len(r1) != 16 {
			t.Errorf("Length mismatch")
		}
	})

	t.Run("Subsequent generations are unique", func(t *testing.T) {
		r2, err = GenerateRandomBytes(16)
		if err != nil {
			t.Errorf("No, '%v'", err.Error())
		}

		if len(r2) != 16 {
			t.Errorf("Length mismatch")
			return
		}

		if bytes.Compare(r1, r2) == 0 {
			t.Errorf("Not unique!")
		}
	})
}

func TestRandomString(t *testing.T) {
	var s1 string = ""
	var s2 string = ""
	var err error = nil

	t.Run("Generates 32 characters", func(t *testing.T) {
		s1, err = GenerateRandomString(32)
		if err != nil {
			t.Errorf("No, '%v'", err.Error())
			return
		}

		if len(s1) != 32 {
			t.Errorf("Length mismatch")
		}
	})

	t.Run("Subsequent generations are unique", func(t *testing.T) {
		s2, err = GenerateRandomString(32)
		if err != nil {
			t.Errorf("No, '%v'", err.Error())
		}

		if len(s2) != 32 {
			t.Errorf("Length mismatch")
			return
		}

		if s1 == s2 {
			t.Errorf("Not unique!")
		}
	})
}

/* random_test.go ends here. */
