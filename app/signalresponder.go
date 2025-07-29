// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// signalresponder.go --- Signal responder type.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
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

package app

// * Imports:

import (
	"os"
	"sync"

	"github.com/Asmodai/gohacks/events"
)

// * Constants:

const (
	signalResponderName string = "signal_responder"
	signalResponderType string = "app.SignalResponder"
)

// * Variables:

// * Code:

// ** Types:

// Callback for the signal responder.
type OnSignalFn func(os.Signal)

// Signal responder.
type SignalResponder struct {
	mu       sync.RWMutex
	callback OnSignalFn
}

// *** Methods

// Returns the name of the responder.
func (sr *SignalResponder) Name() string {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	return signalResponderName
}

// Returns the type of the responder.
func (sr *SignalResponder) Type() string {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	return signalResponderType
}

// Returns whether the responder can respond to a given event.
func (sr *SignalResponder) RespondsTo(evt events.Event) bool {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	switch evt.(type) {
	case *events.Signal:
		return true

	default:
		return false
	}
}

// Invokes the given event.
func (sr *SignalResponder) Invoke(evt events.Event) events.Event {
	sr.mu.RLock()
	defer sr.mu.RUnlock()

	sigevt, ok := evt.(*events.Signal)
	if !ok {
		return evt
	}

	sr.callback(sigevt.Signal())

	return evt
}

// Sets the callback function.
func (sr *SignalResponder) SetOnSignal(callback OnSignalFn) {
	sr.mu.Lock()
	defer sr.mu.Unlock()

	sr.callback = callback
}

// *** Functions:

func defaultSignalResponderCallback(_ os.Signal) {
}

func NewSignalResponder() *SignalResponder {
	return &SignalResponder{
		callback: defaultSignalResponderCallback,
	}
}

// * signalresponder.go ends here.
