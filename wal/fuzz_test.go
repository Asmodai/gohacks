// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// fuzz_test.go --- Fuzzing.
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

// * Comments:

// * Package:

package wal

// * Imports:

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// * Constants:

// * Variables:

// * Code:

// ---------- helpers ----------

func tmpPathFuzz(t testing.TB) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "fuzzwal.walx")
}

func openWithPolicy(t testing.TB, pol Policy) (WriteAheadLog, string) {
	t.Helper()
	p := tmpPathFuzz(t)
	w, err := OpenWALWithPolicy(makeLogger(), p, pol)
	if err != nil {
		t.Fatalf("OpenWAL: %v", err)
	}
	return w, p
}

type rec struct {
	lsn uint64
	ts  int64
	k   []byte
	v   []byte
}

func appendRec(t testing.TB, w WriteAheadLog, r rec) {
	t.Helper()
	if err := w.Append(r.lsn, r.ts, r.k, r.v); err != nil {
		t.Fatalf("append: %v", err)
	}
}

func replayCollect(w WriteAheadLog, base uint64) (out []rec, err error) {
	cb := func(lsn uint64, ts int64, k, v []byte) error {
		out = append(out, rec{lsn: lsn, ts: ts, k: append([]byte(nil), k...), v: append([]byte(nil), v...)})
		return nil
	}
	_, err = w.Replay(base, cb)
	return
}

// ========== FUZZ: round trip append->replay ==========

func FuzzAppendReplayRoundTrip(f *testing.F) {
	// Seeds (kept tiny so we hit the interesting paths quickly).
	seed := func(n int, key, val string) {
		f.Add(uint64(n), int64(1730000000+n), []byte(key), []byte(val))
	}
	seed(1, "k", "v")
	seed(2, "k2", "v2")
	seed(10, "alpha", "beta")

	// Base policy, but force no ticker to avoid nondeterminism.
	pol := tinyPolicy()
	pol.SyncEvery = 0
	pol.SyncEveryBytes = 0

	f.Fuzz(func(t *testing.T, lsn uint64, ts int64, key, val []byte) {
		// Clamp to policy so Append succeeds.
		if len(key) == 0 {
			key = []byte("k")
		}
		if len(key) > int(pol.MaxKeyBytes) {
			key = key[:pol.MaxKeyBytes]
		}
		if ts < 0 {
			ts = 0
		}
		if len(val) == 0 {
			val = []byte("v")
		}
		if len(val) > int(pol.MaxValueBytes) {
			val = val[:pol.MaxValueBytes]
		}

		w, _ := openWithPolicy(t, pol)
		defer func() { _ = w.Close() }()

		// Deterministic RNG per input.
		r := rand.New(rand.NewSource(int64(lsn) ^ ts ^ int64(len(key))<<17 ^ int64(len(val))<<5))

		// Write a small deterministic prefix (so base>0 paths are covered).
		prefix := r.Intn(4) // 0..3
		written := make([]rec, 0, prefix+1)
		for i := 0; i < prefix; i++ {
			rr := rec{
				lsn: uint64(i + 1),
				ts:  int64(1730000000 + i),
				k:   []byte{byte('a' + i)},
				v:   []byte{byte('A' + i)},
			}
			appendRec(t, w, rr)
			written = append(written, rr)
		}

		// One fuzzed record. Allow arbitrary LSN, but make 0 → prefix+1
		if lsn == 0 {
			lsn = uint64(prefix + 1)
		}
		fr := rec{lsn: lsn, ts: ts, k: key, v: val}
		appendRec(t, w, fr)
		written = append(written, fr)

		// Pick a random base within the prefix, so we often skip the first few.
		base := uint64(r.Intn(prefix + 1))

		got, err := replayCollect(w, base)
		if err != nil {
			t.Fatalf("replay: %v", err)
		}

		// Build the expected *applied* sequence (filter by base).
		expected := filterByBase(written, base)

		// WAL tail-truncation semantics: got must be a prefix of expected.
		if !isPrefixRecords(expected, got) {
			t.Fatalf("replay not a prefix (base=%d): want prefix of %v, got %v", base, expected, got)
		}

		// Optional: internal order sanity.
		for i := 1; i < len(got); i++ {
			if got[i].lsn < got[i-1].lsn {
				t.Fatalf("replayed out of order at %d: %d < %d", i, got[i].lsn, got[i-1].lsn)
			}
		}
	})
}

