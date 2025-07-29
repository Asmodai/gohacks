// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
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
//go:build amd64 || arm64 || riscv64

// * Comments:
//
//

// * Package:

package dynworker

// * Imports:

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/tozd/go/errors"

	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/math"
)

// * Constants:

const (
	// Default average process time.
	defaultAverageProcessTime = 100 * time.Millisecond

	// Default number of worker channels.
	defaultWorkerChannels int64 = 1000

	defaultScaleCooldown = 3 * time.Second

	defaultHystersisThreshold = 2

	defaultMaxScaleDown = 4

	smoothingFactor = 0.2
)

// * Variables:

var (
	//nolint:gochecknoglobals
	activeWorkers = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dynworker_active_workers",
			Help: "Number of active workers",
		},
		[]string{"pool"},
	)

	//nolint:gochecknoglobals
	tasksTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynworker_tasks_total",
			Help: "Total tasks processed",
		},
		[]string{"pool"},
	)

	//nolint:gochecknoglobals
	taskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "dynworker_task_duration_seconds",
			Help: "Histogram of task processing durations",
		},
		[]string{"pool"},
	)

	//nolint:gochecknoglobals
	totalScaledUp = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynworker_scaled_up_total",
			Help: "Total times the worker pool has scaled up",
		},
		[]string{"pool"},
	)

	//nolint:gochecknoglobals
	totalScaledDown = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "dynworker_scaled_down_total",
			Help: "Total times the worker pool has scaled down",
		},
		[]string{"pool"},
	)

	//nolint:gochecknoglobals
	prometheusInitOnce sync.Once

	ErrNotTask error = errors.Base("task pool entity is not a task")
)

// * Code:

// ** Interfaces:

// Worker pool interface.
type WorkerPool interface {
	// Start the worker pool.
	Start()

	// Stop the worker pool.
	Stop()

	// Submit a task to the worker pool.
	Submit(UserData) error

	// Return the number of current workers in the pool.
	WorkerCount() int64

	// Return the minimum number of workers in the pool.
	MinWorkers() int64

	// Return the maximum number of workers in the pool.
	MaxWorkers() int64

	// Set the minimum number of workers to the given value.
	SetMinWorkers(int64)

	// Set the maximum number of workers to the given value.
	SetMaxWorkers(int64)

	// Set the task callback function.
	SetTaskFunction(TaskFn)

	// Set the task scaler function.
	SetScalerFunction(ScalerFn)
}

// ** Types:

type workerPool struct {
	name       string
	input      chan *Task
	minWorkers int64
	maxWorkers int64

	scaleUpCh   chan struct{}
	scaleDownCh chan struct{}

	shutdownChans []chan struct{}
	shutdownLock  sync.Mutex

	lastScaleTime       time.Time
	scaleCooldown       time.Duration
	smoothedRequired    atomic.Int64
	hysteresisThreshold int64
	maxScaleDown        int64

	processFn TaskFn
	scalerFn  ScalerFn

	wg       sync.WaitGroup
	taskPool *sync.Pool

	ctx    context.Context
	cancel context.CancelFunc
	lgr    logger.Logger
	config *Config

	workerCount atomic.Int64
	avgProcTime atomic.Int64

	activeWorkersMetric   prometheus.Gauge
	tasksTotalMetric      prometheus.Counter
	taskDurationMetric    prometheus.Observer
	totalScaledUpMetric   prometheus.Counter
	totalScaledDownMetric prometheus.Counter
}

// ** Methods:

// Return the number of current workers in the pool.
func (obj *workerPool) WorkerCount() int64 {
	return obj.workerCount.Load()
}

// Return the minimum number of workers in the pool.
func (obj *workerPool) MinWorkers() int64 {
	return obj.minWorkers
}

// Set the minimum number of workers to the given value.
func (obj *workerPool) SetMinWorkers(val int64) {
	obj.minWorkers = val
}

// Set the maximum number of workers to the given value.
func (obj *workerPool) SetMaxWorkers(val int64) {
	obj.maxWorkers = val
}

