// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// bound.go --- Bound objects.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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

package validator

// * Code:

// ** Types:

type BoundObject struct {
	Descriptor *StructDescriptor
	Binding    any
}

// Get the value for the given key from the bound object.
//
// This works by using the accessor obtained via reflection during the
// predicate building phase.
func (bo *BoundObject) GetValue(key string) (any, bool) {
	field, ok := bo.Descriptor.Get(key)
	if !ok {
		return nil, false
	}

	finfo, finfoOk := field.(*FieldInfo)
	if !finfoOk {
		return nil, false
	}

	val := finfo.Accessor(bo.Binding)

	return val, true
}

func (bo *BoundObject) Get(key string) (any, bool) {
	return bo.Descriptor.Get(key)
}

func (bo *BoundObject) Set(_ string, _ any) bool {
	return false
}

func (bo *BoundObject) Keys() []string {
	return bo.Descriptor.Keys()
}

func (bo *BoundObject) String() string {
	return bo.Descriptor.String()
}

// * bound.go ends here.
