// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// compiler_test.go --- Compiler tests.
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

package dag

// * Imports:

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/Asmodai/gohacks/contextdi"
	"github.com/Asmodai/gohacks/logger"
	mlogger "github.com/Asmodai/gohacks/mocks/logger"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/mock/gomock"
)

// * Constants:

const (
	rulesYAML string = `
- name: warm_weather
  conditions:
    - attribute: type
      operator: string-equal
      value: weather
    - attribute: temp
      operator: '>='
      value: 22
  action:
    name: warm_weather_action
    perform: log
    params:
      message: It's warm!
- name: cold_weather
  conditions:
    - attribute: type
      operator: string-equal
      value: weather
    - attribute: temp
      operator: '<='
      value: 10
  action:
    name: cold_weather_action
    perform: log
    params:
      message: It's cold!`

	badRulesYAML string = `
- name: derpy_weather
  conditions:
    - attribute: type
      operator: dance-please
      value: lots
  action:
    name: woo_woo
    perform: dance
    params:
      type: fandango
- name: chungus_weather
  conditions:
    - attribute: type
      operator: string-equal
      value: lots
  action:
    name: woo_woo
    perform: dance
    params:
      type: fandango
- name: blibble_weather
  conditions:
    - attribute: type
      operator: string-equal
      value: lots
  action:
    name: slibblewibbles
    perform: log
    params:
      type: log`

	rulesJSON string = `[
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

// ** Test structure:

type BenchTest struct {
	one int
	two string
}

func (bt *BenchTest) Get(key string) (any, bool) {
	switch key {
	case "One":
		return bt.one, true
	case "Two":
		return bt.two, true
	default:
		return nil, false
	}
}

func (bt *BenchTest) Set(key string, val any) bool {
	switch key {
	case "One":
		bt.one = val.(int)
		return true

	case "Two":
		bt.two = val.(string)
		return true

	default:
		return false
	}
}

func (bt *BenchTest) Keys() []string {
	return []string{"One", "Two"}
}

func (bt *BenchTest) String() string {
	return fmt.Sprintf("One:%v  Two:%v", bt.one, bt.two)
}

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
		return nil, errors.WithMessagef(ErrUnknownBuiltin, "%q", fn)
	}
}

func (obj *MockActions) logAction(params ActionParams) (afn ActionFn, err error) {
	msgparam, ok := params["message"]
	if !ok {
		afn = nil
		err = errors.WithMessage(ErrMissingParam,
			`Parameter "message"`)

		return
	}

	msg, ok := msgparam.(string)
	if !ok {
		afn = nil
		err = errors.WithStack(ErrExpectedString)

		return
	}

	err = nil
	afn = func(ctx context.Context, input Filterable) {
		log.Printf("LOG ACTION: Run with: %#v\n", input)

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

	stdlgr := logger.NewDefaultLogger()
	lgr := mlogger.NewMockLogger(mocker)

	lgr.EXPECT().
		Error(gomock.Any(), gomock.Any()).
		DoAndReturn(func(message string, rest ...any) {
			stdlgr.Error(message, rest...)
		}).
		MaxTimes(0)

	lgr.EXPECT().
		Fatal(gomock.Any(), gomock.Any()).
		MaxTimes(0)

	lgr.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		DoAndReturn(func(message string, rest ...any) {
			stdlgr.Info(message, rest...)
		}).
		AnyTimes()

	lgr.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		DoAndReturn(func(message string, rest ...any) {
			stdlgr.Debug(message, rest...)
		}).
		AnyTimes()

	mact := &MockActions{
		hasRunLog: false,
	}

	ctx, err := logger.SetLogger(context.TODO(), lgr)
	if err != nil {
		t.Fatalf("Could not set logger DI: %#v", err)
	}

	ctx, err = contextdi.SetDebugMode(ctx, true)
	if err != nil {
		t.Fatalf("Could not set debug flag to DI: %#v", err)
	}

	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		t.Errorf("JSON: %v", err.Error())
		return
	}

	compiler := NewCompiler(ctx, mact)
	issues := compiler.Compile(rules)

	compiler.Export(os.Stdout)

	if len(issues) > 0 {
		t.Logf("Compiler issues:")
		for idx := range issues {
			t.Logf("%d: %s", idx+1, issues[idx])
		}
	}

	t.Run("Warm alert", func(t *testing.T) {
		mact.hasRunLog = false

		input := NewDataInputFromMap(WarmWeather)
		compiler.Evaluate(input)

		if !mact.hasRunLog {
			t.Error("Warm alert was not triggered.")
		}
	})

	t.Run("Cold alert", func(t *testing.T) {
		mact.hasRunLog = false

		input := NewDataInputFromMap(ColdWeather)
		compiler.Evaluate(input)

		if !mact.hasRunLog {
			t.Error("Cold alert was not triggered.")
		}
	})

	t.Run("No alert", func(t *testing.T) {
		mact.hasRunLog = false

		input := NewDataInputFromMap(NormalWeather)
		compiler.Evaluate(input)

		if mact.hasRunLog {
			t.Error("Unexpected alert(s) triggered!")
		}
	})

	t.Run("Test compiler warnings", func(t *testing.T) {
		jank, err := ParseFromYAML(badRulesYAML)
		if err != nil {
			t.Fatalf("YAML: %#v", err)
		}

		compiler := NewCompiler(ctx, mact)
		issues := compiler.Compile(jank)

		if len(issues) == 0 {
			t.Fatalf("How did this jank compile?")
		}

		if len(issues) > 3 {
			t.Fatalf("Too many issues: %d != 3", len(issues))
		}

		t.Logf("Compiler issues:")
		for idx := range issues {
			t.Logf("%d: %s", idx+1, issues[idx])
		}
	})
}

func TestUtilities(t *testing.T) {
	t.Run("ParseFromYAML", func(t *testing.T) {
		_, err := ParseFromYAML(rulesYAML)
		if err != nil {
			t.Fatalf("YAML: %#v", err)
		}
	})

	t.Run("ParseFromJSON", func(t *testing.T) {
		_, err := ParseFromJSON(rulesJSON)
		if err != nil {
			t.Fatalf("YAML: %#v", err)
		}
	})
}

// ** Benchmarks:

func BenchmarkCompiler(b *testing.B) {
	data := &BenchTest{one: 42, two: "Foo"}

	mocker := gomock.NewController(b)
	defer mocker.Finish()

	lgr := mlogger.NewMockLogger(mocker)
	ctx, err := logger.SetLogger(context.TODO(), lgr)
	if err != nil {
		b.Fatalf("Could not set logger DI: %#v", err)
	}

	rules, err := ParseFromYAML(rulesYAML)
	if err != nil {
		b.Fatalf("YAML: %#v", err)
	}

	mact := &MockActions{
		hasRunLog: false,
	}

	compiler := NewCompiler(ctx, mact)

	b.Run("Compile()", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			compiler.Compile(rules)
		}
	})

	b.Run("Evaluate()", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			compiler.Evaluate(data)
		}
	})
}

// * compiler_test.go ends here.
