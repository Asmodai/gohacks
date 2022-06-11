/*
 * queue.go --- Simple queue.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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
	n := &Queue{
		queue:  make([]interface{}, 0),
		bounds: 0,
	}

	if bounds == 0 {
		n.bounded = false

		return n
	}

	n.bounded = true
	n.bounds = bounds

	return n
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
	var length int = 0

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
