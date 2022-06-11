/*
 * jsondoc.go --- JSON documents.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
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
