// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// types.go --- Write Ahead Log types.
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
//
//mock:yes

// * Comments:

// * Package:

package wal

// * Code:

type ApplyCallbackFn func(lsn uint64, tstamp int64, key, value []byte) error

// WriteAheadLog is a single-writer, crash-safe, append-only log.
//
// Records are written in this binary layout:
//
//	[size:u32][lsn:u64][ts:u64][klen:u32][vlen:u32][key][value][crc:u32]
//
// Where `size` is the number of bytes after the size field (including the
// CRC), and `lsn` is a caller-supplied logical sequence number, `ts` is a
// non-negative UNIX timestamp (seconds), and `crc` is a CRC32C of the
// payload.
//
// Concurrency & safety:
//   - Append/Sync/Reset/SetPolicy/Close are safe to call from multiple
//     goroutines, but only one append will make progress at a time.
//   - Replay may run concurrently with appends; it snapshots the current
//     end and stops at the first incomplete/corrupt record (safe truncation
//     behavior).
//
// Durability:
//   - Durability is controlled by Policy (time-based fsync via SyncEvery,
//     byte-based via SyncEveryBytes). Sync() forces an fsync immediately.
//   - Close() stops background flush, flushes if dirty, then closes the file.
//
// Limits & validation:
//   - MaxKeyBytes and MaxValueBytes (from Policy) bound record sizes.
//   - Timestamps must be in [0, math.MaxInt64]. Violations return errors.
//
// LSN semantics:
//   - Caller is responsible for monotonically increasing LSNs.
//   - Replay(baseLSN, ...) applies records with LSN > baseLSN and returns the
//     highest LSN applied.
//
// Example:
//
// ```go
//
//	w, _ := wal.OpenWAL(ctx, "data.wal", 4<<20)
//	defer w.Close()
//	_ = w.Append(next, time.Now().Unix(), []byte("k"), []byte("v"))
//	_ = w.Sync()
//	_, _ = w.Replay(0, func(lsn uint64, ts int64, k, v []byte) error {
//		return nil
//	})
//
// ```
//
// This is probably not the best implementation of a WAL.
type WriteAheadLog interface {
	// Append writes one record to the log.
	//
	// Errors:
	//   • ErrKeyTooLarge / ErrValueTooLarge / ErrRecordTooLarge when
	//     limits are exceeded.
	//   • ErrTimestampNegative / ErrTimestampTooBig for invalid
	//     timestamps.
	//   • I/O errors wrapped with stack context
	//
	// Durability:
	//   • Subject to Policy; may buffer in page cache until
	//     SyncEvery/SyncEveryBytes hit.
	//   • Call Sync() to force durability.
	Append(lsn uint64, tstamp int64, key, value []byte) error

	// Replay re-applies records with LSN > baseLSN in log order.
	//
	// The supplied callback is invoked for each record; if the callback
	// returns an error, replay aborts and that error is returned. On
	// success, the returned uint64 is the highest LSN applied.
	//
	// Robustness:
	//   • Stops at first malformed/incomplete/CRC-mismatched record
	//     (treats as a safely truncatable tail). Earlier valid records
	//     are preserved.
	Replay(baseLSN uint64, applyCb ApplyCallbackFn) (uint64, error)

	// SetPolicy atomically updates durability and size limits at runtime.
	//
	// Notes:
	//   • Time-based flushing may start/stop a background ticker.
	//   • Limits (MaxKeyBytes/MaxValueBytes) apply to subsequent appends
	//     only.
	SetPolicy(Policy)

	// Sync forces an fsync if there are dirty bytes; otherwise it is a
	// no-op.
	//
	// Returns any I/O error encountered.
	Sync() error

	// Reset truncates the log back to just the header (discarding all
	// records), then fsyncs.
	//
	// Use with care.
	Reset() error

	// Close stops background flushing, flushes if dirty, and closes the
	// file.
	//
	// Returns any flush/close error encountered. Close is
	// idempotent-safe to call once; do not use the instance after Close.
	Close() error
}

// * types.go ends here.
