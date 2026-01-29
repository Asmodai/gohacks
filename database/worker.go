// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// worker.go --- Database dynamic worker.
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

// * Comments:

// * Package:

package database

// * Imports:

import (
	"context"
	gomath "math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Asmodai/gohacks/dynworker"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/math"
	"github.com/Asmodai/gohacks/types"
	"gitlab.com/tozd/go/errors"
)

// * Code:
// ** TxnProvider:

// Any object that contains a `Txn` function can be used for callbacks.
type TxnProvider interface {
	Txn(context.Context, Runner) error
}

// ** WorkerJob:
// *** Interface:

// WorkerJob is a user-supplied unit of work which will be executed inside
// a database transaction.
//
// Implement `Run' with your SQL using the provided `Runner' (`*sqlx.DB` or
// `*sqlx.Tx`).
type WorkerJob interface {
	Run(ctx context.Context, runner Runner) error
}

// *** Type:

// Adapts a `WorkerJob` to a named `TxnFn` method.
type jobRunner struct {
	Job WorkerJob
}

// *** Methods:

// Execute the job's transaction code.
func (jr jobRunner) Txn(ctx context.Context, runner Runner) error {
	return errors.WithStack(jr.Job.Run(ctx, runner))
}

// ** BatchJob:
// *** Interface:

// BatchJob provides a means to invoke a user-supplied function with a batch
// of jobs.
type BatchJob interface {
	Run(ctx context.Context, runner Runner, data []dynworker.UserData) error
}

// *** Type:

// Adapts a `BatchJob` to a named `TxnFn` method.
type batchRunner struct {
	UserFn BatchJob
	Batch  []dynworker.UserData
	Logger logger.Logger
}

// *** Methods:

// Execute the batch's transaction code.
func (br batchRunner) Txn(ctx context.Context, runner Runner) error {
	return errors.WithStack(br.UserFn.Run(ctx, runner, br.Batch))
}

// ** workerTask:
// *** Type:

// The named `TaskFn` receiver used by the dynamic worker.
type workerTask struct {
	db      Database
	handler BatchJob
}

// *** Methods:

// Execute a job on the queue.
func (wt *workerTask) Work(task *dynworker.Task) error {
	ctx := task.Parent()

	lgr := task.Logger()
	if lgr == nil {
		// Put up with the risk of a panic here.
		//
		// If you get this far without setting up a logger, then
		// congratulations... you managed to somehow avoid everything
		// else (such as dynworker itself) from breaking violently.
		lgr = logger.MustGetLogger(ctx)
	}

	raw := task.Data()

	// Wrappers that provide a `Txn` method.
	if tprov, ok := raw.(TxnProvider); ok {
		return errors.WithStack(wt.db.WithTransaction(ctx, tprov.Txn))
	}

	// Bare batch.
	if batch, ok := raw.([]dynworker.UserData); ok && len(batch) > 0 {
		brn := batchRunner{
			UserFn: wt.handler,
			Batch:  batch,
			Logger: lgr,
		}

		return errors.WithStack(wt.db.WithTransaction(ctx, brn.Txn))
	}

	// Bare worker job.
	if job, ok := raw.(WorkerJob); ok && job != nil {
		jrn := jobRunner{Job: job}

		return errors.WithStack(wt.db.WithTransaction(ctx, jrn.Txn))
	}

	return nil
}

// ** Worker:
// *** Interface:

type Worker interface {
	Name() string
	Database() Database
	Start()
	Stop()
	SubmitBatch(dynworker.UserData) error
	SubmitJob(WorkerJob) error
}

// *** Type:

type worker struct {
	db           Database
	batchHandler BatchJob
	pool         dynworker.WorkerPool
	queue        dynworker.TaskQueue
	started      atomic.Bool
	inputCh      chan dynworker.UserData
	waitg        sync.WaitGroup
	cancel       context.CancelFunc
	ctx          context.Context
	minWorkers   int
	maxWorkers   int
	batchSize    int
	drainTarget  time.Duration
	batchTimeout time.Duration
}

// *** Methods:

func (w *worker) Name() string {
	return w.pool.Name()
}

func (w *worker) Database() Database {
	return w.db
}

func (w *worker) Start() {
	if !w.started.CompareAndSwap(false, true) {
		return
	}

	// Start pool.
	w.pool.Start()

	// Start batcher.
	w.waitg.Add(1)

	go w.batcher()
}

