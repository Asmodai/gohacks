// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// header.go --- WAL header.
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

import (
	"os"
	"time"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	// Size of the write-ahead log header.
	HeaderSize = 4 + 4 + 8 + 8

	// Write-ahead log magic number.
	//
	// This is the string `WALX` expressed as a little-endian integer.
	MagicNumber = 0x584C4157

	// Write-ahead log version number.
	VersionNumber = 1

	// WAL makes use of 32-bit CRC Castagnoli.
	FeatureCRC32C uint64 = 1 << 0

	// Version number of the latest write-ahead log facility.
	currentVersion = 1

	// Current feature flags.
	currentFeatures = FeatureCRC32C
)

// * Variables:

// * Code:

type Header struct {
	// Write-Ahead Log magic number.
	//
	// This equates to the string `WALX`.
	Magic uint32

	// File format version number.
	Version uint32

	// Features.
	Features uint64

	// When the write-ahead log was created.
	CreatedAt time.Time
}

// ** Methods:

func (h *Header) IsLatestVersion() bool {
	return h.Version == currentVersion
}

func (h *Header) HasCRC32C() bool {
	return (h.Features&FeatureCRC32C != 0)
}

func (h *Header) sanity() bool {
	if h.Magic != MagicNumber {
		return false
	}

	if h.Version == 0 && h.Version > VersionNumber {
		return false
	}

	return true
}

// ** Functions:

func ReadHeader(fptr *os.File) (Header, error) {
	// This is just a wrapper.
	//
	// It's a wrapper because the public version exposes less.
	_, header, err := readHeader(fptr)

	return header, errors.WithStack(err)
}

func readHeader(fptr *os.File) (bool, Header, error) {
	var (
		hdr    = make([]byte, HeaderSize)
		header = Header{}
		good   bool
		tstamp uint64
	)

	if _, err := fptr.ReadAt(hdr[:HeaderSize], 0); err != nil {
		return false, Header{}, errors.WithStack(err)
	}

	dec := newDecoder(hdr)

	if header.Magic, good = dec.u32(); !good {
		return false, Header{}, nil
	}

	if header.Version, good = dec.u32(); !good {
		return false, Header{}, nil
	}

	if header.Features, good = dec.u64(); !good {
		return false, Header{}, nil
	}

	if tstamp, good = dec.u64(); !good {
		return false, Header{}, nil
	}

	now, err := u64Tstamp(tstamp)
	if err != nil {
		return false, Header{}, errors.WithStack(err)
	}

	header.CreatedAt = time.Unix(now, 0)

	// Check sanity.
	sane := header.sanity()

	return sane, header, nil
}

func writeHeader(fptr *os.File) (int64, error) {
	var hdr = make([]byte, HeaderSize)

	now, err := tstampU64(time.Now().Unix())
	if err != nil {
		return 0, errors.WithStack(err)
	}

	enc := newEncoder(hdr)

	enc.u32(MagicNumber)
	enc.u32(VersionNumber)
	enc.u64(currentFeatures)
	enc.u64(now)

	if _, err := fptr.WriteAt(hdr[:HeaderSize], 0); err != nil {
		fptr.Close()

		return 0, errors.WithStack(err)
	}

	if err := fptr.Sync(); err != nil {
		fptr.Close()

		return 0, errors.WithStack(err)
	}

	return HeaderSize, nil
}

// * header.go ends here.
