// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config.go --- Default application config.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

package app

import (
	"github.com/Asmodai/gohacks/v1/config"
	"github.com/Asmodai/gohacks/v1/logger"
	"github.com/Asmodai/gohacks/v1/process"
	"github.com/Asmodai/gohacks/v1/semver"
)

// Application configuration.
type Config struct {
	// The application's pretty name.
	Name string

	// The application's version number.
	Version *semver.SemVer

	// The application's logger instance.
	Logger logger.Logger

	// The application's process manager instance.
	ProcessManager process.Manager

	// The application's configuration.
	AppConfig any

	// Validators used to validate the application's configuration.
	Validators config.ValidatorsMap
}

// Create a new empty configuration.
func NewConfig() *Config {
	return &Config{}
}

// Validate the configuration.
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

// config.go ends here.
