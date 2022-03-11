/*
 * service.go --- Cheap nasty dependency injection.
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
	"sync"
)

// Constructor function for creating new service records.
type ServiceCtorFn func() interface{}

/*

Service structure.

To use:

1) Invoke `di.GetInstance` to access the singleton:
```go
  svc := di.GetInstance
```

2a) Add your required service:
```go
   svc.Add("SomeName", someInstance)
```

2b) Create your required service:
```go
   svc.Create("SomeName", func() interface{} { return NewThing() })
```

*/
type Service struct {
	sync.RWMutex

	services map[string]interface{}
	classes  map[string]ServiceCtorFn
}

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
			classes:  make(map[string]ServiceCtorFn),
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
func (s *Service) AddClass(name string, ctor ServiceCtorFn) {
	s.Lock()
	s.classes[name] = ctor
	s.Unlock()
}

// Get a service with the given name.
func (s *Service) Get(name string) (interface{}, bool) {
	s.Lock()
	t := s.services[name]
	s.Unlock()

	if t == nil {
		return nil, false
	}

	return t, true
}

// Create a new instance of the given class by invoking its registered
// constructor.
func (s *Service) CreateNew(name string) (interface{}, bool) {
	s.Lock()
	c := s.classes[name]
	s.Unlock()

	if c == nil {
		return nil, false
	}

	return c(), true
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

/* service.go ends here. */
