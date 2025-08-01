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

func (pred *FTEQPredicate) Eval(input dag.Filterable) bool {
	want, wantOk := pred.MetaPredicate.GetValueAsString()
	fInfo, fInfoOk := pred.MetaPredicate.GetKeyAsFieldInfo(input)

	if !(wantOk && fInfoOk) {
		return false
	}

	return strings.EqualFold(want, fInfo.TypeName)
}

// ** Builder:

type FTEQBuilder struct{}

func (bld *FTEQBuilder) Token() string {
	return fteqToken
}

func (bld *FTEQBuilder) Build(key string, val any) dag.Predicate {
	return &FTEQPredicate{
		MetaPredicate: MetaPredicate{key: key, val: val},
	}
}

// * predicate_fteq.go ends here.
