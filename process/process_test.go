// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// process_test.go --- Process state tests.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

package process

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/responder"
	"github.com/Asmodai/gohacks/types"
)

var (
	manager_inst      Manager
	testProc          *Process
	testEProc         *Process
	testResponderProc *Process

	testBlockingSend bool
	fromNonblocking  interface{}
	fromBlocking     interface{}

	EveryVal int = 0
)

type DummyResponder struct {
}

func (r *DummyResponder) Name() string { return "Test Responder" }
func (r *DummyResponder) Type() string { return "DummyResponder" }

func (r *DummyResponder) RespondsTo(event events.Event) bool {
	switch val := event.(type) {
	case *events.Message:
		switch val.Command() {
		case "test1":
			return true

		case "test2":
			return true

		default:
			return false
		}

	default:
		return false
	}
}

func (r *DummyResponder) Invoke(event events.Event) events.Event {
	cmd, ok := event.(*events.Message)
	if !ok {
		return events.NewError(fmt.Errorf("not events.Message"))
	}

	switch cmd.Command() {
	case "test1":
		return events.NewResponse(cmd, "Ok, closing the pod bay doors")

	case "test2":
		return events.NewResponse(cmd, "Uh, Tuesday?")
	}

	return events.NewError(fmt.Errorf("unknown command: %s", cmd.Command()))
}

// Test `every` function.
func EveryFn(state *State) {
	EveryVal++
}

// Create a new test config.
func NewTestConfig() *Config {
	cnf := NewConfig()

	cnf.Name = "Test"
	cnf.Interval = 0
	cnf.Function = nil

	return cnf
}

func NewResponderConfig(object responder.Respondable) *Config {
	cnf := NewConfig()

	cnf.Name = "Responder Test"
	cnf.Interval = 0
	cnf.Function = nil
	cnf.Responder = object

	return cnf
}

// Create a new config for the `every` event test.
func NewEveryConfig() *Config {
	return &Config{
		Name:     "Every",
		Interval: types.Duration(1 * time.Second),
		Function: EveryFn,
	}
}

// Main testing function.
func TestMain(m *testing.M) {
	var err error

	log.Println("Setting up processes.")

	lgr := logger.NewDefaultLogger()
	ctx := context.Background()

	ctx, err = logger.SetLogger(ctx, lgr)
	if err != nil {
		log.Printf("Could not set up DI: %#v", err)
		os.Exit(128)
	}

	manager_inst = NewManagerWithContext(ctx)
	testProc = manager_inst.Create(NewTestConfig())
	testEProc = manager_inst.Create(NewEveryConfig())

	resp := &DummyResponder{}
	testResponderProc = manager_inst.Create(NewResponderConfig(resp))

	log.Println("Starting processes.")

	go manager_inst.Run("Test")
	defer testProc.Stop()

	go manager_inst.Run("Every")
	defer testEProc.Stop()

	log.Println("Running tests.")
	val := m.Run()

	log.Println("Shutting down.")
	os.Exit(val)
}

func TestResponder(t *testing.T) {
	good := []events.Event{
		events.NewMessage("test1", nil),
		events.NewMessage("test2", "Yes..."),
	}

	bad := []events.Event{
		events.NewMessage("Invalid", nil),
		events.NewInterrupt(nil),
	}

	t.Run("RespondsTo", func(t *testing.T) {
		t.Run("Unhandled should return false", func(t *testing.T) {
			for _, evt := range bad {
				ret := testResponderProc.RespondsTo(evt)
				if ret {
					t.Fatalf("Responds to bad event: %#v",
						evt)
				}
			}
		})

		t.Run("Handled should return true", func(t *testing.T) {
			for _, evt := range good {
				ret := testResponderProc.RespondsTo(evt)
				if !ret {
					t.Fatalf("Doesn't respond to good event: %#v",
						evt)
				}
			}
		})
	})

	t.Run("Invoke", func(t *testing.T) {
		t.Run("Handled return properly", func(t *testing.T) {
			ret := testResponderProc.Invoke(good[0])
			want := "Ok, closing the pod bay doors"

			rsp, ok := ret.(*events.Response)
			if !ok {
				t.Fatalf("Unexpected response: %#v", ret)
			}

			if rsp.Response() != want {
				t.Errorf("Unexpected result: #%v", rsp)
			}
		})

		t.Run("Unhandled should return error", func(t *testing.T) {
			ret := testResponderProc.Invoke(bad[0])

			if ret != nil {
				t.Errorf("Unexpected result: %#v", ret)
			}
		})
	})
}

