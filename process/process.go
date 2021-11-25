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
	"context"
	"log"
	"runtime"
	"sync"
	"time"
)

var (
	procs = []*Process{}
)

const (
	EventLoopSleep time.Duration = 250 * time.Millisecond
)

// Process callback function.
type CallbackFn func(context.Context, *IPC)

// Process stopping callback function.
type OnStopFn func()

/*

Process structure.

To use:

1) Create a config:

  conf := &process.Config{
    Name:     "Windows 95",
    Interval: 10,        // 10 seconds.
    Function: func(_ context.Context) {
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

	config   *Config
	running  bool               // Is the process running?
	interval time.Duration      // Real interval, as a time.Duration.
	period   time.Duration      // Remaining interval.
	ctx      context.Context    // Processes' context.
	cancelFn context.CancelFunc // Context 'cancel' function.
	ipc      *IPC               // IPC mechanism
}

// Create a new process with the given configuration.
func NewProcess(config *Config) *Process {
	duration := time.Duration(config.Interval) * time.Second

	// Set up a new cancelable context.
	ctx, cancel := context.WithCancel(context.Background())

	p := &Process{
		config:   config,
		running:  false,
		interval: duration,
		period:   duration,
		ipc:      NewIPC(ctx),
		ctx:      ctx,
		cancelFn: cancel,
	}

	// Set default callback if none is provided.
	if p.config.ActionFn == nil {
		p.config.ActionFn = p.nilFunction
	}

	return p
}

// Run the process with its action taking place on a continuous loop.
//
// Returns 'true' if the process has been started, or 'false' if it is
// already running.
func (p *Process) Run() bool {
	if p.running {
		return false
	}

	p.running = true

	// Are we to run on an interval?
	if p.config.Interval > 0 {
		log.Printf(
			"PROCESS: %s started, will invoke every %d second(s).\n",
			p.config.Name,
			p.interval/time.Second,
		)
		p.everyAction()

		return true
	}

	log.Printf("PROCESS: %s started.\n", p.config.Name)
	p.runAction()

	return true
}

// Stop the process.
//
// Returns 'true' if the process was successfully stopped, or 'false'
// if it was not running.
func (p *Process) Stop() bool {
	if !p.running {
		return false
	}

	p.cancelFn()

	p.running = false
	log.Printf("PROCESS: %s stopped.\n", p.config.Name)

	return true
}

// Is the process running?
func (p *Process) Running() bool {
	return p.running
}

// Send data to the process with blocking.
func (p *Process) Send(data interface{}) {
	p.ipc.ClientSend(data)
}

// Receive data from the process with blocking.
func (p *Process) Receive() (interface{}, bool) {
	return p.ipc.ClientReceive()
}

// Default action callback.
func (p *Process) nilFunction(_ context.Context, _ *IPC) {
	// noop.
}

// Run the configured action for this process.
func (p *Process) runAction() {
	p.Lock()
	defer p.Unlock()

	for {
		select {
		case <-p.ctx.Done():
			// Invoked when context is cancelled.
			{
				if p.config.StopFn != nil {
					p.config.StopFn()
				}

				return
			}

		default:
			// Non-blocking select!
		}

		if p.config.ActionFn != nil {
			ctx, cancel := context.WithCancel(p.ctx)

			p.config.ActionFn(ctx, p.ipc)
			cancel()
		}

		// Give time to both Go and OS for scheduler.
		runtime.Gosched()
		time.Sleep(EventLoopSleep)
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
			// Invoked when context is cancelled.
			{
				if p.config.StopFn != nil {
					p.config.StopFn()
				}

				return
			}

		case <-time.After(p.period):
			// Invoked when timer fires.
			break
		}

		started := time.Now()
		if p.config.ActionFn != nil {
			ctx, cancel := context.WithCancel(p.ctx)

			p.config.ActionFn(ctx, p.ipc)
			cancel()
		}
		finished := time.Now()

		duration := finished.Sub(started)
		p.period = p.interval - duration

		if p.period < 0 {
			log.Printf(
				"PROCESS: WARNING: Event loop for %s took longer than interval of %d!\n",
				p.config.Name,
				p.interval/time.Second,
			)
		}
	}
}

/* process.go ends here. */
