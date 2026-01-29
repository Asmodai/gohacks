// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// actions.go --- Validator actions.
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
	"fmt"
	"strings"

	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/logger"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	ErrFailInvoked = errors.Base("fail action was invoked")
)

// * Code:

// ** Types:

// Implementation of validator actions.
//
// Conforms to `dag.Actions`.
type actions struct {
	errors []error
}

// ** Methods:

// Build a compiled action from the given function name and parameters.
func (act *actions) Builder(funame string, params dag.ActionParams) (dag.ActionFn, error) {
	normalised := strings.ToLower(funame) // Normalise to lower case.

	switch normalised {
	case "none":
		return act.noneAction(params)

	case "log":
		return act.logAction(params)

	case "error":
		return act.errorAction(params)

	default:
		return nil, errors.WithMessagef(
			dag.ErrUnknownBuiltin,
			"%q",
			funame)
	}
}

// Compile an `none` action.
//
// There are no required parameters.
func (act *actions) noneAction(_ dag.ActionParams) (dag.ActionFn, error) {
	afn := func(_ context.Context, _ dag.Filterable) {}

	return afn, nil
}

// Compile an `error` action.
//
// Required parameters are:
//
//	`message`: The message to print in the log.
//
// If no valid parameters are passed then `ErrExpectedParams` is returned.
//
// If no `message` parameter is provided then `ErrMissingParam` is returned.
//
// If the given `message` parameter is not a string then `ErrExpectedString`
// is returned.
//
// Upon success a function object containing the error generator is returned.
func (act *actions) errorAction(params dag.ActionParams) (dag.ActionFn, error) {
	if params == nil {
		return nil, errors.WithStack(dag.ErrExpectedParams)
	}

	msg, okay := params["message"]
	if !okay {
		return nil, errors.WithMessage(dag.ErrMissingParam,
			`Parameter "message"`)
	}

	smsg, okay := msg.(string)
	if !okay {
		return nil, errors.WithStack(dag.ErrExpectedString)
	}

	afn := func(_ context.Context, _ dag.Filterable) {
		if act.errors == nil {
			act.errors = []error{}
		}

		//nolint:err113
		act.errors = append(act.errors, fmt.Errorf("%s", smsg))
	}

	return afn, nil
}

// Compile a `log` action.
//
// Required parameters are:
//
//	`message`: The message to print in the log.
//
// If no valid parameters are passed then `ErrExpectedParams` is returned.
//
// If no `message` parameter is provided then `ErrMissingParam` is returned.
//
// If the given `message` parameter is not a string then `ErrExpectedString`
// is returned.
//
// Upon success a function object containing the compiled logger is returned.
func (act *actions) logAction(params dag.ActionParams) (dag.ActionFn, error) {
	if params == nil {
		return nil, errors.WithStack(dag.ErrExpectedParams)
	}

	msg, okay := params["message"]
	if !okay {
		return nil, errors.WithMessage(dag.ErrMissingParam,
			`Parameter "message"`)
	}

	smsg, okay := msg.(string)
	if !okay {
		return nil, errors.WithStack(dag.ErrExpectedString)
	}

	afn := func(ctx context.Context, input dag.Filterable) {
		lgr := logger.MustGetLogger(ctx)

		lgr.Info(
			smsg,
			"src", "log_action",
			"structure", input)
	}

	return afn, nil
}

// * actions.go ends here.
