// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// chain.go --- Objective-C-like responder chain.
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
//go:generate go run github.com/Asmodai/gohacks/cmd/digen -pattern .
//di:gen basename=ResponderChain key=gohacks/responder@v1 type=*Chain fallback=NewChain("unnamed")

// * Comments:

// This is loosely based on old Objective-C from the NeXT days.
//
// The basic theory is that rather than IPC via function call, we invoke
// object methods by sending them a message.
//
// Objects may or may not implement (and thus respond to) messages, and
// objects can register themselves in a "responder chain", allowing a single
// message to propagate to multiple objects, or specific objects within that
// chain.

// * Package:

package responder

// * Imports:

import (
	"fmt"
	"sort"
	"sync"

	"github.com/Asmodai/gohacks/events"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	chainTypeName string = "responder.Chain"

	defaultResponderPriority int = 0
	defaultPreallocate       int = 5 // Just a guess.
)

// * Variables:

var (
	// Error condition signalled when an attempt is made to add a
	// non-unique responder to a responder chain.
	ErrDuplicateResponder error = errors.Base("duplicate responder")

	// Error condition signalled when a responder's name via `Name()` is
	// invalid or zero length.
	ErrResponderNameInvalid error = errors.Base("responder name invalid")
)

// * Code:

// ** Types:

// Responder chain structure.
//
// This attempts to bring a little bit of Smalltalk and Objective-C to the
// wonderful world of Go.
//
// It might also attempt to bring a bit of MIT Flavors, too... but don't
// expect to see crazy like `defwrapper` and `defwhopper`.
type Chain struct {
	nameIndex  map[string]int
	typeIndex  map[string][]int
	name       string
	responders []chainEntry
	mu         sync.RWMutex
}

// ** Methods:

// Internal helper method to find a responder with a given name in the chain.
func (chain *Chain) findNamed(name string) (chainEntry, bool) {
	idx, ok := chain.nameIndex[name]
	if !ok || idx >= len(chain.responders) {
		return chainEntry{}, false
	}

	return chain.responders[idx], true
}

// Internal helper to rebuild indices.
func (chain *Chain) rebuildIndices() {
	chain.nameIndex = make(map[string]int)
	chain.typeIndex = make(map[string][]int)

	for idx, ent := range chain.responders {
		chain.nameIndex[ent.name] = idx

		typ := ent.responder.Type()
		chain.typeIndex[typ] = append(chain.typeIndex[typ], idx)
	}
}

// Internal helper to sort responders by priority.
func (chain *Chain) sortResponders() {
	sort.SliceStable(chain.responders, func(i, j int) bool {
		return chain.responders[i].priority > chain.responders[j].priority
	})
}

// Internal helper to append responders to the chain.
func (chain *Chain) appendResponder(name string, responder Respondable, priority int) {
	chain.responders = append(
		chain.responders,
		chainEntry{name: name, responder: responder, priority: priority},
	)
}

// Adds the supplied responder to the chain using a user-specified name and
// priority.
//
// The given name overrides that provided by the `Name()` method in the
// responder, thus should only be used in use-cases where you need a specific
// identifier for a responder.
//
// Returns `ErrDuplicateResponder` if a non-unique responder is added.
func (chain *Chain) AddNamedWithPriority(
	name string,
	responder Respondable,
	priority int,
) (Respondable, error) {
	chain.mu.Lock()
	defer chain.mu.Unlock()

	if len(name) == 0 {
		return nil, errors.WithStack(ErrResponderNameInvalid)
	}

	if _, found := chain.findNamed(name); found {
		return nil, errors.WithMessagef(
			ErrDuplicateResponder,
			"duplicate responder: %q",
			name,
		)
	}

	// Append responder.
	chain.appendResponder(name, responder, priority)

	// Sort the responders array by priority.
	chain.sortResponders()

	// Rebuild indices.
	chain.rebuildIndices()

	return responder, nil
}

// Adds the supplied responder to the chain using a user-specified name.
//
// The given name overrides that provided by the `Name()` method in the
// responder, thus should only be used in use-cases where you need a specific
// identifier for a responder.
//
// The responder will have a default priority of 0.
//
// Returns `ErrDuplicateResponder` if a non-unique responder is added.
func (chain *Chain) AddNamed(name string, responder Respondable) (Respondable, error) {
	return chain.AddNamedWithPriority(
		name,
		responder,
		defaultResponderPriority,
	)
}

