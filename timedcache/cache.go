// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// cache.go --- Timed cache.
//
// Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
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

	//  Returns `true` if the cache has expired.
	Expired() bool

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
	Object any
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
func (obj *timedCache) Keys() []any {
	obj.mutex.RLock()
	defer obj.mutex.RUnlock()

	keys := make([]any, 0, len(obj.items))
	for k := range obj.items {
		keys = append(keys, k)
	}

	return keys
}

// Set the value for the given key.
func (obj *timedCache) Set(key any, value any) {
	var size int

	now := time.Now()

	obj.mutex.Lock()
	// CRITICAL SECTION START.
	{
		obj.updated = now
		obj.items[key] = Item{Object: value}
		size = len(obj.items)
	}
	// CRITICAL SECTION END.
	obj.mutex.Unlock()

	obj.cacheSetMetric.Inc()
	obj.cacheItemsMetric.Set(float64(size))
	obj.cacheUpdatedMetric.Set(float64(now.Unix()))
}

// Ge the value for the given key.
func (obj *timedCache) Get(key any) (any, bool) {
	obj.mutex.RLock()
	itm, found := obj.items[key]
	obj.mutex.RUnlock()

	get := obj.cacheGetMetric

	if found {
		get.Inc()
		obj.cacheHitMetric.Inc()

		return itm.Object, true
	}

	get.Inc()
	obj.cacheMissMetric.Inc()

	return itm, found
}

// Add a new key/value pair to the cache.
//
// Triggers `ErrKeyExists` if the given key already exists.
func (obj *timedCache) Add(key any, value any) error {
	var (
		size int
		now  time.Time
	)

	obj.mutex.Lock()
	// CRITICAL SECTION START.
	{
		if _, exists := obj.items[key]; exists {
			obj.mutex.Unlock() // Exit critical section here.

			return errors.WithMessagef(
				ErrKeyExists,
				"key %q already exists",
				key)
		}

		now = time.Now()
		obj.items[key] = Item{Object: value}
		obj.updated = now
		size = len(obj.items)
	}
	// CRITICAL SECTION END.
	obj.mutex.Unlock()

	obj.cacheSetMetric.Inc()
	obj.cacheItemsMetric.Set(float64(size))
	obj.cacheUpdatedMetric.Set(float64(now.Unix()))

	return nil
}

// Replace the value for the given key.
//
// Triggers `ErrKeyNotExist` if the key does not exist.
func (obj *timedCache) Replace(key any, value any) error {
	var now time.Time

	obj.mutex.Lock()
	// CRITICAL SECTION START.
	{
		if _, exists := obj.items[key]; !exists {
			obj.mutex.Unlock() // Exit critical section here.

			return errors.WithMessagef(
				ErrKeyNotExist,
				"key %q does not exist",
				key)
		}

		now = time.Now()
		obj.items[key] = Item{Object: value}
		obj.updated = now
	}
	// CRITICAL SECTION END.
	obj.mutex.Unlock()

	obj.cacheSetMetric.Inc()
	obj.cacheUpdatedMetric.Set(float64(now.Unix()))

	return nil
}

// Delete the given key from the cache.
func (obj *timedCache) Delete(key any) (any, bool) {
	var (
		val      any
		canEvict bool
		evict    OnEvictFn
	)

	obj.mutex.Lock()
	// CRITICAL SECTION START.
	{
		if itm, found := obj.items[key]; found {
			delete(obj.items, key)

			val = itm.Object
			canEvict = true
			evict = obj.onEvicted
		}
	}
	// CRITICAL SECTION END.
	obj.mutex.Unlock()

	if canEvict && evict != nil {
		obj.cacheEvictMetric.Inc()
		obj.cacheDeleteMetric.Inc()
		obj.cacheUpdatedMetric.Set(float64(time.Now().Unix()))

		evict(key, val)
	}

	return val, canEvict
}

// Set the "on eviction" callback function.
func (obj *timedCache) OnEvicted(fn OnEvictFn) {
	obj.mutex.Lock()
	obj.onEvicted = fn
	obj.mutex.Unlock()
}

// Return a count of the elements in the cache.
func (obj *timedCache) Count() int {
	obj.mutex.RLock()
	itms := len(obj.items)
	obj.mutex.RUnlock()

	return itms
}

// Flush all elements from the cache.
func (obj *timedCache) Flush() {
	var (
		items CacheItems
		evict OnEvictFn
		now   time.Time
	)

	obj.mutex.Lock()
	// CRITICAL SECTION START.
	{
		items = obj.items
		evict = obj.onEvicted
		now = time.Now()

		obj.updated = now
		obj.items = CacheItems{}
	}
	// CRITICAL SECTION END.
	obj.mutex.Unlock()

	if evict != nil {
		for k, v := range items {
			obj.cacheEvictMetric.Inc()

			evict(k, v.Object)
		}
	}

	obj.cacheFlushMetric.Inc()
	obj.cacheUpdatedMetric.Set(float64(now.Unix()))
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
