// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// service.go --- Cheap nasty service management.
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

package service

import (
	"sync"
)

// Constructor function for creating new service records.
type ConstructorFn func() interface{}

/*
Service structure.

To use:

1) Invoke `service.GetInstance` to access the singleton:

```go

	svc := service.GetInstance

```

2a) Add your required service:

```go

	svc.Add("SomeName", someInstance)

```

2b) Create your required service:

```go

	svc.Create("SomeName", func() interface{} { return NewThing() })

```

Profit.
*/
type Service struct {
	sync.RWMutex

	services map[string]interface{}
	classes  map[string]ConstructorFn
}

//nolint:gochecknoglobals
var (
	once     sync.Once
	instance *Service
)

// Debugging aid -- do *not* use.
func DumpInstance() *Service {
	return instance
}

// Return the service manager's singleton instance.
func GetInstance() *Service {
	once.Do(func() {
		instance = &Service{
			services: make(map[string]interface{}),
			classes:  make(map[string]ConstructorFn),
		}
	})

	return instance
}

// Add a new service instance with the given name.
func (s *Service) Add(name string, thing interface{}) {
	s.Lock()
	s.services[name] = thing
	s.Unlock()
}

// Add a new class with the given name.
func (s *Service) AddClass(name string, ctor ConstructorFn) {
	s.Lock()
	s.classes[name] = ctor
	s.Unlock()
}

// Get a service with the given name.
func (s *Service) Get(name string) (interface{}, bool) {
	s.Lock()
	thing := s.services[name]
	s.Unlock()

	if thing == nil {
		return nil, false
	}

	return thing, true
}

// Create a new instance of the given class by invoking its registered
// constructor.
func (s *Service) CreateNew(name string) (interface{}, bool) {
	s.Lock()
	ctor := s.classes[name]
	s.Unlock()

	if ctor == nil {
		return nil, false
	}

	return ctor(), true
}

// Get a list of registered services.
func (s *Service) Services() []string {
	s.Lock()
	keys := make([]string, len(s.services))
	i := 0

	for k := range s.services {
		keys[i] = k
		i++
	}

	s.Unlock()

	return keys
}

// Get a list of registered classes.
func (s *Service) Classes() []string {
	s.Lock()
	keys := make([]string, len(s.classes))
	i := 0

	for k := range s.classes {
		keys[i] = k
		i++
	}

	s.Unlock()

	return keys
}

// Return a count of registered services.
func (s *Service) CountServices() int {
	s.Lock()
	c := len(s.services)
	s.Unlock()

	return c
}

// Return a count of registered classes.
func (s *Service) CountClasses() int {
	s.Lock()
	c := len(s.classes)
	s.Unlock()

	return c
}

// service.go ends here.
