/*
 * process.go --- Dispatcher process.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
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

package apiserver

import (
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"

	"github.com/gin-gonic/gin"

	"fmt"
	"sync"
)

const (
	getRouter int = 1
)

type DispatcherProc struct {
	sync.Mutex

	inst *Dispatcher
}

func NewDispatcherProc(lgr logger.ILogger, config *Config) *DispatcherProc {
	return &DispatcherProc{
		inst: NewDispatcher(lgr, config),
	}
}

func (p *DispatcherProc) start() {
	p.inst.Start()
}

func (p *DispatcherProc) Action(state **process.State) {
	var ps *process.State = *state

	cmd := ps.ReceiveBlocking()
	if cmd == nil {
		return
	}

	p.Lock()
	switch cmd.(*types.Pair).First {
	case getRouter:
		ps.Send(p.inst.GetRouter())
		break
	}
	p.Unlock()
}

func (p *DispatcherProc) stop(state **process.State) {
	p.inst.Stop()
}

func Spawn(mgr process.IManager, lgr logger.ILogger, config *Config) (*process.Process, error) {
	name := "Dispatcher"

	if mgr == nil {
		return nil, fmt.Errorf("Process manager not provided!")
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
	pr := mgr.Create(conf)

	dispatch.start()
	go pr.Run()

	return pr, nil
}

func GetRouter(mgr process.IManager) (*gin.Engine, error) {
	inst, found := mgr.Find("Dispatcher")
	if !found {
		return nil, fmt.Errorf("Could not find Dispatcher process!")
	}

	inst.Send(types.NewPair(getRouter, nil))
	result := inst.Receive()

	return result.(*gin.Engine), nil
}

/* process.go ends here. */
