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
	//nolint:gochecknoglobals,gomnd,mnd
	intSizeMap = map[reflect.Kind]int{
		reflect.Int8:      8,
		reflect.Int16:     16,
		reflect.Int32:     32,
		reflect.Int64:     64,
		reflect.Int:       wordSize,
		reflect.Uint8:     8,
		reflect.Uint16:    16,
		reflect.Uint32:    32,
		reflect.Uint64:    64,
		reflect.Uint:      wordSize,
		reflect.Interface: -1,
	}
)

// * Code:

// ** Predicate:

// Field Value Equality.
//
// This predicate compares the value to that in the structure.  If they
// are equal then the predicate returns true.
//
// The predicate will take various circumstances into consideration while
// checking the value:
//
// If the field is `any` then the comparison will match just the type of
// the value rather than using the type of the field along with the value.
//
// If the field is integer, then the structure's field must have a bid
// width large enough to hold the value.
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
		return dispatchInt(pred, field.TypeKind, val)

	case uint, uint8, uint16, uint32, uint64:
		return dispatchUint(pred, field.TypeKind, val)

	case float32:
		return checkFloat(pred,
			field.TypeKind,
			reflect.Float32,
			float64(val))

	case float64:
		return checkFloat(pred,
			field.TypeKind,
			reflect.Float64,
			val)

	case complex64:
		return checkComplex(pred,
			field.TypeKind,
			reflect.Complex64,
			complex128(val))

	case complex128:
		return checkComplex(pred,
			field.TypeKind,
			reflect.Complex128,
			val)
	case string:
		return checkString(pred,
			field.TypeKind,
			reflect.String,
			val)

	case bool:
		return checkBool(pred,
			field.TypeKind,
			reflect.Bool,
			val)

	default:
		return false
	}
}

// ** Builder:

type FVEQBuilder struct{}

func (bld *FVEQBuilder) Token() string {
	return fveqToken
}

func (bld *FVEQBuilder) Build(key string, val any) (dag.Predicate, error) {
	pred := &FVEQPredicate{
		MetaPredicate: MetaPredicate{key: key, val: val},
	}

	return pred, nil
}

// ** Functions:

// Return true if two `float32` values are approximately equal, using
// relative error comparison.
//
// This method avoids false negatives caused by `float32` rounding and scale.
// The comparison is tolerant to small relative differences, but should catch
// genuinely different values.
func compareFloat32(want, have float32) bool {
	const epsilon = 1e-6

	diff := math.Abs(float64(want) - float64(have))
	maxAbs := math.Max(math.Abs(float64(want)), math.Abs(float64(have)))

	return diff < epsilon*maxAbs
}

// Returns true if two `float64` values are approximately equal, using both
// absolute error and a ULP-based "next representable value" check.
//
// This comparison allows for very small absolute difference (under epsilon)
// and also considers values that differ by just one floating-point step.
//
// It's suitable for high-precision float comparisons where minor rounding
// differences are expected.
func compareFloat64(want, have float64) bool {
	const epsilon = 1e-9

	diff := math.Abs(want - have)

	return math.Nextafter(want, have) == have || diff < epsilon
}

// Compare two 64 bit complex numbers.
func compareComplex64(want, have complex64) bool {
	rwant := compareFloat32(real(want), real(have))
	rhave := compareFloat32(imag(want), imag(have))

	return rwant && rhave
}

// Compare two 128 bit complex numbers.
func compareComplex128(want, have complex128) bool {
	rwant := compareFloat64(real(want), real(have))
	rhave := compareFloat64(imag(want), imag(have))

	return rwant && rhave
}

// Check that the variable's declared type is big enough for a given
// integer type.
func checkBitWidth(kind reflect.Kind, bits int) bool {
	width, ok := intSizeMap[kind]
	if !ok {
		return false
	}

	// If width == -1, then it is `any`.
	return width == -1 || width >= bits
}

// Ensure that the underlying type is the same as the given type.
//
// This will also match fields of type `any`.
func checkType(kind, want reflect.Kind) bool {
	return kind == reflect.Interface || kind == want
}

// Perform a check on a boolean.
func checkBool(pred *FVEQPredicate, kind, want reflect.Kind, value bool) bool {
	if !checkType(kind, want) {
		return false
	}

	have, ok := pred.MetaPredicate.GetValueAsBool()
	if !ok {
		return false
	}

	return have == value
}

// Perform a check on an string.
func checkString(pred *FVEQPredicate, kind, want reflect.Kind, value string) bool {
	if !checkType(kind, want) {
		return false
	}

	have, ok := pred.MetaPredicate.GetValueAsString()
	if !ok {
		return false
	}

	return strings.EqualFold(have, value)
}

// Perform a check on a floating-point number.
func checkFloat(pred *FVEQPredicate, kind, want reflect.Kind, value float64) bool {
	if !checkType(kind, want) {
		return false
	}

	have, ok := pred.MetaPredicate.GetValueAsFloat64()
	if !ok {
		return false
	}

	if want == reflect.Float32 {
		return compareFloat32(float32(have), float32(value))
	}

	return compareFloat64(have, value)
}

// Perform a check on a complex number.
func checkComplex(pred *FVEQPredicate, kind, want reflect.Kind, value complex128) bool {
	if !checkType(kind, want) {
		return false
	}

	have, ok := pred.MetaPredicate.GetValueAsComplex128()
	if !ok {
		return false
	}

	if want == reflect.Complex64 {
		return compareComplex64(complex64(have), complex64(value))
	}

	return compareComplex128(have, value)
}

// Perform a check on the unsigned integer value.
//
// This checks that the type of the structure's field is big enough.
func checkUint(pred *FVEQPredicate, kind reflect.Kind, value uint64, bits int) bool {
	if !checkBitWidth(kind, bits) {
		return false
	}

	have, ok := pred.MetaPredicate.GetValueAsUint64()
	if !ok {
		return false
	}

	return have == value
}

// Perform a check on the signed integer value.
//
// This checks that the type of the structure's field is big enough.
func checkInt(pred *FVEQPredicate, kind reflect.Kind, value int64, bits int) bool {
	if !checkBitWidth(kind, bits) {
		return false
	}

	have, ok := pred.MetaPredicate.GetValueAsInt64()
	if !ok {
		return false
	}

	return have == value
}

// Dispatch on signed integer type and check the value accordingly.
//
//nolint:mnd,gomnd
func dispatchInt(pred *FVEQPredicate, kind reflect.Kind, value any) bool {
	switch val := value.(type) {
	case int:
		return checkInt(pred, kind, int64(val), wordSize)

	case int8:
		return checkInt(pred, kind, int64(val), 8)

	case int16:
		return checkInt(pred, kind, int64(val), 16)

	case int32:
		return checkInt(pred, kind, int64(val), 32)

	case int64:
		return checkInt(pred, kind, val, 64)

	default:
		return false
	}
}

// Dispatch on unsigned integer type and check the value accordingly.
//
//nolint:mnd,gomnd
func dispatchUint(pred *FVEQPredicate, kind reflect.Kind, value any) bool {
	switch val := value.(type) {
	case uint:
		return checkUint(pred, kind, uint64(val), wordSize)

	case uint8:
		return checkUint(pred, kind, uint64(val), 8)

	case uint16:
		return checkUint(pred, kind, uint64(val), 16)

	case uint32:
		return checkUint(pred, kind, uint64(val), 32)

	case uint64:
		return checkUint(pred, kind, val, 64)

	default:
		return false
	}
}

// * predicate_fveq.go ends here.
