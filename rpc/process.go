/*
 * process.go --- RPC process.
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

package rpc

import (
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"

	"context"
	"errors"
	"reflect"
	"sync"
)

const (
	addRpcable int = 1
)

type ManagerProc struct {
	sync.Mutex

	inst *Manager
}

func NewManagerProc(lgr logger.Logger, ctx context.Context, config *Config) *ManagerProc {
	return &ManagerProc{
		inst: NewManager(config, ctx, lgr),
	}
}

func (p *ManagerProc) start() {
	p.inst.Start()
}

func (p *ManagerProc) shutdown() {
	p.inst.Shutdown()
}

func (p *ManagerProc) stop(_ **process.State) {
	p.shutdown()
}

func (p *ManagerProc) Action(state **process.State) {
	var ps *process.State = *state

	cmd := ps.ReceiveBlocking()
	if cmd == nil {
		return
	}

	p.Lock()
	switch cmd.(*types.Pair).First {
	case addRpcable:
		thing := cmd.(*types.Pair).Second.(reflect.Type)
		ps.Send(p.inst.Add(thing))
	}
	p.Unlock()
}

func Spawn(mgr process.Manager, lgr logger.Logger, ctx context.Context, cnf *Config) (*process.Process, error) {
	name := "RPC Manager"

	if mgr == nil {
		return nil, errors.New("Process manager not provided!")
	}

	inst, found := mgr.Find(name)
	if found {
		return inst, nil
	}

	r := NewManagerProc(lgr, ctx, cnf)
	pr := mgr.Create(&process.Config{
		Name:     name,
		Interval: 0,
		Function: r.Action,
		OnStop:   r.stop,
	})

	go r.start()
	go pr.Run()

	return pr, nil
}

func Add(mgr process.Manager, t reflect.Type) (bool, error) {
	inst, found := mgr.Find("RPC Manager")
	if !found {
		return false, errors.New("Could not find RPC Manager process!")
	}

	inst.Send(types.NewPair(addRpcable, t))
	result := inst.Receive()

	return result.(bool), nil
}

/* process.go ends here. */
