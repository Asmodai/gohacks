/*
 * process_test.go --- Process state tests.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package process

import (
	"log"
	"os"
	"testing"
	"time"
)

var (
	manager    *Manager
	testProc   *Process
	testNBProc *Process
	testBProc  *Process
	testEProc  *Process

	testBlockingSend bool
	fromNonblocking  interface{}
	fromBlocking     interface{}

	EveryVal int = 0
)

// Test `blocking` function.
func BlockingFn(state **State) {
	var ps *State = *state

	fromBlocking = ps.ReceiveBlocking()

	time.Sleep(1 * time.Second)
	ps.Send(fromBlocking)

	// Sneakily test the `default` clause.
	testBlockingSend = ps.Send("Nope")
}

// Test `Nonblocking` function.
func NonblockingFn(state **State) {
	var ps *State = *state

	data, ok := ps.Receive()
	if ok {
		fromNonblocking = data
		ps.SendBlocking(data)
	}
}

// Test `every` function.
func EveryFn(state **State) {
	EveryVal++
}

// Create a new test config.
func NewTestConfig() *Config {
	cnf := NewDefaultConfig()

	cnf.Name = "Test"
	cnf.Interval = 0
	cnf.Function = nil

	return cnf
}

// Create a new config for blocking send test.
func NewBlockingConfig() *Config {
	return &Config{
		Name:     "Blocking",
		Interval: 0,
		Function: BlockingFn,
	}
}

// Create a new config for nonblocking send test.
func NewNonblockingConfig() *Config {
	return &Config{
		Name:     "Nonblocking",
		Interval: 0,
		Function: NonblockingFn,
	}
}

// Create a new config for the `every` event test.
func NewEveryConfig() *Config {
	return &Config{
		Name:     "Every",
		Interval: 1,
		Function: EveryFn,
	}
}

// Main testing function.
func TestMain(m *testing.M) {
	log.Println("Setting up processes.")

	manager = NewManager()
	testProc = manager.Create(NewTestConfig())
	testNBProc = manager.Create(NewNonblockingConfig())
	testBProc = manager.Create(NewBlockingConfig())
	testEProc = manager.Create(NewEveryConfig())

	log.Println("Starting processes.")

	go manager.Run("Test")
	defer testProc.Stop()

	go manager.Run("Blocking")
	defer testBProc.Stop()

	go manager.Run("Nonblocking")
	defer testNBProc.Stop()

	go manager.Run("Every")
	defer testEProc.Stop()

	log.Println("Running tests.")
	val := m.Run()

	log.Println("Shutting down.")
	os.Exit(val)
}

// Test Nonblocking functions.
func TestNonblocking(t *testing.T) {
	t.Log("Does non-blocking `Receive` work as expected?")

	testNBProc.Send("test")
	time.Sleep(1 * time.Second)

	t.Logf("Function got '%v'", fromNonblocking)
	if fromNonblocking == "test" {
		time.Sleep(1 * time.Second)
		result := testNBProc.Receive()

		t.Logf("Process got '%v'", result)
		if result == "test" {
			t.Log("Expected result.")
		} else {
			t.Error("Unexpected result.")
			return
		}
	} else {
		t.Error("Unexpected result.")
		return
	}

	t.Log("Did the second `send` fail because of a channel buffer?")
	if !testBlockingSend {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

// Test blocking functions.
func TestBlocking(t *testing.T) {
	t.Log("Does blocking `Receive` work as expected?")

	testBProc.Send("test")
	time.Sleep(1 * time.Second)

	t.Logf("Function got '%v'", fromBlocking)
	if fromBlocking == "test" {
		t.Log("Expected result.")
		return
	}

	t.Error("Unexpected result.")
}

// Test `every` repeating processes.
func TestEvery(t *testing.T) {
	t.Log("Does `Every` fire as expected?")

	time.Sleep(2 * time.Second)
	if EveryVal > 1 {
		t.Log("Yes.")
	} else {
		t.Error("No.")
	}
}

// Test process manager.
func TestManager(t *testing.T) {
	pm := NewManager()

	t.Log("Does `Manager.Add` do nothing if given no process?")
	pm.Add(nil)
	if pm.Count() > 0 {
		t.Errorf("Somehow we have %d processes!", pm.Count())
		return
	}
	t.Log("Yes.")

	t.Log("Does `Manager.Stop` do nothing if given an invalid process?")
	if !pm.Stop("chickens") {
		t.Log("Yes.")
	} else {
		t.Error("No.")
		return
	}

	t.Log("Does `Manager.Run` do nothing if the process is invalid?")
	if !pm.Run("chickens") {
		t.Log("Yes.")
	} else {
		t.Error("No.")
	}
}

// Test finding invalid processes.
func TestInfalidFind(t *testing.T) {
	t.Log("Does `Manager.Find` do the right thing when no process is found?`")

	_, found := manager.Find("nope")
	if found {
		t.Error("No, found a non-existing process!")
		return
	}

	t.Log("Yes.")
}

// Test finding processes.
func TestFind(t *testing.T) {
	t.Log("Can I find my instance?")
	i, found := manager.Find("Test")
	if !found {
		t.Error("Could not find my instance!")
		return
	}
	t.Log("Yes.")

	t.Log("Did we get the *right* process?")
	if i != testProc {
		t.Error("Returned process was not ours.")
		return
	}
	t.Log("Yes.")

	t.Log("Is it running?")
	if !i.Running {
		t.Error("Process is not running!")
		return
	}
	t.Log("Yes.")
}

// Test `RunEvery` when process is already running.
func TestEveryAlreadyRunning(t *testing.T) {
	t.Log("Does `RunEvery` return `false` if already running?")

	res := testEProc.Run()
	if testEProc.Running {
		if res {
			t.Error("No.")
			return
		}

		t.Log("Yes.")
		return
	}

	t.Error("Process was not running!")
}

// Test dumping.
func TestDump(t *testing.T) {
	t.Log("Can we list processes?")

	res := manager.Processes()
	if res != nil {
		if len(*res) == 4 {
			t.Log("Yes.")
			return
		}
	}

	t.Error("No.")
}

// Test stopping processes.
func TestStop(t *testing.T) {
	t.Log("Does `Stop` work as expected?")

	res1 := testNBProc.Stop()
	if !testNBProc.Running {
		if res1 {
			t.Log("Yes.")
		} else {
			t.Error("No.")
			return
		}
	} else {
		t.Error("Process did not shut down!")
		return
	}

	t.Log("Does `Stop` return `false` if process not running?")
	time.Sleep(1 * time.Second)
	res2 := testNBProc.Stop()
	if !testNBProc.Running {
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

	t.Log("Does `Manager.StopAll` work as expected?")
	res3 := manager.StopAll()
	if res3 {
		t.Log("Yes.")
	} else {
		t.Error("No.")
	}
}

/* process_test.go ends here. */
