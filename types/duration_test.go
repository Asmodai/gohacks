// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// duration_test.go --- Duration tests.
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
	"testing"
	"time"

	"gitlab.com/tozd/go/errors"
	"gopkg.in/yaml.v3"
)

// * Constants:

const (
	TestJSON1 string = `{
		"test1": "2s",
		"test2": "2000ms"
	}`

	TestJSON2 string = `{
		"test1": "2s",
		"test2": 2000
	}`

	// Do not indent this.
	TestYAML1 string = `
test1: 2s
test2: 2000ms`

	// Do not indent this, either.
	TestYAML2 string = `
test1: 2s
test2: 2000`
)

// * Code:

// ** Types:

type TestStruct struct {
	Test1 Duration `json:"test1"  yaml:"test1"`
	Test2 Duration `json:"test2"  yaml:"test2"`
}

// ** Tests:

// *** Base tests:

// Test the base Duration methods.
func TestDuration(t *testing.T) {
	dur := time.Duration(2_000_000_000)
	str := "2s"

	// Test we can initialise with a `time.Duration`.
	t.Run("From time.Duration", func(t *testing.T) {
		val := Duration(dur)

		if val.Duration() != dur {
			t.Errorf("Unexpected value: %#v != %#v", val, dur)
		}
	})

	// Test we can initialise from a string using the `Set` method.
	t.Run("From string", func(t *testing.T) {
		val := Duration(0)

		// Test with a valid string duration.
		t.Run("Valid", func(t *testing.T) {
			err := val.Set(str)
			if err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if val.Duration() != dur {
				t.Errorf("Unexpected value: %#v != %#v", val, dur)
			}
		})

		// Test with an invalid time duration.
		t.Run("Invalid", func(t *testing.T) {
			err := val.Set("2000 glorkles, Rick!")
			if err == nil {
				t.Fatal("Expected an error")
			}

			if !errors.Is(err, ErrInvalidDuration) {
				t.Errorf("Unexpected error: %#v", err)
			}
		})
	}) // From string.

	t.Run("Coercion", func(t *testing.T) {
		// Test we can coerce to a string via `String`.
		t.Run("String", func(t *testing.T) {
			val := Duration(dur)
			sval := val.String()

			if sval != str {
				t.Errorf("Unexpected value: %#v != %#v", sval, str)
			}
		})

		// Test we can coerce to a `time.Duration` via `Duration`.
		t.Run("Duration", func(t *testing.T) {
			val := Duration(dur)
			dval := val.Duration()

			if dval != dur {
				t.Errorf("Unexpected value: %#v != %#v", dval, dur)
			}
		})
	})
} // Base.

// *** JSON tests:

// Test JSON-specific methods.
func TestDurationJSON(t *testing.T) {
	// Test we can unmarshal from valid JSON.
	t.Run("Good data", func(t *testing.T) {
		dur := time.Duration(2_000_000_000)
		obj := &TestStruct{}

		err := json.Unmarshal([]byte(TestJSON1), &obj)
		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		if obj.Test1.Duration() != dur {
			t.Errorf(
				"Unexpected value: %#v != %#v",
				obj.Test1.Duration(),
				dur,
			)
		}
	})

	// Test what happens when we unmarshal from bad JSON.
	t.Run("Bad data", func(t *testing.T) {
		obj := &TestStruct{}

		err := json.Unmarshal([]byte(TestJSON2), &obj)
		if err == nil {
			t.Fatal("Expected an error")
		}

		if !errors.Is(err, ErrInvalidDuration) {
			t.Errorf("Unexpected error: %#v", err)
		}
	})

	// Test JSON marshalling.
	t.Run("Marshal", func(t *testing.T) {
		dur := Duration(2_000_000_000)
		str := `"2s"`

		val, err := json.Marshal(dur)
		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		sval := string(val)
		if sval != str {
			t.Errorf("Unexpected value: %#v != %#v", sval, str)
		}
	})
} // JSON.

// *** YAML tests:

