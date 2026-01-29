// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// manager.go --- Process manager.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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
//
// mock:yes
//go:generate go run github.com/Asmodai/gohacks/cmd/digen -pattern .
//di:gen basename=Manager key=gohacks/process@v1 type=Manager fallback=NewManager()

// * Comments:

// * Package:

package process

// * Imports:

import (
	"context"
	"fmt"
	"sync"

	"github.com/Asmodai/gohacks/logger"
)

// * Code:

// ** Interface:

/*
Process manager structure.

To use,

1) Create a new process manager:

```go

	procmgr := process.NewManager()

```

2) Create your process configuration:

```go

	conf := &process.Config{
	  Name:     "Windows 95",
	  Interval: 10, // seconds
	  Function: func(state *State) {
	    // Crash or something.
	  }
	}

```

3) Create the process itself.

```go

	proc := procmgr.Create(conf)

```

4) Run the process.

```go

	procmgr.Run("Windows 95")

```

/or/

```go

	proc.Run()

```

Manager is optional, as you can create processes directly.
*/
type Manager interface {
	Logger() logger.Logger
	SetContext(context.Context)
	SetLogger(logger.Logger)
	Context() context.Context
	Create(*Config) *Process
	Add(*Process)
	Find(string) (*Process, bool)
	Run(string) bool
	Stop(string) bool
	StopAll() StopAllResults
	Processes() []*Process
	Count() int
}

// ** Types:

type StopAllResults map[string]bool

// Process map type.
type processMap map[string]*Process

type manager struct {
	mu sync.RWMutex

	processes processMap
	logger    logger.Logger
	parent    context.Context
	ctx       context.Context
	cancel    context.CancelFunc
	cwg       *sync.WaitGroup
}

// ** Methods:

// Set the process manager's logger.
func (pm *manager) SetLogger(lgr logger.Logger) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.logger = lgr
}

// Return the manager's logger.
func (pm *manager) Logger() logger.Logger {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.logger
}

// Set the process manager's context.
func (pm *manager) SetContext(parent context.Context) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	ctx, cancel := context.WithCancel(parent)

	pm.parent = parent
	pm.ctx = ctx
	pm.cancel = cancel
}

// Get the process manager's context.
func (pm *manager) Context() context.Context {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.ctx
}

// Create a new process with the given configuration.
func (pm *manager) Create(config *Config) *Process {
	proc := NewProcessWithContext(pm.ctx, config)
	proc.setWaitGroup(pm.cwg)

	pm.Add(proc)

	return proc
}

// Add an existing process to the manager.
func (pm *manager) Add(proc *Process) {
	if proc == nil {
		return
	}

	pm.mu.Lock()
	defer pm.mu.Unlock()

	proc.setLogger(pm.logger)

	_, found := pm.processes[proc.name]
	if found {
		// No, absolutely do not allow this... and be violent about
		// it.  Processes need to be unique.
		panic(fmt.Sprintf(
			"Attempt made to replace an existing process '%s'",
			proc.name))
	}

	pm.processes[proc.name] = proc
}

// Find and return the given process, or nil if not found.
func (pm *manager) Find(name string) (*Process, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	proc, found := pm.processes[name]
	if !found {
		return nil, false
	}

	return proc, true
}

// Run the named process.
//
// Returns 'false' if the process is not found;  otherwise returns
// the result of the process execution.
func (pm *manager) Run(name string) bool {
	proc, found := pm.Find(name)
	if !found {
		return false
	}

	pm.mu.Lock()
	proc.setContext(pm.ctx)
	pm.mu.Unlock()

	return proc.Run()
}

// Stop the given process.
//
// Returns 'true' if the process has been stopped; otherwise 'false'.
func (pm *manager) Stop(name string) bool {
	proc, found := pm.Find(name)
	if !found {
		return false
	}

	// Stopping one process doesn't require us to wait for the group.
	return proc.Stop()
}

// Stop all processes.
//
// Returns a `map[string]bool` value where the keys are the names of the
// currently-managed processes and the result of invoking `Stop` on them.
func (pm *manager) StopAll() StopAllResults {
	res := StopAllResults{}

	pm.logger.Info(
		"Stopping all processes.",
		"type", "stop",
	)

	// Block is being used here to highlight the lock.
	pm.mu.RLock()
	{
		// This is better than invoking the context's cancel, as
		// it allows cleanup to be executed.
		for _, proc := range pm.processes {
			pm.logger.Info(
				"Stopping process.",
				"type", "stop",
				"name", proc.name,
			)

			res[proc.name] = proc.Stop()
		}
	}
	pm.mu.RUnlock()

	// Stopping all process requires us to wait.
	pm.cwg.Wait()

	pm.logger.Info(
		"All processes stopped.",
		"type", "stop",
	)

	return res
}

// Return a list of all processes.
func (pm *manager) Processes() []*Process {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var processes = make([]*Process, 0, len(pm.processes))

	for _, proc := range pm.processes {
		processes = append(processes, proc)
	}

	return processes
}

// Return the number of processes that we are managing.
func (pm *manager) Count() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return len(pm.processes)
}

// ** Functions:

// Create a new process manager with a given parent context.
func NewManagerWithContext(parent context.Context) Manager {
	lgr, err := logger.GetLogger(parent)
	if err != nil {
		lgr = logger.NewDefaultLogger()
	}

	ctx, cancel := context.WithCancel(parent)

	return &manager{
		processes: make(processMap),
		logger:    lgr,
		ctx:       ctx,
		cancel:    cancel,
		cwg:       &sync.WaitGroup{},
	}
}

// Create a new process manager.
func NewManager() Manager {
	return NewManagerWithContext(context.Background())
}

// * manager.go ends here.
