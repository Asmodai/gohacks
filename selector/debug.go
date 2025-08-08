// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// debug.go --- Debugging.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
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

package selector

// * Imports:

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// * Variables:

//nolint:gochecknoglobals
var (
	traceMu   sync.RWMutex
	enabled   = false
	traceOut  = os.Stderr
	traceTime = true
)

// * Code:

func EnableTrace() {
	traceMu.Lock()
	defer traceMu.Unlock()

	enabled = true
}

func DisableTrace() {
	traceMu.Lock()
	defer traceMu.Unlock()

	enabled = false
}

func SetTraceOutput(out *os.File) {
	traceMu.Lock()
	defer traceMu.Unlock()

	traceOut = out
}

func shortFile(name string) string {
	dir, file := filepath.Split(name)
	dir = strings.TrimRight(dir, string(filepath.Separator))
	base := filepath.Base(dir)

	if len(base) == 0 || base == "." {
		return file
	}

	return base + string(filepath.Separator) + file
}

//nolint:unused
func shortFunc(fun string) string {
	if idx := strings.LastIndex(fun, "/"); idx >= 0 {
		fun = fun[idx+1:]
	}

	// drop package prefix
	if idx := strings.Index(fun, "."); idx >= 0 {
		fun = fun[idx+1:]
	}

	return fun
}

func Trace(format string, args ...any) {
	TraceWithWrapper(false, format, args...)
}

func TraceWithWrapper(isWrapped bool, format string, args ...any) {
	traceMu.RLock()
	en := enabled
	out := traceOut
	showTime := traceTime
	traceMu.RUnlock()

	if !en {
		return
	}

	caller := 1
	if isWrapped {
		caller = 2 // 2 == caller's caller.
	}

	// 0 = Trace, 1 = its caller. Bump if you add wrappers.
	_, file, line, good := runtime.Caller(caller)

	tstamp := ""
	if showTime {
		tstamp = time.Now().Format(time.RFC3339) + " "
	}

	msg := fmt.Sprintf(format, args...)
	if good {
		fmt.Fprintf(out, "%s%s:%d\t[TRACE] %s\n",
			tstamp,
			shortFile(file),
			line,
			msg)
	} else {
		fmt.Fprintf(out, "%s[TRACE] %s\n",
			tstamp,
			msg)
	}
}

// * debug.go ends here.
