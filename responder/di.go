// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// di.go --- Dependency injection.
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

package responder

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/contextdi"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	ContextKeyResponderChain = "_DI_RESPONDER"
)

// * Variables:

var (
	ErrValueNotResponderChain = errors.Base("value is not responder.Chain")
)

// * Code:

// ** Functions:

// Set the responder chain value in the context map.
func SetResponderChain(ctx context.Context, inst *Chain) (context.Context, error) {
	val, err := contextdi.PutToContext(ctx, ContextKeyResponderChain, inst)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return val, nil
}

// Get the responder chain from the given context.
//
// Will return `ErrValueNotResponderChain` if the value in the context is not
// of type `responder.Chain`.
//
// Please be aware that this responder chain should be treated as immutable,
// as we can't really propagate changes down the context hierarchy.
func GetResponderChain(ctx context.Context) (*Chain, error) {
	val, err := contextdi.GetFromContext(ctx, ContextKeyResponderChain)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	inst, ok := val.(*Chain)
	if !ok {
		return nil, errors.WithStack(ErrValueNotResponderChain)
	}

	return inst, nil
}

// Attempt to get the responder chain from the given context.  Panics if the
// operation fails.
func MustGetResponderChain(ctx context.Context) *Chain {
	inst, err := GetResponderChain(ctx)
	if err != nil {
		panic("Could not get responder chain instance from context")
	}

	return inst
}

// * di.go ends here.
