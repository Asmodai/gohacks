// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// priority_test.go --- Priority scheduler tests.
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
	"sync"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/health"
	"github.com/Asmodai/gohacks/logger"
)

// * Code:

// ** Helpers:

// Receive exactly n jobs from `ch` before `deadline`, otherwise fail.
//
// Returns the received jobs in order.
func recvN[T any](t *testing.T, ch <-chan T, n int, deadline time.Duration) []T {
	t.Helper()

	out := make([]T, 0, n)
	timer := time.NewTimer(deadline)
	defer timer.Stop()

	for len(out) < n {
		select {
		case v := <-ch:
			out = append(out, v)

		case <-timer.C:
			t.Fatalf("timeout waiting for %d items; got %d",
				n,
				len(out))
		}
	}

	return out
}

// Ensure that no jobs arrive within the given duration.
func expectNone[T any](t *testing.T, ch <-chan T, dur time.Duration) {
	t.Helper()

	timer := time.NewTimer(dur)
	defer timer.Stop()

	select {
	case v := <-ch:
		t.Fatalf("unexpected item received: %#v", v)

	case <-timer.C:
		// ok
	}
}

// ** Types:

// *** Fake ticker:

type fakeTicker struct {
	ch chan time.Time
}

func (t *fakeTicker) Channel() <-chan time.Time { return t.ch }
func (t *fakeTicker) Stop()                     { close(t.ch) }

func (t *fakeTicker) Tick() {
	t.ch <- time.Now()
}

func newFakeTicker() *fakeTicker {
	return &fakeTicker{ch: make(chan time.Time, 16)}
}

// *** Test job:

type testTimedJob struct {
	id    string
	runAt time.Time
}

func (t *testTimedJob) RunAt() time.Time { return t.runAt }

// Satisfy `Job`.
//
// We don't use these in the scheduler loop, but need them for interface
// conformance.
func (t *testTimedJob) Validate() error               { return nil }
func (t *testTimedJob) Resolve(context.Context) error { return nil }
func (t *testTimedJob) Object() Task                  { return nil }
func (t *testTimedJob) Function() JobFn               { return nil }

// Compile-time assertion.
var _ TimedJob = (*testTimedJob)(nil)

// ** Tests:

func TestSchedulerRun_BasicOrdering(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, _ = logger.SetLogger(ctx, logger.NewDefaultLogger())
	config := NewDefaultConfig()
	sched := NewPriority(ctx, config)
	now := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.Start()
	}()

	// Three jobs.
	j1 := &testTimedJob{id: "a", runAt: now.Add(30 * time.Millisecond)}
	j2 := &testTimedJob{id: "b", runAt: now.Add(60 * time.Millisecond)}
	j3 := &testTimedJob{id: "c", runAt: now.Add(45 * time.Millisecond)}

	// Add jobs out of order.
	sched.addCh <- j2
	sched.addCh <- j1
	sched.addCh <- j3

	// Receive jobs.
	got := recvN(t, sched.workCh, 3, 2*time.Second)

	// IDs of received jobs.
	ids := []string{
		got[0].(*testTimedJob).id,
		got[1].(*testTimedJob).id,
		got[2].(*testTimedJob).id,
	}

	// Order we want.
	want := []string{"a", "c", "b"}

	// Check if we have the right order.
	for idx := range want {
		if ids[idx] != want[idx] {
			t.Fatalf("order mismatch: got %v want %v", ids, want)
		}
	}

	cancel()
	wg.Wait()
}

func TestSchedulerRun_FIFOForSameTime(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, _ = logger.SetLogger(ctx, logger.NewDefaultLogger())
	config := NewDefaultConfig()
	sched := NewPriority(ctx, config)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.run()
	}()

	runAt := time.Now().Add(50 * time.Millisecond)

	// Same time; FIFO means insertion order should be preserved.
	j1 := &testTimedJob{id: "a", runAt: runAt}
	j2 := &testTimedJob{id: "b", runAt: runAt}
	j3 := &testTimedJob{id: "c", runAt: runAt}

	// Add jobs.
	sched.addCh <- j1
	sched.addCh <- j2
	sched.addCh <- j3

	// Receive jobs.
	got := recvN(t, sched.workCh, 3, 2*time.Second)

	// IDs of received jobs.
	ids := []string{
		got[0].(*testTimedJob).id,
		got[1].(*testTimedJob).id,
		got[2].(*testTimedJob).id,
	}

	// Order we want.
	want := []string{"a", "b", "c"}

	// Check if we have the right order.
	for idx := range want {
		if ids[idx] != want[idx] {
			t.Fatalf("order mismatch: got %v want %v", ids, want)
		}
	}

	cancel()
	wg.Wait()
}

