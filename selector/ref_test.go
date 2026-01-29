// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// ref_test.go --- Reference tests.
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

import (
	"fmt"
	"testing"
)

// * Imports:

// * Constants:

// * Variables:

// * Code:

func TestRef(t *testing.T) {
	data := []struct {
		input    string
		pkg      string
		name     string
		version  string
		internal bool
	}{
		{"foo", "", "foo", "", false},
		{"bar:foo", "bar", "foo", "", false},
		{"bar::foo", "bar", "foo", "", true},
		{"bar:foo@v1", "bar", "foo", "v1", false},
		{"bar::foo@v1", "bar", "foo", "v1", true},
	}

	for idx, elt := range data {
		t.Run(fmt.Sprintf("%02d %s", idx, elt.input), func(t *testing.T) {
			result, ok := ParseRef(elt.input)
			if !ok {
				t.Fatalf("Could not parse!")
			}

			if result.Package != elt.pkg {
				t.Errorf("Package: expected %q got %q",
					elt.pkg,
					result.Package)
			}

			if result.Name != elt.name {
				t.Errorf("NBame: expected %q got %q",
					elt.name,
					result.Name)
			}

			if result.Version != elt.version {
				t.Errorf("Version: expected %q got %q",
					elt.version,
					result.Version)
			}

			if result.Internal != elt.internal {
				t.Errorf("Internal: expected %v got %v",
					elt.internal,
					result.Internal)
			}
		})
	}
}

// * ref_test.go ends here.
