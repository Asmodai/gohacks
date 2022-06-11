/*
 * rfc3339.go --- RFC3339 support.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
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

package rfc3339

import (
	"github.com/btubbs/datetime"

	"strings"
	"time"
)

// An RFC3339 object.
type JsonRFC3339 time.Time

// Unmarshal an RFC3339 timestamp from JSON.
func (j *JsonRFC3339) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	t, err := RFC3339Parse(s)
	if err != nil {
		return err
	}

	*j = JsonRFC3339(t)

	return nil
}

// Marshal an RFC3339 object to JSON.
func (j JsonRFC3339) MarshalJSON() ([]byte, error) {
	return []byte("\"" + j.Format(time.RFC3339) + "\""), nil
}

// Format an RFC3339 object as a string.
func (j JsonRFC3339) Format(s string) string {
	return j.Time().Format(s)
}

// Convert an RFC3339 time to UTC.
func (j JsonRFC3339) UTC() time.Time {
	return j.Time().UTC()
}

// Convert an RFC3339 time to Unix time.
func (j JsonRFC3339) Unix() int64 {
	return j.Time().Unix()
}

// convert an RFC3339 time to time.Time.
func (j JsonRFC3339) Time() time.Time {
	return time.Time(j)
}

// Convert an RFC3339 time to a MySQL timestamp.
func (j JsonRFC3339) MySQL() string {
	return TimeToMySQL(j.Time())
}

// Parse a string to an RFC3339 timestamp.
func RFC3339Parse(data string) (time.Time, error) {
	tzchar := strings.ToUpper(data[len(data)-1:])
	tzoff := data[len(data)-5:]

	if tzchar == "Z" || tzoff == "+" || tzoff == "-" {
		return time.Parse(time.RFC3339, data)
	}

	temp, _ := datetime.Parse(data, time.UTC)
	if temp.After(time.Now()) {
		return RFC3339Parse(time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	}

	return datetime.Parse(data, time.Local)
}

// Return the current time zone.
func CurrentZone() (string, int) {
	return time.Now().Zone()
}

// Convert a time to a MySQL string.
func TimeToMySQL(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// Convert a Unix `time_t` value to a time.
func FromUnix(t int64) time.Time {
	unix := time.Unix(t, 0)

	return unix.UTC()
}

/* rfc3339.go ends here. */
