// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// config.go --- Dispatcher configuration.
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

package apiserver

import (
	"gitlab.com/tozd/go/errors"

	"net"
	"strconv"
)

// API server configuration.
type Config struct {
	// Network address that the server will bind to.
	//
	// This is in the format of <address>:<port>.
	// To bind to all available addresses, specify ":<port>" only.
	Addr string `json:"address"`

	// Path to an SSL certificate file if TLS is required.
	Cert string `json:"cert_file"`

	// Path to an SSL key file if TLS is required.
	Key string `json:"key_file"`

	// Should the server use TLS?
	UseTLS bool `json:"use_tls"`

	// Path to the log file for the server.
	LogFile string `json:"log_file"`

	cachedHost string `config_hide:"true"`
	cachedPort int    `config_hide:"true"`
}

// Create a new default configuration.
//
// This will create a configuration that has default values.
func NewDefaultConfig() *Config {
	return &Config{}
}

// Create a new configuration.
func NewConfig(addr, log, cert, key string, tls bool) *Config {
	return &Config{
		Addr:    addr,
		Cert:    cert,
		Key:     key,
		UseTLS:  tls,
		LogFile: log,
	}
}

// Return the hostname on which the server is bound.
func (c *Config) Host() (string, error) {
	if c.cachedHost == "" {
		host, port, err := net.SplitHostPort(c.Addr)
		if err != nil {
			return "", errors.WithStack(err)
		}

		// Cache the host.
		c.cachedHost = host

		// Cache the port.
		c.cachedPort, err = strconv.Atoi(port)
		if err != nil {
			return "", errors.WithStack(err)
		}
	}

	return c.cachedHost, nil
}

// Return the port number on which the server is listening.
func (c *Config) Port() (int, error) {
	if c.cachedHost == "" {
		if _, err := c.Host(); err != nil {
			return 0, errors.WithStack(err)
		}
	}

	return c.cachedPort, nil
}

// config.go ends here.
