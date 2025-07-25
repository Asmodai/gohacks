// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// state.go --- Internal process state.
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

package process

import (
	"github.com/Asmodai/gohacks/logger"

	"context"
)

// Internal state for processes.
type State struct {
	parent *Process
}

func NewState() *State {
	return &State{}
}

// Return the context for the parent process.
func (ps *State) Context() context.Context {
	return ps.parent.Context()
}

// Send data from a process to an external entity.
func (ps *State) Send(data any) bool {
	select {
	case ps.parent.chanFromState <- data:
		return true

	default:
	}

	return false
}

// Send data from a process to an external entity with blocking.
func (ps *State) SendBlocking(data any) {
	ps.parent.chanFromState <- data
}

// Read data from an external entity.
func (ps *State) Receive() (any, bool) {
	select {
	case data := <-ps.parent.chanToState:
		return data, true

	default:
	}

	return nil, false
}

// Read data from an external entity with blocking.
func (ps *State) ReceiveBlocking() any {
	return <-ps.parent.chanToState
}

func (ps *State) Logger() logger.Logger {
	return ps.parent.logger
}

// state.go ends here.
