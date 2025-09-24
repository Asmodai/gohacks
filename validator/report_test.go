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
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/Asmodai/gohacks/contextdi"
	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/logger"
	mlogger "github.com/Asmodai/gohacks/mocks/logger"
	"go.uber.org/mock/gomock"
)

// * Constants:

const (
	rulesFTEQ string = `
- name: "'One' must be valid"
  conditions:
    - attribute: one
      operator: field-type-in
      value: [int8, int16, int32, int64]
    - attribute: one
      operator: field-value-in
      value: [40, 41, 42, 43]
  failure:
    perform: error
    params:
      message: "'One' is not valid"

- name: "'Two' must be map[string]int and not empty"
  conditions:
    - attribute: two
      operator: field-type-equal
      value: map[string]int
    - attribute: two
      operator: field-value-is-true
  failure:
    perform: error
    params:
      message: "'Two' is not valid"

- name: "'three' must be nil"
  conditions:
    - attribute: three
      operator: field-value-is-nil
  failure:
    perform: error
    params:
      message: "'Three' is not valid"

- name: "'four' must be string and member"
  conditions:
    - attribute: four
      operator: field-type-equal
      value: string
    - attribute: four
      operator: field-value-in
      value: [OK, CRITICAL, WARNING]
  failure:
    perform: error
    params:
      message: "'Four' is not valid"

- name: "'five' must match regex"
  conditions:
    - attribute: five
      operator: field-type-equal
      value: string
    - attribute: five
      operator: field-value-regex-match
      value: ".*coffee.*"
  failure:
    perform: error
    params:
      message: "'Five' is not valid"
`
)

// * Variables:

var (
	dummyStructureType reflect.Type
	dummyStructureOnce sync.Once
)

// * Code:

// ** Test structure:

type UnusedStruct struct {
}

type DummyStructure struct {
	One   any            `json:"one"`
	Two   map[string]int `json:"name_to_id"`
	Three *UnusedStruct
	Four  string
	Five  string
}

func (ds *DummyStructure) ReflectType() reflect.Type {
	dummyStructureOnce.Do(func() {
		dummyStructureType = reflect.TypeOf(ds).Elem()
	})

	return dummyStructureType
}

var (
	testData = &DummyStructure{
		One:   int64(42),
		Two:   map[string]int{"one": 1, "two": 2, "three": 3},
		Three: nil,
		Four:  "CRITICAL",
		Five:  "Must contain coffee in here",
	}

	testInvalid = &DummyStructure{
		One:   int64(76),
		Two:   map[string]int{},
		Three: &UnusedStruct{},
		Four:  "FINE",
		Five:  "Must contain tea in here",
	}
)

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

	compiler := NewValidator(ctx)
	issues := compiler.Compile(rules)

	compiler.Export(os.Stdout)

	if len(issues) > 0 {
		t.Logf("Compiler issues:")
		for idx := range issues {
			t.Logf("%d: %s", idx+1, issues[idx])
		}
		t.Fatal("Fix the compiler issues in your test!")
	}

	bindings := NewBindings()
	_, ok := bindings.Build(&DummyStructure{})
	if !ok {
		t.Fatalf("Could not create descriptor")
	}

	t.Run("valid data", func(t *testing.T) {
		obj, _ := bindings.Bind(testData)
		compiler.Evaluate(obj)
	})

	t.Run("invalid data", func(t *testing.T) {
		obj, _ := bindings.Bind(testInvalid)
		compiler.Evaluate(obj)

		if len(compiler.Failures()) == 0 {
			t.Fatal("No errors generated")
		}

		for idx, val := range compiler.Failures() {
			t.Logf("%02d - %q", idx, val)
		}
	})

}

// ** Benchmarks:

func BenchmarkCompiler(b *testing.B) {
	data := BuildDescriptor(reflect.TypeOf(&DummyStructure{}))
	bound := &BoundObject{
		Descriptor: data,
		Binding:    testData,
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

	compiler := NewValidator(ctx)

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