// Test `every` repeating processes.
func TestEvery(t *testing.T) {
	time.Sleep(2 * time.Second)
	if EveryVal < 1 {
		t.Errorf("Unexpected value: %#v", EveryVal)
	}
}

// Test process manager.
func TestManager(t *testing.T) {
	pm := NewManager()

	t.Run("Does `Manager.Add` do nothing if given no process?",
		func(t *testing.T) {
			pm.Add(nil)
			if pm.Count() > 0 {
				t.Errorf("Somehow we have %d processes!",
					pm.Count())
			}
		})

	t.Run("Does `Manager.Stop` do nothing if given an invalid process?",
		func(t *testing.T) {
			if pm.Stop("chickens") {
				t.Error("Stopped a non-existing process.")
			}
		})

	t.Run("Does `Manager.Run` do nothing if the process is invalid?",
		func(t *testing.T) {
			if pm.Run("chickens") {
				t.Error("Started an invalid process.")
			}
		})
}

// Test finding invalid processes.
func TestInfalidFind(t *testing.T) {
	_, found := manager_inst.Find("nope")
	if found {
		t.Error("No, found a non-existing process!")
	}
}

// Test finding processes.
func TestFind(t *testing.T) {
	var inst *Process
	var found bool

	t.Run("Finds own instance", func(t *testing.T) {
		inst, found = manager_inst.Find("Test")
		if !found {
			t.Error("Could not find my instance!")
		}
	})

	t.Run("Did we get the *right* process?", func(t *testing.T) {
		if inst != testProc {
			t.Error("Returned process was not ours.")
		}
	})

	t.Run("Is it running?", func(t *testing.T) {
		if !inst.Running() {
			t.Error("Process is not running!")
		}
	})
}

// Test `RunEvery` when process is already running.
func TestEveryAlreadyRunning(t *testing.T) {
	res := testEProc.Run()
	if res {
		t.Error("`Run` returned true.")
	}
}

// Test dumping.
func TestDump(t *testing.T) {
	res := manager_inst.Processes()
	if res != nil && len(res) != 3 {
		t.Errorf("Unexpected process count: %#v", res)
	}
}

// Test stopping processes.
func TestStop(t *testing.T) {
	t.Run("Does `Stop` work as expected?", func(t *testing.T) {
		res := testEProc.Stop()

		if !testEProc.Running() {
			if !res {
				t.Error("Process did not stop properly.")
			}
		} else {
			t.Error("Process did not shut down!")
			return
		}
	})

	t.Run("Does `Stop` return `false` if process not running?", func(t *testing.T) {
		time.Sleep(1 * time.Second)
		res2 := testEProc.Stop()
		if !testEProc.Running() {
			if !res2 {
				t.Log("Yes.")
			} else {
				t.Error("No.")
				return
			}
		} else {
			t.Error("Process did not shut down!")
			return
		}
	})

	t.Run("Does `Manager.StopAll` work as expected?", func(t *testing.T) {
		trueCount := 0

		res := manager_inst.StopAll()
		if len(res) != 3 {
			t.Fatalf("Unexpected results, should be 3: %#v", res)
		}

		for _, val := range res {
			if val {
				trueCount++
			}
		}

		if trueCount != 1 {
			t.Errorf("Unexpectyed result, should be 1: %#v", trueCount)
		}
	})
}

// process_test.go ends here.
