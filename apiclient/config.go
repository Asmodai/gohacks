// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// config.go --- API client configuration.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package apiclient

// API client configuration.
type Config struct {
	// The number of requests per second should rate limiting be required.
	RequestsPerSecond int `json:"requests_per_second"`

	// HTTP connection timeout value.
	Timeout int `json:"timeout"`
}

// Create a new API client configuration.
func NewConfig(reqsPerSec, timeout int) *Config {
	return &Config{
		RequestsPerSecond: reqsPerSec,
		Timeout:           timeout,
	}
}

// Return a new default API client configuration.
//
//nolint:mnd
func NewDefaultConfig() *Config {
	return NewConfig(5, 5)
}

// config.go ends here.
