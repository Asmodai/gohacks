// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_sieq.go --- SIEG - String (Insensitive) Equality.
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
	sieqIsn   = "SIEQ"
	sieqToken = "string-ci-equal"
)

// * Code:

// ** Predicate:

// SIEG - String (Insensitive) Equality predicate.
//
// Returns true if the filter value matches the input value.
//
// This predicate does not care about case.
type SIEQPredicate struct {
	MetaPredicate
}

func (pred *SIEQPredicate) String() string {
	if val, ok := pred.MetaPredicate.val.(string); ok {
		return FormatIsnf(sieqIsn,
			"%s %s %#v",
			pred.MetaPredicate.key,
			sieqToken,
			val)
	}

	return FormatIsnf(sieqIsn, invalidTokenString)
}

func (pred *SIEQPredicate) Eval(_ context.Context, input Filterable) bool {
	lhs, rhs, ok := pred.MetaPredicate.GetStringValues(input)

	return ok && strings.EqualFold(lhs, rhs)
}

// ** Builder:

type SIEQBuilder struct{}

func (bld *SIEQBuilder) Token() string {
	return sieqToken
}

func (bld *SIEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error) {
	pred := &SIEQPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_sieq.go ends here.
