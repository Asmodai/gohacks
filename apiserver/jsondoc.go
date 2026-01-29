// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// jsondoc.go --- JSON documents.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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

package apiserver

import (
	"github.com/gin-gonic/gin"

	"fmt"
	"reflect"
	"time"
)

// JSON document.
type Document struct {
	status  int               `json:"-"`
	headers map[string]string `json:"-"`
	start   time.Time         `json:"-"`

	// JSON document data.
	Data any `json:"data,omitempty"`

	// Number of elements present should `Data` be an array of some kind.
	Count int64 `json:"count"`

	// Error document.
	Error *ErrorDocument `json:"error,omitempty"`

	// Time taken to generate the JSON document.
	Elapsed string `json:"elapsed_time,omitempty"`
}

// Generate a new JSON document.
func NewDocument(status int, data any) *Document {
	var length int64

	if data != nil {
		//nolint:exhaustive
		switch reflect.TypeOf(data).Kind() {
		case reflect.Slice, reflect.Map:
			s := reflect.ValueOf(data)
			length = int64(s.Len())
		}
	}

	return &Document{
		status: status,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		start: time.Now(),
		Data:  data,
		Count: length,
		Error: nil,
	}
}

// Set the `Error` component of the document.
func (d *Document) SetError(err *ErrorDocument) {
	d.Data = nil
	d.status = err.Status
	d.Error = err
}

// Add a header to the document's HTTP response.
func (d *Document) AddHeader(key, value string) {
	d.headers[key] = value
}

// Return the document's HTTP status code response.
func (d *Document) Status() int {
	return d.status
}

// Write the document to the given gin-gonic context.
func (d *Document) Write(ctx *gin.Context) {
	if len(d.headers) > 0 {
		for key, val := range d.headers {
			ctx.Header(key, val)
		}
	}

	t := ctx.GetTime("start_time")
	if !t.IsZero() {
		d.Elapsed = fmt.Sprintf("%v", time.Since(t))
	}

	ctx.JSON(
		d.status,
		d,
	)
}

// jsondoc.go ends here.
