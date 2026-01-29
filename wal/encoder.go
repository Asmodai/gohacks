// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// encoder.go --- Encoder helpers.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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

package wal

// * Imports:

import "encoding/binary"

// * Code:

// ** Type:

type encoder struct {
	data   []byte
	offset int
}

// ** Methods:

func (e *encoder) u32(val uint32) {
	binary.LittleEndian.PutUint32(e.data[e.offset:e.offset+i32Size], val)
	e.offset += i32Size
}

func (e *encoder) u64(val uint64) {
	binary.LittleEndian.PutUint64(e.data[e.offset:e.offset+i64Size], val)
	e.offset += i64Size
}

func (e *encoder) copy(data []byte) {
	length := len(data)

	copy(e.data[e.offset:e.offset+length], data)
	e.offset += length
}

// ** Function:

func newEncoder(data []byte) encoder {
	return encoder{
		data:   data,
		offset: 0}
}

// * encoder.go ends here.
