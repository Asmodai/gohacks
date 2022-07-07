/*
 * sysinfo_test.go --- Tests for sysinfo package.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
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

	sinfo.UpdateStats()
	rtgo := runtime.NumGoroutine()
	sigo := sinfo.GoRoutines()

	if rtgo == sigo {
		t.Log("Yes.")
		return
	}

	t.Errorf("No, runtime:%v, sysinfo:%v", rtgo, sigo)
}

/* sysinfo_test.go ends here. */
