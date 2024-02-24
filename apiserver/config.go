/*
 * config.go --- Dispatcher configuration.
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

package apiserver

import (
	"net"
	"strconv"
)

type Config struct {
	Addr    string `json:"address"`
	Cert    string `json:"cert_file"`
	Key     string `json:"key_file"`
	UseTLS  bool   `json:"use_tls"`
	LogFile string `json:"log_file"`

	cachedHost string `config_hide:"true"`
	cachedPort int    `config_hide:"true"`
}

func NewDefaultConfig() *Config {
	return &Config{}
}

func NewConfig(addr, log, cert, key string, tls bool) *Config {
	return &Config{
		Addr:    addr,
		Cert:    cert,
		Key:     key,
		UseTLS:  tls,
		LogFile: log,
	}
}

func (c *Config) Host() (string, error) {
	if c.cachedHost == "" {
		host, port, err := net.SplitHostPort(c.Addr)
		if err != nil {
			return "", err
		}

		c.cachedHost = host
		c.cachedPort, err = strconv.Atoi(port)
		if err != nil {
			return "", err
		}
	}

	return c.cachedHost, nil
}

func (c *Config) Port() (int, error) {
	if c.cachedHost == "" {
		if _, err := c.Host(); err != nil {
			return 0, err
		}
	}

	return c.cachedPort, nil
}

/* config.go ends here. */
