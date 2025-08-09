// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_reim.go --- REIM - Regular Expression (Insensitive) Match.
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
	"strings"

	"github.com/Asmodai/gohacks/logger"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	reimIsn   = "REIM"
	reimToken = "regex-ci-match" //nolint:gosec
)

// * Code:

// ** Predicate:

// REIM - Regular Expression (Insensitive) Match predicate.
//
// Returns true if the regular expression in the filter matches against the
// input value.
//
// The regular expression will be compiled with a prefix denoting that it
// does not care about case.
type REIMPredicate struct {
	compiled *regexp.Regexp

	MetaPredicate
}

func (pred *REIMPredicate) Instruction() string {
	return reimIsn
}

func (pred *REIMPredicate) Token() string {
	return reimToken
}

func (pred *REIMPredicate) String() string {
	return pred.MetaPredicate.String(reimToken)
}

func (pred *REIMPredicate) Debug() string {
	return pred.MetaPredicate.Debug(reimIsn, reimToken)
}

func (pred *REIMPredicate) Eval(_ context.Context, input Filterable) bool {
	data, _, ok := pred.MetaPredicate.GetStringValues(input)
	if !ok {
		return false
	}

	return pred.compiled.MatchString(data)
}

// ** Builder:

type REIMBuilder struct{}

func (bld *REIMBuilder) Token() string {
	return reimToken
}

func (bld *REIMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error) {
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

	if !strings.HasPrefix(strVal, "(?i)") {
		strVal = "(?i)" + strVal
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

	pred := &REIMPredicate{
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

// * predicate_reim.go ends here.
