/*
 * error.go --- Better Errors(tm)
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
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

package types

import (
	"github.com/goccy/go-json"

	"fmt"
	"log"
)

/*

Custom error structure.

This is compatible with the `error` interface and provides `Unwrap`
support.

*/
type Error struct {
	Module  string
	Message string
}

// Create a new error object.
func NewError(module string, format string, args ...interface{}) *Error {
	msg := fmt.Sprintf(format, args...)

	return &Error{
		Module:  module,
		Message: msg,
	}
}

// Create a new error object and immediately log it.
func NewErrorAndLog(module string, format string, args ...interface{}) *Error {
	err := NewError(module, format, args...)

	err.Log()

	return err
}

// Return a human-readable string representation of the error.
func (e *Error) String() string {
	return e.Error()
}

// Return a human-readable string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Module, e.Message)
}

// Unwrap the error.
func (e *Error) Unwrap() error {
	return fmt.Errorf(e.Message)
}

// Log the error.
func (e *Error) Log() {
	log.Printf("%s: %s", e.Module, e.Message)
}

// Convert the error to a JSON string.
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Error())
}

/* error.go ends here. */
