// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// databasemgr.go --- Database manager.
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
//
// mock:yes
//go:generate go run github.com/Asmodai/gohacks/cmd/digen -pattern .
//di:gen basename=Manager key=gohacks/database@v1 type=Manager fallback=NewManager()

package database

import (
	// This is the MySQL driver, it must be blank.
	_ "github.com/go-sql-driver/mysql"

	"gitlab.com/tozd/go/errors"
)

/*
Database management.

This is a series of wrappers around Go's internal DB stuff to ensure
that we set up max idle/open connections et al.
*/
type Manager interface {
	Open(string, string) (Database, error)
	OpenConfig(*Config) (Database, error)
	CheckDB(Database) error
}

// Internal implementation.
type manager struct {
}

// Create a new manager.
func NewManager() Manager {
	return &manager{}
}

// Open a connection to the database specified in the DSN string.
func (dbm *manager) Open(driver string, dsn string) (Database, error) {
	return Open(driver, dsn)
}

// Open and configure a database connection.
func (dbm *manager) OpenConfig(conf *Config) (Database, error) {
	dbase, err := dbm.Open(conf.Driver, conf.ToDSN())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if conf.SetPoolLimits {
		dbase.SetMaxIdleConns(conf.MaxIdleConns)
		dbase.SetMaxOpenConns(conf.MaxOpenConns)
	}

	return dbase, nil
}

// Check the db connection.
func (dbm *manager) CheckDB(db Database) error {
	return errors.WithStack(db.Ping())
}

// databasemgr.go ends here.
