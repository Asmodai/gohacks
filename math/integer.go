// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// integer.go --- Integer math functions.
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

func MaxI64(lhs, rhs int64) int64 {
	if lhs < rhs {
		return rhs
	}

	return lhs
}

func MaxI32(lhs, rhs int32) int32 {
	return int32(MaxI64(int64(lhs), int64(rhs)))
}

func MaxI(lhs, rhs int) int {
	return int(MaxI64(int64(lhs), int64(rhs)))
}

func MinI64(lhs, rhs int64) int64 {
	if lhs < rhs {
		return lhs
	}

	return rhs
}

func MinI32(lhs, rhs int32) int32 {
	return int32(MinI64(int64(lhs), int64(rhs)))
}

func MinI(lhs, rhs int) int {
	return int(MinI64(int64(lhs), int64(rhs)))
}

// * integer.go ends here.
