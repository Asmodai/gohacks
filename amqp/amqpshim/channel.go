// -*- Mode: Go -*-
//
// channel.go --- Channel implementation.
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

package amqpshim

// * Imports:

import (
	"context"

	goamqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Types:

type channel struct {
	ch *goamqp.Channel
}

// ** Methods:

func (c *channel) Ack(tag uint64, multiple bool) error {
	err := c.ch.Ack(tag, multiple)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) Cancel(consumer string, noWait bool) error {
	err := c.ch.Cancel(consumer, noWait)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) Close() error {
	err := c.ch.Close()
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) Confirm(noWait bool) error {
	err := c.ch.Confirm(noWait)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) Consume(
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args goamqp.Table,
) (<-chan goamqp.Delivery, error) {
	chnl, err := c.ch.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return chnl, err
}

func (c *channel) ConsumeWithContext(
	ctx context.Context,
	queue, consumer string,
	autoAck, exclusive, noLocal, noWait bool,
	args goamqp.Table,
) (<-chan goamqp.Delivery, error) {
	chnl, err := c.ch.ConsumeWithContext(
		ctx,
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return chnl, err
}

func (c *channel) ExchangeBind(
	destination, key, source string,
	noWait bool,
	args goamqp.Table,
) error {
	err := c.ch.ExchangeBind(destination, key, source, noWait, args)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) ExchangeDeclare(
	name, kind string,
	durable, autoDelete, internal, noWait bool,
	args goamqp.Table,
) error {
	err := c.ch.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) ExchangeDeclarePassive(
	name, kind string,
	durable, autoDelete, internal, noWait bool,
	args goamqp.Table,
) error {
	err := c.ch.ExchangeDeclarePassive(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		noWait,
		args,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) ExchangeDelete(name string, ifUnused, noWait bool) error {
	err := c.ch.ExchangeDelete(name, ifUnused, noWait)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) ExchangeUnbind(
	destination, key, source string,
	noWait bool,
	args goamqp.Table,
) error {
	err := c.ch.ExchangeUnbind(destination, key, source, noWait, args)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) Flow(active bool) error {
	err := c.ch.Flow(active)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) Get(queue string, autoAck bool) (goamqp.Delivery, bool, error) {
	msg, ok, err := c.ch.Get(queue, autoAck)
	if err != nil {
		err = errors.WithStack(err)
	}

	return msg, ok, err
}

func (c *channel) GetNextPublishSeqNo() uint64 {
	return c.ch.GetNextPublishSeqNo()
}

func (c *channel) IsClosed() bool {
	return c.ch.IsClosed()
}

func (c *channel) Nack(tag uint64, multiple bool, requeue bool) error {
	err := c.ch.Nack(tag, multiple, requeue)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) NotifyCancel(cChan chan string) chan string {
	return c.ch.NotifyCancel(cChan)
}

func (c *channel) NotifyClose(cChan chan *goamqp.Error) chan *goamqp.Error {
	return c.ch.NotifyClose(cChan)
}

func (c *channel) NotifyConfirm(ack, nack chan uint64) (chan uint64, chan uint64) {
	return c.ch.NotifyConfirm(ack, nack)
}

func (c *channel) NotifyFlow(cChan chan bool) chan bool {
	return c.ch.NotifyFlow(cChan)
}

func (c *channel) NotifyPublish(confirm chan goamqp.Confirmation) chan goamqp.Confirmation {
	return c.ch.NotifyPublish(confirm)
}

func (c *channel) NotifyReturn(cChan chan goamqp.Return) chan goamqp.Return {
	return c.ch.NotifyReturn(cChan)
}

func (c *channel) Publish(
	exchange, key string,
	mandatory, immediate bool,
	msg goamqp.Publishing,
) error {
	err := c.ch.Publish(
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) PublishWithContext(
	ctx context.Context,
	exchange, key string,
	mandatory, immediate bool,
	msg goamqp.Publishing,
) error {
	err := c.ch.PublishWithContext(
		ctx,
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) PublishWithDeferredConfirm(
	exchange, key string,
	mandatory, immediate bool,
	msg goamqp.Publishing,
) (*goamqp.DeferredConfirmation, error) {
	cnfm, err := c.ch.PublishWithDeferredConfirm(
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return cnfm, err
}

func (c *channel) PublishWithDeferredConfirmWithContext(
	ctx context.Context,
	exchange, key string,
	mandatory, immediate bool,
	msg goamqp.Publishing,
) (*goamqp.DeferredConfirmation, error) {
	cnfm, err := c.ch.PublishWithDeferredConfirmWithContext(
		ctx,
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return cnfm, err
}

func (c *channel) Qos(prefetchCount, prefetchSize int, global bool) error {
	err := c.ch.Qos(prefetchCount, prefetchSize, global)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) QueueBind(
	name, key, exchange string,
	noWait bool,
	args goamqp.Table,
) error {
	err := c.ch.QueueBind(
		name,
		key,
		exchange,
		noWait,
		args,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) QueueDeclare(
	name string,
	durable, autoDelete, exclusive, noWait bool,
	args goamqp.Table,
) (goamqp.Queue, error) {
	queue, err := c.ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return queue, err
}

func (c *channel) QueueDeclarePassive(
	name string,
	durable, autoDelete, exclusive, noWait bool,
	args goamqp.Table,
) (goamqp.Queue, error) {
	queue, err := c.ch.QueueDeclarePassive(
		name,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)

	if err != nil {
		err = errors.WithStack(err)
	}

	return queue, err
}

func (c *channel) QueueDelete(
	name string,
	ifUnused, ifEmpty, noWait bool,
) (int, error) {
	purged, err := c.ch.QueueDelete(name, ifUnused, ifEmpty, noWait)
	if err != nil {
		err = errors.WithStack(err)
	}

	return purged, err
}

func (c *channel) QueuePurge(name string, noWait bool) (int, error) {
	purged, err := c.ch.QueuePurge(name, noWait)
	if err != nil {
		err = errors.WithStack(err)
	}

	return purged, err
}

func (c *channel) QueueUnbind(name, key, exchange string, args goamqp.Table) error {
	err := c.ch.QueueUnbind(name, key, exchange, args)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) Reject(tag uint64, requeue bool) error {
	err := c.ch.Reject(tag, requeue)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) Tx() error {
	err := c.ch.Tx()
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) TxCommit() error {
	err := c.ch.TxCommit()
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *channel) TxRollback() error {
	err := c.ch.TxRollback()
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

// * channel.go ends here.
