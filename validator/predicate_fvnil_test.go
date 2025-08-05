// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvnil_test.go --- FVNIL tests.
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
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
)

// * Variables:

var (
	testFVNILStructType reflect.Type
	testFVNILStructOnce sync.Once
)

// * Code:

// ** Types:

type testFVNILSubStruct struct {
	Unused bool
}

type testFVNILStruct struct {
	Primitive int
	Map1      map[int]string
	Map2      map[int]string
	Slice1    []int
	Slice2    []int
	Struct1   *testFVNILSubStruct
	Struct2   *testFVNILSubStruct
}

func (t *testFVNILStruct) ReflectType() reflect.Type {
	testFVNILStructOnce.Do(func() {
		testFVNILStructType = reflect.TypeOf(t).Elem()
	})

	return testFVNILStructType
}

// ** Tests:

func TestFVNILPredicate(t *testing.T) {
	input := &testFVNILStruct{
		Primitive: 42,
		Map1:      map[int]string{},
		Map2:      nil,
		Slice1:    []int{},
		Slice2:    nil,
		Struct1:   &testFVNILSubStruct{},
		Struct2:   nil,
	}

	tests := []struct {
		field string
		want  bool
	}{
		{"Primitive", false},
		{"Map1", false},
		{"Map2", true},
		{"Slice1", false},
		{"Slice2", true},
		{"Struct1", false},
		{"Struct2", true},
	}

	inst := &testFVNILStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FVNIL(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FVNILBuilder{}).Build(tt.field, nil, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FVNIL(%s) = %v, want %v",
					tt.field,
					result,
					tt.want)
			}
		})
	}
}

// ** Benchmarks:

func BenchmarkFVNILPredicate(b *testing.B) {
	input := &testFVNILStruct{
		Primitive: 42,
		Map1:      map[int]string{},
		Map2:      nil,
		Slice1:    []int{},
		Slice2:    nil,
		Struct1:   &testFVNILSubStruct{},
		Struct2:   nil,
	}

	field := "Struct2"

	inst := &testFVNILStruct{}
	bindings := NewBindings()
	bindings.Build(inst)

	pred, err := (&FVNILBuilder{}).Build(field, nil, nil, false)
	if err != nil {
		b.Fatal(err.Error())
	}

	b.Run("Eval", func(b *testing.B) {
		b.ReportAllocs()

		obj, _ := bindings.Bind(input)
		_ = pred.Eval(context.TODO(), obj)
	})
}

// * predicate_fvnil_test.go ends here.
