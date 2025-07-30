// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_ssm.go --- SSM - String (Sensitive) Member.
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

// * Constants:

const (
	ssmIsn   = "SSM"
	ssmToken = "<"
)

// * Code:

// ** Predicate:

type SSMPredicate struct {
	MetaPredicate
}

func (pred *SSMPredicate) String() string {
	val, ok := pred.MetaPredicate.GetPredicateStringArray()
	if !ok {
		return FormatIsnf(ssmIsn, invalidTokenString)
	}

	return FormatIsnf(ssmIsn,
		"%s %s %#v",
		pred.MetaPredicate.key,
		ssmToken,
		val)
}

func (pred *SSMPredicate) Eval(input DataMap) bool {
	return pred.MetaPredicate.EvalStringMember(input, false)
}

// ** Builder:

type SSMBuilder struct{}

func (bld *SSMBuilder) Token() string {
	return ssmToken
}

func (bld *SSMBuilder) Build(key string, val any) Predicate {
	return &SSMPredicate{
		MetaPredicate: MetaPredicate{key: key, val: val},
	}
}

// * predicate_ssm.go ends here.