// Adds the supplied responder to the chain using a user-specified name and
// priority.
//
// If the responder already exists in the chain then it is replaced with the
// provided responder.
//
// The given name overrides that provided by the `Name()` method in the
// responder, thus should only be used in use-cases where you need a specific
// identifier for a responder.
func (chain *Chain) AddOrReplaceNamedWithPriority(
	name string,
	responder Respondable,
	priority int,
) Respondable {
	chain.mu.Lock()
	defer chain.mu.Unlock()

	if len(name) == 0 {
		return nil
	}

	found := false

	for idx := range chain.responders {
		if chain.responders[idx].name == name {
			chain.responders[idx] = chainEntry{
				name:      name,
				responder: responder,
				priority:  priority}
			found = true

			break
		}
	}

	if !found {
		// Responder doesn't exist, add it.
		chain.appendResponder(name, responder, priority)
	}

	// Sort the responders array by priority.
	chain.sortResponders()

	// Rebuild indices.
	chain.rebuildIndices()

	return responder
}

// Adds the supplied responder to the chain using a user-specified name.
//
// The responder will have a default priority of 0.
//
// If the responder already exists in the chain then it is replaced with the
// provided responder.
//
// The given name overrides that provided by the `Name()` method in the
// responder, thus should only be used in use-cases where you need a specific
// identifier for a responder.
func (chain *Chain) AddOrReplaceNamed(name string, responder Respondable) Respondable {
	return chain.AddOrReplaceNamedWithPriority(
		name,
		responder,
		defaultResponderPriority,
	)
}

// Adds a responder to the responder chain.
//
// The responder will have a default priority of 0.
//
// Returns `ErrDuplicateResponder` if a non-unique responder is added.
func (chain *Chain) Add(responder Respondable) (Respondable, error) {
	return chain.AddNamedWithPriority(
		responder.Name(),
		responder,
		defaultResponderPriority,
	)
}

// Adds the supplied responder to the chain using the given priority.
//
// Returns `ErrDuplicateResponder` if a non-unique responder is added.
func (chain *Chain) AddWithPriority(responder Respondable, priority int) (Respondable, error) {
	return chain.AddNamedWithPriority(
		responder.Name(),
		responder,
		priority,
	)
}

// Adds a responder to the responder chain.
//
// The responder will have a default priority of 0.
//
// If the responder already exists in the chain then it is replaced with the
// provided responder.
func (chain *Chain) AddOrReplace(responder Respondable) Respondable {
	return chain.AddOrReplaceNamedWithPriority(
		responder.Name(),
		responder,
		defaultResponderPriority,
	)
}

// Adds a responder to the responder chain.
//
// If the responder already exists in the chain then it is replaced with the
// provided responder.
func (chain *Chain) AddOrReplaceWithPriority(responder Respondable, priority int) Respondable {
	return chain.AddOrReplaceNamedWithPriority(
		responder.Name(),
		responder,
		priority,
	)
}

// Remove the named responder from the responder chain.
//
// Returns false if no such responder was found.
func (chain *Chain) RemoveNamed(name string) bool {
	idx, ok := chain.nameIndex[name]
	if !ok || idx >= len(chain.responders) {
		return false
	}

	chain.responders = append(
		chain.responders[:idx],
		chain.responders[idx+1:]...,
	)

	chain.sortResponders()
	chain.rebuildIndices()

	return true
}

func (chain *Chain) Remove(responder Respondable) bool {
	if len(responder.Name()) == 0 {
		return false
	}

	return chain.RemoveNamed(responder.Name())
}

// Internal unsafe helper function for sending to named receivers.
//
// This does not do any locking, and assumes that the caller set up a
// read lock.
//
// See documentation for `SendNamed` for details on the return values.
func (chain *Chain) sendNamedUnsafe(name string, event events.Event) (events.Event, bool, bool) {
	responder, ok := chain.findNamed(name)
	if !ok {
		return nil, false, false
	}

	if !responder.Responder().RespondsTo(event) {
		return nil, false, true
	}

	return responder.Responder().Invoke(event), true, true
}

// Send a message to a specific responder.
//
// Returns an object of interface `events.Event`, which may be the same
// event that was passed to it.  Doing sanity on the return value is up to
// you.
// Returns the resulting event from the responder, a boolean value that
// states if the responder was able to respond, and a boolean value that
// states whether the responder was found.
//
// Example usage would be:
//
// ```go
//
//	result, responds, found := someChain.SendNamed("something", evt)
//	if !found {
//		log.Warn("Responder `something` not found!")
//	} else if !responds {
//		log.Info("Responder 'something' ignored the event.")
//	} else {
//		log.Debug("Event was handled.")
//	}
//
// ```
//
// This method is thread-safe.
func (chain *Chain) SendNamed(name string, event events.Event) (events.Event, bool, bool) {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	return chain.sendNamedUnsafe(name, event)
}

// Send a message to a specific responder.  Panics if the responder does
// not exist or does not respond to the event.
func (chain *Chain) MustSendNamed(name string, event events.Event) events.Event {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	evt, responds, found := chain.sendNamedUnsafe(name, event)

	if !found {
		panic(fmt.Sprintf(
			"No responder named %v found for event %T",
			name,
			event))
	}

	if !responds {
		panic(fmt.Sprintf(
			"Responder named %v does not respond to event %T",
			name,
			event))
	}

	return evt
}

