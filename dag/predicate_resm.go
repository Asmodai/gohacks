// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_resm.go --- RESM - Regular Expression (Sensitive) Match.
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
	"regexp"

	"github.com/Asmodai/gohacks/logger"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	resmIsn   = "RESM"
	resmToken = "regex-match" //nolint:gosec
)

// * Variables:

var (
	ErrInvalidRegexp  = errors.Base("invalid regexp")
	ErrRegexpParse    = errors.Base("error parsing regexp")
	ErrValueNotString = errors.Base("value is not a string")
)

// * Code:

// ** Predicate:

// RESM - Regular Expression (Sensitive) Match predicate.
//
// Returns true if the regular expression in the filter value matches against
// the input value.
//
// The regular expression will not be forced into being case-insensitive.
type RESMPredicate struct {
	compiled *regexp.Regexp

	MetaPredicate
}

func (pred *RESMPredicate) Instruction() string {
	return resmIsn
}

func (pred *RESMPredicate) Token() string {
	return resmToken
}

func (pred *RESMPredicate) String() string {
	return pred.MetaPredicate.String(resmToken)
}

func (pred *RESMPredicate) Debug() string {
	return pred.MetaPredicate.Debug(resmIsn, resmToken)
}

func (pred *RESMPredicate) Eval(_ context.Context, input Filterable) bool {
	data, _, ok := pred.MetaPredicate.GetStringValues(input) // TODO: simplify.
	if !ok {
		return false
	}

	return pred.compiled.MatchString(data)
}

// ** Builder:

type RESMBuilder struct{}

func (bld *RESMBuilder) Token() string {
	return resmToken
}

func (bld *RESMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error) {
	strVal, strOk := val.(string)
	if !strOk {
		return nil, errors.WithMessagef(
			ErrValueNotString,
			"%s: value %q",
			resmToken,
			val)
	}

	if len(strVal) == 0 {
		return nil, errors.WithMessagef(
			ErrInvalidRegexp,
			"%s: %q",
			resmToken,
			strVal)
	}

	compiled, err := regexp.Compile(strVal)
	if err != nil {
		return nil, errors.WithMessagef(
			errors.WrapWith(err, ErrRegexpParse),
			"%s: regex %q: %s",
			resmToken,
			strVal,
			err.Error())
	}

	pred := &RESMPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    val,
			logger: lgr,
			debug:  dbg,
		},
		compiled: compiled,
	}

	return pred, nil
}

// * predicate_resm.go ends here.
