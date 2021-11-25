/*
 * process.go --- Dispatcher process.
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

package apiserver

import (
	"github.com/Asmodai/gohacks/di"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"

	"github.com/gin-gonic/gin"

	"context"
	"fmt"
	"sync"
	"time"
)

const (
	getRouter int = 1
)

type DispatcherProc struct {
	sync.Mutex

	inst *Dispatcher
}

func NewDispatcherProc(config *Config) *DispatcherProc {
	return &DispatcherProc{
		inst: NewDispatcher(config),
	}
}

func (p *DispatcherProc) start() {
	p.inst.Start()
}

func (p *DispatcherProc) Action(ctx context.Context, ipc *process.IPC) {
	/*
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
	*/

	cmd, ok := ipc.Receive()
	if ok {
		p.Lock()
		switch cmd.(*types.Pair).First {
		case getRouter:
			ipc.Send(p.inst.GetRouter())
			break
		}
		p.Unlock()
	}
}

func (p *DispatcherProc) stop() {
	p.inst.Stop()
}

func Spawn(config *Config) (*process.Process, error) {
	name := "Dispatcher"

	dism := di.GetInstance()

	mgr, found := dism.Get("ProcMgr")
	if !found {
		return nil, types.NewError(
			"DISPATCHER",
			"Could not locate 'ProcMgr' service.",
		)
	}

	inst, found := mgr.(*process.Manager).Find(name)
	if found {
		return inst, nil
	}

	dispatch := NewDispatcherProc(config)
	conf := &process.Config{
		Name:     name,
		Interval: 0,
		ActionFn: dispatch.Action,
		StopFn:   dispatch.stop,
	}
	pr := mgr.(*process.Manager).Create(conf)

	dispatch.start()
	go pr.Run()

	return pr, nil
}

func GetRouter(mgr process.IManager) (*gin.Engine, error) {
	inst, found := mgr.Find("Dispatcher")
	if !found {
		return nil, fmt.Errorf("Could not find Dispatcher process!")
	}

	// Ugh
	for !inst.Running() {
		time.Sleep(10 * time.Millisecond)
	}

	inst.Send(types.NewPair(getRouter, nil))
	result, _ := inst.Receive()

	return result.(*gin.Engine), nil
}

/* process.go ends here. */
