// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// response.go --- Response event.
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

package events

import (
	"fmt"
	"time"
)

type Response struct {
	Time

	received time.Time
	index    uint64
	command  string
	response any
}

func NewResponse(msg *Message, rsp any) *Response {
	return &Response{
		Time: Time{
			TStamp: time.Now(),
		},

		received: msg.When(),
		index:    msg.Index(),
		command:  msg.Command(),
		response: rsp,
	}
}

func (e *Response) Received() time.Time { return e.received }
func (e *Response) Index() uint64       { return e.index }
func (e *Response) Command() string     { return e.command }
func (e *Response) Response() any       { return e.response }

func (e *Response) String() string {
	return fmt.Sprintf(
		"Response Event: index:%d duration:%s",
		e.index,
		e.When().Sub(e.Received()).String(),
	)
}

// response.go ends here.
