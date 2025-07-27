// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// db_manager.go --- Database manager context value.
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

	"github.com/Asmodai/gohacks/v1/database"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	ContextKeyDBManager = "_DI_DB_MGR"
)

// * Variables:

var (
	ErrValueNotDBManager = errors.Base("value is not database.Manager")
)

// * Code:

// ** Functions:

// Set the database manager value to the context map.
func SetDBManager(ctx context.Context, inst database.Manager) (context.Context, error) {
	return PutToContext(ctx, ContextKeyDBManager, inst)
}

// Get the database manager from the given context.
//
// Will return `ErrValueNoDBManager` if the value in the context is not of type
// `database.Manager`.
func GetDBManager(ctx context.Context) (database.Manager, error) {
	val, err := GetFromContext(ctx, ContextKeyDBManager)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	inst, ok := val.(database.Manager)
	if !ok {
		return nil, errors.WithStack(ErrValueNotDBManager)
	}

	return inst, nil
}

// Attempt to get the database manager from the given context.  Panics if the
// operation fails.
func MustGetDBManager(ctx context.Context) database.Manager {
	inst, err := GetDBManager(ctx)

	if err != nil {
		panic("Could not get database manager instance from context")
	}

	return inst
}

// * db_manager.go ends here.
