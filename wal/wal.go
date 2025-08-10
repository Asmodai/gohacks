// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// wal.go --- Write Ahead Log facility.
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
//go:build amd64 || arm64 || riscv64

// * Comments:

// * Package:

package wal

// * Imports:

import (
	"context"
	"encoding/binary"
	"hash/crc32"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"github.com/Asmodai/gohacks/logger"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	magicNumber   = 0x584C4157 // WAL magic.  'W' 'A' 'L' 'X'
	versionNumber = 1

	defaultFileMode = 0o644

	i32Size    = 4 // Size of a 32-bit integer.
	i64Size    = 8 // Size of a 64-bit integer.
	headerSize = 8 // Size of the WAL header.

	// Size of the WAL payload prefix (not including the CRC).
	payloadSizePrefix = i64Size + i64Size + i32Size + i32Size

	// Default max number of bytes in a record.
	defaultMaxValueBytes = uint32(2048)

	// Default max number of bytes in a key.
	defaultMaxKeyBytes = uint32(256)

	defaultInitialCapacity = 1024
)

// * Variables:

var (
	ErrKeyTooLarge       = errors.Base("key too large")
	ErrValueTooLarge     = errors.Base("value too large")
	ErrRecordTooLarge    = errors.Base("record too large")
	ErrShortWrite        = errors.Base("short write")
	ErrInvalidHeader     = errors.Base("invalid WAL header")
	ErrInvalidLog        = errors.Base("WAL log file is invalid")
	ErrTimestampNegative = errors.Base("negative value timestamp")
	ErrTimestampTooBig   = errors.Base("timestamp too big")
)

// * Code:

// ** Pooled buffer:

type bufPool struct {
	capacity int
	pool     sync.Pool
}

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

// ** Encoding and decoding helpers:

// *** Encoding:

type encoder struct {
	data   []byte
	offset int
}

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

func newEncoder(data []byte) encoder {
	return encoder{
		data:   data,
		offset: 0}
}

// *** Decoding:

type decoder struct {
	data   []byte
	length int64
	offset int64
}

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

func newDecoder(data []byte) decoder {
	return decoder{
		data:   data,
		length: int64(len(data)),
		offset: 0,
	}
}

// ** Record size structure:

type recSize struct {
	sizeU32 uint32
	sizeI64 int64
	nextPos int64
}

// ** Record fields structure:

type recFields struct {
	lsn    uint64
	tstamp int64
	klen   uint32
	vlen   uint32
	key    []byte
	val    []byte
}

// ** Types:

type writeAheadLog struct {
	lgr         logger.Logger
	ctx         context.Context
	cancel      context.CancelFunc
	path        string
	fptr        *os.File
	bytes       int64
	lastSyncAt  int64
	mu          sync.Mutex
	policy      Policy
	flushTicker *time.Ticker
	stopCh      chan struct{}
	dirty       bool
	crcTab      *crc32.Table
	pool        *bufPool
}

// ** Methods:

func (wal *writeAheadLog) validateKV(tstamp int64, key, val []byte) (uint32, uint32, uint64, error) {
	if err := wal.ctx.Err(); err != nil {
		return 0, 0, 0, errors.WithStack(err)
	}

	klen, klenOk := lenU32(key)
	if !klenOk || klen > wal.policy.MaxKeyBytes {
		return 0, 0, 0, errors.WithMessagef(
			ErrKeyTooLarge,
			"%d bytes",
			len(key))
	}

	vlen, vlenOk := lenU32(val)
	if !vlenOk || vlen > wal.policy.MaxValueBytes {
		return 0, 0, 0, errors.WithMessagef(
			ErrValueTooLarge,
			"%d bytes",
			len(val))
	}

	tsu, err := tstampU64(tstamp)
	if err != nil {
		return 0, 0, 0, errors.WithStack(err)
	}

	return klen, vlen, tsu, nil
}

