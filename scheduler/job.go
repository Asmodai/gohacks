// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// job.go --- Timing wheel job definition.
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
//
//mock:yes

// * Comments:

// * Package:

package scheduler

// * Imports:

import (
	"context"

	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	ErrJobRunAtZero       = errors.Base("job run-at is zero")
	ErrJobNoTarget        = errors.Base("no job target specified")
	ErrJobAmbiguousTarget = errors.Base("job target is ambiguous")
)

// * Code:

// ** Interface:

type Job interface {
	Validate() error
	Resolve(context.Context) error

	Object() Task
	Function() JobFn
}

// ** Types:

// Job function.
type JobFn func(context.Context) error

// Job structure.
type job struct {
	obj Task  // Object to execute.
	fn  JobFn // Function to execute.
}

// ** Methods:

func (j job) Object() Task {
	return j.obj
}

func (j job) Function() JobFn {
	return j.fn
}

// Validate a job.
//
// A job is valid if it has either an object or a function.
//
// If the job has both an object and a function, it is invalid.
func (j job) Validate() error {
	// Check identifier.
	switch {
	case j.obj == nil && j.fn == nil:
		return errors.WithStack(ErrJobNoTarget)

	case j.obj != nil && j.fn != nil:
		return errors.WithStack(ErrJobAmbiguousTarget)

	default:
		return nil
	}
}

// Resolve and call the job's function.
//
// This will either invoke the `Execute` method on a job, or call the job's
// function.
func (j job) Resolve(ctx context.Context) error {
	if j.obj != nil {
		return errors.WithStack(j.obj.Execute(ctx))
	}

	// If we get here and `fn` is nil, complain loudly about it.
	if j.fn == nil {
		return errors.WithStack(ErrJobNoTarget)
	}

	return errors.WithStack(j.Function()(ctx))
}

// ** Functions:

// Create a new job.
func MakeJob(obj Task, fn JobFn) Job {
	return &job{
		obj: obj,
		fn:  fn,
	}
}

// Insert a job into a list of jobs.
func InsertJob(jobs []Job, njob Job) ([]Job, error) {
	if err := njob.Validate(); err != nil {
		// Do not return nil, ensure caller doesn't nuke their jobs
		// list.
		return jobs, errors.WithStack(err)
	}

	return append(jobs, njob), nil
}

// * job.go ends here.
