// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvneq_test.go --- FVNEQ tests.
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
	testFVNEQStructType reflect.Type
	testFVNEQStructOnce sync.Once
)

// * Code:

// ** Tests:

type testFVNEQStruct struct {
	I    int
	I8   int8
	I16  int16
	I32  int32
	I64  int64
	U    uint
	U8   uint8
	U16  uint16
	U32  uint32
	U64  uint64
	F32  float32
	F64  float64
	C64  complex64
	C128 complex128
	S    string
	B    bool
	Any  any
}

func (t *testFVNEQStruct) ReflectType() reflect.Type {
	testFVNEQStructOnce.Do(func() {
		testFVNEQStructType = reflect.TypeOf(t).Elem()
	})

	return testFVNEQStructType
}

func TestFVNEQPredicate(t *testing.T) {
	input := &testFVNEQStruct{
		I:    42,
		I8:   42,
		I16:  42,
		I32:  42,
		I64:  42,
		U:    42,
		U8:   42,
		U16:  42,
		U32:  42,
		U64:  42,
		F32:  3.14159,
		F64:  3.14159,
		C64:  complex64(3.14 + 1.72i),
		C128: complex128(3.0 + 1.0i),
		S:    "Hello",
		B:    true,
		Any:  int(42),
	}

	tests := []struct {
		field string
		value any
		want  bool
	}{
		{"I", int(42), false},
		{"I", uint(95), true},
		{"I8", int8(42), false},
		{"I8", int16(42), false},
		{"I16", int8(42), false},
		{"I16", int32(94), true},
		{"I32", int64(42), false},
		{"I32", uint64(12), true},
		{"I64", int32(42), false},
		{"I64", uint32(12), true},
		{"U", uint32(42), false},
		{"U", int(-42), true},
		{"U8", uint8(42), false},
		{"U8", int8(-42), true},
		{"U16", uint8(42), false},
		{"U16", int8(-42), true},
		{"U32", uint8(42), false},
		{"U32", int8(-42), true},
		{"U64", uint64(42), false},
		{"U64", int8(-42), true},
		{"F32", float64(3.14159), false},
		{"F64", float64(3.1416), true}, // epsilon boundary?
		{"C64", complex64(3.14 + 1.72i), false},
		{"C64", complex64(3.014 + 1.70i), true},
		{"C128", complex128(3.0 + 1.0i), false},
		{"C128", complex128(3.001 + 1.0i), true},
		{"S", "hello", false}, // case-insensitive
		{"S", "world", true},
		{"B", true, false},
		{"B", false, true},
		{"Any", int(42), false},
		{"Any", float64(42), true},
		{"I64", "not an int", true}, // type mismatch
	}

	inst := &testFVNEQStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FVNEQ(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FVNEQBuilder{}).Build(tt.field, tt.value, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FVNEQ(%s == %#v) = %v, want %v",
					tt.field,
					tt.value,
					result,
					tt.want)
			}
		})
	}
}

// * predicate_fvneq_test.go ends here.
