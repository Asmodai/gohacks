/*
 * mocklogger.go --- Mock logger for tests.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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

package logger

import (
	"fmt"
	"testing"
)

/*

Mock logger for Go testing framework.

To use, be sure to set `Test` to your test's `testing.T` instance.

*/
type MockLogger struct {
	Test      *testing.T
	LastFatal string
}

// Create a new mock logger.
func NewMockLogger(_ string) *MockLogger {
	return &MockLogger{}
}

// Set debug mode.
func (l *MockLogger) SetDebug(junk bool) {
}

// Set the log file to use.
func (l *MockLogger) SetLogFile(junk string) {
}

// Write a debug message to the log.
func (l *MockLogger) Debug(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	l.Test.Logf("DEBUG: %s  %v", msg, rest)
}

// Write an error message to the log.
func (l *MockLogger) Error(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	l.Test.Logf("ERROR:  %s  %v", msg, rest)
}

// Write a warning message to the log.
func (l *MockLogger) Warn(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	l.Test.Logf("WARN:  %s  %v", msg, rest)
}

// Write an information message to the log.
func (l *MockLogger) Info(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	l.Test.Logf("INFO:  %s  %v", msg, rest)
}

// Write a fatal message to the log and then exit.
func (l *MockLogger) Fatal(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	// Don't ever use `Fatal` or whatever here, we don't want to exit.
	l.LastFatal = fmt.Sprintf("FATAL: %s  %v", msg, rest)
	l.Test.Log(l.LastFatal)
}

// Write a debug message to the log.
func (l *MockLogger) Debugf(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	l.Test.Logf("DEBUG: %s", fmt.Sprintf(msg, rest...))
}

// Write a warning message to the log.
func (l *MockLogger) Warnf(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	l.Test.Logf("WARN:  %s", fmt.Sprintf(msg, rest...))
}

// Write an information message to the log.
func (l *MockLogger) Infof(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	l.Test.Logf("INFO:  %s", fmt.Sprintf(msg, rest...))
}

// Write a fatal message to the log and then exit.
func (l *MockLogger) Fatalf(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	// Don't ever use `Fatal` or whatever here, we don't want to exit.
	l.LastFatal = fmt.Sprintf("FATAL: %s", fmt.Sprintf(msg, rest...))
	l.Test.Log(l.LastFatal)
}

// Write a fatal message to the log and then exit.
func (l *MockLogger) Errorf(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	// Don't ever use `Fatal` or whatever here, we don't want to exit.
	l.LastFatal = fmt.Sprintf("ERROR: %s", fmt.Sprintf(msg, rest...))
	l.Test.Log(l.LastFatal)
}

// Write a fatal message to the log and then exit.
func (l *MockLogger) Panicf(msg string, rest ...interface{}) {
	if l.Test == nil {
		return
	}

	// Don't ever use `Fatal` or whatever here, we don't want to exit.
	l.LastFatal = fmt.Sprintf("PANIC: %s", fmt.Sprintf(msg, rest...))
	l.Test.Log(l.LastFatal)
}

func (l *MockLogger) WithFields(_ Fields) ILogger {
	return NewMockLogger("")
}

/* mocklogger.go ends here. */
