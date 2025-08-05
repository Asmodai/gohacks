// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: NONE
//
// predicate_sineq.go --- SINEQ - String (Insensitive) Inequality.
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

import (
	"context"
	"strings"

	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	sineqIsn   = "SINEQ"
	sineqToken = "string-ci-not-equal" //nolint:gosec
)

// * Code:

// ** Predicate:

// SINEQ - String (Insensitive) Inequality predicate.
//
// Returns true if the input string is not the same as the filter string.
//
// Case is not taken into account.
type SINEQPredicate struct {
	MetaPredicate
}

func (pred *SINEQPredicate) String() string {
	if val, ok := pred.MetaPredicate.val.(string); ok {
		return FormatIsnf(sineqIsn,
			"%s %s %#v",
			pred.MetaPredicate.key,
			sineqToken,
			val)
	}

	return FormatIsnf(sineqIsn, invalidTokenString)
}

func (pred *SINEQPredicate) Eval(_ context.Context, input Filterable) bool {
	lhs, rhs, ok := pred.MetaPredicate.GetStringValues(input)

	return ok && !strings.EqualFold(lhs, rhs)
}

// ** Builder:

type SINEQBuilder struct{}

func (bld *SINEQBuilder) Token() string {
	return sineqToken
}

func (bld *SINEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error) {
	pred := &SINEQPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_sineq.go ends here.
