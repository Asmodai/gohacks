// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// event_response.go --- Selector response event.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

// * Constants:

// * Variables:

// * Code:
// ** Type:

// Selector response event.
//
// NOTE: `golangci-lint` will want this to be called `Response', and that is
// not what we want.  This is an explicit event, not to be confused with
// `events.Response`.
//
//nolint:revive
type SelectorResponse struct {
	events.Time

	received time.Time
	selector string
	response any
}

// ** Methods:

func (r *SelectorResponse) When() time.Time  { return r.Time.TStamp }
func (r *SelectorResponse) Response() any    { return r.response }
func (r *SelectorResponse) Selector() string { return r.selector }

func (r *SelectorResponse) String() string {
	return r.selector + " response"
}

// ** Functions:

func NewSelectorResponse(sel SelectorEvent, data any) SelectorEvent {
	return &SelectorResponse{
		Time:     events.Time{TStamp: time.Now()},
		selector: sel.Selector(),
		received: sel.When(),
		response: data,
	}
}

// * event_response.go ends here.