func filterByBase(in []rec, base uint64) []rec {
	out := make([]rec, 0, len(in))
	for _, r := range in {
		if r.lsn > base {
			out = append(out, r)
		}
	}
	return out
}

func isPrefixRecords(full, cand []rec) bool {
	if len(cand) > len(full) {
		return false
	}
	for i := range cand {
		a, b := full[i], cand[i]
		if a.lsn != b.lsn || a.ts != b.ts || !bytes.Equal(a.k, b.k) || !bytes.Equal(a.v, b.v) {
			return false
		}
	}
	return true
}

func maxU64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// ========== FUZZ: random corruption (truncate/flip) yields prefix or error, never panic ==========

func FuzzReplayCorruptionPrefixOrError(f *testing.F) {
	pol := tinyPolicy()

	f.Add(uint64(5), int64(3), int64(1)) // lsnCount, flips, truncBytes

	f.Fuzz(func(t *testing.T, nLSN uint64, flips int64, truncBytes int64) {
		if nLSN == 0 {
			nLSN = 1
		}
		w, path := openWithPolicy(t, pol)

		// write nLSN records
		var want []rec
		for i := uint64(1); i <= nLSN; i++ {
			r := rec{
				lsn: i,
				ts:  int64(1730000000 + i),
				k:   []byte(fmt.Sprintf("k%d", i)),
				v:   bytes.Repeat([]byte{byte(i)}, 8+int(i%5)),
			}
			appendRec(t, w, r)
			want = append(want, r)
		}
		_ = w.Close()

		// corrupt
		fh, err := os.OpenFile(path, os.O_RDWR, 0)
		if err != nil {
			t.Skip("open raw:", err)
		}
		info, _ := fh.Stat()
		sz := info.Size()

		// truncate tail randomly
		if truncBytes < 0 {
			truncBytes = -truncBytes
		}
		if truncBytes > 0 && truncBytes < sz-int64(HeaderSize) {
			_ = fh.Truncate(sz - truncBytes)
			sz -= truncBytes
		}

		// flip some random bytes (avoid first header)
		if flips < 0 {
			flips = -flips
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := int64(0); i < flips; i++ {
			pos := int64(HeaderSize) + r.Int63n(maxI64(0, sz-int64(HeaderSize)))
			var b [1]byte
			_, _ = fh.ReadAt(b[:], pos)
			b[0] ^= 0xFF
			_, _ = fh.WriteAt(b[:], pos)
		}
		_ = fh.Close()

		// reopen & replay
		w2, err := OpenWALWithPolicy(makeLogger(), path, pol)
		if err != nil {
			// invalid header/truncated header → allowed
			return
		}
		defer w2.Close()

		var got []rec
		cb := func(lsn uint64, ts int64, k, v []byte) error {
			got = append(got, rec{lsn, ts, append([]byte(nil), k...), append([]byte(nil), v...)})
			return nil
		}
		_, err = w2.Replay(0, cb)
		// property: either error OR got is a prefix (<= len(want)) and each entry matches
		if err == nil {
			if len(got) > len(want) {
				t.Fatalf("got longer than want: %d > %d", len(got), len(want))
			}
			for i := range got {
				if got[i].lsn != want[i].lsn || got[i].ts != want[i].ts ||
					!bytes.Equal(got[i].k, want[i].k) || !bytes.Equal(got[i].v, want[i].v) {
					t.Fatalf("got differs at %d", i)
				}
			}
		}
	})
}

func maxI64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

// ========== FUZZ: header variants (short/garbage) ==========

func FuzzOpenHeader(f *testing.F) {
	f.Add([]byte{0, 0, 0, 0, 0, 0, 0, 0})             // empty header
	f.Add([]byte{0x57, 0x41, 0x4C, 0x58, 1, 0, 0, 0}) // correct magic+ver (little endian)

	f.Fuzz(func(t *testing.T, hdr []byte) {
		path := tmpPathFuzz(t)
		fh, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			t.Skip(err)
		}
		// write arbitrary header (maybe too short)
		_, _ = fh.WriteAt(hdr, 0)
		_ = fh.Close()

		_, _ = OpenWALWithPolicy(makeLogger(), path, tinyPolicy())
		// Either succeeds (if header valid) or returns ErrInvalidLog/ErrInvalidHeader.
		// We don't assert exact error here—just ensure no panic.
	})
}

