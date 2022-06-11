/*
 * semver.go --- Nasty semantic version hack.
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
