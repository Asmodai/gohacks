// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// timedcache.go --- Timed cache context value.
//
// Copyright (c) 2023-2025 Paul Ward <paul@lisphacker.uk>
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

package contextdi

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/timedcache"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	ContextKeyTimedCache = "_DI_TIMEDCACHE"
)

// * Variables:

var (
	ErrValueNotTimedCache = errors.Base("value is not timedcache.TimedCache")
)

// * Code:

// ** Functions:

// Set the timed cache value in the context map.
func SetTimedCache(ctx context.Context, inst timedcache.TimedCache) (context.Context, error) {
	return PutToContext(ctx, ContextKeyTimedCache, inst)
}

// Get the timed cache value from the given context.
//
// WIll return `ErrValueNotTimedCache` if the value in the context is not
// of type `timedcache.TimedCache`.
func GetTimedCache(ctx context.Context) (timedcache.TimedCache, error) {
	val, err := GetFromContext(ctx, ContextKeyTimedCache)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	inst, ok := val.(timedcache.TimedCache)
	if !ok {
		return nil, errors.WithStack(ErrValueNotTimedCache)
	}

	return inst, nil
}

// Attempt to get the timed cache value from the given context.  Panics if
// the operation fails.
func MustGetTimedCache(ctx context.Context) timedcache.TimedCache {
	inst, err := GetTimedCache(ctx)
	if err != nil {
		panic("Could not get timed cache instance from context")
	}

	return inst
}

// * timedcache.go ends here.
