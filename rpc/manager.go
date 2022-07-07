/*
 * manager.go --- RPC manager.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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
	"context"
	"net"
	"net/http"
	gorpc "net/rpc"
	"reflect"

	"github.com/Asmodai/gohacks/logger"
)

type Manager struct {
	ctx     context.Context
	started bool

	entities map[reflect.Type]IRPCAble

	conf     *Config
	listener net.Listener
	logger   logger.ILogger
}

func NewManager(cnf *Config, ctx context.Context, lgr logger.ILogger) *Manager {
	return &Manager{
		ctx:      ctx,
		started:  false,
		entities: map[reflect.Type]IRPCAble{},
		conf:     cnf,
		listener: nil,
		logger:   lgr,
	}
}

func (m *Manager) Add(t reflect.Type) bool {
	if _, ok := m.entities[t]; ok {
		return false
	}

	m.entities[t] = NewRPCAble(t)

	if err := gorpc.Register(m.entities[t]); err != nil {
		m.logger.Warn(
			"Failed to register RPC-able.",
			"entity", t,
			"err", err.Error(),
		)

		return false
	}

	return true
}

func (m *Manager) Start() {
	if m.started {
		return
	}

	var err error

	gorpc.HandleHTTP()

	m.listener, err = net.Listen(m.conf.Proto, m.conf.Addr)
	if err != nil {
		m.logger.Fatal(
			"Fatal RPC error",
			"err", err.Error(),
		)
	}

	if err = http.Serve(m.listener, nil); err != nil {
		m.logger.Fatal(
			"Could not serve HTTP RPC handler.",
			"err", err.Error(),
		)
	}

	m.started = true
}

func (m *Manager) Shutdown() {
	if !m.started {
		return
	}

	m.listener.Close()
	m.started = false
}

/* rpcmanager.go ends here. */