func (wal *writeAheadLog) allocRecordBuf(klen, vlen uint32) (*[]byte, []byte, uint32, error) {
	size := recordSize(klen, vlen)

	bufLen, ok := fitsAlloc(size)
	if !ok {
		return nil, nil, 0, errors.WithMessagef(
			ErrRecordTooLarge,
			"%d bytes",
			size)
	}

	if wal.pool == nil {
		return nil, make([]byte, bufLen), bufLen, nil
	}

	pool, buf := wal.pool.get(int(bufLen))

	return pool, buf, bufLen, nil
}

func (wal *writeAheadLog) Append(lsn uint64, tstamp int64, key, val []byte) error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	klen, vlen, tsu, err := wal.validateKV(tstamp, key, val)
	if err != nil {
		return errors.WithStack(err)
	}

	pbuf, buf, buflenU32, err := wal.allocRecordBuf(klen, vlen)
	if err != nil {
		return errors.WithStack(err)
	}

	buflen := int(buflenU32)

	encodeRecord(buf, buflenU32, lsn, tsu, klen, vlen, key, val)
	finaliseCRC(buf, buflen, wal.crcTab)

	if err := writeFullAt(wal.fptr, buf, wal.bytes); err != nil {
		if wal.pool != nil {
			wal.pool.put(pbuf)
		}

		return errors.WithStack(err)
	}

	if wal.pool != nil {
		wal.pool.put(pbuf)
	}

	wal.bytes += int64(buflen)
	wal.dirty = true

	if wal.policy.SyncEveryBytes > 0 && (wal.bytes-wal.lastSyncAt) >= wal.policy.SyncEveryBytes {
		if err := wal.syncLocked(); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (wal *writeAheadLog) readRecSizeAt(pos, end int64) (recSize, bool, error) {
	var out recSize

	if pos+int64(i32Size) > end {
		return out, true, nil
	}

	var sizeb [i32Size]byte

	if _, err := wal.fptr.ReadAt(sizeb[:], pos); err != nil {
		if errors.Is(err, io.EOF) {
			return out, true, nil
		}

		return out, false, errors.WithStack(err)
	}

	out.sizeU32 = binary.LittleEndian.Uint32(sizeb[:])
	out.sizeI64 = int64(out.sizeU32)
	out.nextPos = pos + i32Size

	if out.sizeI64 < payloadSizePrefix+i32Size {
		return out, true, nil
	}

	return out, false, nil
}

func (wal *writeAheadLog) ensureScratch(scratch []byte, need int) []byte {
	if cap(scratch) >= need {
		return scratch[:need]
	}

	capNew := cap(scratch)
	if capNew == 0 {
		capNew = defaultInitialCapacity
	}

	for capNew < need {
		if capNew > math.MaxInt/2 {
			capNew = need

			break
		}

		capNew *= 2
	}

	slab := make([]byte, capNew)

	return slab[:need]
}

func (wal *writeAheadLog) readRecAt(pos int64, size int, scratch []byte) ([]byte, error) {
	buf := wal.ensureScratch(scratch, size)

	if _, err := wal.fptr.ReadAt(buf, pos); err != nil {
		return nil, errors.WithStack(err)
	}

	return buf, nil
}

//nolint:cyclop
func (wal *writeAheadLog) decodeFields(full []byte) (recFields, bool, error) {
	payload, ok := wal.validateCRC(full)
	if !ok {
		return recFields{}, false, nil
	}

	dec := newDecoder(payload)
	lsn, ok1 := dec.u64()
	tsu, ok2 := dec.u64()

	if !ok1 || !ok2 {
		return recFields{}, false, nil
	}

	tstamp, err := u64Tstamp(tsu)
	if err != nil {
		return recFields{}, false, errors.WithStack(err)
	}

	klen, ok3 := dec.u32()
	vlen, ok4 := dec.u32()

	if !ok3 || !ok4 {
		return recFields{}, false, nil
	}

	if klen > wal.policy.MaxKeyBytes {
		return recFields{}, false, errors.WithMessagef(
			ErrKeyTooLarge,
			"key length %d",
			klen)
	}

	if vlen > wal.policy.MaxValueBytes {
		return recFields{}, false, errors.WithMessagef(
			ErrValueTooLarge,
			"value length %d",
			vlen)
	}

	if dec.offset+int64(klen+vlen) > int64(len(payload)) {
		return recFields{}, false, nil
	}

	key, ok5 := dec.bytes(klen)
	val, ok6 := dec.bytes(vlen)

	if !ok5 || !ok6 {
		return recFields{}, false, nil
	}

	return recFields{
			lsn:    lsn,
			tstamp: tstamp,
			klen:   klen,
			vlen:   vlen,
			key:    key,
			val:    val,
		},
		true,
		nil
}

func (wal *writeAheadLog) validateCRC(buf []byte) ([]byte, bool) {
	if len(buf) < i32Size {
		return nil, false
	}

	payload := buf[:len(buf)-i32Size]
	want := binary.LittleEndian.Uint32(buf[len(buf)-i32Size:])
	have := crc32.Checksum(payload, wal.crcTab)

	return payload, have == want
}

func (wal *writeAheadLog) Sync() error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	if !wal.dirty || wal.bytes == wal.lastSyncAt {
		return nil
	}

	result := wal.syncLocked()

	return result
}

