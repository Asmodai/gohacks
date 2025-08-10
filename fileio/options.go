// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// options.go --- Options.
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
	"math"
	"os"
	"time"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	FsyncNever      FsyncPolicy = iota // Never
	FsyncOnClose                       // Synchronise on close.
	FsyncEveryWrite                    // Synchronise every write.

	// Default temporary file prefix.
	DefaultTempFilePrefix = "."

	// Default temporary file suffix.
	DefaultTempFileSuffix = ".tmp"
)

// * Code:

// ** Read Options:

// *** Type:

// File I/O options.
type ReadOptions struct {
	// Maximum number of bytes to read.
	//
	// If this value is set to a negative number then it is interpreted
	// as "read an unlimited number of bytes".
	//
	// However "unlimited", in this case, equates to `math.MaxInt` bytes.
	//
	// The reason for this is that we return a slice of bytes, and the
	// maximum number of elements in a Go slice is `math.MaxInt`.
	//
	// The default value is 0.
	MaxReadBytes int64

	// Should symbolic links be followed?
	//
	// If false, then symbolic links are not followed.
	//
	// The default value is false.
	FollowSymlinks bool
}

// *** methods:

// Sanitise options.
func (opt *ReadOptions) sanity() error {
	if opt.MaxReadBytes < 0 {
		opt.MaxReadBytes = 0
	}

	if opt.MaxReadBytes > int64(math.MaxInt) {
		return errors.WithMessagef(
			ErrInvalidSize,
			"MaxReadByts %d",
			opt.MaxReadBytes)
	}

	return nil
}

// ** Write Options:

// *** Constants:

const (
	CreateModeTruncate CreateMode = iota // Truncate files before writing.
	CreateModeAppend                     // Append to file during writing.
)

// *** Types:

type FsyncPolicy int

type CreateMode int

type WriteOptions struct {
	// File mode.
	//
	// Default is 0o644.
	Mode os.FileMode

	// File creation mode.
	//
	// Can be one of "truncate" or "append".
	CreateMode CreateMode

	// Create directories should they not exist?
	CreateDirs bool

	// Buffer size.
	//
	// If the value is zero, the value in`DefaultWriteBufferSize` shall
	// be used.
	BufferSize int

	// File sync policy.
	Fsync FsyncPolicy

	// Synchronise file every n writes.
	//
	// If zero, no syncs will be performed.
	FsyncEveryN int64

	// Context deadline.
	Timeout time.Duration
}

// *** Methods:

func (opt *WriteOptions) sanity() error {
	if opt.Mode == 0 {
		opt.Mode = DefaultFileMode
	}

	if opt.BufferSize <= 0 {
		opt.BufferSize = defaultWriteBufferSize
	}

	if opt.CreateMode != CreateModeAppend && opt.CreateMode != CreateModeTruncate {
		return errors.WithMessagef(
			ErrInvalidWriteMode,
			"invalid create mode %q",
			opt.CreateMode)
	}

	return nil
}

// * options.go ends here.
