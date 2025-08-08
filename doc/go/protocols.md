<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# protocols -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/protocols"
```

## Usage

```go
var (
	ErrNoVerifierFunction = errors.Base("no verifier function")
)
```

#### type Protocol

```go
type Protocol struct {
	Name      string
	Selectors []string
}
```


#### type Registry

```go
type Registry struct {
}
```


#### func  NewRegistry

```go
func NewRegistry() *Registry
```

#### func (*Registry) Register

```go
func (r *Registry) Register(proto *Protocol)
```

#### func (*Registry) RegisterWithVerifier

```go
func (r *Registry) RegisterWithVerifier(proto *Protocol, verifier Verifier)
```

#### func (*Registry) Validate

```go
func (r *Registry) Validate(name string, rbl hasMethodsIntrospector) bool
```

#### func (*Registry) Verify

```go
func (r *Registry) Verify(name string, obj selector.Introspectable) error
```

#### type Verifier

```go
type Verifier func(selector.Introspectable) error
```
