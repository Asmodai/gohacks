// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// bindings.go --- Bindings manager.
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

//
//
//

// * Package:

package validator

// * Imports:

import (
	"reflect"
	"sync"
)

// * Constants:

// * Variables:

// * Code:

// ** Types:

type Bindings struct {
	mu sync.RWMutex

	Bindings map[reflect.Type]*StructDescriptor
}

func (b *Bindings) getType(object any) reflect.Type {
	typeOf := reflect.TypeOf(object)

	if typeOf.Kind() == reflect.Pointer {
		typeOf = typeOf.Elem()
	}

	return typeOf
}

func (b *Bindings) Register(object *StructDescriptor) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	typ := object.Type

	if _, found := b.Bindings[typ]; found {
		return false
	}

	b.Bindings[typ] = object

	return true
}

func (b *Bindings) build(object any, typeOf reflect.Type) (*StructDescriptor, bool) {
	if object == nil {
		panic("Attempt made to build descriptor for nil.")
	}

	desc := BuildDescriptor(typeOf)
	if desc == nil {
		return nil, false
	}

	return desc, b.Register(desc)
}

func (b *Bindings) Build(object Reflectable) (*StructDescriptor, bool) {
	return b.build(object, object.ReflectType())
}

func (b *Bindings) BuildWithReflection(object any) (*StructDescriptor, bool) {
	if object == nil {
		return nil, false
	}

	return b.build(object, b.getType(object))
}

func (b *Bindings) bind(object any, typeOf reflect.Type) (*BoundObject, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	descriptor, found := b.Bindings[typeOf]
	if !found {
		return nil, false
	}

	bound := &BoundObject{
		Descriptor: descriptor,
		Binding:    object,
	}

	return bound, true
}

func (b *Bindings) Bind(object Reflectable) (*BoundObject, bool) {
	return b.bind(object, object.ReflectType())
}

func (b *Bindings) BindWithReflection(object any) (*BoundObject, bool) {
	if object == nil {
		return nil, false
	}

	return b.bind(object, b.getType(object))
}

func NewBindings() *Bindings {
	return &Bindings{
		Bindings: make(map[reflect.Type]*StructDescriptor),
	}
}

// * bindings.go ends here.
