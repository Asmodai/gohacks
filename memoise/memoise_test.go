// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// memoise_test.go --- Memoisation tests.
//
// Copyright (c) 2024-2025 Paul Ward <paul@lisphacker.uk>
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

//
//
//

// * Package:

package memoise

// * Imports:

import (
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	testError error = errors.Base("this is an error")
)

// * Code:

// ** Tests:

func TestMemoise(t *testing.T) {
	t.Run("Works as expected", func(t *testing.T) {
		memo := NewMemoise(NewDefaultConfig())

		r1, _ := memo.Check(
			"Test1",
			func() (any, error) { return 42, nil },
		)

		// r2 is sufficiently different from r1 so we can detect
		// whether there's a miss or not.
		r2, _ := memo.Check(
			"Test1",
			func() (any, error) { return 84, nil },
		)

		if r1 != r2 {
			t.Errorf("Result mismatch: %v != %v", r1, r2)
		}
	})

	t.Run("Handles errors", func(t *testing.T) {
		memo := NewMemoise(NewDefaultConfig())

		_, err := memo.Check(
			"Test2",
			func() (any, error) { return nil, testError },
		)

		if !errors.Is(err, testError) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("Errors are not memoised", func(t *testing.T) {
		memo := NewMemoise(NewDefaultConfig())
		callCount := 0

		callback := func() (any, error) {
			callCount++
			return nil, testError
		}

		for i := 0; i < 3; i++ {
			_, _ = memo.Check("faily", callback)
		}

		if callCount != 3 {
			t.Errorf(
				"Expected callback to run 3 times, ran %d times",
				callCount,
			)
		}
	})

	t.Run("Resets", func(t *testing.T) {
		memo := NewMemoise(NewDefaultConfig())

		r1, _ := memo.Check(
			"Test1",
			func() (any, error) { return 100, nil },
		)

		// `Reset` should nuke the existing `Test` result.
		memo.Reset()

		r2, _ := memo.Check(
			"Test1",
			func() (any, error) { return 200, nil },
		)

		if r1 == r2 {
			t.Errorf("Result mismatch: %#v != %#v", r1, r2)
		}
	})
}

func TestMemoise_Concurrency(t *testing.T) {
	memo := NewMemoise(NewDefaultConfig())

	const goroutines = 100
	var wg sync.WaitGroup
	wg.Add(goroutines)

	counter := 0
	callback := func() (any, error) {
		time.Sleep(10 * time.Millisecond) // Simulate real work
		counter++
		return "hello", nil
	}

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			_, err := memo.Check("key", callback)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	if counter != 1 {
		t.Errorf("Expected callback to be called once, got %d", counter)
	}
}

// ** Benchmarks:

// a tiny value to avoid heap noise
var benchVal any = struct{}{}

func BenchmarkMemo_Hit_Serial(b *testing.B) {
	m := NewMemoise(NewDefaultConfig())
	// seed once
	_, err := m.Check("k", func() (any, error) { return benchVal, nil })
	if err != nil {
		b.Fatalf("seed: %v", err)
	}

	var calls int64
	cb := func() (any, error) {
		atomic.AddInt64(&calls, 1)
		return benchVal, nil
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		v, err := m.Check("k", cb)
		if err != nil {
			b.Fatal(err)
		}
		if v == nil {
			b.Fatal("nil")
		}
	}

	b.StopTimer()
	if atomic.LoadInt64(&calls) != 0 {
		b.Fatalf("callback should not be called on hits; got %d", calls)
	}
}

func BenchmarkMemo_Miss_Serial(b *testing.B) {
	m := NewMemoise(NewDefaultConfig())
	var ctr int64
	cb := func() (any, error) {
		// pretend to compute
		return atomic.AddInt64(&ctr, 1), nil
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		v, err := m.Check(key, cb)
		if err != nil {
			b.Fatal(err)
		}
		if v == nil {
			b.Fatal("nil")
		}
	}
}

func BenchmarkMemo_Hit_RunParallel(b *testing.B) {
	m := NewMemoise(NewDefaultConfig())
	// seed once
	_, err := m.Check("hot", func() (any, error) { return benchVal, nil })
	if err != nil {
		b.Fatalf("seed: %v", err)
	}

	var calls int64
	cb := func() (any, error) {
		atomic.AddInt64(&calls, 1)
		return benchVal, nil
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v, err := m.Check("hot", cb)
			if err != nil {
				b.Fatal(err)
			}
			_ = v
		}
	})

	b.StopTimer()
	if atomic.LoadInt64(&calls) != 0 {
		b.Fatalf("callback should not run on parallel hits; got %d", calls)
	}
}

// Contention singleflight: P goroutines hit the same cold key, callback runs once.
func BenchmarkMemo_SingleFlight_Contention(b *testing.B) {
	P := runtime.GOMAXPROCS(0)

	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m := NewMemoise(NewDefaultConfig())

		start := make(chan struct{})
		var calls int64

		cb := func() (any, error) {
			atomic.AddInt64(&calls, 1)
			<-start // all goroutines pile up here; release once
			return benchVal, nil
		}

		var wg sync.WaitGroup
		wg.Add(P)
		for g := 0; g < P; g++ {
			go func() {
				defer wg.Done()
				v, err := m.Check("key", cb)
				if err != nil {
					b.Error(err)
				}
				_ = v
			}()
		}

		// Measure only the critical region
		b.StopTimer()
		close(start) // let exactly one compute run, others wait on ready
		b.StartTimer()

		wg.Wait()

		b.StopTimer()
		if atomic.LoadInt64(&calls) != 1 {
			b.Fatalf("singleflight broken: callback ran %d times", calls)
		}
		b.StartTimer()
	}
}

// Error path: ensure errors are NOT cached and the callback runs each time.
func BenchmarkMemo_Error_NotCached(b *testing.B) {
	m := NewMemoise(NewDefaultConfig())
	var calls int64
	cb := func() (any, error) {
		atomic.AddInt64(&calls, 1)
		return nil, assertErr("boom")
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		_, err := m.Check("errkey", cb)
		if err == nil {
			b.Fatal("expected error")
		}
	}

	b.StopTimer()
	if atomic.LoadInt64(&calls) != int64(b.N) {
		b.Fatalf("error should re-run each time; calls=%d, N=%d", calls, b.N)
	}
}

// Distinct keys with parallel churn: stresses map growth + write lock path.
func BenchmarkMemo_DistinctKeys_RunParallel(b *testing.B) {
	m := NewMemoise(NewDefaultConfig())
	var ctr int64
	cb := func() (any, error) {
		return atomic.AddInt64(&ctr, 1), nil
	}

	b.ReportAllocs()
	b.ResetTimer()

	var keyCtr uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			i := atomic.AddUint64(&keyCtr, 1)
			_, err := m.Check(strconv.FormatUint(i, 10), cb)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// tiny error type to avoid fmt allocs
type assertErr string

func (e assertErr) Error() string { return string(e) }

// Optional: callback that does "real" work w/o sleep, if you want a heavier miss bench.
func spin(ms int) {
	deadline := time.Now().Add(time.Duration(ms) * time.Millisecond)
	for time.Now().Before(deadline) {
		// busy work: nothing
	}
}

// * memoise_test.go ends here.
