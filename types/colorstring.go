/*
 * colorstring.go --- Coloured text.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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
