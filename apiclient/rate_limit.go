// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// rate_limit.go --- Rate-limiting information.
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

// * Package:

package apiclient

// * Imports:

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// * Constants:

const (
	ReasonNone int = iota
	ReasonRetryAfter
	ReasonRateLimitExceeded
	ReasonHTTPStatusCode
	ReasonMaximum

	ReasonInvalid string = "invalid reason"
)

// * Variables:

var (
	//nolint:gochecknoglobals
	throttleReasons = []string{
		"No reason.",
		"'Retry-After' header present",
		"Rate limit exceeded",
		"HTTP status code infers throttling",
	}

	//nolint:gochecknoglobals
	limitKeys = []string{
		"X-Ratelimit-Limit",
		"X-Rate-Limit-Limit",
		"Ratelimit-Limit",
		"Rate-Limit-Limit",
	}

	//nolint:gochecknoglobals
	remainingKeys = []string{
		"X-Ratelimit-Remaining",
		"X-Rate-Limit-Remaining",
		"Ratelimit-Remaining",
		"Rate-Limit-Remaining",
	}

	//nolint:gochecknoglobals
	resetKeys = []string{
		"X-Ratelimit-Reset",
		"X-Rate-Limit-Reset",
		"Ratelimit-Reset",
		"Rate-Limit-Reset",
		"Retry-After",
	}
)

// * Code:

// ** Types:

// Rate limit information.
type RateLimitInfo struct {
	Limit      int           // Total limit count.
	Remaining  int           // Number of calls remaining.
	ResetTime  time.Time     // Limit reset time.
	RetryAfter time.Duration // Retry delay.
	ReasonCode int           // Inferred throttle reason.
}

// ** Methods:

// Has the rate limit been reached?
func (obj RateLimitInfo) IsRateLimited() bool {
	// We set the reason code when we have a HTTP status code that infers
	// throttling of some kind, so check for that here.
	if obj.ReasonCode > ReasonNone && obj.ReasonCode < ReasonMaximum {
		return true
	}

	// The way this works with the rate-limit headers is that if you are
	// limited, then `RateLimit-Remaining` is 0, and you can use the
	// time in `RateLimit-Reset` as your "try again after."
	//
	// There is zero point in doing any amazing logic here, because
	// the simple matter is that because this code uses the non-standard
	// rate-limit headers *and* the standard `Retry-After` header, it can
	// get into a situation where it always claims to be limited.
	return obj.Remaining <= 0
}

// Return the throttling reason (if any) as a string.
func (obj RateLimitInfo) Reason() string {
	if obj.ReasonCode > ReasonMaximum {
		return ReasonInvalid
	}

	return throttleReasons[obj.ReasonCode]
}

// Return the rate limit information in string format.
func (obj RateLimitInfo) String() string {
	return fmt.Sprintf(
		"Limit=%d  Remaining=%d  Reset=%v  RetryAfter=%v  Reason=%v",
		obj.Limit,
		obj.Remaining,
		obj.ResetTime,
		obj.RetryAfter,
		obj.Reason(),
	)
}

// ** Functions:

// Parse an integer value from HTTP headers.
func parseHeaderInt(headers http.Header, keys ...string) (int, bool) {
	for _, key := range keys {
		if val := headers.Get(key); len(val) > 0 {
			result, err := strconv.Atoi(strings.TrimSpace(val))
			if err == nil {
				return result, true
			}
		}
	}

	return 0, false
}

// Parse a time value from HTTP headers.
func parseResetTime(
	headers http.Header,
	keys ...string,
) (time.Time, time.Duration, bool) {
	for _, key := range keys {
		val := headers.Get(key)
		if val == "" {
			continue
		}

		val = strings.TrimSpace(val)

		if secs, err := strconv.Atoi(val); err == nil {
			return time.Now().Add(time.Duration(secs) * time.Second),
				time.Duration(secs) * time.Second,
				true
		}

		// Try Unix timestamp.
		if ts, err := strconv.ParseInt(val, 10, 64); err == nil {
			when := time.Unix(ts, 0)
			dur := time.Until(when)

			return when, dur, true
		}

		// Try HTTP date.
		if date, err := http.ParseTime(val); err == nil {
			dur := time.Until(date)

			return date, dur, true
		}
	}

	return time.Time{}, time.Duration(0), false
}

// Return a new rate limit info object.
func NewRateLimitInfo(code int, headers http.Header) RateLimitInfo {
	inst := RateLimitInfo{}

	if limit, ok := parseHeaderInt(headers, limitKeys...); ok {
		inst.Limit = limit
	}

	if remain, ok := parseHeaderInt(headers, remainingKeys...); ok {
		inst.Remaining = remain
	}

	if when, until, ok := parseResetTime(headers, resetKeys...); ok {
		inst.ResetTime = when
		inst.RetryAfter = until
	}

	// Why `inst.Limit`?  Some API servers will return a Remaining and
	// Limit value of 0 if you've reached a concurrent API call limit.
	if inst.Remaining <= 0 || inst.Limit <= 0 {
		inst.ReasonCode = ReasonRateLimitExceeded
	}

	// Special case for `Retry-After`.
	if val := headers.Get("Retry-After"); len(val) > 0 {
		inst.Limit = 0
		inst.Remaining = 0
		inst.ReasonCode = ReasonRetryAfter
	}

	// Special case for HTTP 429 if there are no other throttle
	// conditions.
	if code == http.StatusTooManyRequests && inst.ReasonCode == ReasonNone {
		inst.ReasonCode = ReasonHTTPStatusCode
	}

	return inst
}

// * rate_limit.go ends here.
