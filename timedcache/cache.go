// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// cache.go --- Timed cache.
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
//
// mock:yes

package timedcache

import (
	"gitlab.com/tozd/go/errors"

	"fmt"
	"sync"
	"time"
)

var (
	// Triggered when an operation that expects a key to not exist find
	// that the key actually does exist.
	ErrKeyExists = errors.Base("the specified key already exists")

	// Triggered when an operation that expects a key to exist finds that
	// the key actually does not exist.
	ErrKeyNotExist = errors.Base("the specified key does not exist")
)

// Type definition for the "On Eviction" callback function.
type OnEvictFn func(any, any)

// Type definition for a metrics callback function.
type MetricFn func()

// Type definition for the internal cache item structure.
type Item struct {
	Object any
}

// Type definition for the map of items in the cache.
type CacheItems map[any]Item

type TimedCache interface {
	// Sets the value for the given key to the given value.
	//
	// This uses Go map semantics, so if the given key doesn't exist in
	// the cache then one will be created.
	Set(any, any)

	// Gets the value for the given key.
	//
	// If the key exists, then the value and `true` will be returned;
	// otherwise `nil` and `false` will be returned.
	Get(any) (any, bool)

	// Adds the given key/value pair to the cache.
	//
	// This method expects the given key to not be present in the cache
	// and will return `ErrKeyExists` should it be present.
	Add(any, any) error

	// Replace the value for the given key with the given value.
	//
	// This method expects the given key to be present in the cache and
	// will return `ErrKeyNotExist` should it not be present.
	Replace(any, any) error

	// Delete the key/value pair from the cache.
	//
	// If the key exists, then its value and `true` will be returned;
	// otherwise `nil` and `false` will be returned.
	//
	// This method will attempt to invoke the "on eviction" callback.
	Delete(any) (any, bool)

	// Sets the "on eviction" callback to the given function.
	//
	// The function should take two arguments, the key and the value, of
	// type `any` and should not return a value.
	OnEvicted(OnEvictFn)

	// Return a count of the number of items in the cache.
	Count() int

	// Flush all items from the cache.
	Flush()

	// Return the time the cache was last updated.
	LastUpdated() time.Time

	//  Returns `true` if the cache has expired.
	Expired() bool
}

type timedCache struct {
	updated         time.Time     // Last update time.
	expiration      time.Duration // Cache expiration time.
	items           CacheItems    // Cached items.
	mutex           sync.RWMutex  // R/W mutex.
	onEvicted       OnEvictFn     // Callback for eviction.
	cacheHitMetric  MetricFn      // Callback for hit metrics.
	cacheMissMetric MetricFn      // Callback for miss metrics.
	cacheGetMetric  MetricFn      // Callback for get metrics.
	cacheSetMetric  MetricFn      // Callbcak for set metrics.
}

// Create a new timed cache with the given configuration.
func New(config *Config) TimedCache {
	expire := time.Duration(config.ExpirationTime) * time.Second

	return &timedCache{
		expiration:      expire,
		onEvicted:       config.OnEvicted,
		items:           CacheItems{},
		cacheHitMetric:  config.CacheHitMetric,
		cacheMissMetric: config.CacheMissMetric,
		cacheGetMetric:  config.CacheGetMetric,
		cacheSetMetric:  config.CacheSetMetric,
	}
}

// Set the value for the given key.
func (obj *timedCache) Set(key any, value any) {
	obj.mutex.Lock()
	obj.set(key, value)
	obj.mutex.Unlock()
}

// Internal implementation for setting values.
func (obj *timedCache) set(key any, value any) {
	obj.tryCacheSetMetric()

	obj.updated = time.Now()
	obj.items[key] = Item{Object: value}
}

// Ge the value for the given key.
func (obj *timedCache) Get(key any) (any, bool) {
	obj.mutex.RLock()
	itm, found := obj.get(key)
	obj.mutex.RUnlock()

	return itm, found
}