// Test YAML-specific methods.
func TestDurationYAML(t *testing.T) {
	// Teat we can unmarshal from good YAML.
	t.Run("Good data", func(t *testing.T) {
		dur := time.Duration(2_000_000_000)
		obj := &TestStruct{}

		err := yaml.Unmarshal([]byte(TestYAML1), &obj)
		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		if obj.Test1.Duration() != dur {
			t.Errorf(
				"Unexpected value: %#v != %#v",
				obj.Test1.Duration(),
				dur,
			)
		}
	})

	// Test what happens when we unmarshal bad YAML.
	t.Run("Bad data", func(t *testing.T) {
		obj := &TestStruct{}

		err := yaml.Unmarshal([]byte(TestYAML2), &obj)
		if err == nil {
			t.Fatal("Expected an error")
		}

		if !errors.Is(err, ErrInvalidDuration) {
			t.Errorf("Unexpected error: %#v", err)
		}
	})

	// Test YAML marshalling.
	t.Run("Marshal", func(t *testing.T) {
		dur := Duration(2_000_000_000)
		str := "2s\n"

		val, err := yaml.Marshal(dur)
		if err != nil {
			t.Fatalf("Unexpected error: %#v", err)
		}

		sval := string(val)
		if sval != str {
			t.Errorf("Unexpected value: %#v != %#v", sval, str)
		}
	})

	// Test that the `SetYAML` method works.
	t.Run("SetYAML", func(t *testing.T) {
		dur := Duration(0)
		tdur := time.Duration(2_000_000_000)

		// Test with valid data.
		t.Run("Valid data", func(t *testing.T) {
			str := "2s"
			err := dur.SetYAML(str)
			if err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if dur.Duration() != tdur {
				t.Errorf(
					"Unexpected value: %#v != %#v",
					dur.Duration(),
					tdur,
				)
			}
		})

		// Test with invalid duration.
		t.Run("Invalid duration", func(t *testing.T) {
			str := "2 glorks"
			err := dur.SetYAML(str)
			if err == nil {
				t.Fatal("Expected an error")
			}

			if !errors.Is(err, ErrInvalidDuration) {
				t.Fatalf("Unexpected error: %#v", err)
			}
		})

		// Test with a non-string value.
		t.Run("Non-string", func(t *testing.T) {
			str := 63
			err := dur.SetYAML(str)
			if err == nil {
				t.Fatal("Expected an error")
			}

			if !errors.Is(err, ErrDurationNotString) {
				t.Fatalf("Unexpected error: %#v", err)
			}
		})
	}) // SetYAML.
} // YAML.

// *** CLI tests:

// Test the CLI-specific methods.
func TestDurationCLI(t *testing.T) {
	// Test CLI datatype for `flags`.
	t.Run("Type", func(t *testing.T) {
		val := Duration(0)

		if val.Type() != cliDurationType {
			t.Errorf(
				"Unexpected value: %#v != %#v",
				val.Type(),
				cliDurationType,
			)
		}
	})

	// Test the `Validate` method.
	t.Run("Validate", func(t *testing.T) {
		min := time.Duration(2_000_000_000)
		max := time.Duration(20_000_000_000)

		// Test with a good duration.
		t.Run("Valid", func(t *testing.T) {
			val := Duration(13_000_000_000)

			err := val.Validate(min, max)
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}
		})

		// Test with an out-of-bounds duration.
		t.Run("Invalid", func(t *testing.T) {
			val := Duration(420_000_000_000)

			err := val.Validate(min, max)
			if err == nil {
				t.Fatal("Expected an error")
			}

			if !errors.Is(err, ErrOutOfBounds) {
				t.Errorf("Unexpected error: %#v", err)
			}
		})
	})
} // CLI.

// ** Benchmark:

func BenchmarkRFC3339Set(b *testing.B) {
	b.ReportAllocs()

	for range b.N {
		_, err := NewFromDuration("21h")
		if err != nil {
			b.Errorf("Unexpected error: %#v", err)
		}
	}
}

// * duration_test.go ends here.
