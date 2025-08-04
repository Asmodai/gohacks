// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvin_test.go --- FVIN tests.
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
	"fmt"
	"reflect"
	"sync"
	"testing"
)

// * Constants:

// * Variables:

var (
	testFVINStructType reflect.Type
	testFVINStructOnce sync.Once
)

// * Code:

type testFVINStruct struct {
	IntVal   int
	Int64Val int64
	StrVal   string
	BoolVal  bool
	Any      any
}

func (t *testFVINStruct) ReflectType() reflect.Type {
	testFVINStructOnce.Do(func() {
		testFVINStructType = reflect.TypeOf(t).Elem()
	})

	return testFVINStructType
}

func TestFVINPredicate(t *testing.T) {
	input := &testFVINStruct{
		IntVal:   42,
		Int64Val: 9001,
		StrVal:   "Hello",
		BoolVal:  true,
		Any:      uint64(69),
	}

	tests := []struct {
		field string
		vals  []any
		want  bool
	}{
		{"IntVal", []any{1, 42, 99}, true},
		{"Int64Val", []any{1, 42, 9001}, true},
		{"StrVal", []any{"Hello", "world"}, true},
		{"Any", []any{42}, false},
		{"Any", []any{69}, false},
		{"Any", []any{1, 2, 3}, false},
		{"IntVal", []any{100, 200}, false},
	}

	inst := &testFVINStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FVIN(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FVINBuilder{}).Build(tt.field, tt.vals)
			result := pred.Eval(obj)

			if result != tt.want {
				t.Errorf("FVIN(%s IN %#v) = %v, want %v",
					tt.field,
					tt.vals,
					result,
					tt.want)
			}
		})

	}
}

// * predicate_fvin_test.go ends here.
