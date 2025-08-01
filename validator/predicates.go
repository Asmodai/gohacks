// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicates.go --- Predicates.
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
	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/math/conversion"
)

// * Constants:

const (
	invalidTokenString = "<invalid!>"
)

// * Variables:

// * Code:

// ** Types:

type MetaPredicate struct {
	key string
	val any
}

// ** Methods:

//
// XXX
//
// Need a way of looking up live structure values.  This might need us to
// tie in the live structure with its FieldInfo, or StructureDecription
//
// Maybe via pointers (or references) or whatever.
//

func (meta *MetaPredicate) GetValueAsAny() (any, bool) {
	return meta.val, true
}

// Return the condition value as a string.
func (meta *MetaPredicate) GetValueAsString() (string, bool) {
	return conversion.ToString(meta.val)
}

func (meta *MetaPredicate) GetValueAsInt64() (int64, bool) {
	return conversion.ToInt64(meta.val)
}

func (meta *MetaPredicate) GetValueAsUint64() (uint64, bool) {
	return conversion.ToUint64(meta.val)
}

func (meta *MetaPredicate) GetValueAsFloat64() (float64, bool) {
	return conversion.ToFloat64(meta.val)
}

func (meta *MetaPredicate) GetValueAsBool() (bool, bool) {
	return conversion.ToBool(meta.val)
}

func (meta *MetaPredicate) GetValueAsComplex128() (complex128, bool) {
	return conversion.ToComplex128(meta.val)
}

// Return the `Filterable`'s field information.
//
// This is directed through to `BoundObject.Description.Fields`.
func (meta *MetaPredicate) GetKeyAsFieldInfo(input dag.Filterable) (*FieldInfo, bool) {
	// `BoundObject.Get` will return `map[string]*FieldInfo`.
	data, dataOk := input.Get(meta.key)
	finfo, finfoOk := data.(*FieldInfo)

	return finfo, dataOk && finfoOk
}

// Get the value from the `Filterable`.
//
// This equates to `BoundObject.Descriptor.Field[key].Accessor` being called
// with `BoundObject.Binding`.
//
// See `BoundObject.GetValue` for more.
func (meta *MetaPredicate) GetKeyAsValue(input dag.Filterable) (any, bool) {
	bound, boundOk := input.(*BoundObject)
	if !boundOk {
		return nil, false
	}

	return bound.GetValue(meta.key)
}

func (meta *MetaPredicate) GetKeyAsString(input dag.Filterable) (string, bool) {
	data, dataOk := input.Get(meta.key)
	str, strOk := conversion.ToString(data)

	return str, dataOk && strOk
}

// ** Functions:

func BuildPredicateDict() dag.PredicateDict {
	result := make(dag.PredicateDict)
	preds := []dag.PredicateBuilder{
		&FTEQBuilder{},
		&FVEQBuilder{},
	}

	for idx := range preds {
		pred := preds[idx]

		result[pred.Token()] = pred
	}

	return result
}

// * predicates.go ends here.
