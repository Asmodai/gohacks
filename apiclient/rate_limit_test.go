// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// rate_limit_test.go --- Rate limit tests.
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

package apiclient

// * Imports:

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

// * Constants:

const (
	valueThousand    string = "1000"
	valueThousandInt int    = 1000

	valueHundred    string = "100"
	valueHundredInt int    = 100

	valueZero    string = "0"
	valueZeroInt int    = 0

	valueTimestamp string = "1893456001"

	retryAfter         string        = "120"
	retryAfterDuration time.Duration = time.Duration(120_000_000_000)
)

// * Code:

// ** Tests:

func TestRateLimitInfo(t *testing.T) {
	t.Run("Rate limit OK", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("X-RateLimit-Limit", valueThousand)
		headers.Add("X-RateLimit-Remaining", valueHundred)
		headers.Add("X-RateLimit-Reset", valueTimestamp)

		inst := NewRateLimitInfo(http.StatusOK, headers)

		if inst.Limit != valueThousandInt {
			t.Errorf("Rate limit mismatch: %#v != %#v",
				inst.Limit,
				valueThousandInt)
		}

		if inst.Remaining != valueHundredInt {
			t.Errorf("Remaining mismatch: %#v != %#v",
				inst.Remaining,
				valueHundredInt)
		}

		if inst.IsRateLimited() == true {
			t.Error("Claiming rate limited when not limited.")
		}

		if inst.ReasonCode != ReasonNone {
			t.Errorf("Unexpected reason code: %s", inst.Reason())
		}
	})

	t.Run("Only Retry-After", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("Retry-After", retryAfter)

		inst := NewRateLimitInfo(http.StatusOK, headers)

		if inst.Limit > 0 {
			t.Errorf("Limit is non-zero: %#v", inst.Limit)
		}

		if inst.Remaining > 0 {
			t.Errorf("Remaining is non-zero: %#v", inst.Remaining)
		}

		if inst.IsRateLimited() == false {
			t.Errorf("Claiming to not be rate limited.")
		}

		if inst.ReasonCode != ReasonRetryAfter {
			t.Errorf("Unexpected reason code: %s", inst.Reason())
		}
	})

	t.Run("Rate limit reached", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("X-RateLimit-Limit", valueThousand)
		headers.Add("X-RateLimit-Remaining", valueZero)
		headers.Add("x-RateLimit-Reset", valueTimestamp)

		inst := NewRateLimitInfo(http.StatusOK, headers)

		if inst.Limit != valueThousandInt {
			t.Errorf("Rate limit mismatch: %#v != %#v",
				inst.Limit,
				valueThousandInt)
		}

		if inst.Remaining != valueZeroInt {
			t.Errorf("Remaining mismatch: %#v != %#v",
				inst.Remaining,
				valueZeroInt)
		}

		if inst.IsRateLimited() == false {
			t.Errorf("Claiming to not be rate limited.")
		}

		if inst.ReasonCode != ReasonRateLimitExceeded {
			t.Errorf("Unexpected reason code: %s", inst.Reason())
		}
	})

	t.Run("Concurrent throttling", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("X-RateLimit-Limit", valueZero)
		headers.Add("X-RateLimit-Remaining", valueZero)
		headers.Add("X-RateLImit-Reset", valueTimestamp)

		inst := NewRateLimitInfo(http.StatusTooManyRequests, headers)

		if inst.Limit != valueZeroInt {
			t.Errorf("Rate limit mismatch: %#v != %#v",
				inst.Limit,
				valueZeroInt)
		}

		if inst.Remaining != valueZeroInt {
			t.Errorf("Remaining mismatch: %#v != %#v",
				inst.Remaining,
				valueZeroInt)
		}

		if inst.IsRateLimited() == false {
			t.Errorf("Claiming to not be rate limited.")
		}

		if inst.ReasonCode != ReasonRateLimitExceeded {
			t.Errorf("Unexpected reason code: %s", inst.Reason())
		}
	})

	t.Run("429", func(t *testing.T) {
		headers := http.Header{}
		headers.Add("X-RateLimit-Limit", valueThousand)
		headers.Add("X-RateLimit-Remaining", valueHundred)
		headers.Add("X-RateLimit-Reset", valueTimestamp)

		inst := NewRateLimitInfo(http.StatusTooManyRequests, headers)

		if inst.Limit != valueThousandInt {
			t.Errorf("Rate limit mismatch: %#v != %#v",
				inst.Limit,
				valueThousandInt)
		}

		if inst.Remaining != valueHundredInt {
			t.Errorf("Remaining mismatch: %#v != %#v",
				inst.Remaining,
				valueHundredInt)
		}

		if inst.IsRateLimited() == false {
			t.Error("Claiming to not be rate limited.")
		}

		if inst.ReasonCode != ReasonHTTPStatusCode {
			t.Errorf("Unexpected reason code: %s", inst.Reason())
		}
	})
}

