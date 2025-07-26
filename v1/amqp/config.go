// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// config.go --- AMQP client configuration.
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
//
//go:build amd64 || arm64 || riscv64

// * Comments:
//
//

// * Package:

package amqp

// * Imports:

import (
	"context"
	"os"
	"time"

	"github.com/Asmodai/gohacks/amqp/amqpshim"
	"github.com/Asmodai/gohacks/dynworker"
	"github.com/Asmodai/gohacks/logger"
	"github.com/Asmodai/gohacks/types"
)

// * Constants:

const (
	// Default AMQP consumer tag prefix.
	defaultConsumerTag string = "go-amqp"

	// Default AMQP consumer name.
	//
	// This is composed into the consumer tag on setup and will also be
	// used as a label for all metrics for the given AMQP client.
	defaultConsumerName string = "unknown-consumer"

	// Default AMQP prefetch count.
	//
	// The only sane value for a default prefetch count is 1.  Anything
	// more is violating the principle of least surprise.
	defaultPrefetchCount int64 = 1

	// Default AMQP polling interval.
	defaultPollInterval types.Duration = types.Duration(5 * time.Second)

	// Default AMQP reconnect delay.
	defaultReconnectDelay types.Duration = types.Duration(10 * time.Second)

	// Set this to a non-zero value to limit the retry attempts to this
	// value.
	defaultMaxRetryConnect int = 0

	defaultMinWorkerCount int64 = 1

	defaultMaxWorkerCount int64 = 200

	defaultWorkerIdleTimeout types.Duration = types.Duration(30 * time.Second)
)

// * Code:

// ** Types:

type DialFn func(url string) (amqpshim.Connection, error)

type Config struct {
	URL                   string         `json:"url"`
	QueueName             string         `json:"queue_name"`
	QueueIsDurable        bool           `json:"queue_is_durable"`
	QueueDeleteWhenUnused bool           `json:"queue_delete_when_unused"`
	QueueIsExclusive      bool           `json:"queue_is_exclusive"`
	QueueNoWait           bool           `json:"queue_no_wait"`
	PrefetchCount         int64          `json:"prefetch_count"`
	PollInterval          types.Duration `json:"poll_interval"`
	ReconnectDelay        types.Duration `json:"reconnect_delay"`
	ConsumerName          string         `json:"consumer_name"`
	MaxRetryConnect       int            `json:"max_retry_connect"`
	MaxWorkers            int64          `json:"max_workers"`
	MinWorkers            int64          `json:"min_workers"`
	WorkerIdleTimeout     types.Duration `json:"worker_idle_timeout"`

	dialer         DialFn           `json:"-"`
	parent         context.Context  `json:"-"`
	logger         logger.Logger    `json:"-"`
	consumerTag    string           `json:"-"`
	metricsLabel   string           `json:"-"`
	messageHandler dynworker.TaskFn `json:"-"`
}

// ** Methods:

func (obj *Config) SetMessageHandler(callback dynworker.TaskFn) {
	obj.messageHandler = callback
}

func (obj *Config) SetParent(ctx context.Context) {
	obj.parent = ctx
}

func (obj *Config) SetLogger(lgr logger.Logger) {
	obj.logger = lgr
}

func (obj *Config) SetDialer(dialer DialFn) {
	obj.dialer = dialer
}

func (obj *Config) defaultMessageHandler(_ *dynworker.Task) error {
	obj.logger.Warn(
		"No AMQP message handler callback installed.",
		"consumer", obj.ConsumerName,
	)

	return nil
}

func (obj *Config) makeMetricsLabel() {
	if len(obj.metricsLabel) > 0 {
		return
	}

	name := obj.ConsumerName
	if len(name) == 0 {
		name = defaultConsumerName
	}

	obj.metricsLabel = defaultConsumerTag + "-" + name
}

// Generate an AMQP consumer tag.
func (obj *Config) makeConsumerTag() {
	if len(obj.consumerTag) > 0 {
		return
	}

	name := obj.ConsumerName
	if len(name) == 0 {
		name = defaultConsumerName
	}

	tag := defaultConsumerTag + "-" + name

	if host, err := os.Hostname(); err != nil {
		tag += "--" + host
	}

	obj.consumerTag = tag
}

func (obj *Config) ConfigureWorkerPool() *dynworker.Config {
	return &dynworker.Config{
		Name:        obj.ConsumerName,
		MinWorkers:  obj.MinWorkers,
		MaxWorkers:  obj.MaxWorkers,
		Logger:      obj.logger,
		Parent:      obj.parent,
		IdleTimeout: obj.WorkerIdleTimeout.Duration(),
	}
}

//nolint:cyclop
func (obj *Config) Validate() {
	if len(obj.URL) == 0 {
		obj.logger.Fatal(
			"No AMQP server URL has been provided.",
		)
	}

	if obj.dialer == nil {
		obj.dialer = amqpshim.Dial
	}

	if len(obj.consumerTag) == 0 {
		obj.consumerTag = defaultConsumerTag
	}

	if len(obj.ConsumerName) == 0 {
		obj.ConsumerName = defaultConsumerName
	}

	if obj.PrefetchCount == 0 {
		obj.PrefetchCount = defaultPrefetchCount
	}

	if obj.PollInterval == 0 {
		obj.PollInterval = defaultPollInterval
	}

	if obj.ReconnectDelay == 0 {
		obj.ReconnectDelay = defaultReconnectDelay
	}

	if obj.MaxRetryConnect == 0 {
		obj.MaxRetryConnect = defaultMaxRetryConnect
	}

	if obj.MaxWorkers == 0 {
		obj.MaxWorkers = defaultMaxWorkerCount
	}

	if obj.MinWorkers == 0 {
		obj.MinWorkers = defaultMinWorkerCount
	}

	if obj.WorkerIdleTimeout == 0 {
		obj.WorkerIdleTimeout = defaultWorkerIdleTimeout
	}

	// Set up the default handler.
	obj.messageHandler = obj.defaultMessageHandler

	// Build the metrics label and consumer tag.
	obj.makeMetricsLabel()
	obj.makeConsumerTag()
}

// ** Functions:

// Generate a new default configuration object.
func NewDefaultConfig() *Config {
	return NewConfig(
		context.Background(),
		logger.NewDefaultLogger(),
		"amqp://127.0.0.1",
		"",
	)
}

// Generate a new configuration object.
func NewConfig(
	parent context.Context,
	lgr logger.Logger,
	url, queuename string,
) *Config {
	inst := &Config{
		URL:            url,
		QueueName:      queuename,
		PrefetchCount:  defaultPrefetchCount,
		PollInterval:   defaultPollInterval,
		ReconnectDelay: defaultReconnectDelay,
		parent:         parent,
		logger:         lgr,
	}

	inst.Validate()

	return inst
}

// * config.go ends here.
