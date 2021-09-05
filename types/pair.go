/*
 * pair.go --- Basic pair type.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package types

import (
	"fmt"
)

/*

Pair structure.

This is a cheap implementation of a pair (aka two-value tuple).

I wish generics were a thing.

*/
type Pair struct {
	First  interface{}
	Second interface{}
}

// Create a new empty pair.
func NewEmptyPair() *Pair {
	return &Pair{
		First:  nil,
		Second: nil,
	}
}

// Create a new pair.
func NewPair(first interface{}, second interface{}) *Pair {
	return &Pair{
		First:  first,
		Second: second,
	}
}

// Return a string representation of the pair.
func (p *Pair) String() string {
	return fmt.Sprintf("%v : %v", p.First, p.Second)
}

/* pair.go ends here. */
