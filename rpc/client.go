/*
 * client.go --- RPC client.
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

	"errors"
	gorpc "net/rpc"
)

type Client struct {
	conf   *Config
	client *gorpc.Client
	logger logger.Logger
}

func NewClient(cnf *Config, lgr logger.Logger) *Client {
	return &Client{
		conf:   cnf,
		client: nil,
		logger: lgr,
	}
}

func (c *Client) Dial() error {
	clnt, err := gorpc.DialHTTP(c.conf.Proto, c.conf.Addr)
	if err != nil {
		c.client = nil

		return err
	}

	c.client = clnt

	return nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Call(method string, args any, reply any) error {
	if c.client == nil {
		c.logger.Warn(
			"Attempt to call RPC endpoint without client connection.",
		)

		return errors.New("No current RPC client connection.")
	}

	if err := c.client.Call(method, args, reply); err != nil {
		return err
	}

	return nil
}

/* client.go ends here. */
