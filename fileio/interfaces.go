// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// reader.go --- Reader interface.
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
//
//mock:yes

// * Comments:

// * Package:

package fileio

// * Imports:

import (
	"context"
	"io"
	"os"
)

// * Constants:

// * Variables:

// * Code:

// ** Types:

// Data chunk used when dealing with streamed file data.
type Chunk struct {
	Offset int64  // Offset into the file.
	Data   []byte // Data slab.
}

// Result of opening a file for streaming read.
type StreamResult struct {
	ChunkCh <-chan Chunk // Channel for chunks.
	ErrorCh <-chan error // Channel for errors.
	Close   func()       // Cancel the stream.
	Wait    func() error // Wait for completion.
}

// ** Interface:

// *** Reader:

/*
File reader.

A utility that provides file reading functionality wrapped in a mockable
interface.

To use:

 1. Create an instance with the file path you wish to read:

```go

	load := fileio.NewWReaderithFile("/path/to/file")

```

 2. Check it exists (optional):

```go

	found, err := load.Exists()
	if err != nil {
		panic("File does not exist: " + err.Error())
	}

```

 3. Load your file:

```go

	data, err := load.Load()
	if err != nil {
		panic("Could not load file: " + err.Error())
	}

```

The `Load` method returns the file content as a byte array.
*/
type Reader interface {
	// The file name that we wish to load.
	Filename() string

	// Check whether the file exists.
	//
	// If the file exists, then `true` is returned along with no error.
	//
	// If the file does not exist, then `false` is returned along with
	// no error.
	//
	// If the file exists but is not a regular file, then false is
	// returned along with `ErrNotRegular`.
	//
	// If we are following symbolic links and the file exists and is
	// a symbolic link then it is resolved to the symbolic link's target
	// and, if that exists, `true` is returned along with no error.
	Exists() (bool, error)

	// Check if the file is a symbolic link.
	IsSymlink() (bool, error)

	// Read the file and return a byte array of its content.
	//
	// If `MaxReadBytes` in the options is zero then the number of bytes
	// read will be at most `math.MaxInt`.
	//
	// If `MaxReadBytes` in the options is higher than zero then the
	// number of bytes read will be at most `MaxReadBytes`.
	//
	// `MaxReadBytes` can never be negative and can never exceed
	// `math.MaxInt`.
	Load() ([]byte, error)

	// Open the file and return an `io.ReadCloser`.
	//
	// `MaxReadBytes` is ignored.
	Open() (io.ReadCloser, error)

	// Open the file and stream its contents to the specified writer.
	//
	// If `limit` is zero then the entire contents shall be copied.
	//
	// If `limit` is greater than zero, then at most `limit` bytes will
	// be copied.
	CopyTo(writer io.Writer, limit int64) (int64, error)

	// Stream chunks of up to `chunkSize` bytes.
	//
	// `bufSize` will utilise a readahead buffer of the given size.
	//
	// If `limit` is zero then the entirety of the content will be
	// streamed.
	//
	// if `chunkSize` is zero or lower then a default chunk size of
	// 64 * 1024 shall be used.
	Stream(ctx context.Context, chunkSize, bufSize int, limit int64) StreamResult
}

// *** Writer:

/*
File writer.

A utility that provides file writing functionality wrapped in a mockable
interface.

To use:

 1. Create an instance with the file path you wish to write:

```go

	ctx := context.TODO()
	writer := fileio.NewAppendWriter(
		ctx,
		"/path/to/file",
		fileio.WriteOptions{ ... },
	)

```

 2. Write to your file:

```go

	err := writer.Write(someData)
	if err != nil {
		panic("Could not write to file: " + err.Error())
	}

```

 3. Sync and close your file:

```go

	_ = writer.Sync();
	_ = writer.Close()

```

Add error handling to taste.
*/
type Writer interface {
	io.WriteCloser

	// Perform a synchronisation.
	Sync() error

	// Number of bytes written.
	BytesWritten() int64

	// Name of the file to which we are writing.
	Name() string

	// Abort file writing.
	Abort() error
}

// *** Files:

type Files interface {
	// Open a file for writing.
	OpenWriter(context.Context, string, WriteOptions) (Writer, error)

	// Write to a file.
	WriteFile(context.Context, string, []byte, WriteOptions) error

	// Append data to a file.
	AppendFile(context.Context, string, io.Reader, WriteOptions) (int64, error)

	// Remove a file.
	Remove(string) error

	// Rename a file.
	Rename(string, string) error

	// Create directory.
	//
	// If any directory in the path does not exist, then it will be
	// created.
	MkdirAll(string, os.FileMode) error
}

// * reader.go ends here.
