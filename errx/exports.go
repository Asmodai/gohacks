// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// exports.go --- Export hackery.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

package errx

// * Imports:

import (
	tozd "gitlab.com/tozd/go/errors"
)

// * Variables:

// We simply import the interesting symbols that we want from our given
// error handling package.
//
// In this case, the error handling package we use is gitlab.com/tozd/go/errors
//
//nolint:gochecknoglobals
var (
	AllDetails   = tozd.AllDetails
	As           = tozd.As
	Base         = tozd.Base
	BaseWrap     = tozd.BaseWrap
	BaseWrapf    = tozd.BaseWrapf
	Basef        = tozd.Basef
	Cause        = tozd.Cause
	Details      = tozd.Details
	Errorf       = tozd.Errorf
	Is           = tozd.Is
	Join         = tozd.Join
	New          = tozd.New
	Prefix       = tozd.Prefix
	Unjoin       = tozd.Unjoin
	Unwrap       = tozd.Unwrap
	WithDetails  = tozd.WithDetails
	WithMessage  = tozd.WithMessage
	WithMessagef = tozd.WithMessagef
	WithStack    = tozd.WithStack
	Wrap         = tozd.Wrap
	WrapWith     = tozd.WrapWith
	Wrapf        = tozd.Wrapf
)

type Error tozd.E

// * exports.go ends here.
