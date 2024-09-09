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
	ctxvalmap "github.com/Asmodai/gohacks/context"

	// This is the MySQL driver, it must be blank.
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"gitlab.com/tozd/go/errors"

	"context"
	"database/sql"
)

var (
	ErrNoContextKey       error = errors.Base("no context key given")
	ErrValueIsNotDatabase error = errors.Base("not a database")
)

/*
SQL proxy object.

This trainwreck exists so that we can make use of database interfaces.
*/
type Database interface {
	MustBegin() *sqlx.Tx
	Begin() (*sql.Tx, error)
	Beginx() (*sqlx.Tx, error)
	Close() error
	Exec(string, ...any) (sql.Result, error)
	NamedExec(string, any) (sql.Result, error)
	Ping() error
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...any) (Rows, error)
	Queryx(string, ...any) (Rowsx, error)
	QueryRowx(string, ...any) Row
	Select(any, string, ...any) error
	Get(any, string, ...any) error
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
}

type database struct {
	real *sqlx.DB
}

func (db *database) MustBegin() *sqlx.Tx {
	return db.real.MustBegin()
}

func (db *database) Begin() (*sql.Tx, error) {
	rval, err := db.real.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *database) Beginx() (*sqlx.Tx, error) {
	rval, err := db.real.Beginx()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *database) Ping() error  { return errors.WithStack(db.real.Ping()) }
func (db *database) Close() error { return errors.WithStack(db.real.Close()) }

func (db *database) SetMaxIdleConns(limit int) { db.real.SetMaxIdleConns(limit) }
func (db *database) SetMaxOpenConns(limit int) { db.real.SetMaxOpenConns(limit) }

func (db *database) Prepare(query string) (*sql.Stmt, error) {
	rval, err := db.real.Prepare(query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *database) Exec(query string, args ...any) (sql.Result, error) {
	rval, err := db.real.Exec(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *database) NamedExec(query string, args any) (sql.Result, error) {
	rval, err := db.real.NamedExec(query, args)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *database) Query(query string, args ...any) (Rows, error) {
	//nolint:rowserrcheck
	rval, err := db.real.Query(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *database) Queryx(query string, args ...any) (Rowsx, error) {
	rval, err := db.real.Queryx(query, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return rval, nil
}

func (db *database) QueryRowx(query string, args ...any) Row {
	return db.real.QueryRowx(query, args...)
}

func (db *database) Select(what any, query string, args ...any) error {
	return errors.WithStack(db.real.Select(what, query, args...))
}

func (db *database) Get(what any, query string, args ...any) error {
	return errors.WithStack(db.real.Get(what, query, args...))
}

func FromContext(ctx context.Context, key string) (Database, error) {
	vmap, err := ctxvalmap.GetValueMap(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rval, found := vmap.Get(key)
	if !found {
		return nil, errors.WithStack(ErrNoContextKey)
	}

	dbval, ok := rval.(Database)
	if !ok {
		return nil, errors.WithStack(ErrValueIsNotDatabase)
	}

	return dbval, nil
}

func ToContext(ctx context.Context, inst Database, key string) (context.Context, error) {
	var (
		vmap ctxvalmap.ValueMap
		err  error
	)

	vmap, err = ctxvalmap.GetValueMap(ctx)
	if err != nil {
		vmap = ctxvalmap.NewValueMap()
	}

	vmap.Set(key, inst)

	ctx = ctxvalmap.WithValueMap(ctx, vmap)

	return ctx, nil
}

func Open(driver string, dsn string) (Database, error) {
	db, err := sqlx.Open(driver, dsn)

	return &database{
		real: db,
	}, errors.WithStack(err)
}

// database.go ends here.
