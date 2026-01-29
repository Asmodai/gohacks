// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// duration.go --- Enhanced time duration type.
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
	"fmt"
	"strings"
	"time"

	"gitlab.com/tozd/go/errors"
	"gopkg.in/yaml.v3"
)

// * Constants:

const (
	// Type used when reporting an argument's type.
	//
	// This is used for CLI flag processing.
	cliDurationType string = "duration"
)

// * Variables:

var (
	// Error condition that signals an invalid time duration of some kind.
	//
	// This error is usually wrapped around a descriptive message string.
	ErrInvalidDuration error = errors.Base("invalid time duration")

	// Error condition that signals that a duration is not a string value.
	//
	// This error is used by `Set` as well as JSON and YAML methods.
	ErrDurationNotString error = errors.Base("duration must be a string")

	// Error condition that signals that a duration is out of bounds.
	//
	// This is used by `Validate`.
	ErrOutOfBounds error = errors.Base("duration out of bounds")
)

// * Code:

// ** Types:

// Enhanced time duration type.
//
//nolint:recvcheck
type Duration time.Duration

// ** Methods:

// *** Coercion methods:

// Coerce a duration to a string value.
func (obj Duration) String() string {
	tval := time.Duration(obj)

	return tval.String()
}

// Coerce a duration to a `time.Duration` value.
func (obj Duration) Duration() time.Duration {
	return time.Duration(obj)
}

// *** Mutation methods:

// Set the duration to that of the given string.
//
// This method uses `time.ParseDuration`, so any string that `time` understands
// may be used.
//
// If the string value fails parsing, then `ErrInvalidDuration` is returned.
func (obj *Duration) Set(str string) error {
	dur, err := time.ParseDuration(str)
	if err != nil {
		return errors.WithMessagef(
			ErrInvalidDuration,
			"invalid duration %q: %s",
			str,
			err.Error(),
		)
	}

	*obj = Duration(dur)

	return nil
}

// *** CLI methods:

// Return the data type name for CLI flag parsing purposes.
func (obj Duration) Type() string {
	return cliDurationType
}

// Validate a duration.
//
// This ensures a duration is within a given range.
//
// If validation fails, then `ErrOutOfBounds` is returned.
func (obj Duration) Validate(minDuration, maxDuration time.Duration) error {
	dur := time.Duration(obj)

	if dur < minDuration || dur > maxDuration {
		return errors.WithMessagef(
			ErrOutOfBounds,
			"duration %s is out of bounds [%s, %s]",
			dur,
			minDuration,
			maxDuration,
		)
	}

	return nil
}

// *** JSON methods:

// JSON unmarshalling method.
func (obj *Duration) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)

	return obj.Set(str)
}

// JSON marshalling method.
func (obj Duration) MarshalJSON() ([]byte, error) {
	res, err := json.Marshal(obj.String())

	if err != nil {
		err = errors.WithStack(err)
	}

	return res, err
}

// *** YAML methods:

// YAML unmarshalling method.
func (obj *Duration) UnmarshalYAML(value *yaml.Node) error {
	var str string

	if err := value.Decode(&str); err != nil {
		return errors.WrapWith(err, ErrDurationNotString)
	}

	return obj.Set(str)
}

// YAML marshalling method.
func (obj Duration) MarshalYAML() (any, error) {
	return obj.String(), nil
}

// Set the duration value from a YAML value.
//
// If the passed YAML value is not a string, then `ErrDurationNotString` is
// returned.
//
// Will also return any error condition from the `Set` method.
func (obj *Duration) SetYAML(value any) error {
	str, ok := value.(string)
	if !ok {
		return errors.WithMessagef(
			ErrDurationNotString,
			"expected string, got %T",
			value,
		)
	}

	str = strings.TrimRight(str, "\r\n")

	return obj.Set(str)
}

// ** Functions:

func NewFromDuration(duration string) (Duration, error) {
	obj := Duration(0)

	err := obj.Set(duration)
	if err != nil {
		return obj, errors.WithStack(err)
	}

	return obj, nil
}

// Format a time duration in pretty format.
//
// Example, a duration of 72 minutes becomes "1 hour(s), 12 minute(s)".
func PrettyFormat(dur Duration) string {
	val := dur.Duration().Round(time.Minute)

	hour := val / time.Hour
	val -= hour * time.Hour //nolint:durationcheck
	minute := val / time.Minute

	return fmt.Sprintf("%0d hour(s), %0d minute(s)", hour, minute)
}

// * duration.go ends here.
