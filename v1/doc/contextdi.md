<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# contextdi -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/contextdi"
```

## Usage

```go
const (
	ContextKeyDBManager = "_DI_DB_MGR"
)
```

```go
const (
	ContextKeyLogger = "_DI_LOGGER"
)
```

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
	ErrValueNotDBManager = errors.Base("value is not database.Manager")
)
```

```go
var (
	ErrValueNotLogger = errors.Base("value is not logger.Logger")
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

#### func  GetDBManager

```go
func GetDBManager(ctx context.Context) (database.Manager, error)
```
Get the database manager from the given context.

Will return `ErrValueNoDBManager` if the value in the context is not of type
`database.Manager`.

#### func  GetFromContext

```go
func GetFromContext(ctx context.Context, key string) (any, error)
```
Get a value from a context.

Will signal `contextext.ErrInvalidContext` if the context is not valid. Will
signal `contextext.ErrValueMapNotFound` if there is no value map. Will signal
`ErrKeyNotFound` if the value map does not contain the key.

#### func  GetLogger

```go
func GetLogger(ctx context.Context) (logger.Logger, error)
```
Get the logger from the given context.

Will return `ErrValueNotLogger` if the value in the context is not of type
`logger.Logger`.

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

#### func  MustGetDBManager

```go
func MustGetDBManager(ctx context.Context) database.Manager
```
Attempt to get the database manager from the given context. Panics if the
operation fails.

#### func  MustGetLogger

```go
func MustGetLogger(ctx context.Context) logger.Logger
```
Attempt to get the logger from the given context. Panics if the operation fails.

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

#### func  SetDBManager

```go
func SetDBManager(ctx context.Context, inst database.Manager) (context.Context, error)
```
Set the database manager value to the context map.

#### func  SetLogger

```go
func SetLogger(ctx context.Context, inst logger.Logger) (context.Context, error)
```
Set the logger value to the context map.

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
