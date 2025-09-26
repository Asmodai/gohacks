// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config.go --- SQL configuration.
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

// * Comments:

// * Package:

package database

// * Imports:

import (
	"fmt"
	"time"

	"github.com/Asmodai/gohacks/types"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	// Default maximum idle connections.
	defaultMaxIdleConns = 20

	// Default maximum open connections.
	defaultMaxOpenConns = 20

	// Default MySQL database port.
	defaultDatabasePort = 3306

	// Default minimum worker count.
	defaultMinWorkerCount = 1

	// Default maximum worker count.
	defaultMaxWorkerCount = 8

	// Default batch timeout.
	defaultBatchTimeout types.Duration = types.Duration(100 * time.Millisecond)

	// Default worker pool idle timeout.
	defaultPoolIdleTimeout types.Duration = types.Duration(30 * time.Second)

	// Default worker pool drain timeout.
	defaultPoolDrainTimeout types.Duration = types.Duration(500 * time.Millisecond)

	// Default batch size.
	defaultBatchSize = 200
)

// * Variables:

var (
	// Signalled when there is no SQL driver.
	ErrNoDriver error = errors.Base("no driver provided")

	// Signalled when there is no SQL username given.
	ErrNoUsername error = errors.Base("no username provided")

	// Signalled when there is no SQL password given.
	ErrNoPassword error = errors.Base("no password provided")

	// Signalled when there is no SQL server hostname given.
	ErrNoHostname error = errors.Base("no hostname provided")

	// Signalled when there is no SQL database name provided.
	ErrNoDatabase error = errors.Base("no database name provided")
)

// * Code:
// ** Type:

/*
SQL configuration structure.
*/
type Config struct {
	Driver           string         `json:"driver"`
	Username         string         `json:"username"`
	Password         string         `config_obscure:"true"     json:"password"`
	Hostname         string         `json:"hostname"`
	Port             int            `json:"port"`
	Database         string         `json:"database"`
	BatchSize        int            `json:"batch_size"`
	BatchTimeout     types.Duration `json:"batch_timeout"`
	SetPoolLimits    bool           `json:"set_pool_limits"`
	MaxIdleConns     int            `json:"max_idle_conns"`
	MaxOpenConns     int            `json:"max_open_conns"`
	UsePool          bool           `json:"use_worker_pool"`
	PoolMinWorkers   int            `json:"pool_min_workers"`
	PoolMaxWorkers   int            `json:"pool_max_workers"`
	PoolIdleTimeout  types.Duration `json:"pool_idle_timeout"`
	PoolDrainTimeout types.Duration `json:"pool_drain_timeout"`

	cachedDsn string `config_hide:"true" json:"-"`
}

// ** Methods:

// Validate the configuration.
//
//nolint:cyclop,funlen
func (c *Config) Validate() []error {
	errs := []error{}

	if c.Driver == "" {
		errs = append(errs, errors.WithStack(ErrNoDriver))
	}

	if c.Username == "" {
		errs = append(errs, errors.WithStack(ErrNoUsername))
	}

	if c.Password == "" {
		errs = append(errs, errors.WithStack(ErrNoPassword))
	}

	if c.Hostname == "" {
		errs = append(errs, errors.WithStack(ErrNoHostname))
	}

	if c.Port <= 0 {
		c.Port = defaultDatabasePort
	}

	if c.Database == "" {
		errs = append(errs, errors.WithStack(ErrNoDatabase))
	}

	if c.BatchSize <= 0 {
		c.BatchSize = defaultBatchSize
	}

	if c.BatchTimeout <= 0 {
		c.BatchTimeout = defaultBatchTimeout
	}

	if c.MaxIdleConns <= 0 {
		c.MaxIdleConns = defaultMaxIdleConns
	}

	if c.MaxOpenConns <= 0 {
		c.MaxOpenConns = defaultMaxOpenConns
	}

	if c.PoolMinWorkers <= 0 {
		c.PoolMinWorkers = defaultMinWorkerCount
	}

	if c.PoolMaxWorkers <= 0 {
		c.PoolMaxWorkers = defaultMaxWorkerCount
	}

	if c.PoolDrainTimeout <= 0 {
		c.PoolDrainTimeout = defaultPoolDrainTimeout
	}

	if c.PoolIdleTimeout <= 0 {
		c.PoolIdleTimeout = defaultPoolIdleTimeout
	}

	if c.UsePool {
		c.SetPoolLimits = true
	}

	return errs
}

// Return the DSN for this database configuration.
func (c *Config) ToDSN() string {
	if c.cachedDsn != "" {
		return c.cachedDsn
	}

	var buf = ""

	if c.Username != "" {
		buf += c.Username

		if c.Password != "" {
			buf += ":" + c.Password
		}

		buf += "@"
	}

	buf += "tcp(" + c.Hostname
	if c.Port > 0 {
		buf += fmt.Sprintf(":%d", c.Port)
	}

	buf += ")/" + c.Database
	buf += "?parseTime=True&loc=UTC&time_zone='-00:00'"

	c.cachedDsn = buf

	return buf
}

// ** Functions:

// Create a new configuration object.
func NewConfig() *Config {
	return &Config{
		Driver:          "mysql",
		Username:        "",
		Password:        "",
		Hostname:        "",
		Port:            0,
		Database:        "",
		cachedDsn:       "",
		BatchSize:       defaultBatchSize,
		SetPoolLimits:   true,
		MaxIdleConns:    defaultMaxIdleConns,
		MaxOpenConns:    defaultMaxOpenConns,
		UsePool:         false,
		PoolMinWorkers:  defaultMinWorkerCount,
		PoolMaxWorkers:  defaultMaxWorkerCount,
		PoolIdleTimeout: defaultPoolIdleTimeout,
	}
}

/* config.go ends here. */
