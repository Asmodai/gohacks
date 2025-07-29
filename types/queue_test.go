// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// queue_test.go --- Queue tests.
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
	"context"
	"sync"
	"testing"
	"time"
)

// * Code:

// ** Tests:

type dummyDatum struct {
	val int
}

// satisfy the Datum interface if any
var _ Datum = (*dummyDatum)(nil)

func TestQueueBasicPutGet(t *testing.T) {
	q := NewBoundedQueue(2)
	d1 := &dummyDatum{1}
	d2 := &dummyDatum{2}

	q.Put(d1)
	q.Put(d2)

	if q.Len() != 2 {
		t.Fatalf("expected length 2, got %d", q.Len())
	}

	got := q.Get()
	if got.(*dummyDatum).val != 1 {
		t.Errorf("expected 1, got %v", got)
	}

	got = q.Get()
	if got.(*dummyDatum).val != 2 {
		t.Errorf("expected 2, got %v", got)
	}
}

func TestQueuePutBlocksWhenFull(t *testing.T) {
	q := NewBoundedQueue(1)
	q.Put(&dummyDatum{1})

	done := make(chan struct{})
	go func() {
		// this should block until Get frees space
		q.Put(&dummyDatum{2})
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Put did not block when queue was full")
	case <-time.After(100 * time.Millisecond):
		// expected, still blocked
	}

	// Free up space
	_ = q.Get()

	select {
	case <-done:
		// success, Put unblocked
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Put did not unblock after space freed")
	}
}

func TestQueueGetBlocksWhenEmpty(t *testing.T) {
	q := NewBoundedQueue(1)

	done := make(chan struct{})
	var got Datum

	go func() {
		got = q.Get()
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Get did not block when queue was empty")
	case <-time.After(100 * time.Millisecond):
		// expected, still blocked
	}

	val := &dummyDatum{42}
	q.Put(val)

	select {
	case <-done:
		if got.(*dummyDatum) != val {
			t.Errorf("expected val %v, got %v", val, got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Get did not unblock after Put")
	}
}

func TestQueuePutWithContextCancellation(t *testing.T) {
	q := NewBoundedQueue(1)
	q.Put(&dummyDatum{1})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := q.PutWithContext(ctx, &dummyDatum{2})
	if err == nil {
		t.Fatal("expected timeout error from PutWithContext")
	}
}

func TestQueueGetWithContextCancellation(t *testing.T) {
	q := NewBoundedQueue(1)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := q.GetWithContext(ctx)
	if err == nil {
		t.Fatal("expected timeout error from GetWithContext")
	}
}

func TestQueuePutWithoutBlock(t *testing.T) {
	q := NewBoundedQueue(1)
	ok := q.PutWithoutBlock(&dummyDatum{1})
	if !ok {
		t.Fatal("expected PutWithoutBlock to succeed")
	}

	ok = q.PutWithoutBlock(&dummyDatum{2})
	if ok {
		t.Fatal("expected PutWithoutBlock to fail due to full queue")
	}
}

func TestQueueGetWithoutBlock(t *testing.T) {
	q := NewBoundedQueue(1)
	q.Put(&dummyDatum{1})

	d, ok := q.GetWithoutBlock()
	if !ok {
		t.Fatal("expected GetWithoutBlock to succeed")
	}
	if d.(*dummyDatum).val != 1 {
		t.Errorf("expected value 1, got %v", d)
	}

	d, ok = q.GetWithoutBlock()
	if ok {
		t.Fatal("expected GetWithoutBlock to fail on empty queue")
	}
}

func TestQueueConcurrentPutGet(t *testing.T) {
	q := NewBoundedQueue(100)
	wg := sync.WaitGroup{}
	n := 1000

	// Producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			q.Put(&dummyDatum{i})
		}
	}()

	// Consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < n; i++ {
			d := q.Get()
			if d == nil {
				t.Errorf("got nil datum at index %d", i)
			}
		}
	}()

	wg.Wait()
}

// * queue_test.go ends here.
