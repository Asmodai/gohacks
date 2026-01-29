// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// span.go --- Source code spans.
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

package lucette

// * Imports:

import (
	"strings"
)

// * Variables:

//nolint:gochecknoglobals
var (
	// A span with zero values.
	ZeroSpan = &Span{}
)

// * Code:

// ** Type:

// Span within source code.
type Span struct {
	start  Position // Starting position of span.
	end    Position // Ending position of span.
	cache  string   // Cache for stringer.
	cached bool     // Are we cached?
}

// ** Methods:

// Return the span's starting position.
func (s *Span) Start() Position {
	return s.start
}

// Return the span's ending position.
func (s *Span) End() Position {
	return s.end
}

// Return the string representation of a span.
func (s *Span) String() string {
	if !s.cached {
		var sbld strings.Builder

		sbld.WriteString(s.start.String())
		sbld.WriteRune('-')
		sbld.WriteString(s.end.String())

		s.cache = sbld.String()
		s.cached = true
	}

	return s.cache
}

// ** Functions:

// Create a new span with the given start and end positions.
func NewSpan(start, end Position) *Span {
	return &Span{start: start, end: end}
}

// Create a new empty span.
func NewEmptySpan() *Span {
	pos := NewEmptyPosition()

	return NewSpan(pos, pos)
}

// * span.go ends here.
