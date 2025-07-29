<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# contextdi -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/contextdi"
```

## Usage

```go
const (
	ContextKeyMemoise = "_DI_MEMO"
)
```

```go
const (
	ContextKeyResponderChain = "_DI_RESPONDER"
)
```

```go
const (
	ContextKeyTimedCache = "_DI_TIMEDCACHE"
)
```

```go
var (
	ErrKeyNotFound = errors.Base("value map key not found")
)
```

```go
var (
	ErrValueNotMemoise = errors.Base("value is not memoise.Memoise")
)
```

```go
var (
	ErrValueNotResponderChain = errors.Base("value is not responder.Chain")
)
```

```go
var (
	ErrValueNotTimedCache = errors.Base("value is not timedcache.TimedCache")
)
```

#### func  GetFromContext

```go
func GetFromContext(ctx context.Context, key string) (any, error)
```
Get a value from a context.

Will signal `contextext.ErrInvalidContext` if the context is not valid. Will
signal `contextext.ErrValueMapNotFound` if there is no value map. Will signal
`ErrKeyNotFound` if the value map does not contain the key.

#### func  GetMemoise

```go
func GetMemoise(ctx context.Context) (memoise.Memoise, error)
```
Get the memoiser from the given context.

Will return `ErrValueNoMemoise` if the value in the context is not of type
`memoise.Memoise`.

#### func  GetResponderChain

```go
func GetResponderChain(ctx context.Context) (*responder.Chain, error)
```
Get the responder chain from the given context.

Will return `ErrValueNotResponderChain` if the value in the context is not of
type `responder.Chain`.

Please be aware that this responder chain should be treated as immutable, as we
can't really propagate changes down the context hierarchy.

#### func  GetTimedCache

```go
func GetTimedCache(ctx context.Context) (timedcache.TimedCache, error)
```
Get the timed cache value from the given context.

WIll return `ErrValueNotTimedCache` if the value in the context is not of type
`timedcache.TimedCache`.

#### func  MustGetMemoise

```go
func MustGetMemoise(ctx context.Context) memoise.Memoise
```
Attempt to get the memoiser from the given context. Panics if the operation
fails.

#### func  MustGetResponderChain

```go
func MustGetResponderChain(ctx context.Context) *responder.Chain
```
Attempt to get the responder chain from the given context. Panics if the
operation fails.

#### func  MustGetTimedCache

```go
func MustGetTimedCache(ctx context.Context) timedcache.TimedCache
```
Attempt to get the timed cache value from the given context. Panics if the
operation fails.

#### func  PutToContext

```go
func PutToContext(ctx context.Context, key string, value any) (context.Context, error)
```
Place a value in a context.

If there is no value map in the context then one will be created.

Returns a new context with the value map.

#### func  SetMemoise

```go
func SetMemoise(ctx context.Context, inst memoise.Memoise) (context.Context, error)
```
Set the memoiser value to the context map.

#### func  SetResponderChain

```go
func SetResponderChain(ctx context.Context, inst *responder.Chain) (context.Context, error)
```
Set the responder chain value in the context map.

#### func  SetTimedCache

```go
func SetTimedCache(ctx context.Context, inst timedcache.TimedCache) (context.Context, error)
```
Set the timed cache value in the context map.
