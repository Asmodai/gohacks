// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// queue.go --- Simple queue.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

package types

// * Imports:

import (
	"sync"
)

// * Code:

// ** Types:

/*
Queue structure.

This is a cheap implementation of a FIFO queue.
*/
type Queue struct {
	sync.Mutex

	queue   []Datum
	bounds  int
	bounded bool
}

// ** Methods:

// Append an element to the queue.  Returns `false` if there is no
// more room in the queue.
func (q *Queue) Put(elem Datum) bool {
	q.Lock()
	defer q.Unlock()

	// If you're new to synchronisation:
	// This is the same as `Full`, yes... but doesn't allow races...
	// so leave this alone.
	if q.bounded && q.helperLen() == q.bounds {
		return false
	}

	q.queue = append(q.queue, elem)

	return true
}

// Remove an element from the end of the queue and return it.
func (q *Queue) Get() (Datum, bool) {
	q.Lock()
	defer q.Unlock()

	if q.helperLen() == 0 {
		return nil, false
	}

	elem := q.queue[0]
	q.queue[0] = nil
	q.queue = q.queue[1:]

	return elem, true
}

// Internal helper that returns the number of elements in the queue.
//
// Should never lock. Callers must ensure lock is held.
func (q *Queue) helperLen() int {
	return len(q.queue)
}

// Return the number of elements in the queue.
func (q *Queue) Len() int {
	q.Lock()
	defer q.Unlock()

	return len(q.queue)
}

// Is the queue full?
func (q *Queue) Full() bool {
	q.Lock()
	defer q.Unlock()

	return q.bounded && q.helperLen() == q.bounds
}

// Is the queue empty?
func (q *Queue) Empty() bool {
	q.Lock()
	defer q.Unlock()

	return q.helperLen() == 0
}

// ** Functions:

// Create a new empty queue.
func NewQueue() *Queue {
	return NewBoundedQueue(0)
}

// Create a queue that is bounded to a specific size.
func NewBoundedQueue(bounds int) *Queue {
	var capacity int

	// Do not accept negative bounds.
	if bounds < 0 {
		bounds = 0
	}

	if bounds > 0 {
		capacity = bounds
	}

	return &Queue{
		queue:   make([]Datum, 0, capacity),
		bounds:  bounds,
		bounded: bounds > 0,
	}
}

// * queue.go ends here.
