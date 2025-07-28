// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// rfc3339_test.go --- RFC3339 time tests.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
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
	"testing"
	"time"

	"gitlab.com/tozd/go/errors"
	"gopkg.in/yaml.v3"
)

// * Constants:

const (
	TestRFC3339JSON1 string = `{ "test": "2012-12-12T12:12:12Z" }`
	TestRFC3339JSON2 string = `{ "test": "2025-07-25 10:33" }`
	TestRFC3339YAML1 string = `test: 2012-12-12T12:12:12Z`
	TestRFC3339YAML2 string = `test: 2025-07-25 10:33`

	TestRFC3339String  string = "2012-12-12T12:12:12Z"
	TestRFC3339NoTZ    string = "2012-12-12T12:12:12"
	TestRFC3339Encoded string = `"2012-12-12T12:12:12Z"`
	TestRFC3339Unix    int64  = 1355314332
)

// * Variables:

// * Code:

// ** Types:

type TestRFC3339Struct struct {
	Test RFC3339 `json:"test"  yaml:"test"`
}

// ** Tests:

// *** Base tests:

func TestRFC3339(t *testing.T) {
	str := TestRFC3339String
	then := time.Unix(TestRFC3339Unix, 0).UTC()

	t.Run("From Time", func(t *testing.T) {
		val := RFC3339(then)

		if val.Time() != then.UTC() {
			t.Errorf("Unexpected value: %#v != %#v",
				val.Time(),
				then.UTC())
		}
	})

	t.Run("From string", func(t *testing.T) {
		val := RFC3339{}

		t.Run("Valid", func(t *testing.T) {
			err := val.Set(str)
			if err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			// RFC3339 time values should *always* be UTC.
			if val.Time() == then.Local() {
				t.Errorf("Not UTC!  %#v != %#v",
					val.Time(),
					then.Local())
			}

			if val.Time() != then.UTC() {
				t.Errorf("Unexpected value: %#v != %#v",
					val.Time(),
					then.UTC())
			}
		})

		t.Run("Invalid", func(t *testing.T) {
			err := val.Set("Get schwifty")
			if err == nil {
				t.Fatal("Expecting an error")
			}

			if !errors.Is(err, ErrInvalidRFC3339) {
				t.Errorf("Unexpected error: %#v", err)
			}
		})

		t.Run("No TZ", func(t *testing.T) {
			err := val.Set(TestRFC3339NoTZ)
			if err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if val.Time() != then.UTC() {
				t.Errorf("Unexpected value: %#v != %#v",
					val.Time(),
					then.UTC())
			}
		})
	}) // From string.

	t.Run("Coercion", func(t *testing.T) {
		t.Run("String", func(t *testing.T) {
			val := RFC3339(then)
			sval := val.String()

			if sval != str {
				t.Errorf("Unexpected value: %#v != %#v",
					sval,
					str)
			}
		})

		t.Run("Time", func(t *testing.T) {
			val := RFC3339(then)
			tval := val.Time()

			if tval != then {
				t.Errorf("unexpected value: %#v != %#v",
					tval,
					then)
			}
		})
	}) // Coercion

	t.Run("Shims", func(t *testing.T) {
		obj := RFC3339(then)

		t.Run("IsZero", func(t *testing.T) {
			if ok := obj.IsZero(); ok {
				t.Errorf("%v is Zero!", obj)
			}
		})

		t.Run("IsDST", func(t *testing.T) {
			if ok := obj.IsDST(); ok {
				t.Errorf("%v is DST", obj)
			}
		})
	}) // Shims.
} // Base.

// *** JSON tests:

