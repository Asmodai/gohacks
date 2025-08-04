// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fteq.go --- FTEQ - Field Type Equals.
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
	"reflect"
	"strings"

	"github.com/Asmodai/gohacks/dag"
)

// * Constants:

const (
	fteqIsn   = "FTEQ"
	fteqToken = "field-type-equal"
)

// * Code:

// ** Predicate:

// Field Type Equality.
//
// This predicate compares the type of the structure's field.  If it is
// equal then the predicate returns true.
type FTEQPredicate struct {
	MetaPredicate
}

func (pred *FTEQPredicate) String() string {
	val, ok := pred.MetaPredicate.GetValueAsString()
	if !ok {
		return dag.FormatIsnf(fteqIsn, invalidTokenString)
	}

	return dag.FormatIsnf(
		fteqIsn,
		"%q %s %q",
		pred.MetaPredicate.key,
		fteqToken,
		val,
	)
}

func (pred *FTEQPredicate) checkSigned(value any, want string) bool {
	switch value.(type) {
	case int:
		return strings.EqualFold(want, "int")

	case int8:
		return strings.EqualFold(want, "int8")

	case int16:
		return strings.EqualFold(want, "int16")

	case int32:
		return strings.EqualFold(want, "int32")

	case int64:
		return strings.EqualFold(want, "int64")

	default:
		return false
	}
}

func (pred *FTEQPredicate) checkUnsigned(value any, want string) bool {
	switch value.(type) {
	case uint:
		return strings.EqualFold(want, "uint")

	case uint8:
		return strings.EqualFold(want, "uint8")

	case uint16:
		return strings.EqualFold(want, "uint16")

	case uint32:
		return strings.EqualFold(want, "uint32")

	case uint64:
		return strings.EqualFold(want, "uint64")

	default:
		return false
	}
}

func (pred *FTEQPredicate) checkFloat(value any, want string) bool {
	switch value.(type) {
	case float32:
		return strings.EqualFold(want, "float32")

	case float64:
		return strings.EqualFold(want, "float64")
	default:
		return false
	}
}

func (pred *FTEQPredicate) checkComplex(value any, want string) bool {
	switch value.(type) {
	case complex64:
		return strings.EqualFold(want, "complex64")

	case complex128:
		return strings.EqualFold(want, "complex128")
	default:
		return false
	}
}

//nolint:cyclop
func (pred *FTEQPredicate) resolveAny(value any, want string) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return pred.checkSigned(value, want)

	case uint, uint8, uint16, uint32, uint64:
		return pred.checkUnsigned(value, want)

	case float32, float64:
		return pred.checkFloat(value, want)

	case complex64, complex128:
		return pred.checkComplex(value, want)

	case bool:
		return strings.EqualFold(want, "bool")

	case string:
		return strings.EqualFold(want, "string")

	case []byte:
		return strings.EqualFold(want, "[]byte")

	case []any:
		return strings.EqualFold(want, "[]any")

	case any:
		return strings.EqualFold(want, "any")

	default:
		return false
	}
}

func (pred *FTEQPredicate) Eval(input dag.Filterable) bool {
	want, wantOk := pred.MetaPredicate.GetValueAsString()
	fInfo, fInfoOk := pred.MetaPredicate.GetKeyAsFieldInfo(input)

	if !(wantOk && fInfoOk) {
		return false
	}

	if fInfo.TypeKind == reflect.Interface {
		val, valok := pred.MetaPredicate.GetKeyAsValue(input)

		if valok {
			return pred.resolveAny(val, want)
		}
	}

	return strings.EqualFold(want, fInfo.TypeName)
}

// ** Builder:

type FTEQBuilder struct{}

func (bld *FTEQBuilder) Token() string {
	return fteqToken
}

func (bld *FTEQBuilder) Build(key string, val any) (dag.Predicate, error) {
	pred := &FTEQPredicate{
		MetaPredicate: MetaPredicate{key: key, val: val},
	}

	return pred, nil
}

// * predicate_fteq.go ends here.
