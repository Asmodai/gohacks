// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// files.go --- File functions.
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
	"io"
	"os"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

// * Variables:

// * Code:

// ** Types:

type files struct{}

// ** Methods:

// Helper to write to a file.
func (f *files) writeFrom(parent context.Context, path string, rdr io.Reader, opts WriteOptions) (int64, error) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	ctx = parent
	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(parent, opts.Timeout)
		defer cancel()
	}

	writer, err := f.OpenWriter(ctx, path, opts)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	buf := make([]byte, opts.BufferSize)
	crdr := &ctxReader{ctx: ctx, rdr: rdr}

	num, copyErr := io.CopyBuffer(writer, crdr, buf)
	if copyErr != nil {
		_ = writer.Abort()

		return num, errors.WithStack(copyErr)
	}

	if err := writer.Close(); err != nil {
		return num, errors.WithStack(err)
	}

	return num, nil
}

// Open a file for writing.
func (f *files) OpenWriter(ctx context.Context, path string, opts WriteOptions) (Writer, error) {
	return NewWriter(ctx, path, opts)
}

// Append data to a file.
func (f *files) AppendFile(ctx context.Context, path string, rdr io.Reader, opts WriteOptions) (int64, error) {
	opts.CreateMode = CreateModeAppend

	result, err := f.writeFrom(ctx, path, rdr, opts)

	return result, errors.WithStack(err)
}

// Write to a file.
func (f *files) WriteFile(ctx context.Context, path string, data []byte, opts WriteOptions) error {
	_, err := f.writeFrom(ctx, path, bytes.NewReader(data), opts)

	return errors.WithStack(err)
}

// Remove a file.
func (f *files) Remove(path string) error {
	return errors.WithStack(os.Remove(path))
}

// Rename a file.
func (f *files) Rename(oldPath, newPath string) error {
	return errors.WithStack(os.Rename(oldPath, newPath))
}

// Create directory.
func (f *files) MkdirAll(dir string, mode os.FileMode) error {
	if len(dir) == 0 {
		return nil
	}

	err := os.MkdirAll(dir, mode)

	return errors.WithStack(err)
}

// ** Functions:

func NewFiles() Files {
	return &files{}
}

// * files.go ends here.
