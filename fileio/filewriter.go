// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// filewriter.go --- Writer.
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
	"bufio"
	"context"
	"os"
	"path/filepath"

	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Types:

type fileWriter struct {
	ctx  context.Context
	fptr *os.File
	buf  *bufio.Writer
	opts WriteOptions
	name string
}

// ** Methods:

func (w *fileWriter) Name() string        { return w.name }
func (w *fileWriter) BytesWritten() int64 { return int64(w.buf.Buffered()) }

// Write data.
func (w *fileWriter) Write(data []byte) (int, error) {
	if err := w.ctx.Err(); err != nil {
		return 0, errors.WithStack(err)
	}

	num, err := w.buf.Write(data)

	if w.opts.Fsync == FsyncEveryWrite {
		if ferr := w.Sync(); ferr != nil {
			return num, errors.WithStack(ferr)
		}
	}

	return num, errors.WithStack(err)
}

// Perform a synchronisation.
func (w *fileWriter) Sync() error {
	if err := w.buf.Flush(); err != nil {
		return errors.WithStack(err)
	}

	if w.opts.Fsync == FsyncNever {
		return nil
	}

	return errors.WithStack(w.fptr.Sync())
}

// Abort file writing.
func (w *fileWriter) Abort() error {
	w.buf.Reset(nil)

	return errors.WithStack(w.fptr.Close())
}

// Close the file.
func (w *fileWriter) Close() error {
	if err := w.Sync(); err != nil {
		_ = w.fptr.Close()

		return errors.WithStack(err)
	}

	return errors.WithStack(w.fptr.Close())
}

// ** Functions:

func NewWriter(ctx context.Context, path string, opts WriteOptions) (Writer, error) {
	var cmode int

	if err := opts.sanity(); err != nil {
		return nil, errors.WithStack(err)
	}

	if opts.CreateDirs {
		err := os.MkdirAll(filepath.Dir(path), DefaultDirectoryMode)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	switch opts.CreateMode {
	case CreateModeAppend:
		cmode = os.O_CREATE | os.O_APPEND | os.O_WRONLY

	case CreateModeTruncate:
		cmode = os.O_CREATE | os.O_TRUNC | os.O_WRONLY

	default:
		return nil, errors.WithMessagef(
			ErrInvalidWriteMode,
			"%v is not a valid write mode",
			opts.CreateMode)
	}

	fptr, err := os.OpenFile(path, cmode, opts.Mode)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	inst := &fileWriter{
		ctx:  ctx,
		fptr: fptr,
		buf:  bufio.NewWriterSize(fptr, opts.BufferSize),
		opts: opts,
		name: path}

	return inst, nil
}

// * filewriter.go ends here.
