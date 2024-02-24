/*
 * process.go --- Exec process.
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

package exec

import (
	"context"
	"errors"

	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/types"
)

const (
	CMD_SET_CTX = iota
	CMD_SET_PATH
	CMD_SET_ARGS
	CMD_SPAWN
	CMD_CHECK
	CMD_KILL_ALL
)

type Process struct {
	manager *Manager
	logger  logger.ILogger
	conf    *Config
}

func NewProcess(lgr logger.ILogger, cnf *Config) *Process {
	return &Process{
		manager: NewManager(lgr, cnf.Count, cnf.Base),
		logger:  lgr,
		conf:    cnf,
	}
}

func (p *Process) Action(state **process.State) {
	var ps *process.State = *state

	cmd := ps.ReceiveBlocking()
	if cmd == nil {
		return
	}

	switch cmd.(*types.Pair).First {
	case CMD_SET_CTX:
		p.SetContext(cmd.(*types.Pair).Second.(context.Context))

	case CMD_SET_PATH:
		p.SetPath(cmd.(*types.Pair).Second.(string))

	case CMD_SET_ARGS:
		p.SetArgs(cmd.(*types.Pair).Second.([]string)...)

	case CMD_SPAWN:
		p.manager.Spawn()

	case CMD_CHECK:
		p.manager.Check()

	case CMD_KILL_ALL:
		p.manager.KillAll()
	}
}

func (p *Process) SetPath(val string) {
	p.manager.SetPath(val)
}

func (p *Process) SetArgs(args ...string) {
	p.manager.SetArgs(NewArgs(p.conf.Base, args))
}

func (p *Process) SetCount(val int) {
	p.manager.SetCount(val)
}

func (p *Process) SetContext(val context.Context) {
	p.manager.SetContext(val)
}

func Spawn(mgr process.IManager, lgr logger.ILogger, cnf *Config) (*process.Process, error) {
	name := "ExecManager"

	inst, found := mgr.Find(name)
	if found {
		return inst, nil
	}

	p := NewProcess(lgr, cnf)
	conf := &process.Config{
		Name:     name,
		Function: p.Action,
	}
	pr := mgr.Create(conf)

	go pr.Run()

	return pr, nil
}

func sendSimple(mgr process.IManager, arg *types.Pair) error {
	inst, found := mgr.Find("ExecManager")
	if !found {
		return errors.New("Could not find exec manager process!")
	}

	inst.Send(arg)

	return nil
}

func SetContext(mgr process.IManager, ctx context.Context) error {
	return sendSimple(mgr, types.NewPair(CMD_SET_CTX, ctx))
}

func SetPath(mgr process.IManager, path string) error {
	return sendSimple(mgr, types.NewPair(CMD_SET_PATH, path))
}

func SetArgs(mgr process.IManager, args ...string) error {
	return sendSimple(mgr, types.NewPair(CMD_SET_ARGS, args))
}

func SpawnProcs(mgr process.IManager) error {
	return sendSimple(mgr, types.NewPair(CMD_SPAWN, nil))
}

func CheckProcs(mgr process.IManager) error {
	return sendSimple(mgr, types.NewPair(CMD_CHECK, nil))
}

func KillAllProcs(mgr process.IManager) error {
	return sendSimple(mgr, types.NewPair(CMD_KILL_ALL, nil))
}

/* process.go ends here. */
