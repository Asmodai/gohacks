// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// selector_respondable.go --- Selector respondable type.
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

// * Package:

package selector

import (
	"slices"

	"github.com/Asmodai/gohacks/events"
	"gitlab.com/tozd/go/errors"
)

// * Imports:

// * Constants:

// * Code:

type Respondable struct {
	name      string
	typeName  string
	methods   *Table
	protocols map[string]bool
}

// ** Methods:

func (sr *Respondable) Name() string { return sr.name }
func (sr *Respondable) Type() string { return sr.typeName }

func (sr *Respondable) AddProtocol(name string) {
	if sr.protocols == nil {
		sr.protocols = make(map[string]bool)
	}

	sr.protocols[name] = true
}

// *** `Introspectable` methods:

func (sr *Respondable) Selectors() []string {
	selectors := make([]string, 0)

	sr.methods.mu.RLock()
	defer sr.methods.mu.RUnlock()

	for sel := range sr.methods.selectors {
		selectors = append(selectors, sel)
	}

	return selectors
}

func (sr *Respondable) SortedSelectors() []string {
	selectors := sr.Selectors()

	slices.Sort(selectors)

	return selectors
}

func (sr *Respondable) Methods() *Table {
	return sr.methods
}

func (sr *Respondable) ConformsTo(name string) bool {
	if sr.protocols == nil {
		return false
	}

	return sr.protocols[name]
}

func (sr *Respondable) ListProtocols() []string {
	names := make([]string, 0, len(sr.protocols))

	for proto := range sr.protocols {
		names = append(names, proto)
	}

	return names
}

func (sr *Respondable) MetadataForSelector(selector string) (map[string]string, error) {
	return sr.methods.ListMetadata(selector)
}

// *** `Respondable` methods:

func (sr *Respondable) RespondsTo(evt events.Event) bool {
	selEvt, selEvtOk := evt.(SelectorEvent)
	if !selEvtOk {
		return false
	}

	return sr.methods.HasSelector(selEvt.Selector())
}

func (sr *Respondable) Invoke(evt events.Event) events.Event {
	selEvt, selEvtOk := evt.(SelectorEvent)
	if !selEvtOk {
		return NewSelectorError(errors.WithMessagef(
			ErrHasNoSelector,
			"event %T",
			evt))
	}

	if result, ok := sr.methods.InvokeSelector(selEvt.Selector(), sr, evt); ok {
		return result
	}

	return NewSelectorError(errors.WithMessagef(
		ErrSelectorNotFound,
		"Event %T selector %q",
		evt,
		selEvt.Selector()))
}

// ** Functions:

func NewRespondable(name, typeName string) *Respondable {
	return &Respondable{
		name:     name,
		typeName: typeName,
		methods:  NewTable(),
	}
}

// * selector_respondable.go ends here.
