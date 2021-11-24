/*
 * idatabase.go --- Database interface.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
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
