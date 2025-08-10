// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// filereader_test.go --- File reader tests.
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
//
//

// * Package:

package fileio

// * Imports:

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	TestingFile    string = "test.txt"
	SymlinkFile    string = "symlink.txt"
	TestingBadFile string = "notexists.txt"
	FileContents   string = "This is a test file.\n\n"
)

// * Code:

func writeTemp(t *testing.T, size int) string {
	t.Helper()

	dir := t.TempDir()
	p := filepath.Join(dir, "data.bin")

	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// deterministic content
	buf := make([]byte, 8192)
	for i := range buf {
		buf[i] = byte(i % 251)
	}

	written := 0
	for written < size {
		n := size - written

		if n > len(buf) {
			n = len(buf)
		}

		if _, err := f.Write(buf[:n]); err != nil {
			t.Fatal(err)
		}

		written += n
	}

	return p
}

func collectAll(sr StreamResult) ([]Chunk, error) {
	var chunks []Chunk

	for c := range sr.ChunkCh {
		chunks = append(chunks, c)
	}

	return chunks, sr.Wait()
}

func TestValidReader(t *testing.T) {
	var inst Reader

	t.Run("Construction", func(t *testing.T) {
		var err error

		inst, err = NewReaderWithFile(TestingFile)

		if inst == nil {
			t.Fatal("Instance is invalid")
		}

		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		if inst.Filename() != TestingFile {
			t.Fatalf("Unexpected filename: %#v != %#v",
				inst.Filename(),
				TestingFile)
		}
	})

	t.Run("Exists", func(t *testing.T) {
		ok, err := inst.Exists()

		if err != nil {
			t.Fatalf("Returned error: %#v", err)
		}

		if !ok {
			t.Fatal("No error, but also not ok.")
		}
	})

	t.Run("Load", func(t *testing.T) {
		data, err := inst.Load()

		if err != nil {
			t.Fatalf("Returned error: %#v", err)
		}

		datastr := string(data)
		if datastr != FileContents {
			t.Fatalf("Unexpected data: %#v != %#v",
				datastr,
				FileContents)
		}
	})
}

func TestSymlinkReader(t *testing.T) {
	var (
		inst1 Reader
		inst2 Reader
		err   error
	)

	inst1, err = NewReaderWithFile(SymlinkFile)
	switch {
	case inst1 == nil:
		t.Fatal("Instance is invalid")

	case err != nil:
		t.Fatalf("Unexpected error: %#v", err)
	}
	if inst1.Filename() != SymlinkFile {
		t.Fatalf("Unexpected filename: %#v != %#v",
			inst1.Filename(),
			SymlinkFile)
	}

	inst2, err = NewReaderWithFileAndOptions(SymlinkFile,
		ReadOptions{FollowSymlinks: false})
	switch {
	case inst2 == nil:
		t.Fatal("Instance is invalid")

	case err != nil:
		t.Fatalf("Unexpected error: %#v", err)
	}
	if inst2.Filename() != SymlinkFile {
		t.Fatalf("Unexpected filename: %#v != %#v",
			inst2.Filename(),
			SymlinkFile)
	}

	t.Run("Exists succeeds with FollowSymlinks", func(t *testing.T) {
		res, err := inst1.Exists()

		if err != nil {
			t.Fatalf("Expecting error: %#v", err)

		}

		if !res {
			t.Error("Did not resolve symlink.")
		}
	})

	t.Run("Exists fails with no FollowSymlinks", func(t *testing.T) {
		res, err := inst2.Exists()

		switch {
		case err == nil:
			t.Fatal("Expecting error.")

		case !errors.Is(err, ErrNotRegular):
			t.Fatalf("Unexpected error: %#v", err)
		}

		if res {
			t.Error("Expected an error")
		}
	})

}

func TestInvalidReader(t *testing.T) {
	var inst Reader

	inst, _ = NewReaderWithFile(TestingBadFile)

	t.Run("Exists", func(t *testing.T) {
		ok, err := inst.Exists()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if ok {
			t.Fatal("Non-existent file apparently exists")
		}
	})

	t.Run("Load", func(t *testing.T) {
		_, err := inst.Load()

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	})
}

func TestStream_Happy(t *testing.T) {
	var (
		shouldTotal int   = 10123
		chunkSize   int   = 4096
		offset      int64 = int64(chunkSize)
		total       int   = 0
	)

	path := writeTemp(t, shouldTotal)
	fr, _ := NewReaderWithFile(path) // whatever ctor you have
	sr := fr.Stream(context.Background(), chunkSize, 2, 0)
	defer sr.Close()

	chunks, err := collectAll(sr)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	for _, chunk := range chunks {
		t.Logf("Adding %v", chunk.Offset)
		total += len(chunk.Data)

		if chunk.Offset > 0 && chunk.Offset != offset {
			t.Fatalf("Offset error: %v != %v", chunk.Offset, offset)
		}

		offset += chunk.Offset
	}

	if total != shouldTotal {
		t.Errorf("Total error: %v != %v", total, shouldTotal)
	}
}

