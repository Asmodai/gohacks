/*
 * manager.go --- Management structure.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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

package exec

import (
	"github.com/Asmodai/gohacks/logger"

	"context"
	"fmt"
	goexec "os/exec"
	"time"
)

var (
	CheckDelaySleep time.Duration = 250 * time.Millisecond
)

type Manager struct {
	path   string
	args   Args
	logger logger.ILogger
	procs  []*goexec.Cmd
	number int
	base   int
	ctx    context.Context
}

func NewManager(lgr logger.ILogger, count, base int) *Manager {
	return &Manager{
		path:   "",
		args:   Args{base: base},
		logger: lgr,
		procs:  make([]*goexec.Cmd, count),
		number: count,
		base:   base,
		ctx:    context.TODO(),
	}
}

func (m *Manager) SetPath(val string)             { m.path = val }
func (m *Manager) SetArgs(val Args)               { m.args = val }
func (m *Manager) SetCount(val int)               { m.number = val }
func (m *Manager) SetContext(val context.Context) { m.ctx = val }

func (m *Manager) Dump() {
	for idx := range m.procs {
		fmt.Printf("Process %+#v\n", m.procs[idx])
	}
}

func (m *Manager) KillAll() {
	for idx := range m.procs {
		if m.procs[idx] == nil {
			continue
		}

		// No state and no process info, not running, so skip.
		if m.procs[idx].ProcessState == nil || m.procs[idx].Process == nil {
			continue
		}

		// Process has exited, so skip.
		if m.procs[idx].ProcessState.Exited() {
			continue
		}

		m.logger.Info(
			"Killing process.",
			"index", idx,
			"pid", m.procs[idx].Process.Pid,
			"path", m.procs[idx].Path,
		)

		if err := m.procs[idx].Process.Kill(); err != nil {
			m.logger.Warn(
				"Could not kill process.",
				"index", idx,
				"pid", m.procs[idx].Process.Pid,
				"path", m.procs[idx].Path,
				"err", err.Error(),
			)
		}
	}
}

func (m *Manager) Check() {
	for idx := range m.procs {
		// Check if our context has been cancelled.
		select {
		case <-m.ctx.Done():
			return
		default:
		}

		// We have no process defined, so, continue.
		if m.procs[idx] == nil {
			time.Sleep(CheckDelaySleep)
			continue
		}

		// No process running, spawn one.
		if m.procs[idx].Process == nil {
			m.logger.Info(
				"Starting process.",
				"index", idx,
				"path", m.procs[idx].Path,
			)

			var err error
			go func(e *error) {
				*e = m.procs[idx].Run()
			}(&err)

			// Delay slightly before checking for error, so the go
			// routine has time to do something.
			time.Sleep(CheckDelaySleep)

			if err != nil {
				m.logger.Fatal(
					"Error while starting process.",
					"index", idx,
					"path", m.procs[idx].Path,
					"args", m.procs[idx].Args,
					"err", err.Error(),
				)
			}

			continue
		}

		// Process is running but has no state, it is probably still
		// starting up, so skip.
		if m.procs[idx].ProcessState == nil {
			continue
		}

		exittype := "exited"
		switch m.procs[idx].ProcessState.Exited() {
		case false: // process has not exited.
			if m.procs[idx].ProcessState.ExitCode() != -1 {
				// and has no exit code, so break out of this.
				break
			}

			// If we reach here, then we have apparently not exited but have
			// an exit code of -1... on Unix, this indicates that the
			// process has been killed by a signal.
			exittype = "signal"
			fallthrough

		case true: // process has exited.
			m.logger.Info(
				"Process exited.",
				"index", idx,
				"type", exittype,
				"code", m.procs[idx].ProcessState.ExitCode(),
			)

			m.doSpawn(idx)
		}

		time.Sleep(CheckDelaySleep)
	}
}

func (m *Manager) doSpawn(idx int) {
	m.procs[idx] = nil
	m.procs[idx] = goexec.CommandContext(
		m.ctx,
		m.path,
		m.args.Get(idx)...,
	)
}

func (m *Manager) Spawn() {
	m.logger.Info(
		"Spawning processes.",
		"path", m.path,
		"args", m.args,
	)

	for idx := range m.procs {
		m.doSpawn(idx)
	}
}

/* manager.go ends here. */
