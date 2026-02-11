// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// workingmem.go --- Working memory.
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

package expertsys

// * Imports:

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// * Code:

// ** Types:

// Working memory implementation.
//
// Implements `WorkingMemory`.
type workingMemory struct {
	facts   map[string]any
	version atomic.Uint64
	mu      sync.RWMutex
}

// ** Methods:

// Return the fact (if any) for the given key.
func (wm *workingMemory) Get(key string) (any, bool) {
	var (
		val   any
		found bool
	)

	wm.mu.RLock()
	{
		val, found = wm.facts[key]
	}
	wm.mu.RUnlock()

	return val, found
}

// Update the existing key with the given value.
//
// Returns `true` if the key was changed; otherwise `false` is returned.
//
// If no key is found, then one is created and assigned the given value.
// If a key is found and the values are approximately equal, no change occurs.
//
// The memory's version number is incremented after successful update.
func (wm *workingMemory) Set(key string, val any) bool {
	wm.mu.Lock()
	{
		old, found := wm.facts[key]

		if found && equalish(old, val) {
			wm.mu.Unlock()

			return false
		}

		wm.facts[key] = val
	}
	wm.mu.Unlock()

	wm.version.Add(1)

	return true
}

// Return a sorted list of keys present in the working memory.
func (wm *workingMemory) Keys() []string {
	var keys []string

	wm.mu.RLock()
	{
		keys = make([]string, 0, len(wm.facts))

		for key := range wm.facts {
			keys = append(keys, key)
		}
	}
	wm.mu.RUnlock()

	sort.Strings(keys)

	return keys
}

// Return the string representation of the working memory.
func (wm *workingMemory) String() string {
	keys := wm.Keys()

	var sbld strings.Builder

	sbld.WriteRune('{')

	for idx, key := range keys {
		if idx > 0 {
			sbld.WriteString(", ")
		}

		val, _ := wm.Get(key)

		fmt.Fprintf(&sbld, "%s=%v", key, val)
	}

	sbld.WriteRune('}')

	return sbld.String()
}

// Return the working memory's current version.
func (wm *workingMemory) Version() uint64 {
	return wm.version.Load()
}

// ** Functions:

// Is the lhs sufficiently equal to the rhs?
//
//nolint:cyclop,nlreturn,funlen
func equalish(lhs, rhs any) bool {
	switch lval := lhs.(type) {
	case nil:
		return rhs == nil
	case string:
		rval, ok := rhs.(string)
		return ok && lval == rval

	case bool:
		rval, ok := rhs.(bool)
		return ok && lval == rval

	case int:
		rval, ok := rhs.(int)
		return ok && lval == rval

	case int32:
		rval, ok := rhs.(int32)
		return ok && lval == rval

	case int64:
		rval, ok := rhs.(int64)
		return ok && lval == rval

	case uint:
		rval, ok := rhs.(uint)
		return ok && lval == rval

	case uint32:
		rval, ok := rhs.(uint32)
		return ok && lval == rval

	case uint64:
		rval, ok := rhs.(uint64)
		return ok && lval == rval

	case float32:
		rval, ok := rhs.(float32)
		return ok && lval == rval

	case float64:
		rval, ok := rhs.(float64)
		return ok && lval == rval

	case time.Time:
		rval, ok := rhs.(time.Time)
		return ok && lval.Equal(rval)

	case time.Duration:
		rval, ok := rhs.(time.Duration)
		return ok && lval == rval

	default:
		// Fallback:  treat as changed.
		return false
	}
}

// Create a new working memory instance.
func NewWorkingMemory() WorkingMemory {
	return &workingMemory{
		facts: make(map[string]any),
	}
}

// * workingmem.go ends here.
