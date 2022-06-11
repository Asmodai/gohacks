/*
 * bytes_test.go --- Byte conversion testing.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
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

package conversion

import (
	"testing"
)

func TestKiB(t *testing.T) {
	var b uint64 = 1024
	var kib uint64 = 1

	r1 := BToKiB(b)
	r2 := KiBToB(r1)

	t.Log("Does B -> KiB -> B conversion work?")
	if r1 == kib && r2 == b {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

func TestMiB(t *testing.T) {
	var b uint64 = 1 * 1.049e6
	var mib uint64 = 1

	r1 := BToMiB(b)
	r2 := MiBToB(r1)

	t.Log("Does B -> MiB -> B conversion work?")
	if r1 == mib && r2 == b {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

func TestGiB(t *testing.T) {
	var b uint64 = 1 * 1.074e9
	var gib uint64 = 1

	r1 := BToGiB(b)
	r2 := GiBToB(r1)

	t.Log("Does B -> GiB -> B conversion work?")
	if r1 == gib && r2 == b {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

func TestTiB(t *testing.T) {
	var b uint64 = 1 * 1.1e12
	var tib uint64 = 1

	r1 := BToTiB(b)
	r2 := TiBToB(r1)

	t.Log("Does B -> TiB -> B conversion work?")
	if r1 == tib && r2 == b {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

/* bytes_test.go ends here. */
