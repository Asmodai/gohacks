/*
 * database_test.go --- Database driver tests.
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
	"context"
	"database/sql/driver"
	"testing"
)

// ==================================================================
// {{{ Mock SQL driver:

type MockDriver struct {
}

func (md MockDriver) Open(name string) (driver.Conn, error) {
	return nil, nil
}

// }}}
// ==================================================================

// ==================================================================
// {{{ Mock SQL connection:

type MockConn struct {
}

func (mc MockConn) Connect(ctx context.Context) (driver.Conn, error) {
	return nil, nil
}

func (mc MockConn) Driver() driver.Driver {
	return MockDriver{}
}

// }}}
// ==================================================================

func TestDatabaseOpen(t *testing.T) {
	t.Log("Does `Open` return an error if it cannot connect?")

	_, e1 := Open("nil", "nil")
	if e1 != nil {
		t.Log("Yes.")
		return
	}

	t.Error("No!")
}

/* database_test.go ends here. */
