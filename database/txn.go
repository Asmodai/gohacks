// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// txn.go --- Transaction hackery.
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
//

// * Package:

package database

// * Imports:

import (
	"context"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

// * Variables:

// * Code:

// ** Types:

type TxnFn func(ctx context.Context) error

func WithTransaction(ctx context.Context, dbase Database, callback TxnFn) error {
	var retErr error

	txCtx, err := dbase.Begin(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if prec := recover(); prec != nil {
			// If a panic occurred, attempt a rollback and then
			// re-panic like a boss.
			if rbErr := dbase.Rollback(txCtx); rbErr != nil {
				prec = errors.WithStack(
					errors.WithMessagef(rbErr, "%v", prec),
				)
			}

			panic(prec)
		}

		if retErr != nil {
			if rbErr := dbase.Rollback(txCtx); rbErr != nil {
				retErr = errors.WithStack(errors.WrapWith(retErr, rbErr))
			}
		} else {
			if cErr := dbase.Commit(txCtx); cErr != nil {
				retErr = errors.WithStack(cErr)
			}
		}
	}()

	// Invoke the callback making sure to capture its error.
	retErr = errors.WithStack(callback(txCtx))

	return retErr
}

// * txn.go ends here.
