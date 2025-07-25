// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// interface.go --- AMQP "driver" interface.
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
//
//mock:yes

// * Comments:
//
//

// * Package:

package amqpshim

// * Imports:

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	goamqp "github.com/rabbitmq/amqp091-go"
)

// * Code:

// ** Interfaces:

type Connection interface {
	Channel() (Channel, error)

	Close() error

	CloseDeadline(time.Time) error

	ConnectionState() tls.ConnectionState

	IsClosed() bool

	LocalAddr() net.Addr

	NotifyBlocked(chan goamqp.Blocking) chan goamqp.Blocking

	NotifyClose(chan *goamqp.Error) chan *goamqp.Error

	RemoteAddr() net.Addr

	UpdateSecret(string, string) error
}

type Channel interface {
	Ack(tag uint64, multiple bool) error

	Cancel(string, bool) error

	Close() error

	Confirm(bool) error

	Consume(
		string,
		string,
		bool,
		bool,
		bool,
		bool,
		goamqp.Table,
	) (<-chan goamqp.Delivery, error)

	ConsumeWithContext(
		context.Context,
		string,
		string,
		bool,
		bool,
		bool,
		bool,
		goamqp.Table,
	) (<-chan goamqp.Delivery, error)

	ExchangeBind(
		string,
		string,
		string,
		bool,
		goamqp.Table,
	) error

	ExchangeDeclare(
		string,
		string,
		bool,
		bool,
		bool,
		bool,
		goamqp.Table,
	) error

	ExchangeDeclarePassive(
		string,
		string,
		bool,
		bool,
		bool,
		bool,
		goamqp.Table,
	) error

	ExchangeDelete(string, bool, bool) error

	ExchangeUnbind(string, string, string, bool, goamqp.Table) error

	Flow(bool) error

	Get(string, bool) (goamqp.Delivery, bool, error)

	GetNextPublishSeqNo() uint64

	IsClosed() bool

	Nack(uint64, bool, bool) error

	NotifyCancel(chan string) chan string

	NotifyClose(chan *goamqp.Error) chan *goamqp.Error

	NotifyConfirm(chan uint64, chan uint64) (chan uint64, chan uint64)

	NotifyFlow(chan bool) chan bool

	NotifyPublish(chan goamqp.Confirmation) chan goamqp.Confirmation

	NotifyReturn(chan goamqp.Return) chan goamqp.Return

	Publish(
		string,
		string,
		bool,
		bool,
		goamqp.Publishing,
	) error

	PublishWithContext(
		context.Context,
		string,
		string,
		bool,
		bool,
		goamqp.Publishing,
	) error

	PublishWithDeferredConfirm(
		string,
		string,
		bool,
		bool,
		goamqp.Publishing,
	) (*goamqp.DeferredConfirmation, error)

	PublishWithDeferredConfirmWithContext(
		context.Context,
		string,
		string,
		bool,
		bool,
		goamqp.Publishing,
	) (*goamqp.DeferredConfirmation, error)

	Qos(int, int, bool) error

	QueueBind(string, string, string, bool, goamqp.Table) error

	QueueDeclare(
		string,
		bool,
		bool,
		bool,
		bool,
		goamqp.Table,
	) (goamqp.Queue, error)

	QueueDeclarePassive(
		string,
		bool,
		bool,
		bool,
		bool,
		goamqp.Table,
	) (goamqp.Queue, error)

	QueueDelete(string, bool, bool, bool) (int, error)

	QueuePurge(string, bool) (int, error)

	QueueUnbind(string, string, string, goamqp.Table) error

	Reject(uint64, bool) error

	Tx() error

	TxCommit() error

	TxRollback() error
}

// * interface.go ends here.
