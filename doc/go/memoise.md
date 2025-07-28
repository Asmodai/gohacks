<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

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
	// Check returns the memoised value for the given key if available.
	// Otherwise it calls the provided callback to compute the value,
	// stores the result, and returns it.
	// Thread-safe.
	Check(string, CallbackFn) (any, error)

	// Clear the contents of the memoise map.
	Reset()
}
```

Memoisation type.

#### func  NewMemoise

```go
func NewMemoise() Memoise
```
Create a new memoisation object.
