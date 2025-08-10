// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// wal_test.go --- Write Ahead Log tests.
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
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

// * Variables:

// * Code:

// ** helpers:

func makeLogger() context.Context {
	ctx := context.Background()

	lgr := logger.NewDefaultLogger()
	nctx, _ := logger.SetLogger(ctx, lgr)

	return nctx
}

func tmpPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test.walx")
}

func open(t *testing.T, pol Policy) (WriteAheadLog, string) {
	t.Helper()
	p := tmpPath(t)
	w, err := OpenWALWithPolicy(makeLogger(), p, pol)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	return w, p
}

func mustAppend(t *testing.T, w WriteAheadLog, lsn uint64, ts int64, k, v []byte) {
	t.Helper()
	if err := w.Append(lsn, ts, k, v); err != nil {
		t.Fatalf("append: %v", err)
	}
}

func collectReplay(w WriteAheadLog, base uint64) (applied []struct {
	LSN uint64
	TS  int64
	K   []byte
	V   []byte
}, err error) {
	var mu sync.Mutex
	cb := func(lsn uint64, tstamp int64, key, val []byte) error {
		mu.Lock()
		applied = append(applied, struct {
			LSN uint64
			TS  int64
			K   []byte
			V   []byte
		}{lsn, tstamp, append([]byte(nil), key...), append([]byte(nil), val...)})
		mu.Unlock()
		return nil
	}
	_, err = w.Replay(base, cb)
	return
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}

func tinyPolicy() Policy {
	// very small limits to exercise boundaries quickly
	return Policy{
		MaxKeyBytes:    64,
		MaxValueBytes:  1024,
		SyncEvery:      0,
		SyncEveryBytes: 0,
	}
}

// ** Tests:

// *** happy path:

func TestAppendReplay_Happy(t *testing.T) {
	t.Parallel()

	w, _ := open(t, tinyPolicy())
	defer w.Close()

	now := time.Now().Unix()
	mustAppend(t, w, 1, now, []byte("k1"), []byte("v1"))
	mustAppend(t, w, 2, now+1, []byte("k2"), []byte("v2"))

	got, err := collectReplay(w, 0)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if len(got) != 2 || got[0].LSN != 1 || got[1].LSN != 2 {
		t.Fatalf("unexpected replay: %#v", got)
	}
}

// *** sync policy: bytes threshold:

func TestSyncEveryBytes_AutoFlush(t *testing.T) {
	t.Parallel()

	p := tinyPolicy()
	p.SyncEveryBytes = 1 // tiny threshold, should sync after first append
	w, path := open(t, p)
	defer w.Close()

	now := time.Now().Unix()
	mustAppend(t, w, 10, now, []byte("a"), []byte("b"))

	// Close and reopen; if sync happened, the record persists even if we didn't call Sync().
	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	w2, err := OpenWALWithPolicy(makeLogger(), path, p)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer w2.Close()

	got, err := collectReplay(w2, 0)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if len(got) != 1 || got[0].LSN != 10 {
		t.Fatalf("expected 1 record after auto-sync, got %#v", got)
	}
}

// *** close must flush dirty:

func TestClose_FlushesDirty(t *testing.T) {
	t.Parallel()

	p := tinyPolicy()
	p.SyncEveryBytes = 0 // disable auto
	w, path := open(t, p)

	now := time.Now().Unix()
	mustAppend(t, w, 42, now, []byte("x"), []byte("y"))
	// do not call Sync(), just Close()
	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	// reopen: record should be there
	w2, err := OpenWALWithPolicy(makeLogger(), path, p)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer w2.Close()
	got, err := collectReplay(w2, 0)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if len(got) != 1 || got[0].LSN != 42 {
		t.Fatalf("expected flushed record, got %#v", got)
	}
}

// *** reset:

func TestReset_TruncatesToHeader(t *testing.T) {
	t.Parallel()

	w, path := open(t, tinyPolicy())
	defer w.Close()

	mustAppend(t, w, 1, time.Now().Unix(), []byte("k"), []byte("v"))

	if err := w.Reset(); err != nil {
		t.Fatalf("reset: %v", err)
	}

	// Close + reopen and verify nothing replays
	_ = w.Close()
	w2, err := OpenWALWithPolicy(makeLogger(), path, tinyPolicy())
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer w2.Close()

	got, err := collectReplay(w2, 0)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty after reset, got %#v", got)
	}
}

