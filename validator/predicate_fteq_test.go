// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fteq_test.go --- FTEQ tests.
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

// * Variables:

var (
	testFTEQStructType reflect.Type
	testFTEQStructOnce sync.Once
)

// * Code:

// ** Types:

type testFTEQStruct struct {
	IntField    int
	Int64Field  int64
	StringField string
	AnyField    any
}

func (t *testFTEQStruct) ReflectType() reflect.Type {
	testFTEQStructOnce.Do(func() {
		testFTEQStructType = reflect.TypeOf(t).Elem()
	})

	return testFTEQStructType
}

// ** Tests:

func TestFTEQPredicate(t *testing.T) {
	input := &testFTEQStruct{
		IntField:    42,
		Int64Field:  9001,
		StringField: "Hello",
		AnyField:    uint64(654321),
	}

	tests := []struct {
		field string
		types any
		want  bool
	}{
		{"IntField", "int", true},
		{"IntField", "uint64", false},
		{"Int64Field", "int", false},
		{"Int64Field", "int64", true},
		{"StringField", "string", true},
		{"StringField", "any", false},
		{"AnyField", "uint32", false},
		{"AnyField", "uint64", true},
		{"AnyField", "any", true},
	}

	inst := &testFTEQStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FTEQ(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FTEQBuilder{}).Build(tt.field, tt.types, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FTEQ(%s IN %#v) = %v, want %v",
					tt.field,
					tt.types,
					result,
					tt.want)
			}
		})
	}
}

// ** Benchmarks:

func BenchmarkFTEQPredicate(b *testing.B) {
	input := &testFTEQStruct{
		IntField:    42,
		Int64Field:  9001,
		StringField: "Hello",
		AnyField:    uint64(654321),
	}

	field := "AnyField"
	vals := "uint64"

	inst := &testFTEQStruct{}
	bindings := NewBindings()
	bindings.Build(inst)

	pred, err := (&FTEQBuilder{}).Build(field, vals, nil, false)
	if err != nil {
		b.Fatal(err.Error())
	}

	b.Run("Eval", func(b *testing.B) {
		b.ReportAllocs()

		obj, _ := bindings.Bind(input)
		_ = pred.Eval(context.TODO(), obj)
	})
}

// * predicate_fteq_test.go ends here.
