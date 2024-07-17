/*
 * queue.go --- Simple queue.
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

package types

import (
	"sync"
)

/*
Queue structure.

This is a cheap implementation of a LIFO queue.
*/
type Queue struct {
	sync.Mutex

	queue   []interface{}
	bounds  int
	bounded bool
}

// Create a new empty queue.
func NewQueue() *Queue {
	return NewBoundedQueue(0)
}

func NewBoundedQueue(bounds int) *Queue {
	queue := &Queue{
		queue:  make([]interface{}, 0),
		bounds: 0,
	}

	if bounds == 0 {
		queue.bounded = false

		return queue
	}

	queue.bounded = true
	queue.bounds = bounds

	return queue
}

// Append an element to the queue.  Returns `false` if there is no
// more room in the queue.
func (q *Queue) Put(elem interface{}) bool {
	if q.bounded && q.Len() == q.bounds {
		return false
	}

	q.Lock()
	{
		q.queue = append(q.queue, elem)
	}
	q.Unlock()

	return true
}

// Remove an element from the end of the queue and return it.
func (q *Queue) Get() (interface{}, bool) {
	var elem interface{}

	if q.Len() == 0 {
		return nil, false
	}

	q.Lock()
	{
		elem = q.queue[0]

		q.queue[0] = nil
		q.queue = q.queue[1:]
	}
	q.Unlock()

	return elem, true
}

// Return the number of elements in the queue.
func (q *Queue) Len() int {
	var length int

	q.Lock()
	{
		length = len(q.queue)
	}
	q.Unlock()

	return length
}

// Is the queue full?
func (q *Queue) Full() bool {
	return q.bounded && q.Len() == q.bounds
}

/* queue.go ends here. */