func (wal *writeAheadLog) syncLocked() error {
	if err := wal.fptr.Sync(); err != nil {
		return errors.WithStack(err)
	}

	wal.lastSyncAt = wal.bytes
	wal.dirty = false

	return nil
}

func (wal *writeAheadLog) Reset() error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	// Truncate back to just the header.
	if err := wal.fptr.Truncate(headerSize); err != nil {
		return errors.WithStack(err)
	}

	wal.bytes = headerSize

	if err := wal.syncLocked(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

//nolint:cyclop,funlen
func (wal *writeAheadLog) Replay(baseLSN uint64, applyCb ApplyCallbackFn) (uint64, error) {
	wal.mu.Lock()
	end := wal.bytes
	wal.mu.Unlock()

	var (
		pos     = int64(headerSize)
		maxLSN  = baseLSN
		scratch []byte
	)

	for pos < end {
		size, done, err := wal.readRecSizeAt(pos, end)

		if done {
			break
		}

		if err != nil {
			return maxLSN, err
		}

		total := int(size.sizeU32)

		if pos+int64(i32Size)+int64(total) > end {
			break
		}

		scratch, err = wal.readRecAt(size.nextPos, total, scratch)
		if err != nil {
			wal.lgr.Info(
				"Write Ahead Log CRC failed",
				"err", err.Error())

			return maxLSN, errors.WithStack(err)
		}

		rec, recOk, err := wal.decodeFields(scratch)
		if err != nil {
			wal.lgr.Info(
				"Write Ahead Log truncated tail",
				"err", err.Error())

			return maxLSN, errors.WithStack(err)
		}

		if !recOk {
			break
		}

		if rec.lsn > baseLSN {
			err := applyCb(rec.lsn, rec.tstamp, rec.key, rec.val)
			if err != nil {
				wal.lgr.Info(
					"Write ahead Log callback failed",
					"err", err.Error())

				return maxLSN, errors.WithStack(err)
			}

			if rec.lsn > maxLSN {
				maxLSN = rec.lsn
			}
		}

		pos = size.nextPos + int64(total)
	}

	return maxLSN, nil
}

func (wal *writeAheadLog) Close() error {
	wal.stopSyncTicker()

	if err := wal.flushIfDirty(); err != nil {
		_ = wal.fptr.Close()

		return errors.WithStack(err)
	}

	return errors.WithStack(wal.fptr.Close())
}

func (wal *writeAheadLog) flushIfDirty() error {
	wal.mu.Lock()
	defer wal.mu.Unlock()

	if !wal.dirty || wal.bytes == wal.lastSyncAt {
		return nil
	}

	return errors.WithStack(wal.syncLocked())
}

func (wal *writeAheadLog) SetPolicy(pol Policy) {
	var (
		needStop  bool
		needStart bool
	)

	pol.sanity()

	wal.mu.Lock()
	// CRITICAL SECTION START.
	{
		if pol.SyncEvery <= 0 {
			needStop = wal.flushTicker != nil
		} else {
			if wal.flushTicker != nil {
				wal.flushTicker.Reset(pol.SyncEvery)
			} else {
				needStart = true
			}
		}

		wal.policy = pol
	}
	// CRITICAL SECTION END.
	wal.mu.Unlock()

	if needStop {
		wal.stopSyncTicker()
	}

	if needStart {
		wal.createSyncTicker()
	}
}

func (wal *writeAheadLog) Cancel() {
	wal.cancel()
	wal.stopSyncTicker()
}

func (wal *writeAheadLog) stopSyncTicker() {
	var (
		ticker *time.Ticker
		stop   chan struct{}
	)

	wal.mu.Lock()
	// CRITICAL SECTION START.
	{
		ticker = wal.flushTicker
		stop = wal.stopCh

		wal.flushTicker = nil
		wal.stopCh = nil
	}
	// CRITICAL SECTION END.
	wal.mu.Unlock()

	if ticker != nil {
		ticker.Stop()
	}

	if stop != nil {
		close(stop)
	}
}

func (wal *writeAheadLog) createSyncTicker() {
	var (
		ticker *time.Ticker
		stop   chan struct{}
		ctx    context.Context
	)

	wal.mu.Lock()
	// CRITICAL SECTION START.
	{
		if wal.policy.SyncEvery <= 0 || wal.flushTicker != nil {
			wal.mu.Unlock()

			return
		}

		ticker = time.NewTicker(wal.policy.SyncEvery)
		stop = make(chan struct{})

		wal.flushTicker = ticker
		wal.stopCh = stop

		ctx = wal.ctx
	}
	// CRITICAL SECTION END
	wal.mu.Unlock()

	// Capture the ticker.
	tick := ticker.C

	go func(tick <-chan time.Time, stop <-chan struct{}, ctx context.Context) {
		for {
			select {
			case <-tick:
				_ = wal.flushIfDirty()

			case <-stop:
				return
			case <-ctx.Done():
				return
			}
		}
	}(tick, stop, ctx)
}

// ** Functions:

func finaliseCRC(buf []byte, bufLen int, tab *crc32.Table) {
	payLen := bufLen - i32Size - i32Size
	crc := crc32.Checksum(buf[i32Size:i32Size+payLen], tab)

	binary.LittleEndian.PutUint32(buf[bufLen-i32Size:bufLen], crc)
}

// Fill everything except CRC.
func encodeRecord(buf []byte, bufLen uint32, lsn uint64, tsu uint64, klen, vlen uint32, key, val []byte) {
	enc := newEncoder(buf)

	enc.u32(bufLen - i32Size)
	enc.u64(lsn)
	enc.u64(tsu)
	enc.u32(klen)
	enc.u32(vlen)
	enc.copy(key)
	enc.copy(val)
}

func tstampU64(tstamp int64) (uint64, error) {
	if tstamp < 0 {
		return 0, errors.WithStack(ErrTimestampNegative)
	}

	return uint64(tstamp), nil
}

func u64Tstamp(val uint64) (int64, error) {
	if val > uint64(math.MaxInt64) {
		return 0, errors.WithStack(ErrTimestampTooBig)
	}

	return int64(val), nil
}

func lenU32(thing []byte) (uint32, bool) {
	ilen := len(thing)

	if uint64(ilen) <= math.MaxUint32 {
		return uint32(ilen), true //nolint:gosec
	}

	return 0, false
}

func recordSize(key, val uint32) int64 {
	return int64(i32Size) +
		(payloadSizePrefix + int64(key) + int64(val)) +
		int64(i32Size)
}

// This is pedantic, it will refuse to accept any value that is higher than
// the maximum 32-bit unsigned.  Even though `int` on a 64-bit platform
// is a much bigger number.
//
// This is so that we can ensure 32-bit on everything that needs 32-bit.
func fitsAlloc(num int64) (uint32, bool) {
	if num <= 0 || num > int64(math.MaxUint32) {
		return 0, false
	}

	return uint32(num), true
}

func writeFullAt(fptr *os.File, data []byte, offset int64) error {
	for len(data) > 0 {
		num, err := fptr.WriteAt(data, offset)
		if err != nil {
			return errors.WithStack(err)
		}

		if num == 0 {
			return errors.WithMessage(
				ErrShortWrite,
				"wrote 0 bytes")
		}

		offset += int64(num)
		data = data[num:]
	}

	return nil
}

func readHeader(fptr *os.File) (bool, error) {
	var hdr [headerSize]byte

	if _, err := fptr.ReadAt(hdr[:], 0); err != nil {
		return false, errors.WithStack(err)
	}

	if binary.LittleEndian.Uint32(hdr[0:i32Size]) != magicNumber {
		return false, nil
	}

	if binary.LittleEndian.Uint32(hdr[i32Size:i64Size]) != versionNumber {
		return false, nil
	}

	return true, nil
}

func writeHeader(fptr *os.File) (int64, error) {
	var hdr [headerSize]byte

	binary.LittleEndian.PutUint32(hdr[0:4], magicNumber)
	binary.LittleEndian.PutUint32(hdr[4:8], versionNumber)

	if _, err := fptr.WriteAt(hdr[:], 0); err != nil {
		fptr.Close()

		return 0, errors.WithStack(err)
	}

	if err := fptr.Sync(); err != nil {
		fptr.Close()

		return 0, errors.WithStack(err)
	}

	return headerSize, nil
}

func OpenWAL(ctx context.Context, path string, syncEveryBytes int64) (WriteAheadLog, error) {
	return OpenWALWithPolicy(
		ctx,
		path,
		Policy{SyncEveryBytes: syncEveryBytes})
}

//nolint:cyclop,funlen
func OpenWALWithPolicy(parent context.Context, path string, pol Policy) (WriteAheadLog, error) {
	var pool *bufPool

	lgr := logger.MustGetLogger(parent)

	fptr, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, defaultFileMode)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	finfo, err := fptr.Stat()
	if err != nil {
		fptr.Close()

		return nil, errors.WithStack(err)
	}

	pol.sanity()

	ctx, cancel := context.WithCancel(parent)

	maxRecSize := int(recordSize(
		pol.MaxKeyBytes,
		pol.MaxValueBytes))

	poolSize, poolOk := fitsAlloc(int64(maxRecSize))
	if poolOk {
		pool = newBufPool(int(poolSize))
	}

	wal := &writeAheadLog{
		lgr:        lgr,
		ctx:        ctx,
		cancel:     cancel,
		crcTab:     crc32.MakeTable(crc32.Castagnoli),
		path:       path,
		fptr:       fptr,
		bytes:      finfo.Size(),
		lastSyncAt: finfo.Size(),
		policy:     pol,
		pool:       pool}

	switch {
	case finfo.Size() == 0:
		written, err := writeHeader(fptr)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		wal.bytes = written
		wal.lastSyncAt = written

	case finfo.Size() >= headerSize:
		hdrOk, err := readHeader(fptr)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		if !hdrOk {
			return nil, errors.WithStack(ErrInvalidHeader)
		}

	default:
		fptr.Close()

		return nil, errors.WithStack(ErrInvalidLog)
	}

	if pol.SyncEvery > 0 {
		// Context used by this is part of the structure.
		wal.createSyncTicker()
	}

	return wal, nil
}

// * wal.go ends here.
