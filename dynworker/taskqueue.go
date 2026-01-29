// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// taskqueue.go --- Task queue type.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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
//
//mock:yes

// * Comments:

// * Package:

package dynworker

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/types"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

// * Variables:

var (
	ErrChannelClosed = errors.Base("channel closed")
)

// * Code:

// ** Interface:

type TaskQueue interface {
	Put(context.Context, *Task) error
	Get(context.Context) (*Task, error)
	Len() int
}

// ** Channel-based task queue:

// *** Types:

type chanTaskQueue struct {
	ch chan *Task
}

// *** Methods:

func (c *chanTaskQueue) Put(ctx context.Context, task *Task) error {
	select {
	case c.ch <- task:
		return nil

	case <-ctx.Done():
		return errors.WithStack(ctx.Err())
	}
}

func (c *chanTaskQueue) Get(ctx context.Context) (*Task, error) {
	select {
	case t, ok := <-c.ch:
		if !ok {
			return nil, errors.WithStack(ErrChannelClosed)
		}

		return t, nil

	case <-ctx.Done():
		return nil, errors.WithStack(ctx.Err())
	}
}

func (c *chanTaskQueue) Len() int {
	return len(c.ch)
}

// *** Functions:

func NewChanTaskQueue(size int) TaskQueue {
	return &chanTaskQueue{
		ch: make(chan *Task, size),
	}
}

// ** Queue-based task queue:

// *** Types:

type queueTaskQueue struct {
	q *types.Queue
}

// *** Methods:

func (q *queueTaskQueue) Put(ctx context.Context, task *Task) error {
	if err := q.q.PutWithContext(ctx, task); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (q *queueTaskQueue) Get(ctx context.Context) (*Task, error) {
	result, err := q.q.GetWithContext(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	task, ok := result.(*Task)
	if !ok {
		return nil, errors.WithStack(ErrNotTask)
	}

	return task, nil
}

func (q *queueTaskQueue) Len() int {
	return q.q.Len()
}

// *** Functions:

func NewQueueTaskQueue(q *types.Queue) TaskQueue {
	return &queueTaskQueue{q: q}
}

// * taskqueue.go ends here.
