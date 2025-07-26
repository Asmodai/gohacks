// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// memoise.go --- Memoisation hacks.
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
//
// mock:yes

// * Comments:

//
//
//

// * Package:

package memoise

// * Imports:

import (
	"sync"

	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Types:

// Memoisation function type.
type CallbackFn func() (any, error)

// Memoisation type.
type Memoise interface {
	// Check if we have a memorised value for a given key.  If not, then
	// inovke the callback function and memorise its result.
	Check(string, CallbackFn) (any, error)
}

// Implementation of the memoisation type.
type memoise struct {
	sync.RWMutex

	store map[string]any
}

// ** Methods:

// Check the map of memoised values fo#r the given key.  If the key exists,
// then return its associated value.  Otherwise, obtain the value via the
// given memoisation function.
func (obj *memoise) Check(name string, callback CallbackFn) (any, error) {
	obj.RLock()
	result, ok := obj.store[name]
	obj.RUnlock()

	// If we get a hit, return it.
	if ok {
		return result, nil
	}

	// Miss, so obtain a lock.
	obj.Lock()
	defer obj.Unlock()

	// Sanity check in case we lost a race.
	if result, ok := obj.store[name]; ok {
		return result, nil
	}

	// Still a miss.
	res, err := callback()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Store the result.
	obj.store[name] = res

	return res, errors.WithStack(err)
}

// ** Functions:

// Create a new memoisation object.
func NewMemoise() Memoise {
	return &memoise{
		store: map[string]any{},
	}
}

// * memoise.go ends here.