func TestStream_LimitSmallerThanFile(t *testing.T) {
	path := writeTemp(t, 10_000)
	rd, _ := NewReaderWithFile(path)

	const limit = 5_000
	sr := rd.Stream(context.Background(), 2048, 2, limit)
	defer sr.Close()

	chunks, err := collectAll(sr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var total int
	for _, c := range chunks {
		total += len(c.Data)
	}
	if got := int64(total); got != limit {
		t.Fatalf("total=%d want %d", got, limit)
	}
}

func TestStream_ImmediateCancel_NoRead(t *testing.T) {
	path := writeTemp(t, 1<<20) // 1MiB
	rd, _ := NewReaderWithFile(path)

	ctx, cancel := context.WithCancel(context.Background())
	sr := rd.Stream(ctx, 64*1024, 1, 0)
	cancel()         // cancel immediately
	defer sr.Close() // idempotent

	err := sr.Wait()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
}

func TestStream_CancelMidStream_Backpressure(t *testing.T) {
	path := writeTemp(t, 8<<20) // 8MiB
	rd, _ := NewReaderWithFile(path)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sr := rd.Stream(ctx, 256*1024, 1, 0) // small buf to create backpressure

	// read one chunk to let producer start
	select {
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting first chunk")
	case _, ok := <-sr.ChunkCh:
		if !ok {
			t.Fatal("channel closed too early")
		}
	}

	// now cancel; producer must exit quickly without blocking on send
	sr.Close()

	if err := sr.Wait(); err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
}

func TestStream_ParentContextCancel(t *testing.T) {
	path := writeTemp(t, 4<<20)
	rd, _ := NewReaderWithFile(path)

	parent, cancel := context.WithCancel(context.Background())
	sr := rd.Stream(parent, 128*1024, 2, 0)

	done := make(chan struct{})
	go func() {
		defer close(done)

		select {
		case <-time.After(2 * time.Second):
			t.Fatalf("Timeout waiting for first chunk.")
			cancel()

		case _, ok := <-sr.ChunkCh:
			if ok {
				cancel()
			}
		}
	}()

	err := sr.Wait()
	<-done

	if err == nil || !errors.Is(err, context.Canceled) {
		t.Errorf("Want context.Canceled, got %#v", err)
	}
}

func TestStream_OpenError(t *testing.T) {
	rd, _ := NewReaderWithFile(filepath.Join(t.TempDir(), "nope.dne"))
	sr := rd.Stream(context.Background(), 4096, 1, 0)
	defer sr.Close()

	_, err := collectAll(sr)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// Accept either fs.ErrNotExist or a wrapped variant
	if !errors.Is(err, fs.ErrNotExist) && !strings.Contains(err.Error(), "no such file") {
		t.Fatalf("want not-exist error, got %v", err)
	}
}

func TestStream_EOF_NoError(t *testing.T) {
	path := writeTemp(t, 12345)
	rd, _ := NewReaderWithFile(path)
	sr := rd.Stream(context.Background(), 2048, 2, 0)
	defer sr.Close()

	chunks, err := collectAll(sr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var total int
	for _, c := range chunks {
		total += len(c.Data)
	}
	if total != 12345 {
		t.Fatalf("total=%d want %d", total, 12345)
	}
}

func TestStream_DefaultTunings(t *testing.T) {
	path := writeTemp(t, 50_000)
	rd, _ := NewReaderWithFile(path)
	sr := rd.Stream(context.Background(), 0, 0, 0) // defaults
	defer sr.Close()

	chunks, err := collectAll(sr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var total int
	for _, c := range chunks {
		if len(c.Data) <= 0 {
			t.Fatalf("empty chunk encountered")
		}
		total += len(c.Data)
	}
	if total != 50_000 {
		t.Fatalf("total=%d want %d", total, 50_000)
	}
}

func TestStream_CloseIdempotent(t *testing.T) {
	path := writeTemp(t, 200_000)
	rd, _ := NewReaderWithFile(path)
	sr := rd.Stream(context.Background(), 16*1024, 2, 0)

	// normal drain
	_, err := collectAll(sr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// call Close multiple times; should not panic or change outcome
	sr.Close()
	sr.Close()
}

func TestStream_NoLeak_Smoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skip in -short")
	}
	before := runtime.NumGoroutine()

	path := writeTemp(t, 3<<20)
	rd, _ := NewReaderWithFile(path)
	sr := rd.Stream(context.Background(), 64*1024, 2, 0)
	_, err := collectAll(sr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// allow sched to settle
	time.Sleep(50 * time.Millisecond)
	after := runtime.NumGoroutine()

	// Be generous to avoid flakiness; we just want to catch egregious leaks
	if after-before > 5 {
		t.Fatalf("goroutines before=%d after=%d (possible leak)", before, after)
	}
}

// * filereader_test.go ends here.
