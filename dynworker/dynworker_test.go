// -*- Mode: Go -*-
//
// dynworker_test.go --- Dynamic worker tests.
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
	"github.com/Asmodai/gohacks/mocks/logger"

	"go.uber.org/mock/gomock"

	"context"
	"testing"
	"time"
)

// * Code:

// ** Tests:

func TestDynworker(t *testing.T) {
	mocker := gomock.NewController(t)
	defer mocker.Finish()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	lgr := logger.NewMockLogger(mocker)
	lgr.EXPECT().
		Info(gomock.Any(), gomock.Any()).
		AnyTimes()

	cfg := NewConfig(ctx, lgr, "test", 2, 5)
	pool := NewWorkerPool(cfg, func(task Task) error {
		time.Sleep(50 * time.Millisecond)
		return nil
	})

	pool.Start()
	defer pool.Stop()

	t.Run("Submit tasks", func(t *testing.T) {
		for range 10 {
			if err := pool.Submit(struct{}{}); err != nil {
				t.Fatalf("Submit failed: %v", err)
			}
		}
	})

	t.Run("Worker count", func(t *testing.T) {
		t.Log("Sleeping 200ms")
		time.Sleep(200 * time.Millisecond)

		current := pool.WorkerCount()
		if current < 2 || current > 5 {
			t.Errorf("Expected between 2 and 5 workers, got %d", current)

			return
		}
	})

	t.Run("Scale down", func(t *testing.T) {
		t.Log("Sleeping 35 seconds")
		time.Sleep(35 * time.Second)

		current := pool.WorkerCount()
		if current != pool.MinWorkers() {
			t.Errorf(
				"Expected workers to scale to %d, but got %d",
				pool.MinWorkers(),
				current,
			)

			return
		}
	})
}

// * dynworker_test.go ends here.
