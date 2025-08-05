// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_ftin_test.go --- FTIN tests.
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

// * Package:

package validator

// * Imports:

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

// * Constants:

// * Variables:

var (
	testFTINStructType reflect.Type
	testFTINStructOnce sync.Once
)

// * Code:

type testFTINStruct struct {
	IntField    int
	Int64Field  int64
	StringField string
	AnyField    any
}

func (t *testFTINStruct) ReflectType() reflect.Type {
	testFTINStructOnce.Do(func() {
		testFTINStructType = reflect.TypeOf(t).Elem()
	})

	return testFTINStructType
}

func TestFTINPredicate(t *testing.T) {
	input := &testFTINStruct{
		IntField:    42,
		Int64Field:  9001,
		StringField: "Hello",
		AnyField:    uint64(654321),
	}

	tests := []struct {
		field string
		types []any
		want  bool
	}{
		{"IntField", []any{"int"}, true},
		{"IntField", []any{"uint", "uint32", "uint64"}, false},
		{"IntField", []any{"int16", "int32", "int64"}, false},
		{"StringField", []any{"string"}, true},
		{"StringField", []any{"any", "uint64"}, false},
		{"AnyField", []any{"string", "bool", "uint32"}, false},
		{"AnyField", []any{"uint8", "uint16", "uint32", "uint64"}, true},
		{"AnyField", []any{"any"}, true},
	}

	inst := &testFTINStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FTIN(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FTINBuilder{}).Build(tt.field, tt.types, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FTIN(%s IN %#v) = %v, want %v",
					tt.field,
					tt.types,
					result,
					tt.want)
			}
		})
	}
}

// * predicate_ftin_test.go ends here.
