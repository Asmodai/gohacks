/*
 * process_test.go --- SysInfo process tests.
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
	"github.com/Asmodai/gohacks/process"

	"testing"
	"time"
)

var (
	testSIProc *process.Process
)

func TestProcess(t *testing.T) {
	t.Log("Does the system info process run as expected?")

	mgr := process.NewManager()

	testSIProc, err := Spawn(mgr, 1)
	if err != nil {
		t.Errorf("Spawn: %s", err.Error())
		return
	}

	time.Sleep(2 * time.Second)

	if !testSIProc.Running {
		t.Error("Process is not running.")
		testSIProc.Stop()
		return
	}

	testSIProc.Stop()
	t.Log("Yes.")
}

/* process_test.go ends here. */
