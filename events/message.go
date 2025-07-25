// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// message.go --- Message events.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

//nolint:gochecknoglobals
var counter uint64

type Message struct {
	Time

	index   uint64
	command int
	data    any
}

func updateCounter() {
	if atomic.LoadUint64(&counter) == math.MaxUint64 {
		atomic.StoreUint64(&counter, 0)
	}

	atomic.AddUint64(&counter, 1)
}

func NewMessage(cmd int, data any) *Message {
	updateCounter()

	return &Message{
		Time: Time{
			TStamp: time.Now(),
		},
		index:   atomic.LoadUint64(&counter),
		command: cmd,
		data:    data,
	}
}

func (e *Message) Index() uint64 { return e.index }
func (e *Message) Command() int  { return e.command }
func (e *Message) Data() any     { return e.data }

func (e *Message) String() string {
	return fmt.Sprintf("Message Event: index:%d", e.index)
}

// message.go ends here.
