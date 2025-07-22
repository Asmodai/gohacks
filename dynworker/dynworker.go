// -*- Mode: Go -*-
//
// dynworker.go --- Dynamic worker.
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
//
//mock:yes

// * Comments:
//
//

// * Package:

package dynworker

// * Imports:

import (
	"github.com/Asmodai/gohacks/logger"

	"github.com/prometheus/client_golang/prometheus"

	"context"
	"sync"
	"sync/atomic"
	"time"
)

// * Constants:

const (
	// Default average process time.
	defaultAverageProcessTime time.Duration = 100 * time.Millisecond

	// Default number of worker channels.
	defaultWorkerChannels int = 1000
)

// * Variables:

//nolint:gochecknoglobals
var (
	activeWorkers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dynworker_active_workers",
			Help: "Number of active workers",
		},
		[]string{"pool"},
	)

	tasksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynworker_tasks_total",
			Help: "Total tasks processed",
		},
		[]string{"pool"},
	)

	taskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "dynworker_task_duration_seconds",
			Help: "Histogram of task processing durations",
		},
		[]string{"pool"},
	)

	totalScaledUp = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynworker_scaled_up_total",
			Help: "Total times the worker pool has scaled up",
		},
		[]string{"pool"},
	)

	totalScaledDown = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynworker_scaled_down_total",
			Help: "Total times the worker pool has scaled down",
		},
		[]string{"pool"},
	)
)

// * Code:

// ** Types:

// Task data type.
type Task any

// Type of functions executed by workers.
type TaskFn func(Task) error

// Worker pool interface.
type WorkerPool interface {
	// Start the worker pool.
	Start()

	// Stop the worker pool.
	Stop()

	// Submit a task to the worker pool.
	Submit(Task) error

	// Return the number of current workers in the pool.
	WorkerCount() int32

	// Return the minimum number of workers in the pool.
	MinWorkers() int32

	// Return the maximum number of workers in the pool.
	MaxWorkers() int32
}

type workerPool struct {
	name          string
	input         chan Task
	minWorkers    int32
	maxWorkers    int32
	scaleUpStep   int32
	scaleDownStep int32

	scaleUpCh   chan struct{}
	scaleDownCh chan struct{}

	processFn TaskFn

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	lgr    logger.Logger
	config *Config

	workerCount atomic.Int32
	avgProcTime atomic.Int64

	activeWorkersMetric   prometheus.Gauge
	tasksTotalMetric      prometheus.Counter
	taskDurationMetric    prometheus.Observer
	totalScaledUpMetric   prometheus.Counter
	totalScaledDownMetric prometheus.Counter
}

// ** Methods:

// Return the number of current workers in the pool.
func (obj *workerPool) WorkerCount() int32 {
	return obj.workerCount.Load()
}

// Return the minimum number of workers in the pool.
func (obj *workerPool) MinWorkers() int32 {
	return obj.minWorkers
}

// Return the maximum number of workers in the pool.
func (obj *workerPool) MaxWorkers() int32 {
	return obj.maxWorkers
}

// Start the worker pool.
func (obj *workerPool) Start() {
	for range obj.minWorkers {
		obj.spawnWorker()
	}

	go obj.scaler()
}

// Stop the worker pool.
func (obj *workerPool) Stop() {
	obj.cancel()
	obj.wg.Wait()
	close(obj.input)
}

// Submit a task to the worker pool.
func (obj *workerPool) Submit(task Task) error {
	select {
	case obj.input <- task:
		return nil

	case <-obj.ctx.Done():
		return context.Canceled
	}
}

