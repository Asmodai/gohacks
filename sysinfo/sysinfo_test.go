/*
 * sysinfo_test.go --- Tests for sysinfo package.
 *
 * Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
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

	"os"
	"runtime"
	"testing"
)

var (
	sinfo *SysInfo
)

func TestSysInfo(t *testing.T) {
	var (
		stats runtime.MemStats
		sinfo = NewSysInfo()
	)

	t.Run("`Hostname`", func(t *testing.T) {
		this, err := os.Hostname()
		if err != nil {
			t.Errorf("%s", err.Error())
		}

		if this != sinfo.Hostname() {
			t.Errorf("Hostname: %s != %s", this, sinfo.Hostname())
		}
	})

	t.Run("`MemStats`", func(t *testing.T) {
		runtime.ReadMemStats(&stats)
		sinfo.UpdateStats()
		runtime.ReadMemStats(&stats)

		t.Run("Stats match", func(t *testing.T) {
			if conversion.BToMiB(stats.Alloc) != sinfo.Allocated() {
				t.Error("`Allocated` does NOT match.")
			}
		})

		if conversion.BToMiB(stats.HeapSys) != sinfo.Heap() {
			t.Errorf("`Heap` does NOT match.")
		}

		if conversion.BToMiB(stats.Sys) != sinfo.System() {
			t.Error("`System` does NOT match.")
		}
	})

	t.Run("GC matches", func(t *testing.T) {
		if stats.NumGC != sinfo.GC() {
			t.Error("`GC` does NOT match.")
		}
	})

	t.Run("GoRoutines match", func(t *testing.T) {
		sinfo.UpdateStats()
		rtgo := runtime.NumGoroutine()
		sigo := sinfo.GoRoutines()

		if rtgo != sigo {
			t.Errorf("No, runtime:%v, sysinfo:%v", rtgo, sigo)
		}
	})
}

/* sysinfo_test.go ends here. */
