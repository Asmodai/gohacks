// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// filereader.go --- File input I/O.
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

package fileio

// * Imports:

import (
	"context"
	"io"
	"os"
	"sync"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	// Default stream chunk size.
	defaultChunkSize = 64 * 1024

	defaultBufSize = 8
)

// * Code:

// ** Types:

type fileReader struct {
	filename string
	opts     Options
}

// ** Methods:

// The file path that we wish to load.
func (frdr *fileReader) Filename() string {
	return frdr.filename
}

// Check if the file exists.
func (frdr *fileReader) Exists() (bool, error) {
	var (
		finfo os.FileInfo
		err   error
	)

	if frdr.opts.FollowSymlinks {
		finfo, err = os.Stat(frdr.filename)
	} else {
		finfo, err = os.Lstat(frdr.filename)
	}

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, errors.WithStack(err)
	}

	// If we're not following symlinks, then a symlink won't be
	// Mode().IsRegular()
	//
	// If we are following, then `Stat` resolve it and we can check the
	// target.
	if !finfo.Mode().IsRegular() {
		return false, errors.WithStack(ErrNotRegular)
	}

	return true, nil
}

// Check if the file is a symbolic link.
func (frdr *fileReader) IsSymlink() (bool, error) {
	finfo, err := os.Lstat(frdr.filename)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return (finfo.Mode() & os.ModeSymlink) != 0, nil
}

// Read the file and return a byte array of its content.
//
//nolint:cyclop
func (frdr *fileReader) Load() ([]byte, error) {
	file, err := os.Open(frdr.filename)

	switch {
	case os.IsNotExist(err):
		return nil, nil

	case err != nil:
		return nil, errors.WithStack(err)
	}

	defer func() {
		if cerr := file.Close(); cerr != nil && err != nil {
			err = errors.WithStack(err)
		}
	}()

	finfo, statErr := file.Stat()
	if statErr != nil {
		return nil, errors.WithStack(statErr)
	}

	if !finfo.Mode().IsRegular() {
		return nil, errors.WithStack(ErrNotRegular)
	}

	if frdr.opts.MaxReadBytes > 0 && finfo.Size() > frdr.opts.MaxReadBytes {
		return nil, errors.WithMessagef(ErrTooLarge, "%q", frdr.filename)
	}

	if size := finfo.Size(); size > 0 {
		buf := make([]byte, size)

		if _, err = io.ReadFull(file, buf); err != nil {
			return nil, errors.WithStack(err)
		}

		return buf, nil
	}

	// Fallback for rare cases of "unknown size".
	var rdr io.Reader = file

	if frdr.opts.MaxReadBytes > 0 {
		rdr = io.LimitReader(file, frdr.opts.MaxReadBytes+1)
	}

	buf, readErr := io.ReadAll(rdr)
	if readErr != nil {
		return nil, errors.WithStack(readErr)
	}

	if frdr.opts.MaxReadBytes > 0 && int64(len(buf)) > frdr.opts.MaxReadBytes {
		return nil, errors.WithMessagef(ErrTooLarge, "%q", frdr.filename)
	}

	return buf, nil
}

// Open the file and return an `io.ReadCloser`.
func (frdr *fileReader) Open() (io.ReadCloser, error) {
	if !frdr.opts.FollowSymlinks {
		issymlink, err := frdr.IsSymlink()

		switch {
		case err != nil:
			return nil, errors.WithStack(err)

		case issymlink:
			return nil, errors.WithMessagef(
				ErrSymlinkDenied,
				"%q",
				frdr.filename)
		}
	}

	file, err := os.Open(frdr.filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return file, nil
}

func (frdr *fileReader) CopyTo(writer io.Writer, limit int64) (int64, error) {
	rdr, err := frdr.Open()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	defer rdr.Close()

	var reader io.Reader = rdr

	if limit > 0 {
		// CopyN returns EOF when source is shorter than limit.
		num, err := io.CopyN(writer, reader, limit)
		if errors.Is(err, io.EOF) {
			return num, nil
		}

		return num, errors.WithStack(err)
	}

	num, err := io.Copy(writer, reader)
	if err != nil {
		return num, errors.WithStack(err)
	}

	return num, nil
}

// Stream chunks of up to `chunkSize` bytes.
//
//nolint:cyclop,gocognit,funlen
func (frdr *fileReader) Stream(parent context.Context, chunkSize, bufSize int, limit int64) StreamResult {
	if chunkSize <= 0 {
		chunkSize = defaultChunkSize
	}

	if bufSize <= 0 {
		bufSize = defaultBufSize
	}

	out := make(chan Chunk, bufSize)
	errc := make(chan error, 1)

	// Provide a context cancel function.
	ctx, cancel := context.WithCancel(parent)

	// Provide a drain function that can drain chunks so the producer
	// can't block trying to send.
	wait := func() error {
		for range out { //nolint:revive
			// This should not be removed, it drains `out`.
		}

		if err, ok := <-errc; ok {
			return errors.WithStack(err)
		}

		return nil
	}

	go func() {
		defer close(errc)
		defer close(out)

		var sentErr sync.Once

		trySendErr := func(err error) {
			sentErr.Do(func() {
				select {
				case errc <- errors.WithStack(err):
				default:
				}
			})
		}

		fil, err := frdr.Open()
		if err != nil {
			errc <- errors.WithStack(err)

			return
		}
		defer fil.Close()

		rdr := io.Reader(fil)
		if limit > 0 {
			rdr = io.LimitReader(fil, limit)
		}

		buf := make([]byte, chunkSize)

		var off int64

		for {
			if err := ctx.Err(); err != nil {
				trySendErr(err)

				return
			}

			num, rerr := rdr.Read(buf)
			if num > 0 {
				slab := make([]byte, num)
				copy(slab, buf[:num])

				select {
				case out <- Chunk{Offset: off, Data: slab}:
					off += int64(num)

				case <-ctx.Done():
					trySendErr(ctx.Err())

					return
				}
			}

			if rerr != nil {
				if errors.Is(rerr, io.EOF) {
					return
				}

				errc <- errors.WithStack(rerr)

				return
			}
		}
	}()

	return StreamResult{
		ChunkCh: out,
		ErrorCh: errc,
		Close:   cancel,
		Wait:    wait}
}

// ** Functions:

// Create a new FileReader with the given file name.
func NewReaderWithFile(filename string) (Reader, error) {
	return NewReaderWithFileAndOptions(filename, Options{FollowSymlinks: true})
}

// Create a new FileReader with the given file name and options.
func NewReaderWithFileAndOptions(filename string, opts Options) (Reader, error) {
	if err := opts.sanity(); err != nil {
		return nil, errors.WithStack(err)
	}

	inst := &fileReader{
		filename: filename,
		opts: Options{
			MaxReadBytes:   opts.MaxReadBytes,
			FollowSymlinks: opts.FollowSymlinks,
		},
	}

	return inst, nil
}

// * filereader.go ends here.
