// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// mailbox_test.go --- Mailbox tests.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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
	"time"
)

// * Code:

// ** Tests:

func TestMailboxNoContext(t *testing.T) {
	var mbox *Mailbox = nil

	//
	// Test creation.
	t.Run("Can create new mailbox", func(t *testing.T) {
		mbox = NewMailbox()

		if mbox.Full() {
			t.Error("Somehow the mailbox is already full!")
			return
		}
	})

	//
	// Test single write/read.
	t.Run("Can write and read in a single routine", func(t *testing.T) {
		mbox.Put("test")

		val, ok := mbox.Get()
		if !ok {
			t.Error("Get failed")
			return
		}

		if val.(string) != "test" {
			t.Errorf("Unexpected value '%v' returned.", val)
			return
		}
	})

	//
	// Test with goroutines.
	t.Run("Can write and read with concurrency", func(t *testing.T) {
		writeChan := make(chan *Pair, 1)
		readChan := make(chan *Pair, 1)

		go func(mb *Mailbox) {
			mb.Put("concurrency")
			writeChan <- NewPair(true, nil)
		}(mbox)

		go func(mb *Mailbox) {
			val, ok := mb.Get()
			for ok == false {
				val, ok = mb.Get()
				time.Sleep(50 * time.Millisecond)
			}

			readChan <- NewPair(val, ok)
		}(mbox)

		<-writeChan

		select {
		case readRes := <-readChan:
			{
				if readRes.Second.(bool) != true {
					t.Error("Problem with reading value!")
					return
				}

				if readRes.First.(string) != "concurrency" {
					t.Errorf("Unexpected value '%v' returned.", readRes.First)
					return
				}
			}

		case <-time.After(5 * time.Second):
			t.Error("Timeout after 5 seconds!")
			return
		}
	})
}

// ** Benchmarks:

var benchMsg Datum = struct{}{} // adjust if Datum isn't `any`

// Single goroutine round-trip (no contention).
func BenchmarkMailbox_Serial_PutGet(b *testing.B) {
	m := NewMailbox()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m.Put(benchMsg)
		_, _ = m.Get()
	}
}

// 1 producer / 1 consumer, steady flow.
func BenchmarkMailbox_1P1C(b *testing.B) {
	m := NewMailbox()
	var wg sync.WaitGroup
	wg.Add(1)

	b.ReportAllocs()
	b.ResetTimer()

	go func(n int) {
		defer wg.Done()
		for i := 0; i < n; i++ {
			_, _ = m.Get()
		}
	}(b.N)

	for i := 0; i < b.N; i++ {
		_ = m.Put(benchMsg)
	}

	wg.Wait()
}

// Many producers & consumers, exercises wakeups and contention.
func BenchmarkMailbox_ManyP_ManyC(b *testing.B) {
	type cfg struct{ P, C int }
	cases := []cfg{
		{P: 1, C: 1},
		{P: runtime.GOMAXPROCS(0), C: runtime.GOMAXPROCS(0)},
		{P: 8, C: 8},
	}

	for _, tc := range cases {
		name := benchName("P", tc.P, "C", tc.C)
		b.Run(name, func(b *testing.B) {
			m := NewMailbox()
			var got int64
			want := int64(b.N)

			// Consumers
			var cg sync.WaitGroup
			cg.Add(tc.C)
			ctx, cancel := context.WithCancel(context.Background())
			for i := 0; i < tc.C; i++ {
				go func() {
					defer cg.Done()
					for {
						v, ok := m.GetWithContext(ctx)
						if !ok {
							return
						}
						_ = v
						if atomic.AddInt64(&got, 1) >= want {
							// keep draining until cancelled
						}
					}
				}()
			}

			// Producers share the load
			var pg sync.WaitGroup
			pg.Add(tc.P)
			per := b.N / tc.P
			rem := b.N % tc.P

			b.ReportAllocs()
			b.ResetTimer()

			for p := 0; p < tc.P; p++ {
				n := per
				if p == 0 {
					n += rem
				}
				go func(n int) {
					defer pg.Done()
					for i := 0; i < n; i++ {
						_ = m.Put(benchMsg)
					}
				}(n)
			}

			pg.Wait()

			// Wait until all observed
			for atomic.LoadInt64(&got) < want {
				runtime.Gosched()
			}
			cancel()
			cg.Wait()
		})
	}
}

// Hot path: each goroutine does Put then Get (balanced).
func BenchmarkMailbox_RunParallel_Pairs(b *testing.B) {
	m := NewMailbox()
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = m.Put(benchMsg)
			_, _ = m.Get()
		}
	})
}

// Non-blocking try paths under churn.
func BenchmarkMailbox_TryPutTryGet(b *testing.B) {
	m := NewMailbox()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !m.TryPut(benchMsg) {
			_, _ = m.TryGet()
		}
	}
}

// Context-cancel microbenches (fast-fail).
func BenchmarkMailbox_GetWithContext_Canceled(b *testing.B) {
	m := NewMailbox()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.GetWithContext(ctx)
	}
}

func BenchmarkMailbox_PutWithContext_Canceled(b *testing.B) {
	m := NewMailbox()
	// Fill the mailbox so Put will need to wait.
	_ = m.Put(benchMsg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m.PutWithContext(ctx, benchMsg)
	}
}

// tiny, alloc-free name helper
func benchName(k1 string, v1 int, k2 string, v2 int) string {
	var buf [64]byte
	n := 0
	n += copy(buf[n:], k1)
	n += itoaInto(&buf, n, v1)
	n += copy(buf[n:], "_")
	n += copy(buf[n:], k2)
	n += itoaInto(&buf, n, v2)
	return string(buf[:n])
}

func itoaInto(b *[64]byte, off int, n int) int {
	if n == 0 {
		b[off] = '0'
		return 1
	}
	start := off
	if n < 0 {
		b[off] = '-'
		off++
		n = -n
	}
	// write digits into temp
	var tmp [20]byte
	i := len(tmp)
	for n > 0 {
		i--
		tmp[i] = byte('0' + (n % 10))
		n /= 10
	}
	c := copy(b[off:], tmp[i:])
	return (off - start) + c
}

// * mailbox_test.go ends here.
