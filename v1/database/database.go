// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// database.go --- Mockable SQL interface.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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
//
// mock:yes

package database

import (
	"github.com/Asmodai/gohacks/v1/contextext"

	// This is the MySQL driver, it must be blank.
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/tozd/go/errors"

	"context"
	"database/sql"
	"fmt"
	"strings"
)

var (
	errKeyNotFound error = errors.Base("key not found")

	ErrTxnKeyNotFound error = errors.Base("transaction key not found")
	ErrTxnKeyNotTxn   error = errors.Base("key value is not a transaction")
	ErrTxnContext     error = errors.Base("could not create transaction context")
	ErrTxnStart       error = errors.Base("could not start transaction")

	ErrTxnDeadlock      error = errors.Base("deadlock found when trying to get lock")
	ErrServerConnClosed error = errors.Base("server connection closed")
	ErrLostConn         error = errors.Base("lost connection during query")
)

const (
	KeyTransaction   string = "_DB_TXN"
	StringDeadlock   string = "Error 1213" // Deadlock detected.
	StringConnClosed string = "Error 2006" // MySQL server connection closed.
	StringLostConn   string = "Error 2013" // Lost connection during query.
)

type Database interface {
	// Pings the database connection to ensure it is alive and connected.
	Ping() error

	// Close a database connection.  This does nothing if the connection
	// is already closed.
	Close() error

	// Set the maximum idle connections.
	SetMaxIdleConns(int)

	// Set the maximum open connections.
	SetMaxOpenConns(int)

	// Return the transaction (if any) from the given context.
	Tx(context.Context) (*sqlx.Tx, error)

	// Initiate a transaction.  Returns a new context that contains the
	// database transaction session as a value.
	Begin(context.Context) (context.Context, error)

	// Initiate a transaction commit.
	Commit(context.Context) error

	// Initiate a transaction rollback.
	Rollback(context.Context) error

	// Parses the given error looking for common MySQL error conditions.
	//
	// If one is found, then a Golang error describing the condition is
	// raised.
	//
	// If nothing interesting is found, then the original error is
	// returned.
	GetError(error) error
}

type database struct {
	real   *sqlx.DB
	driver string
}

// Ping a database connection.
func (obj *database) Ping() error {
	return errors.WithStack(obj.real.Ping())
}

// Close an open database connection.
func (obj *database) Close() error {
	return errors.WithStack(obj.real.Close())
}

// Set the maximum number of idle connections that will be supported by the
// database connection.
func (obj *database) SetMaxIdleConns(limit int) {
	obj.real.SetMaxIdleConns(limit)
}

// Set the maximum number of open connections that will be supported by the
// database connection.
func (obj *database) SetMaxOpenConns(limit int) {
	obj.real.SetMaxOpenConns(limit)
}

// Obtain a transaction from a context (if any).
func (obj *database) Tx(ctx context.Context) (*sqlx.Tx, error) {
	return getTx(ctx)
}

// Begin a transaction.
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

// Initiate a transaction commit.
func (obj *database) Commit(ctx context.Context) error {
	txn, err := getTx(ctx)

	if err != nil {
		return errors.Wrap(
			err,
			"could not get transaction from context",
		)
	}

	if err := txn.Commit(); err != nil {
		return obj.GetError(err)
	}

	return nil
}

// Initiate a transaction rollback.
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

// Translate specific MySQL errors into distinct error conditions.
func (obj *database) GetError(err error) error {
	var nerr = err

	switch {
	case strings.Contains(err.Error(), StringDeadlock):
		nerr = ErrTxnDeadlock

	case strings.Contains(err.Error(), StringConnClosed):
		nerr = ErrServerConnClosed

	case strings.Contains(err.Error(), StringLostConn):
		nerr = ErrLostConn
	}

	return errors.WithStack(nerr)
}

// Helper function for obtaining a value from a context's value map.
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

// Get the context's transaction value in the value map.
//
// If there is no value map in the context then contextext's
// `ErrValueMapNotFound` is returned.
//
// If the value for the transaction key is not of type `*sql.Tx` (or cannot be
// coerced to that type) then `ErrTxnKeyNotTxn` is returned.
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

func sqlError(err error, sql string, args ...any) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("error executing '%s' [%v]: %w", sql, args, err)
}

// Create a new database object using an existing `sql` object.
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
