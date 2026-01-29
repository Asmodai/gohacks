// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// tunables.go --- WAL tunables.
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

import "time"

// * Code:

type Policy struct {
	// SyncEveryBytes: fsync after this many bytes appended since the
	// last sync.
	//
	// 0 disables byte-based syncing (only time-based or explicit
	// sync/reset/close).
	SyncEveryBytes int64

	// SyncEvery: also fsync on this cadence if there is unsynced data.
	//
	// 0 disables time-based syncing.
	SyncEvery time.Duration

	// Maximum size of a value in bytes.
	MaxValueBytes uint32

	// Maximum size of a key in bytes.
	MaxKeyBytes uint32
}

func (pol *Policy) sanity() {
	if pol.MaxKeyBytes <= 0 {
		pol.MaxKeyBytes = defaultMaxKeyBytes
	}

	if pol.MaxValueBytes <= 0 {
		pol.MaxValueBytes = defaultMaxValueBytes
	}

	if pol.SyncEvery < 0 {
		pol.SyncEvery = 0
	}

	if pol.SyncEveryBytes < 0 {
		pol.SyncEveryBytes = 0
	}
}

// * tunables.go ends here.
