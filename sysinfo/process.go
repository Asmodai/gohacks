// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// process.go --- SysInfo process.
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

// * Comments:

// * Package:

package sysinfo

// * Imports:

import (
	"context"
	"time"

	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"
)

// * Constants:

const (
	// Name for the process.
	processName string = "SysInfo"
)

// * Code:

// ** Types:

type Proc struct {
	si *SysInfo
}

// ** Methods:

// Function that runs every tick in the sysinfo process.
//
// Simply prints out the Go runtime stats via the process's logger.
func (sip *Proc) Action(state *process.State) {
	sip.si.UpdateStats()

	state.Logger().Debug(
		"System Information.",
		"runtime", sip.si.RunTime().Round(time.Second),
		"allocated_mib", sip.si.Allocated(),
		"heap_mib", sip.si.Heap(),
		"system_mib", sip.si.System(),
		"collections", sip.si.GC(),
		"goroutines", sip.si.GoRoutines(),
	)
}

// ** Functions:

// Create a new system information process with default values.
func NewProc() *Proc {
	return &Proc{
		si: NewSysInfo(),
	}
}

// Spawn a system information process.
//
// The provided context must have a `process.Manager` entry in its user
// value.  See `contextdi` and `process.SetProcessManager`.
func Spawn(ctx context.Context, interval types.Duration) (*process.Process, error) {
	mgr := process.MustGetProcessManager(ctx)

	inst, found := mgr.Find(processName)
	if found {
		return inst, nil
	}

	sip := NewProc()
	conf := &process.Config{
		Name:     processName,
		Interval: interval,
		Function: sip.Action,
	}
	pr := mgr.Create(conf)

	go pr.Run()

	return pr, nil
}

// * process.go ends here.