// Send a message to the responder chain.
//
// The first responder capable of responding to the event will consume the
// event.
//
// Returns an object of interface `events.Event`, which may be the same event
// that was passed to it.  Doing sanity on the return value is up to you.
//
// Unlike `SendNamed`, there is no indicator that either receivers were not
// found or that there were no receivers.  So it would be wise to assume that
// a `false` means that absolutely nothing in the chain received your event.
//
// Example usage would be:
//
// ```go
//
//	result, ok := someChain.Send(evt)
//	if !ok {
//		log.Warn("No responders have responded to the event.")
//	} else {
//		log.Info("At least one responder has responded to the event.")
//	}
//
// ```
//
// This method is thread-safe.
func (chain *Chain) SendFirst(event events.Event) (events.Event, bool) {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	for idx := range chain.responders {
		responder := chain.responders[idx]

		ret, ok, _ := chain.sendNamedUnsafe(responder.name, event)
		if ok {
			return ret, ok
		}
	}

	return nil, false
}

// Send a message to the responder chain.  Panics if no responder is able
// to respond to the event.
func (chain *Chain) MustSendFirst(event events.Event) events.Event {
	// Do not add locks here or you will deadlock due to `Send` locking.
	evt, ok := chain.SendFirst(event)

	if !ok {
		panic(fmt.Sprintf("No response for event %q", event))
	}

	return evt
}

// Send a message to all responders in the chain.
//
// All responders capable of responding to the event will receive the event.
//
// Returns a list of objects of interface `events.Event`, which may be the
// same event as was passed.  Doing sanity on the return values is up to you
//
// This method is thread-safe.
func (chain *Chain) SendAll(event events.Event) []events.Event {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	results := make([]events.Event, 0, len(chain.responders))

	for idx := range chain.responders {
		responder := chain.responders[idx]

		ret, ok, _ := chain.sendNamedUnsafe(responder.name, event)
		if ok {
			results = append(results, ret)
		}
	}

	return results
}

// Send a message to all responders of a given type in the chain.
//
// All responders of the given type that are capable of responding to the
// event will receive the event.
//
// Returns a list of objects of interface `events.Event`, which may be the
// same event as was passed.  Doing sanity on the return values is up to you.
//
// This method is thread-safe.
func (chain *Chain) SendType(typeName string, event events.Event) []events.Event {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	// No type name?
	if len(typeName) == 0 {
		// XXX log? error condition?
		return []events.Event{}
	}

	// Do we have a type index?
	indices, ok := chain.typeIndex[typeName]
	if !ok {
		// Assume we don't handle the type and bail.
		return []events.Event{}
	}

	results := make([]events.Event, 0, len(indices))

	for _, idx := range indices {
		if idx >= len(chain.responders) {
			// If this happens, we have larger problems.
			panic(fmt.Sprintf(
				"Index for type %#v is out of bounds: %d",
				typeName,
				idx))
		}

		responder := chain.responders[idx]

		ret, ok, _ := chain.sendNamedUnsafe(responder.name, event)
		if ok {
			results = append(results, ret)
		}
	}

	return results
}

// Return the name of the chain.
//
// Implements `Respondable`.
func (chain *Chain) Name() string { return chain.name }

// Return the type name of the chain.
//
// Implements `Respondable`.
func (chain *Chain) Type() string { return chainTypeName }

// Iterate over responders checking if any implement the given event.
//
// The first responder found that responds to the event will result in `true`
// being returned.
//
// Implements `Respondable`.
func (chain *Chain) RespondsTo(event events.Event) bool {
	for _, responder := range chain.responders {
		if responder.responder.RespondsTo(event) {
			return true
		}
	}

	return false
}

// Send an event to the chain.
//
// The first object that can respond to the event will consume it.
//
// Implements `Respondable`.
func (chain *Chain) Invoke(event events.Event) events.Event {
	result, _ := chain.SendFirst(event)

	return result
}

// Return a list of names for the responders currently in the chain.
func (chain *Chain) Names() []string {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	names := make([]string, len(chain.responders))

	for idx := range chain.responders {
		names[idx] = chain.responders[idx].name
	}

	return names
}

// Clear all responders from the chain.
func (chain *Chain) Clear() {
	chain.mu.Lock()
	defer chain.mu.Unlock()

	chain.responders = make([]chainEntry, 0, defaultPreallocate)
	chain.nameIndex = make(map[string]int)
	chain.typeIndex = make(map[string][]int)
}

// Return the number of responders in the chain.
func (chain *Chain) Count() int {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	return len(chain.responders)
}

// Is the chain empty?
func (chain *Chain) IsEmpty() bool {
	chain.mu.RLock()
	defer chain.mu.RUnlock()

	return len(chain.responders) == 0
}

// ** Functions:

// Create a new responder chain object.
func NewChain(name string) *Chain {
	chain := &Chain{
		name: name,
	}

	chain.Clear()

	return chain
}

// * chain.go ends here.
