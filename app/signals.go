// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// signals.go --- Signal handler.
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

// * Comments:

// * Package:

package app

// * Imports:

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Asmodai/gohacks/events"
)

// * Code:

// ** Methods:

// Install signal handler.
func (app *application) installSignals() {
	sigs := make(chan os.Signal, 1)

	// We don't care for the following signals:
	signal.Ignore(syscall.SIGURG)

	// Notify when a signal we care for is received.
	signal.Notify(sigs)

	go func() {
		for {
			sig := <-sigs
			event := events.NewSignal(sig)

			if sig != syscall.SIGURG {
				app.Logger().Info(
					"Received signal",
					"signal", sig.String(),
				)
			}

			switch sig {
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM:
				// Handle termination.
				app.Terminate()

				return

			case syscall.SIGHUP:
				// Handle SIGHUP.
				app.responders.SendAll(event)

			case syscall.SIGWINCH:
				// Handle WINCH.
				// Note: Do not bother logging this one.
				app.responders.SendAll(event)

			case syscall.SIGUSR1:
				// Handle user-defined signal #1.
				app.responders.SendAll(event)

			case syscall.SIGUSR2:
				// Handle user-defined signal #2.
				app.responders.SendAll(event)

			case syscall.SIGCHLD:
				// Handle SIGCHLD.
				app.responders.SendAll(event)

			default:
			}
		}
	}()
}

// * signals.go ends here.
