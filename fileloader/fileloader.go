// -*- Mode: Go -*-
//
// fileloader.go --- File loader utility.
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
// mock:yes

// * Comments:
//
//

// * Package:

package fileloader

// * Imports:

import (
	"gitlab.com/tozd/go/errors"

	"io"
	"os"
)

// * Code:

// ** Types:

// *** Interface:

/*
File loader.

A utility that provides file opening functionality wrapped in a mockable
interface.

To use:

 1. Create an instance with the file path you wish to load:

```go

	load := fileloader.NewWithFile("/path/to/file")

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
*/
type FileLoader interface {
	// The file path that we wish to load.
	Filename() string

	// Check if the file exists.
	Exists() (bool, error)

	// Read the file and return a byte array of its content.
	Load() ([]byte, error)
}

// *** Implementation:

type fileloader struct {
	filename string
}

// **** Methods:

// The file path that we wish to load.
func (obj *fileloader) Filename() string {
	return obj.filename
}

// Check if the file exists.
func (obj *fileloader) Exists() (bool, error) {
	_, err := os.Stat(obj.filename)
	if err == nil {
		return true, nil
	}

	return false, errors.WithStack(err)
}

// Read the file and return a byte array of its content.
func (obj *fileloader) Load() ([]byte, error) {
	_, err := obj.Exists()
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}

	file, err := os.Open(obj.filename)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}

	return bytes, nil
}

// ** Functions:

// Create a new FileLoader object with the given file name.
func NewWithFile(filename string) FileLoader {
	return &fileloader{
		filename: filename,
	}
}

// * fileloader.go ends here.
