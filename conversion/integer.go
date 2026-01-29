// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// integer.go --- Integer conversions.
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

package conversion

// * Imports:

import "math"

// * Code:

func safeInt64ToUint64(val int64) (uint64, bool) {
	if val < 0 {
		return 0, false
	}

	return uint64(val), true
}

func safeUint64ToInt64(val uint64) (int64, bool) {
	if val > math.MaxInt64 {
		return 0, false
	}

	return int64(val), true
}

//nolint:cyclop
func ToInt64(value any) (int64, bool) {
	switch val := value.(type) {
	case int:
		return int64(val), true
	case int8:
		return int64(val), true
	case int16:
		return int64(val), true
	case int32:
		return int64(val), true
	case int64:
		return val, true

	case uint:
		return safeUint64ToInt64(uint64(val))
	case uint8:
		return safeUint64ToInt64(uint64(val))
	case uint16:
		return safeUint64ToInt64(uint64(val))
	case uint32:
		return safeUint64ToInt64(uint64(val))
	case uint64:
		return safeUint64ToInt64(val)

	case bool:
		if val {
			return 1, true
		}

		return 0, true

	default:
		return 0, false
	}
}

//nolint:cyclop
func ToUint64(value any) (uint64, bool) {
	switch val := value.(type) {
	case int:
		return safeInt64ToUint64(int64(val))
	case int8:
		return safeInt64ToUint64(int64(val))
	case int16:
		return safeInt64ToUint64(int64(val))
	case int32:
		return safeInt64ToUint64(int64(val))
	case int64:
		return safeInt64ToUint64(val)

	case uint:
		return uint64(val), true
	case uint8:
		return uint64(val), true
	case uint16:
		return uint64(val), true
	case uint32:
		return uint64(val), true
	case uint64:
		return val, true

	case bool:
		if val {
			return 1, true
		}

		return 0, true

	default:
		return 0, false
	}
}

// * integer.go ends here.
