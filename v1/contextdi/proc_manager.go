// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// proc_manager.go --- Process manager context value.
//
// Copyright (c) 2023-2025 Paul Ward <paul@lisphacker.uk>
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

package contextdi

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/v1/process"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	ContextKeyProcManager = "_DI_PROC_MGR"
)

// * Variables:

var (
	ErrValueNotProcessManager = errors.Base("value is not process.Manager")
)

// * Code:

// ** Functions:

// Set the process manager value to the context map.
func SetProcessManager(ctx context.Context, inst process.Manager) (context.Context, error) {
	return PutToContext(ctx, ContextKeyProcManager, inst)
}

// Get the process manager from the given context.
//
// Will return `ErrValueNotProcessManager` if the value in the context is
// not of type `process.Manager`.
func GetProcessManager(ctx context.Context) (process.Manager, error) {
	val, err := GetFromContext(ctx, ContextKeyProcManager)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	inst, ok := val.(process.Manager)
	if !ok {
		return nil, errors.WithStack(ErrValueNotProcessManager)
	}

	return inst, nil
}

// Attempt to get the process manager from the given context.  Panics if the
// operation fails.
func MustGetProcessManager(ctx context.Context) process.Manager {
	inst, err := GetProcessManager(ctx)

	if err != nil {
		panic("Could not get process manager instance from context")
	}

	return inst
}

// * proc_manager.go ends here.
