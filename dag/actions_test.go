// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// actions_test.go --- Actions tests.
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

// * Package:

package dag

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/logger"
	mlogger "github.com/Asmodai/gohacks/mocks/logger"
	"gitlab.com/tozd/go/errors"
	"go.uber.org/mock/gomock"

	"testing"
)

// * Constants

const (
	MutateAttr   string = "test"
	MutateOldVal string = "old value"
	MutateNewVal string = "new value"
	LogMessage   string = "Out of coffee"
)

// * Code:

func TestActions(t *testing.T) {
	var fn func(context.Context, Filterable)
	acts := &actions{}

	mocker := gomock.NewController(t)

	lgr := mlogger.NewMockLogger(mocker)
	lgr.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()

	parent := context.TODO()
	ctx, err := logger.SetLogger(parent, lgr)
	if err != nil {
		t.Fatalf("Could not set logger DI: %#v", err)
	}

	// Test `mutate` builtin.
	t.Run("Mutate", func(t *testing.T) {
		t.Run("No params", func(t *testing.T) {
			_, err := acts.Builder("mutate", nil)

			if !errors.Is(err, ErrExpectedParams) {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("No attribute", func(t *testing.T) {
			_, err := acts.Builder("mutate", ActionParams{})

			if !errors.Is(err, ErrMissingParam) {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("Attribute not string", func(t *testing.T) {
			_, err := acts.Builder(
				"mutate",
				ActionParams{
					"attribute": 42,
				},
			)

			if !errors.Is(err, ErrExpectedString) {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("No new_value", func(t *testing.T) {
			_, err := acts.Builder(
				"mutate",
				ActionParams{
					"attribute": "test",
				},
			)

			if !errors.Is(err, ErrMissingParam) {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("Valid", func(t *testing.T) {
			var err error

			fn, err = acts.Builder(
				"mutate",
				ActionParams{
					"attribute": MutateAttr,
					"new_value": MutateNewVal,
				},
			)

			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("Executes", func(t *testing.T) {
			data := NewDataInput()
			data.fields[MutateAttr] = MutateOldVal

			fn(ctx, data)

			value, ok := data.Get(MutateAttr)
			if !ok {
				t.Fatalf("Could not find %q in map.  %#v",
					MutateAttr,
					data)
			}

			if value != MutateNewVal {
				t.Errorf(
					"Unexpected value: %#v != %#v",
					value,
					MutateNewVal,
				)

				return
			}
		})
	}) // mutate

	// Test `log` builtin.
	t.Run("Log", func(t *testing.T) {
		t.Run("No params", func(t *testing.T) {
			_, err := acts.Builder("log", nil)

			if !errors.Is(err, ErrExpectedParams) {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("No message", func(t *testing.T) {
			_, err := acts.Builder("log", ActionParams{})

			if !errors.Is(err, ErrMissingParam) {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("Message not string", func(t *testing.T) {
			_, err := acts.Builder(
				"log",
				ActionParams{
					"message": 42,
				},
			)

			if !errors.Is(err, ErrExpectedString) {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("Valid", func(t *testing.T) {
			var err error

			fn, err = acts.Builder(
				"Log",
				ActionParams{
					"message": LogMessage,
				},
			)

			if err != nil {
				t.Errorf("Unexpected error: %#v", err)
				return
			}
		})

		t.Run("Executes", func(t *testing.T) {
			data := NewDataInput()

			lgr.EXPECT().
				Info(gomock.Any(), gomock.Any()).
				MinTimes(1).
				MaxTimes(1)

			fn(ctx, data)
		})
	}) // log

}

// * actions_test.go ends here.
