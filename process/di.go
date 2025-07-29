// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// di.go --- Dependency injection.
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

package process

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/contextdi"
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
func SetManager(ctx context.Context, inst Manager) (context.Context, error) {
	result, err := contextdi.PutToContext(ctx, ContextKeyProcManager, inst)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return result, nil
}

// Get the process manager from the given context.
//
// Will return `ErrValueNotProcessManager` if the value in the context is
// not of type `process.Manager`.
func GetManager(ctx context.Context) (Manager, error) {
	val, err := contextdi.GetFromContext(ctx, ContextKeyProcManager)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	inst, ok := val.(Manager)
	if !ok {
		return nil, errors.WithStack(ErrValueNotProcessManager)
	}

	return inst, nil
}

// Attempt to get the process manager from the given context.  Panics if the
// operation fails.
func MustGetManager(ctx context.Context) Manager {
	inst, err := GetManager(ctx)

	if err != nil {
		panic("Could not get process manager instance from context")
	}

	return inst
}

// * di.go ends here.
