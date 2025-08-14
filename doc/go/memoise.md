<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# memoise -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/memoise"
```

## Usage

```go
const ContextKeyMemoise = "gohacks/memoise@v1"
```
Key used to store the instance in the context's user value.

```go
var ErrValueNotMemoise = errors.Base("value is not Memoise")
```
Signalled if the instance associated with the context key is not of type
Memoise.

#### func  InitPrometheus

```go
func InitPrometheus(reg prometheus.Registerer)
```
Initialise Prometheus metrics.

#### func  SetMemoise

```go
func SetMemoise(ctx context.Context, inst Memoise) (context.Context, error)
```
Set Memoise stores the instance in the context map.

#### func  SetMemoiseIfAbsent

```go
func SetMemoiseIfAbsent(ctx context.Context, inst Memoise) (context.Context, error)
```
SetMemoiseIfAbsent sets only if not already present.

#### func  WithMemoise

```go
func WithMemoise(ctx context.Context, fn func(Memoise))
```
WithMemoise calls fn with the instance or fallback.

#### type CallbackFn

```go
type CallbackFn func() (any, error)
```

Memoisation function type.

#### type Config

```go
type Config struct {
	// Prometheus registerer.
	Prometheus prometheus.Registerer `json:"-"`

	// Instance name.
	//
	// This gets used for Prometheus metrics should you have more than
	// one memoiser in your application.
	Name string `json:"-"`
}
```


#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```

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

#### func  FromMemoise

```go
func FromMemoise(ctx context.Context) Memoise
```
FromMemoise returns the instance or the fallback.

#### func  GetMemoise

```go
func GetMemoise(ctx context.Context) (Memoise, error)
```
Get the logger from the given context.

Will return ErrValueNotMemoise if the value in the context is not of type
Memoise.

#### func  MustGetMemoise

```go
func MustGetMemoise(ctx context.Context) Memoise
```
Attempt to get the instance from the given context. Panics if the operation
fails.

#### func  NewDefaultMemoise

```go
func NewDefaultMemoise() Memoise
```

#### func  NewMemoise

```go
func NewMemoise(cfg *Config) Memoise
```
Create a new memoisation object.

#### func  TryGetMemoise

```go
func TryGetMemoise(ctx context.Context) (Memoise, bool)
```
TryGetMemoise returns the instance and true if present and typed.