func TestRFC3339JSON(t *testing.T) {
	then := time.Unix(TestRFC3339Unix, 0).UTC()

	t.Run("Good data", func(t *testing.T) {
		obj := TestRFC3339Struct{}

		err := json.Unmarshal([]byte(TestRFC3339JSON1), &obj)
		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		if obj.Test.Time() != then {
			t.Errorf("Unexpected value: %#v != %#v",
				obj.Test.Time(),
				then)
		}
	})

	t.Run("Bad data", func(t *testing.T) {
		obj := TestRFC3339Struct{}

		err := json.Unmarshal([]byte(TestRFC3339JSON2), &obj)
		if err == nil {
			t.Fatal("Expected an error")
		}

		if !errors.Is(err, ErrInvalidRFC3339) {
			t.Errorf("Unexpected error: %#v", err)
		}
	})

	t.Run("Marshal", func(t *testing.T) {
		obj := RFC3339(then)

		val, err := json.Marshal(obj)
		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		sval := string(val)
		if sval != TestRFC3339Encoded {
			t.Errorf("Unexpected value: %#v != %#v",
				sval,
				TestRFC3339Encoded)
		}
	})
} // JSON.

// *** YAML tests:

func TestRFC3339YAML(t *testing.T) {
	then := time.Unix(TestRFC3339Unix, 0).UTC()

	t.Run("Good data", func(t *testing.T) {
		obj := TestRFC3339Struct{}

		err := yaml.Unmarshal([]byte(TestRFC3339JSON1), &obj)
		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		if obj.Test.Time() != then {
			t.Errorf("Unexpected value: %#v != %#v",
				obj.Test.Time(),
				then)
		}
	})

	t.Run("Bad data", func(t *testing.T) {
		obj := TestRFC3339Struct{}

		err := yaml.Unmarshal([]byte(TestRFC3339JSON2), &obj)
		if err == nil {
			t.Fatal("Expected an error")
		}

		if !errors.Is(err, ErrInvalidRFC3339) {
			t.Errorf("Unexpected error: %#v", err)
		}
	})

	t.Run("Marshal", func(t *testing.T) {
		obj := RFC3339(then)
		yaml_is_crap_and_I_hate_it := TestRFC3339Encoded + "\n"

		val, err := yaml.Marshal(obj)
		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		sval := string(val)
		if sval != yaml_is_crap_and_I_hate_it {
			t.Errorf("Unexpected value: %#v != %#v",
				sval,
				yaml_is_crap_and_I_hate_it)
		}
	})

	// Test that the `SetYAML` method works.
	t.Run("SetYAML", func(t *testing.T) {
		obj := RFC3339{}

		// Test with valid data.
		t.Run("Valid data", func(t *testing.T) {
			err := obj.SetYAML(TestRFC3339String)
			if err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if obj.Time() != then {
				t.Errorf(
					"Unexpected value: %#v != %#v",
					obj.Time(),
					then,
				)
			}
		})

		// Test with invalid duration.
		t.Run("Invalid duration", func(t *testing.T) {
			str := "2 glorks"
			err := obj.SetYAML(str)
			if err == nil {
				t.Fatal("Expected an error")
			}

			if !errors.Is(err, ErrInvalidRFC3339) {
				t.Fatalf("Unexpected error: %#v", err)
			}
		})

		// Test with a non-string value.
		t.Run("Non-string", func(t *testing.T) {
			str := 63
			err := obj.SetYAML(str)
			if err == nil {
				t.Fatal("Expected an error")
			}

			if !errors.Is(err, ErrRFC3339NotString) {
				t.Fatalf("Unexpected error: %#v", err)
			}
		})
	}) // SetYAML.
}

// ** Benchmark:

func BenchmarkRFC3339Parse(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		ParseRFC3339(TestRFC3339String)
	}
}

// ** Fuzzing:

