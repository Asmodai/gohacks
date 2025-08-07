// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_eir.go --- EIR - Exclusive In Range.
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

// ** Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	eirIsn   = "EIR"
	eirToken = "in-range-exclusive" //nolint:gosec
)

// * Code:

// ** Predicate:

// EIR - Exclusive In Range predicate.
//
// Returne true if the input value is in the filter range inclusive.
type EIRPredicate struct {
	MetaPredicate
}

func (pred *EIRPredicate) Instruction() string {
	return eirIsn
}

func (pred *EIRPredicate) Token() string {
	return eirToken
}

func (pred *EIRPredicate) String() string {
	return pred.MetaPredicate.String(eirToken)
}

func (pred *EIRPredicate) Debug() string {
	return pred.MetaPredicate.Debug(eirIsn, eirToken)
}

func (pred *EIRPredicate) Eval(_ context.Context, input Filterable) bool {
	return pred.MetaPredicate.EvalExclusiveRange(input)
}

// ** Builder:

type EIRBuilder struct{}

func (bld *EIRBuilder) Token() string {
	return eirToken
}

func (bld *EIRBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error) {
	pred := &EIRPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_eir.go ends here.
