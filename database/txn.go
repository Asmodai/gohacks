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
	"time"

	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Types:

type TxnFn func(context.Context, Runner) error

// ** Methods:

// Runs `fn` in a DB transaction and commits/rolls back.
func (obj *database) WithTransaction(ctx context.Context, txfn TxnFn) error {
	const (
		maxRetries = 3
		backoff    = 100 * time.Millisecond
	)

	for attempt := 0; ; attempt++ {
		// one attempt = one scoped function, so defer is not stacked
		attemptErr := func() error {
			txn, err := obj.real.BeginTxx(ctx, nil)
			if err != nil {
				return errors.WithStack(err)
			}

			defer func() {
				if p := recover(); p != nil {
					_ = txn.Rollback()

					panic(p)
				}
			}()

			// Execute the transaction function.
			if err := txfn(ctx, txn); err != nil {
				_ = txn.Rollback()

				return obj.GetError(err)
			}

			if err := txn.Commit(); err != nil {
				_ = txn.Rollback()

				return obj.GetError(err)
			}

			return nil
		}()

		if attemptErr == nil {
			return nil
		}

		inerror := (errors.Is(attemptErr, ErrTxnDeadlock) ||
			errors.Is(attemptErr, ErrTxnSerialization)) &&
			attempt < maxRetries

		if inerror {
			// Back off and retry.
			time.Sleep(backoff * time.Duration(attempt+1))

			continue
		}

		return attemptErr
	}
}

// * txn.go ends here.