func FuzzRFC3339Parse(f *testing.F) {
	// Valid.
	f.Add("2023-03-01T12:34:56Z")
	f.Add("2006-01-02T15:04:05+07:00")
	f.Add("1999-12-31T23:59:59-08:00")

	// Valid but spicy.
	f.Add("0000-01-01T00:00:00Z")           // Year 0 (Go handles it, ISO does not).
	f.Add("9999-12-31T23:59:59Z")           // Max time.
	f.Add("2006-01-02T15:04:05+00:00")      // Common format.
	f.Add("2012-02-29T12:00:00Z")           // Leap day.
	f.Add("2012-12-12T12:12:12.999999999Z") // Max fractional second.

	// Technically valid but seriously sus.
	f.Add("2024-04-31T00:00:00Z")      // April has 30 days, but some parsers eat it.
	f.Add("2012-13-01T00:00:00Z")      // Month 13.
	f.Add("2012-12-32T00:00:00Z")      // Day 32.
	f.Add("2012-12-12T24:00:00Z")      // 24:00:00 is allowed by ISO 8601 but rare.
	f.Add("2012-12-12T12:60:00Z")      // 60 minutes.
	f.Add("2012-12-12T12:59:60Z")      // Leap second.
	f.Add("2012-12-12T12:12:12.")      // Trailing dot.
	f.Add("2012-12-12T12:12:12.Z")     // Dot before Z.
	f.Add("2012-12-12T12:12:12+14:00") // Max valid offset.
	f.Add("2012-12-12T12:12:12-12:00") // Min valid offset.

	// Cursed.
	f.Add("ಠ_ಠ")                               // Unicode test.
	f.Add("💣2023-01-01T00:00:00Z💥")            // Valid in middle of garbage.
	f.Add("2001-02-03T04:05:06.7.8.9Z")        // WTF fractional nonsense.
	f.Add("9999-99-99T99:99:99Z")              // Max jank.
	f.Add("this-is-not-a-date")                // Random string.
	f.Add("")                                  // Empty string.
	f.Add(" ")                                 // Single space.
	f.Add("    ")                              // Whitespace.
	f.Add("\x00\x01\x02")                      // Control chars.
	f.Add("2012-12-12t12:12:12z")              // Lowercase `t` and `z`.
	f.Add("2012-12-12T12:12:12Z\n")            // Trailing newline.
	f.Add("2012-12-12T12:12:12Z\x00")          // Null-terminated.
	f.Add("2025-13-40T25:61:61Z")              // Absurd values.
	f.Add("2025-07-25T08:00:00")               // Missing timezone.
	f.Add("2025-07-25")                        // No time part.
	f.Add("T08:00:00Z")                        // No date part.
	f.Add("2025/07/25T08:00:00Z")              // Wrong separator.
	f.Add("2025-07-25 08:00:00Z")              // Space instead of T.
	f.Add("2025-07-25T08:00:00.000000000000Z") // Stupidly high precision.
	f.Add("2025-07-25T08:00:00z")              // Lowercase Z.
	f.Add("2025-07-25T08:00:00+99:99")         // Bonkers offset.
	f.Add("20🦀25-07-25T08:00:00Z")             // Emoji in the year.
	f.Add("2025-07-25T08:00:00\nZ")            // Newline in the middle.
	f.Add("2025-07-25T08:00:00Z💀")             // Valid with garbage trailing.
	f.Add("2025-07-25T08:00:00Z\000")          // Null byte.

	// Aaaaaand... hold on to your pants.
	f.Add(strings.Repeat("9", 10_000)) // 10KB of 9s — stress string handling.

	f.Fuzz(func(t *testing.T, input string) {
		_, err := time.Parse(time.RFC3339, input)
		if err != nil {
			// Lots of stuff should fail to parse.
			// What we don't want, though, are panics.
			return
		}

		// Round-trip test.
		parsed, err := time.Parse(time.RFC3339, input)
		if err != nil {
			t.Fatalf("unexpected parse fail on valid input %q: %v",
				input,
				err)
		}

		out := parsed.Format(time.RFC3339)
		_, err = time.Parse(time.RFC3339, out)
		if err != nil {
			t.Errorf("round-trip failed on %q -> %q", input, out)
		}
	})
}

// * rfc3339_test.go ends here.
