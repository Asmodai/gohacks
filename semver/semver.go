/*
 * semver.go --- Nasty semantic version hack.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type SemVer struct {
	Major  int
	Minor  int
	Patch  int
	Commit string
}

func MakeSemVer(info string) (*SemVer, error) {
	semver := &SemVer{}
	err := semver.FromString(info)
	if err != nil {
		semver = nil
	}

	return semver, err
}

func NewSemVer() *SemVer {
	return &SemVer{}
}

// Convert numeric version to components
func (s *SemVer) FromString(info string) error {
	var str string = info
	var commit string = "<local>"

	if strings.Contains(info, ":") {
		arr := strings.Split(info, ":")

		if len(arr) != 2 {
			return errors.New(fmt.Sprintf("Invalid version string '%s'", info))
		}

		str = arr[0]
		commit = arr[1]
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		return err
	}

	nmaj := i / 10000000
	nmin := (i / 10000) - (nmaj * 1000)
	npatch := i - ((nmaj * 10000000) + (nmin * 10000))

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
	return ((s.Major * 10000000) +
		(s.Minor * 10000) +
		s.Patch)
}

/* semver.go ends here. */
