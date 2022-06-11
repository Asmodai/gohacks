/*
 * sysinfo.go --- System information.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
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
	si := &SysInfo{
		rt:    runtime.MemStats{},
		start: time.Now(),
	}

	si.initHostname()

	return si
}

// Initialize hostname field.
func (si *SysInfo) initHostname() {
	host, err := os.Hostname()
	if err != nil {
		log.Printf("os.Hostname(): %s", err.Error())
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
	return time.Now().Sub(si.start)
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

/* sysinfo.go ends here. */
