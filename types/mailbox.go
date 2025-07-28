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
	"sync"

	"golang.org/x/sync/semaphore"

	"context"
	"time"
)

// * Constants:

const (
	// Default deadline for context timeouts.
	defaultCtxDeadline time.Duration = 5 * time.Second
)

// * Code:

// ** Types:

type Datum = any

/*
Mailbox structure.

This is a cheap implementation of a mailbox.

It uses two semaphores to control read and write access, and contains
a single datum.

This is *not* a queue!
*/
type Mailbox struct {
	mu      sync.Mutex
	element Datum

	// The `writeAvailable` semaphore, when acquired, will prevent writes.
	writeAvailable *semaphore.Weighted

	// The `readAvailable` semaphore, when acquired, will prevent reads.
	readAvailable *semaphore.Weighted
}

// ** Methods:

// Put an element into the mailbox.
func (m *Mailbox) Put(elem Datum) bool {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultCtxDeadline,
	)
	defer cancel()

	return m.PutWithContext(ctx, elem)
}

// Put an element into the mailbox using a context.
func (m *Mailbox) PutWithContext(ctx context.Context, elem Datum) bool {
	// Attempt to acquire `writeAvailable` semaphore.
	if err := m.writeAvailable.Acquire(ctx, 1); err != nil {
		return false
	}

	// Semaphore acquired, put elem on queue.
	m.mu.Lock()
	m.element = elem
	m.mu.Unlock()

	// We're no longer full, so release the `readAvailable` semaphore.
	m.readAvailable.Release(1)

	return true
}

// Try to put an element into the mailbox.
func (m *Mailbox) TryPut(item Datum) bool {
	if !m.writeAvailable.TryAcquire(1) {
		return false
	}

	m.mu.Lock()
	m.element = item
	m.mu.Unlock()

	m.readAvailable.Release(1)

	return true
}

// Get an element from the mailbox.  Defaults to using a context with
// a deadline of 5 seconds.
func (m *Mailbox) Get() (Datum, bool) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		defaultCtxDeadline,
	)
	defer cancel()

	return m.GetWithContext(ctx)
}

// Get an element from the mailbox using the provided context.
//
// It is recommended to use a context that has a timeout deadline.
func (m *Mailbox) GetWithContext(ctx context.Context) (Datum, bool) {
	// Attempt to acquire the `readAvailable` semaphore
	if err := m.readAvailable.Acquire(ctx, 1); err != nil {
		return nil, false
	}

	m.mu.Lock()
	result := m.element
	m.element = nil
	m.mu.Unlock()

	// Release the `writeAvailable` semaphore.
	m.writeAvailable.Release(1)

	return result, true
}

// Try to get an element from the mailbox.
func (m *Mailbox) TryGet() (Datum, bool) {
	if !m.readAvailable.TryAcquire(1) {
		return nil, false
	}

	m.mu.Lock()
	result := m.element
	m.element = nil
	m.mu.Unlock()

	m.writeAvailable.Release(1)

	return result, true
}

// Does the mailbox contain a value?
func (m *Mailbox) Full() bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.element != nil
}

// Is the mailbox empty like my heart?
func (m *Mailbox) Empty() bool {
	return !m.Full()
}

// Reset the mailbox.
func (m *Mailbox) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.element = nil
	m.writeAvailable = semaphore.NewWeighted(1)
	m.readAvailable = semaphore.NewWeighted(1)

	//nolint:errcheck
	m.readAvailable.Acquire(context.Background(), 1)
}

// ** Functions:

// Create and return a new empty mailbox.
//
// Note: this acquires the `readAvailable` semaphore.
//
//nolint:errcheck
func NewMailbox() *Mailbox {
	// Please note that the context given here should never be one
	// passed in by the user, we want a TODO context because *we* are
	// setting up this initial context.
	readAvailable := semaphore.NewWeighted(int64(1))
	readAvailable.Acquire(context.TODO(), 1)

	return &Mailbox{
		element:        nil,
		writeAvailable: semaphore.NewWeighted(int64(1)),
		readAvailable:  readAvailable,
	}
}

// * mailbox.go ends here.
