/*
 * error.go --- Better Errors(tm)
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package types

import (
	"encoding/json"
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
