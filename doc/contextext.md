<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# contextext -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/contextext"
```

## Usage

```go
var (
	ErrInvalidContext   = errors.Base("invalid context")
	ErrInvalidValueMap  = errors.Base("invalid value map")
	ErrValueMapNotFound = errors.Base("value map not found")
)
```

#### func  WithValueMap

```go
func WithValueMap(ctx context.Context, valuemap ValueMap) context.Context
```
Create a context with the value map using a default key.

#### func  WithValueMapWithKey

```go
func WithValueMapWithKey(ctx context.Context, key string, valuemap ValueMap) context.Context
```
Create a context with the value map using the specified key.

#### type ValueMap

```go
type ValueMap interface {
	Get(string) (key any, ok bool)
	Set(key string, value any)
}
```

A map-based storage structure to pass multiple values via contexts rather than
many invocations of `context.WithValue` and their respective copy operations.

The main caveat with this approach is that as contexts are copied by the various
`With` functions we have no means of passing changes to child contexts once the
context with the value map is copied.

This is not the main aim of this type, so such functionality should not be
considered. The main usage is to provide a means of passing a lot of values to
some top-level context in order to avoid a lot of `WithValue` calls and a
somewhat slow lookup.

#### func  GetValueMap

```go
func GetValueMap(ctx context.Context) (ValueMap, error)
```
Get the value map (if any) from the context.

Returns nil if there is no value map.

#### func  GetValueMapWithKey

```go
func GetValueMapWithKey(ctx context.Context, key string) (ValueMap, error)
```
Get the value map (if any) from the context with the specified value key.

#### func  NewValueMap

```go
func NewValueMap() ValueMap
```
Create a new value map with no data.

#### type ValueMapKey

```go
type ValueMapKey string
```

ValueMap key type for `WithValue`.
