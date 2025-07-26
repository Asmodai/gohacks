// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// queries_test.go --- DB query tests.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>

// * Comments:
//
//
// * End of Comments.

// The DB query tests. package.
package database

// * Imports:
import (
	"github.com/DATA-DOG/go-sqlmock"

	"context"
	"fmt"
	"testing"
)

// * Code:

// ** Query helpers:

func helperMakeError(msg, query string, args ...any) string {
	return fmt.Sprintf(
		"error executing '%s' [%v]: %s",
		query,
		args,
		msg,
	)
}

// *** Simple Query:
// **** Types:

type SimpleQuery struct {
	Id   int
	Name string
}

// **** Variables:

var (
	BasicQueryString = "SELECT id, name FROM foo"
	BasicQueryResult = sqlmock.NewRows([]string{"id", "name"}).
				AddRow(1, "Dave").
				AddRow(2, "Kevin")
)

// **** Methods:

func (q *SimpleQuery) Make() string           { return BasicQueryString }
func (q *SimpleQuery) Results() *sqlmock.Rows { return BasicQueryResult }

func (q *SimpleQuery) Error(msg string, args ...any) string {
	return helperMakeError(msg, BasicQueryString, args...)
}

// *** Query with arguments:
// **** Types:

type ArgQuery struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

// **** Variables:

var (
	ArgQueryString = "SELECT id, name FROM foo WHERE id=:id"
	ArgQueryExec   = "SELECT id, name FROM foo WHERE id=?"
	ArgQueryResult = sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Dave")
)

// **** Functions:

func (q *ArgQuery) Make() string           { return ArgQueryString }
func (q *ArgQuery) Exec() string           { return ArgQueryExec }
func (q *ArgQuery) Results() *sqlmock.Rows { return ArgQueryResult }

func (q *ArgQuery) Error(msg string, args ...any) string {
	return helperMakeError(msg, ArgQueryString, args...)
}

// *** Transaction helpers:
// **** Types:

type sqlFn func(context.Context) error

// **** Functions:

func doSqlTxn(ctx context.Context, obj Database, fn sqlFn) error {
	txctx, err := obj.Begin(ctx)
	if err != nil {
		return err
	}
	if txctx == nil {
		return fmt.Errorf("Transaction context is nil")
	}

	err = fn(txctx)
	if err != nil {
		return err
	}

	err = obj.Commit(txctx)
	if err != nil {
		return fmt.Errorf("Unexpected error: %v", err.Error())
	}

	return nil
}

// ** Query wrappers:

func DoSelect(ctx context.Context, obj Database, dest any, query string, args ...any) error {
	return doSqlTxn(ctx, obj, func(tx context.Context) error {
		return Select(tx, dest, query, args...)
	})
}

func DoQueryx(ctx context.Context, obj Database, query string, args ...any) error {
	return doSqlTxn(ctx, obj, func(tx context.Context) error {
		_, err := Queryx(tx, query, args...)
		return err
	})
}

func DoQueryxContext(ctx context.Context, obj Database, query string, args ...any) error {
	return doSqlTxn(ctx, obj, func(tx context.Context) error {
		_, err := QueryxContext(tx, query, args...)
		return err
	})
}

func DoNamedExec(ctx context.Context, obj Database, query string, arg any) error {
	return doSqlTxn(ctx, obj, func(tx context.Context) error {
		_, err := NamedExec(tx, query, arg)
		return err
	})
}

func DoExec(ctx context.Context, obj Database, query string, args ...any) error {
	return doSqlTxn(ctx, obj, func(tx context.Context) error {
		_, err := Exec(tx, query, args...)
		return err
	})
}

func DoPrepare(ctx context.Context, obj Database, query string, args ...any) error {
	return doSqlTxn(ctx, obj, func(tx context.Context) error {
		_, err := Prepare(tx, query, args...)
		return err
	})
}

func DoExecStmt(ctx context.Context, obj Database, stmt *stmt, args ...any) error {
	return doSqlTxn(ctx, obj, func(tx context.Context) error {
		_, err := ExecStmt(tx, stmt, args...)
		return err
	})
}

// ** Tests:

type TestFn func(context.Context, Database, sqlmock.Sqlmock) error

