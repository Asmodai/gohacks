// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvfalse.go --- FVFALSE - Field Value is Logically False.
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
	"context"
	"reflect"

	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	fvfalseIsn   = "FVFALSE"
	fvfalseToken = "field-value-is-false" //nolint:gosec
)

// * Code:

// ** Predicate:

// Field Value is Logically False.
//
// This predicate returns true if the value of the filtered field is
// logically false.
//
// A logical false value is any value that is empty or zero.  The following
// are examples of this:
//
//	"" string, 0 numeric, [] array
//
// Logical falsehood is not the same as `nil`, so if you are looking for
// nil values then you should look at `FVNIL` instead.
//
// Structures are a special case.  They are never logically false.  This is
// because the validator does not recurse into structures.  If you wish to
// deal with structures within structures, then those sub-structures require
// validation by themselves.  How you do that is up to you.
//
// Interfaces are also a special case.  An interface can be considered
// logically false if it is `nil`, but it can also be considered logically
// false if the wrapped value is zero or empty.
type FVFALSEPredicate struct {
	MetaPredicate
}

func (pred *FVFALSEPredicate) String() string {
	if val, ok := pred.MetaPredicate.GetValueAsString(); ok {
		return dag.FormatIsnf(
			fvfalseIsn,
			"%q %s %v",
			pred.MetaPredicate.key,
			fvfalseToken,
			val)
	}

	return dag.FormatIsnf(fvfalseIsn, invalidTokenString)
}

//nolint:cyclop,exhaustive
func (pred *FVFALSEPredicate) Eval(_ context.Context, input dag.Filterable) bool {
	finfo, finfoOk := pred.MetaPredicate.GetKeyAsFieldInfo(input)
	if !finfoOk {
		return false
	}

	// Structure... never false.
	if finfo.TypeKind == reflect.Struct {
		return false
	}

	val, valOk := pred.MetaPredicate.GetKeyAsValue(input)
	if !valOk {
		return false
	}

	valOf := reflect.ValueOf(val)
	if !valOf.IsValid() {
		return true // Not valid, thus false.
	}

	switch finfo.Kind {
	case KindPrimitive:
		return valOf.IsZero()

	case KindSlice, KindMap, KindInterface:
		switch valOf.Kind() {
		case reflect.Struct:
			return false

		case reflect.Map, reflect.Slice:
			return valOf.IsZero() || valOf.Len() == 0

		default:
			return valOf.IsZero()
		}

	default:
		return false
	}
}

// ** Builder:

type FVFALSEBuilder struct{}

func (bld *FVFALSEBuilder) Token() string {
	return fvfalseToken
}

func (bld *FVFALSEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	pred := &FVFALSEPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_fvfalse.go ends here.
