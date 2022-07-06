/*
 * colorstring.go --- Coloured text.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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

package types

import (
	"fmt"
)

const (
	// CSI Pm m -- Character Attributes (SGR).
	NORMAL        = 0 // Normal (default), VT100.
	BOLD          = 1 // Bold, VT100.
	FAINT         = 2 // Faint, decreased intensity, ECMA-48 2nd.
	ITALICS       = 3 // Italicized, ECMA-48 2nd.
	UNDERLINE     = 4 // Underlined, VT100.
	BLINK         = 5 // Blink, VT100.
	INVERSE       = 7 // Inverse, VT100.
	STRIKETHROUGH = 9 // Crossed-out characters, ECMA-48 3rd.

	BLACK   = 0
	RED     = 1
	GREEN   = 2
	YELLOW  = 3
	BLUE    = 4
	MAGENTA = 5
	CYAN    = 6
	WHITE   = 7
	DEFAULT = 9
)

type ColorString struct {
	data string
	attr int
	fg   int
	bg   int
}

func MakeColorString() *ColorString {
	return &ColorString{
		attr: 0,
		fg:   9,
		bg:   9,
	}
}

func (cs *ColorString) SetString(val string) {
	cs.data = val
}

func (cs *ColorString) SetAttr(attr int) {
	cs.attr = attr
}

func (cs *ColorString) SetFG(col int) {
	cs.fg = col
}

func (cs *ColorString) SetBG(col int) {
	cs.bg = col
}

func (cs *ColorString) String() string {
	return fmt.Sprintf("\x1b[%d;%d;%dm%s\x1b[0m", cs.attr, 30+cs.fg, 40+cs.bg, cs.data)
}

/* colorstring.go ends here. */
