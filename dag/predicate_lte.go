// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_lte.go --- LTE - Numeric Less-Than-or-Equal-To.
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

//nolint:dupl
package dag

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/conversion"
	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	lteIsn   = "LTE"
	lteToken = "<="
)

// * Variables:

// * Code:

// ** Predicate:

type LTEPredicate struct {
	MetaPredicate
}

func (pred *LTEPredicate) String() string {
	if val, ok := conversion.ToFloat64(pred.MetaPredicate.val); ok {
		return FormatIsnf(lteIsn,
			"%s %s %g",
			pred.MetaPredicate.key,
			lteToken,
			val)
	}

	return FormatIsnf(lteIsn, invalidTokenString)
}

func (pred *LTEPredicate) Eval(_ context.Context, input Filterable) bool {
	lhs, rhs, ok := pred.MetaPredicate.GetFloatValues(input)

	return ok && lhs <= rhs
}

// ** Builder:

type LTEBuilder struct{}

func (bld *LTEBuilder) Token() string {
	return lteToken
}

func (bld *LTEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error) {
	pred := &LTEPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_lte.go ends here.
