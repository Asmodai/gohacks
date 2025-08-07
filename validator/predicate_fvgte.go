// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvgte.go --- FVGTE - Field Value is Lesser Than.
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
package validator

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	fvgteIsn   = "FVGTE"
	fvgteToken = "field-value-<"
)

// * Code:

// ** Predicate:

type FVGTEPredicate struct {
	MetaPredicate
}

func (pred *FVGTEPredicate) Instruction() string {
	return fvgteIsn
}

func (pred *FVGTEPredicate) Token() string {
	return fvgteToken
}

func (pred *FVGTEPredicate) String() string {
	return pred.MetaPredicate.String(fvgteToken)
}

func (pred *FVGTEPredicate) Debug() string {
	return pred.MetaPredicate.Debug(fvgteIsn, fvgteToken)
}

func (pred *FVGTEPredicate) Eval(_ context.Context, input dag.Filterable) bool {
	lhs, lhsOk := pred.MetaPredicate.GetKeyAsFloat64(input)
	rhs, rhsOk := pred.MetaPredicate.GetValueAsFloat64()

	return (lhsOk && rhsOk) && lhs >= rhs
}

// ** Builder:

type FVGTEBuilder struct{}

func (bld *FVGTEBuilder) Token() string {
	return fvgteToken
}

func (bld *FVGTEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	pred := &FVGTEPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_fvgte.go ends here.
