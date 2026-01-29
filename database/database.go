// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// database.go --- Mockable SQL interface.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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

// * Package:

package database

// * Imports:

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	ErrServerConnClosed = errors.Base("server connection closed")
	ErrLostConn         = errors.Base("lost connection during query")
	ErrTxnDeadlock      = errors.Base("deadlock found when trying to get lock")
	ErrTxnSerialization = errors.Base("serialization failure")
)

// * Code:
// ** Interface:

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

	// Rebind query placeholders to the chosen SQL backend.
	Rebind(string) string

	// Parses the given error looking for common MySQL error conditions.
	//
	// If one is found, then a Golang error describing the condition is
	// raised.
	//
	// If nothing interesting is found, then the original error is
	// returned.
	GetError(error) error

	// Run a query function within the context of a database transaction.
	//
	// If there is no error, then the transaction is committed.
	//
	// If there is an error, then the transaction is rolled back.
	WithTransaction(context.Context, TxnFn) error

	// Exposes the database's pool as a `Runner`.
	Runner() Runner
}

// ** Type:

// Concrete database type.
type database struct {
	real   *sqlx.DB
	driver string
}

// ** Methods:

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

// Rebind query placeholders to the chosen SQL backend.
func (obj *database) Rebind(query string) string {
	bind := sqlx.QUESTION

	switch obj.driver {
	case "postgres", "pgx", "pgx/v5":
		bind = sqlx.DOLLAR
	}

	return sqlx.Rebind(bind, query)
}

// Expose the database's pool as a runner.
func (obj *database) Runner() Runner {
	return obj.real
}

// ** Functions:

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
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &database{real: db, driver: driver}, nil
}

// database.go ends here.
