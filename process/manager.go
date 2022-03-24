/*
 * manager.go --- Process manager.
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

package process

import (
	"github.com/Asmodai/gohacks/logger"

	"context"
	"sync"
)

/*

Process manager structure.

To use,

1) Create a new process manager:

  procmgr := process.NewManager()

2) Create your process configuration:

  conf := &process.Config{
    Name:     "Windows 95",
    Interval: 10, // seconds
    Function: func(state **State) {
      // Crash or something.
    }
  }

3) Create the process itself.

  proc := procmgr.Create(conf)

4) Run the process.

  procmgr.Run("Windows 95")

/or/

  proc.Run()

*/
type Manager struct {
	processes []*Process
	logger    logger.ILogger
	parent    context.Context
	ctx       context.Context
	cancel    context.CancelFunc
	cwg       *sync.WaitGroup
}

// Create a new process manager with a given parent context.
func NewManagerWithContext(parent context.Context) *Manager {
	ctx, cancel := context.WithCancel(parent)

	return &Manager{
		processes: []*Process{},
		logger:    logger.NewDefaultLogger(),
		ctx:       ctx,
		cancel:    cancel,
		cwg:       &sync.WaitGroup{},
	}
}

// Create a new process manager.
func NewManager() *Manager {
	return NewManagerWithContext(context.TODO())
}

// Set the process manager's logger.
func (pm *Manager) SetLogger(lgr logger.ILogger) {
	pm.logger = lgr
}

// Set the process manager's context.
func (pm *Manager) SetContext(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)

	pm.parent = parent
	pm.ctx = ctx
	pm.cancel = cancel
}

// Create a new process with the given configuration.
func (pm *Manager) Create(config *Config) *Process {
	proc := NewProcess(config)
	proc.SetLogger(pm.logger)
	proc.SetWaitGroup(pm.cwg)

	pm.processes = append(pm.processes, proc)

	return proc
}

// Add an existing process to the manager.
func (pm *Manager) Add(proc *Process) {
	if proc == nil {
		return
	}

	proc.SetLogger(pm.logger)
	pm.processes = append(pm.processes, proc)
}

// Find and return the given process, or nil if not found.
func (pm *Manager) Find(name string) (*Process, bool) {
	for _, p := range pm.processes {
		if p.Name == name {
			return p, true
		}
	}

	return nil, false
}

// Run the named process.
//
// Returns 'false' if the process is not found;  otherwise returns
// the result of the process execution.
func (pm *Manager) Run(name string) bool {
	proc, found := pm.Find(name)
	if !found {
		return false
	}

	proc.SetContext(pm.ctx)

	return proc.Run()
}

// Stop the given process.
//
// Returns 'true' if the process has been stopped; otherwise 'false'.
func (pm *Manager) Stop(name string) bool {
	proc, found := pm.Find(name)
	if !found {
		return false
	}

	// Stopping one process doesn't require us to wait for the group.
	return proc.Stop()
}

// Stop all processes.
//
// Returns 'true' if *all* processes have been stopped; otherwise
// 'false' is returned.
func (pm *Manager) StopAll() bool {
	res := true

	pm.logger.Info(
		"Stopping all processes.",
		"type", "stop",
	)

	// This is better than invoking the context's cancel, as it allows
	// cleanup to be executed.
	for _, proc := range pm.processes {
		pm.logger.Info(
			"Stopping process.",
			"type", "stop",
			"name", proc.Name,
		)
		res = proc.Stop()
	}

	// Stopping all process requires us to wait.
	pm.cwg.Wait()

	pm.logger.Info(
		"All processes stopped.",
		"type", "stop",
	)

	return res
}

// Return a list of all processes
func (pm *Manager) Processes() *[]*Process {
	return &pm.processes
}

// Return the number of processes that we are managing.
func (pm *Manager) Count() int {
	return len(pm.processes)
}

/* manager.go ends here. */
