// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// errors.go --- Error definitions.
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

import "gitlab.com/tozd/go/errors"

// * Variables:

var (

	// Signalled when a file is not a regular file.
	//
	// That is not a symlink, pipe, socket, device, etc.
	ErrNotRegular = errors.Base("not a regular file")

	// Signalled when an attempt is made to process a file that is
	// too large.
	//
	// The size limit is configurable via `Options`.
	ErrTooLarge = errors.Base("file exceeds size limit")

	// Signalled if a size option is invalid.
	ErrInvalidSize = errors.Base("invalid size")

	// Signalled if we're trying to operate on a symbolic link without
	// `FollowSymlinks` enabled.
	ErrSymlinkDenied = errors.Base("symlink not allowed")
)

// * Code:

// * errors.go ends here.
