// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvtrue.go --- FVTRUE - Field Value is Logically True.
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
	fvtrueIsn   = "FVTRUE"
	fvtrueToken = "field-value-is-true" //nolint:gosec
)

// * Code:

// ** Predicate:

// Field Value is Logically True.
//
// This predicate returns true if the value of the filtered field is
// logically true.
//
// A logical true value is any value that is not empty or zero.
//
// For more details on how this works, see `FVFALSE`.
type FVTRUEPredicate struct {
	FVFALSEPredicate
}

func (pred *FVTRUEPredicate) String() string {
	if val, ok := pred.MetaPredicate.GetValueAsString(); ok {
		return dag.FormatIsnf(
			fvtrueIsn,
			"%q %s %v",
			pred.MetaPredicate.key,
			fvtrueToken,
			val)
	}

	return dag.FormatIsnf(fvtrueIsn, invalidTokenString)
}

func (pred *FVTRUEPredicate) Eval(ctx context.Context, input dag.Filterable) bool {
	finfo, finfoOk := pred.MetaPredicate.GetKeyAsFieldInfo(input)
	if !finfoOk {
		return false
	}

	// Structure... never true.
	if finfo.Kind == KindStruct || finfo.TypeKind == reflect.Struct {
		return false
	}

	// We also need to check the value to make sure it's not a
	// structure.
	val, valOk := pred.MetaPredicate.GetKeyAsValue(input)
	if !valOk {
		return false
	}

	kindOf := reflect.ValueOf(val).Kind()
	if kindOf == reflect.Struct || kindOf == reflect.Ptr {
		return false
	}

	// Simply negate the result of FVFALSE now.
	return !pred.FVFALSEPredicate.Eval(ctx, input)
}

// ** Builder:

type FVTRUEBuilder struct{}

func (bld *FVTRUEBuilder) Token() string {
	return fvtrueToken
}

func (bld *FVTRUEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	pred := &FVTRUEPredicate{
		FVFALSEPredicate: FVFALSEPredicate{
			MetaPredicate: MetaPredicate{
				key:    key,
				val:    val,
				logger: lgr,
				debug:  dbg,
			},
		},
	}

	return pred, nil
}

// * predicate_fvtrue.go ends here.
