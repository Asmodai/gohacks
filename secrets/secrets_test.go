// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// secrets_test.go --- Secrets tests.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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

package secrets

import (
	"github.com/Asmodai/gohacks/utils"

	"gitlab.com/tozd/go/errors"

	"os"
	"runtime"
	"testing"
)

const (
	TestDataPath  = "/../testing/"
	TestDataFile  = "secrets"
	TestZeroFile  = "zerolength"
	TestDataValue = "testing\n"
)

func GenerateDirectory() string {
	switch runtime.GOOS {
	case "windows":
		return utils.GetEnv("SYSTEMROOT", "C:\\Windows")

	default:
		return "/tmp"
	}
}

func GetPath() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}

	return path, nil
}

func TestSecrets(t *testing.T) {
	base, err := GetPath()
	if err != nil {
		t.Errorf("Could not get base directory: %s", err.Error())

		return
	}

	// Create a new secret object.
	secret := Make(TestDataFile)

	// Fix up the path
	secret.setPath(base + TestDataPath + TestDataFile)

	t.Run("Has path", func(t *testing.T) {
		switch secret.Path() {
		case "":
			t.Error("Path has zero length.")

		case base + TestDataPath + TestDataFile:
		// Everything is good.

		default:
			t.Errorf("Unexpected path: %#v", secret.Path())
		}
	})

	t.Run("Can probe", func(t *testing.T) {
		if err := secret.Probe(); err != nil {
			t.Errorf("Probe failed: %s", err.Error())
		}
	})

	t.Run("Has value", func(t *testing.T) {
		switch secret.Value() {
		case "":
			t.Error("Value has zero length.")

		case TestDataValue:
			// Everything is good.

		default:
			t.Errorf("Unexpected value: %#v", secret.Value())
		}
	})
}

func TestFailureConditions(t *testing.T) {
	base, err := GetPath()
	if err != nil {
		t.Errorf("Could not get base directory: %s", err.Error())

		return
	}

	// Create a new secret objects.
	secret := New()
	zerolength := New()

	// Fix up the path
	zerolength.setPath(base + TestDataPath + TestZeroFile)

	t.Run("Error for no path set", func(t *testing.T) {
		secret.setPath("")

		err := secret.Probe()

		if err == nil {
			t.Error("No error returned.")
		}

		if !errors.Is(err, ErrNoPathSet) {
			t.Errorf("Error is invalid for condition: %#v", err)
		}
	})

	t.Run("Error on bad path", func(t *testing.T) {
		secret.setPath("/totally/broken/path")

		err := secret.Probe()

		if !errors.Is(err, ErrFileNotFound) {
			t.Errorf("Error is invalid for condition: %#v", err)
		}
	})

	t.Run("Error on directory", func(t *testing.T) {
		secret.setPath(GenerateDirectory())

		err := secret.Probe()

		if !errors.Is(err, ErrNotAFile) {
			t.Errorf("Error is invalid for condition: %#v", err)
		}
	})

	t.Run("Error on zero-length file", func(t *testing.T) {
		err := zerolength.Probe()

		if !errors.Is(err, ErrZeroLengthFile) {
			t.Errorf("Error is invalid for condition: %#v", err)
		}
	})
}

// secrets_test.go ends here.
