// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// bool.go --- Boolean conversion.
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

package conversion

import "strings"

// * Imports:

// * Code:

func booliseInt64(value int64) (bool, bool) {
	if value == 0 {
		return false, true
	} else if value == 1 {
		return true, true
	}

	return false, false
}

func booliseUint64(value uint64) (bool, bool) {
	if value == 0 {
		return false, true
	} else if value == 1 {
		return true, true
	}

	return false, false
}

func booliseFloat64(value float64) (bool, bool) {
	if value == 0 {
		return false, true
	} else if value == 1 {
		return true, true
	}

	return false, false
}

//nolint:cyclop
func ToBool(value any) (bool, bool) {
	switch val := value.(type) {
	case bool:
		return val, true

	case int:
		return booliseInt64(int64(val))
	case int8:
		return booliseInt64(int64(val))
	case int16:
		return booliseInt64(int64(val))
	case int32:
		return booliseInt64(int64(val))
	case int64:
		return booliseInt64(val)

	case uint:
		return booliseUint64(uint64(val))
	case uint8:
		return booliseUint64(uint64(val))
	case uint16:
		return booliseUint64(uint64(val))
	case uint32:
		return booliseUint64(uint64(val))
	case uint64:
		return booliseUint64(val)

	case float32:
		return booliseFloat64(float64(val))

	case float64:
		return booliseFloat64(val)

	case string:
		lower := strings.ToLower(val)
		if lower == "true" || lower == "yes" {
			return true, true
		} else if lower == "false" || lower == "no" {
			return false, true
		}

		return false, false

	default:
		return false, false
	}
}

// * bool.go ends here.
