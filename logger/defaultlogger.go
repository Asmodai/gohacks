// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// defaultlogger.go --- Default logger.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package logger

import (
	"gitlab.com/tozd/go/errors"

	"fmt"
	"log"
)

/*
Default logging structure.

This is a simple implementation of the `ILogger` interface that simply
redirects messages to `log.Printf`.

It is used in the same way as the main `Logger` implementation.
*/
type defaultLogger struct {
}

// Create a new default logger.
func NewDefaultLogger() Logger {
	return &defaultLogger{}
}

// Set debug mode.
func (l *defaultLogger) SetDebug(_ bool) {
}

// Set the log file to use.
func (l *defaultLogger) SetLogFile(_ string) {
}

// Write a Go error to the log.
func (l *defaultLogger) GoError(err error, rest ...any) {
	nerr, ok := err.(errors.E) //nolint:errorlint
	if !ok {
		goto done
	}

	if cause := errors.Cause(nerr); cause != nil {
		args := []any{"cause", cause}
		args = append(args, rest...)
		l.Error(err.Error(), args...)

		return
	}

done:
	l.Error(err.Error(), rest...)
}

// Write a debug message to the log.
func (l *defaultLogger) Debug(msg string, rest ...any) {
	log.Printf("DEBUG: %s  %v", msg, rest)
}

// Write an error message to the log.
func (l *defaultLogger) Error(msg string, rest ...any) {
	log.Printf("ERROR:  %s  %v", msg, rest)
}

// Write a warning message to the log.
func (l *defaultLogger) Warn(msg string, rest ...any) {
	log.Printf("WARN:  %s  %v", msg, rest)
}

// Write an information message to the log.
func (l *defaultLogger) Info(msg string, rest ...any) {
	log.Printf("INFO:  %s  %v", msg, rest)
}

// Write a fatal message to the log and then exit.
func (l *defaultLogger) Fatal(msg string, rest ...any) {
	log.Fatalf("FATAL: %s  %v", msg, rest)
}

// Write a Panic message to the log and then exit.
func (l *defaultLogger) Panic(msg string, rest ...any) {
	log.Fatalf("PANIC: %s  %v", msg, rest)
}

// Write a debug message to the log.
func (l *defaultLogger) Debugf(format string, args ...any) {
	log.Printf("DEBUG: %s", fmt.Sprintf(format, args...))
}

// Write an error message to the log and then exit.
func (l *defaultLogger) Errorf(msg string, args ...any) {
	log.Fatalf("ERROR: %s", fmt.Sprintf(msg, args...))
}

// Write a warning message to the log.
func (l *defaultLogger) Warnf(msg string, args ...any) {
	log.Printf("WARN:  %s", fmt.Sprintf(msg, args...))
}

// Write an information message to the log.
func (l *defaultLogger) Infof(msg string, args ...any) {
	log.Printf("INFO:  %s", fmt.Sprintf(msg, args...))
}

// Write a fatal message to the log and then exit.
func (l *defaultLogger) Fatalf(msg string, args ...any) {
	log.Fatalf("FATAL: %s", fmt.Sprintf(msg, args...))
}

// Write a fatal message to the log and then exit.
func (l *defaultLogger) Panicf(msg string, args ...any) {
	log.Fatalf("PANIC: %s", fmt.Sprintf(msg, args...))
}

func (l *defaultLogger) WithFields(_ Fields) Logger {
	return NewDefaultLogger()
}

// defaultlogger.go ends here.
