// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// database_test.go --- Database driver tests.
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

package database

import (
	"github.com/Asmodai/gohacks/contextext"
	"github.com/DATA-DOG/go-sqlmock"
	"gitlab.com/tozd/go/errors"

	"context"
	"testing"
)

func TestDatabaseOpen(t *testing.T) {
	t.Run("Returns error if cannot connect", func(t *testing.T) {
		_, e1 := Open("nil", "nil")
		if e1 == nil {
			t.Error("Expected an error condition")
		}
	})
}

func TestBasic(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Errorf("Could not mock DB: %v", err.Error())
		return
	}

	obj := FromDB(db, "sqlmock")
	defer obj.Close()

	t.Run("Ping", func(t *testing.T) {
		mock.ExpectPing()

		obj.Ping()
	})
}

func TestContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	t.Run("getTx() with bad transaction", func(t *testing.T) {
		vmap := contextext.NewValueMap()
		vmap.Set(KeyTransaction, 42)
		nctx := contextext.WithValueMap(ctx, vmap)

		_, err := getTx(nctx)
		if !errors.Is(err, ErrTxnKeyNotTxn) {
			t.Errorf("Unexpected error: %#v", errors.Unwrap(err))
		}
	})

	t.Run("getTx() with missing key", func(t *testing.T) {
		vmap := contextext.NewValueMap()
		nctx := contextext.WithValueMap(ctx, vmap)

		_, err := getTx(nctx)
		if !errors.Is(err, errKeyNotFound) {
			t.Errorf("Unexpected error: %#v", errors.Unwrap(err))
		}
	})
}

func TestTransactionFailures(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Could not mock DB: %v", err.Error())
		return
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	obj := FromDB(db, "sqlmock")
	defer obj.Close()

	t.Run("Begin() ErrTxnStart", func(t *testing.T) {
		mock.ExpectBegin()

		txctx, err := obj.Begin(ctx)
		_, err = obj.Begin(txctx)

		if !errors.Is(err, ErrTxnStart) {
			t.Errorf("Unexpected error: %#v", err)
		}
	})

	t.Run("Begin() error from sqlx", func(t *testing.T) {
		mock.ExpectBegin().WillReturnError(errors.Base("no cheese"))

		_, err := obj.Begin(ctx)

		if err.Error() != "no cheese" {
			t.Errorf("Unexpected error: %#v", err)
		}
	})

	t.Run("Commit() context error", func(t *testing.T) {
		mock.ExpectBegin()

		_, err := obj.Begin(ctx)
		err = obj.Commit(ctx)

		if !errors.Is(err, contextext.ErrValueMapNotFound) {
			t.Errorf("Unexpected error: %#v", errors.Unwrap(err))
		}
	})

	t.Run("Commit() sqlx error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(errors.Base("no cheese"))

		txctx, err := obj.Begin(ctx)
		err = obj.Commit(txctx)

		if err.Error() != "no cheese" {
			t.Errorf("Unexpected error: %#v", errors.Unwrap(err))
		}
	})

	t.Run("Rollback() context error", func(t *testing.T) {
		mock.ExpectBegin()

		_, err := obj.Begin(ctx)
		err = obj.Rollback(ctx)

		if !errors.Is(err, contextext.ErrValueMapNotFound) {
			t.Errorf("Unexpected error: %#v", errors.Unwrap(err))
		}
	})

	t.Run("Rollback() sqlx error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectRollback().WillReturnError(errors.Base("no cheese"))

		txctx, err := obj.Begin(ctx)
		err = obj.Rollback(txctx)

		if err.Error() != "no cheese" {
			t.Errorf("Unexpected error: %#v", errors.Unwrap(err))
		}
	})
}

func TestTransactionCommit(t *testing.T) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		txctx  context.Context
	)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Could not mock DB: %v", err.Error())
		return
	}

	ctx, cancel = context.WithCancel(context.TODO())
	defer cancel()

	obj := FromDB(db, "sqlmock")
	defer obj.Close()

	// Query, mock rows and results.
	query := "SELECT id,name FROM foo"
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Dave").
		AddRow(2, "Charles")

	// Mock our SQL transaction.
	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM foo").WillReturnRows(rows)
	mock.ExpectCommit()

	// Create our transaction.
	txctx, err = obj.Begin(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
		return
	}

	// Is the transaction ok?
	if txctx == nil {
		t.Error("Transaction context is nil")
	}

	// Extract the sqlx transaction from the context.
	txn, err := obj.Tx(txctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
		return
	}

	// Transaction ok?
	if txn == nil {
		t.Error("Transaction is nil")
		return
	}

	// Run our query.
	txn.Query(query)

	// Commit our query.
	err = obj.Commit(txctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
}

func TestTransactionRollback(t *testing.T) {
	var (
		ctx    context.Context
		cancel context.CancelFunc
		txctx  context.Context
	)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Could not mock DB: %v", err.Error())
		return
	}

	ctx, cancel = context.WithCancel(context.TODO())
	defer cancel()

	obj := FromDB(db, "sqlmock")
	defer obj.Close()

	// Query, mock rows and results.
	query := "SELECT id,name FROM foo"
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Dave").
		AddRow(2, "Charles")

	// Mock our SQL transaction.
	mock.ExpectBegin()
	mock.ExpectQuery("^SELECT (.+) FROM foo").WillReturnRows(rows)
	mock.ExpectRollback()

	// Create our transaction.
	txctx, err = obj.Begin(ctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
		return
	}

	// Is the transaction ok?
	if txctx == nil {
		t.Error("Transaction context is nil")
	}

	// Extract the sqlx transaction from the context.
	txn, err := obj.Tx(txctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
		return
	}

	// Transaction ok?
	if txn == nil {
		t.Error("Transaction is nil")
		return
	}

	// Run our query.
	txn.Query(query)

	// Rollback our query.
	err = obj.Rollback(txctx)
	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}
}

func TestSQLErrors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("Could not mock DB: %v", err.Error())
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	obj := FromDB(db, "sqlmock")
	defer obj.Close()

	t.Run("1213: Deadlock found", func(t *testing.T) {
		sqlerr := "Error 1213: Deadlock found when trying to get lock; try restarting transaction"

		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(errors.Base(sqlerr))

		txctx, err := obj.Begin(ctx)
		err = obj.Commit(txctx)

		if !errors.Is(err, ErrTxnDeadlock) {
			t.Errorf("Unexpected error: %#v", err)
		}
	})
}

// database_test.go ends here.
