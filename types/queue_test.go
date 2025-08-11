// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// queue_test.go --- Queue tests.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

package types

// * Imports:

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

// * Code:

// ** Benchmarks:

var benchDatum Datum = struct{}{} // adjust if Datum isn't `any`

// Serial put->get round trip in one goroutine (no contention).
func BenchmarkQueue_Serial_PutGet(b *testing.B) {
	q := NewBoundedQueue(1024)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		q.Put(benchDatum)
		_ = q.Get()
	}
}

// One producer, one consumer, bounded queue to exercise cond wakeups.
func BenchmarkQueue_1P1C_Bounded(b *testing.B) {
	q := NewBoundedQueue(256)
	var wg sync.WaitGroup
	wg.Add(1)

	b.ReportAllocs()
	b.ResetTimer()

	go func(n int) {
		defer wg.Done()
		for i := 0; i < n; i++ {
			_ = q.Get()
		}
	}(b.N)

	for i := 0; i < b.N; i++ {
		q.Put(benchDatum)
	}

	wg.Wait()
}

// Many producers/consumers with context-cancellable consumers.
// Stresses lock contention + cond.Broadcast paths.
func BenchmarkQueue_ManyP_ManyC_Bounded(b *testing.B) {
	type cfg struct{ P, C, Cap int }
	cases := []cfg{
		{P: 1, C: 1, Cap: 64},
		{P: runtime.GOMAXPROCS(0), C: runtime.GOMAXPROCS(0), Cap: 128},
		{P: 8, C: 8, Cap: 32},
	}

	for _, tc := range cases {
		name := func(c cfg) string {
			return b.Name() + "_P" + itoa(c.P) + "_C" + itoa(c.C) + "_Cap" + itoa(c.Cap)
		}(tc)

		b.Run(name, func(b *testing.B) {
			q := NewBoundedQueue(tc.Cap)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			var got int64
			want := int64(b.N)

			var wg sync.WaitGroup
			wg.Add(tc.C)
			for c := 0; c < tc.C; c++ {
				go func() {
					defer wg.Done()
					for {
						v, err := q.GetWithContext(ctx)
						if err != nil {
							return // ctx cancelled
						}
						_ = v
						if atomic.AddInt64(&got, 1) >= want {
							// We’ve consumed everything; other consumers will exit on cancel.
						}
					}
				}()
			}

			b.ReportAllocs()
			b.ResetTimer()

			// Producers share the work: b.N total items.
			var pg sync.WaitGroup
			pg.Add(tc.P)
			per := b.N / tc.P
			rem := b.N % tc.P
			for p := 0; p < tc.P; p++ {
				n := per
				if p == 0 {
					n += rem // account for remainder
				}
				go func(n int) {
					defer pg.Done()
					for i := 0; i < n; i++ {
						q.Put(benchDatum)
					}
				}(n)
			}

			pg.Wait()

			// Wait until all items are consumed, then cancel consumers cleanly.
			for atomic.LoadInt64(&got) < want {
				runtime.Gosched()
			}
			cancel()
			wg.Wait()
		})
	}
}

// Contended hot-path with goroutines doing Put+Get pairs.
// Good for measuring lock/cond overhead under mixed ops.
func BenchmarkQueue_RunParallel_Pairs(b *testing.B) {
	q := NewBoundedQueue(256)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Put(benchDatum)
			_ = q.Get()
		}
	})
}

// Microbench the cancellation path: empty queue, already-cancelled ctx.
func BenchmarkQueue_GetWithContext_Canceled(b *testing.B) {
	q := NewBoundedQueue(0)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already canceled

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = q.GetWithContext(ctx)
	}
}

// Helper: tiny int -> string without fmt allocation in inner benches.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	// small, allocation-free int to ascii
	var a [20]byte
	i := len(a)
	neg := n < 0
	u := uint64(n)
	if neg {
		u = uint64(-n)
	}
	for u > 0 {
		i--
		a[i] = byte('0' + u%10)
		u /= 10
	}
	if neg {
		i--
		a[i] = '-'
	}
	return string(a[i:])
}

// * queue_test.go ends here.
