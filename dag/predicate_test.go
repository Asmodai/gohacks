// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// compiler.go --- DAG compiler.
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

// * Package:

package dag

// * Imports:

import (
	"testing"
)

// * Code:

// ** Tests:

func TestPredicateEvaluation(t *testing.T) {
	const (
		regexval  string  = "Peach"
		stringval string  = "Fortra"
		floatval  float64 = 42.0
	)

	input := DataMap{
		"Numeric": floatval,
		"String":  stringval,
		"Regex":   regexval,
	}

	//
	// Test string predicates.
	t.Run("String", func(t *testing.T) {
		const (
			good   string = "Fortra"
			goodCI string = "FORTRA"
			bad    string = "Wabbajack"
		)

		// String equality test.
		t.Run("evalStringEQ", func(t *testing.T) {
			pred1 := buildStringEQ("String", good)
			pred2 := buildStringEQ("String", bad)

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %s != %s", stringval, good)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %s == %s", stringval, bad)
			}
		})

		// String inequality test.
		t.Run("evalStringNEQ", func(t *testing.T) {
			pred1 := buildStringNEQ("String", bad)  // good
			pred2 := buildStringNEQ("String", good) // bad

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %s == %s", stringval, good)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %s != %s", stringval, bad)
			}
		})

		// String case-insensitive equality test.
		t.Run("evalStringCIEQ", func(t *testing.T) {
			pred1 := buildStringCIEQ("String", good)
			pred2 := buildStringCIEQ("String", bad)

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %s != %s", stringval, good)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %s == %s", stringval, bad)
			}
		})

		// String case-insensitive inequality test.
		t.Run("evalStringCINEQ", func(t *testing.T) {
			pred1 := buildStringCINEQ("String", bad)  // good
			pred2 := buildStringCINEQ("String", good) // bad

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %s == %s", stringval, good)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %s != %s", stringval, bad)
			}
		})

	}) // String

	t.Run("Regex", func(t *testing.T) {
		const (
			good string = "(?i)p([a-z]+)ch"
			bad  string = "c([oa]*)ee"
		)

		t.Run("evalregex_Match", func(t *testing.T) {
			pred1 := buildRegexMatch("Regex", good)
			pred2 := buildRegexMatch("Regex", bad)

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %s !~= %s", regexval, good)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %s ~= %s", regexval, bad)
			}
		})
	})

	//
	// Test numeric predicates.
	t.Run("Numeric", func(t *testing.T) {
		const (
			good    int = 42
			bad     int = 76
			greater int = 96
			lesser  int = 24
		)

		// Equality test.
		t.Run("evalEQ", func(t *testing.T) {
			pred1 := buildEQ("Numeric", good)
			pred2 := buildEQ("Numeric", bad)

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %d != %.1f", good, floatval)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %d == %.1f", bad, floatval)
			}
		})

		// Inequality test.
		t.Run("evalNEQ", func(t *testing.T) {
			pred1 := buildNEQ("Numeric", bad)  // good!
			pred2 := buildNEQ("Numeric", good) // bad!

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %d == %.1f", bad, floatval)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %d != %.1f", good, floatval)
			}
		})

		// Greater-than test.  Is value > lesser?
		t.Run("evalGT", func(t *testing.T) {
			pred1 := buildGT("Numeric", lesser)  // good!
			pred2 := buildGT("Numeric", greater) // bad!

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %d <= %.1f", lesser, floatval)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %d > %.1f", greater, floatval)
			}
		})

		// Lesser-than test.  Is value < greater?
		t.Run("evalLT", func(t *testing.T) {
			pred1 := buildLT("Numeric", greater)
			pred2 := buildLT("Numeric", lesser)

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %d >= %.1f", greater, floatval)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %d < %.1f", lesser, floatval)
			}
		})

		// Greater-than-or-equal test.  Is value >= lesser?
		t.Run("evalGTE", func(t *testing.T) {
			pred1 := buildGTE("Numeric", lesser)  // good!
			pred2 := buildGTE("Numeric", greater) // bad!

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %d < %.1f", lesser, floatval)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %d >= %.1f", greater, floatval)
			}
		})

		// Lesser-than-or-equal test.  Is value <= greater?
		t.Run("evalLTE", func(t *testing.T) {
			pred1 := buildLTE("Numeric", greater)
			pred2 := buildLTE("Numeric", lesser)

			if !pred1.Eval(&pred1, input) {
				t.Errorf("Failure, %d >= %.1f", greater, floatval)
			}

			if pred2.Eval(&pred2, input) {
				t.Errorf("Failure, %d < %.1f", lesser, floatval)
			}
		})

	}) //Numeric
}

// * predicate_test.go ends here.
