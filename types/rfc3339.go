// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// rfc3339.go --- RFC 3339 time format.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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

package types

// * Imports:

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/btubbs/datetime"
	"gitlab.com/tozd/go/errors"
	"gopkg.in/yaml.v3"
)

// * Constants:

const (
	// Type used when reporting an argument's type.
	//
	// This is used for CLI flag processing.
	cliRFC3339Type string = "rfc3339"

	// Minimum length of an RFC3339 timestamp string.
	minRFC3339Length = 5
)

// * Variables:

var (
	// Error condition that signals an invalid RFC3339 timestamp  of some
	// kind.
	//
	// This error is usually wrapped around a descriptive message string.
	ErrInvalidRFC3339 error = errors.Base("invalid RFC3339 timestamp")

	// Error condition that signals that an RFC3339 timestamp is not a
	// string format.
	//
	// This error is used by `Set` as well as JSON and YAML methods.
	ErrRFC3339NotString error = errors.Base("RFC3339 timestamp  must be a string")
)

// * Code:

// ** Types

// RFC 3339 time type.
//
//nolint:recvcheck
type RFC3339 time.Time

// ** Methods:

// *** Coercion methods:

// Coerce an RFC3339 time value to a string.
func (obj RFC3339) String() string {
	tval := time.Time(obj)

	return tval.UTC().Format(time.RFC3339)
}

// Coerce an RFC3339 time value to a `time.Time` value.
func (obj RFC3339) Time() time.Time {
	return time.Time(obj).UTC()
}

// *** Time methods:

// Return the UTC time for the timestamp.
//
// RFC3339 timestamps are always UTC internally, so `UTC` is provided as a
// courtesy.
func (obj RFC3339) UTC() time.Time {
	// NOTE: `Time` returns UTC.
	return obj.Time()
}

// Return the Unix time for the timestamp.
func (obj RFC3339) Unix() int64 {
	// Unix time is based on UTC.
	return obj.Time().Unix()
}

// Return a string that can be used in MySQL queries.
func (obj RFC3339) MySQL() string {
	return TimeToMySQL(obj.Time())
}

// Format the timestamp with the given format.
func (obj RFC3339) Format(format string) string {
	return obj.UTC().Format(format)
}

// *** Shim methods:

// Is the timestamp a zero value?
func (obj RFC3339) IsZero() bool {
	return obj.Time().IsZero()
}

// Does the timestamp correspond to a time where DST is in effect?
func (obj RFC3339) IsDST() bool {
	return obj.Time().IsDST()
}

// Is the given time before the time in the timestamp?
func (obj RFC3339) Before(t time.Time) bool {
	return obj.Time().Before(t)
}

// Is the given time after the time in the timestamp?
func (obj RFC3339) After(t time.Time) bool {
	return obj.Time().After(t)
}

// Is the given time equal to the time in the timestamp?
func (obj RFC3339) Equal(t time.Time) bool {
	return obj.Time().Equal(t)
}

// ** Manipulation methods:

// Add a `time.Duration` value to the timestamp, returning a new timestamp.
func (obj RFC3339) Add(d time.Duration) RFC3339 {
	return RFC3339(obj.Time().Add(d))
}

// Subtract a `time.Time` value from the timestamp, returning a
// `time.Duration`.
func (obj RFC3339) Sub(t time.Time) time.Duration {
	return obj.Time().Sub(t)
}

// *** Mutation methods:

// Set the RFC3339 timestamp to that of the given string.
//
// If the string value fails to parse, then `ErrInvalidRFC3339` is returned.
func (obj *RFC3339) Set(str string) error {
	parsed, err := ParseRFC3339(str)
	if err != nil {
		return errors.WithMessagef(
			ErrInvalidRFC3339,
			"invalid RFC3339 timestamp %q: %s",
			str,
			err.Error(),
		)
	}

	*obj = parsed

	return nil
}

// *** CLI methods:

// Return the data type name for CLI flag parsing purposes.
func (obj RFC3339) Type() string {
	return cliRFC3339Type
}

// *** JSON Methods:

// JSON unmarshalling method.
func (obj *RFC3339) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)

	return obj.Set(str)
}

// JSON marshalling method.
func (obj RFC3339) MarshalJSON() ([]byte, error) {
	res, err := json.Marshal(obj.String())

	if err != nil {
		err = errors.WithStack(err)
	}

	return res, err
}

// *** YAML methods:

// YAML unmarshalling method.
func (obj *RFC3339) UnmarshalYAML(value *yaml.Node) error {
	var str string

	if err := value.Decode(&str); err != nil {
		return errors.WrapWith(err, ErrRFC3339NotString)
	}

	return obj.Set(str)
}

// YAML marshalling method.
func (obj RFC3339) MarshalYAML() (any, error) {
	return obj.String(), nil
}

// Set the RFC3339 value from a YAML value.
//
// If the passed YAML value is not a string, then `ErrRFC3339NotString` is
// returned.
//
// Will also return any error condition from the `Set` method.
func (obj *RFC3339) SetYAML(value any) error {
	str, ok := value.(string)
	if !ok {
		return errors.WithMessagef(
			ErrRFC3339NotString,
			"expected string, got %T",
			value,
		)
	}

	str = strings.TrimRight(str, "\r\n")

	return obj.Set(str)
}

// ** Functions:

// *** Time:

// Return the current timezone for the host.
func CurrentZone() (string, int) {
	return time.Now().Zone()
}

// Convert a Unix timestamp to an RFC3339 timestamp.
func RFC3339FromUnix(unix int64) RFC3339 {
	return RFC3339(time.Unix(unix, 0).UTC())
}

// *** Conversion:

// Convert a `time.Time` value to a MySQL timestamp for queries.
func TimeToMySQL(val time.Time) string {
	return val.Format("2006-01-02 15:04:05")
}

// *** Parsing:

// Parse the given string for an RFC3339 timestamp.
//
// If the timestamp is not a valid RFC3339 timestamp, then `ErrInvalidRFC3339`
// is returned.
func ParseRFC3339(data string) (RFC3339, error) {
	if len(data) < minRFC3339Length {
		return RFC3339{}, errors.WithMessagef(
			ErrInvalidRFC3339,
			"timestamp too short: %q",
			data,
		)
	}

	tzChar := strings.ToUpper(data[len(data)-1:])
	tzOff := data[len(data)-5:]

	if tzChar == "Z" || tzOff == "+" || tzOff == "-" {
		rval, err := time.Parse(time.RFC3339, data)
		if err != nil {
			return RFC3339{}, errors.WithStack(err)
		}

		return RFC3339(rval), nil
	}

	// We're not bothering with error handling here.  This code is
	// meant to catch cases where people send us timestamps from the
	// future.  We don't care about the future.
	temp, _ := datetime.Parse(data, time.UTC)
	if temp.After(time.Now()) {
		return ParseRFC3339(time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	}

	// `time.Local` exists here because of some cases where there is
	// no offset identifier, and the sending system is sending
	// timestamps in local time.  Because people are dumb.
	//
	//nolint:gosmopolitan
	rval, err := datetime.Parse(data, time.Local)
	if err != nil {
		return RFC3339{}, errors.WithStack(err)
	}

	return RFC3339(rval), nil
}

// * rfc3339.go ends here.
