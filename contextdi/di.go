// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// di.go --- DI via context user values.
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

	"github.com/Asmodai/gohacks/contextext"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	ErrKeyNotFound = errors.Base("value map key not found")
)

// * Code:

// ** Functions:

// Get a value from a context.
//
// Will signal `contextext.ErrInvalidContext` if the context is not valid.
// Will signal `contextext.ErrValueMapNotFound` if there is no value map.
// Will signal `ErrKeyNotFound` if the value map does not contain the key.
func GetFromContext(ctx context.Context, key string) (any, error) {
	vmap, err := contextext.GetValueMap(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rval, found := vmap.Get(key)
	if !found {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

// Place a value in a context.
//
// If there is no value map in the context then one will be created.
//
// Returns a new context with the value map.
func PutToContext(ctx context.Context, key string, value any) (context.Context, error) {
	var (
		vmap contextext.ValueMap
		err  error
	)

	vmap, err = contextext.GetValueMap(ctx)
	if err != nil {
		vmap = contextext.NewValueMap()
	}

	vmap.Set(key, value)

	return contextext.WithValueMap(ctx, vmap), nil
}

// * di.go ends here.
