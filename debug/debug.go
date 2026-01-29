// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// debug.go --- Debugging tools.
//
// Copyright (c) 2022-2026 Paul Ward <paul@lisphacker.uk>
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
//
//

// * Package:

package debug

// * Imports:

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// * Constants:

const (
	asciiSpace          string = " "
	asciiNewLine        string = "\n"
	unicodeLineHorz     string = "│"
	unicodeLineVert     string = "─"
	unicodeLeaderTop    string = "╭─┤"
	unicodeLineStartTop string = "├"
	unicodeLeaderBot    string = "╰"

	SpacesPerIndent int = 4  // Number of spaces to use for indentation.
	LineFullLength  int = 69 // Length of debug info line.
	LineTitleLength int = 64 // Length of title segment.
)

// * Code:

// ** Interface:

// Abstract interface that defines whether an object is debuggable.
type Debugable interface {
	Debug(...any) *Debug
}

// ** Types:

// Debug information.
type Debug struct {
	title       string
	note        string
	indent      int
	lead        string
	indexed     bool
	builder     strings.Builder
	parent      **Debug
	widthTitle  int
	widthFull   int
	widthIndent int
}

// ** Methods:

// Extract parameters for arity-1 arguments.
func (obj *Debug) extractParamsArity1(params ...any) {
	switch val := params[0].(type) {
	case int:
		obj.indent = val

	case **Debug:
		obj.parent = val
		obj.indent = (*val).indent

	case string:
		obj.note = val
	}
}

// Extract parameters for arity-2 arguments.
func (obj *Debug) extractParamsArity2(params ...any) {
	obj.indexed = true

	switch val := params[1].(type) {
	case int:
		obj.title = fmt.Sprintf("%s #%d", obj.title, val)

	case uuid.UUID:
		obj.title = fmt.Sprintf("%s id:%s", obj.title, val.String())

	default:
		obj.title = fmt.Sprintf("%s id:%s", obj.title, val)
	}
}

// Extract parameters for arity-3 arguments.
func (obj *Debug) extractParamsArity3(params ...any) {
	// This is a switch because it might be extended in the future.
	//
	//nolint:gocritic
	switch val := params[2].(type) {
	case string:
		obj.note = val
	}
}

// Initialise the debug object.
func (obj *Debug) Init(params ...any) {
	// Bail if already initialised.
	if obj.builder.Len() > 0 {
		return
	}

	// Set up indentation.
	obj.indent = SpacesPerIndent
	obj.widthTitle = LineTitleLength
	obj.widthFull = LineFullLength

	if len(params) > 0 {
		obj.extractParamsArity1(params...)

		//nolint:gomnd
		if len(params) > 1 {
			obj.extractParamsArity2(params...)
		}

		//nolint:all
		if len(params) > 2 {
			obj.extractParamsArity3(params...)
		}
	}

	obj.widthIndent = obj.indent + 1
	obj.lead = strings.Repeat(asciiSpace, obj.indent)

	if obj.parent != nil && *obj.parent != nil {
		obj.widthTitle = (*obj.parent).widthTitle - obj.widthIndent
		obj.widthFull = (*obj.parent).widthFull - obj.widthIndent
	}

	pad := strings.Repeat(
		unicodeLineVert,
		imax(0, obj.widthTitle-len(obj.title)))

	if obj.indexed {
		obj.append(asciiNewLine)
	}

	obj.append(fmt.Sprintf(
		"%s %s %s%s\n",
		unicodeLeaderTop,
		obj.title,
		unicodeLineStartTop,
		pad))

	if len(obj.note) > 0 {
		obj.Printf("NOTE: %s", obj.note)
		obj.Printf("")
	}
}

// Finalise the debug object.
func (obj *Debug) End() {
	obj.append(fmt.Sprintf(
		"%s%s\n",
		unicodeLeaderBot,
		strings.Repeat(unicodeLineVert, obj.widthFull)))
}

// Append some text to the debug information.
func (obj *Debug) append(line string) {
	if obj.parent != nil && *obj.parent != nil {
		(*obj.parent).append(unicodeLineHorz + obj.lead + line)

		return
	}

	obj.builder.WriteString(line)
}

// Print to the debug information.
func (obj *Debug) Printf(format string, args ...any) {
	obj.append(fmt.Sprintf(
		"%s %s\n",
		unicodeLineHorz,
		fmt.Sprintf(format, args...)))
}

// Print out the debug information to standard output.
func (obj *Debug) Print() {
	if obj.parent == nil {
		//nolint:forbidigo
		fmt.Print(obj.String())
	}
}

// Return the string representation of the debug object.
func (obj *Debug) String() string {
	return obj.builder.String()
}

// ** Functions:

// Cheap `max` implementation.  Integer-only .
func imax(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// Create a new debug object.
func NewDebug(title string) *Debug {
	return &Debug{
		title: title,
	}
}

// Print out a `Debugable` thing to standard output.
//
//nolint:revive
func DebugPrint(thing any, params ...any) {
	impl, ok := thing.(Debugable)
	if !ok {
		return
	}

	impl.Debug(params...).Print()
}

// Return the string value of a `Debugable` thing.
//
//nolint:revive
func DebugString(thing any, params ...any) (string, bool) {
	impl, ok := thing.(Debugable)
	if !ok {
		return "", false
	}

	return impl.Debug(params...).String(), true
}

// * debug.go ends here.
