/*
 * values.go --- Nastiness.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"math"
)

type Numeric interface {
	int64 | float64
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
	switch thing.(type) {
	case float64:
		return Number(thing.(float64))
	}

	return thing
}

/* values.go ends here. */
