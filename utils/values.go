/*
 * values.go --- Nastiness.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package utils

import (
	"math"
)

type Numeric interface {
	~int64 | ~float64
}

// Does the given float have a faction?
func HasFraction(num float64) bool {
	const epsilon = 1e-9

	if _, frac := math.Modf(math.Abs(num)); frac < epsilon || frac > 1.0-epsilon {
		return false
	}

	return true
}

// If the given thing does not have a fraction, convert it to an int.
func Number[V Numeric](val V) interface{} {
	if HasFraction(float64(val)) {
		return val
	}

	return int64(val)
}

/*
 * This pretty much only exists because, as far as Go's JSON parser is
 * concerned, all numbers are floats.
 *
 * This will check if the thing passed to it is `float64`.  If it is, then
 * the float will be passed off to `Number`, and converted to an int if
 * there is no fraction.
 *
 * This could be used for other things too, I guess.
 */
func ValueOf(thing interface{}) interface{} {
	//nolint:gosimple
	switch thing.(type) {
	case float64:
		return Number(thing.(float64))
	}

	return thing
}

/* values.go ends here. */
