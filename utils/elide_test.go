/*
 * elide_test.go --- `elide` tests.
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
