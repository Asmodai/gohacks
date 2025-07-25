// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// colorstring.go --- Coloured text.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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
//

// * Package:

package types

// * Imports:

import (
	"strings"
)

// * Constants:

// CSI Pm [; Pm ...] m -- Character Attributes (SGR).
const (
	Black   Colour = '0' // Black colour -- ANSI.
	Red     Colour = '1' // Red colour -- ANSI.
	Green   Colour = '2' // Green colour -- ANSI.
	Yellow  Colour = '3' // Yellow colour -- ANSI.
	Blue    Colour = '4' // Blue colour -- ANSI.
	Magenta Colour = '5' // Magenta colour -- ANSI.
	Cyan    Colour = '6' // Cyan colour -- ANSI.
	White   Colour = '7' // White colour -- ANSI.
	Default Colour = '9' // Default colour -- ANSI.

	Normal     int = 0 // Reset all attributes.
	Bold       int = 1 // Bold -- VT100.
	Faint      int = 2 // Faint, decreased intensity -- ECMA-48 2e.
	Italic     int = 3 // Italicizsed -- ECMA-48 2e.
	Underline  int = 4 // Underlined -- VT100.
	Blink      int = 5 // Blinking -- VT100.
	Inverse    int = 6 // Inverse video -- VT100.
	Strikethru int = 7 // Crossed-out characters -- ECMA-48 3e.

	fgPrefix   rune   = '3'       // Prefix for foreground colours.
	bgPrefix   rune   = '4'       // Offset for background colours.
	csiSep     rune   = ';'       // CSI sequence param separator.
	csiCommand rune   = 'm'       // CSI command for attributes.
	csiPrefix  string = "\x1b["   // CSI prefix.
	attrReset  string = "\x1b[0m" // CSI sequence to reset attributes.

	saneLengthGuess int = 24
)

// * Variables:

//nolint:gochecknoglobals
var (
	attrs = [8]rune{'0', '1', '2', '3', '4', '5', '7', '9'}
)

// * Code:

// ** Types:

// ECMA-48 colour descriptor type.
type Colour rune

// ECMA-48 attributes array type.
type attributes [8]bool

// Colour string.
//
// Generates a string that, with the right terminal type, display text using
// various character attributes.
//
// To find out more, consult your nearest DEC VT340 programmer's manual or the
// latest ECMA-48 standard.
type ColorString struct {
	data   string     // User-supplied data.
	attrs  attributes // The 10 main ECMA-48 attributes.
	fg     Colour     // Foreground colour.
	bg     Colour     // Background colour.
	dirty  bool       // Has the object been changed?
	cached string     // Cached string for speed.
}

// ** Methods:

// Set the string to display.
func (cs *ColorString) SetString(val string) *ColorString {
	cs.data = val
	cs.dirty = true

	return cs
}

// Helper function to set attributes.
func (cs *ColorString) setAttr(index int) *ColorString {
	cs.dirty = true

	if index == Normal {
		cs.attrs = attributes{}
		cs.attrs[0] = true
	} else {
		cs.attrs[0] = false
		cs.attrs[index] = true
	}

	return cs
}

// Reset all attributes.
func (cs *ColorString) Clear() *ColorString {
	cs.fg = Default
	cs.bg = Default

	return cs.SetNormal()
}

// Set the attribute to normal.
//
// This removes all other attributes.
func (cs *ColorString) SetNormal() *ColorString {
	return cs.setAttr(Normal)
}

// Set the bold attribute.
func (cs *ColorString) SetBold() *ColorString {
	return cs.setAttr(Bold)
}

// Set the faint attribute.
func (cs *ColorString) SetFaint() *ColorString {
	return cs.setAttr(Faint)
}

// Set the italic attribute.
func (cs *ColorString) SetItalic() *ColorString {
	return cs.setAttr(Italic)
}

// Set the underline attribute.
func (cs *ColorString) SetUnderline() *ColorString {
	return cs.setAttr(Underline)
}

// Set the blink attribute.
func (cs *ColorString) SetBlink() *ColorString {
	return cs.setAttr(Blink)
}

// Set the inverse video attribute.
func (cs *ColorString) SetInverse() *ColorString {
	return cs.setAttr(Inverse)
}

// Set the strikethrough attribute.
func (cs *ColorString) SetStrikethru() *ColorString {
	return cs.setAttr(Strikethru)
}

// Add a specific attribute.
//
// This ignores the "Normal" attribute.  To clear attributes, use
// `SetNormal` or `Clear`.
func (cs *ColorString) AddAttr(index int) *ColorString {
	if index < 1 || index >= len(cs.attrs) {
		return cs
	}

	cs.attrs[0] = false
	cs.attrs[index] = true
	cs.dirty = true

	return cs
}

// Remove a specific attribute.
//
// This ignores the "Normal" attribute.  To clear attributes, use
// `SetNormal` or `Clear`.
func (cs *ColorString) RemoveAttr(index int) *ColorString {
	if index < 1 || index >= len(cs.attrs) {
		return cs
	}

	cs.attrs[index] = false
	cs.dirty = true

	// If we've removed all useful attributes, then reset "normal".
	anySet := false

	for idx, set := range cs.attrs {
		if idx != Normal && set {
			anySet = true

			break
		}
	}

	if !anySet {
		cs.attrs[Normal] = true
	}

	return cs
}

// Set the foreground colour.
func (cs *ColorString) SetFG(col Colour) *ColorString {
	cs.fg = col
	cs.dirty = true

	return cs
}

// Set the background colour.
func (cs *ColorString) SetBG(col Colour) *ColorString {
	cs.bg = col
	cs.dirty = true

	return cs
}

// Convert to a string.
//
// Warning: the resulting string will contain escape sequences for use with a
// compliant terminal or terminal emulator.
func (cs *ColorString) String() string {
	if cs.dirty {
		cs.cached = ""
		cs.dirty = false
	}

	if len(cs.cached) > 0 {
		return cs.cached
	}

	var bldr strings.Builder

	bldr.Grow(saneLengthGuess + len(cs.data))

	bldr.WriteString(csiPrefix)

	if cs.attrs[0] {
		bldr.WriteRune(attrs[0])
		bldr.WriteRune(csiSep)
	} else {
		// Must not emit a `0` when other attrs are set.
		for idx := 1; idx < len(attrs); idx++ {
			if cs.attrs[idx] {
				bldr.WriteRune(attrs[idx])
				bldr.WriteRune(csiSep)
			}
		}
	}

	bldr.WriteRune(fgPrefix)
	bldr.WriteRune(rune(cs.fg))
	bldr.WriteRune(csiSep)

	bldr.WriteRune(bgPrefix)
	bldr.WriteRune(rune(cs.bg))

	bldr.WriteRune(csiCommand)
	bldr.WriteString(cs.data)

	bldr.WriteString(attrReset)

	cs.cached = bldr.String()

	return cs.cached
}

// ** Functions:

// Make a new coloured string.
func NewColorString() *ColorString {
	return NewColorStringWithColors("", Default, Default)
}

// Make a new coloured string with the given attributes.
func NewColorStringWithColors(
	data string,
	foreg, backg Colour,
) *ColorString {
	return &ColorString{
		attrs: attributes{
			true, // Default
			false,
			false,
			false,
			false,
			false,
			false,
			false,
		},
		fg:    foreg,
		bg:    backg,
		dirty: true,
		data:  data,
	}
}

// * colorstring.go ends here.
