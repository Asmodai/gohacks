/*
 * process.go --- SysInfo process.
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

	"time"
)

type SysInfoProc struct {
	si *SysInfo
}

func NewSysInfoProc() *SysInfoProc {
	return &SysInfoProc{
		si: NewSysInfo(),
	}
}

func (sip *SysInfoProc) Action(state **process.State) {
	var ps *process.State = *state

	sip.si.UpdateStats()

	ps.Logger().Info(
		"System Information.",
		"runtime", sip.si.RunTime().Round(time.Second),
		"allocated_mib", sip.si.Allocated(),
		"heap_mib", sip.si.Heap(),
		"system_mib", sip.si.System(),
		"collections", sip.si.GC(),
		"goroutines", sip.si.GoRoutines(),
	)
}

func Spawn(mgr process.IManager, interval int) (*process.Process, error) {
	name := "SysInfo"

	inst, found := mgr.Find(name)
	if found {
		return inst, nil
	}

	si := NewSysInfoProc()
	conf := &process.Config{
		Name:     name,
		Interval: interval,
		Function: si.Action,
	}
	pr := mgr.Create(conf)

	go pr.Run()

	return pr, nil
}

/* process.go ends here. */
