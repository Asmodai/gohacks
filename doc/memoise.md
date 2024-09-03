-*- Mode: gfm -*-

# memoise -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/memoise"
```

## Usage

#### type CallbackFn

```go
type CallbackFn func() (any, error)
```

Memoisation function type.

#### type Memoise

```go
type Memoise interface {
	// Check if we have a memorised value for a given key.  If not, then
	// inovke the callback function and memorise its result.
	Check(string, CallbackFn) (any, error)
}
```

Memoisation type.

#### func  NewMemoise

```go
func NewMemoise() Memoise
```
Create a new memoisation object.
