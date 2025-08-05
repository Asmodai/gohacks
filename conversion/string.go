// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// string.go --- String conversion.
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

// * Imports:

import (
	"fmt"
	"strconv"
)

// * Code:

// Convert a value to a string.
//
//nolint:cyclop
func ToString(value any) (string, bool) {
	switch val := value.(type) {
	case string:
		return val, true

	case fmt.Stringer:
		return val.String(), true

	case []byte:
		return string(val), true

	case bool:
		return strconv.FormatBool(val), true

	case int:
		return strconv.Itoa(val), true

	case int8:
		return strconv.Itoa(int(val)), true

	case int16:
		return strconv.Itoa(int(val)), true

	case int32:
		return strconv.Itoa(int(val)), true

	case int64:
		return strconv.FormatInt(val, 10), true

	case uint:
		return strconv.FormatUint(uint64(val), 10), true

	case uint8:
		return strconv.FormatUint(uint64(val), 10), true

	case uint16:
		return strconv.FormatUint(uint64(val), 10), true

	case uint32:
		return strconv.FormatUint(uint64(val), 10), true

	case uint64:
		return strconv.FormatUint(val, 10), true

	case float32:
		return fmt.Sprintf("%v", val), true

	case float64:
		return fmt.Sprintf("%v", val), true

	default:
		return "", false
	}
}

// * string.go ends here.
