// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// respondable.go --- Interface for `Respondable` things.
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
//
//mock:yes

// * Comments:

//
//
//

// * Package:

package responder

// * Imports:

import (
	"github.com/Asmodai/gohacks/v1/events"
)

// * Code:

// ** Interfaces:

// Objects the implement these methods are considered `respondable` and are
// deemed capable of being sent messages directly or via responder chain.
type Respondable interface {
	// A unique name for the respondable object.
	//
	// As this allows us to send events to a specific thing the value
	// returned here must be unique.
	Name() string

	// The type of the respondable object.
	//
	// This can be the internal Go type, or some arbitrary user-specified
	// value that makes sense to you.
	//
	// This is used to implement a "send to all of type" system.
	Type() string

	// Does the receiver respond to a specific event or event type?
	//
	// There is no definition for what `RespondsTo` should do other than
	// return a boolean that states whether an object responds to an
	// event or not.
	RespondsTo(events.Event) bool

	// Send an event to the object.
	//
	// There is no second return value to indicate success or whether
	// the event was handled or not.  The idea being that the receiver
	// will send an `events.Response` event back.
	Invoke(events.Event) events.Event
}

// * respondable.go ends here.
