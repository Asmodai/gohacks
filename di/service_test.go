/*
 * service_test.go --- Di service tests.
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
