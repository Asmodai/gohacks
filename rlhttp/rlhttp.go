// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// rlhttp.go --- Rate-limited HTTP client.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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
//go:generate go run github.com/Asmodai/gohacks/cmd/digen -pattern .
//di:gen basename=RLHTTP key=gohacks/rlhttp@v1 type=*Client fallback=NewDefault()

// * Comments:

// * Package:

package rlhttp

// * Imports:

import (
	"context"
	"net/http"
	"time"

	"github.com/Asmodai/gohacks/errx"
	"golang.org/x/time/rate"
)

// * Variables:

var (
	ErrNilRequest = errx.Base("nil request")
)

// * Code:

// ** Types:

type Client struct {
	client  *http.Client
	limiter *rate.Limiter
	timeout time.Duration
}

// ** Methods:

// Perform a HTTP request.
//
// Please note that as this is essentially meant to be used as a middle man
// between `apiclient.Client` and `http.Client`, we do not need to drain
// and close the body here.  That is firmly the responsibility of whatever
// consumes `rlhttp.Client`.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	if req == nil {
		return nil, errx.WithMessage(ErrNilRequest, "nil *http.Request")
	}

	// Make the timeout cover both the limiter wait and request
	// execution.
	if c.timeout > 0 {
		if dl, ok := ctx.Deadline(); !ok || time.Until(dl) > c.timeout {
			var cancel context.CancelFunc

			ctx, cancel = context.WithTimeout(ctx, c.timeout)
			defer cancel()

			req = req.WithContext(ctx)
		}
	}

	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, errx.WithStack(err)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errx.WithStack(err)
	}

	return resp, nil
}

// ** Functions:

func NewDefault() *Client {
	return NewClient(NewDefaultConfig())
}

// Create a new rate-limited HTTP client instance.
func NewClient(cnf *Config) *Client {
	var limiter *rate.Limiter

	// Harsh, but need to drive home the point that mis-configured
	// rate limiting is evil.
	if cnf == nil {
		panic("No rate limiting configuration given")
	}

	if cnf.Enabled {
		// This should be done upstream, but let's do it again for
		// the same reason.
		if errs := cnf.Validate(); len(errs) > 0 {
			panic(errs)
		}

		perReq := rate.Every(cnf.Every.Duration() / time.Duration(cnf.Max))
		if perReq <= 0 {
			panic("rate limit interval <= 0")
		}

		limiter = rate.NewLimiter(
			perReq,
			cnf.Burst,
		)
	}

	// Safety net here.
	client := &http.Client{Timeout: cnf.Timeout.Duration()}

	return &Client{
		client:  client,
		limiter: limiter,
		timeout: cnf.Timeout.Duration(),
	}
}

// * rlhttp.go ends here.
