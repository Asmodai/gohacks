// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_ftin.go --- FTIN - Field Type In.
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
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	ftinIsn   = "FTIN"
	ftinToken = "field-type-in"
)

// * Code:

// ** Predicate:

// Field Type In.
//
// This predicate returns true of the type of a field in the input structure
// is one of the provided values in the predicate.
type FTINPredicate struct {
	valueSet map[string]struct{}

	MetaPredicate
}

func (pred *FTINPredicate) Instruction() string {
	return ftinIsn
}

func (pred *FTINPredicate) Token() string {
	return ftinToken
}

func (pred *FTINPredicate) String() string {
	return pred.MetaPredicate.String(ftinToken)
}

func (pred *FTINPredicate) Debug() string {
	return pred.MetaPredicate.Debug(ftinIsn, ftinToken)
}

func (pred *FTINPredicate) Eval(_ context.Context, input dag.Filterable) bool {
	have, haveOk := pred.MetaPredicate.GetKeyAsFieldInfo(input)
	if !haveOk {
		return false
	}

	typeName := have.TypeName

	// If the field is `any`, then look at the actual value of the data.
	if have.TypeKind == reflect.Interface {
		if _, haveAny := pred.valueSet["interface {}"]; haveAny {
			// We match `any`, so we don't need to do more.
			return true
		}

		// Get the real value.
		value, valueOk := pred.MetaPredicate.GetKeyAsValue(input)
		if !valueOk {
			return false
		}

		// Resolve it.
		actual, actualOk := resolveAnyType(value)
		if !actualOk {
			return false
		}

		// Use it.
		typeName = actual
	}

	_, exists := pred.valueSet[typeName]

	return exists
}

// ** Builder:

type FTINBuilder struct{}

func (bld *FTINBuilder) Token() string {
	return ftinToken
}

func (bld *FTINBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	valSlice, sliceOk := val.([]any)
	if !sliceOk {
		return nil, errors.WithMessagef(
			ErrInvalidSlice,
			"%s: syntax error",
			ftinToken)
	}

	valMap := make(map[string]struct{}, len(valSlice))

	// Iterate over all the values and coerce them to string.
	for _, elt := range valSlice {
		// NOTE: We are not using `conversion.ToString()` here as
		// we don't want to coerce non-strings to string.
		typeName, converted := elt.(string)
		if !converted {
			return nil, errors.WithMessagef(
				ErrValueNotString,
				"%s: value %q",
				ftinToken,
				elt)
		}

		typeName = normaliseTypeName(typeName)
		valMap[typeName] = struct{}{}
	}

	pred := &FTINPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
		valueSet: valMap,
	}

	return pred, nil
}

// * predicate_ftin.go ends here.
