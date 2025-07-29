<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# memoise -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/memoise"
```

## Usage

```go
const (
	ContextKeyMemoise = "_DI_MEMO"
)
```

```go
var (
	ErrValueNotMemoise = errors.Base("value is not memoise.Memoise")
)
```

#### func  SetMemoiser

```go
func SetMemoiser(ctx context.Context, inst Memoise) (context.Context, error)
```
Set the memoiser value to the context map.

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

#### func  GetMemoiser

```go
func GetMemoiser(ctx context.Context) (Memoise, error)
```
Get the memoiser from the given context.

Will return `ErrValueNoMemoise` if the value in the context is not of type
`memoise.Memoise`.

#### func  MustGetMemoiser

```go
func MustGetMemoiser(ctx context.Context) Memoise
```
Attempt to get the memoiser from the given context. Panics if the operation
fails.

#### func  NewMemoise

```go
func NewMemoise() Memoise
```
Create a new memoisation object.
