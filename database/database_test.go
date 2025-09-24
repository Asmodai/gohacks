// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// database_test.go --- Database driver tests.
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

// * Package:

package database

// * Imports:

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// * Code:

// helper to build a *database from a sqlmock DB with a given driver string
func newTestDB(t *testing.T, driver string) (*database, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new: %v", err)
	}

	// wrap in sqlx and our database type
	sqlxDB := sqlx.NewDb(db, driver)

	return &database{real: sqlxDB, driver: driver}, mock, db
}

func TestRebind_MySQL(t *testing.T) {
	d, _, db := newTestDB(t, "mysql")
	defer db.Close()

	q := "INSERT INTO x(a,b,c) VALUES (?, ?, ?)"
	got := d.Rebind(q)
	if got != q {
		t.Fatalf("mysql rebind changed query: %q -> %q", q, got)
	}
}

func TestRebind_Postgres(t *testing.T) {
	d, _, db := newTestDB(t, "postgres")
	defer db.Close()

	q := "INSERT INTO x(a,b,c) VALUES (?, ?, ?)"
	got := d.Rebind(q)
	want := "INSERT INTO x(a,b,c) VALUES ($1, $2, $3)"
	if got != want {
		t.Fatalf("postgres rebind = %q, want %q", got, want)
	}
}

func TestRebind_PgxV5(t *testing.T) {
	d, _, db := newTestDB(t, "pgx/v5")
	defer db.Close()

	q := "UPDATE t SET a=?, b=? WHERE id=?"
	got := d.Rebind(q)
	want := "UPDATE t SET a=$1, b=$2 WHERE id=$3"
	if got != want {
		t.Fatalf("pgx/v5 rebind = %q, want %q", got, want)
	}
}

func TestGetError_MySQL_Mappings(t *testing.T) {
	d, _, db := newTestDB(t, "mysql")
	defer db.Close()

	tests := []struct {
		err  error
		want error
	}{
		{&mysql.MySQLError{Number: 1213, Message: "deadlock"}, ErrTxnDeadlock},
		{&mysql.MySQLError{Number: 2006, Message: "server gone"}, ErrServerConnClosed},
		{&mysql.MySQLError{Number: 2013, Message: "lost conn"}, ErrLostConn},
	}

	for _, tt := range tests {
		got := d.GetError(tt.err)

		if !errors.Is(got, tt.want) {
			t.Fatalf("GetError(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestGetError_Postgres_Mappings_pgx(t *testing.T) {
	d, _, db := newTestDB(t, "pgx")
	defer db.Close()

	tests := []struct {
		err  error
		want error
	}{
		{&pgconn.PgError{Code: "40P01", Message: "deadlock"}, ErrTxnDeadlock},
		{&pgconn.PgError{Code: "40001", Message: "serialization"}, ErrTxnSerialization},
		{&pgconn.PgError{Code: "08006", Message: "conn failure"}, ErrLostConn},
		{&pgconn.PgError{Code: "57P01", Message: "admin shutdown"}, ErrServerConnClosed},
	}

	for _, tt := range tests {
		got := d.GetError(tt.err)

		if !errors.Is(got, tt.want) {
			t.Fatalf("GetError(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestGetError_Postgres_Mappings_libpq(t *testing.T) {
	d, _, db := newTestDB(t, "postgres")
	defer db.Close()

	tests := []struct {
		err  error
		want error
	}{
		{&pq.Error{Code: "40P01"}, ErrTxnDeadlock},
		{&pq.Error{Code: "40001"}, ErrTxnSerialization},
		{&pq.Error{Code: "08006"}, ErrLostConn},
		{&pq.Error{Code: "57P01"}, ErrServerConnClosed},
	}

	for _, tt := range tests {
		got := d.GetError(tt.err)

		if !errors.Is(got, tt.want) {
			t.Fatalf("GetError(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestWithTransaction_CommitHappyPath(t *testing.T) {
	d, mock, raw := newTestDB(t, "mysql")
	defer raw.Close()

	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(name) VALUES (?)")).
		WithArgs("Ada").
		WillReturnResult(sqlmock.NewResult(123, 1))
	mock.ExpectCommit()

	err := d.WithTransaction(ctx, func(ctx context.Context, r Runner) error {
		q := d.Rebind("INSERT INTO users(name) VALUES (?)")
		_, err := r.ExecContext(ctx, q, "Ada")
		return err
	})

	if err != nil {
		t.Fatalf("WithTransaction commit path: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWithTransaction_RollbackOnError(t *testing.T) {
	d, mock, raw := newTestDB(t, "mysql")
	defer raw.Close()

	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users(name) VALUES (?)")).
		WithArgs("Ada").
		WillReturnError(errors.New("boom"))
	mock.ExpectRollback()

	err := d.WithTransaction(ctx, func(ctx context.Context, r Runner) error {
		q := d.Rebind("INSERT INTO users(name) VALUES (?)")
		_, execErr := r.ExecContext(ctx, q, "Ada")
		if execErr != nil {
			return execErr
		}
		return nil
	})

	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWithTransaction_RetryOnMySQLDeadlock(t *testing.T) {
	d, mock, raw := newTestDB(t, "mysql")
	defer raw.Close()

	ctx := context.Background()

	// Attempt 1: deadlock -> rollback
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE t SET v=? WHERE id=?")).
		WithArgs(1, 42).
		WillReturnError(&mysql.MySQLError{Number: 1213, Message: "deadlock"})
	mock.ExpectRollback()

	// Attempt 2: success -> commit
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE t SET v=? WHERE id=?")).
		WithArgs(1, 42).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	start := time.Now()
	err := d.WithTransaction(ctx, func(ctx context.Context, r Runner) error {
		q := d.Rebind("UPDATE t SET v=? WHERE id=?")
		_, execErr := r.ExecContext(ctx, q, 1, 42)
		return execErr
	})
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}

	// sanity: ensure backoff didn’t explode (should be small)
	if elapsed > 2*time.Second {
		t.Fatalf("unexpected long retry backoff: %v", elapsed)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWithTransaction_RetryOnPgSerialization(t *testing.T) {
	d, mock, raw := newTestDB(t, "pgx")
	defer raw.Close()

	ctx := context.Background()

	// Attempt 1: serialization failure -> rollback
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE t SET v=$1 WHERE id=$2")).
		WithArgs(1, 42).
		WillReturnError(&pgconn.PgError{Code: "40001", Message: "serialization failure"})
	mock.ExpectRollback()

	// Attempt 2: success -> commit
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE t SET v=$1 WHERE id=$2")).
		WithArgs(1, 42).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := d.WithTransaction(ctx, func(ctx context.Context, r Runner) error {
		q := d.Rebind("UPDATE t SET v=? WHERE id=?")
		_, execErr := r.ExecContext(ctx, q, 1, 42)
		return execErr
	})

	if err != nil {
		t.Fatalf("expected retry to succeed, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestWithTransaction_PanicPathRollsBack(t *testing.T) {
	d, mock, raw := newTestDB(t, "mysql")
	defer raw.Close()

	ctx := context.Background()

	mock.ExpectBegin()
	mock.ExpectRollback()

	defer func() {
		if p := recover(); p == nil {
			t.Fatalf("expected panic to propagate")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	}()

	_ = d.WithTransaction(ctx, func(ctx context.Context, r Runner) error {
		panic("kaboom")
	})
}

// database_test.go ends here.
