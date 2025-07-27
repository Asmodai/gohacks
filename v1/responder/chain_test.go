// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// chain_test.go --- Responder chain tests.
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

package responder

// * Imports:

import (
	"fmt"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/v1/events"
)

// * Constants:

// * Variables:

// * Code:

// ** Dummy event:

type dummyEvent struct {
	events.Time

	kind string
}

func (e *dummyEvent) Kind() string   { return e.kind }
func (e *dummyEvent) String() string { return "Dummy Event: " + e.Kind() }

// ** Dummy responder:

type dummyResponder struct {
	name        string
	typ         string
	log         *[]string
	accepts     string
	useResponse bool
	response    any
}

func (d *dummyResponder) Name() string { return d.name }
func (d *dummyResponder) Type() string { return d.typ }

func (d *dummyResponder) SetUseResponse(val bool) { d.useResponse = val }
func (d *dummyResponder) SetResponse(val any)     { d.response = val }

func (d *dummyResponder) RespondsTo(evt events.Event) bool {
	e, ok := evt.(*dummyEvent)

	return ok && e.Kind() == d.accepts
}

func (d *dummyResponder) Send(evt events.Event) events.Event {
	var result events.Event

	if d.log != nil {
		*d.log = append(*d.log, fmt.Sprintf("Handled by %s", d.name))
	}

	if d.useResponse {
		msg := events.NewMessage(42, d.name)

		result = events.NewResponse(msg, d.response)
	} else {
		result = evt
	}

	return result
}

// ** Tests:

func TestAddAndSendFirst(t *testing.T) {
	log := []string{}

	r2resp := "shamon"
	r1 := &dummyResponder{"one", "dummy", &log, "ping", true, "foo"}
	r2 := &dummyResponder{"two", "dummy", &log, "ping", true, r2resp}

	c := NewChain("test")
	_, _ = c.AddWithPriority(r1, 5)
	_, _ = c.AddWithPriority(r2, 10)

	evt := &dummyEvent{
		Time: events.Time{
			TStamp: time.Now(),
		},
		kind: "ping",
	}

	response, ok := c.SendFirst(evt)
	if !ok {
		t.Fatal("Expected a responder to handle the first event")
	}

	if len(log) != 1 || log[0] != "Handled by two" {
		t.Fatalf("Expected responder 'two' to handle event: %v", log)
	}

	switch val := response.(type) {
	case *events.Response:
		if val.Response() != r2resp {
			t.Errorf("Expected response to be %v, got %#v",
				r2resp,
				val.Response(),
			)
		}

	default:
		t.Errorf("Invalid type.")
	}
}

func TestSendAll(t *testing.T) {
	log := []string{}

	r1 := &dummyResponder{"one", "dummy", &log, "ping", false, nil}
	r2 := &dummyResponder{"two", "dummy", &log, "ping", false, nil}

	c := NewChain("test")
	_, _ = c.Add(r1)
	_, _ = c.Add(r2)

	evt := &dummyEvent{
		Time: events.Time{
			TStamp: time.Now(),
		},
		kind: "ping",
	}

	res := c.SendAll(evt)

	if len(res) != 2 {
		t.Fatalf("Expected 2 responses, got %d", len(res))
	}

	if len(log) != 2 {
		t.Errorf("Expected both responders to log event, got %v", log)
	}
}

func TestSendNamed(t *testing.T) {
	log := []string{}
	r := &dummyResponder{"blip", "dummy", &log, "ping", false, nil}

	c := NewChain("named")
	_, _ = c.Add(r)

	evt := &dummyEvent{
		Time: events.Time{
			TStamp: time.Now(),
		},
		kind: "ping",
	}

	_, ok, found := c.SendNamed("blip", evt)

	if !found {
		t.Fatal("Expected named responder to be found")
	}

	if !ok {
		t.Fatal("Expected named responder to respond")
	}

	if len(log) != 1 || log[0] != "Handled by blip" {
		t.Errorf("Expected 'blip' to handle event, got %v", log)
	}
}

func TestRemove(t *testing.T) {
	r := &dummyResponder{"gone", "dummy", nil, "ping", false, nil}

	c := NewChain("remove")
	_, _ = c.Add(r)

	if !c.Remove(r) {
		t.Fatal("Expected Remove to return true")
	}

	if c.Count() != 0 {
		t.Errorf("Expected 0 responders, got %d", c.Count())
	}
}

func TestRespondsTo(t *testing.T) {
	r := &dummyResponder{"hi", "dummy", nil, "foo", false, nil}

	c := NewChain("check")
	_, _ = c.Add(r)

	evt1 := &dummyEvent{
		Time: events.Time{
			TStamp: time.Now(),
		},
		kind: "foo",
	}

	evt2 := &dummyEvent{
		Time: events.Time{
			TStamp: time.Now(),
		},
		kind: "bar",
	}

	if !c.RespondsTo(evt1) {
		t.Error("Expected chain to respond to event")
	}

	if c.RespondsTo(evt2) {
		t.Error("Expected chain to *not* respond to event")
	}
}

func TestSendType(t *testing.T) {
	log := []string{}
	r1 := &dummyResponder{"t1", "alpha", &log, "foo", false, nil}
	r2 := &dummyResponder{"t2", "beta", &log, "foo", false, nil}
	r3 := &dummyResponder{"t3", "alpha", &log, "foo", false, nil}

	c := NewChain("types")
	_, _ = c.Add(r1)
	_, _ = c.Add(r2)
	_, _ = c.Add(r3)

	evt := &dummyEvent{
		Time: events.Time{
			TStamp: time.Now(),
		},
		kind: "foo",
	}

	res := c.SendType("alpha", evt)

	if len(res) != 2 {
		t.Errorf("Expected 2 responses, got %d", len(res))
	}

	if len(log) != 2 {
		t.Errorf("Expected 2 log entries, got %d", len(log))
	}
}

func TestChainInChain(t *testing.T) {
	log := []string{}
	inner := NewChain("inner")
	outer := NewChain("outer")

	r := &dummyResponder{"foo", "dummy", &log, "bar", false, nil}
	_, _ = inner.Add(r)
	_, _ = outer.Add(inner)

	evt := &dummyEvent{
		Time: events.Time{
			TStamp: time.Now(),
		},
		kind: "bar",
	}
	_, ok := outer.SendFirst(evt)

	if !ok || len(log) != 1 {
		t.Errorf("Expected event to bubble through nested chain, got log: %v", log)
	}
}

// * chain_test.go ends here.
