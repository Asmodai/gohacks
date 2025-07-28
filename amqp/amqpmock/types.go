// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// types.go --- AMQP mock types.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
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
//
//

// * Package:

package amqpmock

// * Imports:

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Asmodai/gohacks/amqp/amqpshim"
	goamqp "github.com/rabbitmq/amqp091-go"
)

// * Code:

// ** Call Log Types:

type CallLogList []any

type CallLog map[string]CallLogList

func (obj CallLog) Dump() string {
	var sbld strings.Builder

	for fname, calls := range obj {
		fmt.Fprintf(&sbld, "%s:\n", fname)

		for idx, args := range calls {
			arglst, ok := args.(CallLogList)
			if !ok {
				fmt.Fprintf(
					&sbld,
					"  #%d: <invalid call args: %#v>\n",
					idx,
					args,
				)

				continue
			}

			fmt.Fprintf(&sbld, "  #%d: (", idx)

			for argno, arg := range arglst {
				if argno > 0 {
					sbld.WriteString(", ")
				}

				fmt.Fprintf(&sbld, "%#v", arg)
			}

			sbld.WriteString(")\n")
		}
	}

	return sbld.String()
}

// ** Mock Return Types:

type ErrorResults struct {
	Error error
}

type BoolResults struct {
	Value bool
}

type UInt64Results struct {
	Value uint64
}

type AddrResults struct {
	Addr net.Addr
}

type ChannelResults struct {
	Channel amqpshim.Channel
	Error   error
}

type ConnectionStateResults struct {
	TLS tls.ConnectionState
}

type ConsumeResults struct {
	Channel <-chan goamqp.Delivery
	Error   error
}

type GetResults struct {
	Message goamqp.Delivery
	Ok      bool
	Error   error
}

type NotifyBlockedResults struct {
	BlockingChan chan goamqp.Blocking
}

type NotifyCancelResults struct {
	StringChan chan string
}

type NotifyCloseResults struct {
	ErrorChan chan *goamqp.Error
}

type NotifyConfirmResults struct {
	AckChan  chan uint64
	NackChan chan uint64
}

type NotifyFlowResults struct {
	FlowChan chan bool
}

type NotifyPublishResults struct {
	ConfirmChan chan goamqp.Confirmation
}

type NotifyReturnResults struct {
	ReturnChan chan goamqp.Return
}

type PublishDeferredResults struct {
	Confirmation *goamqp.DeferredConfirmation
	Error        error
}

type QueueDeclareResults struct {
	Queue goamqp.Queue
	Error error
}

type QueueDeleteResults struct {
	Purged int
	Error  error
}

// ** Function Types:

type SimpleErrorFn func() error

type SimpleBoolFn func() bool

type SimpleUInt64Fn func() uint64

type SimpleAddrFn func() net.Addr

type AckFn func(uint64, bool) error

type CancelFn func(string, bool) error

type ChannelFn func() (amqpshim.Channel, error)

type CloseDeadlineFn func(time.Time) error

type ConfirmFn func(bool) error

type ConnectionStateFn func() tls.ConnectionState

type ConsumeFn func(
	string,
	string,
	bool,
	bool,
	bool,
	bool,
	goamqp.Table,
) (<-chan goamqp.Delivery, error)

type ConsumeContextFn func(
	context.Context,
	string,
	string,
	bool,
	bool,
	bool,
	bool,
	goamqp.Table,
) (<-chan goamqp.Delivery, error)

type ExchangeBindFn func(string, string, string, bool, goamqp.Table) error

type ExchangeDeclareFn func(
	string,
	string,
	bool,
	bool,
	bool,
	bool,
	goamqp.Table,
) error

type ExchangeDeleteFn func(string, bool, bool) error

type ExchangeUnbindFn func(string, string, string, bool, goamqp.Table) error

type FlowFn func(bool) error

type GetFn func(string, bool) (goamqp.Delivery, bool, error)

type NackFn func(uint64, bool, bool) error

type NotifyBlockedFn func(chan goamqp.Blocking) chan goamqp.Blocking

type NotifyCancelFn func(chan string) chan string

type NotifyCloseFn func(chan *goamqp.Error) chan *goamqp.Error

type NotifyConfirmFn func(chan uint64, chan uint64) (chan uint64, chan uint64)

type NotifyFlowFn func(chan bool) chan bool

type NotifyPublishFn func(chan goamqp.Confirmation) chan goamqp.Confirmation

type NotifyReturnFn func(chan goamqp.Return) chan goamqp.Return

type PublishFn func(string, string, bool, bool, goamqp.Publishing) error

type PublishContextFn func(
	context.Context,
	string,
	string,
	bool,
	bool,
	goamqp.Publishing,
) error

type PublishDeferredFn func(
	string,
	string,
	bool,
	bool,
	goamqp.Publishing,
) (*goamqp.DeferredConfirmation, error)

type PublishDeferredContextFn func(
	context.Context,
	string,
	string,
	bool,
	bool,
	goamqp.Publishing,
) (*goamqp.DeferredConfirmation, error)

type QosFn func(int, int, bool) error

type QueueBindFn func(string, string, string, bool, goamqp.Table) error

type QueueDeclareFn func(
	string,
	bool,
	bool,
	bool,
	bool,
	goamqp.Table,
) (goamqp.Queue, error)

type QueueDeleteFn func(string, bool, bool, bool) (int, error)

type QueuePurgeFn func(string, bool) (int, error)

type QueueUnbindFn func(string, string, string, goamqp.Table) error

type RejectFn func(uint64, bool) error

type UpdateSecretFn func(string, string) error

// * types.go ends here.
