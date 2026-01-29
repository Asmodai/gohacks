// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// interface.go --- Metadata interface.
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
//
//mock:yes

// * Comments:

// * Package:

package metadata

// * Imports:

// * Constants:

// * Variables:

// * Code:

// ** Interface:

type Metadata interface {
	List() map[string]string

	Set(string, string) error

	Get(string) (string, bool)

	SetDoc(string) Metadata
	SetSince(string) Metadata
	SetVersion(string) Metadata
	SetDeprecated(string) Metadata
	SetProtocol(string) Metadata
	SetVisibility(string) Metadata
	SetExample(string) Metadata
	SetTags(string) Metadata
	SetTagsFromSlice([]string) Metadata
	SetAuthor(string) Metadata

	GetDoc() string
	GetSince() string
	GetVersion() string
	GetDeprecated() string
	GetProtocol() string
	GetVisibility() string
	GetExample() string
	GetTags() string
	GetAuthor() string

	Clone() Metadata
	Merge(Metadata, bool) Metadata

	Tags() []string
	TagsNormalised() []string
}

// * interface.go ends here.
