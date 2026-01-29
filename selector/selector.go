// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// selector.go --- Selectors.
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
	"context"
	"slices"
	"sync"

	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/metadata"
	"github.com/Asmodai/gohacks/responder"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	DefaultPriority        = int(100)
	defaultMaxForwardDepth = 8
)

// * Code:

// ** Types:

// Selector method function signature type.
type Method func(responder.Respondable, events.Event) events.Event

type AuxiliaryMethod struct {
	priority int
	method   Method
}

// Selector table entry.
type Entry struct {
	primary Method            // Primary method.
	before  []AuxiliaryMethod // Methods to invoke before primary.
	after   []AuxiliaryMethod // Methods to invoke after primary.
	mdata   metadata.Metadata // Metadata.
}

// Map selector names to method implementations.
type Table struct {
	mu              sync.RWMutex
	selectors       map[string]*Entry
	defaultSelector string
	maxForwardDepth int
}

// ** Methods:

// Register a method for a selector.
func (st *Table) Register(selector string, method Method) {
	st.mu.Lock()
	defer st.mu.Unlock()

	entry, exists := st.selectors[selector]
	if !exists {
		entry = newEntry()
		st.selectors[selector] = entry
	}

	entry.primary = method
}

func (st *Table) Unregister(selector string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	delete(st.selectors, selector)
}

func (st *Table) SetPrimary(selector string, method Method) (Method, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	if method == nil {
		return method, errors.WithStack(ErrNoMethodSpecified)
	}

	entry, exists := st.selectors[selector]
	if !exists {
		return method, errors.WithMessagef(
			ErrSelectorNotFound,
			"%q",
			selector)
	}

	entry.primary = method

	return method, nil
}

func (st *Table) AddBefore(selector string, method Method) error {
	return st.AddBeforeWithPriority(DefaultPriority, selector, method)
}

func (st *Table) AddBeforeWithPriority(priority int, selector string, method Method) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	entry, exists := st.selectors[selector]
	if !exists {
		return errors.WithMessagef(
			ErrNoMethodToWrap,
			"selector %q",
			selector)
	}

	// Copy to preserve old slice.
	before := append([]AuxiliaryMethod(nil), entry.before...)
	before = append(before, AuxiliaryMethod{
		priority: priority,
		method:   method})

	sortAux(before)

	// Write.
	entry.before = before

	return nil
}

func (st *Table) AddAfter(selector string, method Method) error {
	return st.AddAfterWithPriority(DefaultPriority, selector, method)
}

func (st *Table) AddAfterWithPriority(priority int, selector string, method Method) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	entry, exists := st.selectors[selector]
	if !exists {
		return errors.WithMessagef(
			ErrNoMethodToWrap,
			"selector %q",
			selector)
	}

	// Copy to preserve old slice.
	after := append([]AuxiliaryMethod(nil), entry.after...)
	after = append(after, AuxiliaryMethod{
		priority: priority,
		method:   method})

	sortAux(after)

	// Write.
	entry.after = after

	return nil
}

func (st *Table) SetMaxForwardDepth(val int) {
	st.mu.Lock()
	defer st.mu.Unlock()

	st.maxForwardDepth = val
}

func (st *Table) SetDefault(selector string) {
	st.mu.Lock()
	st.defaultSelector = selector
	st.mu.Unlock()
}

func (st *Table) resolveSelector(selector string) (string, *Entry, bool) {
	if entry, found := st.selectors[selector]; found {
		return selector, entry, true
	}

	if len(st.defaultSelector) > 0 {
		return st.defaultSelector, st.selectors[st.defaultSelector], true
	}

	return "", nil, false
}

// Check whether a selector is defined.
func (st *Table) HasSelector(selector string) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	_, found := st.selectors[selector]

	return found
}

// Return an entry for a selector.
func (st *Table) Get(name string) (*Entry, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	entry, found := st.selectors[name]

	return entry, found
}

