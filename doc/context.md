-*- Mode: gfm -*-

# context -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/context"
```

## Usage

#### func  WithValueMap

```go
func WithValueMap(ctx Context, valuemap ValueMap) Context
```
Create a context with the value map using a default key.

#### func  WithValueMapWithKey

```go
func WithValueMapWithKey(ctx Context, key string, valuemap ValueMap) Context
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

#### func  GetValueMap

```go
func GetValueMap(ctx Context) ValueMap
```
Get the value map (if any) from the context.

Returns nil if there is no value map.

#### func  GetValueMapWithKey

```go
func GetValueMapWithKey(ctx Context, key string) ValueMap
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
