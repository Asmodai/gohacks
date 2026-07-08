// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// cache.go --- Timed cache.
//
// Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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
//
//mock:yes
//go:generate go run github.com/Asmodai/gohacks/cmd/digen -pattern .
//di:gen basename=TimedCache key=gohacks/timedcache@v1 type=TimedCache fallback=NewDefault()

// * Comments:

//
// TODO: Replace weird metrics callbacks with Prometheus.
// TODO: Add a config element to give the cache a name so that we can have
// multiple caches with Prometheus labels to track them.
//

// * Package:

package timedcache

// * Imports:

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	// Triggered when an operation that expects a key to not exist find
	// that the key actually does exist.
	ErrKeyExists = errors.Base("the specified key already exists")

	// Triggered when an operation that expects a key to exist finds that
	// the key actually does not exist.
	ErrKeyNotExist = errors.Base("the specified key does not exist")

	//nolint:gochecknoglobals
	itemsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "timedcache_items",
			Help: "Number of items in the cache"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	updatedGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "timedcache_updates",
			Help: "Last time cache was updated"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	getTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "timedcache_get_total",
			Help: "Number of cache gets"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	setTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "timedcache_set_total",
			Help: "Number of cache sets"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	hitTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "timedcache_hit_total",
			Help: "Number of cache hits"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	missTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "timedcache_miss_total",
			Help: "Number of cache misses"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	evictTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "timedcache_evict_total",
			Help: "Number of cache evictions"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	deleteTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "timedcache_delete_total",
			Help: "Number of cache deletions"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	flushTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "timedcache_flush_total",
			Help: "Number of cache flushes"},
		[]string{"timedcache"})

	//nolint:gochecknoglobals
	prometheusInitOnce sync.Once
)

// * Code:

// ** Interfaces:

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

	//  Returns `true` if the cached key has expired.
	//
	// The second return value will be `true` if the item was found.
	Expired(any) (bool, bool)

	// Return a list of all keys in the cache.
	Keys() []any
}

// ** Types:

// Type definition for the "On Eviction" callback function.
type OnEvictFn func(any, any)

// Type definition for a metrics callback function.
type MetricFn func()

// Type definition for the internal cache item structure.
type Item struct {
	Object    any
	ExpiresAt time.Time
}

// Type definition for the map of items in the cache.
type CacheItems map[any]Item

// Timed cache implementation.
type timedCache struct {
	updated            time.Time          // Last update time.
	cacheItemsMetric   prometheus.Gauge   // Item count metric.
	cacheUpdatedMetric prometheus.Gauge   // Update time metric.
	cacheHitMetric     prometheus.Counter // Cache hit metric.
	cacheMissMetric    prometheus.Counter // Cache miss metric.
	cacheGetMetric     prometheus.Counter // Cache get metric.
	cacheSetMetric     prometheus.Counter // Cache set metric.
	cacheEvictMetric   prometheus.Counter // Cache evict metric.
	cacheDeleteMetric  prometheus.Counter // Cache delete metric.
	cacheFlushMetric   prometheus.Counter // Cache flush metric.
	items              CacheItems         // Cached items.
	onEvicted          OnEvictFn          // Callback for eviction.
	name               string             // Name for the cache.
	expiration         time.Duration      // Cache expiration time.
	mutex              sync.RWMutex       // R/W mutex.
}

// ** Methods:

// Return a list of all keys in the cache.
func (tc *timedCache) Keys() []any {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	keys := make([]any, 0, len(tc.items))
	for k := range tc.items {
		keys = append(keys, k)
	}

	return keys
}

// Set the value for the given key.
func (tc *timedCache) Set(key any, value any) {
	var size int

	now := time.Now()

	tc.mutex.Lock()
	// CRITICAL SECTION START.
	{
		tc.updated = now
		tc.items[key] = Item{
			Object:    value,
			ExpiresAt: now.Add(tc.expiration),
		}
		size = len(tc.items)
	}
	// CRITICAL SECTION END.
	tc.mutex.Unlock()

	tc.cacheSetMetric.Inc()
	tc.cacheItemsMetric.Set(float64(size))
	tc.cacheUpdatedMetric.Set(float64(now.Unix()))
}

// Get the value for the given key.
func (tc *timedCache) Get(key any) (any, bool) {
	tc.mutex.RLock()
	itm, found := tc.items[key]
	tc.mutex.RUnlock()

	get := tc.cacheGetMetric
	get.Inc()

	if found && time.Now().Before(itm.ExpiresAt) {
		tc.cacheHitMetric.Inc()

		return itm.Object, true
	}

	if time.Now().After(itm.ExpiresAt) {
		tc.Delete(key)
	}

	tc.cacheMissMetric.Inc()

	return nil, false
}

// Add a new key/value pair to the cache.
//
// Triggers `ErrKeyExists` if the given key already exists.
func (tc *timedCache) Add(key any, value any) error {
	var (
		size int
		now  time.Time
	)

	tc.mutex.Lock()
	// CRITICAL SECTION START.
	{
		if _, exists := tc.items[key]; exists {
			tc.mutex.Unlock() // Exit critical section here.

			return errors.WithMessagef(
				ErrKeyExists,
				"key %q already exists",
				key)
		}

		now = time.Now()
		tc.items[key] = Item{
			Object:    value,
			ExpiresAt: now.Add(tc.expiration),
		}
		tc.updated = now
		size = len(tc.items)
	}
	// CRITICAL SECTION END.
	tc.mutex.Unlock()

	tc.cacheSetMetric.Inc()
	tc.cacheItemsMetric.Set(float64(size))
	tc.cacheUpdatedMetric.Set(float64(now.Unix()))

	return nil
}

