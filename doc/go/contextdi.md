<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# contextdi -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/contextdi"
```

## Usage

```go
const (
	ContextKeyDebugMode string = "_DI_FLG_DEBUG"
)
```

```go
var (
	ErrKeyNotFound = errors.Base("value map key not found")
)
```

#### func  GetDebugMode

```go
func GetDebugMode(ctx context.Context) (bool, error)
```
Get the debug mode flag from the DI context.

#### func  GetFromContext

```go
func GetFromContext(ctx context.Context, key string) (any, error)
```
Get a value from a context.

Will signal `contextext.ErrInvalidContext` if the context is not valid. Will
signal `contextext.ErrValueMapNotFound` if there is no value map. Will signal
`ErrKeyNotFound` if the value map does not contain the key.

#### func  PutToContext

```go
func PutToContext(ctx context.Context, key string, value any) (context.Context, error)
```
Place a value in a context.

If there is no value map in the context then one will be created.

Returns a new context with the value map.

#### func  SetDebugMode

```go
func SetDebugMode(ctx context.Context, debugMode bool) (context.Context, error)
```
Set the debug mode flag in the DI context to the given value.
