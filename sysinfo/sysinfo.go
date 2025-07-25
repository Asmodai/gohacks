// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// sysinfo.go --- System information.
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

package sysinfo

import (
	"github.com/Asmodai/gohacks/math/conversion"

	"log"
	"os"
	"runtime"
	"sync"
	"time"
)

type SysInfo struct {
	sync.Mutex

	hostname   string
	start      time.Time
	rt         runtime.MemStats
	goroutines int
}

// Create a new System Information instance.
func NewSysInfo() *SysInfo {
	siproc := &SysInfo{
		rt:    runtime.MemStats{},
		start: time.Now(),
	}

	siproc.initHostname()

	return siproc
}

// Initialize hostname field.
func (si *SysInfo) initHostname() {
	host, err := os.Hostname()
	if err != nil {
		log.Printf("os.Hostname(): %s", err.Error())

		// Nod to BSD systems right here.
		host = "amnesiac"
	}

	si.Lock()
	si.hostname = host
	si.Unlock()
}

// Update runtime statistics.
func (si *SysInfo) UpdateStats() {
	si.Lock()

	runtime.ReadMemStats(&si.rt)
	si.goroutines = runtime.NumGoroutine()

	si.Unlock()
}

// Return this system's hostname.
func (si *SysInfo) Hostname() string {
	return si.hostname
}

// Return the time running.
func (si *SysInfo) RunTime() time.Duration {
	return time.Since(si.start)
}

// Return number of MiB currently allocated.
func (si *SysInfo) Allocated() uint64 {
	return conversion.BToMiB(si.rt.Alloc)
}

// Return number of MiB used by the heap.
func (si *SysInfo) Heap() uint64 {
	return conversion.BToMiB(si.rt.HeapSys)
}

// Return number of MiB allocated from the system.
func (si *SysInfo) System() uint64 {
	return conversion.BToMiB(si.rt.Sys)
}

// Return the number of collections performed.
func (si *SysInfo) GC() uint32 {
	return si.rt.NumGC
}

// Return the number of Go routines.
func (si *SysInfo) GoRoutines() int {
	return si.goroutines
}

// sysinfo.go ends here.
