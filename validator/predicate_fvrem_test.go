// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// predicate_fvrem_test.go --- FVREM tests.
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
	"fmt"
	"reflect"
	"sync"
	"testing"

	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	testFVREMStructType reflect.Type
	testFVREMStructOnce sync.Once
)

// * Code:

// ** Types:

type testFVREMStruct struct {
	String string
}

func (t *testFVREMStruct) ReflectType() reflect.Type {
	testFVREMStructOnce.Do(func() {
		testFVREMStructType = reflect.TypeOf(t).Elem()
	})

	return testFVREMStructType
}

// ** Tests:

func TestFVREMPredicate(t *testing.T) {
	input := &testFVREMStruct{String: "this is a string"}

	field := "String"
	tests := []struct {
		regex string
		want  bool
	}{
		{"(?i)str[iao]ng", true},
		{"STRANG", false},
		{"STRING", false},
		{"\\w+", true},
	}

	inst := &testFVREMStruct{}
	bindings := NewBindings()
	bindings.Build(inst)
	obj, _ := bindings.Bind(input)

	for idx, tt := range tests {
		t.Run(fmt.Sprintf("%02d FVREM(%s)", idx, field), func(t *testing.T) {
			pred, _ := (&FVREMBuilder{}).Build(field, tt.regex, nil, false)
			result := pred.Eval(context.TODO(), obj)

			if result != tt.want {
				t.Errorf("FVREM(%s) = %v, want %v",
					field,
					result,
					tt.want)
			}
		})
	}

	t.Run("Compile error", func(t *testing.T) {
		_, err := (&FVREMBuilder{}).Build(field, "(?!\b*)", nil, false)

		if err == nil {
			t.Fatal("Expecting an error")
		}

		if !errors.Is(err, ErrRegexpParse) {
			t.Errorf("Unexpected error: %#v", err)
		}
	})

	t.Run("Invalid regexp", func(t *testing.T) {
		_, err := (&FVREMBuilder{}).Build(field, "", nil, false)

		if err == nil {
			t.Fatal("Expecting an error")
		}

		if !errors.Is(err, ErrInvalidRegexp) {
			t.Errorf("Unexpected error: %#v", err)
		}
	})
}

// * predicate_fvrem_test.go ends here.