func RunTxTest(tfn TestFn) error {
	db, mock, err := sqlmock.New()
	if err != nil {
		return fmt.Errorf("could not mockDB: %w", err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	obj := FromDB(db, "sqlmock")
	defer obj.Close()

	return tfn(ctx, obj, mock)
}

// *** `Get`:

func DoGet(ctx context.Context, obj Database, dest any, query string, args ...any) error {
	return doSqlTxn(ctx, obj, func(tx context.Context) error {
		return Get(tx, dest, query, args...)
	})
}

func TestQueryGet(t *testing.T) {
	result := sqlmock.NewRows([]string{"name"}).AddRow("Dave")
	qstr := "SELECT name FROM foo WHERE id = 1"
	data := ""

	// Take note that `Get' is designed for simple `SELECT' queries and
	// does not really like scannables.
	t.Run("Success", func(t *testing.T) {
		RunTxTest(func(ctx context.Context, db Database, mock sqlmock.Sqlmock) error {
			mock.ExpectBegin()
			mock.ExpectQuery(qstr).WillReturnRows(result)
			mock.ExpectCommit()

			err := DoGet(ctx, db, &data, qstr, nil)
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			}

			return err
		})
	})

	t.Run("Failure", func(t *testing.T) {
	})

	t.Run("Invalid context", func(t *testing.T) {
	})
}

func TestQuerySuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Could not mock DB: %v", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	obj := FromDB(db, "sqlmock")
	defer obj.Close()

	t.Run("Get", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnRows(query.Results())
		mock.ExpectCommit()

		DoGet(ctx, obj, &query, qstr, nil)
	})

	t.Run("Select", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnRows(query.Results())
		mock.ExpectCommit()

		DoSelect(ctx, obj, &query, qstr, nil)
	})

	t.Run("Queryx", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnRows(query.Results())
		mock.ExpectCommit()

		DoQueryx(ctx, obj, qstr, nil)
	})

	t.Run("QueryxContext", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnRows(query.Results())
		mock.ExpectCommit()

		DoQueryxContext(ctx, obj, qstr, nil)
	})

	t.Run("NamedExec", func(t *testing.T) {
		query := &ArgQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnRows(query.Results())
		mock.ExpectCommit()

		DoNamedExec(ctx, obj, qstr, 1)
	})

	t.Run("Exec", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnRows(query.Results())
		mock.ExpectCommit()

		DoExec(ctx, obj, qstr, nil)
	})

	t.Run("Prepare", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnRows(query.Results())
		mock.ExpectCommit()

		DoPrepare(ctx, obj, qstr, nil)
	})

	t.Run("ExecStmt", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnRows(query.Results())
		mock.ExpectCommit()

		stmt := &stmt{
			query: qstr,
		}

		DoExecStmt(ctx, obj, stmt, nil)
	})
}

func TestQueryFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Could not mock DB: %v", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	obj := FromDB(db, "sqlmock")
	defer obj.Close()

	errmsg := "out of cheese"

	t.Run("Get", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnError(fmt.Errorf(errmsg))

		err := DoGet(ctx, obj, &query, qstr, nil)
		if err.Error() != query.Error(errmsg, nil) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("Select", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnError(fmt.Errorf(errmsg))

		err := DoSelect(ctx, obj, &query, qstr, nil)
		if err.Error() != query.Error(errmsg, nil) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("Queryx", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnError(fmt.Errorf(errmsg))

		err := DoQueryx(ctx, obj, qstr, nil)
		if err.Error() != query.Error(errmsg, nil) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("QueryxContext", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnError(fmt.Errorf(errmsg))

		err := DoQueryxContext(ctx, obj, qstr, nil)
		if err.Error() != query.Error(errmsg, nil) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("NamedExec", func(t *testing.T) {
		query := &ArgQuery{}
		qstr := query.Make()
		args := map[string]any{"id": 1}

		mock.ExpectBegin()
		mock.ExpectExec(query.Exec()).
			WillReturnError(fmt.Errorf(errmsg))

		err := DoNamedExec(ctx, obj, qstr, args)
		if err.Error() != query.Error(errmsg, args) {
			t.Errorf("Unexpected error: %v != %v", err.Error(), query.Error(errmsg, args))
		}
	})

	t.Run("Exec", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectExec(qstr).
			WillReturnError(fmt.Errorf(errmsg))

		err := DoExec(ctx, obj, qstr, nil)
		if err.Error() != query.Error(errmsg, nil) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})

	t.Run("Prepare", func(t *testing.T) {
		query := &SimpleQuery{}
		qstr := query.Make()

		mock.ExpectBegin()
		mock.ExpectQuery(qstr).
			WillReturnError(fmt.Errorf(errmsg))

		err := DoQueryx(ctx, obj, qstr, nil)
		if err.Error() != query.Error(errmsg, nil) {
			t.Errorf("Unexpected error: %v", err.Error())
		}
	})
}

// * queries_test.go ends here.
