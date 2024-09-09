// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// tx.go --- Transactions.
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
//
// mock:yes

package database

import (
	"github.com/jmoiron/sqlx"
	"gitlab.com/tozd/go/errors"

	"database/sql"
)

type Tx interface {
	NamedExec(string, any) (sql.Result, error)
	Commit() error
}

type tx struct {
	real *sqlx.Tx
}

func NewTx() Tx {
	return &tx{}
}

func (obj *tx) NamedExec(query string, arg any) (sql.Result, error) {
	rval, err := obj.real.NamedExec(query, arg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (obj *tx) Commit() error {
	return errors.WithStack(obj.real.Commit())
}

// tx.go ends here.
