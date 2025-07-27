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

// * Comments:

// * Package:

package process

// * Imports:

import (
	"context"
	"sync"

	"github.com/Asmodai/gohacks/v1/events"
	"github.com/Asmodai/gohacks/v1/logger"
	"github.com/Asmodai/gohacks/v1/responder"
)

// * Code:

// ** Types:

// Internal state for processes.
type State struct {
	mu sync.RWMutex

	parent     *Process
	responders *responder.Chain
}

// ** Methods:

func (ps *State) RespondsTo(event events.Event) bool {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.responders.RespondsTo(event)
}

func (ps *State) Invoke(event events.Event) (events.Event, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.responders.SendFirst(event)
}

// Return the context for the parent process.
func (ps *State) Context() context.Context {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.parent.Context()
}

func (ps *State) Logger() logger.Logger {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return ps.parent.logger
}

// ** Functions:

func newState(name string) *State {
	return &State{
		responders: responder.NewChain(name),
	}
}

// * state.go ends here.
