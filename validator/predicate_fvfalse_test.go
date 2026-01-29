// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvfalse_test.go --- FVFALSE tests.
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
	testFVFALSEStructType reflect.Type
	testFVFALSEStructOnce sync.Once
)

// * Code:

// ** Types:

type testFVFALSESubStruct struct {
	Unused bool
}

type testFVFALSEStruct struct {
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
	Struct1    *testFVFALSESubStruct
}

func (t *testFVFALSEStruct) ReflectType() reflect.Type {
	testFVFALSEStructOnce.Do(func() {
		testFVFALSEStructType = reflect.TypeOf(t).Elem()
	})

	return testFVFALSEStructType
}

// ** Tests:

func TestFVFALSEPredicate(t *testing.T) {
	input := &testFVFALSEStruct{
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
		Interface5: &testFVFALSESubStruct{},
		Interface6: testFVFALSESubStruct{},
		Struct1:    &testFVFALSESubStruct{},
	}

	tests := []struct {
		field string
		want  bool
	}{
		{"Primitive1", false},
		{"Primitive2", true},
		{"String1", false},
		{"String2", true},
		{"Map1", false},
		{"Map2", true},
		{"Slice1", false},
		{"Slice2", true},
		{"Interface1", false},
		{"Interface2", true},
		{"Interface3", false},
		{"Interface4", true},
		{"Interface5", false},
		{"Interface6", false},
		{"Struct1", false},
	}

	inst := &testFVFALSEStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FVFALSE(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FVFALSEBuilder{}).Build(tt.field, nil, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FVFALSE(%s) = %v, want %v",
					tt.field,
					result,
					tt.want)
			}
		})
	}
}

// * predicate_fvfalse_test.go ends here.
