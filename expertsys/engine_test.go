// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// engine_test.go --- Engine tests.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

package expertsys

// * Imports:

import (
	"context"
	"os"
	"testing"

	"github.com/Asmodai/gohacks/contextdi"
	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/errx"
	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

// * Variables:

// * Code:

// ** Mock actions:

// *** Type:

type MockActions struct{}

// *** Methods:

func (ma *MockActions) Builder(fn string, params dag.ActionParams) (dag.ActionFn, error) {
	switch fn {
	case "assert":
		keyAny, ok := params["key"]
		if !ok {
			return nil, errx.New(`missing param "key"`)
		}

		key, ok := keyAny.(string)
		if !ok || len(key) == 0 {
			return nil, errx.New(`param "key" must be a non-empty string"`)
		}

		val, ok := params["value"]
		if !ok {
			return nil, errx.New(`missing param "value"`)
		}

		return func(ctx context.Context, input dag.Filterable) {
			_ = input.Set(key, val)
		}, nil

	default:
		return nil, errx.New("unknown action: " + fn)
	}
}

// ** Tests:

func TestEngine(t *testing.T) {
	ctx, err := logger.SetLogger(context.TODO(), logger.NewDefaultLogger())
	if err != nil {
		t.Fatalf("set logger DI: %#v", err)
	}

	ctx, err = contextdi.SetDebugMode(ctx, true)
	if err != nil {
		t.Fatalf("Could not set debug flag to DI: %#v", err)
	}

	cmplr := dag.NewCompiler(ctx, &MockActions{})

	rules := []dag.RuleSpec{
		{
			Name: "rule_b_depends_on_foo",
			Conditions: []dag.ConditionSpec{
				{
					Attribute: "foo",
					Operator:  "string-equal",
					Value:     "yes",
				},
			},
			Action: &dag.ActionSpec{
				Name:    "set_bar",
				Reason:  "foo is `yes`.",
				Perform: "assert",
				Params: dag.ActionParams{
					"key":   "bar",
					"value": "done",
				},
			},
		}, {
			Name: "rule_a_sets_foo",
			Conditions: []dag.ConditionSpec{
				{
					Attribute: "type",
					Operator:  "string-equal",
					Value:     "event",
				},
			},
			Action: &dag.ActionSpec{
				Name:    "set_foo",
				Reason:  "`foo` is set.",
				Perform: "assert",
				Params: dag.ActionParams{
					"key":   "foo",
					"value": "yes",
				},
			},
		},
	}

	issues := cmplr.Compile(rules)
	if len(issues) > 0 {
		t.Log("Compiler issues:")
		for idx := range issues {
			t.Logf("%02d: %s", idx+1, issues[idx])
		}
		t.Fatal("compiler issues present, see log output")
	}

	eng := &Engine{cmplr: cmplr}

	t.Run("maxIters=1 should not stabilise", func(t *testing.T) {
		wm := NewWorkingMemory()
		wm.Set("type", "event")

		_, err := eng.RunToFixpoint(wm, 1)

		if err == nil {
			t.Fatal("expected non-stable error")
		}
	})

	t.Run("stabilises and produces derived fact", func(t *testing.T) {
		wm := NewWorkingMemory()
		wm.Set("type", "event")

		iter, err := eng.RunToFixpoint(wm, 8)
		if err != nil {
			t.Fatalf("unexpected error: %#v", err)
		}

		v, ok := wm.Get("bar")
		if !ok {
			t.Fatal("expected `bar` to exist")
		}

		if s, ok := v.(string); !ok || s != "done" {
			t.Fatalf(`expected bar="done", got %#v`, v)
		}

		t.Logf("Final state after %d iterations:", iter)
		t.Logf("%#v", wm.(*workingMemory).facts)
	})

	cmplr.Export(os.Stdout)

}

// * engine_test.go ends here.
