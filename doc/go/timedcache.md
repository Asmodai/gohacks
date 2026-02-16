<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# timedcache -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/timedcache"
```

## Usage

```go
const ContextKeyTimedCache = "gohacks/timedcache@v1"
```
Key used to store the instance in the context's user value.

```go
const (
	DefaultCacheExpiration int = 1000
)
```

```go
var (
	// Triggered when an operation that expects a key to not exist find
	// that the key actually does exist.
	ErrKeyExists = errors.Base("the specified key already exists")

	// Triggered when an operation that expects a key to exist finds that
	// the key actually does not exist.
	ErrKeyNotExist = errors.Base("the specified key does not exist")
)
```

```go
var ErrValueNotTimedCache = errors.Base("value is not TimedCache")
```
Signalled if the instance associated with the context key is not of type
TimedCache.

#### func  InitPrometheus

```go
func InitPrometheus(reg prometheus.Registerer)
```
Initialise Prometheus metrics.

#### func  SetTimedCache

```go
func SetTimedCache(ctx context.Context, inst TimedCache) (context.Context, error)
```
Set TimedCache stores the instance in the context map.

#### func  SetTimedCacheIfAbsent

```go
func SetTimedCacheIfAbsent(ctx context.Context, inst TimedCache) (context.Context, error)
```
SetTimedCacheIfAbsent sets only if not already present.

#### func  WithTimedCache

```go
func WithTimedCache(ctx context.Context, fn func(TimedCache))
```
WithTimedCache calls fn with the instance or fallback.

#### type CacheItems

```go
type CacheItems map[any]Item
```

Type definition for the map of items in the cache.

#### type Config

```go
type Config struct {
	OnEvicted      OnEvictFn             `config_hide:"true"     json:"-"`
	ExpirationTime int                   `json:"expiration_time"`
	Name           string                `json:"-"`
	Prometheus     prometheus.Registerer `json:"-"`
}
```


#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Create a timed cache with a default configuration.

#### type Item

```go
type Item struct {
	Object any
}
```

Type definition for the internal cache item structure.

#### type MetricFn

```go
type MetricFn func()
```

Type definition for a metrics callback function.

#### type OnEvictFn

```go
type OnEvictFn func(any, any)
```

Type definition for the "On Eviction" callback function.

#### type TimedCache

```go
type TimedCache interface {
	// Sets the value for the given key to the given value.
	//
	// This uses Go map semantics, so if the given key doesn't exist in
	// the cache then one will be created.
	Set(any, any)

	// Gets the value for the given key.
	//
	// If the key exists, then the value and `true` will be returned;
	// otherwise `nil` and `false` will be returned.
	Get(any) (any, bool)

	// Adds the given key/value pair to the cache.
	//
	// This method expects the given key to not be present in the cache
	// and will return `ErrKeyExists` should it be present.
	Add(any, any) error

	// Replace the value for the given key with the given value.
	//
	// This method expects the given key to be present in the cache and
	// will return `ErrKeyNotExist` should it not be present.
	Replace(any, any) error

	// Delete the key/value pair from the cache.
	//
	// If the key exists, then its value and `true` will be returned;
	// otherwise `nil` and `false` will be returned.
	//
	// This method will attempt to invoke the "on eviction" callback.
	Delete(any) (any, bool)

	// Sets the "on eviction" callback to the given function.
	//
	// The function should take two arguments, the key and the value, of
	// type `any` and should not return a value.
	OnEvicted(OnEvictFn)

	// Return a count of the number of items in the cache.
	Count() int

	// Flush all items from the cache.
	Flush()

	// Return the time the cache was last updated.
	LastUpdated() time.Time

	//  Returns `true` if the cache has expired.
	Expired() bool

	// Return a list of all keys in the cache.
	Keys() []any
}
```


#### func  FromTimedCache

```go
func FromTimedCache(ctx context.Context) TimedCache
```
FromTimedCache returns the instance or the fallback.

#### func  GetTimedCache

```go
func GetTimedCache(ctx context.Context) (TimedCache, error)
```
Get the instance from the given context.

Will return ErrValueNotTimedCache if the value in the context is not of type
TimedCache.

#### func  MustGetTimedCache

```go
func MustGetTimedCache(ctx context.Context) TimedCache
```
Attempt to get the instance from the given context. Panics if the operation
fails.

#### func  New

```go
func New(config *Config) TimedCache
```
Create a new timed cache with the given configuration.

#### func  NewDefault

```go
func NewDefault() TimedCache
```

#### func  TryGetTimedCache

```go
func TryGetTimedCache(ctx context.Context) (TimedCache, bool)
```
TryGetTimedCache returns the instance and true if present and typed.
