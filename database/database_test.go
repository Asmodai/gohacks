// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// database_test.go --- Database driver tests.
//
// Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>
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

package database

import (
	ctxvalmap "github.com/Asmodai/gohacks/context"
	"gitlab.com/tozd/go/errors"

	"context"
	"testing"
)

func TestDatabaseOpen(t *testing.T) {
	t.Run("Returns error if cannot connect", func(t *testing.T) {
		_, e1 := Open("nil", "nil")
		if e1 == nil {
			t.Error("Expected an error condition")
		}
	})
}

func TestContext(t *testing.T) {
	var nctx context.Context

	db1 := &database{}
	db2 := &database{}
	ctx := context.TODO()

	t.Run("Errors when no value map", func(t *testing.T) {
		_, err := FromContext(ctx, "testingdb")
		if !errors.Is(err, ctxvalmap.ErrValueMapNotFound) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("Writes a value map", func(t *testing.T) {
		nctx, _ = ToContext(ctx, db1, "testingdb")
	})

	t.Run("Errors when key not found", func(t *testing.T) {
		_, err := FromContext(nctx, "derp")
		if !errors.Is(err, ErrNoContextKey) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("Reads from value maps", func(t *testing.T) {
		val, err := FromContext(nctx, "testingdb")
		if err != nil {
			t.Errorf("Unexpected error: %v", err.Error())
		}

		if val == nil {
			t.Error("Got NIL back")
		}

		if val != db1 {
			t.Error("Did not get the right instance")
		}
	})

	t.Run("Modifies value maps", func(t *testing.T) {
		nctx, _ = ToContext(nctx, db2, "testingdb")

		val, err := FromContext(nctx, "testingdb")
		if err != nil {
			t.Errorf("Unexpected error: %v", err.Error())
		}

		if val == nil {
			t.Error("Got NIL back")
		}

		if val != db2 {
			t.Error("Did not get the right instance")
		}
	})

	t.Run("Errors when a key is not a database", func(t *testing.T) {
		vmap, _ := ctxvalmap.GetValueMap(nctx)
		vmap.Set("notadb", 42)
		nctx = ctxvalmap.WithValueMap(ctx, vmap)

		_, err := FromContext(nctx, "notadb")
		if !errors.Is(err, ErrValueIsNotDatabase) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})
}

// database_test.go ends here.
