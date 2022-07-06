/*
 * database.go --- Mockable SQL interface.
 *
 * Copyright (c) 2021-2022 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be
 * included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package database

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

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

func (db *Database) MustBegin() *sqlx.Tx { return db.real.MustBegin() }

func (db *Database) Begin() (*sql.Tx, error)   { return db.real.Begin() }
func (db *Database) Beginx() (*sqlx.Tx, error) { return db.real.Beginx() }

func (db *Database) Ping() error  { return db.real.Ping() }
func (db *Database) Close() error { return db.real.Close() }

func (db *Database) SetMaxIdleConns(limit int) { db.real.SetMaxIdleConns(limit) }
func (db *Database) SetMaxOpenConns(limit int) { db.real.SetMaxOpenConns(limit) }

func (db *Database) Prepare(query string) (*sql.Stmt, error) {
	return db.real.Prepare(query)
}

func (db *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.real.Exec(query, args...)
}

func (db *Database) NamedExec(query string, args interface{}) (sql.Result, error) {
	return db.real.NamedExec(query, args)
}

func (db *Database) Query(query string, args ...interface{}) (IRows, error) {
	return db.real.Query(query, args...)
}

func (db *Database) Queryx(query string, args ...interface{}) (IRowsx, error) {
	return db.real.Queryx(query, args...)
}

func (db *Database) QueryRowx(query string, args ...interface{}) IRow {
	return db.real.QueryRowx(query, args...)
}

func (db *Database) Select(what interface{}, query string, args ...interface{}) error {
	return db.real.Select(what, query, args...)
}

func (db *Database) Get(what interface{}, query string, args ...interface{}) error {
	return db.real.Get(what, query, args...)
}

func Open(driver string, dsn string) (IDatabase, error) {
	db, err := sqlx.Open(driver, dsn)

	return &Database{
		real: db,
	}, err
}

type Tx struct {
	real *sqlx.Tx
}

func (tx *Tx) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return tx.real.NamedExec(query, arg)
}

func (tx *Tx) Commit() {
	tx.real.Commit()
}

/* database.go ends here. */
