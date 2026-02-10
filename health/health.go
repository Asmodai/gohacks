// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// health.go --- Health structure.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

// * Comments:

// * Package:

package health

// * Imports:

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Asmodai/gohacks/errx"
)

// * Constants:

const (
	// Default health timeout in minutes.
	DefaultHealthTimeoutMinutes int64 = 25
)

// * Code:

// ** Types:

// Structure for health time type constraint.
type heartbeatTime struct {
	t time.Time
}

// Structure used for marshalling health to JSON.
type healthMarshal struct {
	Healthy   bool           `json:"is_healthy"`
	Heartbeat time.Time      `json:"last_heartbeat,omitempty"`
	UserData  map[string]any `json:"userdata,omitempty"`
}

/*
Health structure.

Here we provide a means of signalling health for various services.

Our health system is simply this:

1) Process invokes `Tick` to update a heartbeat timestamp,
2) `Healthy` is used to determine health, and
3) `LastHeartbeat` returns the time of the last heartbeat.
*/
type Health struct {
	timeout   time.Duration  // Timeout until transition to unhealthy.
	heartbeat atomic.Value   // Current heartbeat.  Must be monotonic.
	udmutex   sync.RWMutex   // Mutex for user data.
	userdata  map[string]any // User data.
}

// ** Methods:

// Are we healthy?
//
// We are considered healthy if the amount of time since the last heartbeat
// is within the timeout.
func (h *Health) Healthy() bool {
	last := h.heartbeat.Load()
	if last == nil {
		return false
	}

	//nolint:forcetypeassert
	tstamp := last.(heartbeatTime).t

	return time.Since(tstamp) <= h.timeout
}

// Return the timestamp of the last heartbeat.
//
// This is wall-clock time suitable for biologicals.  Do not use it for logic.
func (h *Health) LastHeartbeat() time.Time {
	last := h.heartbeat.Load()
	if last == nil {
		return time.Time{}
	}

	//nolint:forcetypeassert
	return last.(heartbeatTime).t
}

// Store current timestamp as the heartbeat value.
//
// The stored time includes a monotonic component.
// This method is atomic.
func (h *Health) Tick() {
	h.heartbeat.Store(heartbeatTime{t: time.Now()})
}

// Get the value for a given key from the user data.
func (h *Health) UserGet(key string) (any, bool) {
	var (
		val   any
		found bool
	)

	h.udmutex.RLock()
	{
		val, found = h.userdata[key]
	}
	h.udmutex.RUnlock()

	return val, found
}

// Set the value for the given key in the user data.
func (h *Health) UserSet(key string, value any) {
	h.udmutex.Lock()
	{
		if h.userdata == nil {
			h.userdata = make(map[string]any)
		}

		h.userdata[key] = value
	}
	h.udmutex.Unlock()
}

// Encode the health object as JSON.
func (h *Health) MarshalJSON() ([]byte, error) {
	if h == nil {
		return []byte{}, nil
	}

	// Snapshot userdata safely.
	var udata map[string]any

	h.udmutex.RLock()
	{
		if len(h.userdata) > 0 {
			udata = make(map[string]any, len(h.userdata))

			for key, val := range h.userdata {
				udata[key] = val
			}
		}
	}
	h.udmutex.RUnlock()

	tmp := &healthMarshal{
		Healthy:   h.Healthy(),
		Heartbeat: h.LastHeartbeat(),
		UserData:  udata,
	}

	result, err := json.Marshal(tmp)

	return result, errx.WithStack(err)
}

// ** Functions:

// Create a new health instance with the timeout set to the given duration.
func NewHealthWithDuration(duration time.Duration) *Health {
	inst := &Health{
		timeout:  duration,
		userdata: make(map[string]any),
	}

	// Perform an initial tick.
	inst.Tick()

	return inst
}

// Create a new health instance with the timeout set to the given minutes.
//
// The argument here is minutes, and is converted to a duration in minutes.
// This is possibly not the method you want to use.
func NewHealth(timeoutMinutes int64) *Health {
	return NewHealthWithDuration(time.Duration(timeoutMinutes) * time.Minute)
}

// Create a new health instance with the default timeout value.
func NewDefaultHealth() *Health {
	return NewHealth(DefaultHealthTimeoutMinutes)
}

// * health.go ends here.
