// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvrem.go --- FVREM - Field Value Regular Expression Match.
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
	"regexp"

	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/logger"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	fvremIsn   = "FVREM"
	fvremToken = "field-value-regex-match" //nolint:gosec
)

// * Variables:

var (
	ErrInvalidRegexp = errors.Base("invalid regexp")
	ErrRegexpParse   = errors.Base("error parsing regexp")
)

// * Code:

// ** Predicate:

type FVREMPredicate struct {
	MetaPredicate

	compiled *regexp.Regexp
}

func (pred *FVREMPredicate) Instruction() string {
	return fvremIsn
}

func (pred *FVREMPredicate) Token() string {
	return fvremToken
}

func (pred *FVREMPredicate) String() string {
	return pred.MetaPredicate.String(fvremToken)
}

func (pred *FVREMPredicate) Debug() string {
	return pred.MetaPredicate.Debug(fvremIsn, fvremToken)
}

func (pred *FVREMPredicate) Eval(_ context.Context, input dag.Filterable) bool {
	data, dataOk := pred.MetaPredicate.GetKeyAsString(input)
	if !dataOk {
		return false
	}

	return pred.compiled.MatchString(data)
}

// ** Builder:

type FVREMBuilder struct{}

func (bld *FVREMBuilder) Token() string {
	return fvremToken
}

func (bld *FVREMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	strVal, strOk := val.(string)
	if !strOk {
		return nil, errors.WithMessagef(
			ErrValueNotString,
			"%s: value %q",
			fvremToken,
			val)
	}

	if len(strVal) == 0 {
		return nil, errors.WithMessagef(
			ErrInvalidRegexp,
			"%s: %q",
			fvremToken,
			strVal)
	}

	compiled, err := regexp.Compile(strVal)
	if err != nil {
		return nil, errors.WithMessagef(
			errors.WrapWith(err, ErrRegexpParse),
			"%s: regex %q: %s",
			fvremToken,
			strVal,
			err.Error())
	}

	pred := &FVREMPredicate{
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

// * predicate_fvrem.go ends here.
