// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// options.go --- Options.
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

package fileio

// * Imports:

import "math"

// * Code:

// File I/O options.
type Options struct {
	// Maximum number of bytes to read.
	//
	// If this value is set to a negative number then it is interpreted
	// as "read an unlimited number of bytes".
	//
	// However "unlimited", in this case, equates to `math.MaxInt` bytes.
	//
	// The reason for this is that we return a slice of bytes, and the
	// maximum number of elements in a Go slice is `math.MaxInt`.
	//
	// The default value is 0.
	MaxReadBytes int64

	// Should symbolic links be followed?
	//
	// If false, then symbolic links are not followed.
	//
	// The default value is false.
	FollowSymlinks bool
}

// ** methods:

func (opt *Options) sanity() error {
	if opt.MaxReadBytes < 0 {
		opt.MaxReadBytes = 0
	}

	if opt.MaxReadBytes > int64(math.MaxInt) {
		return ErrInvalidSize
	}

	return nil
}

// * options.go ends here.
