/*
 * process.go --- Managed processes.
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

package process

import (
	"github.com/Asmodai/gohacks/logger"

	"context"
	"sync"
	"time"
)

const (
	EventLoopSleep time.Duration = 250 * time.Millisecond
)

// Callback function.
type CallbackFn func(**State)
type QueryFn func(interface{}) interface{}

/*

Process structure.

To use:

1) Create a config:

```go
  conf := &process.Config{
    Name:     "Windows 95",
    Interval: 10,        // 10 seconds.
    Function: func(state **State) {
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

*/
type Process struct {
	sync.Mutex

	Name     string        // Pretty name.
	Function CallbackFn    // `Action` callback.
	OnStart  CallbackFn    // `Start` callback.
	OnStop   CallbackFn    // `Stop` callback.
	OnQuery  QueryFn       // `Query` callback.
	Running  bool          // Is the process running?
	Interval time.Duration // `RunEvery` time interval.

	logger logger.Logger

	ctx    context.Context
	cancel context.CancelFunc
	wg     *sync.WaitGroup

	chanToState   chan interface{}
	chanFromState chan interface{}

	period time.Duration
	state  *State
}

// Create a new process with the given configuration and parent context.
func NewProcessWithContext(config *Config, parent context.Context) *Process {
	if config.Logger == nil {
		config.Logger = logger.NewDefaultLogger()
	}

	ctx, cancel := context.WithCancel(parent)

	p := &Process{
		Name:          config.Name,
		Function:      config.Function,
		OnStart:       config.OnStart,
		OnStop:        config.OnStop,
		OnQuery:       config.OnQuery,
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
func (p *Process) SetLogger(lgr logger.Logger) {
	p.logger = lgr
}

// Set the process's context.
func (p *Process) SetContext(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)

	p.ctx = ctx
	p.cancel = cancel
}

// Return the context for the process.
func (p *Process) Context() context.Context {
	return p.ctx
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

	// Execute startup callback if available.
	if p.OnStart != nil {
		p.OnStart(&p.state)
	}

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

// Query the running process.
//
// This allows interaction with the process's base object without using
// `Action`.
func (p *Process) Query(arg interface{}) interface{} {
	if p.OnQuery == nil {
		return nil
	}

	return p.OnQuery(arg)
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
