// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// valuemap_test.go --- Value map tests.
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

package contextext

import (
	"context"
	"reflect"
	"testing"
)

// valueMap structure tests.
func TestValueMap(t *testing.T) {
	var vmap ValueMap = nil

	t.Run("Constructs properly", func(t *testing.T) {
		vmap = NewValueMap()
		tnam := reflect.TypeOf(vmap).Name()

		// Remember, this wants the internal type.
		if tnam != "valueMap" {
			t.Errorf("Got unexpected type %s (%v)", tnam, vmap)
		}
	})

	t.Run("Get returns nil for non-existing key", func(t *testing.T) {
		if val, ok := vmap.Get("foo"); ok {
			t.Errorf("Key has unexpected value: foo = %#v", val)
		}
	})

	t.Run("Get returns non-nil for existing key", func(t *testing.T) {
		vmap.Set("testing", "yes")
		val, ok := vmap.Get("testing")

		if !ok {
			t.Error(`Key "testing" lacks a value.`)
		}

		if val != "yes" {
			t.Errorf("Value is incorrect: testing = %#v", val)
		}
	})
}

// WithValueMap/GetValueMap tests.
func TestValueMapDefaultKey(t *testing.T) {
	vmap := NewValueMap()
	vmap.Set("test", 42)

	ctx := WithValueMap(context.TODO(), vmap)

	t.Run("Returns existing value", func(t *testing.T) {
		if _, err := GetValueMap(ctx); err != nil {
			t.Error("No value map key was found!")
		}
	})

	t.Run("Returns nil for non-existent value", func(t *testing.T) {
		ctx := context.TODO()
		if res, err := GetValueMap(ctx); res != nil && err != nil {
			t.Errorf("Unexpected value returned: %#v", res)
		}
	})

	t.Run("Returns value for key within map", func(t *testing.T) {
		res, err := GetValueMap(ctx)
		if err != nil {
			t.Error("No value map was found!")
		}

		if _, ok := res.Get("test"); !ok {
			t.Error("Value for key 'test' was not found!")
		}
	})

	t.Run("Returns false for value without key in map", func(t *testing.T) {
		res, err := GetValueMap(ctx)
		if err != nil {
			t.Error("No value map was found!")
		}

		if _, ok := res.Get("nope"); ok {
			t.Error("Somehow a value for a non-existent key was found.")
		}
	})
}

// WithValueMapWithKey/GetValueMapWithKey tests.
func TestValueMapCustomKey(t *testing.T) {
	vmap := NewValueMap()
	vmap.Set("test", 42)

	ctx := WithValueMapWithKey(context.TODO(), "testing", vmap)

	t.Run("Returns existing value", func(t *testing.T) {
		if _, err := GetValueMapWithKey(ctx, "testing"); err != nil {
			t.Error("No value map key was found!")
		}
	})

	t.Run("Returns nil for non-existent key", func(t *testing.T) {
		if res, err := GetValueMapWithKey(ctx, "nope"); res != nil && err != nil {
			t.Error("Somehow found a key that shouldn't exist.")
		}
	})

	t.Run("Returns value for key within map", func(t *testing.T) {
		res, err := GetValueMapWithKey(ctx, "testing")
		if err != nil {
			t.Error("No value map was found!")
		}

		if _, ok := res.Get("test"); !ok {
			t.Error("Value for key 'test' was not found!")
		}
	})

	t.Run("Returns false for value without key in map", func(t *testing.T) {
		res, err := GetValueMapWithKey(ctx, "testing")
		if err != nil {
			t.Error("No value map was found!")
		}

		if _, ok := res.Get("nope"); ok {
			t.Error("Somehow a value for a non-existent key was found.")
		}
	})
}

func TestChildCopy(t *testing.T) {

	vmap := NewValueMap()
	vmap.Set("test1", "One")
	vmap.Set("test2", 2)

	parent := WithValueMap(context.TODO(), vmap)
	child, _ := context.WithCancel(parent)

	// Test the parent here, just to be sure.
	t.Run("Parent has values", func(t *testing.T) {
		vals, err := GetValueMap(parent)
		if vals == nil || err != nil {
			t.Error("Parent has no value map.")
		}

		t.Run("test1 is ok", func(t *testing.T) {
			val, ok := vals.Get("test1")
			if !ok {
				t.Error("Parent does not have 'test1'.")
			}
			if val != "One" {
				t.Errorf(`Unexpected value, val = #%v != "One"`, val)
			}
		})

		t.Run("test2 is ok", func(t *testing.T) {
			val, ok := vals.Get("test2")
			if !ok {
				t.Error("Parent does not have 'test2'.")
			}
			if val != 2 {
				t.Errorf(`Unexpected value, val = #%v != 2`, val)
			}
		})
	})

	// Now test the child.
	t.Run("Child has values", func(t *testing.T) {
		vals, err := GetValueMap(child)
		if vals == nil || err != nil {
			t.Error("Child has no value map.")
		}

		t.Run("test1 is ok", func(t *testing.T) {
			val, ok := vals.Get("test1")
			if !ok {
				t.Error("Child does not have 'test1'.")
			}
			if val != "One" {
				t.Errorf(`Unexpected value, val = #%v != "One"`, val)
			}
		})

		t.Run("test2 is ok", func(t *testing.T) {
			val, ok := vals.Get("test2")
			if !ok {
				t.Error("Child does not have 'test2'.")
			}
			if val != 2 {
				t.Errorf(`Unexpected value, val = #%v != 2`, val)
			}
		})
	})
}

// valuemap_test.go ends here.
