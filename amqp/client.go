// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// client.go --- AMQP client.
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
//mock:yes
//go:build amd64 || arm64 || riscv64

// * Comments:

//
//
//

// * Package:

package amqp

// * Imports:

import (
	"context"
	"sync"
	"time"

	"github.com/Asmodai/gohacks/amqp/amqpshim"
	"github.com/Asmodai/gohacks/dynworker"
	"github.com/Asmodai/gohacks/logger"
	"github.com/prometheus/client_golang/prometheus"
	goamqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	ErrNoWorkerPool error = errors.Base("no worker pool available")

	//nolint:gochecknoglobals
	disconnectTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "amqp_disconnect_total",
			Help: "Total number of times client has been disconnected",
		},
		[]string{"consumer"},
	)

	//nolint:gochecknoglobals
	reconnectAttemptTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "amqp_reconnect_attempt_total",
			Help: "Total number of times a reconnect has been attempted",
		},
		[]string{"consumer"},
	)

	//nolint:gochecknoglobals
	consumeTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "amqp_consume_total",
			Help: "Number of consumed messages",
		},
		[]string{"consumer"},
	)

	//nolint:gochecknoglobals
	consumeErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "amqp_consume_error_total",
			Help: "Number of errors during message consumption",
		},
		[]string{"consumer"},
	)

	//nolint:gochecknoglobals
	rejectTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "amqp_reject_total",
			Help: "Number of rejects during message consumption",
		},
		[]string{"consumer"},
	)

	//nolint:gochecknoglobals
	publishTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "amqp_publish_total",
			Help: "Number of messages published",
		},
		[]string{"consumer"},
	)

	//nolint:gochecknoglobals
	publishErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "amqp_publish_error_total",
			Help: "Number of errors during message publishing",
		},
		[]string{"consumer"},
	)

	//nolint:gochecknoglobals
	prometheusInitOnce sync.Once
)

// * Code:

// ** Interface:

type Client interface {
	Connect() error
	IsConnected() bool
	Consume() error
	Publish(goamqp.Publishing) error
	QueueStats() (goamqp.Queue, error)
	GetMessageCount() int64
	Disconnect()
	Close() error
}

// ** Types:

type client struct {
	conn                   amqpshim.Connection // was *goamqp.Connection
	channel                amqpshim.Channel    // was *goamqp.Channel
	lgr                    logger.Logger
	ctx                    context.Context
	pool                   dynworker.WorkerPool
	disconnectMetric       prometheus.Counter
	reconnectAttemptMetric prometheus.Counter
	consumeMetric          prometheus.Counter
	consumeErrorMetric     prometheus.Counter
	rejectMetric           prometheus.Counter
	publishMetric          prometheus.Counter
	publishErrorMetric     prometheus.Counter
	cfg                    *Config
	dialFn                 DialFn
	cancel                 context.CancelFunc
	lastMsgCount           int64
	connMux                sync.Mutex
	chanMux                sync.Mutex
	metricsMux             sync.Mutex
	connected              bool
}

// ** Methods:

func (obj *client) Connect() error {
	obj.connMux.Lock()
	defer obj.connMux.Unlock()

	return obj.connectLocked()
}

