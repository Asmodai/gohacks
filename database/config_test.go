/*
 * config_test.go --- SQL config tests.
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
	t.Log("Does `ToDSN` work as expected?")

	sql := MakeSQL()
	dsn := sql.ToDSN()

	if dsn == "user:pass@tcp(localhost:1337)/db?parseTime=True&loc=UTC&time_zone='-00:00'" {
		t.Log("Yes.")
		return
	}

	t.Errorf("No, got '%v'", dsn)
}

/* config_test.go ends here. */
