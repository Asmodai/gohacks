// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// namespace.go --- Namespaces
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

import (
	"fmt"

	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/responder"
)

// * Constants:

const (
	reasonDirect        = "direct"
	reasonDefault       = "pkg-default"
	reasonUses          = "uses"
	reasonGlobalDefault = "global-default"
)

// * Code:

// ** Types:

type Namespace struct {
	Uses   []*Package
	Shadow map[string]bool
}

// ** Methods:

// Resolve a reference.
//
//nolint:cyclop,funlen,gocognit,nestif
func (ns *Namespace) Resolve(reg *Registry, ref Ref) (ResolveResult, bool) {
	//
	// XXX Please, please split this apart.
	//
	// I'm having to nolint all kinds of stuff I haven't seen before.
	//
	Trace("Attempting to resolve %q (%q)", ref.Name, ref.Package)

	if len(ref.Package) > 0 {
		Trace("... Looking for %q in package %q",
			ref.Name,
			ref.Package)

		pkg, found := reg.GetPackage(ref.Package)
		if !found {
			return makeFail(nil,
				ref.Name,
				fmt.Sprintf("Package %s not found",
					ref.Package))
		}

		// Internal -- Not allowed.
		if ref.Internal {
			return makeFail(pkg, ref.Name, "Internal")
		}

		name := pkg.ResolveAlias(ref.Name)

		// Version mapping.
		if len(ref.Version) > 0 {
			test := name + "@" + ref.Version

			if _, found := pkg.Table.Get(test); found {
				name = test
			}
		}

		// Must be exported.
		if !ref.Internal && !pkg.IsExported(name) {
			// Try default.
			if def, ok := pkg.GetDefault(ref.Name); ok && pkg.IsExported(def) {
				return makeSucceed(
					pkg,
					pkg.Table,
					def,
					reasonDefault)
			}

			return makeFail(pkg, name, "Not exported")
		}

		if entry, ok := pkg.Table.Get(name); ok && entry != nil {
			return makeSucceed(
				pkg,
				pkg.Table,
				name,
				reasonDirect)
		}

		// Package-level defaults.
		if def, ok := pkg.GetDefault(ref.Name); ok {
			if _, ok := pkg.Table.Get(def); ok {
				return makeSucceed(
					pkg,
					pkg.Table,
					def,
					reasonDefault)
			}
		}

		return makeFail(pkg, name, "Not found in package")
	}

	// Unqualified.
	if ns.Shadow != nil && ns.Shadow[ref.Name] {
		return makeFail(nil, ref.Name, "Is shadowed")
	}

	// Check in imports.
	for _, pkg := range ns.Uses {
		Trace("... Checking imported package %q",
			pkg.Name)

		name := pkg.ResolveAlias(ref.Name)

		if len(ref.Version) > 0 {
			test := name + "@" + ref.Version

			if _, ok := pkg.Table.Get(test); ok {
				name = test
			}
		}

		if !pkg.IsExported(name) {
			if def, ok := pkg.GetDefault(ref.Name); ok && pkg.IsExported(def) {
				if _, ok := pkg.Table.Get(def); ok {
					return makeSucceed(
						pkg,
						pkg.Table,
						def,
						reasonDefault)
				}
			}

			makeFail(
				pkg,
				name,
				"Is not exported.  Skipping package.")

			continue
		}

		if entry, ok := pkg.Table.Get(name); ok && entry != nil {
			return makeSucceed(
				pkg,
				pkg.Table,
				name,
				reasonUses)
		}

		if def, ok := pkg.GetDefault(ref.Name); ok {
			if _, ok := pkg.Table.Get(def); ok {
				return makeSucceed(
					pkg,
					pkg.Table,
					def,
					reasonDefault)
			}
		}
	}

	if reg.GlobalDefault != nil {
		Trace("... Looking for %q in global default.", ref.Name)

		if def, ok := reg.GlobalDefault.GetDefault("_"); ok {
			if _, ok := reg.GlobalDefault.Table.Get(def); ok {
				return makeSucceed(
					reg.GlobalDefault,
					reg.GlobalDefault.Table,
					def,
					reasonGlobalDefault)
			}
		}
	}

	return makeFail(nil, ref.Name, "Giving up")
}

func (ns *Namespace) Dispatch(
	reg *Registry,
	raw string, target responder.Respondable, evt events.Event,
) (events.Event, bool, string) {
	ref, refOk := ParseRef(raw)
	if !refOk {
		return nil, false, "parse-error"
	}

	res, resOk := ns.Resolve(reg, ref)
	if !resOk {
		return nil, false, "unresolved"
	}

	out, outOk := res.Table.InvokeSelector(res.Name, target, evt)

	return out, outOk, res.Why
}

// ** Functions:

func makeFail(pkg *Package, name, reason string) (ResolveResult, bool) {
	var pname string

	if pkg != nil {
		pname = pkg.Name
	}

	TraceWithWrapper(true, "Failed.  Package:%s  Name:%q  Reason:%s",
		pname,
		name,
		reason)

	return ResolveResult{}, false
}

func makeSucceed(pkg *Package, table *Table, name, why string) (ResolveResult, bool) {
	var pname string

	if pkg != nil {
		pname = pkg.Name
	}

	TraceWithWrapper(true, "Found.  Package:%q  Name:%q  Why:%s",
		pname,
		name,
		why)

	ret := ResolveResult{
		Pkg:   pkg,
		Table: table,
		Name:  name,
		Why:   why}

	return ret, true
}

// * namespace.go ends here.
