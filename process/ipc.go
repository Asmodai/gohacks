/*
 * ipc.go --- Process IPC mechanism.
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
	"github.com/Asmodai/gohacks/types"

	"context"
	"time"
)

const (
	ReceiveDelaySleep time.Duration = 5 * time.Second
)

type IPC struct {
	ctx    context.Context
	input  *types.Mailbox
	output *types.Mailbox
}

func NewIPC(ctx context.Context) *IPC {
	return &IPC{
		ctx:    ctx,
		input:  types.NewMailbox(),
		output: types.NewMailbox(),
	}
}

// Server --> Client
func (i *IPC) ClientSend(data interface{}) {
	i.input.Put(data)
}

// Client --> Server
func (i *IPC) Send(data interface{}) {
	i.output.Put(data)
}

// Client <-- Server
func (i *IPC) ClientReceive() (interface{}, bool) {
	return i.doReceive(i.output)
}

// Server <-- Client
func (i *IPC) Receive() (interface{}, bool) {
	return i.doReceive(i.input)
}

// Do the actual work for a receive.
//
// This will try to be nice to both the Go scheduler and the underlaying
// OS scheduler.
func (i *IPC) doReceive(mbox *types.Mailbox) (interface{}, bool) {
	ctx, cancel := context.WithTimeout(i.ctx, types.DefaultCtxDeadline)
	defer cancel()

	var val interface{}
	var ok bool = false

	for !ok {
		select {
		case <-ctx.Done():
			return nil, false

		case <-i.ctx.Done():
			return nil, false

		default:
		}

		time.Sleep(ReceiveDelaySleep)

		val, ok = mbox.GetWithContext(ctx)
	}

	return val, ok
}

/* ipc.go ends here. */
