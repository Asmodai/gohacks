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

type Predicate interface {
	Eval(context.Context, Filterable) bool
	String() string
}

type PredicateBuilder interface {
	Token() string
	Build(string, any, logger.Logger, bool) (Predicate, error)
}

// ** Types:

type PredicateDict map[string]PredicateBuilder

// *** Meta predicate:

type MetaPredicate struct {
	key    string
	val    any
	logger logger.Logger
	debug  bool
}

func (meta *MetaPredicate) GetFloatValueFromInput(input Filterable) (float64, bool) {
	data, dataOk := input.Get(meta.key)

	val, valOk := conversion.ToFloat64(data)

	return val, dataOk && valOk
}

func (meta *MetaPredicate) GetFloatValues(input Filterable) (float64, float64, bool) {
	data, dataOk := input.Get(meta.key)

	lhs, lhsOk := conversion.ToFloat64(data)
	rhs, rhsOk := conversion.ToFloat64(meta.val)

	return lhs, rhs, dataOk && lhsOk && rhsOk
}

func (meta *MetaPredicate) GetStringValues(input Filterable) (string, string, bool) {
	data, dataOk := input.Get(meta.key)

	lhs, lhsOk := conversion.ToString(data)
	rhs, rhsOk := conversion.ToString(meta.val)

	return lhs, rhs, dataOk && lhsOk && rhsOk
}

func (meta *MetaPredicate) GetPredicateFloatArray() ([]float64, bool) {
	return conversion.AnyArrayToFloat64Array(meta.val)
}

func (meta *MetaPredicate) GetPredicateStringArray() ([]string, bool) {
	return conversion.AnyArrayToStringArray(meta.val)
}

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

func FormatIsnf(isn, message string, rest ...any) string {
	padded := utils.Pad(isn, isnWidth)

	return fmt.Sprintf("%s: %s", padded, fmt.Sprintf(message, rest...))
}

// * predicate.go ends here.
