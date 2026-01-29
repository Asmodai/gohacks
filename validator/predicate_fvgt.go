// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvgt.go --- FVGT - Field Value is Lesser Than.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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
	fvgtIsn   = "FVGT"
	fvgtToken = "field-value->"
)

// * Code:

// ** Predicate:

// This predicate returns true if the value in the structure is greater than
// the value in the predicate.
type FVGTPredicate struct {
	MetaPredicate
}

func (pred *FVGTPredicate) Instruction() string {
	return fvgtIsn
}

func (pred *FVGTPredicate) Token() string {
	return fvgtToken
}

func (pred *FVGTPredicate) String() string {
	return pred.MetaPredicate.String(fvgtToken)
}

func (pred *FVGTPredicate) Debug() string {
	return pred.MetaPredicate.Debug(fvgtIsn, fvgtToken)
}

func (pred *FVGTPredicate) Eval(_ context.Context, input dag.Filterable) bool {
	lhs, lhsOk := pred.MetaPredicate.GetKeyAsFloat64(input)
	rhs, rhsOk := pred.MetaPredicate.GetValueAsFloat64()

	return (lhsOk && rhsOk) && lhs > rhs
}

// ** Builder:

type FVGTBuilder struct{}

func (bld *FVGTBuilder) Token() string {
	return fvgtToken
}

func (bld *FVGTBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	pred := &FVGTPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_fvgt.go ends here.
