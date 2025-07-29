// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// app_test.go --- Application tests.
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

//
//
//

// * Package:

package app

// * Imports:

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"syscall"
	"testing"

	"github.com/Asmodai/gohacks/config"
	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/logger"
	mlogger "github.com/Asmodai/gohacks/mocks/logger"
	mprocess "github.com/Asmodai/gohacks/mocks/process"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/semver"
	"go.uber.org/mock/gomock"
)

// * Constants:

// * Variables:

// * Code:

// ** Types:

type CallbackTester struct {
	DidNaughty  bool
	DidOnStart  bool
	DidOnExit   bool
	DidMainLoop bool
}

func (cb *CallbackTester) BeNaughty(_ Application) { cb.DidNaughty = true }
func (cb *CallbackTester) OnStart(_ Application)   { cb.DidOnStart = true }
func (cb *CallbackTester) OnExit(_ Application)    { cb.DidOnExit = true }

func (cb *CallbackTester) MainLoop(a Application) {
	log.Println("In main loop...")
	cb.DidMainLoop = true

	a.Terminate()
}

func (cb *CallbackTester) Reset() {
	cb.DidNaughty = false
	cb.DidOnStart = false
	cb.DidOnExit = false
	cb.DidMainLoop = false
}

type ResponderTester struct {
	name string
	typ  string
	log  *[]string
}

func (d *ResponderTester) Name() string { return d.name }
func (d *ResponderTester) Type() string { return d.typ }

func (d *ResponderTester) RespondsTo(_ events.Event) bool {
	return true
}

func (d *ResponderTester) Invoke(evt events.Event) events.Event {
	if d.log != nil {
		*d.log = append(
			*d.log,
			fmt.Sprintf(
				"Received %s",
				reflect.TypeOf(evt).Elem().Name(),
			),
		)
	}

	return evt
}

// ** Tests:

