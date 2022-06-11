/*
 * database.go --- Mockable SQL interface.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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
