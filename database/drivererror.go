// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// drivererror.go --- Database driver error handling.
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

package database

// * Imports:

import (
	"context"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn" // for pgx
	"github.com/lib/pq"       // for lib/pq

	"gitlab.com/tozd/go/errors"
)

// * Code:

// ** Methods:

// Translate specific MySQL errors into distinct error conditions.
func (obj *database) GetError(err error) error {
	if err == nil {
		return nil
	}

	// Normalise context errors.
	switch {
	case errors.Is(err, context.Canceled):
		return errors.WithStack(err)

	case errors.Is(err, context.DeadlineExceeded):
		return errors.WithStack(err)
	}

	// Dispatch on the database driver.
	switch obj.driver {
	case "mysql":
		if mapped := mapMySQLError(err); mapped != nil {
			return mapped
		}

	case "postgres", "pgx", "pgx/v5":
		if mapped := mapPostgresError(err); mapped != nil {
			return mapped
		}

	default:
		// Unknown driver:  try both mappings in case.
		if mapped := mapMySQLError(err); mapped != nil {
			return mapped
		}

		if mapped := mapPostgresError(err); mapped != nil {
			return mapped
		}
	}

	// Nothing matched, just return original error.
	return errors.WithStack(err)
}

// ** Functions:

// Map MySQL errors.
func mapMySQLError(err error) error {
	var merr *mysql.MySQLError

	if errors.As(err, &merr) {
		switch merr.Number {
		case 1213: // ER_LOCK_DEADLOCK
			return errors.WithStack(ErrTxnDeadlock)

		case 2006: // CR_SERVER_GONE_ERROR
			return errors.WithStack(ErrServerConnClosed)

		case 2013: // CR_SERVER_LOST
			return errors.WithStack(ErrLostConn)
		}
	}

	return nil
}

// Map PostgreSQL errors.
func mapPostgresError(err error) error {
	// pgx first
	var pgerr *pgconn.PgError

	if errors.As(err, &pgerr) {
		return errors.WithStack(mapPgSQLState(pgerr.Code))
	}

	// lib/pq
	var pqerr *pq.Error

	if errors.As(err, &pqerr) {
		return errors.WithStack(mapPgSQLState(string(pqerr.Code)))
	}

	return nil
}

// map SQLSTATE (five-char string) to domain errors.
func mapPgSQLState(code string) error {
	switch code {
	// Concurrency
	case "40P01": // deadlock_detected
		return ErrTxnDeadlock

	case "40001": // serialization_failure
		return ErrTxnSerialization

	// Connection exceptions (class 08)
	case "08000", // connection_exception (class)
		"08003", // connection_does_not_exist
		"08006", // connection_failure
		"08001": // sqlclient_unable_to_establish_sqlconnection
		return ErrLostConn

	// Server shutdowns
	case "57P01": // admin_shutdown
		return ErrServerConnClosed
	}

	// No mapping? nil means "don’t override"
	return nil
}

// * drivererror.go ends here.
