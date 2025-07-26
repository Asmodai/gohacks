// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// mailbox.go --- Cheap mailbox data type.
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

//
//
//

// * Package:

package types

// * Imports:

import (
	"golang.org/x/sync/semaphore"

	"context"
	"time"
)

// * Constants:

const (
	// Amount of time to delay semaphore acquisition loops.
	MailboxDelaySleep time.Duration = 50 * time.Millisecond

	// Default deadline for context timeouts.
	DefaultCtxDeadline time.Duration = 5 * time.Second
)

// * Code:

// ** Types:

/*
Mailbox structure.

This is a cheap implementation of a mailbox.

It uses two semaphores to control read and write access, and contains
a single datum.

This is *not* a queue!
*/
type Mailbox struct {
	element any

	// The `preventWrite` semaphore, when acquired, will prevent writes.
	preventWrite *semaphore.Weighted

	// The `preventRead` semaphore, when acquired, will prevent reads.
	preventRead *semaphore.Weighted
}

// ** Methods:

// Put an element into the mailbox.
func (m *Mailbox) Put(elem any) {
	// Attempt to acquire `preventWrite` semaphore.
	for !m.preventWrite.TryAcquire(1) {
		time.Sleep(MailboxDelaySleep)
	}

	// Semaphore acquired, put elem on queue.
	m.element = elem

	// We're no longer full, so release the `preventRead` semaphore.
	m.preventRead.Release(1)
}

// Get an element from the mailbox.  Defaults to using a context with
// a deadline of 5 seconds.
func (m *Mailbox) Get() (any, bool) {
	ctx, cancel := context.WithTimeout(
		context.TODO(),
		DefaultCtxDeadline,
	)

	// Blocks.
	val, ok := m.GetWithContext(ctx)

	// Cancel the context.
	cancel()

	return val, ok
}

// Get an element from the mailbox using the provided context.
//
// It is recommended to use a context that has a timeout deadline.
func (m *Mailbox) GetWithContext(ctx context.Context) (any, bool) {
	var result any

	if m.element == nil {
		return nil, false
	}

	// Attempt to acquire the `preventRead` semaphore
	for !m.preventRead.TryAcquire(1) {
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

// ** Functions:

// Create and return a new empty mailbox.
//
// Note: this acquires the `preventRead` semaphore.
//
//nolint:errcheck
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

// * mailbox.go ends here.
