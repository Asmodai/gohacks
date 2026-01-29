// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fveq_test.go --- Tests for FVEQ predicate.
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
	testFVEQStructType reflect.Type
	testFVEQStructOnce sync.Once
)

// * Code:

// ** Types:

type testFVEQStruct struct {
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

func (t *testFVEQStruct) ReflectType() reflect.Type {
	testFVEQStructOnce.Do(func() {
		testFVEQStructType = reflect.TypeOf(t).Elem()
	})

	return testFVEQStructType
}

// ** Tests:

func TestFVEQPredicate(t *testing.T) {
	input := &testFVEQStruct{
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
		{"I", int(42), true},
		{"I", uint(95), false},
		{"I8", int8(42), true},
		{"I8", int16(42), true},
		{"I16", int8(42), true},
		{"I16", int32(94), false},
		{"I32", int64(42), true},
		{"I32", uint64(12), false},
		{"I64", int32(42), true},
		{"I64", uint32(12), false},
		{"U", uint32(42), true},
		{"U", int(-42), false},
		{"U8", uint8(42), true},
		{"U8", int8(-42), false},
		{"U16", uint8(42), true},
		{"U16", int8(-42), false},
		{"U32", uint8(42), true},
		{"U32", int8(-42), false},
		{"U64", uint64(42), true},
		{"U64", int8(-42), false},
		{"F32", float64(3.14159), true},
		{"F64", float64(3.1416), false}, // epsilon boundary?
		{"C64", complex64(3.14 + 1.72i), true},
		{"C64", complex64(3.014 + 1.70i), false},
		{"C128", complex128(3.0 + 1.0i), true},
		{"C128", complex128(3.001 + 1.0i), false},
		{"S", "hello", true}, // case-insensitive
		{"S", "world", false},
		{"B", true, true},
		{"B", false, false},
		{"Any", int(42), true},
		{"Any", float64(42), false},
		{"I64", "not an int", false}, // type mismatch
	}

	inst := &testFVEQStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FVEQ(%s)", idx, tt.field), func(t *testing.T) {
			pred, _ := (&FVEQBuilder{}).Build(tt.field, tt.value, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FVEQ(%s == %#v) = %v, want %v",
					tt.field,
					tt.value,
					result,
					tt.want)
			}
		})
	}
}

// ** Benchmarks:

func BenchmarkFVEQPredicate(b *testing.B) {
	input := &testFVEQStruct{
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

	field := "C128"
	vals := complex128(3.0 + 1.0i)

	inst := &testFVEQStruct{}
	bindings := NewBindings()
	bindings.Build(inst)

	pred, err := (&FVEQBuilder{}).Build(field, vals, nil, false)
	if err != nil {
		b.Fatal(err.Error())
	}

	b.Run("Eval", func(b *testing.B) {
		b.ReportAllocs()

		obj, _ := bindings.Bind(input)
		_ = pred.Eval(context.TODO(), obj)
	})
}

// * predicate_fveq_test.go ends here.
