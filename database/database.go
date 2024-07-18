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

package database

import (

	// This is the MySQL driver, it must be blank.
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/tozd/go/errors"

	"database/sql"
)

/*
SQL proxy object.

This trainwreck exists so that we can make use of database interfaces.

It might be 100% useless, as `sql.DB` will most likely conform to `IDatabase`,
so this file might vanish at some point.
*/
type Database struct {
	real *sqlx.DB
}

func (db *Database) MustBegin() *sqlx.Tx {
	return db.real.MustBegin()
}

func (db *Database) Begin() (*sql.Tx, error) {
	rval, err := db.real.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *Database) Beginx() (*sqlx.Tx, error) {
	rval, err := db.real.Beginx()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *Database) Ping() error  { return errors.WithStack(db.real.Ping()) }
func (db *Database) Close() error { return errors.WithStack(db.real.Close()) }

func (db *Database) SetMaxIdleConns(limit int) { db.real.SetMaxIdleConns(limit) }
func (db *Database) SetMaxOpenConns(limit int) { db.real.SetMaxOpenConns(limit) }

func (db *Database) Prepare(query string) (*sql.Stmt, error) {
	rval, err := db.real.Prepare(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *Database) Exec(query string, args ...any) (sql.Result, error) {
	rval, err := db.real.Exec(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *Database) NamedExec(query string, args any) (sql.Result, error) {
	rval, err := db.real.NamedExec(query, args)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *Database) Query(query string, args ...any) (IRows, error) {
	//nolint:rowserrcheck
	rval, err := db.real.Query(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *Database) Queryx(query string, args ...any) (IRowsx, error) {
	rval, err := db.real.Queryx(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *Database) QueryRowx(query string, args ...any) IRow {
	return db.real.QueryRowx(query, args...)
}

func (db *Database) Select(what any, query string, args ...any) error {
	return errors.WithStack(db.real.Select(what, query, args...))
}

func (db *Database) Get(what any, query string, args ...any) error {
	return errors.WithStack(db.real.Get(what, query, args...))
}

func Open(driver string, dsn string) (IDatabase, error) {
	db, err := sqlx.Open(driver, dsn)

	return &Database{
		real: db,
	}, errors.WithStack(err)
}

type Tx struct {
	real *sqlx.Tx
}

func (tx *Tx) NamedExec(query string, arg any) (sql.Result, error) {
	rval, err := tx.real.NamedExec(query, arg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (tx *Tx) Commit() error {
	return errors.WithStack(tx.real.Commit())
}

// database.go ends here.
