/*
 * idatabase.go --- Database interface.
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
	"github.com/jmoiron/sqlx"

	"database/sql"
)

type ITx interface {
	NamedExec(string, interface{}) (sql.Result, error)
	Commit()
}

// Interface for `sql.Row` objects.
type IRow interface {
	Err() error
	Scan(dest ...interface{}) error
}

// Interface for `sqlx.Rows` objects.
type IRowsx interface {
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(...interface{}) error
	StructScan(interface{}) error
}

// Interface for `sql.Rows` objects.
type IRows interface {
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(...interface{}) error
}

// Interface for `sql.DB` objects.
type IDatabase interface {
	MustBegin() *sqlx.Tx
	Begin() (*sql.Tx, error)
	Beginx() (*sqlx.Tx, error)
	Close() error
	Exec(string, ...interface{}) (sql.Result, error)
	NamedExec(string, interface{}) (sql.Result, error)
	Ping() error
	Prepare(string) (*sql.Stmt, error)
	Query(string, ...interface{}) (IRows, error)
	Queryx(string, ...interface{}) (IRowsx, error)
	QueryRowx(string, ...interface{}) IRow
	Select(interface{}, string, ...interface{}) error
	Get(interface{}, string, ...interface{}) error
	SetMaxIdleConns(int)
	SetMaxOpenConns(int)
}

/* idatabase.go ends here. */
