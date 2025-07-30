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
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Asmodai/gohacks/logger"
	mlogger "github.com/Asmodai/gohacks/mocks/logger"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/mock/gomock"
)

// * Constants:

const (
	RulesJSON string = `[
		{
			"name": "warm_weather",
			"conditions": [
				{
					"attribute": "type",
					"operator":  "string-equal",
					"value":     "weather"
				}, {
					"attribute": "temp",
					"operator":  ">=",
					"value":     22
				}
			],
			"action": {
				"name":    "warm_weather_action",
				"perform": "log",
				"params": {
					"message": "It's warm!"
				}
			}
		}, {
			"name": "cold_weather",
			"conditions": [
				{
					"attribute": "type",
					"operator":  "string-equal",
					"value":     "weather"
				}, {
					"attribute": "temp",
					"operator":  "<=",
					"value":     10
				}
			],
			"action": {
				"name":    "cold_weather_action",
				"perform": "log",
				"params": {
					"message": "It's cold!"
				}
			}
		}
	]`
)

// * Variables:

var (
	ColdWeather map[string]any = map[string]any{
		"type": "weather",
		"temp": 8,
	}

	WarmWeather map[string]any = map[string]any{
		"type": "weather",
		"temp": 27,
	}

	NormalWeather map[string]any = map[string]any{
		"type": "weather",
		"temp": 19,
	}
)

// * Code:

// ** Mock actions:

// *** Type:

type MockActions struct {
	hasRunLog bool
}

// *** Methods:

func (obj *MockActions) Builder(fn string, params ActionParams) (ActionFn, error) {
	normalised := strings.ToLower(fn)

	switch normalised {
	case "log":
		return obj.logAction(params)

	default:
		return nil, errors.WithMessagef(
			ErrUnknownBuiltin,
			"unknown builtin function '%s'",
			fn,
		)
	}
}

func (obj *MockActions) logAction(params ActionParams) (afn ActionFn, err error) {
	msgparam, ok := params["message"]
	if !ok {
		afn = nil
		err = errors.WithStack(ErrMissingParam)

		return
	}

	msg, ok := msgparam.(string)
	if !ok {
		afn = nil
		err = errors.WithStack(ErrExpectedString)

		return
	}

	err = nil
	afn = func(ctx context.Context, input DataMap) {
		lgr := logger.MustGetLogger(ctx)
		lgr.Info(
			msg,
			"src", "log_action",
			"structure", input,
		)

		obj.hasRunLog = true
	}

	return
}

// ** Tests:

func TestCompiler(t *testing.T) {
	var (
		rules []RuleSpec
	)

	mocker := gomock.NewController(t)
	defer mocker.Finish()

	lgr := mlogger.NewMockLogger(mocker)

	lgr.EXPECT().
		Error(gomock.Any(), gomock.Any()).
		MaxTimes(0)

	lgr.EXPECT().
		Fatal(gomock.Any(), gomock.Any()).
		MaxTimes(0)

	lgr.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	mact := &MockActions{
		hasRunLog: false,
	}

	ctx, err := logger.SetLogger(context.TODO(), lgr)
	if err != nil {
		t.Fatalf("Could not set logger DI: %#v", err)
	}

	if err := json.Unmarshal([]byte(RulesJSON), &rules); err != nil {
		t.Errorf("JSON: %v", err.Error())
		return
	}

	compiler := NewCompiler(ctx, mact)
	issues := compiler.Compile(rules)

	if len(issues) > 0 {
		t.Logf("Compiler issues:")
		for idx := range issues {
			t.Logf("%d: %s", idx+1, issues[idx])
		}
	}

	t.Run("Warm alert", func(t *testing.T) {
		mact.hasRunLog = false

		compiler.Evaluate(WarmWeather)

		if !mact.hasRunLog {
			t.Error("Warm alert was not triggered.")
		}
	})

	t.Run("Cold alert", func(t *testing.T) {
		mact.hasRunLog = false

		compiler.Evaluate(ColdWeather)

		if !mact.hasRunLog {
			t.Error("Cold alert was not triggered.")
		}
	})

	t.Run("No alert", func(t *testing.T) {
		mact.hasRunLog = false

		compiler.Evaluate(NormalWeather)

		if mact.hasRunLog {
			t.Error("Unexpected alert(s) triggered.")
		}
	})
}

// * compiler_test.go ends here.
