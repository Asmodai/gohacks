// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// integer.go --- Integer math functions.
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

package math

// * Code:

// ** Types:

// ** Functions:

// Return the absolute value of a 64-bit signed integer value.
func AbsI64(val int64) int64 {
	if val < 0 {
		return -val
	}

	return val
}

// Clamp a signed 64-bit value to a minima and maxima.
func ClampI64(val, minVal, maxVal int64) int64 {
	if val < minVal {
		return minVal
	}

	if val > maxVal {
		return maxVal
	}

	return val
}

// Return the maximum value of the 64-bit integer values.
func MaxI64(lhs, rhs int64) int64 {
	if lhs < rhs {
		return rhs
	}

	return lhs
}

// Return the maximum value of the 32-bit integer values.
func MaxI32(lhs, rhs int32) int32 {
	// There is no overflow going on here, so...
	//
	//nolint:gosec
	return int32(MaxI64(int64(lhs), int64(rhs)))
}

// Return the maximum value of the integer values.
func MaxI(lhs, rhs int) int {
	return int(MaxI64(int64(lhs), int64(rhs)))
}

// Return the minimum value of the 64-bit integer values.
func MinI64(lhs, rhs int64) int64 {
	if lhs < rhs {
		return lhs
	}

	return rhs
}

// Return the minimum value of the 32-bit integer values.
func MinI32(lhs, rhs int32) int32 {
	// There is no overflow going on here, so...
	//
	//nolint:gosec
	return int32(MinI64(int64(lhs), int64(rhs)))
}

// Return the minimum value of the integer values.
func MinI(lhs, rhs int) int {
	return int(MinI64(int64(lhs), int64(rhs)))
}

// Ensure that the given value is within the limit of the platform-specific
// integer type and, if it is, multiply it by two.
//
// If the value would be larger than the platform integer, then the default
// value in `defValue` is returned.
//
// If `defValue` is too large, then the maximum integer size for the platform
// is returned.
//
//nolint:mnd
func WithinPlatform(value, defValue int64) int {
	const maxInt = int(^uint(0) >> 1) // platform MaxInt

	if value < int64(maxInt/2) {
		return int(value * 2)
	}

	if defValue > int64(maxInt/2) {
		return maxInt
	}

	return int(defValue * 2)
}

// * integer.go ends here.
