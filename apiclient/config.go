// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config.go --- API client configuration.
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

// * Comments:

//
//
//

// * Package:

package apiclient

// * Imports:

import (
	"time"

	"github.com/Asmodai/gohacks/types"
)

// * Constants:

const (
	defaultRequestsPerSecond int            = 60
	defaultTimeout           types.Duration = types.Duration(30 * time.Second)
)

// * Code:

// ** Types:

// API client configuration.
type Config struct {
	// The number of requests per second should rate limiting be required.
	RequestsPerSecond int `json:"requests_per_second"`

	// HTTP connection timeout value.
	Timeout types.Duration `json:"timeout"`

	// Callback to call in order to check request success.
	SuccessCheck SuccessCheckFn `json:"-"`
}

// ** Functions:

// Create a new API client configuration.
func NewConfig(reqsPerSec int, timeout types.Duration) *Config {
	return &Config{
		RequestsPerSecond: reqsPerSec,
		Timeout:           timeout,
	}
}

// Return a new default API client configuration.
func NewDefaultConfig() *Config {
	return NewConfig(defaultRequestsPerSecond, defaultTimeout)
}

// * config.go ends here.
