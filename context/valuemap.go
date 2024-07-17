/* mock:yes */
/*
 * valuemap.go --- Value map structure.
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
}

// Internal structure.
type valueMap struct {
	data map[string]any // Map of string keys to any value type.
}

// Create a new value map with no data.
func NewValueMap() ValueMap {
	return valueMap{
		data: map[string]any{},
	}
}

// Returns a value associated with the given key.
//
// If the key is present, then it is returned along with an `ok` value of
// true.
//
// Otherwise, nil is returned with an `ok` value of false.
func (obj valueMap) Get(key string) (any, bool) {
	value, ok := obj.data[key]

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
func (obj valueMap) Set(key string, value any) {
	obj.data[key] = value
}

/* valuemap.go ends here. */
