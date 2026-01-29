// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// dynworker.go --- Dynamic worker.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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
	gomath "math"
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

	//nolint:gochecknoglobals,mnd
	taskDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "dynworker_task_duration_seconds",
			Help:    "Histogram of task processing durations",
			Buckets: prometheus.ExponentialBuckets(0.005, 2, 12),
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

	// Return the name of the pool.
	Name() string
}

// ** Types:

type workerPool struct {
	lastScaleTime         time.Time
	input                 TaskQueue
	ctx                   context.Context
	lgr                   logger.Logger
	activeWorkersMetric   prometheus.Gauge
	tasksTotalMetric      prometheus.Counter
	taskDurationMetric    prometheus.Observer
	totalScaledUpMetric   prometheus.Counter
	totalScaledDownMetric prometheus.Counter
	processFn             TaskFn
	scalerFn              ScalerFn
	taskPool              *sync.Pool
	cancel                context.CancelFunc
	config                *Config
	name                  string
	shutdownChans         []chan struct{}
	wg                    sync.WaitGroup
	minWorkers            atomic.Int64
	maxWorkers            atomic.Int64
	scaleCooldown         time.Duration
	smoothedRequired      atomic.Int64
	hysteresisThreshold   int64
	maxScaleDown          int64
	workerCount           atomic.Int64
	avgProcTime           atomic.Int64
	shutdownLock          sync.Mutex
}

// ** Methods:

// Return the name of the pool.
func (obj *workerPool) Name() string {
	return obj.name
}

// Return the number of current workers in the pool.
func (obj *workerPool) WorkerCount() int64 {
	return obj.workerCount.Load()
}

// Return the minimum number of workers in the pool.
func (obj *workerPool) MinWorkers() int64 {
	return obj.minWorkers.Load()
}

// Return the maximum number of workers in the pool.
func (obj *workerPool) MaxWorkers() int64 {
	return obj.maxWorkers.Load()
}

// Set the minimum number of workers to the given value.
func (obj *workerPool) SetMinWorkers(val int64) {
	obj.minWorkers.Store(val)
}

// Set the maximum number of workers to the given value.
func (obj *workerPool) SetMaxWorkers(val int64) {
	obj.maxWorkers.Store(val)
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
	for range obj.minWorkers.Load() {
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

	//	close(obj.input)
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

	if err := obj.input.Put(obj.ctx, task); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Spawn a new worker.
func (obj *workerPool) spawnWorker() {
	obj.wg.Add(1)
	obj.workerCount.Add(1)
	obj.activeWorkersMetric.Inc()

	// KillChan?  Japanese mascot for... uh...
	killChan := make(chan struct{})
	obj.registerShutdownChannel(killChan)

	// Per-worker context.
	wctx, cancel := context.WithCancel(obj.ctx)

	go obj.startWorkerLoop(wctx, cancel, killChan)
}

func (obj *workerPool) registerShutdownChannel(killChan chan struct{}) {
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	obj.shutdownChans = append(obj.shutdownChans, killChan)
}

func (obj *workerPool) unregisterShutdownChannel(killChan chan struct{}) {
	obj.shutdownLock.Lock()
	defer obj.shutdownLock.Unlock()

	for idx, kchan := range obj.shutdownChans {
		if kchan == killChan {
			obj.shutdownChans = append(
				obj.shutdownChans[:idx],
				obj.shutdownChans[idx+1:]...,
			)

			break
		}
	}
}

func (obj *workerPool) startWorkerLoop(wctx context.Context, cancel context.CancelFunc, killChan chan struct{}) {
	defer func() {
		cancel()
		obj.wg.Done()
		obj.workerCount.Add(-1)
		obj.activeWorkersMetric.Dec()
		obj.unregisterShutdownChannel(killChan)
	}()

	for {
		select {
		case <-obj.ctx.Done():
			return

		case <-killChan:
			cancel()

			return

		default:
			if obj.handleWorkerLifecycle(wctx) {
				obj.lgr.Info(
					"Worker timed out.",
					"type", "dynworker",
					"pool", obj.name,
				)

				return
			}
		}
	}
}

func (obj *workerPool) handleWorkerLifecycle(wctx context.Context) bool {
	tctx, tcancel := context.WithTimeout(wctx, obj.config.IdleTimeout)
	defer tcancel()

	task, err := obj.input.Get(tctx)
	if err != nil {
		// Channel closed, exit.
		if errors.Is(err, ErrChannelClosed) {
			return true
		}

		// If we timed out and we have more than min workers, exit.
		if errors.Is(errors.Unwrap(err), context.DeadlineExceeded) {
			return obj.workerCount.Load() > obj.minWorkers.Load()
		}

		// Treat cancellation or context shutdown as an exit.
		return true
	}

	if task == nil {
		return true
	}

	obj.processTask(task)

	return false
}

func (obj *workerPool) processTask(task *Task) {
	start := time.Now()

	_ = obj.processFn(task)
	obj.tasksTotalMetric.Inc()

	elapsed := time.Since(start)

	runtime.Gosched() // Give some time to Go.

	obj.updateAvgProcTime(elapsed.Nanoseconds())
	obj.taskDurationMetric.Observe(float64(elapsed) / float64(time.Second))

	task.reset()
	obj.taskPool.Put(task)
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
			obj.minWorkers.Load(),
			obj.maxWorkers.Load(),
		)
	}

	queued := obj.input.Len()

	avg := time.Duration(obj.avgProcTime.Load())
	if avg == 0 {
		avg = defaultAverageProcessTime
	}

	target := obj.config.DrainTarget
	if target <= 0 {
		target = time.Second
	}

	// workesr ~ ceil(queued * avg / target)
	num := int64(gomath.Ceil(float64(queued) *
		(float64(avg) / float64(target))))

	if num < 1 {
		num = 1
	}

	return math.ClampI64(num, obj.minWorkers.Load(), obj.maxWorkers.Load())
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
		"delta", num,
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
		"delta", num,
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

	if config.WorkerFunc == nil {
		panic("dynworker: WorkerFunc must not be nil.")
	}

	if config.Prometheus == nil {
		config.Prometheus = prometheus.DefaultRegisterer
	}

	lgr := logger.MustGetLogger(ctx)
	nctx, cancel := context.WithCancel(ctx)

	InitPrometheus(config.Prometheus)

	label := prometheus.Labels{"pool": config.Name}

	taskPool := &sync.Pool{
		New: func() any {
			return &Task{}
		},
	}

	obj := &workerPool{
		name:                  config.Name,
		input:                 config.InputQueue,
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

	obj.SetMinWorkers(config.MinWorkers)
	obj.SetMaxWorkers(config.MaxWorkers)

	return obj
}

// Initialise Prometheus metrics for this module.
func InitPrometheus(reg prometheus.Registerer) {
	prometheusInitOnce.Do(func() {
		reg.MustRegister(
			activeWorkers,
			tasksTotal,
			taskDuration,
			totalScaledUp,
			totalScaledDown,
		)
	})
}

// * dynworker.go ends here.
