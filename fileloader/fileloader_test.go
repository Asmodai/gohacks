// -*- Mode: Go -*-
//
// fileloader_test.go --- File loader tests.
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

package fileloader

// * Imports:

import (
	"testing"
)

// * Constants:

const (
	TestingFile    string = "test.txt"
	TestingBadFile string = "notexists.txt"
	FileContents   string = "This is a test file for FileLoader.\n\n"
)

// * Code:

func TestValidFileLoader(t *testing.T) {
	var inst FileLoader

	t.Run("Construction", func(t *testing.T) {
		inst = NewWithFile(TestingFile)

		if inst == nil {
			t.Error("Instance is invalid")
			return
		}

		if inst.Filename() != TestingFile {
			t.Errorf(
				"Unexpected filename: %#v != %#v",
				inst.Filename(),
				TestingFile,
			)
			return
		}
	})

	t.Run("Exists", func(t *testing.T) {
		ok, err := inst.Exists()

		if err != nil {
			t.Errorf("Returned error: %#v", err)
			return
		}

		if !ok {
			t.Error("No error, but also not ok.")
			return
		}
	})

	t.Run("Load", func(t *testing.T) {
		data, err := inst.Load()

		if err != nil {
			t.Errorf("Returned error: %#v", err)
			return
		}

		datastr := string(data)
		if datastr != FileContents {
			t.Errorf(
				"Unexpected data: %#v != %#v",
				datastr,
				FileContents,
			)
			return
		}
	})
}

func TestInvalidFileLoader(t *testing.T) {
	var inst FileLoader

	inst = NewWithFile(TestingBadFile)

	t.Run("Exists", func(t *testing.T) {
		ok, err := inst.Exists()

		if err == nil {
			t.Error("Did not return an error")
			return
		}

		if ok {
			t.Error("Non-existent file apparently exists")
			return
		}
	})

	t.Run("Load", func(t *testing.T) {
		_, err := inst.Load()

		if err == nil {
			t.Error("Did not return an error")
			return
		}
	})
}

// * fileloader_test.go ends here.
