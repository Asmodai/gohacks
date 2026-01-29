// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// event_error.go --- Selector error event.
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

// * Comments:

// * Package:

package selector

// * Imports:

import (
	"time"

	"github.com/Asmodai/gohacks/events"
)

// * Code:

// Selector error event.
//
// NOTE: `golangci-lint` will want this to be called `Error', and that is
// not what we want.  This is an explicit event, not to be confused with
// `events.Error`.
//
//nolint:revive
type SelectorError struct {
	err events.Error
}

func (e *SelectorError) When() time.Time  { return e.err.Time.When() }
func (e *SelectorError) String() string   { return e.err.Err.Error() }
func (e *SelectorError) Error() error     { return e.err.Err }
func (e *SelectorError) Selector() string { return "error" }

func NewSelectorError(err error) *SelectorError {
	return &SelectorError{
		err: events.Error{
			Time: events.Time{TStamp: time.Now()},
			Err:  err,
		},
	}
}

// * event_error.go ends here.
