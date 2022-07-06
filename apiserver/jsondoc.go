/*
 * jsondoc.go --- JSON documents.
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

package apiserver

import (
	"github.com/gin-gonic/gin"

	"fmt"
	"time"
)

type Document struct {
	status  int               `json:"-"`
	headers map[string]string `json:"-"`
	start   time.Time         `json:"-"`

	Data    interface{}    `json:"data"`
	Error   *ErrorDocument `json:"error"`
	Elapsed string         `json:"elapsed_time",omitempty`
}

func NewDocument(status int, data interface{}) *Document {
	return &Document{
		status: status,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
		start: time.Now(),
		Data:  data,
		Error: nil,
	}
}

func (d *Document) SetError(err *ErrorDocument) {
	d.Data = nil
	d.status = err.Status
	d.Error = err
}

func (d *Document) AddHeader(key, value string) {
	d.headers[key] = value
}

func (d *Document) Status() int {
	return d.status
}

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

/* jsondoc.go ends here. */
