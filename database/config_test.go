// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// config_test.go --- SQL config tests.
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
	"gitlab.com/tozd/go/errors"

	"fmt"
	"testing"
)

const (
	username  string = "user"
	password  string = "pass"
	hostname  string = "localhost"
	dbname    string = "db"
	portno    int    = 1337
	batchsize int    = 10
)

func MakeDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=True&loc=UTC&time_zone='-00:00'",
		username,
		password,
		hostname,
		portno,
		dbname,
	)
}

func MakeSQL() *Config {
	sql := NewConfig()

	sql.Driver = "test"
	sql.Username = username
	sql.Password = password
	sql.Hostname = hostname
	sql.Port = portno
	sql.Database = dbname
	sql.BatchSize = batchsize

	return sql
}

func TestSQLDSN(t *testing.T) {
	var dsn1 string

	sql := MakeSQL()

	t.Run("Does `ToDSN` work as expected?", func(t *testing.T) {
		dsn1 = sql.ToDSN()

		if dsn1 != MakeDSN() {
			t.Errorf("No, got '%v'", dsn1)
		}
	})

	t.Run("Do subsequent calls work?", func(t *testing.T) {
		dsn2 := sql.ToDSN()

		if dsn2 != dsn1 {
			t.Errorf("No, got '%v'", dsn2)
		}
	})
}

func CheckError(cnf *Config) error {
	err := cnf.Validate()

	if err == nil {
		return fmt.Errorf("no error generated")
	}

	return err
}

func TestValidate(t *testing.T) {
	t.Run("Works as expected", func(t *testing.T) {
		sql := MakeSQL()

		err := sql.Validate()
		if err != nil {
			t.Errorf("Unxepected error: %v", err)
		}
	})

	t.Run("Errors with no driver", func(t *testing.T) {
		sql := MakeSQL()
		sql.Driver = ""

		err := CheckError(sql)
		if !errors.Is(err, ErrNoDriver) {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Errors with no username", func(t *testing.T) {
		sql := MakeSQL()
		sql.Username = ""

		err := CheckError(sql)
		if !errors.Is(err, ErrNoUsername) {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Errors with no password", func(t *testing.T) {
		sql := MakeSQL()
		sql.Password = ""

		err := CheckError(sql)
		if !errors.Is(err, ErrNoPassword) {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Errors with no hostname", func(t *testing.T) {
		sql := MakeSQL()
		sql.Hostname = ""

		err := CheckError(sql)
		if !errors.Is(err, ErrNoHostname) {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Errors with no database", func(t *testing.T) {
		sql := MakeSQL()
		sql.Database = ""

		err := CheckError(sql)
		if !errors.Is(err, ErrNoDatabase) {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("Sets default port", func(t *testing.T) {
		sql := MakeSQL()
		sql.Port = 0

		err := sql.Validate()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if sql.Port != defaultDatabasePort {
			t.Errorf(
				"Unexpected port number, %d!=%d",
				defaultDatabasePort,
				sql.Port,
			)
		}
	})
}

//* config_test.go ends here.