func TestApp(t *testing.T) {
	var (
		err    error
		theapp Application
	)

	// Oh this is potentially very painful.
	os.Args = []string{os.Args[0]}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	controller := gomock.NewController(t)
	defer controller.Finish()

	ver, err := semver.MakeSemVer("10020003:commit")
	if err != nil {
		t.Fatalf("Could not make version: %#v", err)
	}

	lgr := mlogger.NewMockLogger(controller)
	lgr.EXPECT().SetLogFile(gomock.Any()).AnyTimes()
	lgr.EXPECT().SetDebug(gomock.Any()).AnyTimes()
	lgr.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	lgr.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()

	pmgr := mprocess.NewMockManager(controller)
	pmgr.EXPECT().SetContext(gomock.Any()).AnyTimes()
	pmgr.EXPECT().SetLogger(gomock.Any()).AnyTimes()
	pmgr.EXPECT().StopAll().AnyTimes()

	ctx, err = logger.SetLogger(ctx, lgr)
	if err != nil {
		t.Fatalf("Error when setting logger DI: %#v", err)
	}

	ctx, err = process.SetManager(ctx, pmgr)
	if err != nil {
		t.Fatalf("Error when setting process manager DI: %#v", err)
	}

	cbtests := &CallbackTester{}

	cnf := &Config{
		Name:       "Test App",
		Version:    ver,
		Validators: config.ValidatorsMap{},
		RequireCLI: false,
	}

	t.Run("Constructs", func(t *testing.T) {
		theapp = NewApplication(cnf)
		theapp.ParseConfig()

		if theapp == nil {
			t.Fatal("Could not construct application!")
		}

		if theapp.Name() != cnf.Name {
			t.Errorf("Name mismatch: %#v != %#v",
				theapp.Name(),
				cnf.Name)
		}
	})

	t.Run("Init()", func(t *testing.T) {
		theapp.SetContext(ctx)
		theapp.SetOnStart(cbtests.OnStart)
		theapp.SetOnExit(cbtests.OnExit)
		theapp.SetMainLoop(cbtests.MainLoop)

		theapp.Init()

		if !theapp.(*application).initialised.Load() {
			t.Fatal("Application did not initialise!")
		}

		t.Run("Cannot set things after", func(t *testing.T) {
			theapp.SetOnStart(nil)

			if theapp.(*application).onStart == nil {
				t.Fatalf(
					"Able to set onStart: %#v",
					theapp.(*application).onStart,
				)
			}
		})
	})

	t.Run("Accessors", func(t *testing.T) {
		t.Run("Name()", func(t *testing.T) {
			if theapp.Name() != cnf.Name {

				t.Errorf("%#v != %#v", theapp.Name(), cnf.Name)
			}
		})

		t.Run("Version()", func(t *testing.T) {
			if theapp.Version() != ver {
				t.Errorf("%#v != %#v", theapp.Version(), ver)
			}
		})

		t.Run("Commit()", func(t *testing.T) {
			want := "commit"
			if theapp.Commit() != want {
				t.Errorf("%#v != %#v", theapp.Commit(), want)
			}
		})

		t.Run("Context()", func(t *testing.T) {
			want := theapp.(*application).ctx
			if theapp.Context() != want {
				t.Errorf("%#v != %#v", theapp.Context(), want)
			}
		})

		t.Run("ProcessManager()", func(t *testing.T) {
			if theapp.ProcessManager() != pmgr {
				t.Errorf("%#v != %#v", theapp.ProcessManager(), pmgr)
			}
		})

		t.Run("Logger()", func(t *testing.T) {
			if theapp.Logger() != lgr {
				t.Errorf("%#v != %#v", theapp.Logger(), lgr)
			}
		})

		t.Run("IsDebug()", func(t *testing.T) {
			if theapp.IsDebug() != false {
				t.Errorf("%#v != %#v", theapp.IsDebug(), false)
			}
		})

		t.Run("IsRunning()", func(t *testing.T) {
			if theapp.IsRunning() != false {
				t.Errorf("%#v != %#v", theapp.IsRunning(), false)
			}
		})
	})

	t.Run("Starting", func(t *testing.T) {
		theapp.Run()

		// Should terminate in main loop.
		if cbtests.DidOnStart == false {
			t.Error("Did not execute `onStart`.")
		}

		if cbtests.DidOnExit == false {
			t.Error("Did not execute `onExit'.")
		}

		if cbtests.DidMainLoop == false {
			t.Error("Did not execute `mainLoop'.")
		}

		cbtests.Reset()
	})

	t.Run("Responder", func(t *testing.T) {
		log1 := []string{}
		log2 := []string{}

		r1 := &ResponderTester{"one", "test", &log1}
		r2 := &ResponderTester{"two", "test", &log2}

		t.Run("AddResponder()", func(t *testing.T) {
			var err error

			if _, err = theapp.AddResponder(r1); err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}

			if _, err = theapp.AddResponder(r2); err != nil {
				t.Fatalf("Unexpected error: %#v", err)
			}
		})

		t.Run("SendAllResponders()", func(t *testing.T) {
			want := "Received Message"
			evt := events.NewMessage("testing", "Spawn more overlords!")

			theapp.SendAllResponders(evt)

			t.Run("Received by 1", func(t *testing.T) {
				if len(log1) != 1 {
					t.Fatalf("Length mismatch: %d", len(log1))
				}

				if log1[0] != want {
					t.Errorf("Event mismatch: %s", log1[0])
				}
			})

			t.Run("Received by 2", func(t *testing.T) {
				if len(log2) != 1 {
					t.Fatalf("Length mismatch: %d", len(log2))
				}

				if log2[0] != want {
					t.Errorf("Event mismatch: %s", log2[0])
				}
			})
		})

		t.Run("SendFirstResponder()", func(t *testing.T) {
			want := "Received Message"
			evt := events.NewMessage("testing", "Spawn more overlords!")

			theapp.SendFirstResponder(evt)

			t.Logf("Log1: %#v", log1)
			t.Logf("Log2: %#v", log2)

			t.Run("Received by 1", func(t *testing.T) {
				if len(log1) != 2 {
					t.Fatalf("Length mismatch: %d", len(log1))
				}

				if log1[1] != want {
					t.Errorf("Event mismatch: %s", log1[1])
				}
			})

			t.Run("Received by 2", func(t *testing.T) {
				if len(log2) == 0 {
					t.Fatalf("Length mismatch: %d", len(log2))
				}

				if log2[0] != want {
					t.Errorf("Event mismatch: %s", log2[0])
				}
			})
		})
	})

	t.Run("Signals", func(t *testing.T) {
		var sresp *SignalResponder

		thesignal := syscall.SIGUSR1
		goodevt := events.NewSignal(thesignal)
		badevt := events.NewMessage("cheese", "Nope")

		t.Run("Constructs", func(t *testing.T) {
			sresp = NewSignalResponder()

			if sresp == nil {
				t.Fatal("Did not construct!")
			}

			if sresp.callback == nil {
				t.Error("No default callback.")
			}
		})

		t.Run("RespondsTo()", func(t *testing.T) {
			if ok := sresp.RespondsTo(goodevt); !ok {
				t.Error("Responder doesn't respond to good event.")
			}

			if ok := sresp.RespondsTo(badevt); ok {
				t.Error("Responder responds to bad events.")
			}
		})

		t.Run("Receives", func(t *testing.T) {
			var (
				called bool
				with   os.Signal
			)

			sresp.SetOnSignal(func(sig os.Signal) {
				called = true
				with = sig
			})

			theapp.AddResponder(sresp)
			theapp.(*application).responders.SendAll(goodevt)

			if !called {
				t.Fatal("Responder was not invoked")
			}

			if with != thesignal {
				t.Errorf("Result mismatch: %#v != %#v", with, thesignal)
			}
		})
	})
}

// * app_test.go ends here.
