/*
 * params_test.go --- Query params tests.
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
