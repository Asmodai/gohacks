// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// registry_test.go --- Protocol unit tests.
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

package protocols

// * Imports:

import (
	"errors"
	"testing"

	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/responder"
	"github.com/Asmodai/gohacks/selector"
)

// * Code:

// ** Types:

type selEvt struct {
	name string
}

func (e selEvt) When() (t events.Time) { return }
func (e selEvt) String() string        { return e.name }
func (e selEvt) Selector() string      { return e.name }

func noOpMethod(res responder.Respondable, evt events.Event) events.Event {
	return evt
}

// ** Tests:

func TestRegistry_ValidateAndVerify(t *testing.T) {
	// Build a protocol which requires two selectors
	p := &Protocol{
		Name:      "fs.readable",
		Selectors: []string{"fs.open", "fs.read"},
	}

	reg := NewRegistry()

	// Register with a verifier that requires a "fs." prefix on all
	// selectors.
	reg.RegisterWithVerifier(p, func(obj selector.Introspectable) error {
		for _, s := range obj.Selectors() {
			if len(s) < 3 || s[:3] != "fs." {
				return errors.New("selector not in fs.* namespace")
			}
		}
		return nil
	})

	// Build a respondable that implements both required selectors.
	sr := selector.NewRespondable("file0", "fs.File")

	// Register the required selectors on the respondable.
	srSel := []string{"fs.open", "fs.read"}
	for _, s := range srSel {
		srSel := s
		sr.Methods().Register(srSel, noOpMethod)
		// Methods() returns *Table in your code via sr.methods,
		// exposed by New accessor.
	}

	// Validate should pass
	if ok := reg.Validate("fs.readable", sr); !ok {
		t.Fatalf("expected Validate to pass when all selectors exist")
	}

	// Verify should pass (verifier sees both selectors start with fs.)
	if err := reg.Verify("fs.readable", sr); err != nil {
		t.Fatalf("expected Verify to succeed, got %v", err)
	}

	// Now remove one selector: easiest way is build a new respondable
	// with only one.
	sr2 := selector.NewRespondable("file1", "fs.File")
	sr2.Methods().Register("fs.open", noOpMethod)

	if ok := reg.Validate("fs.readable", sr2); ok {
		t.Fatalf("expected Validate to fail when a selector is missing")
	}

	// Register a protocol without a verifier and ensure Verify returns
	// ErrNoVerifierFunction.
	p2 := &Protocol{Name: "no.verifier", Selectors: []string{"foo"}}
	reg.Register(p2)

	if err := reg.Verify("no.verifier", sr); err == nil {
		t.Fatalf("expected Verify to error for protocol without verifier")
	} else if !errors.Is(err, ErrNoVerifierFunction) {
		t.Fatalf("expected ErrNoVerifierFunction, got %v", err)
	}
}

// * registry_test.go ends here.
