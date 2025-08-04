// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// array.go --- Array conversion hacks.
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

package conversion

// * Imports:

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// * Constants:

// * Variables:

// * Code:

type Number interface {
	constraints.Integer | constraints.Float
}

func NumericArrayToFloat64[T Number](in []T) []float64 {
	out := make([]float64, len(in))

	for idx, val := range in {
		out[idx] = float64(val)
	}

	return out
}

//nolint:cyclop
func AnyArrayToFloat64Array(input any) ([]float64, bool) {
	switch val := input.(type) {
	case []int:
		return NumericArrayToFloat64(val), true
	case []int8:
		return NumericArrayToFloat64(val), true

	case []int16:
		return NumericArrayToFloat64(val), true

	case []int32:
		return NumericArrayToFloat64(val), true

	case []int64:
		return NumericArrayToFloat64(val), true

	case []uint:
		return NumericArrayToFloat64(val), true

	case []uint8:
		return NumericArrayToFloat64(val), true

	case []uint16:
		return NumericArrayToFloat64(val), true

	case []uint32:
		return NumericArrayToFloat64(val), true

	case []uint64:
		return NumericArrayToFloat64(val), true

	case []float32:
		return NumericArrayToFloat64(val), true

	case []float64:
		return val, true

	default:
		return nil, false
	}
}

//nolint:cyclop,nlreturn,wsl,varnamelen
func AnyArrayToStringArray(input any) ([]string, bool) {
	switch val := input.(type) {
	case []string:
		return val, true

	case []fmt.Stringer:
		out := make([]string, len(val))
		for i, v := range val {
			out[i] = v.String()
		}
		return out, true

	case [][]byte:
		out := make([]string, len(val))
		for i, v := range val {
			out[i] = string(v)
		}
		return out, true

	case [][]rune:
		out := make([]string, len(val))
		for i, v := range val {
			out[i] = string(v)
		}
		return out, true

	case []any:
		out := make([]string, len(val))
		for i, v := range val {
			switch s := v.(type) {
			case string:
				out[i] = s
			case fmt.Stringer:
				out[i] = s.String()
			case []byte:
				out[i] = string(s)
			case []rune:
				out[i] = string(s)
			default:
				return nil, false
			}
		}
		return out, true

	default:
		return nil, false
	}
}

// * array.go ends here.
