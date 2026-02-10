// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// success.go --- Success envelope.
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
	"fmt"
	"net/http"
	"time"

	"github.com/Asmodai/gohacks/errx"
)

// * Code:

// ** Type:

type Success struct {
	status  int
	headers http.Header

	Data    any
	Count   int64
	Elapsed time.Duration
}

// ** Methods:

func (se *Success) Status() int          { return se.status }
func (se *Success) Headers() http.Header { return se.headers }
func (se *Success) Body() any            { return se }

// *** JSON marshaller:

type marshalSuccessStruct struct {
	Data    any    `json:"data,omitempty"`
	Count   int64  `json:"count,omitempty"`
	Elapsed string `json:"elapsed,omitempty"`
}

func (se *Success) MarshalJSON() ([]byte, error) {
	if se == nil {
		return []byte{}, nil
	}

	tmp := &marshalSuccessStruct{Data: se.Data}

	if se.Count > 0 {
		tmp.Count = se.Count
	}

	if se.Elapsed > 0 {
		tmp.Elapsed = fmt.Sprintf("%v", se.Elapsed)
	}

	result, err := json.Marshal(tmp)

	return result, errx.WithStack(err)
}

// ** Functions:

func NewSuccess(status int, data any) *Success {
	return &Success{
		status:  status,
		headers: defaultEnvelopeHeaders(),
		Data:    data,
		Count:   getContainerLength(data),
		Elapsed: time.Duration(0),
	}
}

// * success.go ends here.
