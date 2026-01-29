// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_sim.go --- SIM - String (Insensitive) Member.
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
package dag

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	simIsn   = "SIM"
	simToken = "string-ci-member"
)

// * Code:

// ** Predicate:

// SIM - String (Insensitive) Member predicate.
//
// Returns true if the input value is a member of the string array in the
// filter value.
type SIMPredicate struct {
	MetaPredicate
}

func (pred *SIMPredicate) Instruction() string {
	return simIsn
}

func (pred *SIMPredicate) Token() string {
	return simToken
}

func (pred *SIMPredicate) String() string {
	return pred.MetaPredicate.String(simToken)
}

func (pred *SIMPredicate) Debug() string {
	return pred.MetaPredicate.Debug(simIsn, simToken)
}

func (pred *SIMPredicate) Eval(_ context.Context, input Filterable) bool {
	return pred.MetaPredicate.EvalStringMember(input, true)
}

// ** Builder:

type SIMBuilder struct{}

func (bld *SIMBuilder) Token() string {
	return simToken
}

func (bld *SIMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error) {
	pred := &SIMPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_sim.go ends here.
