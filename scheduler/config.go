// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config.go --- Scheduler configuration.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

// * Package:

package scheduler

// * Imports:

import (
	"time"

	"github.com/Asmodai/gohacks/health"
	"github.com/prometheus/client_golang/prometheus"
)

// * Constants:

const (
	// Default health ticker period.
	DefaultHealthTickPeriod time.Duration = 1 * time.Second

	// Default channel buffer size.
	//
	// This should be sufficiently big enough to prevent blocking.
	DefaultChannelBufferSize int = 256
)

// * Code:

// ** Types:

// Scheduler configuration.
type Config struct {
	Name             string                // Scheduler name.
	HealthTickPeriod time.Duration         // Health tick period.
	Health           health.Reporter       // Health status reporter.
	Prometheus       prometheus.Registerer // Prometheus registerer.
	AddBuffer        int                   // Size of `add` buffer.
	WorkBuffer       int                   // Size of `work` buffer.
}

// Create a new scheduler configuration instance.
func NewConfig(
	name string,
	tick time.Duration,
	hlth health.Reporter,
	addBuff, workBuff int,
) *Config {
	if tick == 0 {
		tick = DefaultHealthTickPeriod
	}

	if hlth == nil {
		hlth = health.NewDefaultHealth()
	}

	return &Config{
		Name:             name,
		HealthTickPeriod: tick,
		Health:           hlth,
		AddBuffer:        addBuff,
		WorkBuffer:       workBuff,
		Prometheus:       prometheus.DefaultRegisterer,
	}
}

// Create a new default scheduler configuration instance.
func NewDefaultConfig() *Config {
	return NewConfig(
		"default",
		DefaultHealthTickPeriod,
		health.NewDefaultHealth(),
		DefaultChannelBufferSize,
		DefaultChannelBufferSize,
	)
}

// * config.go ends here.
