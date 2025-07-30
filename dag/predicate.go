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
	"fmt"

	"github.com/Asmodai/gohacks/math/conversion"
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
	Eval(DataMap) bool
	String() string
}

type PredicateBuilder interface {
	Token() string
	Build(string, any) Predicate
}

// ** Types:

type PredicateDict map[string]PredicateBuilder

// *** Meta predicate:

type MetaPredicate struct {
	key string
	val any
}

func (meta *MetaPredicate) GetFloatValues(input DataMap) (float64, float64, bool) {
	data, dataOk := input[meta.key]

	lhs, lhsOk := conversion.ToFloat64(data)
	rhs, rhsOk := conversion.ToFloat64(meta.val)

	return lhs, rhs, dataOk && lhsOk && rhsOk
}

func (meta *MetaPredicate) GetStringValues(input DataMap) (string, string, bool) {
	data, dataOk := input[meta.key]

	lhs, lhsOk := data.(string)
	rhs, rhsOk := meta.val.(string)

	return lhs, rhs, dataOk && lhsOk && rhsOk
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
