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
	"context"
	"sync"

	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Types:

/*
Queue structure.

This is a cheap implementation of a FIFO queue.
*/
type Queue struct {
	sync.Mutex

	notEmpty *sync.Cond // Queue is not empty.
	notFull  *sync.Cond // Queue is not full.

	queue   []Datum
	bounds  int
	bounded bool
}

// ** Methods:

// Validate that the queue has been initialised.
//
// Panics if the conds aren't initialised.
func (q *Queue) validate() {
	if q.notEmpty == nil || q.notFull == nil {
		panic("Attempt to use uninitialised queue")
	}
}

// Unlock the main queue mutex while waiting on cond.Wait().
//
// Locks the mutex when done.
// A goroutine is used to avoid deadlocks.
// If the context is cancelled, it will exit and lock the mutex on the way
// out.
func (q *Queue) waitWithContextUnlocked(ctx context.Context, cond *sync.Cond) error {
	q.validate()

	done := make(chan struct{})

	go func() {
		q.Unlock()
		{
			// Notice me, senpai!
			// Notice how we're not locked in here!
			// Notice how I have a key... to your heart.
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
		}
		q.Lock()
		close(done)
	}()

	select {
	case <-ctx.Done():
		// Relock if the context is done.
		q.Lock()

		return errors.WithStack(ctx.Err())
	case <-done:
		return nil
	}
}

// Put an element on to the queue.
//
// This blocks.
func (q *Queue) Put(elem Datum) {
	q.validate()

	if q.notFull == nil {
		panic("Attempt to use uninitialised queue")
	}

	q.Lock()
	defer q.Unlock()

	for q.bounded && len(q.queue) == q.bounds {
		q.notFull.Wait()
	}

	q.queue = append(q.queue, elem)
	q.notEmpty.Signal()
}

// Put an element on to the queue.
//
// Will exit should the context time out or be cancelled.
//
// This blocks.
func (q *Queue) PutWithContext(ctx context.Context, elem Datum) error {
	q.validate()

	q.Lock()
	defer q.Unlock()

	for q.bounded && len(q.queue) == q.bounds {
		if err := q.waitWithContextUnlocked(ctx, q.notFull); err != nil {
			return errors.WithStack(err)
		}
	}

	q.queue = append(q.queue, elem)
	q.notEmpty.Signal()

	return nil
}

// Append an element to the queue.
//
// Returns `false` if there is no more room in the queue.
//
// This does not block.
func (q *Queue) PutWithoutBlock(elem Datum) bool {
	q.validate()

	q.Lock()
	defer q.Unlock()

	// If you're new to synchronisation:
	// This is the same as `Full`, yes... but doesn't allow races...
	// so leave this alone.
	if q.bounded && len(q.queue) == q.bounds {
		return false
	}

	q.queue = append(q.queue, elem)
	q.notEmpty.Signal()

	return true
}

// Remove an element from the start of the queue and return it.
//
// This blocks.
func (q *Queue) Get() Datum {
	q.validate()

	q.Lock()
	defer q.Unlock()

	for len(q.queue) == 0 {
		q.notEmpty.Wait()
	}

	elem := q.queue[0]
	q.queue[0] = nil
	q.queue = q.queue[1:]
	q.notFull.Signal()

	return elem
}

// Remove an element from the start of the queue and return it.
//
// Will exit should the context time out or be cancelled.
//
// This blocks.
func (q *Queue) GetWithContext(ctx context.Context) (Datum, error) {
	q.validate()

	q.Lock()
	defer q.Unlock()

	for len(q.queue) == 0 {
		if err := q.waitWithContextUnlocked(ctx, q.notEmpty); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	elem := q.queue[0]
	q.queue[0] = nil
	q.queue = q.queue[1:]
	q.notFull.Signal()

	return elem, nil
}

// Remove an element from the start of the queue and return it.
//
// This does not block.
func (q *Queue) GetWithoutBlock() (Datum, bool) {
	q.validate()

	q.Lock()
	defer q.Unlock()

	if len(q.queue) == 0 {
		return nil, false
	}

	elem := q.queue[0]
	q.queue[0] = nil
	q.queue = q.queue[1:]
	q.notFull.Signal()

	return elem, true
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

	return q.bounded && len(q.queue) == q.bounds
}

// Is the queue empty?
func (q *Queue) Empty() bool {
	q.Lock()
	defer q.Unlock()

	return len(q.queue) == 0
}

// ** Functions:

// Create a new empty queue.
func NewQueue() *Queue {
	return NewBoundedQueue(0)
}

// Create a queue that is bounded to a specific size.
func NewBoundedQueue(bounds int) *Queue {
	// Do not accept negative bounds.
	if bounds < 0 {
		bounds = 0
	}

	queue := &Queue{
		queue:   make([]Datum, 0, bounds),
		bounds:  bounds,
		bounded: bounds > 0,
	}

	queue.notEmpty = sync.NewCond(&queue.Mutex)
	queue.notFull = sync.NewCond(&queue.Mutex)

	return queue
}

// * queue.go ends here.