func TestParseHeaderInt(t *testing.T) {
}

func TestParseResetTime(t *testing.T) {
}

// ** Fuzzing:

func FuzzParseHeaderInt(f *testing.F) {
	// Corpus of interesting values
	seeds := []string{
		"123",                   // simple good case
		" 456 ",                 // with whitespace
		"-789",                  // negative number
		"0",                     // zero
		"00000000000000123",     // padded
		"",                      // empty
		" ",                     // just space
		"\t\n",                  // whitespace chars
		"abc",                   // pure garbage
		"123abc",                // prefix numeric
		"abc123",                // suffix numeric
		"0xFF",                  // hex (invalid)
		"1.23",                  // float (nope)
		"+42",                   // positive sign
		"--5",                   // double sign
		"999999999999999999999", // overflow
		"１２３",                   // full-width Unicode digits
		"٣٤٥",                   // Arabic-Indic digits (٣ = 3)
		"٠",                     // Arabic zero
		"\u200B123",             // zero-width space prefix
		"123\u200B",             // zero-width space suffix
		"\uFEFF123",             // BOM prefix
		"123\000",               // NUL in string
		"\u202E123",             // RTL override
		"💩",                     // emoji!
		"💩123💩",                 // emoji-wrapped number
	}

	for _, seed := range seeds {
		f.Add("X-RateLimit-Remaining", seed)
	}

	f.Fuzz(func(t *testing.T, key, val string) {
		headers := http.Header{}
		headers.Set(key, val)
		_, _ = parseHeaderInt(headers, key)
	})
}

func FuzzParseResetTime(f *testing.F) {
	// Seeds: the Good, the Bad, and the What The Hell
	seeds := []string{
		"1628784000",                             // Epoch timestamp
		fmt.Sprint(time.Now().Unix()),            // Current epoch
		"Mon, 02 Jan 2006 15:04:05 GMT",          // RFC1123
		time.Now().UTC().Format(http.TimeFormat), // Proper RFC1123 with real timestamp
		"Tue, 25 Jul 2023 15:04:05 GMT",          // Valid RFC1123 date
		"",                                       // Empty string
		" ",                                      // Whitespace only
		"abcdef",                                 // Random string
		"１２３４５",                                  // Fullwidth unicode digits
		"١٢٣٤٥",                                  // Arabic-indic numerals
		"🚀🔥🛸",                                    // Emojis, because why not
		"-1",                                     // Negative epoch
		"999999999999999999999999",               // Too big
		"0x1a",                                   // Hex
		"2023-07-25T15:04:05Z",                   // ISO8601
		"02 Jan 06 15:04 MST",                    // RFC822
		"2 Jan 2006 15:04:05 -0700",              // RFC850-alike
		"Mon Jan 2 15:04:05 2006",                // ANSIC
		"foobar 123 Mon, 02 Jan 2006",            // Corrupt garbage
		"Mon, 02 Jan 06 15:04:05 GMT",            // RFC1123Z but with 2-digit year
		"9999999999",                             // Future epoch (2286)
	}

	for _, seed := range seeds {
		f.Add("Retry-After", seed)
	}

	f.Fuzz(func(t *testing.T, key, val string) {
		headers := http.Header{}
		headers.Set(key, val)

		tm, _, ok := parseResetTime(headers, key)
		if ok {
			if tm.IsZero() {
				t.Errorf("parseResetTime(%q, %q) returned ok but zero time", val, key)
			}

			if tm.Before(time.Unix(0, 0)) || tm.After(time.Now().Add(100*365*24*time.Hour)) {
				t.Logf("parseResetTime(%q) returned implausible time: %v", val, tm)
			}
		}
	})
}

// ** Benchmarks:

func BenchmarkParseHeaderInt(b *testing.B) {
	key := "Rate-Limit-Limit"
	val := "1000"
	hdr := http.Header{}

	hdr.Set(key, val)

	b.ReportAllocs()

	for range b.N {
		parseHeaderInt(hdr, limitKeys...)
	}
}

func BenchmarkParseResetTime(b *testing.B) {
	key := "Retry-After"
	now := time.Now().Format(http.TimeFormat)
	hdr := http.Header{}

	hdr.Set(key, now)

	b.ReportAllocs()

	for range b.N {
		parseResetTime(hdr, resetKeys...)
	}
}

// * rate_limit_test.go ends here.
