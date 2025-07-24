// -*- Mode: Go -*-
//
// config_test.go --- AMQP configuration tests.
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

package amqp

// * Imports:

import (
	"encoding/json"
	"testing"

	"github.com/Asmodai/gohacks/dynworker"
	mlogger "github.com/Asmodai/gohacks/mocks/logger"
	goamqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/mock/gomock"
)

// * Constants:

const (
	TestConfigJSON string = `{
		"url":             "amqp://127.0.0.1",
		"queue_name":      "test_queue",
		"prefetch_count":  9001,
		"poll_interval":   "20ms",
		"reconnect_delay": "20s",
		"consumer_name":   "CoffeeTron"
	}`

	TestHandlerParam string = "this is a test"
	TestScalerParam  int    = 42
)

// * Code:

// ** Tests:

func TestConfig(t *testing.T) {
	var (
		inst *Config

		calledTestHandler      bool = false
		calledTestHandlerValue string
	)

	mocker := gomock.NewController(t)
	defer mocker.Finish()

	lgr := mlogger.NewMockLogger(mocker)

	testHandlerCB := func(param *dynworker.Task) error {
		delivery := param.Data().(goamqp.Delivery)

		calledTestHandler = true
		calledTestHandlerValue = string(delivery.Body)

		return nil
	}

	err := json.Unmarshal([]byte(TestConfigJSON), &inst)
	if err != nil {
		t.Fatalf("JSON error: %#v", err)
	}

	inst.Validate()
	inst.SetMessageHandler(testHandlerCB)
	inst.SetLogger(lgr)

	t.Run("HandleMessage", func(t *testing.T) {
		delivery := goamqp.Delivery{Body: []byte(TestHandlerParam)}
		inst.messageHandler(dynworker.NewTask(nil, nil, delivery))

		if !calledTestHandler {
			t.Error("Test handler was not called!")
			return
		}

		if calledTestHandlerValue != TestHandlerParam {
			t.Errorf(
				"Unexpected value: %v != %v",
				string(calledTestHandlerValue),
				TestHandlerParam,
			)
		}
	})

	t.Run("dynworker config", func(t *testing.T) {
		want := "CoffeeTron"
		cnf := inst.ConfigureWorkerPool()

		if cnf.Name != want {
			t.Errorf("Unexpected Name: %#v != %#v", cnf.Name, want)
		}

		if cnf.MinWorkers != inst.MinWorkers {
			t.Errorf(
				"Unexpected MinWorkers: %#v != %#v",
				cnf.MinWorkers,
				inst.MinWorkers,
			)
		}

		if cnf.MaxWorkers != inst.MaxWorkers {
			t.Errorf(
				"Unexpected MaxWorkers: %#v != %#v",
				cnf.MaxWorkers,
				inst.MaxWorkers,
			)
		}

		if cnf.IdleTimeout != inst.WorkerIdleTimeout.Duration() {
			t.Errorf(
				"Unexpected IdleTimeout: %#v != %#v",
				cnf.IdleTimeout,
				inst.WorkerIdleTimeout.Duration(),
			)
		}

	})
} // TestConfig

func TestConfigCallbacks(t *testing.T) {
	inst := NewDefaultConfig()

	t.Run("defaultHandleMessage", func(t *testing.T) {
		mocker := gomock.NewController(t)
		defer mocker.Finish()

		lgr := mlogger.NewMockLogger(mocker)
		inst.logger = lgr

		lgr.EXPECT().
			Warn(gomock.Any(), gomock.Any()).
			MaxTimes(1).
			MinTimes(1)

		inst.messageHandler(dynworker.NewTask(nil, nil, goamqp.Delivery{}))
	})
} // TestConfigCallbacks

// * config_test.go ends here.
