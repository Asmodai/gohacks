// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// fmtduration.go --- Format a Time duration.
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

package utils

import (
	"fmt"
	"time"
)

// Format a time duration in pretty format.
//
// Example, a duration of 72 minutes becomes "1 hour(s), 12 minute(s)".
func FormatDuration(dur time.Duration) string {
	dur = dur.Round(time.Minute)

	// Compute hours, and then subtract from the duration.
	hour := dur / time.Hour
	dur -= hour * time.Hour //nolint:durationcheck

	// Compute minutes.
	minute := dur / time.Minute

	return fmt.Sprintf("%0d hour(s), %0d minute(s)", hour, minute)
}

// fmtduration.go ends here.
