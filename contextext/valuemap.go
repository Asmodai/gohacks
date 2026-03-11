// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// valuemap.go --- Value map structure.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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

package contextext

import (
	"fmt"
	"sync"
)

// A map-based storage structure to pass multiple values via contexts
// rather than many invocations of `context.WithValue` and their respective
// copy operations.
//
// The main caveat with this approach is that as contexts are copied by the
// various `With` functions we have no means of passing changes to child
// contexts once the context with the value map is copied.
//
// This is not the main aim of this type, so such functionality should not
// be considered.  The main usage is to provide a means of passing a lot of
// values to some top-level context in order to avoid a lot of `WithValue`
// calls and a somewhat slow lookup.
type ValueMap interface {
	Get(string) (key any, ok bool)
	Set(key string, value any)
	Immutable() bool
	Finalise()
}

// Internal structure.
type valueMap struct {
	sync.Mutex
	data      map[string]any // Map of string keys to any value type.
	finalised bool           // Can we write further elements to the map?
}

// Create a new value map with no data.
func NewValueMap() ValueMap {
	return &valueMap{
		data: map[string]any{},
	}
}

// Is the value map immutable?
func (vm *valueMap) Immutable() bool {
	vm.Lock()
	defer vm.Unlock()

	return vm.finalised
}

// Finalise the value map.
//
// Once finalised, the map is treated as immutable.  Further `Set' operations
// will silently return.
func (vm *valueMap) Finalise() {
	vm.Lock()
	vm.finalised = true
	vm.Unlock()
}

// Returns a value associated with the given key.
//
// If the key is present, then it is returned along with an `ok` value of
// true.
//
// Otherwise, nil is returned with an `ok` value of false.
func (vm *valueMap) Get(key string) (any, bool) {
	vm.Lock()
	defer vm.Unlock()

	value, ok := vm.data[key]

	return value, ok
}

// Set the value of the given key.
//
// This will overwrite any existing value.
//
// Be aware that using this method within a context's scope will only
// affect the value within that scope.  The map's contents in the parent
// will not be affected and only children that inherit from the context
// *after* any `set` operation will see the changes.  This is due to
// the context's value field being copied.
func (vm *valueMap) Set(key string, value any) {
	vm.Lock()
	defer vm.Unlock()

	if vm.finalised {
		panic(fmt.Sprintf("attempted to set %q on immutable value map", key))
	}

	vm.data[key] = value
}

// valuemap.go ends here.
