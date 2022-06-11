/*
 * semver_test.go --- SemVer test file.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package semver

import (
	"testing"
)

func TestBuild(t *testing.T) {
	var semver *SemVer = NewSemVer()

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

		err := semver.FromString(str)
		if err == nil {
			t.Error("Expected an error, but got none.")
		}
	})

	t.Run("Returns an error for string with invalid commit", func(t *testing.T) {
		str := "1:invalid:wheee"

		err := semver.FromString(str)
		if err == nil {
			t.Error("Expected an error, but got none.")
		}
	})
}

func TestStringify(t *testing.T) {
	var semver *SemVer = NewSemVer()

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
	var semver *SemVer = NewSemVer()

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

/* semver_test.go ends here. */
