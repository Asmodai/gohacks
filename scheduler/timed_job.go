// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: NONE
//
// timed_job.go --- Timed jobs.
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
	"sort"
	"time"

	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Interface:

type TimedJob interface {
	Job

	RunAt() time.Time
}

// ** Types:

type timedjob struct {
	job

	runAt time.Time
}

// ** Methods:

func (t timedjob) RunAt() time.Time {
	return t.runAt
}

// Validate a job.
//
// A job is valid if it has a non-zero time, and either an object or a
// function.
//
// If the job has both an object and a function, it is invalid.
func (t timedjob) Validate() error {
	if err := t.job.Validate(); err != nil {
		return err
	}

	if t.runAt.IsZero() {
		return errors.WithStack(ErrJobRunAtZero)
	}

	return nil
}

// ** Functions:

// Create a new timed job.
func MakeTimedJob(runAt time.Time, obj Task, fn JobFn) TimedJob {
	return &timedjob{
		job: job{
			obj: obj,
			fn:  fn,
		},
		runAt: runAt,
	}
}

// Insert a timed job into a list of jobs.
func InsertTimedJob(jobs []TimedJob, njob TimedJob) ([]TimedJob, error) {
	if err := njob.Validate(); err != nil {
		// Do not return nil, ensure caller doesn't nuke their jobs
		// list.
		return jobs, errors.WithStack(err)
	}

	place := sort.Search(len(jobs), func(idx int) bool {
		// Ensure FIFO behaviour.
		return jobs[idx].RunAt().After(njob.RunAt())
	})

	jobs = append(jobs, nil)
	copy(jobs[place+1:], jobs[place:])
	jobs[place] = njob

	return jobs, nil
}

// * timed_job.go ends here.
