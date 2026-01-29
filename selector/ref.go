// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// ref.go --- Selector references.
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

package selector

// * Imports:

import "strings"

// * Code:

// ** Types:

type Ref struct {
	Package  string // Package name.
	Name     string // Name.
	Version  string // Version, e.g. "v1", "v1.2" et al
	Internal bool   // If true, then thing is package internal.
}

// ** Functions:

// Parse a reference.
//
// Supports:
//
//	name
//	name@version
//	package:name
//	package::name
//	package:name@version
//	package::name@version
func ParseRef(ref string) (Ref, bool) {
	out := Ref{}

	if at := strings.LastIndexByte(ref, '@'); at > 0 {
		out.Version = ref[at+1:]
		ref = ref[:at]
	}

	// Do we have a package qualifier?
	if idx := strings.Index(ref, "::"); idx > 0 {
		out.Package = ref[:idx]
		out.Name = ref[idx+2:]
		out.Internal = true
	} else if idx := strings.IndexByte(ref, ':'); idx > 0 {
		out.Package = ref[:idx]
		out.Name = ref[idx+1:]
	} else {
		out.Name = ref
	}

	if len(out.Name) == 0 {
		return Ref{}, false
	}

	return out, true
}

// * ref.go ends here.
