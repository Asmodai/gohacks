// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// client_test.go --- AMQP client tests.
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

// * Comments:
//
//

// * Package:

package amqp

// * Imports:

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Asmodai/gohacks/amqp/amqpmock"
	"github.com/Asmodai/gohacks/logger"
	mamqpshim "github.com/Asmodai/gohacks/mocks/amqpshim"
	mdynworker "github.com/Asmodai/gohacks/mocks/dynworker"
	mlogger "github.com/Asmodai/gohacks/mocks/logger"
	goamqp "github.com/rabbitmq/amqp091-go"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/mock/gomock"
)

// * Constants:

const (
	ConfigJSON = `{
		"url":                 "amqp://127.0.0.1",
		"queue_name":          "/test",
		"prefetch_count":      10,
		"poll_interval":       "2s",
		"reconnect_delay":     "5s",
		"consumer_name":       "test_consumer",
		"max_retry_connect":   10,
		"max_workers":         10,
		"min_workers":         10,
		"worker_idle_timeout": "5s"
	}`
)

// ** Variables

var (
	ErrTestConnect error = errors.Base("connected string has broken")
	ErrTestError   error = errors.Base("oh no")
)

// * Code:

// ** Tests:

func TestClient(t *testing.T) {
	var (
		inst Client
		cnf  *Config = &Config{}
	)

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	dialer := &amqpmock.MockDialer{}
	mockconn := &amqpmock.MockConnection{}
	mockchan := &amqpmock.MockChannel{}

	mockconn.Init()
	mockchan.Init()

	mocker := gomock.NewController(t)
	defer mocker.Finish()

	amqpChan := mamqpshim.NewMockChannel(mocker)
	amqpConn := mamqpshim.NewMockConnection(mocker)

	mockchan.MockClose(amqpChan).AnyTimes()
	mockchan.MockQueueDeclare(amqpChan).AnyTimes()
	mockchan.MockQos(amqpChan).AnyTimes()

	mockconn.BuildChannel(amqpmock.ChannelResults{
		Channel: amqpChan,
		Error:   nil,
	})

	mockconn.MockChannel(amqpConn).AnyTimes()
	mockconn.MockNotifyClose(amqpConn).AnyTimes()
	mockconn.MockClose(amqpConn).AnyTimes()

	pool := mdynworker.NewMockWorkerPool(mocker)
	pool.EXPECT().SetScalerFunction(gomock.Any()).AnyTimes()

	lgr := mlogger.NewMockLogger(mocker)
	lgr.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	dictx, err := logger.SetLogger(ctx, lgr)
	if err != nil {
		t.Fatalf("Could not set DI logger: %#v", err)
	}

	err = json.Unmarshal([]byte(ConfigJSON), &cnf)
	if err != nil {
		t.Fatalf("Unexpected error: %#v", err)
	}

	cnf.Hostname = "127.0.0.1"
	cnf.VirtualHost = "/"
	cnf.SetDialer(dialer.Dial)

	errs := cnf.Validate()
	if len(errs) > 0 {
		t.Error("Errors from Validate():")
		for i, e := range errs {
			t.Errorf("  %d: %v", i+1, e.Error())
		}
		t.Fatal("Cannot continue.")
	}

	t.Run("Constructs", func(t *testing.T) {
		inst = NewClient(dictx, cnf, pool)

		if inst == nil {
			t.Fatal("Construction failed!")
		}
	})

	t.Run("Connect", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			dialer.Connection = amqpConn
			dialer.Error = nil

			err := inst.Connect()
			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
			}

			// give `monitorConnection` time to spin up.
			time.Sleep(500 * time.Millisecond)

			if !inst.IsConnected() {
				t.Error("AMQP client reports not connected!")
			}

			inst.Disconnect()
			inst.Close()
		})

		t.Run("Failure", func(t *testing.T) {
			dialer.Connection = nil
			dialer.Error = ErrTestConnect

			err := inst.Connect()
			if !errors.Is(err, ErrTestConnect) {
				t.Errorf("Unexpected error: %#v", err)
			}
		})

		t.Run("AMQP failures", func(t *testing.T) {
			dialer.Connection = amqpConn
			dialer.Error = nil

			t.Run("Channel", func(t *testing.T) {
				mockconn.BuildChannel(
					amqpmock.ChannelResults{
						Channel: nil,
						Error:   ErrTestError,
					},
				)

				err := inst.Connect()

				if err == nil {
					t.Fatal("Expecting an error")
				}

				if !errors.Is(err, ErrTestError) {
					t.Fatalf("Unexpected error: %#v", err)
				}
			})

			t.Run("QueueDeclare", func(t *testing.T) {
				mockconn.BuildChannel(amqpmock.ChannelResults{
					Channel: amqpChan,
					Error:   nil,
				})

				mockchan.BuildQueueDeclare(
					amqpmock.QueueDeclareResults{
						Queue: goamqp.Queue{},
						Error: ErrTestError,
					},
				)

				err := inst.Connect()

				if err == nil {
					t.Fatal("Expecting an error")
				}

				if !errors.Is(err, ErrTestError) {
					t.Fatalf("Unexpected error: %#v", err)
				}
			})

			t.Run("Qos", func(t *testing.T) {
				mockchan.Init()

				mockchan.BuildQos(amqpmock.ErrorResults{
					Error: ErrTestError,
				})

				err := inst.Connect()

				if err == nil {
					t.Fatal("Expecting an error")
				}

				if !errors.Is(err, ErrTestError) {
					t.Fatalf("Unexpected error: %#v", err)
				}
			})
		})
	})

	t.Logf("Connection call log:\n%s", mockconn.CallLog.Dump())
	t.Logf("Channel call log:\n%s", mockchan.CallLog.Dump())
}

// * client_test.go ends here.
