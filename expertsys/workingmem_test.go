// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// workingmem_test.go --- Working memory tests.
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

package expertsys

// * Imports:

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// * Code:

// ** Tests:

func TestWorkingMem(t *testing.T) {
	t.Run("Set and Version", func(t *testing.T) {
		wm := NewWorkingMemory()

		if got := wm.Version(); got != 0 {
			t.Fatalf("expected version 0, got %d", got)
		}

		if changed := wm.Set("a", 1); !changed {
			t.Fatal("expected change on first set")
		}

		if got := wm.Version(); got != 1 {
			t.Fatalf("expected version 1, got %d", got)
		}

		if changed := wm.Set("a", 1); changed {
			t.Fatal("expected no change on subsequent set")
		}

		if got := wm.Version(); got != 1 {
			t.Fatalf("expected version unchanged, got %d", got)
		}

		if changed := wm.Set("a", 2); !changed {
			t.Fatal("expected change when value differs")
		}

		if got := wm.Version(); got != 2 {
			t.Fatalf("expected version 2, got %d", got)
		}
	})

	t.Run("Time equality", func(t *testing.T) {
		wm := NewWorkingMemory()

		t1 := time.Now().UTC().Truncate(time.Nanosecond)
		t2 := t1

		if !wm.Set("t", t2) {
			t.Fatal("expected first set to change")
		}

		v := wm.Version()

		if wm.Set("t", t2) {
			t.Fatal("expected no change for equal time values")
		}

		if got := wm.Version(); got != v {
			t.Fatalf("expected version unchanged, got %d want %d",
				got,
				v)
		}
	})

	t.Run("Keys sorted", func(t *testing.T) {
		wm := NewWorkingMemory()

		wm.Set("z", 1)
		wm.Set("a", 2)
		wm.Set("m", 3)

		keys := wm.Keys()
		want := []string{"a", "m", "z"}

		if len(keys) != len(want) {
			t.Fatalf("unexpected keys:  %v != %v", keys, want)
		}

		for idx := range want {
			if keys[idx] != want[idx] {
				t.Fatalf("unexpected key: %v != %v",
					keys[idx],
					want)
			}
		}
	})

	t.Run("Concurrent access", func(t *testing.T) {
		wm := NewWorkingMemory()

		var wg sync.WaitGroup

		for g := 0; g < 8; g++ {
			wg.Add(1)

			go func(id int) {
				defer wg.Done()

				for i := 0; i < 1000; i++ {
					key := "k"
					wm.Set(key, i)

					_, _ = wm.Get(key)
					_ = wm.Keys()
				}
			}(g)
		}

		wg.Wait()
	})
}

// ** Benchmarks:

func BenchmarkWorkingMemory_SetHotKeyChangingInt(b *testing.B) {
	wm := NewWorkingMemory()
	wm.Set("x", 0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = wm.Set("x", i)
	}
}

func BenchmarkWorkingMemory_SetHotKeyNoopInt(b *testing.B) {
	wm := NewWorkingMemory()
	wm.Set("x", 123)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = wm.Set("x", 123) // Should be a no-op.
	}
}

func BenchmarkWorkingMemory_SetHotKeyNoopTime(b *testing.B) {
	wm := NewWorkingMemory()
	t0 := time.Now().UTC()
	wm.Set("t", t0)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = wm.Set("t", t0) // Equal() path.
	}
}

func BenchmarkWorkingMemory_SetManyKeys(b *testing.B) {
	wm := NewWorkingMemory()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "k" + strconv.Itoa(i) // Allocates, but that's ok.
		_ = wm.Set(key, i)
	}
}

func BenchmarkWorkingMemory_SetKeyPool(b *testing.B) {
	const pool = 1024

	keys := make([]string, pool)
	for i := range keys {
		keys[i] = "k" + strconv.Itoa(i)
	}

	wm := NewWorkingMemory()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = wm.Set(keys[i%pool], i)
	}
}

func BenchmarkWorkingMemory_MixedGetSetParallel(b *testing.B) {
	wm := NewWorkingMemory()
	wm.Set("x", 0)

	var counter atomic.Int64

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v := counter.Add(1)
			_ = wm.Set("x", int(v))
			_, _ = wm.Get("x")
		}
	})
}

// * workingmem_test.go ends here.
