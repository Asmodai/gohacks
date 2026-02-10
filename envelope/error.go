// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// error.go --- Error envelope.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

package envelope

// * Imports:

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Asmodai/gohacks/errx"
)

// * Code:

// ** Type:

type Error struct {
	status  int
	headers http.Header

	Error   error
	Elapsed time.Duration
}

// ** Methods:

func (ee *Error) Status() int          { return ee.status }
func (ee *Error) Headers() http.Header { return ee.headers }
func (ee *Error) Body() any            { return ee }

// *** JSON marshaller:

type marshalError struct {
	Error   string         `json:"error,omitempty"`
	Cause   string         `json:"cause,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

type marshalErrorStruct struct {
	Error   *marshalError `json:"error,omitempty"`
	Elapsed time.Duration `json:"elapsed_ns,omitempty"`
}

func (ee *Error) MarshalJSON() ([]byte, error) {
	if ee == nil {
		return []byte{}, nil
	}

	tmp := &marshalErrorStruct{
		Error: &marshalError{
			Error:   ee.Error.Error(),
			Cause:   errx.Cause(ee.Error).Error(),
			Details: errx.AllDetails(ee.Error),
		},
		Elapsed: ee.Elapsed,
	}

	result, err := json.Marshal(tmp)

	return result, errx.WithStack(err)
}

// ** Functions:

func NewError(status int, err error) *Error {
	return &Error{
		status:  status,
		headers: defaultEnvelopeHeaders(),
		Error:   err,
		Elapsed: time.Duration(0),
	}
}

// * error.go ends here.