// Internal implementation for getting values.
func (obj *timedCache) get(key any) (any, bool) {
	obj.tryCacheGetMetric()

	itm, found := obj.items[key]
	if !found {
		obj.tryCacheMissMetric()

		return nil, false
	}

	obj.tryCacheHitMetric()

	return itm.Object, true
}

// Add a new key/value pair to the cache.
//
// Triggers `ErrKeyExists` if the given key already exists.
func (obj *timedCache) Add(key any, value any) error {
	obj.mutex.Lock()

	if _, found := obj.get(key); found {
		obj.mutex.Unlock()

		return errors.Wrap(
			ErrKeyExists,
			fmt.Sprintf("key '%s' already exists", key),
		)
	}

	obj.set(key, value)
	obj.mutex.Unlock()

	return nil
}

// Replace the value for the given key.
//
// Triggers `ErrKeyNotExist` if the key does not exist.
func (obj *timedCache) Replace(key any, value any) error {
	obj.mutex.Lock()

	if _, found := obj.get(key); !found {
		obj.mutex.Unlock()

		return errors.Wrap(
			ErrKeyNotExist,
			fmt.Sprintf("key '%s' does not exist", key),
		)
	}

	obj.set(key, value)
	obj.mutex.Unlock()

	return nil
}

// Delete the given key from the cache.
func (obj *timedCache) Delete(key any) (any, bool) {
	obj.mutex.Lock()
	val, found := obj.zap(key)
	obj.mutex.Unlock()

	if found {
		obj.tryOnEvicted(key, val)
	}

	return val, found
}

// Internal implementation of key/value pair deletion.
func (obj *timedCache) zap(key any) (any, bool) {
	val, found := obj.items[key]

	if found {
		delete(obj.items, key)
		obj.tryOnEvicted(key, val)
	}

	return val, found
}

// Set the "on eviction" callback function.
func (obj *timedCache) OnEvicted(fn OnEvictFn) {
	obj.mutex.Lock()
	obj.onEvicted = fn
	obj.mutex.Unlock()
}

// Return a count of the elements in the cache.
func (obj *timedCache) Count() int {
	obj.mutex.Lock()
	itms := len(obj.items)
	obj.mutex.Unlock()

	return itms
}

// Flush all elements from the cache.
func (obj *timedCache) Flush() {
	obj.mutex.Lock()
	obj.updated = time.Now()
	obj.items = CacheItems{}
	obj.mutex.Unlock()
}

// Return the time of the last cache update.
func (obj *timedCache) LastUpdated() time.Time {
	obj.mutex.RLock()
	updated := obj.updated
	obj.mutex.RUnlock()

	return updated
}

// Has the cache expired?
func (obj *timedCache) Expired() bool {
	obj.mutex.RLock()
	end := obj.updated.Add(obj.expiration)
	obj.mutex.RUnlock()

	return time.Now().After(end)
}

// Attempt to invoke the eviction callback.
func (obj *timedCache) tryOnEvicted(key any, value any) {
	if obj.onEvicted == nil {
		return
	}

	obj.onEvicted(key, value)
}

// Attempt to invoke the cache set metric callback.
func (obj *timedCache) tryCacheSetMetric() {
	if obj.cacheSetMetric == nil {
		return
	}

	obj.cacheSetMetric()
}

// Attempt to invoke the cache get metric callback.
func (obj *timedCache) tryCacheGetMetric() {
	if obj.cacheGetMetric == nil {
		return
	}

	obj.cacheGetMetric()
}

// Attempt to invoke the cache hit metric callback.
func (obj *timedCache) tryCacheHitMetric() {
	if obj.cacheHitMetric == nil {
		return
	}

	obj.cacheHitMetric()
}

// Attempt to invoke the cache miss metric callback.
func (obj *timedCache) tryCacheMissMetric() {
	if obj.cacheMissMetric == nil {
		return
	}

	obj.cacheMissMetric()
}

// cache.go ends here.