// *** size limits:

func TestAppend_RejectsTooLargeKeyOrValue(t *testing.T) {
	t.Parallel()

	p := tinyPolicy()
	w, _ := open(t, p)
	defer w.Close()

	now := time.Now().Unix()

	tooBigKey := make([]byte, int(p.MaxKeyBytes)+1)
	if err := w.Append(1, now, tooBigKey, []byte("ok")); err == nil {
		t.Fatalf("expected error for too-large key")
	}

	tooBigVal := make([]byte, int(p.MaxValueBytes)+1)
	if err := w.Append(2, now, []byte("ok"), tooBigVal); err == nil {
		t.Fatalf("expected error for too-large value")
	}

	// boundary is ok
	key := make([]byte, int(p.MaxKeyBytes))
	val := make([]byte, int(p.MaxValueBytes))
	if err := w.Append(3, now, key, val); err != nil {
		t.Fatalf("boundary append failed: %v", err)
	}
}

func TestAppend_RejectsNegativeTimestamp(t *testing.T) {
	t.Parallel()

	w, _ := open(t, tinyPolicy())
	defer w.Close()

	if err := w.Append(1, -1, []byte("k"), []byte("v")); err == nil {
		t.Fatalf("expected error for negative timestamp")
	}
}

// *** replay: truncation tail is ignored:

func TestReplay_IgnoresTruncatedTail(t *testing.T) {
	t.Parallel()

	p := tinyPolicy()
	w, path := open(t, p)
	defer w.Close()

	now := time.Now().Unix()
	mustAppend(t, w, 1, now, []byte("k1"), []byte("v1"))
	mustAppend(t, w, 2, now, []byte("k2"), []byte("v2"))

	// Corrupt the file by truncating last few bytes of the last record.
	_ = w.Close()
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("open raw: %v", err)
	}
	info, _ := f.Stat()
	// chop 3 bytes
	if err := f.Truncate(info.Size() - 3); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	_ = f.Close()

	w2, err := OpenWALWithPolicy(makeLogger(), path, p)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer w2.Close()

	got, err := collectReplay(w2, 0)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if len(got) != 1 || got[0].LSN != 1 {
		t.Fatalf("expected only first record, got %#v", got)
	}
}

// *** replay: CRC mismatch stops at boundary:

func TestReplay_StopsOnCRCMismatch(t *testing.T) {
	t.Parallel()

	p := tinyPolicy()
	w, path := open(t, p)

	now := time.Now().Unix()
	mustAppend(t, w, 1, now, []byte("k1"), []byte("v1"))
	mustAppend(t, w, 2, now, []byte("k2"), []byte("v2"))

	// flip a byte in the CRC of the second record
	_ = w.Close()

	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatalf("open raw: %v", err)
	}
	defer f.Close()

	// Walk file: header + record1 + record2; we’ll flip the last byte of file.
	info, _ := f.Stat()
	if info.Size() <= int64(headerSize)+10 {
		t.Fatalf("file too small to corrupt meaningfully")
	}
	// read last byte, flip, write back
	var b [1]byte
	_, _ = f.ReadAt(b[:], info.Size()-1)
	b[0] ^= 0xFF
	if _, err := f.WriteAt(b[:], info.Size()-1); err != nil {
		t.Fatalf("write flip: %v", err)
	}

	w2, err := OpenWALWithPolicy(makeLogger(), path, p)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer w2.Close()

	got, err := collectReplay(w2, 0)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	// CRC mismatch on second means only first is applied
	if len(got) != 1 || got[0].LSN != 1 {
		t.Fatalf("expected only first record, got %#v", got)
	}
}

// *** open: invalid header / short file:

