/*
 * process.go --- SysInfo process.
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
	"github.com/Asmodai/gohacks/process"

	"time"
)

type Proc struct {
	si *SysInfo
}

func NewProc() *Proc {
	return &Proc{
		si: NewSysInfo(),
	}
}

func (sip *Proc) Action(state **process.State) {
	ps := *state

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

func Spawn(mgr process.Manager, interval int) (*process.Process, error) {
	name := "SysInfo"

	inst, found := mgr.Find(name)
	if found {
		return inst, nil
	}

	sip := NewProc()
	conf := &process.Config{
		Name:     name,
		Interval: interval,
		Function: sip.Action,
	}
	pr := mgr.Create(conf)

	go pr.Run()

	return pr, nil
}

/* process.go ends here. */
