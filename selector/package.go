// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// package.go --- Packages.
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

package selector

// * Imports:

import "sync"

// * Constants:

// * Variables:

// * Code:

// ** Types:

// Packages.
//
// A package provides a table of selectors, aliases, and defaults.
type Package struct {
	mu sync.RWMutex

	Name  string
	Table *Table

	exports  map[string]struct{} // Exported names.
	aliases  map[string]string   // Aliases.
	defaults map[string]string   // Defaults.
}

// ** Methods:

// Export a selector.
func (pkg *Package) Export(name string) bool {
	pkg.mu.Lock()
	defer pkg.mu.Unlock()

	if entry, found := pkg.Table.Get(name); !found || entry == nil {
		return false
	}

	pkg.exports[name] = struct{}{}

	return true
}

func (pkg *Package) Unexport(name string) {
	pkg.mu.Lock()
	defer pkg.mu.Unlock()

	delete(pkg.exports, name)
}

// Is the given selector exported?
func (pkg *Package) IsExported(name string) bool {
	pkg.mu.RLock()
	defer pkg.mu.RUnlock()

	_, exported := pkg.exports[name]

	return exported
}

// Create an alias that maps alias to target.
func (pkg *Package) Alias(alias, target string) bool {
	pkg.mu.Lock()
	defer pkg.mu.Unlock()

	if _, found := pkg.Table.Get(target); !found {
		return false
	}

	pkg.aliases[alias] = target

	return true
}

// Resolve a selector.
//
// If the specified selector is an alias, then its target is returned.
func (pkg *Package) ResolveAlias(name string) string {
	pkg.mu.RLock()
	defer pkg.mu.RUnlock()

	if tgt, found := pkg.aliases[name]; found {
		return tgt
	}

	return name
}

// Sets a default selector.
func (pkg *Package) SetDefault(operator, name string) bool {
	pkg.mu.Lock()
	defer pkg.mu.Unlock()

	if _, found := pkg.Table.Get(name); !found {
		return false
	}

	pkg.defaults[operator] = name

	return true
}

// Returns the default for the given operation.
func (pkg *Package) GetDefault(op string) (string, bool) {
	pkg.mu.RLock()
	defer pkg.mu.RUnlock()

	if len(op) > 0 {
		if name, found := pkg.defaults[op]; found {
			return name, true
		}
	}

	if name, found := pkg.defaults["_"]; found {
		return name, true
	}

	return "", false
}

// ** Functions:

func NewPackage(name string) *Package {
	return &Package{
		Name:     name,
		Table:    NewTable(),
		exports:  make(map[string]struct{}),
		aliases:  make(map[string]string),
		defaults: make(map[string]string),
	}
}

// * package.go ends here.
