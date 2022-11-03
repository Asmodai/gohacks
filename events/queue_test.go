/*
 * eventqueue_test.go --- Event queue tests.
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

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queue := NewQueue()

	t.Run("Constructor", func(t *testing.T) {
		if queue.Events() != 0 {
			t.Errorf("Unexpected events in the queue: %d", queue.Events())
		}

		if queue.Capacity() != eventQueueInitialCapacity {
			t.Errorf("Unexpected initial capacity: %d", queue.Capacity())
		}
	})

	t.Run("Push to capacity", func(t *testing.T) {
		for i := 0; i < queue.Capacity(); i++ {
			evt := NewMessage(i, "Nope")

			queue.Push(evt)
		}

		if queue.Events() != queue.Capacity() {
			t.Errorf(
				"Wrong number of events: %d (%d)",
				queue.Events(),
				queue.Capacity(),
			)
		}
	})

	t.Run("Push beyond capacity", func(t *testing.T) {
		curcap := queue.Capacity()
		curevts := queue.Events()
		newevts := 200

		for i := 0; i < newevts; i++ {
			evt := NewMessage(i, "Yup")

			queue.Push(evt)
		}

		if queue.Events() < curevts+newevts {
			t.Errorf(
				"Wrong number of events want:%d != got:%d [%d]",
				curevts+newevts,
				queue.Events(),
				queue.Capacity(),
			)
		}

		if queue.Capacity() < curcap+newevts {
			t.Errorf(
				"Capacity did not increase: %d < %d",
				queue.Capacity(),
				curcap+newevts,
			)
		}
	})

	t.Run("Pop", func(t *testing.T) {
		var evt Event

		evt = queue.Pop()
		for evt != nil {
			evt = queue.Pop()
		}

		if queue.Events() != 0 {
			t.Errorf("Events remaining in queue: %d", queue.Events())
		}
	})
}

/* eventqueue_test.go ends here. */
