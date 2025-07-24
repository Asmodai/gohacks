// -*- Mode: Go -*-
//
// connection.go --- Connection implementation.
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
	"crypto/tls"
	"net"
	"time"

	goamqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Types:

type connection struct {
	conn *goamqp.Connection
}

// ** Methods:

func (c *connection) Channel() (Channel, error) {
	chnl, err := c.conn.Channel()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &channel{chnl}, nil
}

func (c *connection) Close() error {
	err := c.conn.Close()
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *connection) CloseDeadline(deadline time.Time) error {
	err := c.conn.CloseDeadline(deadline)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

func (c *connection) ConnectionState() tls.ConnectionState {
	return c.conn.ConnectionState()
}

func (c *connection) IsClosed() bool {
	return c.conn.IsClosed()
}

func (c *connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *connection) NotifyBlocked(rec chan goamqp.Blocking) chan goamqp.Blocking {
	return c.conn.NotifyBlocked(rec)
}

func (c *connection) NotifyClose(rec chan *goamqp.Error) chan *goamqp.Error {
	return c.conn.NotifyClose(rec)
}

func (c *connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *connection) UpdateSecret(secret, reason string) error {
	err := c.conn.UpdateSecret(secret, reason)
	if err != nil {
		err = errors.WithStack(err)
	}

	return err
}

// ** Functions:

func Dial(url string) (Connection, error) {
	conn, err := goamqp.Dial(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &connection{conn}, nil
}

func DialTLS(url string, amqps *tls.Config) (Connection, error) {
	conn, err := goamqp.DialTLS(url, amqps)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &connection{conn}, nil
}

// * connection.go ends here.
