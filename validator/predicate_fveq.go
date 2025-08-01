// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fveq.go --- FVEQ - Field Value Equality.
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
	"math"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/Asmodai/gohacks/dag"
)

// * Constants:

const (
	fveqIsn = "FVEQ"

	//nolint:gosec
	fveqToken = "field-value-equal"
)

// * Variables:

var (
	//nolint:gochecknoglobals
	wordSize uintptr

	//nolint:gochecknoglobals
	wordSizeOnce sync.Once
)

// * Code:

// ** Predicate:

type FVEQPredicate struct {
	MetaPredicate
}

func (pred *FVEQPredicate) String() string {
	val, ok := pred.MetaPredicate.GetValueAsAny()
	if !ok {
		return dag.FormatIsnf(fveqIsn, invalidTokenString)
	}

	return dag.FormatIsnf(
		fveqIsn,
		"%q %s %#v",
		pred.MetaPredicate.key,
		fveqToken,
		val,
	)
}

func (pred *FVEQPredicate) checkInt64(want, have int64) bool {
	return want == have
}

func (pred *FVEQPredicate) checkUint64(want, have uint64) bool {
	return want == have
}

func (pred *FVEQPredicate) checkFloat64(want, have float64) bool {
	const epsilon = 1e-9

	diff := math.Abs(want - have)

	return diff < epsilon
}

func (pred *FVEQPredicate) checkComplex128(want, have complex128) bool {
	rwant := pred.checkFloat64(real(want), real(have))
	rhave := pred.checkFloat64(imag(want), imag(have))

	return rwant && rhave
}

func (pred *FVEQPredicate) checkString(want, have string) bool {
	return strings.EqualFold(want, have)
}

func (pred *FVEQPredicate) checkBool(want, have bool) bool {
	return want == have
}

// Obtain the native word size of the host machine.
//
// XXX This could probably live elsewhere.
//
//nolint:mnd,gomnd
func (pred *FVEQPredicate) wordSize() uintptr {
	wordSizeOnce.Do(func() {
		//nolint:mnd
		wordSize = unsafe.Sizeof(uintptr(0)) * 8 // 8 = bits per byte
	})

	return wordSize
}

// Dispatches the right signed integer comparator for the given integer.
//
// This is hairy.  This will check the type of the given value and then
// the type of the structure field.  If the structure field can fit the
// type of the value, then a comparison is performed.
//
// The TL;DR of this is that it ensures that whatever value you're comparing
// against is either smaller or the same size as the value in the structure.
//
//nolint:cyclop,dupl,gocyclo
func (pred *FVEQPredicate) dispatchSignedInt(kind reflect.Kind, value any) bool {
	// My kingdom for a Duff's Device in Go.
	switch val := value.(type) {
	case int:
		valid := (kind == reflect.Int ||
			kind == reflect.Int8 ||
			kind == reflect.Int16 ||
			kind == reflect.Int32)

		//nolint:mnd
		if pred.wordSize() == 64 { // 64 = word size
			valid = valid || kind == reflect.Int64
		}

		have, ok := pred.MetaPredicate.GetValueAsInt64()

		return valid && ok && pred.checkInt64(have, int64(val))

	case int8:
		valid := kind == reflect.Int8
		have, ok := pred.MetaPredicate.GetValueAsInt64()

		return valid && ok && pred.checkInt64(have, int64(val))

	case int16:
		valid := (kind == reflect.Int8 ||
			kind == reflect.Int16)
		have, ok := pred.MetaPredicate.GetValueAsInt64()

		return valid && ok && pred.checkInt64(have, int64(val))
	case int32:
		valid := (kind == reflect.Int8 ||
			kind == reflect.Int16 ||
			kind == reflect.Int32)

		//nolint:mnd
		if pred.wordSize() == 32 { // 32 = word size
			valid = valid || kind == reflect.Int
		}

		have, ok := pred.MetaPredicate.GetValueAsInt64()

		return valid && ok && pred.checkInt64(have, int64(val))
	case int64:
		valid := (kind == reflect.Int8 ||
			kind == reflect.Int16 ||
			kind == reflect.Int32 ||
			kind == reflect.Int64)

		//nolint:mnd
		if pred.wordSize() == 64 { // 64 = word size
			valid = valid || kind == reflect.Int
		}

		have, ok := pred.MetaPredicate.GetValueAsInt64()

		return valid && ok && pred.checkInt64(have, val)

	default:
		return false
	}
}

