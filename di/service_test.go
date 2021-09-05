/*
 * service_test.go --- Di service tests.
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

package di

import (
	"testing"
)

// ==================================================================
// {{{ `Test` service:

type TestService struct {
	ok bool
}

func (t *TestService) Do() {
	t.ok = true
}

func (t *TestService) IsOK() bool {
	return t.ok == true
}

// }}}
// ==================================================================

// Test we can add services.
func TestAddService(t *testing.T) {
	t.Log("Can we add a new service?")

	di := GetInstance()
	if di == nil {
		t.Error("No, service instance was nil.")
		return
	}

	di.Add("test1", &TestService{})

	inst, found := di.Get("test1")
	if !found {
		t.Error("No, service could not be found.")
		return
	}

	if inst != nil {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

// Test service creation.
func TestCreation(t *testing.T) {
	t.Log("Can we create a new service?")

	di := GetInstance()
	if di == nil {
		t.Error("No, service instance was nil.")
		return
	}

	di.Create(
		"test2",
		func() interface{} {
			return &TestService{}
		},
	)

	inst, found := di.Get("test2")
	if !found {
		t.Error("No, service could not be found.")
		return
	}

	if inst != nil {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

// Test if we can get a service.
func TestGet(t *testing.T) {
	t.Log("Can we get services as expected?")

	di := GetInstance()
	if di == nil {
		t.Error("No, service instance was nil.")
		return
	}

	_, f1 := di.Get("test1")
	_, f2 := di.Get("nosuchthing")

	if f1 {
		t.Log("Yes, we can find existing services.")
	} else {
		t.Error("No, it did not find an existing service.")
	}

	if !f2 {
		t.Log("Yes, it did not find a non-existing service.")
	} else {
		t.Error("No, it managed to find a non-existing service.")
	}
}

// Does it return the same instance for multiple calls?
func TestInstance(t *testing.T) {
	t.Log("Do we get the same object for multiple gets?")

	di := GetInstance()
	if di == nil {
		t.Error("No, service instance was nil.")
		return
	}

	// Two 'pointers' to the same instance.
	t1, _ := di.Get("test1")
	t2, _ := di.Get("test1")

	t1.(*TestService).Do()

	if t2.(*TestService).IsOK() == true {
		t.Log("Yes.")
		return
	}

	t.Error("No.")
}

// Test various utilities.
func TestUtils(t *testing.T) {
	t.Log("Do the utility functions work?")

	di := GetInstance()
	if di == nil {
		t.Error("No, service instance was nil.")
		return
	}

	c := di.Count()
	l := di.Services()

	if c == 2 {
		t.Log("`Count` works as expected.")
	} else {
		t.Error("`Count` failed.")
	}

	if len(l) == 2 {
		t.Log("`Services` works as expected.")
	} else {
		t.Error("`Services` failed.")
	}
}

/* service_test.go ends here. */
