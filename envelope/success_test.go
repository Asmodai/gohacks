// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// success_test.go --- Success envelope tests.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * Comments:

// * Package:

package envelope

// * Imports:

import (
	"encoding/json"
	"testing"
	"time"
)

// * Constants:

// * Variables:

// * Code:

// ** Tests:

func TestSuccessBasic(t *testing.T) {
	code := 418
	body := []string{"I am not a coffee pot", "I am just a lowly teapot"}

	env := NewSuccess(code, body)
	env.Elapsed = time.Duration(15 * time.Millisecond)

	t.Run("Status", func(t *testing.T) {
		if env.Status() != code {
			t.Fatalf("unexpected code: %#v != %#v",
				env.Status(),
				code)
		}
	})

	t.Run("Body", func(t *testing.T) {
		if env.Body() != env {
			t.Fatalf("unexpected body: %#v", env.Body())
		}
	})

	t.Run("JSON", func(t *testing.T) {
		want := `{
.."data": [
...."I am not a coffee pot",
...."I am just a lowly teapot"
..],
.."count": 2,
.."elapsed": "15ms"
}`

		res, err := json.MarshalIndent(env.Body(), "", "..")
		if err != nil {
			t.Fatalf("json marshal error: %#v", err.Error())
		}

		sres := string(res)
		if sres != want {
			t.Fatalf("invalid json: %v", sres)
		}
	})

}

// * success_test.go ends here.