// Dispatch a message to a selector.
//
//nolint:cyclop,funlen
func (st *Table) invoke(
	sel string,
	target responder.Respondable, event events.Event,
	depth int, path []string, limit int,
) (events.Event, bool) {
	var (
		result events.Event
		retok  bool
	)

	if depth > limit {
		err := errors.WithMessagef(
			ErrForwardLoop,
			"depth=%d path=%v",
			depth,
			path)

		Trace(err.Error())

		return NewSelectorError(err), false
	}

	// START CRITICAL SECTION.
	st.mu.RLock()

	effsel, entry, found := st.resolveSelector(sel)

	if !found || entry == nil || entry.primary == nil {
		st.mu.RUnlock()

		return event, false
	}

	// Snapshot.
	before := append([]AuxiliaryMethod(nil), entry.before...)
	after := append([]AuxiliaryMethod(nil), entry.after...)
	primary := entry.primary

	st.mu.RUnlock()
	// END CRITICAL SECTION.

	Trace("Invoking selector %q on %s (%s)",
		effsel,
		target.Name(),
		target.Type())

	defer func() {
		if rec := recover(); rec != nil {
			Trace("Selector %q panicked: %v", effsel, rec)

			retok = false
			result = NewSelectorError(
				errors.BaseWrapf(
					ErrSelectorPanic,
					"%q panic: %v",
					effsel,
					rec))
		}
	}()

	// NOTE: We diverge from CLOS here slightly --  we do not have the
	// `:around` auxiliary method.
	//
	// Here, `:before` acts as a gatekeeper to the primary method.
	// This means that if any method in the `:before` chain returns nil
	// or error, then the train stops right there.
	//
	// The `:after`method chain is for mutation and/or cleanup, their
	// return values are ignored.  They are invoked with the result
	// of the primary.

	curr := event

	// Execute the `:before` auxiliaries in order.
	for idx, wrap := range before {
		Trace("Executing :before[%d] for %q [prio=%d]",
			idx,
			effsel,
			wrap.priority)

		out := wrap.method(target, curr)

		if errEvt, isErr := out.(*SelectorError); isErr || out == nil {
			msg := out.String()

			if isErr {
				msg = errEvt.Error().Error()
			}

			Trace(":before[%d] for %q short-circuited: %v",
				idx,
				effsel,
				msg)

			return out, true
		}

		// Treat non-nil/non-error as a transformed event.
		curr = out
	}

	// Try forwarding.
	if fwd, ok := curr.(*events.Forward); ok {
		next := fwd.To()

		return st.invoke(
			next,
			target,
			fwd.Event(),
			depth+1,
			append(path, next),
			limit)
	}

	Trace("Executing :primary for %q", effsel)

	// Execute the primary.
	result = primary(target, curr)
	retok = true

	// Execute the `:after` auxiliaries in order, discarding results.
	for idx, wrap := range after {
		Trace("Executing :after[%d] for %q [prio=%d]",
			idx,
			effsel,
			wrap.priority)

		_ = wrap.method(target, result)
	}

	Trace("Selector %q completed", effsel)

	return result, retok
}

func (st *Table) InvokeSelector(sel string, tgt responder.Respondable, evt events.Event) (events.Event, bool) {
	depthLimit := st.maxForwardDepth

	if depthLimit <= 0 {
		depthLimit = defaultMaxForwardDepth
	}

	return st.invoke(sel, tgt, evt, 0, []string{sel}, depthLimit)
}

func (st *Table) InvokeSelectorAsync(
	ctx context.Context,
	sel string,
	tgt responder.Respondable, evt events.Event,
) <-chan events.Event {
	out := make(chan events.Event, 1)

	go func() {
		defer close(out)

		res, _ := st.InvokeSelector(sel, tgt, evt)

		select {
		case out <- res:
		case <-ctx.Done():
		}
	}()

	return out
}

func (st *Table) Metadata(selector string) (metadata.Metadata, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	entry, found := st.selectors[selector]
	if !found {
		return nil, errors.WithMessagef(
			ErrSelectorNotFound,
			"%q",
			selector)
	}

	return entry.mdata, nil
}

func (st *Table) MustMetadata(selector string) metadata.Metadata {
	meta, err := st.Metadata(selector)
	if err != nil {
		panic(errors.WithStack(err))
	}

	return meta
}

func (st *Table) ListMetadata(selector string) (map[string]string, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()

	entry, found := st.selectors[selector]
	if !found {
		return map[string]string{}, errors.WithMessagef(
			ErrSelectorNotFound,
			"%q",
			selector)
	}

	return entry.mdata.List(), nil
}

func (st *Table) AllMetadata() map[string]map[string]string {
	st.mu.RLock()
	defer st.mu.RUnlock()

	results := make(map[string]map[string]string, len(st.selectors))

	for name, sel := range st.selectors {
		results[name] = sel.mdata.List()
	}

	return results
}

// ** Functions:

func sortAux(methods []AuxiliaryMethod) {
	slices.SortStableFunc(methods, func(a, b AuxiliaryMethod) int {
		switch {
		case a.priority < b.priority:
			return -1

		case a.priority > b.priority:
			return 1

		default:
			return 0
		}
	})
}

func NewTable() *Table {
	return &Table{
		selectors: make(map[string]*Entry),
	}
}

func newEntry() *Entry {
	return &Entry{
		before: []AuxiliaryMethod{},
		after:  []AuxiliaryMethod{},
		mdata:  metadata.NewMetadata(),
	}
}

// * selector.go ends here.
