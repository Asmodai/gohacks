// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// job_test.go --- Job tests.
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

package scheduler

// * Imports:

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

// * Variables:

// * Code:

// ** Types:

type SomeStruct struct {
	ID   uint64
	Name string
}

func (s *SomeStruct) Execute(_ context.Context) error {
	return nil
}

// ** Functions:

func FnExecute(_ context.Context) error {
	return nil
}

// ** Tests:
func TestJobValidation(t *testing.T) {
	tests := []struct {
		jobdef Job
		err    error
	}{
		{
			MakeJob(nil, nil),
			ErrJobNoTarget,
		}, {
			MakeJob(&SomeStruct{}, FnExecute),
			ErrJobAmbiguousTarget,
		}, {
			MakeJob(&SomeStruct{}, nil),
			nil,
		}, {
			MakeJob(nil, FnExecute),
			nil,
		},
	}

	for idx, elt := range tests {
		t.Run(fmt.Sprintf("Test %02d", idx), func(t *testing.T) {
			err := elt.jobdef.Validate()
			if (err != nil && elt.err != nil) && errors.Is(err, errors.Unwrap(elt.err)) {
				t.Errorf("Test %02d: expected %#v but got %#v",
					idx,
					elt.err,
					err)
			}
		})
	}
}

func TestJobExecute(t *testing.T) {
	tests := []struct {
		jobdef Job
		err    error
	}{
		{
			MakeJob(&SomeStruct{}, nil),
			nil,
		}, {
			MakeJob(nil, FnExecute),
			nil,
		}, {
			MakeJob(nil, nil),
			ErrJobNoTarget,
		},
	}

	for idx, elt := range tests {
		t.Run(fmt.Sprintf("Test %02d", idx), func(t *testing.T) {
			err := elt.jobdef.Resolve(context.TODO())
			if (err != nil && elt.err != nil) && errors.Is(err, errors.Unwrap(elt.err)) {
				t.Errorf("Test %02d: expected %#v but got %#v",
					idx,
					elt.err,
					err)
			}
		})
	}
}

func TestJobInsertion(t *testing.T) {
	var (
		jobs  = make([]Job, 0)
		err   error
		names = []string{
			"one",
			"two",
			"three",
		}
	)

	t.Run("Insert", func(t *testing.T) {
		for idx := range names {
			jobs, err = InsertJob(
				jobs,
				job{obj: &SomeStruct{Name: names[idx]}, fn: nil})

			if err != nil {
				t.Fatalf("Got error %#v", err)
			}
		}
	})

	t.Run("Sorted?", func(t *testing.T) {
		for idx := range jobs {
			obj := jobs[idx].Object().(*SomeStruct)
			if obj.Name != names[idx] {
				t.Fatalf("Invalid name: Got %#v want %#v",
					obj,
					names[idx])
			}
		}
	})
}

// * job_test.go ends here.
