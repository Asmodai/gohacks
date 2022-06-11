/*
 * sysinfo_test.go --- Tests for sysinfo package.
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
	"testing"
)

var (
	sinfo *SysInfo
)

// Main testing function.
func TestMain(m *testing.M) {
	log.Println("Setting up.")
	sinfo = NewSysInfo()

	log.Println("Running tests.")
	val := m.Run()

	log.Println("Shutting down.")
	os.Exit(val)
}

// Test hostname resolution.
func TestHostname(t *testing.T) {
	t.Log("Does `Hostname` work as expected?")

	this, err := os.Hostname()
	if err != nil {
		t.Errorf("os.Hostname(): %s", err.Error())
		return
	}

	if this == sinfo.Hostname() {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

// Test runtime statistics.
func TestMemStats(t *testing.T) {
	var stats runtime.MemStats
	var fail bool = false

	runtime.ReadMemStats(&stats)
	sinfo.UpdateStats()
	runtime.ReadMemStats(&stats)

	t.Log("Do our runtime stats match?")
	if conversion.BToMiB(stats.Alloc) != sinfo.Allocated() {
		t.Error("`Allocated` does NOT match.")
		fail = true
	}

	if conversion.BToMiB(stats.HeapSys) != sinfo.Heap() {
		t.Errorf("`Heap` does NOT match.")
		fail = true
	}

	if conversion.BToMiB(stats.Sys) != sinfo.System() {
		t.Error("`System` does NOT match.")
		fail = true
	}

	if stats.NumGC != sinfo.GC() {
		t.Error("`GC` does NOT match.")
		fail = true
	}

	if fail {
		t.Error("One or more runtime stats are incorrect.")
	} else {
		t.Log("Yes.")
	}
}

// Test goroutine count.
func TestGoRoutine(t *testing.T) {
	t.Log("Does `GoRoutines` match?")

	if runtime.NumGoroutine() == sinfo.GoRoutines() {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

/* sysinfo_test.go ends here. */
