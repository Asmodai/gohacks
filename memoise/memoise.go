// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// memoise.go --- Memoisation hacks.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package memoise

import (
	"gitlab.com/tozd/go/errors"
)

// Memoisation function type.
type CallbackFn func() (any, error)

// Memoisation type.
type Memoise interface {
	// Check if we have a memorised value for a given key.  If not, then
	// inovke the callback function and memorise its result.
	Check(string, CallbackFn) (any, error)
}

// Implementation of the memoisation type.
type memoise map[string]any

// Create a new memoisation object.
func NewMemoise() Memoise {
	return memoise{}
}

// Check the map of memoised values fo#r the given key.  If the key exists,
// then return its associated value.  Otherwise, obtain the value via the
// given memoisation function.
func (obj memoise) Check(name string, fn CallbackFn) (any, error) {
	if result, ok := obj[name]; ok {
		return result, nil
	}

	res, err := fn()
	if err == nil {
		obj[name] = res
	}

	return obj[name], errors.WithStack(err)
}

// memoise.go ends here.
