// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// semver.go --- Nasty semantic version hack.
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

package semver

import (
	"gitlab.com/tozd/go/errors"

	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidVersion = errors.Base("invalid version")
)

const (
	MAGICMAJOR        = 10000000
	MAGICMINOR        = 10000
	MAGICMAJORTOMINOR = 1000
)

type SemVer struct {
	Major  int
	Minor  int
	Patch  int
	Commit string
}

func MakeSemVer(info string) (*SemVer, error) {
	semver := &SemVer{}

	if err := semver.FromString(info); err != nil {
		return nil, err
	}

	return semver, nil
}

func NewSemVer() *SemVer {
	return &SemVer{}
}

// Convert numeric version to components.
func (s *SemVer) FromString(info string) error {
	var (
		str    = info
		commit = "<local>"
	)

	if strings.Contains(info, ":") {
		arr := strings.Split(info, ":")

		//nolint:gomnd
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

	nmaj := iver / MAGICMAJOR
	nmin := (iver / MAGICMINOR) - (nmaj * MAGICMAJORTOMINOR)
	npatch := iver - ((nmaj * MAGICMAJOR) + (nmin * MAGICMINOR))

	s.Major = nmaj
	s.Minor = nmin
	s.Patch = npatch
	s.Commit = commit

	return nil
}

func (s *SemVer) String() string {
	return fmt.Sprintf(
		"%d.%d.%d",
		s.Major,
		s.Minor,
		s.Patch,
	)
}

func (s *SemVer) Version() int {
	return ((s.Major * MAGICMAJOR) +
		(s.Minor * MAGICMINOR) +
		s.Patch)
}

// semver.go ends here.
