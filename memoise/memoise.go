// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// memoise.go --- Memoisation hacks.
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
// mock:yes

// * Comments:

// * Package:

package memoise

// * Imports:

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/tozd/go/errors"
)

// * Variables:

var (
	//nolint:gochecknoglobals
	checkTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "memoise_check_total",
			Help: "Number of calls to the memoiser value checker"},
		[]string{"memoise", "result"})

	//nolint:gochecknoglobals
	loadDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "memoise_load_duration_seconds",
			Help:    "Callback runtime duration",
			Buckets: prometheus.DefBuckets},
		[]string{"memoise"})

	//nolint:gochecknoglobals
	inFlightGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memoise_inflight",
			Help: "Callbacks in flight"},
		[]string{"memoise"})

	//nolint:gochecknoglobals
	prometheusInitOnce sync.Once
)

// * Code:

// ** Types:

// Memoisation function type.
type CallbackFn func() (any, error)

// Memoisation type.
type Memoise interface {
	// Check returns the memoised value for the given key if available.
	// Otherwise it calls the provided callback to compute the value,
	// stores the result, and returns it.
	// Thread-safe.
	Check(string, CallbackFn) (any, error)

	// Clear the contents of the memoise map.
	Reset()
}

type entry struct {
	ready chan struct{}
	val   any
	err   error
}

// Implementation of the memoisation type.
type memoise struct {
	checkHit      prometheus.Counter
	checkMiss     prometheus.Counter
	checkErr      prometheus.Counter
	loadHist      prometheus.Observer
	inflightGauge prometheus.Gauge
	store         map[string]*entry
	name          string
	mu            sync.RWMutex
}

// ** Methods:

// Check returns the memoised value for the given key if available.
// Otherwise it calls the provided callback to compute the value,
// stores the result, and returns it.
// Thread-safe.
func (obj *memoise) Check(name string, callback CallbackFn) (any, error) {
	obj.mu.RLock()
	result := obj.store[name]
	obj.mu.RUnlock()

	if result != nil {
		<-result.ready

		return result.val, errors.WithStack(result.err)
	}

	// Miss, create a placeholder.
	result = &entry{ready: make(chan struct{})}

	obj.mu.Lock()
	// CRITICAL SECTION START.
	{
		// Re-check after acquiring write lock.
		if exist := obj.store[name]; exist != nil {
			obj.mu.Unlock() // Exit critical section/
			<-exist.ready

			return exist.val, errors.WithStack(exist.err)
		}

		obj.store[name] = result
	}
	// CRITICAL SECTION END.
	obj.mu.Unlock()

	val, err := callback()
	if err != nil {
		result.err = errors.WithStack(err)
		close(result.ready)

		obj.mu.Lock()
		// CRITICAL SECTION START.
		{
			delete(obj.store, name)
		}
		// CRITICAL SECTION END.
		obj.mu.Unlock()

		return nil, errors.WithStack(result.err)
	}

	result.val = val
	close(result.ready)

	return val, nil
}

func (obj *memoise) Reset() {
	obj.mu.Lock()
	// CRITICAL SECTION START.
	{
		obj.store = make(map[string]*entry)
	}
	// CRITICAL SECTION END.
	obj.mu.Unlock()
}

// ** Functions:

// Create a new memoisation object.
func NewMemoise(cfg *Config) Memoise {
	if cfg.Prometheus == nil {
		cfg.Prometheus = prometheus.DefaultRegisterer
	}

	InitPrometheus(cfg.Prometheus)

	if len(cfg.Name) == 0 {
		cfg.Name = "Default"
	}

	label := prometheus.Labels{"memoise": cfg.Name}
	curried, _ := checkTotal.CurryWith(label)
	hit := curried.WithLabelValues("hit")
	miss := curried.WithLabelValues("miss")
	errc := curried.WithLabelValues("error")

	return &memoise{
		name:          cfg.Name,
		store:         make(map[string]*entry),
		checkHit:      hit,
		checkMiss:     miss,
		checkErr:      errc,
		inflightGauge: inFlightGauge.With(label),
		loadHist:      loadDuration.With(label)}
}

// Initialise Prometheus metrics.
func InitPrometheus(reg prometheus.Registerer) {
	prometheusInitOnce.Do(func() {
		reg.MustRegister(
			checkTotal,
			loadDuration,
			inFlightGauge)
	})
}

// * memoise.go ends here.
