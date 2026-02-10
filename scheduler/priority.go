// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// priority.go --- Priority scheduler.
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
	"sync"
	"time"

	"github.com/Asmodai/gohacks/errx"
	"github.com/Asmodai/gohacks/health"
	"github.com/Asmodai/gohacks/logger"
	"github.com/prometheus/client_golang/prometheus"
)

// * Constants:

const (
	// Number of seconds before a task is considered late.
	LateTaskDelay time.Duration = 5 * time.Second
)

// * Variables:

var (
	//nolint:gochecknoglobals
	activeTasks = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "priority_scheduler_active_tasks",
			Help: "Number of active tasks",
		},
		[]string{"priority_scheduler"},
	)

	//nolint:gochecknoglobals,mnd
	taskLateness = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "priority_scheduler_task_lateness_seconds",
			Help:    "Histogram of task lateness",
			Buckets: prometheus.ExponentialBuckets(0.005, 2, 12),
		},
		[]string{"priority_scheduler"},
	)

	//nolint:gochecknoglobals
	taskDispatchedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "priority_scheduler_task_dispatched_total",
			Help: "Count of tasks dispatched via the scheduler",
		},
		[]string{"priority_scheduler"},
	)

	//nolint:gochecknoglobals
	prometheusInitOnce sync.Once
)

// * Code:

// ** Types:

/*
Priority scheduler

Instance is single-use; create a new instance to restart.
*/
type Priority struct {
	lgr                logger.Logger
	ctx                context.Context
	hlth               health.Reporter
	activeTasksMetric  prometheus.Gauge
	taskLatenessMetric prometheus.Observer
	taskDispatchTotal  prometheus.Counter
	cancel             context.CancelFunc
	hlthTicker         func(time.Duration) health.Ticker
	addCh              chan TimedJob
	workCh             chan TimedJob
	done               chan struct{}
	name               string
	hlthTickPeriod     time.Duration
	stopOnce           sync.Once
	startOnce          sync.Once
}

// ** Methods:

// Return the name of the priority scheduler.
func (s *Priority) Name() string {
	return s.name
}

// Return the health reporter for the priority scheduler.
func (s *Priority) Health() health.Reporter {
	return s.hlth
}

// Add a task to the priority scheduler.
func (s *Priority) Submit(task TimedJob) error {
	select {
	case <-s.ctx.Done():
		return errx.WithStack(s.ctx.Err())

	case s.addCh <- task:
		return nil
	}
}

// Get the current work channel for the priority scheduler.
func (s *Priority) Work() <-chan TimedJob {
	return s.workCh
}

// Get the next task in the work channel.
func (s *Priority) Next(ctx context.Context) (TimedJob, bool) {
	var zero TimedJob

	select {
	case <-s.ctx.Done():
		return zero, false

	case <-ctx.Done():
		return zero, false

	case job, ok := <-s.workCh:
		return job, ok
	}
}

// Start the scheduler.
//
// Must be called once before use.
func (s *Priority) Start() {
	s.startOnce.Do(func() {
		go s.run()
	})
}

// Stop the priority scheduler.
func (s *Priority) Stop() {
	s.stopOnce.Do(func() {
		s.cancel()
	})
}

// Has the priority scheduler done processing?
func (s *Priority) Done() <-chan struct{} {
	return s.done
}

// Wait for the goroutine created by `Start`.
func (s *Priority) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errx.WithStack(ctx.Err())

	case <-s.done:
		return nil
	}
}

