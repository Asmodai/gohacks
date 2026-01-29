// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fteq.go --- FTEQ - Field Type Equals.
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

package validator

// * Imports:

import (
	"context"
	"reflect"
	"strings"

	"github.com/Asmodai/gohacks/conversion"
	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/logger"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	fteqIsn   = "FTEQ"
	fteqToken = "field-type-equal"
)

// * Variables:

var (
	ErrValueNotString = errors.Base("value is not a string")
)

// * Code:

// ** Predicate:

// Field Type Equality.
//
// This predicate compares the type of the structure's field.  If it is
// equal then the predicate returns true.
type FTEQPredicate struct {
	MetaPredicate
}

func (pred *FTEQPredicate) Instruction() string {
	return fteqIsn
}

func (pred *FTEQPredicate) Token() string {
	return fteqToken
}

func (pred *FTEQPredicate) String() string {
	return pred.MetaPredicate.String(fteqToken)
}

func (pred *FTEQPredicate) Debug() string {
	return pred.MetaPredicate.Debug(fteqIsn, fteqToken)
}

func (pred *FTEQPredicate) Eval(_ context.Context, input dag.Filterable) bool {
	want, wantOk := pred.MetaPredicate.GetValueAsString()
	fInfo, fInfoOk := pred.MetaPredicate.GetKeyAsFieldInfo(input)

	if !(wantOk && fInfoOk) {
		return false
	}

	if fInfo.TypeKind == reflect.Interface {
		if want == "any" || want == "interface {}" {
			// We want an `any`, so don't bother with the rest.
			return true
		}

		val, valok := pred.MetaPredicate.GetKeyAsValue(input)
		if valok {
			tname, valid := resolveAnyType(val)
			if !valid {
				return false
			}

			return strings.EqualFold(want, tname)
		}
	}

	return strings.EqualFold(want, fInfo.TypeName)
}

// ** Builder:

type FTEQBuilder struct{}

func (bld *FTEQBuilder) Token() string {
	return fteqToken
}

func (bld *FTEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (dag.Predicate, error) {
	sval, svalOk := conversion.ToString(val)
	if !svalOk {
		return nil, errors.WithMessagef(
			ErrValueNotString,
			"%s: value %q",
			fteqToken,
			val)
	}

	pred := &FTEQPredicate{
		MetaPredicate: MetaPredicate{
			key:    key,
			val:    normaliseTypeName(sval),
			logger: lgr,
			debug:  dbg,
		},
	}

	return pred, nil
}

// ** Functions:

func normaliseTypeName(name string) string {
	const (
		anyStr    = "any"
		iface1Str = "interface{}"
		iface2Str = "interface {}"
	)

	switch strings.ToLower(strings.TrimSpace(name)) {
	case anyStr, iface1Str, iface2Str:
		return iface2Str
	}

	return name
}

//nolint:cyclop,funlen
func resolveAnyType(value any) (string, bool) {
	switch value.(type) {
	case int:
		return "int", true
	case int8:
		return "int8", true
	case int16:
		return "int16", true
	case int32:
		return "int32", true
	case int64:
		return "int64", true
	case uint:
		return "uint", true
	case uint8:
		return "uint8", true
	case uint16:
		return "uint16", true
	case uint32:
		return "uint32", true
	case uint64:
		return "uint64", true
	case float32:
		return "float32", true
	case float64:
		return "float64", true
	case complex64:
		return "complex64", true
	case complex128:
		return "complex128", true
	case bool:
		return "bool", true

	case string:
		return "string", true

	case []byte:
		return "[]byte", true

	case []any:
		return "[]any", true

	case any:
		return "any", true

	default:
		return "", false
	}
}

// * predicate_fteq.go ends here.
