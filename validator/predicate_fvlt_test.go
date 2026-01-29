// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvlt_test.go --- FVLT tests.
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
	testFVLTStructType reflect.Type
	testFVLTStructOnce sync.Once
)

// * Code:

// ** Types:

type testFVLTStruct struct {
	Int    int
	Uint   uint
	Float  float64
	String string
}

func (t *testFVLTStruct) ReflectType() reflect.Type {
	testFVLTStructOnce.Do(func() {
		testFVLTStructType = reflect.TypeOf(t).Elem()
	})

	return testFVLTStructType
}

// ** Tests:

func TestFVLTPredicate(t *testing.T) {
	input := &testFVLTStruct{
		Int:    int(42),
		Uint:   uint(96),
		Float:  float64(12.34),
		String: "cheese",
	}

	tests := []struct {
		field string
		less  float64
		want  bool
	}{
		{"Int", 43, true},
		{"Int", 10, false},
		{"Uint", 100, true},
		{"Uint", 30, false},
		{"Float", float64(34.56), true},
		{"Float", 10, false},
		{"String", 200, false},
	}

	inst := &testFVLTStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FVLT(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FVLTBuilder{}).Build(tt.field, tt.less, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FVLT(%s, %v) = %v, want %v",
					tt.field,
					tt.less,
					result,
					tt.want)
			}
		})
	}
}

// * predicate_fvlt_test.go ends here.
