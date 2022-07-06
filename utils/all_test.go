/*
 * all_test.go --- `All` tests.
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