func TestOpen_InvalidHeader(t *testing.T) {
	t.Parallel()

	path := tmpPath(t)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// write bad magic/version
	var hdr [headerSize]byte
	binary.LittleEndian.PutUint32(hdr[0:4], 0xDEADBEEF)
	binary.LittleEndian.PutUint32(hdr[4:8], 99)
	_, _ = f.WriteAt(hdr[:], 0)
	_ = f.Close()

	_, err = OpenWALWithPolicy(makeLogger(), path, tinyPolicy())
	if !errors.Is(err, ErrInvalidHeader) {
		t.Fatalf("expected ErrInvalidHeader, got %v", err)
	}
}

func TestOpen_FileTooShort(t *testing.T) {
	t.Parallel()

	path := tmpPath(t)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// write only 4 bytes < header
	_, _ = f.WriteAt([]byte{1, 2, 3, 4}, 0)
	_ = f.Close()

	_, err = OpenWALWithPolicy(makeLogger(), path, tinyPolicy())
	if !errors.Is(err, ErrInvalidLog) {
		t.Fatalf("expected ErrInvalidLog, got %v", err)
	}
}

// *** replay: timestamp too big (inject record manually):

func TestReplay_TimestampTooBig(t *testing.T) {
	t.Parallel()

	p := tinyPolicy()
	path := tmpPath(t)
	w, err := OpenWALWithPolicy(makeLogger(), path, p)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	// craft one record with timestamp > MaxInt64
	key := []byte("k")
	val := []byte("v")
	bufLen := uint32(recordSize(uint32(len(key)), uint32(len(val))))
	buf := make([]byte, int(bufLen))
	encodeRecord(buf, bufLen, 7, ^uint64(0), uint32(len(key)), uint32(len(val)), key, val)
	finaliseCRC(buf, int(bufLen), w.(*writeAheadLog).crcTab)

	// write directly after header
	if err := writeFullAt(w.(*writeAheadLog).fptr, buf, int64(headerSize)); err != nil {
		t.Fatalf("inject: %v", err)
	}
	_ = w.Close()

	w2, err := OpenWALWithPolicy(makeLogger(), path, p)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer w2.Close()

	var gotErr atomic.Value
	cb := func(lsn uint64, ts int64, k, v []byte) error { return nil }
	_, err = w2.Replay(0, cb)
	// decodeFields should fail with ErrTimestampTooBig → bubbled up
	if err == nil || !errors.Is(err, ErrTimestampTooBig) {
		_ = gotErr // quiet staticcheck
		t.Fatalf("expected ErrTimestampTooBig, got %v", err)
	}
}

// *** ticker-based flush (best-effort; avoid flakiness):

func TestTickerFlushesEventually(t *testing.T) {
	t.Parallel()

	p := tinyPolicy()
	p.SyncEvery = 15 * time.Millisecond
	p.SyncEveryBytes = 0

	w, path := open(t, p)
	defer w.Close()

	now := time.Now().Unix()
	mustAppend(t, w, 1, now, []byte("k"), bytes.Repeat([]byte("x"), 16))

	// wait > 2 ticks
	time.Sleep(40 * time.Millisecond)

	// close & reopen; data should be flushed by ticker
	_ = w.Close()
	w2, err := OpenWALWithPolicy(makeLogger(), path, p)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer w2.Close()

	got, err := collectReplay(w2, 0)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if len(got) != 1 || got[0].LSN != 1 {
		t.Fatalf("expected 1 record after ticker flush, got %#v", got)
	}
}

func TestTickerRaceDoesNotPanic(t *testing.T) {
	pol := Policy{SyncEvery: 1 * time.Millisecond, MaxKeyBytes: 16, MaxValueBytes: 64}
	w, _ := OpenWALWithPolicy(makeLogger(), tmpPath(t), pol)
	defer w.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		// hammer SetPolicy/stop/start while appending
		for i := 0; i < 200; i++ {
			_ = w.Append(uint64(i+1), time.Now().Unix(), []byte("k"), []byte("v"))
			if i%3 == 0 {
				pol.SyncEvery = time.Duration(i%5) * time.Millisecond
				w.SetPolicy(pol)
			}
			time.Sleep(time.Microsecond)
		}
	}()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	}
}

// * wal_test.go ends here.
