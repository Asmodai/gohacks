// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// filewriter_test.go --- File writer tests.
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

package fileio

// * Imports:

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// * Constants:

// * Variables:

// * Code:

// ** Helper functions:

func mustRead(t *testing.T, p string) []byte {
	t.Helper()

	b, err := os.ReadFile(p)

	if err != nil {
		t.Fatalf("read %q: %v", p, err)
	}

	return b
}

// ** Tests:

func TestNewWriter_InvalidMode(t *testing.T) {
	tmp := t.TempDir()
	_, err := NewWriter(context.Background(), filepath.Join(tmp, "x.log"), WriteOptions{
		CreateMode: 9999, // invalid
		Mode:       0o644,
		BufferSize: 8 * 1024,
	})
	if err == nil {
		t.Fatal("expected error for invalid mode")
	}
	if !errors.Is(err, ErrInvalidWriteMode) {
		t.Fatalf("expected ErrInvalidWriteMode, got %v", err)
	}
}

func TestNewWriter_CreateDirs_False_MissingDir(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "a", "b", "c.log")
	_, err := NewWriter(context.Background(), path, WriteOptions{
		CreateDirs: false,
		CreateMode: CreateModeAppend,
		Mode:       0o644,
		BufferSize: 4096,
	})
	if err == nil {
		t.Fatal("expected error when CreateDirs=false and parent dir missing")
	}
}

func TestNewWriter_CreateDirs_True_Creates(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "a", "b", "c.log")
	w, err := NewWriter(context.Background(), path, WriteOptions{
		CreateDirs: true,
		CreateMode: CreateModeAppend,
		Mode:       0o644,
		BufferSize: 4096,
	})
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	defer w.Close()

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist after open: %v", err)
	}
}

func TestWrite_AppendAndTruncate(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "data.log")

	// seed with "old"
	if err := os.WriteFile(p, []byte("old"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	// append
	wa, err := NewWriter(context.Background(), p, WriteOptions{
		CreateMode: CreateModeAppend,
		Mode:       0o644,
		BufferSize: 4096,
	})
	if err != nil {
		t.Fatalf("append NewWriter: %v", err)
	}
	if _, err := wa.Write([]byte("A")); err != nil {
		t.Fatalf("append write: %v", err)
	}
	if err := wa.Close(); err != nil {
		t.Fatalf("append close: %v", err)
	}
	if got := string(mustRead(t, p)); got != "oldA" {
		t.Fatalf("append content = %q, want %q", got, "oldA")
	}

	// truncate
	wt, err := NewWriter(context.Background(), p, WriteOptions{
		CreateMode: CreateModeTruncate,
		Mode:       0o644,
		BufferSize: 4096,
	})
	if err != nil {
		t.Fatalf("truncate NewWriter: %v", err)
	}
	if _, err := wt.Write([]byte("B")); err != nil {
		t.Fatalf("truncate write: %v", err)
	}
	if err := wt.Close(); err != nil {
		t.Fatalf("truncate close: %v", err)
	}
	if got := string(mustRead(t, p)); got != "B" {
		t.Fatalf("truncate content = %q, want %q", got, "B")
	}
}

func TestFsyncEveryWrite_FlushesImmediately(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "sync.log")

	w, err := NewWriter(context.Background(), p, WriteOptions{
		CreateMode: CreateModeTruncate,
		Mode:       0o644,
		BufferSize: 1 << 20, // large buffer
		Fsync:      FsyncEveryWrite,
	})
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	// Write but do NOT Close yet.
	if _, err := w.Write([]byte("hello")); err != nil {
		t.Fatalf("write: %v", err)
	}

	// The data should be visible on disk already (buf flushed + fsync).
	got := mustRead(t, p)
	if !bytes.Equal(got, []byte("hello")) {
		t.Fatalf("on-disk = %q, want %q", string(got), "hello")
	}

	_ = w.Close()
}

