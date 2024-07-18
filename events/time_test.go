// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// time_test.go --- Time event tests.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package events

import (
	"testing"
	"time"
)

func TestTimeEvent(t *testing.T) {
	evt := NewTime()
	zero := time.Time{}
	now := time.Now()

	s := "Time Event: " + now.Format(time.RFC3339)

	t.Run("Constructor", func(t *testing.T) {
		if evt.When() != zero {
			t.Errorf("Unexpected default time: %#v", evt.When())
		}
	})

	t.Run("SetNow", func(t *testing.T) {
		evt.SetNow()

		if evt.When() == zero {
			t.Errorf("%#v == %#v", evt.When(), zero)
		}
	})

	t.Run("SetWhen", func(t *testing.T) {
		when := time.Time{}.Add(5 * time.Hour)
		evt.SetWhen(when)

		if evt.When() != when {
			t.Errorf("%#v != %#v", evt.When(), when)
		}
	})

	t.Run("String", func(t *testing.T) {
		evt.SetWhen(now)

		if evt.String() != s {
			t.Errorf("Unexpected string: '%s'", evt.String())
		}
	})
}

// time_test.go ends here.
