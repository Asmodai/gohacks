/*
 * message_test.go --- Message event tests.
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

package events

import "testing"

func TestMessageEvent(t *testing.T) {
	command := 42
	message := "Data message"

	evt := NewMessage(command, message)

	t.Run("Accessors", func(t *testing.T) {
		if evt.Command() != command {
			t.Errorf("Unexpected command: '%v'", evt.Command())
		}

		if evt.Data() != message {
			t.Errorf("Unexpected data: '%v'", evt.Data())
		}

		if evt.String() != "Message Event: index:2" {
			t.Errorf("Unexpected string: '%s'", evt.String())
		}
	})

	t.Run("Counter", func(t *testing.T) {
		var e *Message

		for i := 1; i < 2000; i++ {
			e = NewMessage(1, "test")
		}

		if e.Index() != uint64(2001) {
			t.Errorf("Unexpected index: %d", e.Index())
		}
	})
}

/* message_test.go ends here. */