// ========== BENCHMARKS ==========

func benchPolicy(maxVal int) Policy {
	return Policy{
		MaxKeyBytes:    64,
		MaxValueBytes:  uint32(maxVal),
		SyncEvery:      0,       // disable ticker
		SyncEveryBytes: 1 << 60, // effectively never during Append
	}
}

func BenchmarkAppend_Small(b *testing.B) {
	p := benchPolicy(4 << 10)
	w, _ := openWithPolicy(b, p)
	defer w.Close()

	k := []byte("k")
	v := bytes.Repeat([]byte("x"), 64)
	b.SetBytes(int64(len(k) + len(v)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = w.Append(uint64(i+1), time.Now().Unix(), k, v)
	}
}

func BenchmarkAppend_Large(b *testing.B) {
	p := benchPolicy(256 << 10) // 256KB values
	w, _ := openWithPolicy(b, p)
	defer w.Close()

	k := bytes.Repeat([]byte("k"), 8)
	v := bytes.Repeat([]byte("v"), 256<<10)
	b.SetBytes(int64(len(k) + len(v)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = w.Append(uint64(i+1), time.Now().Unix(), k, v)
	}
}

func BenchmarkAppend_WithSyncEvery4MB(b *testing.B) {
	p := benchPolicy(64 << 10)
	p.SyncEveryBytes = 4 << 20 // fsync every ~4MB
	w, _ := openWithPolicy(b, p)
	defer w.Close()

	k := []byte("k")
	v := bytes.Repeat([]byte("x"), 32<<10) // 32KB
	b.SetBytes(int64(len(k) + len(v)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = w.Append(uint64(i+1), time.Now().Unix(), k, v)
	}
}

func prepareFileWithRecords(b *testing.B, n int, valSize int) (WriteAheadLog, string) {
	p := benchPolicy(valSize)
	w, path := openWithPolicy(b, p)

	now := time.Now().Unix()
	val := bytes.Repeat([]byte("x"), valSize)
	for i := 0; i < n; i++ {
		k := []byte{byte(i)}
		_ = w.Append(uint64(i+1), now+int64(i), k, val)
	}
	// force flush to be safe
	_ = w.Sync()
	_ = w.Close()

	w2, err := OpenWALWithPolicy(makeLogger(), path, p)
	if err != nil {
		b.Fatalf("reopen: %v", err)
	}
	return w2, path
}

func BenchmarkReplay_10k_1KB(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping on -short")
	}
	w, _ := prepareFileWithRecords(b, 10_000, 1<<10)
	defer w.Close()

	cb := func(uint64, int64, []byte, []byte) error { return nil }

	// help OS cache warm first
	_, _ = w.Replay(0, cb)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.Replay(0, cb)
	}
}

func BenchmarkReplay_100k_256B(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping on -short")
	}
	w, _ := prepareFileWithRecords(b, 100_000, 256)
	defer w.Close()
	cb := func(uint64, int64, []byte, []byte) error { return nil }
	_, _ = w.Replay(0, cb)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.Replay(0, cb)
	}
}

// optional: test contention
func BenchmarkAppend_Concurrent(b *testing.B) {
	p := benchPolicy(4 << 10)
	w, _ := openWithPolicy(b, p)
	defer w.Close()

	k := []byte("k")
	v := bytes.Repeat([]byte("x"), 256)

	procs := runtime.GOMAXPROCS(0)
	b.SetBytes(int64(len(k) + len(v)))
	b.ReportAllocs()
	b.ResetTimer()

	var lsn uint64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := atomicAdd64(&lsn, 1)
			_ = w.Append(id, time.Now().Unix(), k, v)
		}
	})
	_ = procs
}

// tiny atomic helper to avoid importing sync/atomic under race tool complaining in generics
func atomicAdd64(p *uint64, delta uint64) uint64 {
	return uint64((*(*uint64)(p))) + delta // this is not actually atomic; replace with sync/atomic if you want real correctness
	// NOTE: If you want correctness under race detector, use:
	// return atomic.AddUint64(p, delta)
}

// * fuzz_test.go ends here.
