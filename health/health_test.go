// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// health_test.go --- Health object tests.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

package health

// * Imports:

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

// * Constants:

// * Variables:

// * Code:

// ** Tests:

func TestNewHealthInitialState(t *testing.T) {
	inst := NewHealthWithDuration(50 * time.Millisecond)

	if inst == nil {
		t.Fatal("expected non-nil health")
	}

	if heartbeat := inst.LastHeartbeat(); heartbeat.IsZero() {
		t.Fatal("expected LastHeartbeat to be set in constructor")
	}

	if !inst.Healthy() {
		t.Fatal("expected health to be healthy after construction")
	}
}

func TestHealthyTransitions(t *testing.T) {
	timeout := 30 * time.Millisecond
	inst := NewHealthWithDuration(timeout)

	if !inst.Healthy() {
		t.Fatal("expected healthy after construction")
	}

	t.Run("Unhealthy after timeout delay", func(t *testing.T) {
		// Delay for a little while.
		time.Sleep(timeout + (25 * time.Millisecond))

		if inst.Healthy() {
			t.Fatal("expected unhealthy after timeout elapsed")
		}
	})

	t.Run("Healthy after tick", func(t *testing.T) {
		// Should be healthy again after a tick.
		inst.Tick()

		if !inst.Healthy() {
			t.Fatal("expected healthy after tick")
		}
	})
}

func TestUserDataAccessors(t *testing.T) {
	inst := NewHealthWithDuration(1 * time.Second)

	t.Run("Handles missing key", func(t *testing.T) {
		if _, ok := inst.UserGet("nope"); ok {
			t.Fatal("expected missing key to return ok=false")
		}
	})

	// Set some data.
	inst.UserSet("answer", 42)

	t.Run("Handles value for key", func(t *testing.T) {
		ans, ok := inst.UserGet("answer")

		if !ok {
			t.Fatal("expected key 'answer' to exist")
		}

		if got, want := ans.(int), 42; got != want {
			t.Fatalf("unexpected value: %#v != %#v", got, want)
		}
	})
}

func TestUserSetZeroValue(t *testing.T) {
	var inst Health

	inst.UserSet("x", "y")

	val, ok := inst.UserGet("x")

	if !ok {
		t.Fatal("expected ok=true after set on zero-value health")
	}

	if got, want := val.(string), "y"; got != want {
		t.Fatalf("unexpected value: %#v != %#v", got, want)
	}
}

func TestMarshalJSON(t *testing.T) {
	inst := NewDefaultHealth()

	// Set some user data.
	inst.UserSet("testing", "yes")

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	var out healthMarshal

	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if !out.Healthy {
		t.Fatal("expected is_healthy to be true")
	}

	if out.Heartbeat.IsZero() {
		t.Fatal("expected last_heartbeat to be non-zero")
	}

	if out.UserData == nil {
		t.Fatal("expected userdata to be non-nil")
	}

	if got, ok := out.UserData["testing"]; !ok || got.(string) != "yes" {
		t.Fatalf("unexpected userdata[testing]=yes, got=%v ok=%v",
			got,

			ok)
	}
}

func TestMarshalJSONNil(t *testing.T) {
	var inst *Health

	data, err := json.Marshal(inst)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if string(data) != "null" {
		t.Fatalf("expected JSON null: %v", string(data))
	}
}

func TestConcurrencySmoke(t *testing.T) {
	inst := NewHealthWithDuration(250 * time.Millisecond)

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// Writer: ticks + userdata writes.
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				inst.Tick()
				inst.UserSet("i", i)
				i++
			}
		}
	}()

	// Readers: Healthy/LastHeartbeat/UserGet/Marshal.
	readers := 4
	wg.Add(readers)
	for r := 0; r < readers; r++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_ = inst.Healthy()
					_ = inst.LastHeartbeat()
					_, _ = inst.UserGet("i")
					_, _ = json.Marshal(inst)
				}
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)
	close(stop)
	wg.Wait()
}

// * health_test.go ends here.
