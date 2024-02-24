/*
 * bytes_test.go --- Byte conversion testing.
 *
 * Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
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
