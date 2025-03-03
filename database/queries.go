// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
// SPDX-License-Identifier: MIT
//
// queries.go --- Database queries.
//
// Copyright (c) 2021-2025 Paul Ward <asmodai@gmail.com>
//
// Author:     Paul Ward <asmodai@gmail.com>
// Maintainer: Paul Ward <asmodai@gmail.com>

// * Comments:
//
//
// * End of Comments.

// The Database queries. package.
package database

// * Imports:
import (
	// This is the MySQL driver, it must be blank.
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/tozd/go/errors"

	"context"
	"database/sql"
)

// * Code:

// ** Types:

type stmt struct {
	*sql.Stmt
	query string
}

// ** Functions:

// Wrapper around `Tx.Get`.
//
// The transaction should be passed via a context value.
func Get(ctx context.Context, dest any, query string, args ...any) error {
	tx, err := getTx(ctx)
	if err != nil {
		return errors.WithStack(sqlError(err, query, args...))
	}

	err = tx.Get(dest, query, args...)

	return sqlError(err, query, args...)
}

// Wrapper around `Tx.Select`.
//
// The transaction should be passed via a context value.
func Select(ctx context.Context, dest any, query string, args ...any) error {
	tx, err := getTx(ctx)
	if err != nil {
		return errors.WithStack(sqlError(err, query, args...))
	}

	err = tx.Select(dest, query, args...)

	return sqlError(err, query, args...)
}

// Wrapper around `Tx.Queryx`.
//
// The transaction should be passed via a context value.
func Queryx(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, errors.WithStack(sqlError(err, query, args...))
	}

	ret, err := tx.Queryx(query, args...)

	return ret, errors.WithStack(sqlError(err, query, args...))
}

// Wrapper around `Tx.QueryxContext`.
//
// The transaction should be passed via a context value.
func QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, errors.WithStack(sqlError(err, query, args...))
	}

	ret, err := tx.QueryxContext(ctx, query, args...)

	return ret, errors.WithStack(sqlError(err, query, args...))
}

// Wrapper around `Tx.NamedExec`.
//
// The transaction should be passed via a context value.
func NamedExec(ctx context.Context, query string, arg any) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, errors.WithStack(sqlError(err, query, arg))
	}

	ret, err := tx.NamedExec(query, arg)

	return ret, errors.WithStack(sqlError(err, query, arg))
}

// Wrapper around `Tx.Exec`.
//
// The transaction should be passed via a context value.
func Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, errors.WithStack(sqlError(err, query, args))
	}

	ret, err := tx.Exec(query, args...)

	return ret, errors.WithStack(sqlError(err, query, args...))
}

// Wrapper around `Tx.Prepare`.
//
// The transaction should be passed via a context value.
func Prepare(ctx context.Context, query string, args ...any) (*stmt, error) {
	tx, err := getTx(ctx)
	if err != nil {
		return nil, errors.WithStack(sqlError(err, query, args...))
	}

	ret, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.WithStack(sqlError(err, query, args...))
	}

	return &stmt{
		query: query,
		Stmt:  ret,
	}, nil
}

// Wrapper around `Tx.ExecStmt`.
//
// The transaction should be passed via a context value.
func ExecStmt(ctx context.Context, stmt *stmt, args ...any) (sql.Result, error) {
	_, err := getTx(ctx)
	if err != nil {
		return nil, errors.WithStack(sqlError(err, stmt.query, args...))
	}

	ret, err := stmt.ExecContext(ctx, args...)

	return ret, errors.WithStack(sqlError(err, stmt.query, args...))
}

// * queries.go ends here.