// Return the maximum number of workers in the pool.
func (obj *workerPool) MaxWorkers() int64 {
	return obj.maxWorkers
}

// Set the task callback function.
func (obj *workerPool) SetTaskFunction(workfn TaskFn) {
	obj.processFn = workfn
}

func (obj *workerPool) SetScalerFunction(scalerfn ScalerFn) {
	obj.scalerFn = scalerfn
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

	obj.shutdownLock.Lock()
	killCount := int64(len(obj.shutdownChans))
	obj.shutdownLock.Unlock()

	close(obj.input)
	obj.killWorkers(killCount)
	obj.wg.Wait()
}

// Submit a task to the worker pool.
func (obj *workerPool) Submit(userData UserData) error {
	// Use a pool of task objects.
	task, ok := obj.taskPool.Get().(*Task)
	if !ok {
		return errors.WithStack(ErrNotTask)
	}

	*task = Task{
		parent: obj.ctx,
		logger: obj.lgr,
		data:   userData,
	}

	select {
	case obj.input <- task:
		return nil

	case <-obj.ctx.Done():
		return context.Canceled
	}
}

// Spawn a new worker.
//
//nolint:funlen
func (obj *workerPool) spawnWorker() {
	obj.wg.Add(1)
	obj.workerCount.Add(1)
	obj.activeWorkersMetric.Inc()

	// KillChan?  Japanese mascot for... uh...
	killChan := make(chan struct{})

	obj.shutdownLock.Lock()
	obj.shutdownChans = append(obj.shutdownChans, killChan)
	obj.shutdownLock.Unlock()

	go func() {
		defer func() {
			obj.wg.Done()
			obj.workerCount.Add(-1)
			obj.activeWorkersMetric.Dec()

			obj.shutdownLock.Lock()
			for idx, kchan := range obj.shutdownChans {
				if kchan == killChan {
					obj.shutdownChans = append(
						obj.shutdownChans[:idx],
						obj.shutdownChans[idx+1:]...,
					)

					break
				}
			}
			obj.shutdownLock.Unlock()
		}()

		idleTimer := time.NewTimer(obj.config.IdleTimeout)
		defer idleTimer.Stop()

		for {
			select {
			case <-obj.ctx.Done():
				return

			case <-killChan:
				return

			case task := <-obj.input:
				if task == nil {
					return
				}

				start := time.Now()

				_ = obj.processFn(task)
				obj.tasksTotalMetric.Inc()

				elapsed := time.Since(start).Nanoseconds()

				runtime.Gosched() // Give Go some time.

				// Update metrics.
				obj.updateAvgProcTime(elapsed)
				obj.taskDurationMetric.Observe(float64(elapsed))

				// Reset idle timeout.
				idleTimer.Reset(obj.config.IdleTimeout)

				// Reset and put the task back in the pool.
				task.reset()
				obj.taskPool.Put(task)

			case <-idleTimer.C:
				current := obj.workerCount.Load()
				if current > obj.minWorkers {
					obj.lgr.Info(
						"Worker idle timeout.",
						"type", "dynworker",
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

// Kill the given number of workers.
func (obj *workerPool) killWorkers(num int64) {
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	for idx := int64(0); idx < num && len(obj.shutdownChans) > 0; idx++ {
		kchan := obj.shutdownChans[0]

		close(kchan) // Signal death.

		obj.shutdownChans = obj.shutdownChans[1:]
	}
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

// Compute the required number of workers.
func (obj *workerPool) computeRequiredWorkers() int64 {
	if obj.scalerFn != nil {
		return math.ClampI64(
			int64(obj.scalerFn()),
			obj.minWorkers,
			obj.maxWorkers,
		)
	}

	queued := len(obj.input)

	avg := time.Duration(obj.avgProcTime.Load())
	if avg == 0 {
		avg = defaultAverageProcessTime
	}

	raw := int64(float64(queued)*avg.Seconds()) + 1

	return math.ClampI64(raw, obj.minWorkers, obj.maxWorkers)
}

// Smooth the number of required numbers.
func (obj *workerPool) smoothRequiredWorkers(raw int64) int64 {
	old := obj.smoothedRequired.Load()
	if old == 0 {
		old = raw
	}

	smoothed := int64(float64(raw)*smoothingFactor +
		float64(old)*(1-smoothingFactor))

	obj.smoothedRequired.Store(smoothed)

	return smoothed
}

// Should we scale?
//
// Attempts to prevent hysteresis.
func (obj *workerPool) shouldScale(required, current int64) bool {
	return math.AbsI64(required-current) >= obj.hysteresisThreshold
}

// Scale the number of workers up to the given number.
func (obj *workerPool) scaleUp(num, current, required int64) {
	obj.lgr.Info(
		"Scaling up workers.",
		"type", "dynworker",
		"pool", obj.name,
		"current", current,
		"required", required,
		"new", num,
	)

	obj.totalScaledUpMetric.Inc()

	for range num {
		obj.spawnWorker()
	}
}

// Scale the number of workers down to the given number.
func (obj *workerPool) scaleDown(num, current, required int64) {
	if num > obj.maxScaleDown {
		num = obj.maxScaleDown
	}

	obj.lgr.Info(
		"Scaling down workers.",
		"type", "dynworker",
		"pool", obj.name,
		"current", current,
		"required", required,
		"new", num,
	)

	obj.totalScaledDownMetric.Inc()
	obj.killWorkers(num)
}

// Check if we need to scale the number of workers if required.
//
// Note, this will not actively terminate workers should the number require
// scaling down, rather it will let workers terminate through either completion
// or idle timeout.
func (obj *workerPool) scaleCheck() {
	now := time.Now()

	if now.Sub(obj.lastScaleTime) < obj.scaleCooldown {
		return
	}

	current := obj.workerCount.Load()
	rawRequired := obj.computeRequiredWorkers()
	required := obj.smoothRequiredWorkers(rawRequired)

	if !obj.shouldScale(required, current) {
		return
	}

	obj.lastScaleTime = now
	delta := required - current

	if delta > 0 {
		obj.scaleUp(delta, current, required)
	} else {
		obj.scaleDown(-delta, current, required)
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
//
// The provided context must have `logger.Logger` in its user value.
// See `contextdi` and `logger.SetLogger`.
func NewWorkerPool(ctx context.Context, config *Config) WorkerPool {
	if config == nil {
		panic("invalid worker configuration")
	}

	lgr := logger.MustGetLogger(ctx)

	nctx, cancel := context.WithCancel(ctx)

	label := prometheus.Labels{"pool": config.Name}

	taskPool := &sync.Pool{
		New: func() any {
			return &Task{}
		},
	}

	return &workerPool{
		name:                  config.Name,
		input:                 make(chan *Task, defaultWorkerChannels),
		minWorkers:            config.MinWorkers,
		maxWorkers:            config.MaxWorkers,
		scaleUpCh:             make(chan struct{}, 1),
		scaleDownCh:           make(chan struct{}, 1),
		processFn:             config.WorkerFunc,
		scalerFn:              config.ScalerFunc,
		ctx:                   nctx,
		cancel:                cancel,
		lgr:                   lgr,
		config:                config,
		taskPool:              taskPool,
		lastScaleTime:         time.Now(),
		scaleCooldown:         defaultScaleCooldown,
		hysteresisThreshold:   defaultHystersisThreshold,
		maxScaleDown:          defaultMaxScaleDown,
		activeWorkersMetric:   activeWorkers.With(label),
		tasksTotalMetric:      tasksTotal.With(label),
		taskDurationMetric:    taskDuration.With(label),
		totalScaledUpMetric:   totalScaledUp.With(label),
		totalScaledDownMetric: totalScaledDown.With(label),
	}
}

// Initialise Prometheus metrics for this module.
func InitPrometheus() {
	prometheusInitOnce.Do(func() {
		prometheus.MustRegister(
			activeWorkers,
			tasksTotal,
			taskDuration,
			totalScaledUp,
			totalScaledDown,
		)
	})
}

// * dynworker.go ends here.
