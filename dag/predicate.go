// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate.go --- Predicates.
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
	"strings"

	"github.com/Asmodai/gohacks/math/conversion"
	"github.com/Asmodai/gohacks/memoise"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	//nolint:gochecknoglobals
	predicateBuilders = map[string]PredicateFn{
		// Numeric predicates:
		"==": buildEQ,
		"!=": buildNEQ,
		"<":  buildLT,
		">":  buildGT,
		"<=": buildLTE,
		">=": buildGTE,

		// String predicates:
		"string-equal":        buildStringEQ,
		"string-not-equal":    buildStringNEQ,
		"string-ci-equal":     buildStringCIEQ,
		"string-ci-not-equal": buildStringCINEQ,

		// Regex predicates:
		"regex-match": buildRegexMatch,
	}

	//nolint:gochecknoglobals
	regexMemoiser memoise.Memoise
)

// * Code:

// ** Types:

// Attribute key/value structure.
type attrPair struct {
	Key string // Attribute key,
	Val any    // Value.
}

type Predicate struct {
	Eval func(*Predicate, DataMap) bool // Function to evaluate.
	Data any                            // Data to test against.
}

// ** Functions:

// *** Predicate builders:

// **** Numeric:

func buildEQ(attr string, val any) Predicate {
	return Predicate{Eval: evalEQ, Data: attrPair{Key: attr, Val: val}}
}

func buildNEQ(attr string, val any) Predicate {
	return Predicate{Eval: evalNEQ, Data: attrPair{Key: attr, Val: val}}
}

func buildLT(attr string, val any) Predicate {
	return Predicate{Eval: evalLT, Data: attrPair{Key: attr, Val: val}}
}

func buildGT(attr string, val any) Predicate {
	return Predicate{Eval: evalGT, Data: attrPair{Key: attr, Val: val}}
}

func buildLTE(attr string, val any) Predicate {
	return Predicate{Eval: evalLTE, Data: attrPair{Key: attr, Val: val}}
}

func buildGTE(attr string, val any) Predicate {
	return Predicate{Eval: evalGTE, Data: attrPair{Key: attr, Val: val}}
}

// **** Strings:

func buildStringEQ(attr string, val any) Predicate {
	return Predicate{evalStringEQ, attrPair{attr, val}}
}

func buildStringNEQ(attr string, val any) Predicate {
	return Predicate{evalStringNEQ, attrPair{attr, val}}
}

func buildStringCIEQ(attr string, val any) Predicate {
	return Predicate{evalStringCIEQ, attrPair{attr, val}}
}

func buildStringCINEQ(attr string, val any) Predicate {
	return Predicate{evalStringCINEQ, attrPair{attr, val}}
}

// **** Regex:

func buildRegexMatch(attr string, val any) Predicate {
	return Predicate{evalRegexMatch, attrPair{attr, val}}
}

// *** Predicate evaluators:

// **** Numeric:

// Numeric equality.
func evalEQ(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	fval, fvalok := conversion.ToFloat64(val)
	fdval, fdvalok := conversion.ToFloat64(data.Val)

	return valok && fvalok && fdvalok && fval == fdval
}

// Numeric inequality.
func evalNEQ(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	fval, fvalok := conversion.ToFloat64(val)
	fdval, fdvalok := conversion.ToFloat64(data.Val)

	return valok && fvalok && fdvalok && fval != fdval
}

// Numeric lesser-than.
func evalLT(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	lhs, lok := conversion.ToFloat64(val)
	rhs, rok := conversion.ToFloat64(data.Val)

	return valok && lok && rok && lhs < rhs
}

// Numeric greater-than.
func evalGT(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	lhs, lok := conversion.ToFloat64(val)
	rhs, rok := conversion.ToFloat64(data.Val)

	return valok && lok && rok && lhs > rhs
}

// Numeric lesser-than-or-equal-to.
func evalLTE(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	lhs, lok := conversion.ToFloat64(val)
	rhs, rok := conversion.ToFloat64(data.Val)

	return valok && lok && rok && lhs <= rhs
}

// Numeric greater-than-or-equal-to.
func evalGTE(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	lhs, lok := conversion.ToFloat64(val)
	rhs, rok := conversion.ToFloat64(data.Val)

	return valok && lok && rok && lhs >= rhs
}

// **** String:

// String equality, case-sensitive.
func evalStringEQ(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	sval, svalok := val.(string)
	sdval, sdvalok := data.Val.(string)

	return valok && svalok && sdvalok && sval == sdval
}

// String inequality, case-sensitive.
func evalStringNEQ(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	sval, svalok := val.(string)
	sdval, sdvalok := data.Val.(string)

	return valok && svalok && sdvalok && sval != sdval
}

// String equality, case-insensitive.
func evalStringCIEQ(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	sval, svalok := val.(string)
	sdval, sdvalok := data.Val.(string)

	return valok && svalok && sdvalok && strings.EqualFold(sval, sdval)
}

// String inequality, case-insensitive.
func evalStringCINEQ(p *Predicate, input DataMap) bool {
	data, ok := p.Data.(attrPair)
	if !ok {
		return false
	}

	val, valok := input[data.Key]
	sval, svalok := val.(string)
	sdval, sdvalok := data.Val.(string)

	return valok && svalok && sdvalok && !strings.EqualFold(sval, sdval)
}

// **** Regex:

// Regular expression match.
func evalRegexMatch(p *Predicate, input DataMap) bool {
	data, okay := p.Data.(attrPair)
	if !okay {
		return false
	}

	val, valok := input[data.Key]
	sval, svalok := val.(string)
	sdval, sdvalok := data.Val.(string)

	if !(valok && svalok && sdvalok) {
		return false
	}

	if regexMemoiser == nil {
		regexMemoiser = memoise.NewMemoise()
	}

	// The only time the `err` return is used is when the memoiser
	// callback errors, so we can ignore it here.
	memo, _ := regexMemoiser.Check(
		sdval,
		func() (any, error) {
			pattern := sdval
			if !strings.HasPrefix(pattern, "(?i)") {
				pattern = "(?i)" + pattern
			}

			result, err := regexp.Compile(pattern)
			if err != nil {
				panic(errors.WithMessagef(
					err,
					"regex compile failed: %s",
					err.Error()))
			}

			return result, nil
		},
	)

	compiled, okay := memo.(*regexp.Regexp)
	if !okay {
		panic(errors.Basef(
			"Memoiser did not return a compiled regexp: %q",
			sdval))
	}

	return compiled.MatchString(sval)
}

// * predicate.go ends here.
