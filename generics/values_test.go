// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// values_test.go --- Numeric values tests.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package generics

import (
	"math"
	"reflect"
	"testing"
)

func TestHasFraction(t *testing.T) {
	t.Run("With fraction", func(t *testing.T) {
		if HasFraction(float64(7.62)) != true {
			t.Error("Unexpected result, fraction expected.")
		}
	})

	t.Run("Without fraction", func(t *testing.T) {
		if HasFraction(float64(492)) != false {
			t.Error("Unexpected result, no fraction expected")
		}
	})
}

func TestCoerceInt(t *testing.T) {
	t.Run("8-bit signed", func(t *testing.T) {
		r1 := CoerceInt(int8(math.MinInt8))
		r2 := CoerceInt(int8(math.MaxInt8))

		if typ := reflect.TypeOf(r1).Name(); typ != "int8" {
			t.Errorf("Expected type of int8, got %v", typ)
		}

		// a positive signed integer can be represented as unsigned.
		if typ := reflect.TypeOf(r2).Name(); typ != "uint8" {
			t.Errorf("Expected type of uint8, got %v", typ)
		}
	})

	t.Run("16-bit signed", func(t *testing.T) {
		r1 := CoerceInt(int16(math.MinInt16))
		r2 := CoerceInt(int16(math.MaxInt16))

		if typ := reflect.TypeOf(r1).Name(); typ != "int16" {
			t.Errorf("Expected type of int16, got %v", typ)
		}

		// a positive signed integer can be represented as unsigned.
		if typ := reflect.TypeOf(r2).Name(); typ != "uint16" {
			t.Errorf("Expected type of uint16, got %v", typ)
		}
	})

	t.Run("32-bit signed", func(t *testing.T) {
		r1 := CoerceInt(int32(math.MinInt32))
		r2 := CoerceInt(int32(math.MaxInt32))

		if typ := reflect.TypeOf(r1).Name(); typ != "int32" {
			t.Errorf("Expected type of int8, got %v", typ)
		}

		// a positive signed integer can be represented as unsigned.
		if typ := reflect.TypeOf(r2).Name(); typ != "uint32" {
			t.Errorf("Expected type of uint32, got %v", typ)
		}
	})

	t.Run("64-bit signed", func(t *testing.T) {
		r1 := CoerceInt(int64(math.MinInt64))
		r2 := CoerceInt(int64(math.MaxInt64))

		if typ := reflect.TypeOf(r1).Name(); typ != "int64" {
			t.Errorf("Expected type of int64, got %v", typ)
		}

		// a positive signed integer can be represented as unsigned.
		if typ := reflect.TypeOf(r2).Name(); typ != "uint64" {
			t.Errorf("Expected type of uint64, got %v", typ)
		}
	})
}

func TestValues(t *testing.T) {
	t.Run("32-bit float with fraction", func(t *testing.T) {
		res := ValueOf(float32(5.56))

		if typ := reflect.TypeOf(res).Name(); typ != "float32" {
			t.Errorf("Unexpected type: %v != float32", typ)
		}
	})

	t.Run("32-bit float without fraction", func(t *testing.T) {
		res := ValueOf(float32(42))

		if typ := reflect.TypeOf(res).Name(); typ != "uint8" {
			t.Errorf("Unexpected type: %v != uint8", typ)
		}
	})

	t.Run("64-bit float with fraction", func(t *testing.T) {
		res := ValueOf(float64(5.56))

		if typ := reflect.TypeOf(res).Name(); typ != "float64" {
			t.Errorf("Unexpected type: %v != float64", typ)
		}
	})

	t.Run("64-bit float without fraction", func(t *testing.T) {
		res := ValueOf(float64(42))

		if typ := reflect.TypeOf(res).Name(); typ != "uint8" {
			t.Errorf("Unexpected type: %v != uint8", typ)
		}
	})

	t.Run("8-bit signed integer", func(t *testing.T) {
		res := ValueOf(int8(math.MaxInt8))

		if typ := reflect.TypeOf(res).Name(); typ != "int8" {
			t.Errorf("Unexpected type: %v != int8", typ)
		}
	})

	t.Run("16-bit signed integer", func(t *testing.T) {
		res := ValueOf(int16(math.MaxInt16))

		if typ := reflect.TypeOf(res).Name(); typ != "int16" {
			t.Errorf("Unexpected type: %v != int16", typ)
		}
	})

	t.Run("32-bit signed integer", func(t *testing.T) {
		res := ValueOf(int32(math.MaxInt32))

		if typ := reflect.TypeOf(res).Name(); typ != "int32" {
			t.Errorf("Unexpected type: %v != int32", typ)
		}
	})

	t.Run("64-bit signed integer", func(t *testing.T) {
		res := ValueOf(int64(math.MaxInt64))

		if typ := reflect.TypeOf(res).Name(); typ != "int64" {
			t.Errorf("Unexpected type: %v != int64", typ)
		}
	})

	t.Run("8-bit unsigned integer", func(t *testing.T) {
		res := ValueOf(uint8(math.MaxUint8))

		if typ := reflect.TypeOf(res).Name(); typ != "uint8" {
			t.Errorf("Unexpected type: %v != uint8", typ)
		}
	})

	t.Run("16-bit unsigned integer", func(t *testing.T) {
		res := ValueOf(uint16(math.MaxUint16))

		if typ := reflect.TypeOf(res).Name(); typ != "uint16" {
			t.Errorf("Unexpected type: %v != uint16", typ)
		}
	})

	t.Run("32-bit unsigned integer", func(t *testing.T) {
		res := ValueOf(uint32(math.MaxUint32))

		if typ := reflect.TypeOf(res).Name(); typ != "uint32" {
			t.Errorf("Unexpected type: %v != uint32", typ)
		}
	})

	t.Run("64-bit unsigned integer", func(t *testing.T) {
		res := ValueOf(uint64(math.MaxUint64))

		if typ := reflect.TypeOf(res).Name(); typ != "uint64" {
			t.Errorf("Unexpected type: %v != uint64", typ)
		}
	})

	t.Run("String", func(t *testing.T) {
		res := ValueOf(string("seventeen"))

		if typ := reflect.TypeOf(res).Name(); typ != "string" {
			t.Errorf("Unexpected type: %v != string", typ)
		}
	})
}

// values_test.go ends here.
