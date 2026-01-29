// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// logger_file.go --- File logging support.
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
	"io"
	"os"
	"time"
)

const (
	// Default log file mode used when creating log files.
	defaultFileMode = 0644
)

// Initialise the logging component.
func (d *Dispatcher) initLog() gin.LoggerConfig {
	return gin.LoggerConfig{
		Formatter: d.logFormatter,
		Output:    d.logWriter(),
	}
}

// Create a log writer instance.
func (d *Dispatcher) logWriter() io.Writer {
	if d.config.LogFile == "" || gin.Mode() == "debug" {
		return gin.DefaultWriter
	}

	file, err := os.OpenFile(
		d.config.LogFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		defaultFileMode,
	)
	if err != nil {
		d.lgr.Fatal(
			"Could not open file for writing.",
			"file", d.config.LogFile,
			"err", err.Error(),
		)
	}

	d.lgr.Info(
		"API server logging initialised.",
		"file", d.config.LogFile,
	)

	return file
}

// Log entry formatter.
func (d *Dispatcher) logFormatter(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string

	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	return fmt.Sprintf("%s - [%s] \"%s%s%s %s %s\" %s%d%s %s \"%s\" %s\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		methodColor, param.Method, resetColor,
		param.Path,
		param.Request.Proto,
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

// logger_file.go ends here.
