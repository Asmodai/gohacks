/*
 * config_test.go --- SQL config tests.
 *
 * Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
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
	"testing"
)

func MakeSQL() *Config {
	sql := NewConfig()

	sql.Username = "user"
	sql.Password = "pass"
	sql.Hostname = "localhost"
	sql.Port = 1337
	sql.Database = "db"
	sql.BatchSize = 10

	return sql
}

func TestSQLDSN(t *testing.T) {
	var dsn1 string

	sql := MakeSQL()

	t.Run("Does `ToDSN` work as expected?", func(t *testing.T) {
		dsn1 = sql.ToDSN()

		if dsn1 != "user:pass@tcp(localhost:1337)/db?parseTime=True&loc=UTC&time_zone='-00:00'" {
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

/* config_test.go ends here. */
