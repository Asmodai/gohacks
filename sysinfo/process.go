/*
 * process.go --- SysInfo process.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package sysinfo

import (
	"github.com/Asmodai/gohacks/di"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"

	"log"
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
	sip.si.UpdateStats()

	log.Printf(
		"SYSINFO: Running=%v Alloc=%vMiB  Heap=%vMiB  Sys=%vMiB  NumGC=%v  GoRoutines=%d",
		sip.si.RunTime(),
		sip.si.Allocated(),
		sip.si.Heap(),
		sip.si.System(),
		sip.si.GC(),
		sip.si.GoRoutines(),
	)
}

func Spawn(interval time.Duration) (*process.Process, error) {
	name := "SysInfo"

	dism := di.GetInstance()
	mgr, found := dism.Get("ProcMgr")
	if !found {
		return nil, types.NewError(
			"SYSINFO",
			"Could not locate 'ProcMgr' service.",
		)
	}

	inst, found := mgr.(*process.Manager).Find(name)
	if found {
		return inst, nil
	}

	si := NewSysInfoProc()
	conf := &process.Config{
		Name:     name,
		Interval: interval,
		Function: si.Action,
	}
	pr := mgr.(*process.Manager).Create(conf)

	go pr.Run()

	return pr, nil
}

/* process.go ends here. */
