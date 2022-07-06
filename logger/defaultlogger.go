/*
 * logger.go --- Default logger.
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
