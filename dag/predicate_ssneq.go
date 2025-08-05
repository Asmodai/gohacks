// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_ssneq.go --- SSNEQ - String (Sensitive) Inequality.
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

	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	ssneqIsn   = "SSNEQ"
	ssneqToken = "string-not-equal"
)

// * Code:

// ** Predicate:

// SSNEQ - String (Sensitive) Inequality predicate.
//
// Returns true if the input value is different to the filter value.
type SSNEQPredicate struct {
	MetaPredicate
}

func (pred *SSNEQPredicate) String() string {
	if val, ok := pred.MetaPredicate.val.(string); ok {
		return FormatIsnf(ssneqIsn,
			"%s %s %#v",
			pred.MetaPredicate.key,
			ssneqToken,
			val)
	}

	return FormatIsnf(ssneqIsn, invalidTokenString)
}

func (pred *SSNEQPredicate) Eval(_ context.Context, input Filterable) bool {
	lhs, rhs, ok := pred.MetaPredicate.GetStringValues(input)

	return ok && lhs != rhs
}

// ** Builder:

type SSNEQBuilder struct{}

func (bld *SSNEQBuilder) Token() string {
	return ssneqToken
}

func (bld *SSNEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error) {
	pred := &SSNEQPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// * predicate_ssneq.go ends here.
