// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// valuemapkey.go --- Context value key for value maps.
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

package contextext

import (
	"gitlab.com/tozd/go/errors"

	"context"
)

// ValueMap key type for `WithValue`.
type ValueMapKey string

const (
	mapKey ValueMapKey = ValueMapKey("_CTXVALMAP")
)

var (
	ErrInvalidContext   = errors.Base("invalid context")
	ErrInvalidValueMap  = errors.Base("invalid value map")
	ErrValueMapNotFound = errors.Base("value map not found")
)

func extractValueMap(thing any) (ValueMap, error) {
	if vmap, ok := thing.(ValueMap); ok {
		return vmap, nil
	}

	return nil, errors.WithStack(ErrInvalidValueMap)
}

// Get the value map (if any) from the context.
//
// Returns nil if there is no value map.
func GetValueMap(ctx context.Context) (ValueMap, error) {
	if ctx == nil {
		return nil, errors.WithStack(ErrInvalidContext)
	}

	if v := ctx.Value(mapKey); v != nil {
		return extractValueMap(v)
	}

	return nil, errors.WithStack(ErrValueMapNotFound)
}

// Get the value map (if any) from the context with the specified value
// key.
func GetValueMapWithKey(ctx context.Context, key string) (ValueMap, error) {
	if v := ctx.Value(ValueMapKey(key)); v != nil {
		return extractValueMap(v)
	}

	return nil, errors.WithStack(ErrInvalidValueMap)
}

// Create a context with the value map using a default key.
func WithValueMap(ctx context.Context, valuemap ValueMap) context.Context {
	return context.WithValue(ctx, mapKey, valuemap)
}

// Create a context with the value map using the specified key.
func WithValueMapWithKey(ctx context.Context, key string, valuemap ValueMap) context.Context {
	return context.WithValue(ctx, ValueMapKey(key), valuemap)
}

// valuemapkey.go ends here.
