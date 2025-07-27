// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// process.go --- Dispatcher process.
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

package apiserver

// * Imports:

import (
	"sync"

	"github.com/gin-gonic/gin"
	"gitlab.com/tozd/go/errors"

	"github.com/Asmodai/gohacks/v1/events"
	"github.com/Asmodai/gohacks/v1/logger"
	"github.com/Asmodai/gohacks/v1/process"
)

// * Constants:

const (
	// Responder name.
	responderName string = "gohacks.dispatcher"

	// Responder type.
	responderType string = "apiserver.Dispatcher"

	// "GetRouter" command.
	CmdGetRouter string = "dispatcher.getrouter"
)

// * Variables:

var (
	// Triggered when no process manager instance is provided.
	ErrNoProcessManager = errors.Base("no process manager")

	// Triggered when no API dispatcher instance is provided.
	ErrNoDispatcherProc = errors.Base("no dispatcher process")

	// Triggered if the process responds with an unexpected data type.
	ErrWrongReturnType = errors.Base("wrong return type")
)

// * Code:

// ** Types:

// Dispatcher process.
type DispatcherProc struct {
	mu sync.RWMutex

	inst *Dispatcher
}

// ** Methods:

// Start the dispatcher process.
func (p *DispatcherProc) start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.inst.Start()
}

// Internal method to stop the API dispatcher process.
func (p *DispatcherProc) stop(_ *process.State) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.inst.Stop()
}

func (p *DispatcherProc) Name() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return responderName
}

func (p *DispatcherProc) Type() string {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return responderType
}

func (p *DispatcherProc) RespondsTo(event events.Event) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	switch val := event.(type) {
	case *events.Message:
		switch val.Command() {
		case CmdGetRouter:
			return true

		default:
			return false
		}

	default:
		return false
	}
}

func (p *DispatcherProc) Invoke(event events.Event) events.Event {
	cmd, ok := event.(*events.Message)
	if !ok {
		return event
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	switch cmd.Command() {
	case CmdGetRouter:
		return events.NewResponse(cmd, p.inst.GetRouter())

	default:
		return event
	}
}

// * Functions:

// Create a new dispatcher process.
func NewDispatcherProc(lgr logger.Logger, config *Config) *DispatcherProc {
	return &DispatcherProc{
		inst: NewDispatcher(lgr, config),
	}
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
		Name:           name,
		Interval:       0,
		FirstResponder: dispatch,
		OnStop:         dispatch.stop,
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

	event := inst.Invoke(events.NewMessage(CmdGetRouter, nil))

	result, ok := event.(*events.Response)
	if !ok {
		return nil, errors.WithStack(ErrWrongReturnType)
	}

	engine, okay := result.Response().(*gin.Engine)
	if !okay {
		return nil, errors.WithStack(ErrWrongReturnType)
	}

	return engine, nil
}

// process.go ends here.
