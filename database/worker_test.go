// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// worker_test.go --- Worker tests.
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

// * Comments:

// * Package:

package database

// * Imports:

import (
	"context"
	"database/sql"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/dynworker"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/types"
	"github.com/jmoiron/sqlx"
)

// * Constants:

// * Variables:

// * Code:
// ** Fakes:
// *** fakeRunner:

type fakeRunner struct{}

func (f *fakeRunner) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return nil, nil
}

func (f *fakeRunner) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return nil, nil
}

func (f *fakeRunner) GetContext(ctx context.Context, dest any, query string, args ...any) error {
	return nil
}

func (f *fakeRunner) SelectContext(ctx context.Context, dest any, query string, args ...any) error {
	return nil
}

func (f *fakeRunner) QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	return &sqlx.Rows{}, nil
}

func (f *fakeRunner) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	return &sqlx.Row{}
}

func (f *fakeRunner) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return &sqlx.Stmt{}, nil
}

func (f *fakeRunner) BindNamed(query string, arg any) (string, []any, error) {
	// In tests we don’t care about actual binding—return the query unchanged
	// and no args so callers don’t explode.
	return query, nil, nil
}

func (f *fakeRunner) Rebind(query string) string {
	return query
}

func (f *fakeRunner) DriverName() string {
	return "fake"
}

// *** fakeDB:

type errorHolder struct{ err error }

type fakeDB struct {
	mu      sync.Mutex
	txCalls int
	perr    atomic.Pointer[errorHolder]
	driver  string
}

func newFakeDB() *fakeDB {
	db := &fakeDB{driver: "mysql"}
	return db
}

func (f *fakeDB) Ping() error {
	h := f.perr.Load()
	if h == nil {
		return nil
	}
	return h.err
}

// test code can flip health like this:
func (f *fakeDB) setPingErr(err error) {
	if err == nil {
		f.perr.Store(nil)
	} else {
		f.perr.Store(&errorHolder{err: err})
	}
}

func (f *fakeDB) Close() error {
	return nil
}

func (f *fakeDB) SetMaxIdleConns(int) {
}

func (f *fakeDB) SetMaxOpenConns(int) {
}

func (f *fakeDB) Rebind(q string) string {
	return q
}

func (f *fakeDB) Runner() Runner {
	return &fakeRunner{}
}

func (f *fakeDB) GetError(err error) error {
	return err
}

func (f *fakeDB) WithTransaction(ctx context.Context, fn TxnFn) error {
	f.mu.Lock()
	f.txCalls++
	f.mu.Unlock()
	return fn(ctx, f.Runner())
}
func (f *fakeDB) TxCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.txCalls
}

// *** fakeQueue:

type fakeQueue struct {
	lenVal atomic.Int64
}

func (q *fakeQueue) Put(ctx context.Context, t *dynworker.Task) error {
	return nil
}

func (q *fakeQueue) Get(ctx context.Context) (*dynworker.Task, error) {
	return nil, context.Canceled
}

func (q *fakeQueue) Len() int {
	return int(q.lenVal.Load())
}

// ** Helpers:

func waitFor(cond func() bool, timeout time.Duration, tick time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if cond() {
			return true
		}
		time.Sleep(tick)
	}
	return false
}

// ** Handlers:

type testWorkerJob struct {
	calls *atomic.Int64
}

func (j *testWorkerJob) Run(ctx context.Context, r Runner) error {
	j.calls.Add(1)
	_, _ = r.ExecContext(ctx, "INSERT DUMMY")
	return nil
}

type testBatchHandler struct {
	mu        sync.Mutex
	batchLens []int
	calls     int
}

func (h *testBatchHandler) Run(ctx context.Context, r Runner, data []dynworker.UserData) error {
	h.mu.Lock()
	h.calls++
	h.batchLens = append(h.batchLens, len(data))
	h.mu.Unlock()

	// Convert []dynworker.UserData -> []any for variadic ExecContext.
	args := make([]any, len(data))
	for i := range data {
		args[i] = data[i]
	}

	_, _ = r.ExecContext(ctx, "INSERT BATCH", args...)
	return nil
}

func (h *testBatchHandler) callCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.calls
}
func (h *testBatchHandler) lens() []int {
	h.mu.Lock()
	defer h.mu.Unlock()
	cp := make([]int, len(h.batchLens))
	copy(cp, h.batchLens)
	return cp
}

// ** Tests:

func TestWorkerSubmitJobRunsTransaction(t *testing.T) {
	bctx := context.Background()
	db := newFakeDB()
	handler := &testBatchHandler{}

	ctx, err := logger.SetLogger(bctx, logger.NewDefaultLogger())
	if err != nil {
		t.Fatalf("Logger: %s", err.Error())
	}

	cfg := &Config{
		UsePool:          true,
		Database:         "testdb",
		PoolMinWorkers:   1,
		PoolMaxWorkers:   2,
		PoolIdleTimeout:  types.Duration(time.Second),
		PoolDrainTimeout: types.Duration(500 * time.Millisecond),
		BatchSize:        10,
		BatchTimeout:     types.Duration(100 * time.Millisecond),
	}

	w := NewWorker(ctx, cfg, db, handler)
	if w == nil {
		t.Fatal("NewWorker returned nil")
	}
	w.Start()
	defer w.Stop()

	jobCalls := atomic.Int64{}
	job := &testWorkerJob{calls: &jobCalls}

	if err := w.SubmitJob(job); err != nil {
		t.Fatalf("SubmitJob error: %v", err)
	}

	ok := waitFor(func() bool { return db.TxCount() >= 1 && jobCalls.Load() >= 1 }, 2*time.Second, 10*time.Millisecond)
	if !ok {
		t.Fatalf("expected tx and job Run to be called; got tx=%d job=%d", db.TxCount(), jobCalls.Load())
	}
}

