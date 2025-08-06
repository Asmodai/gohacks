// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvnil.go --- FVNIL - Field Valie Is Nil.
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
	fvnilIsn   = "FVNIL"
	fvnilToken = "field-value-is-nil" //nolint:gosec
)

// * Code:

// ** Predicate:

// Field Value Is Nil.
//
// This predicate returns true if and only if the **reference** value of
// the filtered field is `nil`.
//
// If the field value is a concrete type (e.g. string, int, float, bool etc),
// then the predicate will return false.
//
// It only applies to types that can be `nil` in Go--e.g. pointers, slices,
// maps, interfaces, et al.  If you're looking to test whether a field is
// "logically nil" (e.g. zero, false, empty) then consider using `FVFALSE`
// instead.
type FVNILPredicate struct {
	MetaPredicate
}

func (pred *FVNILPredicate) String() string {
	if val, ok := pred.MetaPredicate.GetValueAsString(); ok {
		return dag.FormatIsnf(
			fvnilIsn,
			"%q %s %v",
			pred.MetaPredicate.key,
			fvnilToken,
			val)
	}

	return dag.FormatIsnf(fvnilIsn, invalidTokenString)
}

func (pred *FVNILPredicate) Eval(_ context.Context, input dag.Filterable) bool {
	finfo, finfoOk := pred.MetaPredicate.GetKeyAsFieldInfo(input)
	if !finfoOk {
		return false
	}

	// We don't care for primitive types that can't be set to `nil`.
	if finfo.Kind == KindPrimitive {
		return false
	}

	if val, valOk := pred.MetaPredicate.GetKeyAsValue(input); valOk {
		valOf := reflect.ValueOf(val)

		if finfo.TypeKind == reflect.Interface && !valOf.IsValid() {
			return false
		}

		return valOf.IsNil()
	}

	return false
}

// ** Builder:

type FVNILBuilder struct{}

func (bld *FVNILBuilder) Token() string {
	return fvnilToken
}

func (bld *FVNILBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	pred := &FVNILPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_fvnil.go ends here.
