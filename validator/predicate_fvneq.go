// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvneq.go --- FVNEQ - Field Value Inequality.
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

	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	fvneqIsn   = "FVNEQ"
	fvneqToken = "field-value-not-equal"
)

// * Code:

// ** Predicate:

// Field Valie Inequality.
//
// This predicate compares the value to that in the structure. If they are
// not equal then the predicate returns true.
//
// The predicate will take various circumstances into consideration while
// checking the value:
//
// If the field is `any` then the comparison will match just the type of the
// value rather than using the type of the field along with the value.
//
// If the field is integer then the structure's field must have a bit
// width large enough to hold the value.
type FVNEQPredicate struct {
	FVEQPredicate
}

func (pred *FVNEQPredicate) Instruction() string {
	return fvneqIsn
}

func (pred *FVNEQPredicate) Token() string {
	return fvneqToken
}

func (pred *FVNEQPredicate) String() string {
	return pred.MetaPredicate.String(fvneqToken)
}

func (pred *FVNEQPredicate) Debug() string {
	return pred.MetaPredicate.Debug(fvneqIsn, fvneqToken)
}

func (pred *FVNEQPredicate) Eval(ctx context.Context, input dag.Filterable) bool {
	return !pred.FVEQPredicate.Eval(ctx, input)
}

// ** Builder:

type FVNEQBuilder struct{}

func (bld *FVNEQBuilder) Token() string {
	return fvneqToken
}

func (bld *FVNEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	pred := &FVNEQPredicate{
		FVEQPredicate: FVEQPredicate{
			MetaPredicate: MetaPredicate{
				key:    key,
				val:    val,
				logger: lgr,
				debug:  dbg,
			},
		},
	}

	return pred, nil
}

// * predicate_fvneq.go ends here.
