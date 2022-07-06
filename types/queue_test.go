/*
 * queue_test.go --- Queue tests.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
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

package types

import (
	"testing"
)

// Unbounded queue tests.
func TestUnboundedQueue(t *testing.T) {
	var queue *Queue = nil
	var elems []string = []string{"Test1", "Test2", "Test3"}

	//
	// Test creation.
	t.Run("Can create a new unbounded queue", func(t *testing.T) {
		queue = NewQueue()

		if queue.Full() != false || queue.Len() != 0 {
			t.Error("Something went wrong, queue has unexpected settings.")
			return
		}
	})

	//
	// Test item appending.
	t.Run("Can put items to the queue", func(t *testing.T) {
		for _, elt := range elems {
			if ok := queue.Put(elt); !ok {
				t.Error("Queue is unexpectedly full!")
				return
			}
		}
	})

	// Test item getting.
	t.Run("Can get items from the queue", func(t *testing.T) {
		if queue.Len() == 0 {
			t.Error("Queue has a length of zero!")
			return
		}

		if queue.Len() > len(elems) {
			t.Error("Queue has more elements than source array!")
			return
		}

		for idx := range elems {
			res, ok := queue.Get()

			if !ok {
				t.Error("Queue unexpectedly has a length of zero!")
				return
			}

			if res.(string) != elems[idx] {
				t.Errorf("Result mismatch: '%s' != '%s'", res.(string), elems[idx])
				return
			}
		}

		if queue.Len() > 0 {
			t.Error("Queue unexpectedly still contains items!")
			return
		}
	})

	t.Run("Get operation on empty queue returns false value", func(t *testing.T) {
		if _, ok := queue.Get(); ok {
			t.Error("Somehow we got a valid item from an empty queue!")
		}
	})
}

// Bounded queue tests.
func TestBoundedQueue(t *testing.T) {
	var queue *Queue = nil
	var elems1 []string = []string{"Test1", "Test2", "Test3"}
	var elems2 []string = []string{"Test4", "Test5", "Test6"}

	//
	// Test creation.
	t.Run("Can create a new bounded queue", func(t *testing.T) {
		nlen := len(elems1)

		queue = NewBoundedQueue(nlen)

		if queue.Full() != false || queue.Len() != 0 {
			t.Error("Something went wrong, queue has unexpected settings.")
			return
		}
	})

	//
	// Test item appending.
	t.Run("Can put items to the queue", func(t *testing.T) {
		for _, elt := range elems1 {

			if ok := queue.Put(elt); !ok {
				t.Error("Queue is unexpectedly full!")
				return
			}
		}
	})

	//
	// Test if any further puts fail or not.
	t.Run("Puts refused at max capacity", func(t *testing.T) {
		ok := queue.Put("Fail")
		if ok {
			t.Error("Able to `put` beyond bounds.")
			return
		}
	})

	//
	// Test item getting.
	t.Run("Can get items from the queue", func(t *testing.T) {
		if queue.Len() == 0 {
			t.Error("Queue has a length of zero!")
			return
		}

		if queue.Len() > len(elems1) {
			t.Error("Queue has more elements than source array!")
			return
		}

		for idx := range elems1 {
			res, ok := queue.Get()

			if !ok {
				t.Error("Queue unexpectedly has a length of zero!")
				return
			}

			if res.(string) != elems1[idx] {
				t.Errorf("Result mismatch: '%s' != '%s'", res.(string), elems1[idx])
				return
			}
		}

		if queue.Len() > 0 {
			t.Error("Queue unexpectedly still contains items!")
			return
		}
	})

	//
	// Test item appending once bounds are cleared.
	t.Run("Can put further items to the queue", func(t *testing.T) {
		for _, elt := range elems2 {

			if ok := queue.Put(elt); !ok {
				t.Error("Queue is unexpectedly full!")
				return
			}
		}
	})

	//
	// Test further item getting.
	t.Run("Can get items from the queue", func(t *testing.T) {
		if queue.Len() == 0 {
			t.Error("Queue has a length of zero!")
			return
		}

		if queue.Len() > len(elems2) {
			t.Error("Queue has more elements than source array!")
			return
		}

		for idx := range elems2 {
			res, ok := queue.Get()

			if !ok {
				t.Error("Queue unexpectedly has a length of zero!")
				return
			}

			if res.(string) != elems2[idx] {
				t.Errorf("Result mismatch: '%s' != '%s'", res.(string), elems1[idx])
				return
			}
		}

		if queue.Len() > 0 {
			t.Error("Queue unexpectedly still contains items!")
			return
		}
	})
}

/* queue_test.go ends here. */
