// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// float.go --- Floating-point utils.
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

package math

// * Imports:

import (
	"fmt"
	gomath "math"
)

// * Code:

// Rounds a 64-bit floating point number to the nearest integer,
// returning it as an int. Values halfway between integers are rounded
// away from zero.
func RoundI(num float64) int {
	return int(num + gomath.Copysign(0.5, num))
}

// Rounds num to the given number of decimal places and returns the result
// as a float64. Unlike RoundF, it uses integer rounding logic (via RoundI),
// which may behave slightly differently around half-values.
func ToFixed(num float64, precision uint) float64 {
	ratio := gomath.Pow(10, float64(precision))

	return float64(RoundI(num*ratio)) / ratio
}

// Rounds num to the given number of decimal places and returns the result
// as a float64, using math.Round for IEEE-754 compliant rounding.
func RoundF(num float64, precision uint) float64 {
	ratio := gomath.Pow(10, float64(precision))

	return gomath.Round(num*ratio) / ratio
}

// Format a float.
//
// If the float is NaN or infinite, then those are explicitly returned.
func FormatFloat64(num float64) string {
	switch {
	case gomath.IsNaN(num):
		return "NaN"

	case gomath.IsInf(num, +1):
		return "+Inf"

	case gomath.IsInf(num, -1):
		return "-Inf"

	default:
		return fmt.Sprintf("%g", num)
	}
}

// * float.go ends here.
