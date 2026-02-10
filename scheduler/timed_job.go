// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: NONE
//
// timed_job.go --- Timed jobs.
//
// Copyright (c) 2026 Fortra, LLC.
//
// Author:     Paul Ward <paul.ward@fortra.com>
// Maintainer: Paul Ward <paul.ward@fortra.com>
//
// PROPRIETARY AND CONFIDENTIAL
//
// This source code and any related files contain proprietary information
// belonging exclusively to Paul Ward (or your company name).
//
// Disclosure, distribution, or reproduction of this code is strictly prohibited
// without prior written consent.
//
// This software contains trade secrets and confidential information
// that is protected under applicable intellectual property laws.
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