// Spawn a new worker.
func (obj *workerPool) spawnWorker() {
	obj.wg.Add(1)
	obj.workerCount.Add(1)
	obj.activeWorkersMetric.Inc()

	go func() {
		defer func() {
			obj.wg.Done()
			obj.workerCount.Add(-1)
			obj.activeWorkersMetric.Dec()
		}()

		idleTimer := time.NewTimer(obj.config.IdleTimeout)
		defer idleTimer.Stop()

		for {
			select {
			case <-obj.ctx.Done():
				return

			case task := <-obj.input:
				start := time.Now()
				_ = obj.processFn(task)
				elapsed := time.Since(start).Nanoseconds()

				obj.updateAvgProcTime(elapsed)
				obj.taskDurationMetric.Observe(float64(elapsed))
				idleTimer.Reset(obj.config.IdleTimeout)

			case <-idleTimer.C:
				current := obj.workerCount.Load()
				if current > obj.minWorkers {
					obj.lgr.Info(
						"dynworker: Worker idle timeout.",
						"pool", obj.name,
						"remaining", current-1,
					)

					return
				}

				idleTimer.Reset(obj.config.IdleTimeout)
			}
		}
	}()
}

// Initiate a scale check at a 1 second interval.
func (obj *workerPool) scaler() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-obj.ctx.Done():
			return

		case <-ticker.C:
			obj.scaleCheck()
		}
	}
}

// Check if we need to scale the number of workers if required.
//
// Note, this will not actively terminate workers should the number require
// scaling down, rather it will let workers terminate through either completion
// or idle timeout.
func (obj *workerPool) scaleCheck() {
	queued := len(obj.input)
	current := obj.workerCount.Load()
	avg := time.Duration(obj.avgProcTime.Load())

	// If we don't have an average process time, set one to 100ms.
	if avg == 0 {
		avg = defaultAverageProcessTime
	}

	// Rough required workers = queue * avg / interval.
	required := int32(float64(queued)*avg.Seconds()) + 1

	// Clamp if lower than minimum workers.
	if required < obj.minWorkers {
		required = obj.minWorkers
	}

	// Clamp if higher than maximum workers.
	if required > obj.maxWorkers {
		required = obj.maxWorkers
	}

	// Can we scale?
	if required > current {
		toSpawn := required - current

		obj.lgr.Info(
			"dynworker: scaling up workers.",
			"pool", obj.name,
			"current", current,
			"required", required,
			"new", toSpawn,
		)
		obj.totalScaledUpMetric.Inc()

		for range toSpawn {
			obj.spawnWorker()
		}
	} else {
		obj.lgr.Info(
			"dynworker: scaling down workers.",
			"pool", obj.name,
			"current", current,
			"required", required,
			"note", "noop, workers will die.",
		)
		obj.totalScaledDownMetric.Inc()
	}
}

// Update the average time spent processing.
func (obj *workerPool) updateAvgProcTime(latest int64) {
	const alpha = 0.2

	old := obj.avgProcTime.Load()
	if old == 0 {
		obj.avgProcTime.Store(latest)

		return
	}

	newAvg := int64(float64(latest)*alpha + float64(old)*(1-alpha))
	obj.avgProcTime.Store(newAvg)
}

// ** Functions:

// Create a new worker pool.
func NewWorkerPool(config *Config, workfn TaskFn) WorkerPool {
	if config == nil {
		panic("invalid worker configuration")
	}

	ctx, cancel := context.WithCancel(config.Parent)

	label := prometheus.Labels{"pool": config.Name}

	return &workerPool{
		name:                  config.Name,
		input:                 make(chan Task, defaultWorkerChannels),
		minWorkers:            config.MinWorkers,
		maxWorkers:            config.MaxWorkers,
		scaleUpStep:           1,
		scaleDownStep:         1,
		scaleUpCh:             make(chan struct{}, 1),
		scaleDownCh:           make(chan struct{}, 1),
		processFn:             workfn,
		ctx:                   ctx,
		cancel:                cancel,
		lgr:                   config.Logger,
		config:                config,
		activeWorkersMetric:   activeWorkers.With(label),
		tasksTotalMetric:      tasksTotal.With(label),
		taskDurationMetric:    taskDuration.With(label),
		totalScaledUpMetric:   totalScaledUp.With(label),
		totalScaledDownMetric: totalScaledDown.With(label),
	}
}

// ** Initialisation:

//nolint:gochecknoinits
func init() {
	prometheus.MustRegister(
		activeWorkers,
		tasksTotal,
		taskDuration,
		totalScaledUp,
		totalScaledDown,
	)
}

// * dynworker.go ends here.
