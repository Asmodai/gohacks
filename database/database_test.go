/*
 * database_test.go --- Database driver tests.
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
