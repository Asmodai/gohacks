// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// logger.go --- Logger interface.
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
//
// mock:yes

package logger

/*
Logging structure.

To use,

1) Create a logger:

```go

	lgr := logger.NewLogger()

```

2) Do things with it:

```go

	lgr.Warn("Not enough coffee!")
	lgr.Info("Water is heating up.")
	// and so on.

```

If an empty string is passed to `NewLogger`, then the log facility will
display messages on standard output.
*/
type Logger interface {
	// Set whether the logger prints in human-readable 'debug' output or
	// machine-readable JSON format.
	SetDebug(bool)

	// Set the log file to which the logger will write.
	SetLogFile(string)

	// Log a Go error.
	GoError(error, ...any)

	// Log a debug message.
	Debug(string, ...any)

	// Log an error message.
	Error(string, ...any)

	// Log a warning message.
	Warn(string, ...any)

	// Log an informational message.
	Info(string, ...any)

	// Log a fatal error condition and exit.
	Fatal(string, ...any)

	// Log a panic condition and exit.
	Panic(string, ...any)

	// Log a debug message with printf formatting.
	Debugf(string, ...any)

	// Log an error message with printf formatting.
	Errorf(string, ...any)

	// Log a warning message with printf formatting.
	Warnf(string, ...any)

	// Log an informational message with printf formatting.
	Infof(string, ...any)

	// Log a fatal message with printf formatting.
	Fatalf(string, ...any)

	// Log a panic message with printf formatting.
	Panicf(string, ...any)

	// Encapsulate user-specified metadata fields.
	WithFields(Fields) Logger
}

// logger.go ends here.
