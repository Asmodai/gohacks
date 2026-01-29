// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// queue.go --- Event queue structure.
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

package events

import (
	"sync"
)

const (
	eventQueueInitialCapacity = 256
	eventQueueIncrement       = 128
)

type EventList []Event

type Queue struct {
	sync.Mutex

	queue EventList
}

func NewQueue() *Queue {
	return &Queue{
		queue: make(EventList, 0, eventQueueInitialCapacity),
	}
}

func (e *Queue) Events() int   { return len(e.queue) }
func (e *Queue) Capacity() int { return cap(e.queue) }

func (e *Queue) Pop() Event {
	e.Lock()
	defer e.Unlock()

	switch e.Events() {
	case 0:
		return nil

	case 1:
		ret := e.queue[0]
		e.queue = make(EventList, 0, eventQueueInitialCapacity)

		return ret

	default:
		ret := e.queue[0]
		e.queue = e.queue[1:]

		return ret
	}
}

func (e *Queue) Push(evt Event) {
	e.Lock()
	defer e.Unlock()

	if len(e.queue) == cap(e.queue) {
		nqueue := make(
			EventList,
			len(e.queue),
			cap(e.queue)+eventQueueIncrement,
		)
		copy(nqueue, e.queue)
		e.queue = nqueue
	}

	e.queue = append(e.queue, evt)
}

// queue.go ends here.
