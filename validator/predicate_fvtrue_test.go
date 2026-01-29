// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvtrue_test.go --- FVTRUE tests.
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
	testFVTRUEStructType reflect.Type
	testFVTRUEStructOnce sync.Once
)

// * Code:

// ** Types:

type testFVTRUESubStruct struct {
	Unused bool
}

type testFVTRUEStruct struct {
	Primitive1 int
	Primitive2 int64
	String1    string
	String2    string
	Map1       map[int]string
	Map2       map[int]string
	Slice1     []int
	Slice2     []int
	Interface1 any
	Interface2 any
	Interface3 any
	Interface4 any
	Interface5 any
	Interface6 any
	Struct1    *testFVTRUESubStruct
}

func (t *testFVTRUEStruct) ReflectType() reflect.Type {
	testFVTRUEStructOnce.Do(func() {
		testFVTRUEStructType = reflect.TypeOf(t).Elem()
	})

	return testFVTRUEStructType
}

// ** Tests:

func TestFVTRUEPredicate(t *testing.T) {
	input := &testFVTRUEStruct{
		Primitive1: 42,
		Primitive2: 0,
		String1:    "Hello",
		String2:    "",
		Map1:       map[int]string{1: "Hello", 2: "World"},
		Map2:       map[int]string{},
		Slice1:     []int{1, 2, 3},
		Slice2:     []int{},
		Interface1: uint64(42),
		Interface2: nil,
		Interface3: []int{1, 2, 3},
		Interface4: map[int]string{},
		Interface5: &testFVTRUESubStruct{},
		Interface6: testFVFALSESubStruct{},
		Struct1:    &testFVTRUESubStruct{},
	}

	tests := []struct {
		field string
		want  bool
	}{
		{"Primitive1", true},
		{"Primitive2", false},
		{"String1", true},
		{"String2", false},
		{"Map1", true},
		{"Map2", false},
		{"Slice1", true},
		{"Slice2", false},
		{"Interface1", true},
		{"Interface2", false},
		{"Interface3", true},
		{"Interface4", false},
		{"Interface5", false},
		{"Interface6", false},
		{"Struct1", false},
	}

	inst := &testFVTRUEStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FVTRUE(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FVTRUEBuilder{}).Build(tt.field, nil, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FVTRUE(%s) = %v, want %v",
					tt.field,
					result,
					tt.want)
			}
		})
	}
}

// * predicate_fvtrue_test.go ends here.