// Only call this with `connMux` held.
//
// Failure honor this will introduce  Mr. Foot will meet Mr. Gun!
func (obj *client) connectLocked() error {
	if obj.connected {
		return nil
	}

	var err error

	// Connect to the AMQP server.
	obj.conn, err = obj.dialFn(obj.cfg.URL())
	if err != nil {
		return errors.WithStack(err)
	}

	// Open hailing frequencies, Lt. Uhura!
	obj.channel, err = obj.conn.Channel()
	if err != nil {
		return errors.WithStack(err)
	}

	// Ensure that the queue exists and has the right options.
	_, err = obj.channel.QueueDeclare(
		obj.cfg.QueueName,
		obj.cfg.QueueIsDurable,
		obj.cfg.QueueDeleteWhenUnused,
		obj.cfg.QueueIsExclusive,
		obj.cfg.QueueNoWait,
		nil,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	// Set up Quality-of-Service.
	err = obj.channel.Qos(int(obj.cfg.PrefetchCount), 0, false)
	if err != nil {
		return errors.WithStack(err)
	}

	obj.connected = true

	go obj.monitorConnection()
	go obj.pollQueueStats()

	return nil
}

func (obj *client) IsConnected() bool {
	obj.connMux.Lock()
	defer obj.connMux.Unlock()

	return obj.connected
}

func (obj *client) monitorConnection() {
	obj.connMux.Lock()
	conn := obj.conn
	consumer := obj.cfg.ConsumerName
	obj.connMux.Unlock()

	notifyClose := conn.NotifyClose(make(chan *goamqp.Error))

	select {
	case <-notifyClose:
		obj.lgr.Warn(
			"AMQP connection closed, attempting reconnection...",
			"type", "amqp",
			"consumer", consumer,
		)

		obj.disconnectMetric.Inc()
		obj.reconnectLoop()

	case <-obj.ctx.Done():
		return
	}
}

func (obj *client) reconnectLoop() {
	tries := 0

	for {
		select {
		case <-obj.ctx.Done():
			return

		default:
			time.Sleep(obj.cfg.ReconnectDelay.Duration())

			obj.reconnectAttemptMetric.Inc()

			obj.connMux.Lock()
			err := obj.connectLocked()
			obj.connMux.Unlock()

			if err == nil {
				obj.lgr.Info(
					"AMQP reconnected.",
					"consumer", obj.cfg.ConsumerName,
				)

				return
			}

			if obj.cfg.MaxRetryConnect > 0 && tries > obj.cfg.MaxRetryConnect {
				obj.lgr.Fatal(
					"Maximum AMQP retry attempts reached.",
					"consumer", obj.cfg.ConsumerName,
					"retries", tries,
				)
			}
		}

		tries++
	}
}

func (obj *client) Consume() error {
	if obj.pool == nil {
		return errors.WithStack(ErrNoWorkerPool)
	}

	obj.chanMux.Lock()
	msgs, err := obj.channel.Consume(
		obj.cfg.QueueName,   // Queue name.
		obj.cfg.consumerTag, // Consumer tag.
		false,               // Auto-ack.
		false,               // Exclusive.
		false,               // No-local.
		false,               // No-wait.
		nil,                 // Arguments.
	)
	obj.chanMux.Unlock()

	if err != nil {
		return errors.WithStack(err)
	}

	go func() {
		for {
			select {
			case <-obj.ctx.Done():
				return

			case msg, ok := <-msgs:
				if !ok {
					obj.consumeErrorMetric.Inc()

					return
				}

				obj.consumeMetric.Inc()

				err := obj.pool.Submit(msg)
				if err != nil {
					// We don't care if this fails.
					obj.chanMux.Lock()
					_ = obj.channel.Reject(
						msg.DeliveryTag,
						true,
					)
					obj.chanMux.Unlock()

					obj.rejectMetric.Inc()
				}
			}
		}
	}()

	return nil
}

func (obj *client) Publish(msg goamqp.Publishing) error {
	obj.chanMux.Lock()
	defer obj.chanMux.Unlock()

	obj.publishMetric.Inc()

	err := obj.channel.PublishWithContext(
		obj.ctx,           // Context.
		"",                // Default exchange.
		obj.cfg.QueueName, // Queue name.
		false,             // Is mandatory.
		false,             // Is immediate.
		msg,
	)

	if err != nil {
		obj.publishErrorMetric.Inc()

		err = errors.WithStack(err)
	}

	return err
}

func (obj *client) scaleWorkers() int {
	num := obj.GetMessageCount() / obj.cfg.PrefetchCount
	if (num % obj.cfg.PrefetchCount) != 0 {
		num++
	}

	if num > obj.cfg.MaxWorkers {
		num = obj.cfg.MaxWorkers
	}

	if num < obj.cfg.MinWorkers {
		num = obj.cfg.MinWorkers
	}

	return int(num)
}

func (obj *client) pollQueueStats() {
	ticker := time.NewTicker(obj.cfg.PollInterval.Duration())
	defer ticker.Stop()

	for {
		select {
		case <-obj.ctx.Done():
			return

		case <-ticker.C:
			queue, err := obj.QueueStats()
			if err != nil {
				obj.lgr.Warn(
					"Failed to inspect AMQP queue.",
					"consumer", obj.cfg.ConsumerName,
					"queue", obj.cfg.QueueName,
					"err", err.Error(),
				)

				continue
			}

			// TODO: Scale worker pool as required here?
			obj.metricsMux.Lock()
			obj.lastMsgCount = int64(queue.Messages)
			obj.metricsMux.Unlock()
		}
	}
}

func (obj *client) QueueStats() (goamqp.Queue, error) {
	obj.chanMux.Lock()
	defer obj.chanMux.Unlock()

	queue, err := obj.channel.QueueDeclarePassive(
		obj.cfg.QueueName,
		false, // Passive!
		false, // Don't care.
		false, // Really don't care.
		false, // Yeah, all this needs to be false.
		nil,
	)

	if err != nil {
		return goamqp.Queue{}, errors.WithStack(err)
	}

	return queue, nil
}

func (obj *client) GetMessageCount() int64 {
	obj.metricsMux.Lock()
	defer obj.metricsMux.Unlock()

	return obj.lastMsgCount
}

func (obj *client) Disconnect() {
	obj.lgr.Info(
		"Disconnecting from AMQP.",
		"consumer", obj.cfg.ConsumerName,
	)

	obj.cancel()
}

func (obj *client) Close() error {
	obj.connMux.Lock()
	defer obj.connMux.Unlock()

	var err error

	if obj.channel != nil {
		err = obj.channel.Close()
	}

	if obj.conn != nil {
		err2 := obj.conn.Close()

		if err == nil {
			err = err2
		}
	}

	obj.connected = false

	return errors.WithStack(err)
}

// ** Functions:

func NewClient(ctx context.Context, cfg *Config, pool dynworker.WorkerPool) Client {
	if !cfg.validated {
		panic("AMQP configuration has not been validated.")
	}

	lgr := logger.MustGetLogger(ctx)
	nctx, cancel := context.WithCancel(ctx)

	InitPrometheus(cfg.Prometheus)

	label := prometheus.Labels{"consumer": cfg.ConsumerName}

	inst := &client{
		cfg:    cfg,
		ctx:    nctx,
		cancel: cancel,
		pool:   pool,
		dialFn: cfg.dialer,
		lgr:    lgr,

		// Prometheus metrics.
		disconnectMetric:       disconnectTotal.With(label),
		reconnectAttemptMetric: reconnectAttemptTotal.With(label),
		consumeMetric:          consumeTotal.With(label),
		consumeErrorMetric:     consumeErrorTotal.With(label),
		rejectMetric:           rejectTotal.With(label),
		publishMetric:          publishTotal.With(label),
		publishErrorMetric:     publishErrorTotal.With(label),
	}

	// Set the scaling function.
	inst.pool.SetScalerFunction(inst.scaleWorkers)

	return inst
}

// Initialise Prometheus metrics for this module.
func InitPrometheus(reg prometheus.Registerer) {
	prometheusInitOnce.Do(func() {
		reg.MustRegister(
			disconnectTotal,
			reconnectAttemptTotal,
			consumeTotal,
			consumeErrorTotal,
			rejectTotal,
			publishTotal,
			publishErrorTotal,
		)
	})
}

// * client.go ends here.