// Replace the value for the given key.
//
// Triggers `ErrKeyNotExist` if the key does not exist.
func (tc *timedCache) Replace(key any, value any) error {
	var now time.Time

	tc.mutex.Lock()
	// CRITICAL SECTION START.
	{
		if _, exists := tc.items[key]; !exists {
			tc.mutex.Unlock() // Exit critical section here.

			return errors.WithMessagef(
				ErrKeyNotExist,
				"key %q does not exist",
				key)
		}

		now = time.Now()
		tc.items[key] = Item{
			Object:    value,
			ExpiresAt: now.Add(tc.expiration),
		}
		tc.updated = now
	}
	// CRITICAL SECTION END.
	tc.mutex.Unlock()

	tc.cacheSetMetric.Inc()
	tc.cacheUpdatedMetric.Set(float64(now.Unix()))

	return nil
}

// Delete the given key from the cache.
func (tc *timedCache) Delete(key any) (any, bool) {
	var (
		val      any
		canEvict bool
		evict    OnEvictFn
		size     int
	)

	tc.mutex.Lock()
	// CRITICAL SECTION START.
	{
		if itm, found := tc.items[key]; found {
			delete(tc.items, key)

			val = itm.Object
			canEvict = true
			evict = tc.onEvicted
		}

		size = len(tc.items)
	}
	// CRITICAL SECTION END.
	tc.mutex.Unlock()

	if canEvict && evict != nil {
		tc.cacheEvictMetric.Inc()
		tc.cacheDeleteMetric.Inc()
		tc.cacheUpdatedMetric.Set(float64(time.Now().Unix()))

		evict(key, val)
	}

	tc.cacheItemsMetric.Set(float64(size))

	return val, canEvict
}

// Set the "on eviction" callback function.
func (tc *timedCache) OnEvicted(fn OnEvictFn) {
	tc.mutex.Lock()
	tc.onEvicted = fn
	tc.mutex.Unlock()
}

// Return a count of the elements in the cache.
func (tc *timedCache) Count() int {
	tc.mutex.RLock()
	itms := len(tc.items)
	tc.mutex.RUnlock()

	return itms
}

// Flush all elements from the cache.
func (tc *timedCache) Flush() {
	var (
		items CacheItems
		evict OnEvictFn
		now   time.Time
	)

	tc.mutex.Lock()
	// CRITICAL SECTION START.
	{
		items = tc.items
		evict = tc.onEvicted
		now = time.Now()

		tc.updated = now
		tc.items = CacheItems{}

		// Update the items count metric inside the lock.
		tc.cacheItemsMetric.Set(float64(0))
	}
	// CRITICAL SECTION END.
	tc.mutex.Unlock()

	if evict != nil {
		for k, v := range items {
			tc.cacheEvictMetric.Inc()

			evict(k, v.Object)
		}
	}

	tc.cacheFlushMetric.Inc()
	tc.cacheUpdatedMetric.Set(float64(now.Unix()))
}

// Return the time of the last cache update.
func (tc *timedCache) LastUpdated() time.Time {
	tc.mutex.RLock()
	updated := tc.updated
	tc.mutex.RUnlock()

	return updated
}

// Has the cached item expired?
func (tc *timedCache) Expired(key any) (bool, bool) {
	tc.mutex.RLock()
	itm, found := tc.items[key]
	tc.mutex.RUnlock()

	// If not found, then bomb out.
	if !found {
		return false, false
	}

	return time.Now().After(itm.ExpiresAt), true
}

// * Functions:

func NewDefault() TimedCache {
	return New(&Config{})
}

// Create a new timed cache with the given configuration.
func New(config *Config) TimedCache {
	if config.Prometheus == nil {
		config.Prometheus = prometheus.DefaultRegisterer
	}

	InitPrometheus(config.Prometheus)

	if len(config.Name) == 0 {
		config.Name = "Default"
	}

	label := prometheus.Labels{"timedcache": config.Name}
	expire := time.Duration(config.ExpirationTime) * time.Second

	return &timedCache{
		name:               config.Name,
		updated:            time.Now(),
		expiration:         expire,
		onEvicted:          config.OnEvicted,
		items:              CacheItems{},
		cacheItemsMetric:   itemsGauge.With(label),
		cacheUpdatedMetric: updatedGauge.With(label),
		cacheGetMetric:     getTotal.With(label),
		cacheSetMetric:     setTotal.With(label),
		cacheHitMetric:     hitTotal.With(label),
		cacheMissMetric:    missTotal.With(label),
		cacheEvictMetric:   evictTotal.With(label),
		cacheDeleteMetric:  deleteTotal.With(label),
		cacheFlushMetric:   flushTotal.With(label),
	}
}

// Initialise Prometheus metrics.
func InitPrometheus(reg prometheus.Registerer) {
	prometheusInitOnce.Do(func() {
		reg.MustRegister(
			itemsGauge,
			updatedGauge,
			getTotal,
			setTotal,
			hitTotal,
			missTotal,
			evictTotal,
			deleteTotal,
			flushTotal)
	})
}

// * cache.go ends here.