func TestSchedulerRun_InsertEarlierJobResetsTimer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, _ = logger.SetLogger(ctx, logger.NewDefaultLogger())
	config := NewDefaultConfig()
	sched := NewPriority(ctx, config)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.run()
	}()

	now := time.Now()

	// Add late job first.
	late := &testTimedJob{id: "late", runAt: now.Add(250 * time.Millisecond)}
	sched.addCh <- late

	// Wait a bit so the scheduler has likely set a timer for "late",
	// then inject an earlier job and ensure it fires first.
	time.Sleep(20 * time.Millisecond)

	// Create earlier job.
	early := &testTimedJob{id: "early", runAt: now.Add(60 * time.Millisecond)}
	sched.addCh <- early

	// Receive jobs.
	got := recvN(t, sched.workCh, 2, 2*time.Second)

	// IDs of received jobs.
	ids := []string{
		got[0].(*testTimedJob).id,
		got[1].(*testTimedJob).id,
	}

	// Order we want.
	want := []string{"early", "late"}

	// Check if we have the right order.
	for idx := range want {
		if ids[idx] != want[idx] {
			t.Fatalf("order mismatch: got %v want %v", ids, want)
		}
	}

	cancel()
	wg.Wait()
}

func TestSchedulerRun_DoesNotEmitBeforeDue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, _ = logger.SetLogger(ctx, logger.NewDefaultLogger())
	config := NewDefaultConfig()
	sched := NewPriority(ctx, config)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.run()
	}()

	j := &testTimedJob{id: "x", runAt: time.Now().Add(200 * time.Millisecond)}
	sched.addCh <- j

	// Should not fire immediately.
	expectNone(t, sched.workCh, 50*time.Millisecond)

	// Receive jobs.
	got := recvN(t, sched.workCh, 1, 2*time.Second)
	if got[0].(*testTimedJob).id != "x" {
		t.Fatalf("unexpected job: %#v", got[0])
	}

	cancel()
	wg.Wait()
}

func TestSchedulerRun_CancelStops(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ctx, _ = logger.SetLogger(ctx, logger.NewDefaultLogger())
	config := NewDefaultConfig()
	sched := NewPriority(ctx, config)

	done := make(chan struct{})
	go func() {
		sched.run()
		close(done)
	}()

	cancel()

	select {
	case <-done:
	// ok

	case <-time.After(2 * time.Second):
		t.Fatal("scheduler did not stop after context cancellation")
	}
}

func TestSchedulerRun_ManyJobs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, _ = logger.SetLogger(ctx, logger.NewDefaultLogger())
	config := NewDefaultConfig()
	sched := NewPriority(ctx, config)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sched.run()
	}()

	base := time.Now().Add(50 * time.Millisecond)
	const n = 1024

	for idx := 0; idx < n; idx++ {
		runAt := base.Add(time.Duration((n-idx)%25) * 10 * time.Millisecond)
		sched.addCh <- &testTimedJob{
			id:    fmt.Sprintf("job-%03d", idx),
			runAt: runAt,
		}
	}

	// Receive jobs.
	got := recvN(t, sched.workCh, n, 2*time.Second)

	// Check if we have the right order.
	for idx := 1; idx < len(got); idx++ {
		prev := got[idx-1].RunAt()
		curr := got[idx].RunAt()

		if curr.Before(prev) {
			t.Fatalf("jobs not ordered by time: idx %d has %v before %v",
				idx,
				curr,
				prev)
		}
	}

	cancel()
	wg.Wait()
}

// ** Benchmarks:

func BenchmarkInsertTimedJob(b *testing.B) {
	jobs := make([]TimedJob, 0, b.N)
	base := time.Now()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		j := &testTimedJob{
			id:    "bench",
			runAt: base.Add(time.Duration(i) * time.Millisecond),
		}
		jobs, _ = InsertTimedJob(jobs, j)
	}
}

func TestHealth_TicksWhenIdle(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx, _ = logger.SetLogger(ctx, logger.NewDefaultLogger())
	config := NewDefaultConfig()
	config.Health = health.NewHealthWithDuration(150 * time.Millisecond)
	sched := NewPriority(ctx, config)

	done := make(chan struct{})
	go func() {
		sched.Start()
		close(done)
	}()

	start := sched.Health().LastHeartbeat()

	time.Sleep(1 * time.Second)

	after := sched.Health().LastHeartbeat()

	if !after.After(start) {
		t.Fatalf("expected heartbeat to advance; start=%v after=%v",
			start,
			after)
	}

	if sched.Health().Healthy() == false {
		t.Fatalf("expected scheduler to be healthy after a tick; start=%v after=%v",
			start,
			after)
	}

	cancel()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("scheduler did not stop after cancellation")
	}
}

func TestHealth_UnhealthyAfterStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ctx, _ = logger.SetLogger(ctx, logger.NewDefaultLogger())
	config := NewDefaultConfig()
	config.Health = health.NewHealthWithDuration(250 * time.Millisecond)
	sched := NewPriority(ctx, config)

	done := make(chan struct{})
	go func() {
		sched.Start()
		close(done)
	}()

	// Give the scheduler a moment to start.
	time.Sleep(50 * time.Millisecond)

	// Immediately kill the scheduler.
	cancel()

	select {
	case <-done:

	case <-time.After(2 * time.Second):
		t.Fatal("scheduler did not stop after cancellation")
	}

	time.Sleep(350 * time.Millisecond)

	if sched.Health().Healthy() == true {
		t.Fatal("expected scheduler health to be unhealthy")
	}
}

// * priority_test.go ends here.
