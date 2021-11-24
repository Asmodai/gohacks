/*
 * queue_test.go --- Queue tests.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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
