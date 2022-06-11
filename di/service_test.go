/*
 * service_test.go --- Di service tests.
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

// ==================================================================
// {{{ Test class:

type TestClass struct {
	SomeVar int
}

func MakeTestClass() *TestClass {
	return &TestClass{SomeVar: 0}
}

func (tc *TestClass) Whee() int {
	return 42
}

func ConstructTestClass() func() interface{} {
	return func() interface{} {
		return MakeTestClass()
	}
}

// }}}
// ==================================================================

func TestClasses(t *testing.T) {
	di := GetInstance()

	t.Run("It registers", func(t *testing.T) {
		di.AddClass("TestClass", ConstructTestClass())
	})

	t.Run("It locates", func(t *testing.T) {
		inst, found := di.CreateNew("TestClass")
		if !found {
			t.Error("Class was not found")
			return
		}

		if inst == nil {
			t.Error("Instance was nil")
			return
		}

		if inst.(*TestClass).Whee() != 42 {
			t.Error("something went wrong!")
		}
	})
}

func TestInstances(t *testing.T) {
	di := GetInstance()
	if di == nil {
		t.Error("Service instance was nil.")
		return
	}

	// Test we can add services.
	t.Run("Can we add services?", func(t *testing.T) {
		di.Add("test1", &TestService{})

		inst, found := di.Get("test1")
		if !found {
			t.Error("No, service could not be found.")
			return
		}

		if inst == nil {
			t.Error("No.")
			return
		}
	})

	t.Run("Can we get the service?", func(t *testing.T) {
		_, f1 := di.Get("test1")
		_, f2 := di.Get("nosuchthing")

		if !f1 {
			t.Error("No, it did not find an existing service.")
		}

		if f2 {
			t.Error("No, it managed to find a non-existing service.")
		}
	})

	t.Run("Does it return the same instance?", func(t *testing.T) {
		// Two 'pointers' to the same instance.
		t1, _ := di.Get("test1")
		t2, _ := di.Get("test1")

		t1.(*TestService).Do()

		if t2.(*TestService).IsOK() == false {
			t.Error("No.")
		}
	})

	t.Run("Does counting work?", func(t *testing.T) {
		s := di.CountServices()
		c := di.CountClasses()

		if s != 1 {
			t.Errorf("`CountServices` failed, got %d, expected 2", s)
		}

		if c != 1 {
			t.Error("`CountClasses` failed.")
		}
	})

	t.Run("Can we get a list of keys?", func(t *testing.T) {
		s := di.Services()
		c := di.Classes()

		if len(s) != 1 {
			t.Error("`Services` failed.")
		}

		if len(c) != 1 {
			t.Error("`Classes` failed.")
		}
	})
}

/* service_test.go ends here. */
