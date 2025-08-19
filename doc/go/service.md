<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# service -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/service"
```

## Usage

#### type ConstructorFn

```go
type ConstructorFn func() any
```

Constructor function for creating new service records.

#### type Service

```go
type Service struct {
	sync.RWMutex
}
```

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

    svc.Create("SomeName", func() any { return NewThing() })

```

Profit.

#### func  DumpInstance

```go
func DumpInstance() *Service
```
Debugging aid -- do *not* use.

#### func  GetInstance

```go
func GetInstance() *Service
```
Return the service manager's singleton instance.

#### func (*Service) Add

```go
func (s *Service) Add(name string, thing any)
```
Add a new service instance with the given name.

#### func (*Service) AddClass

```go
func (s *Service) AddClass(name string, ctor ConstructorFn)
```
Add a new class with the given name.

#### func (*Service) Classes

```go
func (s *Service) Classes() []string
```
Get a list of registered classes.

#### func (*Service) CountClasses

```go
func (s *Service) CountClasses() int
```
Return a count of registered classes.

#### func (*Service) CountServices

```go
func (s *Service) CountServices() int
```
Return a count of registered services.

#### func (*Service) CreateNew

```go
func (s *Service) CreateNew(name string) (any, bool)
```
Create a new instance of the given class by invoking its registered constructor.

#### func (*Service) Get

```go
func (s *Service) Get(name string) (any, bool)
```
Get a service with the given name.

#### func (*Service) Services

```go
func (s *Service) Services() []string
```
Get a list of registered services.