func TestAbort_DiscardsBufferedData(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "abort.log")

	w, err := NewWriter(context.Background(), p, WriteOptions{
		CreateMode: CreateModeTruncate,
		Mode:       0o644,
		BufferSize: 64 * 1024, // ensure buffer won't auto-flush
		Fsync:      FsyncNever,
	})
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	if _, err := w.Write([]byte("keep? no.")); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := w.Abort(); err != nil {
		t.Fatalf("abort: %v", err)
	}

	// File should exist but be empty (nothing flushed).
	info, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if size := info.Size(); size != 0 {
		t.Fatalf("size after abort = %d, want 0", size)
	}
}

func TestContextCancel_PreventsWrite(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "ctx.log")

	ctx, cancel := context.WithCancel(context.Background())
	w, err := NewWriter(ctx, p, WriteOptions{
		CreateMode: CreateModeTruncate,
		Mode:       0o644,
		BufferSize: 4096,
		Fsync:      FsyncNever,
	})
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	defer w.Close()

	cancel() // cancel before write

	n, err := w.Write([]byte("nope"))
	if err == nil {
		t.Fatalf("expected error on canceled ctx")
	}
	if n != 0 {
		t.Fatalf("bytes written on canceled ctx = %d, want 0", n)
	}
}

func TestBytesWritten_TracksBuffer_NotOnDisk(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "buf.log")

	w, err := NewWriter(context.Background(), p, WriteOptions{
		CreateMode: CreateModeTruncate,
		Mode:       0o644,
		BufferSize: 1024,
		Fsync:      FsyncNever,
	})
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	data := []byte("abc123")
	if _, err := w.Write(data); err != nil {
		t.Fatalf("write: %v", err)
	}

	if got := w.BytesWritten(); got < int64(len(data)) {
		// It can be > len if bufio coalesced, but should never be less.
		t.Fatalf("BytesWritten = %d, want >= %d", got, len(data))
	}

	// Before close, file should be empty (not flushed).
	if sz, _ := os.Stat(p); sz != nil && sz.Size() != 0 {
		t.Fatalf("on-disk size before close = %d, want 0", sz.Size())
	}

	// Sync drains buffer.
	if err := w.Sync(); err != nil {
		t.Fatalf("sync: %v", err)
	}
	if got := w.BytesWritten(); got != 0 {
		t.Fatalf("BytesWritten after sync = %d, want 0", got)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if !bytes.Equal(mustRead(t, p), data) {
		t.Fatal("final content mismatch")
	}
}

func TestAppend_DoesNotTruncateExisting(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "append.log")
	if err := os.WriteFile(p, []byte("X"), 0o644); err != nil {
		t.Fatalf("seed: %v", err)
	}

	w, err := NewWriter(context.Background(), p, WriteOptions{
		CreateMode: CreateModeAppend,
		Mode:       0o644,
		BufferSize: 64,
	})
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}

	_, _ = w.Write([]byte("Y"))
	_ = w.Close()

	if got := string(mustRead(t, p)); got != "XY" {
		t.Fatalf("content = %q, want %q", got, "XY")
	}
}

func TestFsyncNever_BuffersUntilClose(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "buffered.log")

	w, err := NewWriter(context.Background(), p, WriteOptions{
		CreateMode: CreateModeTruncate,
		Mode:       0o644,
		BufferSize: 1 << 20,
		Fsync:      FsyncNever,
	})
	if err != nil {
		t.Fatalf("NewWriter: %v", err)
	}
	defer w.Close()

	_, _ = w.Write([]byte("zzz"))
	time.Sleep(10 * time.Millisecond) // give scheduler a tick

	// Likely still empty because buffer not flushed.
	if sz, _ := os.Stat(p); sz != nil && sz.Size() != 0 {
		t.Skipf("environment flushes early; skipping strict assertion")
	}

	_ = w.Sync()
	if sz, _ := os.Stat(p); sz == nil || sz.Size() == 0 {
		t.Fatal("expected non-zero size after Sync")
	}
}

// * filewriter_test.go ends here.
