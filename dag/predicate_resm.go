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
	"regexp"
	"sync"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	resmIsn   = "RESM"
	resmToken = "regex-match" //nolint:gosec
)

// * Code:

// ** Predicate:

type RESMPredicate struct {
	MetaPredicate

	compiled *regexp.Regexp
	once     sync.Once
	err      error
}

func (pred *RESMPredicate) String() string {
	val, ok := pred.MetaPredicate.val.(string)
	if !ok {
		return FormatIsnf(resmIsn, invalidTokenString)
	}

	return FormatIsnf(resmIsn, "%s %s %#v", pred.MetaPredicate.key, resmToken, val)
}

func (pred *RESMPredicate) compilePattern(pattern string) {
	pred.once.Do(func() {
		pred.compiled, pred.err = regexp.Compile(pattern)
	})
}

func (pred *RESMPredicate) Eval(input Filterable) bool {
	data, pattern, ok := pred.MetaPredicate.GetStringValues(input)
	if !ok {
		return false
	}

	pred.compilePattern(pattern)

	if pred.err != nil {
		panic(errors.WithMessagef(
			pred.err,
			"regex compilation failed: %s",
			pred.err.Error()))
	}

	return pred.compiled.MatchString(data)
}

// ** Builder:

type RESMBuilder struct{}

func (bld *RESMBuilder) Token() string {
	return resmToken
}

func (bld *RESMBuilder) Build(key string, val any) Predicate {
	return &RESMPredicate{
		MetaPredicate: MetaPredicate{key: key, val: val},
	}
}

// * predicate_resm.go ends here.
