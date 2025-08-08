// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// selector_test.go --- Tests.
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
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/metadata"
	"github.com/Asmodai/gohacks/responder"
)

// * Constants:

// * Variables:

// * Code:

// ** Types:

// *** Dummy event:

type testEvent struct {
	selector string
	payload  string
}

func (e *testEvent) Selector() string { return e.selector }
func (e *testEvent) String() string   { return e.payload }
func (e *testEvent) When() time.Time  { return time.Time{} }

// *** Dummy respondable:

type testResponder struct {
	*Respondable
	log *strings.Builder
}

func newTestResponder(name string) *testResponder {
	res := &testResponder{
		Respondable: NewRespondable(name, "testResponder"),
		log:         &strings.Builder{},
	}
	return res
}

// ** Tests:

func TestSelectorDispatch(t *testing.T) {
	res := newTestResponder("myObject")
	EnableTrace()

	res.Respondable.Methods().Register(
		"greet",
		func(r responder.Respondable, e events.Event) events.Event {
			ev := e.(*testEvent)
			res.log.WriteString("hello " + ev.payload)
			return e
		})

	evt := &testEvent{selector: "greet", payload: "world"}
	out := res.Invoke(evt)

	if out.(*testEvent).payload != "world" {
		t.Errorf("Expected payload 'world', got %v", out)
	}

	if !strings.Contains(res.log.String(), "hello world") {
		t.Errorf("Log did not contain greeting: %v", res.log.String())
	}
}

func TestBeforeShortCircuit(t *testing.T) {
	res := newTestResponder("gatekeeper")

	res.Respondable.Methods().Register(
		"gate",
		func(r responder.Respondable, e events.Event) events.Event {
			return &testEvent{selector: "gate", payload: "allowed"}
		})

	_ = res.Respondable.Methods().AddBefore(
		"gate",
		func(r responder.Respondable, e events.Event) events.Event {
			return NewSelectorError(errors.New("wrong before"))
		})

	_ = res.Respondable.Methods().AddBeforeWithPriority(
		1,
		"gate",
		func(r responder.Respondable, e events.Event) events.Event {
			return NewSelectorError(errors.New("denied"))
		})

	_ = res.Respondable.Methods().AddAfter(
		"gate",
		func(r responder.Respondable, e events.Event) events.Event {
			// You should *not* see this.
			//
			// Basically, the :before that returns "denied"
			// should stop the entire happening train.
			//
			// e.g. :before signals failure?  no :primary or
			// :after.
			Trace(":after - responder:%s:%s - event:%s\n",
				r.Name(),
				r.Type(),
				e.String())

			return e
		})

	result := res.Invoke(&testEvent{selector: "gate", payload: "request"})

	if err, ok := result.(*SelectorError); !ok || err.Error().Error() != "denied" {
		t.Errorf("Expected short-circuit error, got: %#v", result)
	}
}

func TestMetadataAccess(t *testing.T) {
	res := newTestResponder("annotated")

	res.Respondable.Methods().Register(
		"annotate",
		func(r responder.Respondable, e events.Event) events.Event {
			return e
		})

	md, err := res.Respondable.Methods().Metadata("annotate")
	if err != nil {
		t.Fatalf("Failed to get metadata: %v", err)
	}

	// Try wrong way.
	err = md.Set("version", "1.0")
	switch {
	case err == nil:
		t.Fatal("Expected an error, didn't get one.")
	case !errors.Is(err, metadata.ErrKeyIsReserved):
		t.Fatalf("Unexpected error: %v", err)
	}

	// Try the right way.
	md.SetVersion("1.0")

	data, err := res.MetadataForSelector("annotate")
	if err != nil {
		t.Fatalf("ListMetadata failed: %v", err)
	}

	if data["version"] != "1.0" {
		t.Errorf("Expected version '1.0', got %q", data["version"])
	}
}

func TestDumpIntrospectable(t *testing.T) {
	res := newTestResponder("debuggable")

	res.Respondable.Methods().Register(
		"foo",
		func(r responder.Respondable, e events.Event) events.Event {
			return e
		})

	md, err := res.Respondable.Methods().Metadata("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	md.SetVersion("1.0").
		SetSince("2.0").
		SetAuthor("Paul").
		SetTags("example").
		SetDoc("This method claims to unlock the secrets of the " +
			"universe, but it is probbly lying to you.")

	dump := DumpIntrospectableInfo(res)
	t.Logf("\n%s\n", dump)

	if !strings.Contains(dump, "Object: debuggable") {
		t.Errorf("Dump output incorrect: %s", dump)
	}
}

func TestSetPrimary_Swap(t *testing.T) {
	tbl := NewTable()
	tbl.Register("do", func(r responder.Respondable, e events.Event) events.Event {
		return &testEvent{selector: "do", payload: "old"}
	})

	// Swap OK
	_, err := tbl.SetPrimary("do", func(r responder.Respondable, e events.Event) events.Event {
		return &testEvent{selector: "do", payload: "new"}
	})

	if err != nil {
		t.Fatalf("SetPrimary failed: %v", err)
	}

	out, ok := tbl.InvokeSelector("do", NewRespondable("x", "y"), &testEvent{selector: "do"})
	if !ok || out.(*testEvent).payload != "new" {
		t.Fatalf("primary not swapped: %#v ok=%v", out, ok)
	}

	// Missing selector
	_, err = tbl.SetPrimary("missing", func(_ responder.Respondable, _ events.Event) events.Event {
		return nil
	})

	if !errors.Is(err, ErrSelectorNotFound) {
		t.Fatalf("expected ErrSelectorNotFound, got %v", err)
	}

	// Nil method
	_, err = tbl.SetPrimary("do", nil)
	if !errors.Is(err, ErrNoMethodSpecified) {
		t.Fatalf("expected ErrNoMethodSpecified, got %v", err)
	}
}

// * selector_test.go ends here.
