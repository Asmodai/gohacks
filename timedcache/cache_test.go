// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// cache_test.go --- Timed cache tests.
//
// Copyright (c) 2024-2025 Paul Ward <paul@lisphacker.uk>
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

// * Package:

package timedcache

// * Imports:

import (
	"gitlab.com/tozd/go/errors"

	"fmt"
	"testing"
)

// * Constants:

const (
	TestKey1 string = "testKey1"
	TestKey2 string = "testKey2"
	TestKey3 string = "testKey3" // This key should never exist.
	TestKey4 string = "testKey4"

	TestValue1 int = 1
	TestValue2 int = 2
	TestValue3 int = 3
	TestValue4 int = 4
)

// * Variables:

var (
	testCache  TimedCache
	hasEvicted bool = false
	getMetric  int  = 0
	setMetric  int  = 0
	hitMetric  int  = 0
	missMetric int  = 0
)

// * Code:

// ** Helpers:

func OnEvictedEvent(key any, _ any) { hasEvicted = true }

func NewTestConfig() *Config {
	return &Config{
		Name:           "test",
		ExpirationTime: 2,
		OnEvicted:      nil,
	}
}

func NewTestMetricConfig() *Config {
	return &Config{
		Name:           "test",
		ExpirationTime: 2,
		OnEvicted:      nil,
	}
}

// Check if the given key
//
//	a) exists,
//	b) matches the given value.
func CheckKey(cache TimedCache, key any, value any) error {
	if val, fnd := cache.Get(key); fnd {
		if val != value {
			return fmt.Errorf(
				"Unexpected result: %v != %v",
				val,
				value,
			)
		}

		return nil
	}

	return fmt.Errorf("Key '%v' not found.", key)
}

// ** Tests:

// Test the basic accessor methods.
func TestAccessors(t *testing.T) {
	testCache = New(NewTestConfig())

	testCache.(*timedCache).items[TestKey1] = Item{Object: TestValue1}

	//
	// Test Get() with a valid key.
	//
	t.Run("Get() with valid key", func(t *testing.T) {
		err := CheckKey(testCache, TestKey1, TestValue1)
		if err != nil {
			t.Error(err.Error())
		}
	})

	//
	// Test Get() with an invalid key.
	//
	t.Run("Get() with invalid key", func(t *testing.T) {
		if val, fnd := testCache.Get(TestKey2); fnd {
			t.Errorf("Found key and value '%v'", val)
		}
	})

	//
	// Test that Set() creates a new key when required.
	//
	t.Run("Set() adds new key", func(t *testing.T) {
		if val, fnd := testCache.Get(TestKey3); fnd {
			t.Errorf("Found key and value '%v'", val)
			return
		}

		testCache.Set(TestKey3, TestValue3)

		err := CheckKey(testCache, TestKey3, TestValue3)
		if err != nil {
			t.Error(err.Error())
		}
	})

	//
	// Test that Set() modifies an existing key.
	//
	t.Run("Set() modifies existing key", func(t *testing.T) {
		testCache.Set(TestKey1, TestValue2)

		err := CheckKey(testCache, TestKey1, TestValue2)
		if err != nil {
			t.Error(err.Error())
		}
	})

	//
	// Test that Add() creates a new key.
	//
	t.Run("Add() creates new key", func(t *testing.T) {
		if err := testCache.Add(TestKey4, TestValue4); err != nil {
			t.Errorf("No, error was '%v'", err.Error())
			return
		}

		err := CheckKey(testCache, TestKey4, TestValue4)
		if err != nil {
			t.Error(err.Error())
		}
	})

	//
	// Test that Add() errors out when a key exists.
	//
	t.Run("Add() emits errors", func(t *testing.T) {
		err := testCache.Add(TestKey3, TestValue3)

		if !errors.Is(err, ErrKeyExists) {
			t.Errorf("Unexpected error condition: %v", err.Error())
		}
	})

	//
	// Test that Replace() works with existing keys.
	//
	t.Run("Replace() works", func(t *testing.T) {
		var err error

		err = testCache.Replace(TestKey3, TestValue4)

		if err != nil {
			t.Errorf("No, error was '%v'", err.Error())
			return
		}

		err = CheckKey(testCache, TestKey3, TestValue4)
		if err != nil {
			t.Error(err.Error())
		}
	})

	//
	// Test that Replace() errors out when key does not exsit.
	//
	t.Run("Replace() emits errors", func(t *testing.T) {
		err := testCache.Replace("does_not_exist", 76)

		if !errors.Is(err, ErrKeyNotExist) {
			t.Errorf("Unexpected error condition: %v", err.Error())
		}
	})

	//
	// Test that Delete() works.
	//
	t.Run("Delete() works", func(t *testing.T) {
		_, ok := testCache.Delete(TestKey4)
		if !ok {
			t.Error("Could not delete from cache.")
		}

		_, found := testCache.(*timedCache).items[TestKey4]
		if found {
			t.Errorf("Item '%v' is still in the cache!", TestKey4)
		}
	})

	//
	// Test Count().
	//
	t.Run("Count() counts", func(t *testing.T) {
		val := testCache.Count()
		actual := len(testCache.(*timedCache).items)

		if val != actual {
			t.Errorf("Mismatch, %v != %v", val, actual)
		}
	})

	//
	// Test Flush()
	//
	t.Run("Flush() works", func(t *testing.T) {
		testCache.Flush()

		count := len(testCache.(*timedCache).items)
		if count != 0 {
			t.Errorf("Cache was not cleared, length = %v", count)
		}
	})

	//
	// Test LastUpdated()
	//
	t.Run("LastUpdated() works", func(t *testing.T) {
		val := testCache.LastUpdated()
		actual := testCache.(*timedCache).updated

		if val != actual {
			t.Errorf("Mismatch, %v != %v", val, actual)
		}
	})

	//
	// Test Expired()
	//
	t.Run("Expired() works", func(t *testing.T) {
		val := testCache.Expired()

		// We should be expired by now.
		if val {
			t.Error("Expired!?")
		}
	})
}

func TestCallbacks(t *testing.T) {
	testCache = New(NewTestMetricConfig())

	// Manually set the eviction callback.
	testCache.OnEvicted(OnEvictedEvent)

	// Do some work here.
	testCache.Get(TestKey1)             // miss, get
	testCache.Set(TestKey1, TestValue1) // set
	testCache.Get(TestKey1)             // hit, get
	testCache.Delete(TestKey1)          // evict

	if hasEvicted == false {
		t.Error("Have not evicted!")
	}
}

// ** Benchmarks:

func BenchmarkTimedCache(b *testing.B) {
	const (
		TestKey = "TestKey"
		Val     = "This is a long value maybe."
		Replace = "This is a different key, yo."
	)

	var cache TimedCache

	b.Run("Constructor", func(b *testing.B) {
		b.ReportAllocs()

		cache = New(NewTestConfig())
	})

	b.Run("Add", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			cache.Add(TestKey, Val)
		}
	})

	b.Run("Replace", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			cache.Replace(TestKey, Replace)
		}
	})

	b.Run("Set", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			cache.Set(TestKey, Val)
		}
	})

	b.Run("Get", func(b *testing.B) {
		b.ReportAllocs()

		for range b.N {
			_, _ = cache.Get(TestKey)
		}
	})
}

// * cache_test.go ends here.
