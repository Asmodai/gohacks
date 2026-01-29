// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// namespace_test.go --- Namespace tests.
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

import (
	"strings"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/responder"
)

// * Constants:

// * Variables:

// * Code:

// ** Test event:

type nsTestEvent struct {
	selector string
	payload  string
}

func (e *nsTestEvent) Selector() string { return e.selector }
func (e *nsTestEvent) String() string   { return e.payload }
func (e *nsTestEvent) When() time.Time  { return time.Time{} }

// ** Test target:

type nsTarget struct {
	responder.Respondable
}

func newNSTarget(name string) *nsTarget {
	return &nsTarget{Respondable: NewRespondable(name, "nsTarget")}
}

func mkPkg(name string) *Package {
	return NewPackage(name)
}

// ** Tests:

func TestNamespace_Resolve_Qualified_Internal_And_Export(t *testing.T) {
	EnableTrace()

	fs := mkPkg("fs")
	fs.Table.Register("touch", func(r responder.Respondable, e events.Event) events.Event { return e })
	fs.Export("touch")

	reg := NewRegistry()
	reg.AddPackage(fs)
	ns := &Namespace{} // no Current

	// qualified direct
	res, ok := ns.Resolve(reg, mustParse("fs:touch"))
	if !ok || res.Why != reasonDirect || res.Pkg != fs || res.Name != "touch" {
		t.Fatalf("bad resolve: %+v ok=%v", res, ok)
	}

	// internal denied by default
	if _, ok := ns.Resolve(reg, mustParse("fs::touch")); ok {
		t.Fatal("internal should be denied without AllowInternal")
	}
}

func TestNamespace_Resolve_Uses_Export_Only(t *testing.T) {
	a := mkPkg("a")
	b := mkPkg("b")
	b.Table.Register("pub", func(r responder.Respondable, e events.Event) events.Event { return e })
	b.Export("pub")
	b.Table.Register("priv", func(r responder.Respondable, e events.Event) events.Event { return e })

	reg := NewRegistry()
	reg.AddPackage(a)
	reg.AddPackage(b)
	ns := &Namespace{Uses: []*Package{b}}

	if _, ok := ns.Resolve(reg, mustParse("pub")); !ok {
		t.Fatal("exported use should resolve")
	}
	if _, ok := ns.Resolve(reg, mustParse("priv")); ok {
		t.Fatal("unexported should not resolve via uses")
	}
}

func TestNamespace_Alias_Version_Defaults(t *testing.T) {
	p := mkPkg("p")
	p.Table.Register("read@v1", func(r responder.Respondable, e events.Event) events.Event { return e })
	p.Export("read@v1")
	p.Alias("r", "read@v1")
	p.SetDefault("read", "read@v1")

	reg := NewRegistry()
	reg.AddPackage(p)
	ns := &Namespace{Uses: []*Package{p}}

	res, ok := ns.Resolve(reg, mustParse("r"))
	if !ok || res.Name != "read@v1" {
		t.Fatalf("alias failed: %+v", res)
	}

	res, ok = ns.Resolve(reg, mustParse("read"))
	if !ok || res.Why != reasonDefault || res.Name != "read@v1" {
		t.Fatalf("op default failed: %+v", res)
	}

	res, ok = ns.Resolve(reg, mustParse("p:read@v1"))
	if !ok || res.Why != reasonDirect {
		t.Fatalf("versioned failed: %+v", res)
	}
}

func TestNamespace_Shadow_And_GlobalDefault(t *testing.T) {
	fs := mkPkg("fs")
	net := mkPkg("net")

	net.Table.Register("get", func(r responder.Respondable, e events.Event) events.Event { return e })
	net.Export("get")

	fs.Table.Register("fallback", func(r responder.Respondable, e events.Event) events.Event { return e })
	fs.SetDefault("_", "fallback")

	reg := NewRegistry()
	reg.AddPackage(fs)
	reg.AddPackage(net)
	reg.GlobalDefault = fs
	ns := &Namespace{Uses: []*Package{net}, Shadow: map[string]bool{"get": true}}

	if _, ok := ns.Resolve(reg, mustParse("get")); ok {
		t.Fatal("shadowed should not resolve")
	}
	if _, ok := ns.Resolve(reg, mustParse("net:get")); !ok {
		t.Fatal("qualified should still resolve")
	}
	res, ok := ns.Resolve(reg, mustParse("woot"))
	if !ok || res.Why != reasonGlobalDefault || res.Pkg != fs || res.Name != "fallback" {
		t.Fatalf("global default failed: %+v ok=%v", res, ok)
	}
}

func mustParse(s string) Ref {
	r, ok := ParseRef(s)

	if !ok {
		panic("parse fail: " + s)
	}

	return r
}

func TestNamespace_Dispatch_Invokes_Method(t *testing.T) {
	fs := mkPkg("fs")
	reg := NewRegistry()
	reg.AddPackage(fs)

	var log strings.Builder
	fs.Table.Register("ping", func(r responder.Respondable, e events.Event) events.Event {
		log.WriteString("pong")
		return e
	})
	fs.Export("ping")

	ns := &Namespace{}
	target := newNSTarget("receiver")

	out, ok, why := ns.Dispatch(reg, "fs:ping", target, &nsTestEvent{
		selector: "ping",
		payload:  "x"})
	if !ok || out.(*nsTestEvent).payload != "x" || why != reasonDirect {
		t.Fatalf("dispatch failed: ok=%v why=%s out=%#v", ok, why, out)
	}
	if !strings.Contains(log.String(), "pong") {
		t.Fatalf("method body was not executed")
	}
}

// * namespace_test.go ends here.