func TestWorkerBatchFlushBySize(t *testing.T) {
	bctx := context.Background()
	db := newFakeDB()
	handler := &testBatchHandler{}

	ctx, err := logger.SetLogger(bctx, logger.NewDefaultLogger())
	if err != nil {
		t.Fatalf("Logger: %s", err.Error())
	}

	cfg := &Config{
		UsePool:          true,
		Database:         "testdb",
		PoolMinWorkers:   1,
		PoolMaxWorkers:   2,
		PoolIdleTimeout:  types.Duration(time.Second),
		PoolDrainTimeout: types.Duration(500 * time.Millisecond),
		BatchSize:        3,
		BatchTimeout:     types.Duration(5 * time.Second),
	}

	w := NewWorker(ctx, cfg, db, handler).(*worker)
	w.Start()
	defer w.Stop()

	// Submit exactly BatchSize items
	for i := 0; i < cfg.BatchSize; i++ {
		if err := w.SubmitBatch(struct{ N int }{N: i}); err != nil {
			t.Fatalf("SubmitBatch error: %v", err)
		}
	}

	ok := waitFor(func() bool { return db.TxCount() >= 1 && handler.callCount() >= 1 }, 2*time.Second, 10*time.Millisecond)
	if !ok {
		t.Fatalf("batch not flushed by size; tx=%d calls=%d", db.TxCount(), handler.callCount())
	}

	lens := handler.lens()
	if len(lens) != 1 || lens[0] != cfg.BatchSize {
		t.Fatalf("expected single batch of %d, got %v", cfg.BatchSize, lens)
	}
}

func TestWorkerBatchFlushByTimeout(t *testing.T) {
	bctx := context.Background()
	db := newFakeDB()
	handler := &testBatchHandler{}

	ctx, err := logger.SetLogger(bctx, logger.NewDefaultLogger())
	if err != nil {
		t.Fatalf("Logger: %s", err.Error())
	}

	cfg := &Config{
		UsePool:          true,
		Database:         "testdb",
		PoolMinWorkers:   1,
		PoolMaxWorkers:   2,
		PoolIdleTimeout:  types.Duration(time.Second),
		PoolDrainTimeout: types.Duration(500 * time.Millisecond),
		BatchSize:        10,
		BatchTimeout:     types.Duration(100 * time.Millisecond),
	}

	w := NewWorker(ctx, cfg, db, handler).(*worker)
	w.Start()
	defer w.Stop()

	if err := w.SubmitBatch(struct{ S string }{"a"}); err != nil {
		t.Fatalf("SubmitBatch error: %v", err)
	}

	ok := waitFor(func() bool { return db.TxCount() >= 1 && handler.callCount() >= 1 }, 2*time.Second, 10*time.Millisecond)
	if !ok {
		t.Fatalf("batch not flushed by timeout; tx=%d calls=%d", db.TxCount(), handler.callCount())
	}

	lens := handler.lens()
	if len(lens) != 1 || lens[0] != 1 {
		t.Fatalf("expected single batch of 1 by timeout, got %v", lens)
	}
}

func TestWorkerScalerGatesOnDBHealth(t *testing.T) {
	bctx := context.Background()
	db := newFakeDB()
	handler := &testBatchHandler{}

	ctx, err := logger.SetLogger(bctx, logger.NewDefaultLogger())
	if err != nil {
		t.Fatalf("Logger: %s", err.Error())
	}

	cfg := &Config{
		UsePool:          true,
		Database:         "testdb",
		PoolMinWorkers:   2,
		PoolMaxWorkers:   8,
		PoolIdleTimeout:  types.Duration(time.Second),
		PoolDrainTimeout: types.Duration(500 * time.Millisecond),
		BatchSize:        5,
		BatchTimeout:     types.Duration(100 * time.Millisecond),
	}

	w := NewWorker(ctx, cfg, db, handler).(*worker)

	// Replace the real queue with a fake so we can force a length.
	fq := &fakeQueue{}
	w.queue = fq

	// Healthy DB, non-zero queue -> scaler should be > min
	fq.lenVal.Store(20)
	db.setPingErr(error(nil))
	if got := w.scaler(); got <= w.minWorkers || got > w.maxWorkers {
		t.Fatalf("healthy scaler out of range: got=%d min=%d max=%d", got, w.minWorkers, w.maxWorkers)
	}

	// Sick DB -> scaler should clamp to min
	db.setPingErr(context.DeadlineExceeded)
	fq.lenVal.Store(1_000)
	if got := w.scaler(); got != w.minWorkers {
		t.Fatalf("unhealthy scaler must clamp to min: got=%d want=%d", got, w.minWorkers)
	}
}

// * worker_test.go ends here.
