// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config.go --- Dynamic worker configuration.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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
//go:build amd64 || arm64 || riscv64

// * Comments:
//
//

// * Package:

package dynworker

// * Imports:

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// * Constants:

const (
	// Default minimum worker count.
	defaultMinimumWorkerCount int64 = 1

	// Default maximum worker count.
	defaultMaximumWorkerCount int64 = 10

	// Default worker count multipler.
	//
	// This is used when there is an invalid maximum worker count.
	defaultWorkerCountMult int64 = 4

	// Default worker timeout.
	defaultTimeout time.Duration = 30 * time.Second

	// Default drain target.
	defaultDrainTarget time.Duration = 500 * time.Millisecond
)

// * Code:

// ** Types:

type Config struct {

	// Custom queue to use.
	InputQueue TaskQueue

	// Prometheus registerer.
	Prometheus prometheus.Registerer

	// Function to use as the worker.
	WorkerFunc TaskFn

	// Function to use to determine scaling.
	ScalerFunc ScalerFn

	// Worker pool name for logger and metrics.
	Name string

	// Minimum number of workers.
	MinWorkers int64

	// Maximum number of workers.
	MaxWorkers int64

	// Idle timeout duration.
	IdleTimeout time.Duration

	// Drain target duration.
	DrainTarget time.Duration
}

// ** Methods:

// Set the idle timeout value.
func (obj *Config) SetIdleTimeout(timeout time.Duration) {
	obj.IdleTimeout = timeout
}

// Set the worker function.
func (obj *Config) SetWorkerFunction(workfn TaskFn) {
	obj.WorkerFunc = workfn
}

// Set the scaler function.
func (obj *Config) SetScalerFunction(scalefn ScalerFn) {
	obj.ScalerFunc = scalefn
}

// ** Functions:

// Create a new default configuration.
func NewDefaultConfig() *Config {
	return NewConfig(
		"default",
		defaultMinimumWorkerCount,
		defaultMaximumWorkerCount,
	)
}

func NewConfig(name string, minw, maxw int64) *Config {
	return NewConfigWithQueue(name, minw, maxw, nil)
}

// Create a new configuration with a custom queue.
func NewConfigWithQueue(name string, minw, maxw int64, queue TaskQueue) *Config {
	var inputQueue TaskQueue

	if queue != nil {
		inputQueue = queue
	} else {
		inputQueue = NewChanTaskQueue(int(defaultWorkerChannels))
	}

	if minw < 0 {
		minw = defaultMinimumWorkerCount
	}

	if maxw <= 0 {
		maxw = defaultMaximumWorkerCount
	}

	if maxw < minw {
		maxw = minw * defaultWorkerCountMult
	}

	return &Config{
		Name:        name,
		MinWorkers:  minw,
		MaxWorkers:  maxw,
		IdleTimeout: defaultTimeout,
		InputQueue:  inputQueue,
		Prometheus:  prometheus.DefaultRegisterer,
		DrainTarget: defaultDrainTarget,
	}
}

// * config.go ends here.
