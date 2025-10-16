// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// process.go --- Managed processes.
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

package process

// * Imports:

import (
	"context"
	"fmt"
	godebug "runtime/debug"
	"sync"
	"time"

	"github.com/Asmodai/gohacks/events"
	"github.com/Asmodai/gohacks/logger"
)

// * Constants:

const (
	eventLoopSleep    time.Duration = 150 * time.Millisecond
	channelBufferSize int           = 1
	processTypeString string        = "process.Process"
)

// ** Types:

// Callback function.
type CallbackFn func(*State)

type QueryFn func(any) any

/*
Process structure.

To use:

1) Create a config:

```go

	conf := &process.Config{
	  Name:     "Windows 95",
	  Interval: 10,        // 10 seconds.
	  Function: func(state *State) {
	    // Crash or something.
	  },
	}

```

2) Create a process:

```go

	proc := process.NewProcess(conf)

```

3) Run the process:

```go

	go proc.Run()

```

4) Send data to the process:

```go

	proc.Send("Blue Screen of Death")

```

5) Read data from the process:

```go

	data := proc.Receive()

```

6) Stop the process

```go

	proc.Stop()

```

	will stop the process.
*/
type Process struct {
	mu sync.RWMutex

	name     string        // Pretty name.
	function CallbackFn    // `Action` callback.
	onStart  CallbackFn    // `Start` callback.
	onStop   CallbackFn    // `Stop` callback.
	onQuery  QueryFn       // `Query` callback.
	interval time.Duration // `RunEvery` time interval.

	logger logger.Logger
	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	running bool // Is the process running?
	period  time.Duration
	state   *State
}

// ** Methods:

// Set the process's logger.
//
// This should only be called by the process manager at process startup.
func (p *Process) setLogger(lgr logger.Logger) {
	p.logger = lgr
}

// Set the process's context.
//
// This should only be called by the process manager at process startup.
func (p *Process) setContext(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)

	p.ctx = ctx
	p.cancel = cancel
}

// Return the context for the process.
func (p *Process) Context() context.Context {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.ctx
}

// Set the process's wait group.
//
// This should only be called by the process manager at process startup.
func (p *Process) setWaitGroup(wg *sync.WaitGroup) {
	p.wg = wg
}

// Is the process running?
func (p *Process) Running() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.running
}

// Run the process with its action taking place on a continuous loop.
//
// Returns 'true' if the process has been started, or 'false' if it is
// already running.
func (p *Process) Run() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return false
	}

	p.running = true

	// Add child to wait group
	p.wg.Add(1)

	// Wrap everything up so it can be recovered.
	go func() {
		defer p.wg.Done()

		defer func() {
			if r := recover(); r != nil {
				p.logger.Info(
					"Process panicked!",
					"type", "panic",
					"name", p.name,
					"recovery", r,
					"stack", godebug.Stack(),
				)
			}
		}()

		// Execute startup callback if available.
		if p.onStart != nil {
			p.onStart(p.state)
		}

		// Are we to run on an interval?
		if p.interval > 0 {
			p.logger.Info(
				"Process started.",
				"type", "start",
				"name", p.name,
				"interval", p.interval.Round(time.Second),
			)
			p.everyAction()

			return
		}

		p.logger.Info(
			"Process started.",
			"type", "start",
			"name", p.name,
		)

		p.runAction()
	}()

	return true
}

// Stop the process.
//
// Returns 'true' if the process was successfully stopped, or 'false'
// if it was not running.
func (p *Process) Stop() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return false
	}

	p.cancel()
	p.running = false

	return true
}

// Query the running process.
//
// This allows interaction with the process's base object without using
// `Action`.
func (p *Process) Query(arg any) any {
	if p.onQuery == nil {
		return nil
	}

	return p.onQuery(arg)
}

// Default action callback.
func (p *Process) nilFunction(_ *State) {
}

// Internal callback invoked upon process stop.
func (p *Process) internalStop() {
	p.logger.Info(
		"Process stopped.",
		"type", "stop",
		"name", p.name,
	)
}

// Run the configured action for this process.
func (p *Process) runAction() {
	for {
		select {
		case <-p.ctx.Done():
			p.mu.Lock()
			defer p.mu.Unlock()

			if p.onStop != nil {
				p.onStop(p.state)
			}

			p.internalStop()

			return

		default:
		}

		if p.function != nil {
			p.function(p.state)
		}

		// Give time back to the scheduler.
		time.Sleep(eventLoopSleep)
	}
}

// Run the configured action for this process.
//
// Identical to 'runAction', except for the fact that this sleeps,
// giving the appearance of something that runs on an interval.
func (p *Process) everyAction() {
	for {
		select {
		case <-p.ctx.Done():
			p.mu.Lock()
			defer p.mu.Unlock()

			if p.onStop != nil {
				p.onStop(p.state)
			}

			p.internalStop()

			return

		case <-time.After(p.period):
			break
		}

		started := time.Now()

		if p.function != nil {
			p.function(p.state)
		}

		finished := time.Now()
		duration := finished.Sub(started)

		p.period = p.interval - duration
		if p.period < 0 {
			p.logger.Warn(
				"Period less than 0.  Process doing too much?",
				"name", p.name,
				"period", p.period,
				"start", started,
				"finished", finished,
				"duration", duration,
			)

			p.period = 0
		}
	}
}

func (p *Process) Name() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.name
}

func (p *Process) Type() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return processTypeString
}

func (p *Process) RespondsTo(event events.Event) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.state.RespondsTo(event)
}

func (p *Process) Invoke(event events.Event) events.Event {
	p.mu.Lock()
	defer p.mu.Unlock()

	ret, _ := p.state.Invoke(event)

	return ret
}

// ** Functions:

// Create a new process with the given configuration.
func NewProcess(config *Config) *Process {
	return NewProcessWithContext(context.Background(), config)
}

// Create a new process with the given configuration and parent context.
func NewProcessWithContext(parent context.Context, config *Config) *Process {
	lgr, err := logger.GetLogger(parent)
	if err != nil {
		lgr = logger.NewDefaultLogger()
	}

	ctx, cancel := context.WithCancel(parent)

	proc := &Process{
		name:     config.Name,
		function: config.Function,
		onStart:  config.OnStart,
		onStop:   config.OnStop,
		onQuery:  config.OnQuery,
		running:  false,
		interval: config.Interval.Duration(),
		period:   config.Interval.Duration(),
		logger:   lgr,
		ctx:      ctx,
		cancel:   cancel,
		wg:       &sync.WaitGroup{},
		state:    newState(config.Name),
	}

	if config.Function == nil {
		proc.function = proc.nilFunction
	}

	if config.Responder != nil {
		_, err := proc.state.responders.Add(config.Responder)
		if err != nil {
			panic(fmt.Sprintf(
				"Could not add responder for process %s: %#v",
				config.Name,
				config.Responder))
		}
	}

	proc.state.parent = proc

	return proc
}

// * process.go ends here.
