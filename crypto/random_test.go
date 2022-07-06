/*
 * random_test.go --- RNG tests.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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

		if bytes.Equal(r1, r2) {
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
