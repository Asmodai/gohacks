// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// database.go --- Mockable SQL interface.
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
	"github.com/Asmodai/gohacks/contextext"

	// This is the MySQL driver, it must be blank.
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/tozd/go/errors"

	"context"
	"database/sql"
)

var (
	errKeyNotFound error = errors.Base("key not found")

	ErrTxnKeyNotFound error = errors.Base("transaction key not found")
	ErrTxnKeyNotTxn   error = errors.Base("key value is not a transaction")
	ErrTxnContext     error = errors.Base("could not create transaction context")
	ErrTxnStart       error = errors.Base("could not start transaction")
)

const (
	KeyTransaction string = "_DB_TXN"
)

type Database interface {
	Ping() error
	Close() error

	SetMaxIdleConns(int)
	SetMaxOpenConns(int)

	Tx(context.Context) (*sqlx.Tx, error)
	Begin(context.Context) (context.Context, error)
	Commit(context.Context) error
	Rollback(context.Context) error
}

type database struct {
	real   *sqlx.DB
	driver string
}

func (obj *database) Ping() error {
	return errors.WithStack(obj.real.Ping())
}

func (obj *database) Close() error {
	return errors.WithStack(obj.real.Close())
}

func (obj *database) SetMaxIdleConns(limit int) {
	obj.real.SetMaxIdleConns(limit)
}

func (obj *database) SetMaxOpenConns(limit int) {
	obj.real.SetMaxOpenConns(limit)
}

func (obj *database) Tx(ctx context.Context) (*sqlx.Tx, error) {
	return getTx(ctx)
}

func (obj *database) Begin(ctx context.Context) (context.Context, error) {
	txn, err := getTx(ctx)

	// Did we fail at getting a transaction from the context?
	if err != nil {
		if !errors.Is(err, contextext.ErrValueMapNotFound) {
			//
			// This should never happen, but it is better to be
			// safe than sorry in case the implementation of
			// context value maps should change in the future.
			//
			return nil, errors.WithStack(err)
		}
	}

	// Is the transaction non-NIL?
	if txn != nil {
		return nil, errors.Wrap(
			ErrTxnStart,
			"already in a transaction",
		)
	}

	// Begin the transaction.
	ntx, err := obj.real.BeginTxx(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(ErrTxnStart, err.Error())
	}

	// Set up a new context.
	nctx, err := setTx(ctx, ntx)
	if err != nil {
		return nil, errors.Wrap(ErrTxnContext, err.Error())
	}

	return nctx, nil
}

func (obj *database) Commit(ctx context.Context) error {
	txn, err := getTx(ctx)

	if err != nil {
		return errors.Wrap(
			err,
			"could not get transaction from context",
		)
	}

	if err := txn.Commit(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (obj *database) Rollback(ctx context.Context) error {
	txn, err := getTx(ctx)

	if err != nil {
		return errors.WithStack(err)
	}

	if err := txn.Rollback(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func fromContext(ctx context.Context, key string) (any, error) {
	vmap, err := contextext.GetValueMap(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rval, found := vmap.Get(key)
	if !found {
		return nil, errors.WithStack(errKeyNotFound)
	}

	return rval, nil
}

func getTx(ctx context.Context) (*sqlx.Tx, error) {
	val, err := fromContext(ctx, KeyTransaction)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	txn, ok := val.(*sqlx.Tx)
	if !ok {
		return nil, errors.WithStack(ErrTxnKeyNotTxn)
	}

	return txn, nil
}

// Set the context's transaction value in the value map.
//
// This will never return an error condition in it's current state, but should
// retain the error handling in case the implementation of contextext changes.
//
//nolint:unparam
func setTx(ctx context.Context, txn *sqlx.Tx) (context.Context, error) {
	var (
		vmap contextext.ValueMap
		err  error
	)

	vmap, err = contextext.GetValueMap(ctx)
	if err != nil {
		vmap = contextext.NewValueMap()
	}

	vmap.Set(KeyTransaction, txn)

	return contextext.WithValueMap(ctx, vmap), nil
}

func FromDB(db *sql.DB, driver string) Database {
	return &database{
		real:   sqlx.NewDb(db, driver),
		driver: driver,
	}
}

// Open a connection using the relevant driver to the given data source name.
func Open(driver string, dsn string) (Database, error) {
	db, err := sqlx.Connect(driver, dsn)
	obj := &database{
		real:   db,
		driver: driver,
	}

	return obj, errors.WithStack(err)
}

// database.go ends here.
