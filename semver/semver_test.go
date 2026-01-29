// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// semver_test.go --- SemVer test file.
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

// * Comments:

//
//
//

// * Package:

package semver

// * Imports:

import (
	"testing"
)

// * Code:

// ** Tests:

func TestBuild(t *testing.T) {
	t.Run("Parses proper version string", func(t *testing.T) {
		ok := true
		str := "10020003:commit"

		// Use the `MakeSemVer` utility function here.
		semver, err := MakeSemVer(str)
		if err != nil {
			t.Errorf("Got an error: %s", err.Error())
			return
		}

		if semver.Major != 1 || semver.Minor != 2 || semver.Patch != 3 {
			ok = false
		}

		if semver.Commit != "commit" {
			ok = false
		}

		if !ok {
			t.Errorf("Unexpected result: %v", semver)
		}
	})

	t.Run("Parses proper version string without commit hash", func(t *testing.T) {
		var semver = NewSemVer()

		ok := true
		str := "10020003"

		// Use the `FromString` method here.
		err := semver.FromString(str)
		if err != nil {
			t.Errorf("Got an error: %s", err.Error())
			return
		}

		if semver.Major != 1 || semver.Minor != 2 || semver.Patch != 3 {
			ok = false
		}

		if semver.Commit != "<local>" {
			ok = false
		}

		if !ok {
			t.Errorf("Unexpected result: %v", semver)
		}
	})

	t.Run("Returns an error for invalid string", func(t *testing.T) {
		str := "invalid"

		_, err := MakeSemVer(str)
		if err == nil {
			t.Error("Expected an error, but got none.")
		}
	})

	t.Run("Returns an error for string with invalid commit", func(t *testing.T) {
		str := "1:invalid:wheee"

		_, err := MakeSemVer(str)
		if err == nil {
			t.Error("Expected an error, but got none.")
		}
	})
}

func TestStringify(t *testing.T) {
	var semver = NewSemVer()

	t.Run("Returns a version string", func(t *testing.T) {
		str := "10020003:commit"

		err := semver.FromString(str)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
			return
		}

		if semver.String() != "1.2.3" {
			t.Errorf("Unexpected result: %s", semver.String())
		}
	})
}

func TestVersionNumber(t *testing.T) {
	var semver = NewSemVer()

	t.Run("Returns a version numeric", func(t *testing.T) {
		str := "10020003:commit"

		err := semver.FromString(str)
		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
			return
		}

		if semver.Version() != 10020003 {
			t.Errorf("Unexpected numeric: %v", semver.Version())
		}
	})
}

// * semver_test.go ends here.
