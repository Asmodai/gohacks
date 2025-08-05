// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate.go --- Predicates.
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
//
//mock:yes

// * Comments:

// * Package:

package dag

// * Imports:

import (
	"context"
	"fmt"
	"strings"

	"github.com/Asmodai/gohacks/conversion"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/utils"
)

// * Constants:

const (
	invalidTokenString = "<invalid!>"

	isnWidth = 8
)

// * Code:

// ** Interface:

// Predicate interface.
//
// All predicates must adhere to this interface.
type Predicate interface {
	// Evaluate the predicate against the given `Filterable` object.
	//
	// Returns the result of the predicate.
	Eval(context.Context, Filterable) bool

	// Return the string representation of the predicate.
	String() string
}

// Predicate builder interface.
//
// All predicate builders must adhere to this interface.
type PredicateBuilder interface {
	// Return the token name for the predicate.
	//
	// This isn't used in the current version of the directed acyclic
	// graph, but the theory is that this could be used in a tokeniser
	// or as opcode.
	//
	// The value this returns must be unique.
	Token() string

	// Build a new predicate.
	//
	// This will create a predicate that operates on the given field
	// and data.
	Build(field string, data any, lgr logger.Logger, dbg bool) (Predicate, error)
}

// ** Types:

// Dictionary of available predicate builders.
type PredicateDict map[string]PredicateBuilder

// *** Meta predicate:

// A `meta` predicate used by all predicates.
//
// The meta preducate presents common fields and methods so as to avoid
// duplicate code.
type MetaPredicate struct {
	key    string
	val    any
	logger logger.Logger
	debug  bool
}

// Return the predicate's input value as a 64-bit float.
//
// This will return the value for the key on which the predicate operates.
func (meta *MetaPredicate) GetFloatValueFromInput(input Filterable) (float64, bool) {
	data, dataOk := input.Get(meta.key)

	val, valOk := conversion.ToFloat64(data)

	return val, dataOk && valOk
}

// Return both the predicate's input value and filter value as a 64-bit
// float.
func (meta *MetaPredicate) GetFloatValues(input Filterable) (float64, float64, bool) {
	data, dataOk := input.Get(meta.key)

	lhs, lhsOk := conversion.ToFloat64(data)
	rhs, rhsOk := conversion.ToFloat64(meta.val)

	return lhs, rhs, dataOk && lhsOk && rhsOk
}

// Return both the predicate's input value and filter value as a string.
func (meta *MetaPredicate) GetStringValues(input Filterable) (string, string, bool) {
	data, dataOk := input.Get(meta.key)

	lhs, lhsOk := conversion.ToString(data)
	rhs, rhsOk := conversion.ToString(meta.val)

	return lhs, rhs, dataOk && lhsOk && rhsOk
}

// Return the predicate's filter value as an array of 64-bit floats.
func (meta *MetaPredicate) GetPredicateFloatArray() ([]float64, bool) {
	return conversion.AnyArrayToFloat64Array(meta.val)
}

// Return the predicate's filter value as an array of strings.
func (meta *MetaPredicate) GetPredicateStringArray() ([]string, bool) {
	return conversion.AnyArrayToStringArray(meta.val)
}

// Does the predicate's input value fall within the exclusive range defined
// in the predicate's filter value?
func (meta *MetaPredicate) EvalExclusiveRange(input Filterable) bool {
	array, arrayOk := meta.GetPredicateFloatArray()
	val, valOk := meta.GetFloatValueFromInput(input)

	if !valOk || !arrayOk {
		return false
	}

	//nolint:mnd
	if len(array) > 2 {
		return false
	}

	first := array[0]
	second := array[1]

	return (first < val) && (val < second)
}

// Does the predicate's input value fall within the inclusive range defined
// in the predicate's filter value?
//
//nolint:mnd
func (meta *MetaPredicate) EvalInclusiveRange(input Filterable) bool {
	array, arrayOk := meta.GetPredicateFloatArray()
	val, valOk := meta.GetFloatValueFromInput(input)

	if !valOk || !arrayOk {
		return false
	}

	//nolint:mnd
	if len(array) > 2 {
		return false
	}

	first := array[0]
	second := array[1]

	return (first <= val) && (val <= second)
}

// Is the predicate's input value a member of the array of strings in the
// predicate's filter value?
func (meta *MetaPredicate) EvalStringMember(input Filterable, insens bool) bool {
	valueRaw, okay := input.Get(meta.key)
	if !okay {
		return false
	}

	valueStr, okay := conversion.ToString(valueRaw)
	if !okay {
		return false
	}

	strArray, okay := meta.GetPredicateStringArray()
	if !okay {
		return false
	}

	for _, str := range strArray {
		if insens && strings.EqualFold(str, valueStr) {
			return true
		}

		if str == valueStr {
			return true
		}
	}

	return false
}

// ** Functions:

// Build the predicate dictionary for the directed acyclic graph filter.
func BuildPredicateDict() PredicateDict {
	result := make(PredicateDict)
	preds := []PredicateBuilder{
		//
		// Numeric predicates.
		&EQBuilder{}, &NEQBuilder{},
		&GTBuilder{}, &GTEBuilder{},
		&LTBuilder{}, &LTEBuilder{},
		//
		// String predicates.
		&SIEQBuilder{}, &SINEQBuilder{},
		&SSEQBuilder{}, &SSNEQBuilder{},
		//
		// Regex predicates.
		&REIMBuilder{},
		&RESMBuilder{},
	}

	for idx := range preds {
		pred := preds[idx]

		result[pred.Token()] = pred
	}

	return result
}

// Pretty-print a predicate's token.
func FormatIsnf(isn, message string, rest ...any) string {
	padded := utils.Pad(isn, isnWidth)

	return fmt.Sprintf("%s: %s", padded, fmt.Sprintf(message, rest...))
}

// * predicate.go ends here.
