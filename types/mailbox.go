// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// mailbox.go --- Cheap mailbox data type.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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
)

// * Code:

// ** Types:

type Datum = any

/*
Mailbox structure.

This is a cheap implementation of a mailbox.

It uses two semaphores to control read and write access, and contains
a single datum.
*/
type Mailbox struct {
	ch     chan Datum
	mu     sync.Mutex
	closed bool
}

// ** Methods:

// Put an element into the mailbox.
func (m *Mailbox) Put(elem Datum) bool {
	m.mu.Lock()
	// CRITICAL SECTION START.
	{
		if m.closed {
			m.mu.Unlock() // Exit critical section.

			return false
		}
	}
	// CRITICAL SECTION END.
	m.mu.Unlock()

	m.ch <- elem

	return true
}

// Put an element into the mailbox using a context.
func (m *Mailbox) PutWithContext(ctx context.Context, elem Datum) bool {
	m.mu.Lock()
	// CRITICAL SECTION START.
	{
		if m.closed {
			m.mu.Unlock() // Exit critical section.

			return false
		}
	}
	// CRITICAL SECTION END.
	m.mu.Unlock()

	select {
	case m.ch <- elem:
		return true

	case <-ctx.Done():
		return false
	}
}

// Try to put an element into the mailbox.
func (m *Mailbox) TryPut(elem Datum) bool {
	m.mu.Lock()
	// CRITICAL SECTION START.
	{
		if m.closed {
			m.mu.Unlock() // Exit critical section.

			return false
		}
	}
	// CRITICAL SECTION END.
	m.mu.Unlock()

	select {
	case m.ch <- elem:
		return true

	default:
		return false
	}
}

// Get an element from the mailbox.  Defaults to using a context with
// a deadline of 5 seconds.
func (m *Mailbox) Get() (Datum, bool) {
	value, ok := <-m.ch

	return value, ok
}

// Get an element from the mailbox using the provided context.
//
// It is recommended to use a context that has a timeout deadline.
func (m *Mailbox) GetWithContext(ctx context.Context) (Datum, bool) {
	select {
	case val, ok := <-m.ch:
		return val, ok

	case <-ctx.Done():
		return nil, false
	}
}

// Try to get an element from the mailbox.
func (m *Mailbox) TryGet() (Datum, bool) {
	select {
	case val, ok := <-m.ch:
		return val, ok

	default:
		return nil, false
	}
}

// Does the mailbox contain a value?
func (m *Mailbox) Full() bool {
	return len(m.ch) == 1
}

// Is the mailbox empty like my heart?
func (m *Mailbox) Empty() bool {
	return len(m.ch) == 0
}

// ** Functions:

// Create and return a new empty mailbox.
func NewMailbox() *Mailbox {
	return &Mailbox{ch: make(chan Datum, 1)}
}

// * mailbox.go ends here.