func (w *worker) Stop() {
	if !w.started.CompareAndSwap(true, false) {
		return
	}

	// Stop batcher.
	w.cancel()
	w.waitg.Wait()

	// Stop pool.
	w.pool.Stop()
}

// Enqueue a job for execution.
//
// Blocks if the queue is full until capacity frees.
func (w *worker) SubmitBatch(job dynworker.UserData) error {
	select {
	case w.inputCh <- job:
		return nil

	case <-w.ctx.Done():
		return errors.WithStack(w.ctx.Err())
	}
}

// Enqueue a prebuilt job.
func (w *worker) SubmitJob(job WorkerJob) error {
	return errors.WithStack(w.pool.Submit(jobRunner{Job: job}))
}

func (w *worker) batcher() {
	defer w.waitg.Done()

	buf := make([]dynworker.UserData, 0, w.batchSize)

	timer := time.NewTimer(w.batchTimeout)
	defer timer.Stop()

	flush := func() {
		if len(buf) == 0 {
			return
		}

		// Copy to avoid data races with reuse.
		batch := make([]dynworker.UserData, len(buf))
		copy(batch, buf)
		buf = buf[:0]

		_ = w.pool.Submit(batchRunner{
			UserFn: w.batchHandler,
			Batch:  batch})
	}

	for {
		select {
		case <-w.ctx.Done():
			flush()

			return

		case item := <-w.inputCh:
			buf = append(buf, item)

			if len(buf) >= w.batchSize {
				flush()

				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}

				timer.Reset(w.batchTimeout)

				continue
			}

		case <-timer.C:
			flush()
			timer.Reset(w.batchTimeout)
		}
	}
}

func (w *worker) scaler() int {
	if err := w.db.Ping(); err != nil {
		return w.minWorkers
	}

	qlen := w.queue.Len()
	avg := 100 * time.Millisecond //nolint:mnd
	want := int(gomath.Ceil(float64(qlen) *
		(float64(avg) / float64(w.drainTarget))))

	if want < w.minWorkers {
		want = w.minWorkers
	}

	if want > w.maxWorkers {
		want = w.maxWorkers
	}

	return want
}

// ** Functions:

func computeQueueCaps(cfg *Config) int {
	want := cfg.PoolMaxWorkers
	got := math.WithinPlatform(int64(want), defaultMaxWorkerCount)

	if got < 1 {
		got = int(defaultMaxWorkerCount)
	}

	return got
}

func computeInputCaps(cfg *Config) int {
	maxW := cfg.PoolMaxWorkers

	want := int64(maxW) * int64(cfg.BatchSize) * 2 //nolint:mnd
	least := int64(cfg.BatchSize) * 2              //nolint:mnd
	got := math.WithinPlatform(want, least)

	if got < 1 {
		got = 1
	}

	return got
}

func NewWorker(parent context.Context, cfg *Config, dbase Database, handler BatchJob) Worker {
	if !cfg.UsePool {
		return nil
	}

	icaps := computeInputCaps(cfg)
	inputCh := make(chan dynworker.UserData, icaps)

	qcaps := computeQueueCaps(cfg)
	queue := types.NewBoundedQueue(qcaps)
	dwqueue := dynworker.NewQueueTaskQueue(queue)

	ctx, cancel := context.WithCancel(parent)

	dwcfg := dynworker.NewConfigWithQueue(
		"database_"+cfg.Database,
		int64(cfg.PoolMinWorkers),
		int64(cfg.PoolMaxWorkers),
		dwqueue)

	dwcfg.IdleTimeout = cfg.PoolIdleTimeout.Duration()
	dwcfg.DrainTarget = cfg.PoolDrainTimeout.Duration()
	dwcfg.WorkerFunc = (&workerTask{db: dbase, handler: handler}).Work

	pool := dynworker.NewWorkerPool(ctx, dwcfg)

	inst := &worker{
		db:           dbase,
		ctx:          ctx,
		cancel:       cancel,
		inputCh:      inputCh,
		pool:         pool,
		queue:        dwqueue,
		minWorkers:   cfg.PoolMinWorkers,
		maxWorkers:   cfg.PoolMaxWorkers,
		batchSize:    cfg.BatchSize,
		batchTimeout: cfg.BatchTimeout.Duration(),
		drainTarget:  cfg.PoolDrainTimeout.Duration(),
		batchHandler: handler,
	}

	pool.SetScalerFunction(inst.scaler)

	return inst
}

// * worker.go ends here.
