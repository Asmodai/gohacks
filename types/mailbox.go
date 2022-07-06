/*
 * mailbox.go --- Cheap mailbox data type.
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
	"golang.org/x/sync/semaphore"

	"context"
	"time"
)

const (
	// Amount of time to delay semaphore acquisition loops.
	MailboxDelaySleep  time.Duration = 50 * time.Millisecond
	DefaultCtxDeadline time.Duration = 5 * time.Second
)

/*

Mailbox structure.

This is a cheap implementation of a mailbox.

It uses two semaphores to control read and write access, and contains
a single datum.

This is *not* a queue!

*/
type Mailbox struct {
	element interface{}

	// The `preventWrite` semaphore, when acquired, will prevent writes.
	// The `preventRead` semaphore, when acquired, will prevent reads.
	preventWrite *semaphore.Weighted
	preventRead  *semaphore.Weighted
}

// Create and return a new empty mailbox.
//
// Note: this acquires the `preventRead` semaphore.
func NewMailbox() *Mailbox {
	// Please note that the context given here should never be one
	// passed in by the user, we want a TODO context because *we* are
	// setting up this initial context.
	preventRead := semaphore.NewWeighted(int64(1))
	preventRead.Acquire(context.TODO(), 1)

	return &Mailbox{
		element:      nil,
		preventWrite: semaphore.NewWeighted(int64(1)),
		preventRead:  preventRead,
	}
}

// Put an element into the mailbox.
func (m *Mailbox) Put(elem interface{}) {
	// Attempt to acquire `preventWrite` semaphore.
	for m.preventWrite.TryAcquire(1) == false {
		time.Sleep(MailboxDelaySleep)
	}

	// Semaphore acquired, put elem on queue.
	m.element = elem

	// We're no longer full, so release the `preventRead` semaphore.
	m.preventRead.Release(1)
}

// Get an element from the mailbox.  Defaults to using a context with
// a deadline of 5 seconds.
func (m *Mailbox) Get() (interface{}, bool) {
	ctx, cancel := context.WithTimeout(
		context.TODO(),
		DefaultCtxDeadline,
	)

	val, ok := m.GetWithContext(ctx)
	cancel()

	return val, ok
}

func (m *Mailbox) GetWithContext(ctx context.Context) (interface{}, bool) {
	var result interface{} = nil

	if m.element == nil {
		return nil, false
	}

	// Attempt to acquire the `preventRead` semaphore
	for m.preventRead.TryAcquire(1) == false {
		time.Sleep(MailboxDelaySleep)

		select {
		case <-ctx.Done():
			return nil, false
		default:
		}
	}

	result = m.element
	m.element = nil

	// Release the `preventWrite` semaphore.
	m.preventWrite.Release(1)

	return result, true
}

// Does the mailbox contain a value?
func (m *Mailbox) Full() bool {
	return m.element != nil
}

/* mailbox.go ends here. */
