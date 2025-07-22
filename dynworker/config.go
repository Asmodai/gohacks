// -*- Mode: Go -*-
//
// config.go --- Dynamic worker configuration.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
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

// * Package:

package dynworker

// * Imports:

import (
	"github.com/Asmodai/gohacks/logger"

	"context"
	"time"
)

// * Constants:

const (
	// Default minimum worker count.
	DefaultMinimumWorkerCount int32 = 1

	// Default maximum worker count.
	DefaultMaximumWorkerCount int32 = 10

	// Default worker count multipler.
	//
	// This is used when there is an invalid maximum worker count.
	DefaultWorkerCountMult int32 = 4

	// Default worker timeout.
	DefaultTimeout time.Duration = 30 * time.Second
)

// * Code:

// ** Types:

type Config struct {
	Name        string          // Worker pool name for logger and metrics.
	MinWorkers  int32           // Minimum number of workers.
	MaxWorkers  int32           // Maximum number of workers.
	Logger      logger.Logger   // Logger instance.
	Parent      context.Context // Parent context.
	IdleTimeout time.Duration   // Idle timeout duration.
}

// ** Methods:

// Set the idle timeout value.
func (obj *Config) SetItleTimeout(timeout time.Duration) {
	obj.IdleTimeout = timeout
}

// ** Functions:

// Create a new default configuration.
func NewDefaultConfig() *Config {
	return NewConfig(
		context.Background(),
		logger.NewDefaultLogger(),
		"default",
		DefaultMinimumWorkerCount,
		DefaultMaximumWorkerCount,
	)
}

// Create a new configuration.
func NewConfig(
	ctx context.Context,
	lgr logger.Logger,
	name string,
	minw, maxw int32,
) *Config {
	if minw < 0 {
		minw = DefaultMinimumWorkerCount
	}

	if maxw <= 0 {
		maxw = DefaultMaximumWorkerCount
	}

	if maxw < minw {
		maxw = minw * DefaultWorkerCountMult
	}

	return &Config{
		Name:        name,
		MinWorkers:  minw,
		MaxWorkers:  maxw,
		Logger:      lgr,
		Parent:      ctx,
		IdleTimeout: DefaultTimeout,
	}
}

// * config.go ends here.
