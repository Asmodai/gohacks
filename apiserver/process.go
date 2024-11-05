// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// process.go --- Dispatcher process.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package apiserver

import (
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"

	"github.com/gin-gonic/gin"
	"gitlab.com/tozd/go/errors"

	"sync"
)

const (
	// API command magic numbers.
	getRouter int = 1
)

var (
	// Triggered when no process manager instance is provided.
	ErrNoProcessManager = errors.Base("no process manager")

	// Triggered when no API dispatcher instance is provided.
	ErrNoDispatcherProc = errors.Base("no dispatcher process")

	// Triggered if the process responds with an unexpected data type.
	ErrWrongReturnType = errors.Base("wrong return type")
)

// Dispatcher process.
type DispatcherProc struct {
	sync.Mutex

	inst *Dispatcher
}

// Create a new dispatcher process.
func NewDispatcherProc(lgr logger.Logger, config *Config) *DispatcherProc {
	return &DispatcherProc{
		inst: NewDispatcher(lgr, config),
	}
}

// Start the dispatcher process.
func (p *DispatcherProc) start() {
	p.inst.Start()
}

// Action invoked when the dispatcher process receives a message.
func (p *DispatcherProc) Action(state **process.State) {
	var pst = *state

	cmd, ok := pst.ReceiveBlocking().(*types.Pair)
	if !ok && cmd != nil {
		pst.Logger().Fatal(
			"Received message that is not of type types.Pair",
			"cmd", cmd,
		)
	}

	if cmd == nil {
		return
	}

	p.Lock()
	//nolint:gocritic
	switch cmd.First {
	// Caller wants to know what our router is.
	case getRouter:
		pst.Send(process.NewActionResult(p.inst.GetRouter(), nil))
	}
	p.Unlock()
}

// Internal method to stop the API dispatcher process.
func (p *DispatcherProc) stop(_ **process.State) {
	p.inst.Stop()
}

// Spawn an API dispatcher process.
func Spawn(mgr process.Manager, lgr logger.Logger, config *Config) (*process.Process, error) {
	name := "Dispatcher"

	if mgr == nil {
		return nil, errors.WithStack(ErrNoProcessManager)
	}

	inst, found := mgr.Find(name)
	if found {
		return inst, nil
	}

	dispatch := NewDispatcherProc(lgr, config)
	conf := &process.Config{
		Name:     name,
		Interval: 0,
		Function: dispatch.Action,
		OnStop:   dispatch.stop,
	}
	proc := mgr.Create(conf)

	// Start dispatcher.
	dispatch.start()

	// Start goroutine.
	go proc.Run()

	return proc, nil
}

// Get the router currently in use by the registered dispatcher process.
//
// Should no dispatcher process be in the process manager, then
// `ErrNoDispatcherProc` will be returned.
//
// Should the dispatcher process return an unexpected value type, then
// `ErrWrongReturnType` will be returned.
func GetRouter(mgr process.Manager) (*gin.Engine, error) {
	inst, found := mgr.Find("Dispatcher")
	if !found {
		return nil, errors.WithStack(ErrNoDispatcherProc)
	}

	inst.Send(types.NewPair(getRouter, nil))

	result, okay := inst.Receive().(*process.ActionResult)
	if !okay {
		return nil, errors.WithStack(ErrWrongReturnType)
	}

	engine, okay := result.Value.(*gin.Engine)
	if !okay {
		return nil, errors.WithStack(ErrWrongReturnType)
	}

	return engine, nil
}

// process.go ends here.
