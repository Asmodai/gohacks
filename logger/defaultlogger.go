/*
 * logger.go --- Default logger.
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
	"log"
)

/*

Default logging structure.

This is a simple implementation of the `ILogger` interface that simply redirects
messages to `log.Printf`.

It is used in the same way as the main `Logger` implementation.

*/
type DefaultLogger struct {
}

// Create a new default logger.
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{}
}

// Set debug mode.
func (l *DefaultLogger) SetDebug(junk bool) {
}

// Set the log file to use.
func (l *DefaultLogger) SetLogFile(junk string) {
}

// Write a debug message to the log.
func (l *DefaultLogger) Debug(msg string, rest ...interface{}) {
	log.Printf("DEBUG: %s  %v", msg, rest)
}

// Write an error message to the log.
func (l *DefaultLogger) Error(msg string, rest ...interface{}) {
	log.Printf("ERROR:  %s  %v", msg, rest)
}

// Write a warning message to the log.
func (l *DefaultLogger) Warn(msg string, rest ...interface{}) {
	log.Printf("WARN:  %s  %v", msg, rest)
}

// Write an information message to the log.
func (l *DefaultLogger) Info(msg string, rest ...interface{}) {
	log.Printf("INFO:  %s  %v", msg, rest)
}

// Write a fatal message to the log and then exit.
func (l *DefaultLogger) Fatal(msg string, rest ...interface{}) {
	log.Fatalf("FATAL: %s  %v", msg, rest)
}

// Write a debug message to the log.
func (l *DefaultLogger) Debugf(format string, args ...interface{}) {
	log.Printf("DEBUG: %s", fmt.Sprintf(format, args...))
}

// Write a warning message to the log.
func (l *DefaultLogger) Warnf(msg string, args ...interface{}) {
	log.Printf("WARN:  %s", fmt.Sprintf(msg, args...))
}

// Write an information message to the log.
func (l *DefaultLogger) Infof(msg string, args ...interface{}) {
	log.Printf("INFO:  %s", fmt.Sprintf(msg, args...))
}

// Write a fatal message to the log and then exit.
func (l *DefaultLogger) Fatalf(msg string, args ...interface{}) {
	log.Fatalf("FATAL: %s", fmt.Sprintf(msg, args...))
}

// Write an error message to the log and then exit.
func (l *DefaultLogger) Errorf(msg string, args ...interface{}) {
	log.Fatalf("ERROR: %s", fmt.Sprintf(msg, args...))
}

// Write a fatal message to the log and then exit.
func (l *DefaultLogger) Panicf(msg string, args ...interface{}) {
	log.Fatalf("PANIC: %s", fmt.Sprintf(msg, args...))
}

func (l *DefaultLogger) WithFields(_ Fields) ILogger {
	return NewDefaultLogger()
}

/* logger.go ends here. */