// Run the scheduler.
//
//nolint:cyclop,gocognit,funlen
func (s *Priority) run() {
	var (
		jobs  []TimedJob
		timer *time.Timer
		err   error
	)

	heartbeatTicker := s.hlthTicker(s.hlthTickPeriod)
	resetTimer := func(dur time.Duration) {
		if dur < 0 {
			dur = 0
		}

		if timer == nil {
			timer = time.NewTimer(dur)

			return
		}

		if !timer.Stop() {
			// Drain.
			select {
			case <-timer.C:
			default:
			}
		}

		timer.Reset(dur)
	}

	// Close and clean up on our way out.
	defer func() {
		close(s.workCh)
		close(s.done)

		s.activeTasksMetric.Set(0)

		heartbeatTicker.Stop()

		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		// If we have no jobs, then just check whether we're done or
		// if a job needs adding, and then continue the loop.
		if len(jobs) == 0 {
			select {
			case <-s.ctx.Done():
				return

			case <-heartbeatTicker.Channel():
				s.hlth.Tick()

			case njob, ok := <-s.addCh:
				if !ok {
					return
				}

				jobs, err = InsertTimedJob(jobs, njob)
				if err != nil {
					s.lgr.Fatalf(err.Error())
				}

				s.activeTasksMetric.Set(float64(len(jobs)))
			}

			continue
		}

		next := jobs[0].RunAt()
		resetTimer(time.Until(next))

		select {
		case <-s.ctx.Done():
			return

		case <-heartbeatTicker.Channel():
			s.hlth.Tick()

		case <-timer.C:
			now := time.Now()
			idx := 0

			for idx < len(jobs) {
				runAt := jobs[idx].RunAt()

				if runAt.After(now) {
					break
				}

				select {
				case <-s.ctx.Done():
					return

				case s.workCh <- jobs[idx]:
					lateness := now.Sub(runAt)

					switch {
					case lateness < 0:
						lateness = 0

					case lateness > LateTaskDelay:
						// TODO: Log throttling.
						s.lgr.Warn(
							"job dispatched late",
							"delay", lateness,
							"job", jobs[idx],
						)
					}

					s.taskLatenessMetric.Observe(lateness.Seconds())
					s.taskDispatchTotal.Inc()

					idx++

				case <-heartbeatTicker.Channel():
					s.hlth.Tick()
				}
			}

			copy(jobs, jobs[idx:])
			jobs = jobs[:len(jobs)-idx]
			s.activeTasksMetric.Set(float64(len(jobs)))

		case njob, ok := <-s.addCh:
			if !ok {
				return
			}

			jobs, err = InsertTimedJob(jobs, njob)
			if err != nil {
				s.lgr.Fatal(err.Error())
			}

			s.activeTasksMetric.Set(float64(len(jobs)))
		}
	}
}

// ** Functions:

// Return a new priority scheduler instance.
func NewPriority(ctx context.Context, cnf *Config) *Priority {
	if cnf == nil {
		panic("invalid priority scheduler configuration")
	}

	if cnf.Prometheus == nil {
		cnf.Prometheus = prometheus.DefaultRegisterer
	}

	nctx, cancel := context.WithCancel(ctx)
	lgr := logger.MustGetLogger(nctx)

	// Initialise Prometheus.
	InitPrometheus(cnf.Prometheus)

	// Set up a Prometheus label.
	label := prometheus.Labels{"priority_scheduler": cnf.Name}

	return &Priority{
		name:               cnf.Name,
		lgr:                lgr,
		ctx:                nctx,
		cancel:             cancel,
		hlth:               cnf.Health,
		hlthTicker:         health.NewTicker,
		hlthTickPeriod:     cnf.HealthTickPeriod,
		addCh:              make(chan TimedJob, cnf.AddBuffer),
		workCh:             make(chan TimedJob, cnf.WorkBuffer),
		done:               make(chan struct{}),
		activeTasksMetric:  activeTasks.With(label),
		taskLatenessMetric: taskLateness.With(label),
		taskDispatchTotal:  taskDispatchedTotal.With(label),
	}
}

// Initialise Prometheus metrics for this module.
func InitPrometheus(reg prometheus.Registerer) {
	prometheusInitOnce.Do(func() {
		reg.MustRegister(
			activeTasks,
			taskLateness,
			taskDispatchedTotal,
		)
	})
}

// * priority.go ends here.
