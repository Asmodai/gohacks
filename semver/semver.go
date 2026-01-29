// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// semver.go --- Nasty semantic version hack.
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
	"gitlab.com/tozd/go/errors"

	"fmt"
	"strconv"
	"strings"
)

// * Constants:

const (
	magicMajor        = 10000000
	magicMinor        = 10000
	magicMajorToMinor = 1000
)

// * Variables:

var (
	ErrInvalidVersion = errors.Base("invalid version")
)

// * Code:

// ** Types:

// Semantic version structure.
type SemVer struct {
	Major  int    // Major version number.
	Minor  int    // Minor version number.
	Patch  int    // Patch number.
	Commit string // VCS commit identifier.
}

// * Methods:

// Convert numeric version to components.
func (s *SemVer) FromString(info string) error {
	var (
		str    = info
		commit = "<local>"
	)

	if strings.Contains(info, ":") {
		arr := strings.Split(info, ":")

		//nolint:mnd
		if len(arr) != 2 {
			return errors.WithMessagef(ErrInvalidVersion, "%s", info)
		}

		str = arr[0]
		commit = arr[1]
	}

	iver, err := strconv.Atoi(str)
	if err != nil {
		return errors.WithStack(err)
	}

	nmaj := iver / magicMajor
	nmin := (iver / magicMinor) - (nmaj * magicMajorToMinor)
	npatch := iver - ((nmaj * magicMajor) + (nmin * magicMinor))

	s.Major = nmaj
	s.Minor = nmin
	s.Patch = npatch
	s.Commit = commit

	return nil
}

// Return a string representation of the semantic version.
func (s *SemVer) String() string {
	return fmt.Sprintf(
		"%d.%d.%d",
		s.Major,
		s.Minor,
		s.Patch,
	)
}

// Return an integer version.
func (s *SemVer) Version() int {
	return ((s.Major * magicMajor) +
		(s.Minor * magicMinor) +
		s.Patch)
}

// ** Functions:

// Make a new semantic version from the given string.
func MakeSemVer(info string) (*SemVer, error) {
	semver := &SemVer{}

	if err := semver.FromString(info); err != nil {
		return nil, err
	}

	return semver, nil
}

// Create a new empty semantic version object.
func NewSemVer() *SemVer {
	return &SemVer{}
}

// * semver.go ends here.
