/*
 * mocklogger.go --- Mock logger for tests.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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
