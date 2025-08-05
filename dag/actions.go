// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// actions.go --- DAG actions.
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
//
//mock:yes

// * Comments:

//
// It might be nice if we allowed multiple actions to be attached to a rule.
//

// * Package:

package dag

// * Imports:

import (
	"context"
	"strings"

	"github.com/Asmodai/gohacks/contextdi"
	"github.com/Asmodai/gohacks/logger"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	ErrExpectedParams = errors.Base("expected parameters to be given")
	ErrExpectedString = errors.Base("expected a string value")
	ErrMissingParam   = errors.Base("parameter missing")
	ErrUnknownBuiltin = errors.Base("unknown builtin function")
)

// * Code:

// ** Interface:

// Action builder interface.
//
// The action builder provides a means of compiling JSON or YAML actions
// into explicit function objects
//
// The resulting action is a function that takes `context.Context` and
// `Filterable` arguments and then performs some sort of user-defined action.
//
// There are two default builtins provided for you:
//
//	`log`:    Log the contents of the parameters to a logger.
//	`mutate`: Change value(s) in the parameters.
//
// To use the `log` builtin, you must provide a `logger.Logger` instance in
// the context used with the DAG.  For this, you can see `logger.SetLogger`.
type Actions interface {
	// Build the given builtin functions.
	Builder(string, ActionParams) (ActionFn, error)
}

// ** Types:

// Implementation of the default builtin action builder.
type actions struct {
}

// ** Methods:

// Build a compiled action from the given function name and parameters.
func (act *actions) Builder(funame string, params ActionParams) (ActionFn, error) {
	normalised := strings.ToLower(funame) // Normalise to lower case.

	switch normalised {
	case "mutate":
		return act.mutateAction(params)

	case "log":
		return act.logAction(params)

	default:
		return nil, errors.WithMessagef(ErrUnknownBuiltin, "%q", funame)
	}
}

// Compile a 'mutate' action.
//
// Required parameters are:
//
//	`attribute`: The attribute within `Filterable` that is to be mutated.
//	`new_value`: The value to be stored in the given attribute.
//
// If no valid parameters are provided then `ErrExpectedParams` is returned.
//
// If no `attribute` parameter is provided then `ErrMissingParam` is returned.
//
// If the given `attribute` parameter is not a string then `ErrExpectedString`
// is returned.
//
// If no `new_value` parameter is provided then `ErrMissingParam` is returned.
//
// Upon success a function object containing the compiled mutator is returned.
func (act *actions) mutateAction(params ActionParams) (ActionFn, error) {
	if params == nil {
		return nil, errors.WithStack(ErrExpectedParams)
	}

	attr, okay := params["attribute"]
	if !okay {
		return nil, errors.WithMessage(ErrMissingParam,
			`Parameter "attribute"`)
	}

	sattr, okay := attr.(string)
	if !okay {
		return nil, errors.WithStack(ErrExpectedString)
	}

	val, okay := params["new_value"]
	if !okay {
		return nil, errors.WithMessage(ErrMissingParam,
			`Parameter "new_value"`)
	}

	afn := func(ctx context.Context, input Filterable) {
		lgr := logger.MustGetLogger(ctx)

		dbg, err := contextdi.GetDebugMode(ctx)
		if err != nil {
			dbg = false
		}

		if dbg {
			lgr.Debug(
				"MUTATE action invoked",
				"structure", input,
				"attribute", sattr,
				"new_value", val)
		}

		input.Set(sattr, val)
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
func (act *actions) logAction(params ActionParams) (ActionFn, error) {
	if params == nil {
		return nil, errors.WithStack(ErrExpectedParams)
	}

	msg, okay := params["message"]
	if !okay {
		return nil, errors.WithMessage(ErrMissingParam,
			`Parameter "message"`)
	}

	smsg, okay := msg.(string)
	if !okay {
		return nil, errors.WithStack(ErrExpectedString)
	}

	afn := func(ctx context.Context, input Filterable) {
		lgr := logger.MustGetLogger(ctx)

		lgr.Info(
			smsg,
			"src", "log_action",
			"structure", input)
	}

	return afn, nil
}

// ** Functions

// Create a new empty `Actions` object.
func NewDefaultActions() Actions {
	return &actions{}
}

// * actions.go ends here.
