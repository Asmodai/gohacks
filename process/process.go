/*
 * process.go --- Managed processes.
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
	"time"
)

var (
	procs = []*Process{}
)

const (
	EventLoopSleep time.Duration = 250 * time.Millisecond
)

// Callback function.
type CallbackFn func(**State)

/*

Process structure.

To use:

1) Create a config:

  conf := &process.Config{
    Name:     "Windows 95",
    Interval: 10,        // 10 seconds.
    Function: func(state **State) {
      // Crash or something.
    },
  }

2) Create a process:

  proc := process.NewProcess(conf)

3) Run the process:

  go proc.Run()

4) Send data to the process:

  proc.Send("Blue Screen of Death")

5) Read data from the process:

  data := proc.Receive()

6) Stop the process

  proc.Stop()

*/
type Process struct {
	sync.Mutex

	Name     string        // Pretty name.
	Function CallbackFn    // `Action` callback.
	OnStop   CallbackFn    // `Stop` callback.
	Running  bool          // Is the process running?
	Interval time.Duration // `RunEvery` time interval.

	logger logger.ILogger

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	chanToState   chan interface{}
	chanFromState chan interface{}

	period  time.Duration
	state   *State
	manager *Manager
}

// Create a new process with the given configuration and parent context.
func NewProcessWithContext(config *Config, parent context.Context) *Process {
	if config.Logger == nil {
		config.Logger = logger.NewDefaultLogger("")
	}

	ctx, cancel := context.WithCancel(parent)

	p := &Process{
		Name:          config.Name,
		Function:      config.Function,
		OnStop:        config.OnStop,
		Running:       false,
		Interval:      (time.Duration)(config.Interval) * time.Second,
		period:        (time.Duration)(config.Interval) * time.Second,
		logger:        config.Logger,
		ctx:           ctx,
		cancel:        cancel,
		wg:            &sync.WaitGroup{},
		chanToState:   make(chan interface{}, 1),
		chanFromState: make(chan interface{}, 1),
		state:         &State{},
	}

	if config.Function == nil {
		p.Function = p.nilFunction
	}

	p.state.parent = p

	return p
}

// Create a new process with the given configuration.
func NewProcess(config *Config) *Process {
	return NewProcessWithContext(config, context.TODO())
}

// Set the process's logger.
func (p *Process) SetLogger(lgr logger.ILogger) {
	p.logger = lgr
}

// Set the process's context.
func (p *Process) SetContext(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)

	p.ctx = ctx
	p.cancel = cancel
}

// Set the process's wait group.
func (p *Process) SetWaitGroup(wg *sync.WaitGroup) {
	p.wg = wg
}

// Run the process with its action taking place on a continuous loop.
//
// Returns 'true' if the process has been started, or 'false' if it is
// already running.
func (p *Process) Run() bool {
	if p.Running {
		return false
	}

	p.Running = true

	// Add child to wait group
	p.wg.Add(1)

	// Are we to run on an interval?
	if p.Interval > 0 {
		p.logger.Info(
			"Process started.",
			"type", "start",
			"name", p.Name,
			"interval", p.Interval.Round(time.Second),
		)
		p.everyAction()
		return true
	}

	p.logger.Info(
		"Process started.",
		"type", "start",
		"name", p.Name,
	)

	p.runAction()

	return true
}

// Stop the process.
//
// Returns 'true' if the process was successfully stopped, or 'false'
// if it was not running.
func (p *Process) Stop() bool {
	if !p.Running {
		return false
	}

	p.Send(nil)
	p.cancel()
	p.Running = false

	return true
}

// Send data to the process with blocking.
func (p *Process) Send(data interface{}) {
	p.chanToState <- data
}

// Send data to the process without blocking.
func (p *Process) SendNonBlocking(data interface{}) {
	select {
	case p.chanToState <- data:
	default:
	}
}

// Receive data from the process with blocking.
func (p *Process) Receive() interface{} {
	return <-p.chanFromState
}

// Receive data from the process without blocking.
func (p *Process) ReceiveNonBlocking() (interface{}, bool) {
	select {
	case data := <-p.chanFromState:
		return data, true

	default:
	}

	return nil, false
}

// Default action callback.
func (p *Process) nilFunction(state **State) {
}

// Internal callback invoked upon process stop.
func (p *Process) internalStop() {
	p.logger.Info(
		"Process stopped.",
		"type", "stop",
		"name", p.Name,
	)

	// Set wait as done.
	p.wg.Done()
}

// Run the configured action for this process.
func (p *Process) runAction() {
	p.Lock()
	defer p.Unlock()

	for {
		select {
		case <-p.ctx.Done():
			if p.OnStop != nil {
				p.OnStop(&p.state)
			}
			p.internalStop()
			return

		default:
		}

		if p.Function != nil {
			p.Function(&p.state)
		} else {
			time.Sleep(EventLoopSleep)
		}
	}
}

// Run the configured action for this process.
//
// Identical to 'runAction', except for the fact that this sleeps,
// giving the appearance of something that runs on an interval.
func (p *Process) everyAction() {
	p.Lock()
	defer p.Unlock()

	for {
		select {
		case <-p.ctx.Done():
			if p.OnStop != nil {
				p.OnStop(&p.state)
			}
			p.internalStop()
			return

		case <-time.After(p.period):
			break
		}

		started := time.Now()
		if p.Function != nil {
			p.Function(&p.state)
		}
		finished := time.Now()

		duration := finished.Sub(started)
		p.period = p.Interval - duration
	}
}

/* process.go ends here. */
