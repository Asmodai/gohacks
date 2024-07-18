// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// filesystem_test.go --- Filesystem tests.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package utils

import (
	"fmt"
	"runtime"
	"testing"
)

func InvalidTestFile() string {
	if runtime.GOOS == "windows" {
		return "I:\\Wibbles.dll"
	}

	return "/invalid/file"
}

func ValidTestDirectory() string {
	switch runtime.GOOS {
	case "windows":
		return GetEnv("SYSTEMROOT", "C:\\Windows")

	default:
		return "/tmp"
	}
}

func ValidTestFile() (string, string) {
	switch runtime.GOOS {
	case "aix":
		fallthrough
	case "darwin":
		fallthrough
	case "dragonfly":
		fallthrough
	case "freebsd":
		fallthrough
	case "illumos":
		fallthrough
	case "linux":
		fallthrough
	case "netbsd":
		fallthrough
	case "openbsd":
		fallthrough
	case "solaris":
		return "/bin/ls", "ls"

	case "windows":
		// `GetEnv` comes from utils/env.go.
		prog := GetEnv("SYSTEMROOT", "C:\\Windows") + "\\explorer.exe"
		return prog, "explorer.exe"

	default:
		panic(fmt.Sprintf("unhandled os '%s'", runtime.GOOS))
	}
}

// Test a valid file.
//
// n.b. This depends on the host OS.
func TestValidFile(t *testing.T) {
	file, name := ValidTestFile()
	fsys := NewFilesystem(file)

	t.Run("Exists?", func(t *testing.T) {
		if fsys.Exists() == false {
			t.Errorf("File %s seems to not exist", file)
		}
	})

	t.Run("Name?", func(t *testing.T) {
		if fsys.Name() != name {
			t.Errorf("Name does not match.  %s != %s", fsys.Name(), name)
		}
	})

	t.Run("Size?", func(t *testing.T) {
		if fsys.Size() == 0 {
			t.Error("Zero size")
		}
	})

	t.Run("Mode?", func(t *testing.T) {
		if fsys.Mode() == 0 {
			t.Error("Zero mode")
		}
	})

	t.Run("Not directory?", func(t *testing.T) {
		if fsys.IsDirectory() == true {
			t.Errorf("File %s is a directory", file)
		}
	})

	t.Run("File?", func(t *testing.T) {
		if fsys.IsFile() == false {
			t.Errorf("File %s is not a file", file)
		}
	})

	t.Run("Executable?", func(t *testing.T) {
		if fsys.IsExecutable() == false {
			t.Errorf("File %s is not executable", file)
		}
	})
}

// Test a valid directory
func TestValidDirectory(t *testing.T) {
	dir := ValidTestDirectory()
	fsys := NewFilesystem(dir)

	t.Run("Exists?", func(t *testing.T) {
		if fsys.Exists() == false {
			t.Errorf("Directory %s does not exist", dir)
		}
	})

	t.Run("Directory?", func(t *testing.T) {
		if fsys.IsDirectory() != true {
			t.Error("Apparently not a directory")
		}
	})

	t.Run("File?", func(t *testing.T) {
		if fsys.IsFile() != false {
			t.Error("Apparently a file")
		}
	})
}

// Test a invalid file.
func TestInvalidFile(t *testing.T) {
	file := InvalidTestFile()
	fsys := NewFilesystem(file)

	t.Run("Exists?", func(t *testing.T) {
		if fsys.Exists() == true {
			t.Errorf("File %s seems to exist", file)
		}
	})

	t.Run("Name?", func(t *testing.T) {
		if fsys.Name() != "" {
			t.Errorf("Got an unexpected name: %#v", fsys.Name())
		}
	})

	t.Run("Size?", func(t *testing.T) {
		if fsys.Size() != 0 {
			t.Errorf("Non-zero size: %v", fsys.Size())
		}
	})

	t.Run("Mode?", func(t *testing.T) {
		if fsys.Mode() != 0 {
			t.Errorf("Non-zero mode: %v", fsys.Mode())
		}
	})

	t.Run("Directory?", func(t *testing.T) {
		if fsys.IsDirectory() != false {
			t.Errorf("Unexpected result: %v", fsys.IsDirectory())
		}
	})

	t.Run("File?", func(t *testing.T) {
		if fsys.IsFile() != false {
			t.Errorf("Unexpected result: %v", fsys.IsFile())
		}
	})

	t.Run("Executable?", func(t *testing.T) {
		if fsys.IsExecutable() != false {
			t.Errorf("Unexpected result: %v", fsys.IsExecutable())
		}
	})
}

// filesystem_test.go ends here.
