// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// pool.go --- Pooled buffers
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

package wal

// * Imports:

import "sync"

// * Code:

// ** Type:

type bufPool struct {
	pool     sync.Pool
	capacity int
}

// ** Methods:

func (p *bufPool) get(num int) (*[]byte, []byte) {
	if num > p.capacity {
		data := make([]byte, num)

		return nil, data
	}

	pbuf, ok := p.pool.Get().(*[]byte)
	if !ok {
		panic("pool of wrong type")
	}

	return pbuf, (*pbuf)[:num]
}

func (p *bufPool) put(pbuf *[]byte) {
	if pbuf == nil {
		return
	}

	*pbuf = (*pbuf)[:p.capacity]

	p.pool.Put(pbuf)
}

// ** Functions:

func newBufPool(capacity int) *bufPool {
	return &bufPool{
		capacity: capacity,
		pool: sync.Pool{
			New: func() any {
				data := make([]byte, capacity)

				return &data
			},
		},
	}
}

// * pool.go ends here.
