/*
 * process_test.go --- Process state tests.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package process

import (
	"context"
	"log"
	"os"
	"testing"
	"time"
)

var (
	manager       *Manager
	testRunProc   *Process
	testEveryProc *Process
	testMessage   interface{}
)

func testAction(ctx context.Context, ipc *IPC) {
	if val, ok := ipc.Receive(); ok {
		switch val {
		case "test":
			// This test does not send anything back via the mailbox.
			testMessage = val
			break

		default:
			// Anything else will write a string to the mailbox.
			ipc.Send("result")
			break
		}
	}
}

func testStop() {
	log.Printf("Process is terminating due to context cancel.\n")
}

func newRunConfig() *Config {
	return &Config{
		Name:     "Test Run Process",
		Interval: 0,
		ActionFn: testAction,
		StopFn:   testStop,
	}
}

func newEveryConfig() *Config {
	return &Config{
		Name:     "Test Every Process",
		Interval: 1,
		ActionFn: testAction,
		StopFn:   testStop,
	}
}

// ==================================================================
// {{{ Setup:

func TestMain(m *testing.M) {
	log.Println("Setting up processes.")
	{
		manager = NewManager()
		testRunProc = manager.Create(newRunConfig())
		testEveryProc = manager.Create(newEveryConfig())
	}

	log.Println("Starting processes.")
	{
		go manager.Run("Test Run Process")
		go manager.Run("Test Every Process")

		// Allow the processes to start.
		time.Sleep(2 * time.Second)
	}

	log.Println("Running tests.")
	val := m.Run()

	log.Println("Shutting down.")
	{
		testRunProc.Stop()
		testEveryProc.Stop()
	}

	os.Exit(val)
}

// }}}
// ==================================================================

// ==================================================================
// {{{ Regular non-interval process:

func TestRunProc(t *testing.T) {
	t.Run("Is it running?", func(t *testing.T) {
		if testRunProc.Running() == false {
			t.Error("Process is *not* running.")
			return
		}
	})

	t.Run("Can we write to its mailbox?", func(t *testing.T) {
		payload := "test"

		testMessage = nil
		testRunProc.Send(payload)
		time.Sleep(1 * time.Second)

		if testMessage == nil {
			t.Error("Mailbox Send failed!")
			return
		}

		if payload != testMessage.(string) {
			t.Errorf("Unexpected result '%v'", testMessage)
			return
		}
	})

	t.Run("Can we read from its mailbox?", func(t *testing.T) {
		payload := "read"

		testMessage = nil
		testRunProc.Send(payload)
		time.Sleep(2 * time.Second)

		val, ok := testRunProc.Receive()
		if !ok {
			t.Error("Mailbox Get failed!")
			return
		}

		if val.(string) != "result" {
			t.Errorf("Unexpected result '%v'", val)
			return
		}
	})
}

// }}}
// ==================================================================

// ==================================================================
// {{{ Interval process:

// }}}
// ==================================================================

// ==================================================================
// {{{ Process manager:

// Test main manager functions.
func TestManager(t *testing.T) {
	t.Run("Does `Add` do nothing if given nothing?", func(t *testing.T) {
		manager.Add(nil)
		if manager.Count() > 2 {
			t.Errorf("Somehow we have %d processes!", manager.Count())
		}
	})

	t.Run("Does `Stop` do nothing if given invalid process?", func(t *testing.T) {
		if manager.Stop("chickens") {
			t.Error("Manager reports success for stopping invalid process!")
		}
	})

	t.Run("Does `Run` do nothing if given invalid process?", func(t *testing.T) {
		if manager.Run("chickens") {
			t.Error("Manager reports success for running invalid process!")
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("Does it return false for invalid process?", func(t *testing.T) {
		if _, found := manager.Find("chickens"); found {
			t.Error("Manager claims to be able to find invalid process!")
		}
	})

	t.Run("Does it work as expected?", func(t *testing.T) {
		inst, found := manager.Find("Test Run Process")

		if !found {
			t.Error("Could not find test process!")
			return
		}

		if inst != testRunProc {
			t.Error("Found wrong instance!")
		}

		if !inst.Running() {
			t.Error("Found process is not running!")
		}
	})
}

// Test utility functions.
func TestUtils(t *testing.T) {
	t.Run("Can we list processes?", func(t *testing.T) {
		res := manager.Processes()

		if res == nil {
			t.Error("No process list was returned!")
			return
		}

		if len(*res) != 2 {
			t.Errorf("Incorrect number of processes: %d", len(*res))
			return
		}
	})
}

// Test process killing abilities.
func TestStop(t *testing.T) {
	t.Run("Does it work as expected?", func(t *testing.T) {
		res := testRunProc.Stop()
		time.Sleep(100 * time.Millisecond)

		if testRunProc.Running() {
			t.Error("Process reports as still running.")
		}

		if !res {
			t.Error("Process was already stopped.")
		}
	})

	t.Run("Does it return `false` if already stopped?", func(t *testing.T) {
		res := testRunProc.Stop()
		time.Sleep(100 * time.Millisecond)

		if testRunProc.Running() {
			t.Error("Process reports as still running.")
		}

		if res {
			t.Error("Process could not be stopped apparently.")
		}
	})

	t.Run("Does `Manager.StopAll` work?", func(t *testing.T) {
		if res := manager.StopAll(); !res {
			t.Error("`StopAll` failed!")
		}

		// Sleep here to allow system to catch up.
		time.Sleep(100 * time.Millisecond)
	})
}

// }}}
// ==================================================================

/* process_test.go ends here. */