// Dispatches the right unsigned integer comparator for the given integer.
//
// This is hairy.  This will check the type of the given value and then
// the type of the structure field.  If the structure field can fit the
// type of the value, then a comparison is performed.
//
// The TL;DR of this is that it ensures that whatever value you're comparing
// against is either smaller or the same size as the value in the structure.
//
//nolint:cyclop,dupl,gocyclo
func (pred *FVEQPredicate) dispatchUnsignedInt(kind reflect.Kind, value any) bool {
	// My kingdom for a Duff's Device in Go.
	switch val := value.(type) {
	case uint:
		valid := (kind == reflect.Uint ||
			kind == reflect.Uint8 ||
			kind == reflect.Uint16 ||
			kind == reflect.Uint32)

		//nolint:mnd
		if pred.wordSize() == 64 { // 64 = word size
			valid = valid || kind == reflect.Uint64
		}

		have, ok := pred.MetaPredicate.GetValueAsUint64()

		return valid && ok && pred.checkUint64(have, uint64(val))

	case uint8:
		valid := kind == reflect.Uint8
		have, ok := pred.MetaPredicate.GetValueAsUint64()

		return valid && ok && pred.checkUint64(have, uint64(val))

	case uint16:
		valid := (kind == reflect.Uint8 ||
			kind == reflect.Uint16)
		have, ok := pred.MetaPredicate.GetValueAsUint64()

		return valid && ok && pred.checkUint64(have, uint64(val))
	case uint32:
		valid := (kind == reflect.Uint8 ||
			kind == reflect.Uint16 ||
			kind == reflect.Uint32)

		//nolint:mnd
		if pred.wordSize() == 32 { // 32 = word size
			valid = valid || kind == reflect.Uint
		}

		have, ok := pred.MetaPredicate.GetValueAsUint64()

		return valid && ok && pred.checkUint64(have, uint64(val))
	case uint64:
		valid := (kind == reflect.Uint8 ||
			kind == reflect.Uint16 ||
			kind == reflect.Uint32 ||
			kind == reflect.Uint64)

		//nolint:mnd
		if pred.wordSize() == 64 { // 64 = word size
			valid = valid || kind == reflect.Uint
		}

		have, ok := pred.MetaPredicate.GetValueAsUint64()

		return valid && ok && pred.checkUint64(have, val)

	default:
		return false
	}
}

//nolint:cyclop
func (pred *FVEQPredicate) Eval(input dag.Filterable) bool {
	value, valueOk := pred.MetaPredicate.GetKeyAsValue(input)
	if !valueOk {
		return false
	}

	field, fieldOk := pred.MetaPredicate.GetKeyAsFieldInfo(input)
	if !fieldOk {
		return false
	}

	// Dispatch on the type of the value in the condition.
	switch val := value.(type) {
	case int, int8, int16, int32, int64:
		return pred.dispatchSignedInt(field.TypeKind, val)

	case uint, uint8, uint16, uint32, uint64:
		return pred.dispatchUnsignedInt(field.TypeKind, val)

	case float32:
		valid := field.TypeKind == reflect.Float32
		have, ok := pred.MetaPredicate.GetValueAsFloat64()

		return valid && ok && pred.checkFloat64(have, float64(val))

	case float64:
		valid := field.TypeKind == reflect.Float64
		have, ok := pred.MetaPredicate.GetValueAsFloat64()

		return valid && ok && pred.checkFloat64(have, float64(val))

	case complex64:
		valid := field.TypeKind == reflect.Complex64
		have, ok := pred.MetaPredicate.GetValueAsComplex128()

		return valid && ok && pred.checkComplex128(have, complex128(val))

	case complex128:
		valid := field.TypeKind == reflect.Complex128
		have, ok := pred.MetaPredicate.GetValueAsComplex128()

		return valid && ok && pred.checkComplex128(have, complex128(val))

	case string:
		valid := field.TypeKind == reflect.String
		have, ok := pred.MetaPredicate.GetValueAsString()

		return valid && ok && pred.checkString(have, val)

	case bool:
		valid := field.TypeKind == reflect.Bool
		have, ok := pred.MetaPredicate.GetValueAsBool()

		return valid && ok && pred.checkBool(have, val)

	default:
		return false
	}
}

// ** Builder:

type FVEQBuilder struct{}

func (bld *FVEQBuilder) Token() string {
	return fveqToken
}

func (bld *FVEQBuilder) Build(key string, val any) dag.Predicate {
	return &FVEQPredicate{
		MetaPredicate: MetaPredicate{key: key, val: val},
	}
}

// * predicate_fveq.go ends here.
