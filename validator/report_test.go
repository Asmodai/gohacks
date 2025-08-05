// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// report_test.go --- Report tests.
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
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/Asmodai/gohacks/contextdi"
	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/logger"
	mlogger "github.com/Asmodai/gohacks/mocks/logger"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/mock/gomock"
)

// * Constants:

const (
	rulesFTEQ string = `
- name: "'One' must be int64"
  conditions:
    - attribute: one
      operator: field-type-equal
      value: int64
    - attribute: one
      operator: field-value-equal
      value: 42
    - attribute: one
      operator: field-value-not-equal
      value: 9001
    - attribute: one
      operator: field-value-in
      value: [40, 41, 42, 43]
    - attribute: one
      operator: field-type-in
      value: [int8, int16, int32, int64]
  action:
    perform: ignore`
)

// * Variables:

var (
	dummyStructureType reflect.Type
	dummyStructureOnce sync.Once
)

// * Code:

// ** Test structure:

type DummyStructure struct {
	One any            `json:"one"`
	Two map[string]int `json:"name_to_id"`
}

func (ds *DummyStructure) ReflectType() reflect.Type {
	dummyStructureOnce.Do(func() {
		dummyStructureType = reflect.TypeOf(ds).Elem()
	})

	return dummyStructureType
}

// ** Mock actions:

// *** Type:

type MockActions struct {
	hasRunLog bool
}

// *** Methods:

func (obj *MockActions) Builder(fn string, params dag.ActionParams) (dag.ActionFn, error) {
	normalised := strings.ToLower(fn)

	switch normalised {
	case "ignore":
		return func(_ context.Context, _ dag.Filterable) {}, nil

	case "log":
		return obj.logAction(params)

	default:
		return nil, errors.WithMessagef(dag.ErrUnknownBuiltin, "%q", fn)
	}
}

func (obj *MockActions) logAction(params dag.ActionParams) (afn dag.ActionFn, err error) {
	msgparam, ok := params["message"]
	if !ok {
		afn = nil
		err = errors.WithMessage(dag.ErrMissingParam,
			`Parameter "message"`)

		return
	}

	msg, ok := msgparam.(string)
	if !ok {
		afn = nil
		err = errors.WithStack(dag.ErrExpectedString)

		return
	}

	err = nil
	afn = func(ctx context.Context, input dag.Filterable) {
		lgr := logger.MustGetLogger(ctx)
		lgr.Info(
			msg,
			"src", "log_action",
		)

		obj.hasRunLog = true
	}

	return
}

// ** Tests:

func TestShit(t *testing.T) {
	var rules []dag.RuleSpec

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

	ctx, err := logger.SetLogger(context.TODO(), lgr)
	if err != nil {
		t.Fatalf("Could not set logger DI: %#v", err)
	}

	ctx, err = contextdi.SetDebugMode(ctx, true)
	if err != nil {
		t.Fatalf("Could not set debug flag to DI: %#v", err)
	}

	rules, err = dag.ParseFromYAML(rulesFTEQ)
	if err != nil {
		t.Fatalf("YAML parse error: %#v", err)
	}

	mact := &MockActions{
		hasRunLog: false,
	}

	compiler := dag.NewCompilerWithPredicates(
		ctx,
		mact,
		BuildPredicateDict(),
	)
	issues := compiler.Compile(rules)

	if len(issues) > 0 {
		t.Logf("Compiler issues:")
		for idx := range issues {
			t.Logf("%d: %s", idx+1, issues[idx])
		}
	}

	bindings := NewBindings()
	inst := &DummyStructure{One: int64(42), Two: map[string]int{}}

	_, ok := bindings.Build(inst)
	if !ok {
		t.Fatalf("Could not create descriptor")
	}
	obj, _ := bindings.Bind(inst)

	compiler.Evaluate(obj)

}

// ** Benchmarks:

func BenchmarkCompiler(b *testing.B) {
	inst := &DummyStructure{One: int64(42), Two: map[string]int{}}
	data := BuildDescriptor(reflect.TypeOf(&DummyStructure{}))
	bound := &BoundObject{
		Descriptor: data,
		Binding:    inst,
	}

	mocker := gomock.NewController(b)
	defer mocker.Finish()

	lgr := mlogger.NewMockLogger(mocker)
	lgr.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	ctx, err := logger.SetLogger(context.TODO(), lgr)
	if err != nil {
		b.Fatalf("Could not set logger DI: %#v", err)
	}

	rules, err := dag.ParseFromYAML(rulesFTEQ)
	if err != nil {
		b.Fatalf("YAML: %#v", err)
	}

	mact := &MockActions{
		hasRunLog: false,
	}

	compiler := dag.NewCompilerWithPredicates(
		ctx,
		mact,
		BuildPredicateDict(),
	)

	b.Run("Compile()", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			compiler.Compile(rules)
		}
	})

	b.Run("Evaluate()", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			compiler.Evaluate(bound)
		}
	})
}

// * report_test.go ends here.
