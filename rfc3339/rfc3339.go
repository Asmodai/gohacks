/*
 * rfc3339.go --- RFC3339 support.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
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
