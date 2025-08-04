// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvin.go --- Field Value In.
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
	"github.com/Asmodai/gohacks/conversion"
	"github.com/Asmodai/gohacks/dag"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	fvinIsn = "FVIN"

	fvinToken = "field-value-in"
)

// * Variables:

var (
	ErrInvalidSlice       = errors.Base("invalid slice")
	ErrNotCanonicalisable = errors.Base("value cannot be canonicalised")
	ErrNotComparable      = errors.Base("value cannot be compared")
)

// * Code:

// ** Predicate:

// Field Value In.
//
// This predicate returns true if the value in the structure is one of the
// provided values in the predicate.
type FVINPredicate struct {
	MetaPredicate

	valueSet map[any]struct{}
}

func (pred *FVINPredicate) String() string {
	val, ok := pred.MetaPredicate.GetValueAsAny()
	if !ok {
		return dag.FormatIsnf(fvinIsn, invalidTokenString)
	}

	return dag.FormatIsnf(
		fvinIsn,
		"%q %s %v",
		pred.MetaPredicate.key,
		fvinToken,
		val,
	)
}

func (pred *FVINPredicate) Eval(input dag.Filterable) bool {
	have, haveOk := pred.MetaPredicate.GetKeyAsValue(input)
	if !haveOk {
		return false
	}

	canon, canonOk := conversion.Canonicalise(have)
	if !canonOk {
		return false
	}

	_, exists := pred.valueSet[canon]

	return exists
}

// ** Builder:

type FVINBuilder struct{}

func (bld *FVINBuilder) Token() string {
	return fvinToken
}

func (bld *FVINBuilder) Build(key string, val any) (dag.Predicate, error) {
	valSlice, sliceOk := val.([]any)
	if !sliceOk {
		return nil, errors.WithMessagef(
			ErrInvalidSlice,
			"%s: syntax error",
			fvinToken)
	}

	valMap := make(map[any]struct{}, len(valSlice))

	for _, elt := range valSlice {
		if !isComparable(elt) {
			return nil, errors.WithMessagef(
				ErrNotComparable,
				"%s: value %q",
				fvinToken,
				elt)
		}

		canon, canonOk := conversion.Canonicalise(elt)
		if !canonOk {
			return nil, errors.WithMessagef(
				ErrNotCanonicalisable,
				"%s: value %q",
				fvinToken,
				elt)
		}

		valMap[canon] = struct{}{}
	}

	pred := &FVINPredicate{
		MetaPredicate: MetaPredicate{key: key, val: val},
		valueSet:      valMap,
	}

	return pred, nil
}

// ** Functions:

//nolint:errcheck
func isComparable(val any) bool {
	defer func() {
		recover()
	}()

	_ = map[any]struct{}{val: {}}

	return true
}

// * predicate_fvin.go ends here.
