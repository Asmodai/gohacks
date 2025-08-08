// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// registry.go --- Protocol registry.
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

// * Package:

package protocols

// * Imports:

import (
	"sync"

	"github.com/Asmodai/gohacks/selector"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	ErrNoVerifierFunction = errors.Base("no verifier function")
)

// * Code:

// ** Types:

type Verifier func(selector.Introspectable) error

type Registry struct {
	mu        sync.RWMutex
	protocols map[string]*Protocol
	verifiers map[string]Verifier
}

// ** Methods:

func (r *Registry) Register(proto *Protocol) {
	r.RegisterWithVerifier(proto, nil)
}

func (r *Registry) RegisterWithVerifier(proto *Protocol, verifier Verifier) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.protocols[proto.Name] = proto
	r.verifiers[proto.Name] = verifier
}

func (r *Registry) Verify(name string, obj selector.Introspectable) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if fun, found := r.verifiers[name]; found && fun != nil {
		return fun(obj)
	}

	return errors.WithMessagef(
		ErrNoVerifierFunction,
		"%q",
		name)
}

type hasMethodsIntrospector interface {
	Methods() *selector.Table
}

func (r *Registry) Validate(name string, rbl hasMethodsIntrospector) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	proto, found := r.protocols[name]
	if !found {
		return false
	}

	for _, sel := range proto.Selectors {
		if !rbl.Methods().HasSelector(sel) {
			return false
		}
	}

	return true
}

// ** Functions:

func NewRegistry() *Registry {
	return &Registry{
		protocols: make(map[string]*Protocol),
		verifiers: make(map[string]Verifier),
	}
}

// * registry.go ends here.
