// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// colorstring.go --- Coloured text.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package types

import (
	"fmt"
)

// CSI Pm [; Pm ...] m -- Character Attributes (SGR).
const (
	// Normal (default) attribtues -- VT100.
	NORMAL = 0

	// Bold -- VT100.
	BOLD = 1

	// Faint, decreased intensity -- ECMA-48 2e.
	FAINT = 2

	// Italicizsed -- ECMA-48 2e.
	ITALICS = 3

	// Underlined -- VT100.
	UNDERLINE = 4

	// Blinking -- VT100.
	BLINK = 5

	// Inverse video -- VT100.
	INVERSE = 7

	// Crossed-out characters -- ECMA-48 3e.
	STRIKETHROUGH = 9

	// Black colour -- ANSI.
	BLACK = 0

	// Red colour -- ANSI.
	RED = 1

	// Green colour -- ANSI.
	GREEN = 2

	// Yellow colour -- ANSI.
	YELLOW = 3

	// Blue colour -- ANSI.
	BLUE = 4

	// Magenta colour -- ANSI.
	MAGENTA = 5

	// Cyan colour -- ANSI.
	CYAN = 6

	// White colour -- ANSI.
	WHITE = 7

	// Default colour -- ANSI.
	DEFAULT = 9

	// Offset for foreground colours.
	FGOFFSET = 30

	// Offset for background colours.
	BGOFFSET = 40
)

// Colour string.
//
// Generates a string that, with the right terminal type, display text using
// various character attributes.
//
// To find out more, consult your nearest DEC VT340 programmer's manual or the
// latest ECMA-48 standard.
type ColorString struct {
	data string
	attr int
	fg   int
	bg   int
}

// Make a new coloured string.
func MakeColorString() *ColorString {
	return &ColorString{
		attr: NORMAL,
		fg:   DEFAULT,
		bg:   DEFAULT,
	}
}

// Make a new coloured string with the given attributes.
func MakeColorStringWithAttrs(data string, attr, foreg, backg int) *ColorString {
	str := MakeColorString()

	str.SetString(data)
	str.SetAttr(attr)
	str.SetFG(foreg)
	str.SetBG(backg)

	return str
}

// Set the string to display.
func (cs *ColorString) SetString(val string) {
	cs.data = val
}

// Set the character attribute.
func (cs *ColorString) SetAttr(attr int) {
	cs.attr = attr
}

// Set the foreground colour.
func (cs *ColorString) SetFG(col int) {
	cs.fg = col
}

// Set the background colour.
func (cs *ColorString) SetBG(col int) {
	cs.bg = col
}

// Convert to a string.
//
// Warning: the resulting string will contain escape sequences for use with a
// compliant terminal or terminal emulator.
func (cs *ColorString) String() string {
	return fmt.Sprintf(
		"\x1b[%d;%d;%dm%s\x1b[0m",
		cs.attr,
		FGOFFSET+cs.fg,
		BGOFFSET+cs.bg,
		cs.data,
	)
}

// colorstring.go ends here.
