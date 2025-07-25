// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// values.go --- Nastiness.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

package generics

import (
	"math"
)

// Does the given float have a faction?
func HasFraction(num float64) bool {
	const epsilon = 1e-9

	if _, frac := math.Modf(math.Abs(num)); frac < epsilon || frac > 1.0-epsilon {
		return false
	}

	return true
}

// Coerce a numeric to an integer value of appropriate signage and size.
func CoerceInt[N Numeric](val N) any {
	if val < 0 {
		// Must be signed.
		switch {
		case int64(val) < math.MinInt32:
			return int64(val)
		case int64(val) < math.MinInt16:
			return int32(val)
		case int64(val) < math.MinInt8:
			return int16(val)
		default:
			return int8(val)
		}
	}

	// Can be unsigned.
	switch {
	case uint64(val) < math.MaxUint8:
		return uint8(val)
	case uint64(val) < math.MaxUint16:
		return uint16(val)
	case uint64(val) < math.MaxUint32:
		return uint32(val)
	default:
		return uint64(val)
	}
}

// If the given thing does not have a fraction, convert it to an int.
func Number[V Numeric](val V) any {
	if HasFraction(float64(val)) {
		return val
	}

	return CoerceInt(val)
}

// Attempt to ascertain the numeric value (if any) for a given thing.
//
// If the type of `thing` is a floating-point number, then the value will be
// converted via `Number` to either an explicit float or to an integer type if
// there is no fraction.
//
// If the type of `thing` is any other type, then it will be returned with no
// conversion.
func ValueOf(thing any) any {
	switch val := thing.(type) {
	case float32:
		return Number(val)
	case float64:
		return Number(val)
	}

	return thing
}

// values.go ends here.
