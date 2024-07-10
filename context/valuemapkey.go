/*
 * valuemapkey.go --- Context value key for value maps.
 *
 * Copyright (c) 2024 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package context

import (
	. "context"
)

// ValueMap key type for `WithValue`.
type ValueMapKey string

const (
	mapKey ValueMapKey = ValueMapKey("_CTXVALMAP")
)

// Get the value map (if any) from the context.
//
// Returns nil if there is no value map.
func GetValueMap(ctx Context) ValueMap {
	if v := ctx.Value(mapKey); v != nil {
		return v.(ValueMap)
	}

	return nil
}

// Get the value map (if any) from the context with the specified value
// key.
func GetValueMapWithKey(ctx Context, key string) ValueMap {
	if v := ctx.Value(ValueMapKey(key)); v != nil {
		return v.(ValueMap)
	}

	return nil
}

// Create a context with the value map using a default key.
func WithValueMap(ctx Context, valuemap ValueMap) Context {
	return WithValue(ctx, mapKey, valuemap)
}

// Create a context with the value map using the specified key.
func WithValueMapWithKey(ctx Context, key string, valuemap ValueMap) Context {
	return WithValue(ctx, ValueMapKey(key), valuemap)
}

/* valuemapkey.go ends here. */
