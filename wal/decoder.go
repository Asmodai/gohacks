// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// decoder.go --- Decoder helper.
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

type decoder struct {
	data   []byte
	length int64
	offset int64
}

// ** Methods:

func (e *decoder) u32() (uint32, bool) {
	if e.offset+i32Size > e.length {
		return 0, false
	}

	val := binary.LittleEndian.Uint32(e.data[e.offset : e.offset+i32Size])
	e.offset += i32Size

	return val, true
}

func (e *decoder) u64() (uint64, bool) {
	if e.offset+i64Size > e.length {
		return 0, false
	}

	val := binary.LittleEndian.Uint64(e.data[e.offset : e.offset+i64Size])
	e.offset += i64Size

	return val, true
}

func (e *decoder) bytes(length uint32) ([]byte, bool) {
	len64 := int64(length)

	if e.offset+len64 > e.length {
		return []byte{}, false
	}

	val := e.data[e.offset : e.offset+len64]
	e.offset += len64

	return val, true
}

// ** Functions:

func newDecoder(data []byte) decoder {
	return decoder{
		data:   data,
		length: int64(len(data)),
		offset: 0,
	}
}

// * decoder.go ends here.
