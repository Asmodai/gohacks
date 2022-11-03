/*
 * config.go --- Default application config.
 *
 * Copyright (c) 2022 Paul Ward <asmodai@gmail.com>
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

package app

import (
	"github.com/Asmodai/gohacks/config"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/process"
	"github.com/Asmodai/gohacks/semver"
)

type Config struct {
	Name           string
	Version        *semver.SemVer
	Logger         logger.ILogger
	ProcessManager process.IManager
	AppConfig      any
	Validators     config.ValidatorsMap

	//nolint:staticcheck
	options struct {
		enableRPC bool
	}
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) validate() {
	if c.Name == "" {
		c.Name = "Unnamed Application"
	}

	if c.AppConfig == nil {
		c.AppConfig = &defaultAppConfig{}
	}

	if c.Logger == nil {
		c.Logger = logger.NewDefaultLogger()
	}
}

// Empty struct that gets passed as an application configuration when
// no actual configuration is passed to the constructor.
type defaultAppConfig struct {
}

/* config.go ends here. */
