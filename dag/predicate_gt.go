// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_gt.go --- GT - Numeric Greater-Than.
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

package dag

// * Imports:

import "github.com/Asmodai/gohacks/math/conversion"

// * Constants:

const (
	gtIsn   = "GT"
	gtToken = ">"
)

// * Code:

// ** Predicate:

type GTPredicate struct {
	MetaPredicate
}

func (pred *GTPredicate) String() string {
	val, ok := conversion.ToFloat64(pred.MetaPredicate.val)
	if !ok {
		return FormatIsnf(gtIsn, invalidTokenString)
	}

	return FormatIsnf(gtIsn, "%s %s %g", pred.MetaPredicate.key, gtToken, val)
}

func (pred *GTPredicate) Eval(input Filterable) bool {
	lhs, rhs, ok := pred.MetaPredicate.GetFloatValues(input)

	return ok && lhs > rhs
}

// ** Builder:

type GTBuilder struct{}

func (bld *GTBuilder) Token() string {
	return gtToken
}

func (bld *GTBuilder) Build(key string, val any) Predicate {
	return &GTPredicate{
		MetaPredicate: MetaPredicate{key: key, val: val},
	}
}

// * predicate_gt.go ends here.
