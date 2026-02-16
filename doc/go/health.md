<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# health -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/health"
```

## Usage

```go
const ContextKeyHealth = "gohacks/health@v1"
```
Key used to store the instance in the context's user value.

```go
const (
	// Default health timeout in minutes.
	DefaultHealthTimeoutMinutes int64 = 25
)
```

```go
var ErrValueNotHealth = errors.Base("value is not Reporter")
```
Signalled if the instance associated with the context key is not of type
Reporter.

#### func  SetHealth

```go
func SetHealth(ctx context.Context, inst Reporter) (context.Context, error)
```
Set Health stores the instance in the context map.

#### func  SetHealthIfAbsent

```go
func SetHealthIfAbsent(ctx context.Context, inst Reporter) (context.Context, error)
```
SetHealthIfAbsent sets only if not already present.

#### func  WithHealth

```go
func WithHealth(ctx context.Context, fn func(Reporter))
```
WithHealth calls fn with the instance or fallback.

#### type Health

```go
type Health struct {
}
```

Health structure.

Here we provide a means of signalling health for various services.

Our health system is simply this:

1) Process invokes `Tick` to update a heartbeat timestamp, 2) `Healthy` is used
to determine health, and 3) `LastHeartbeat` returns the time of the last
heartbeat.

#### func  NewDefaultHealth

```go
func NewDefaultHealth() *Health
```
Create a new health instance with the default timeout value.

#### func  NewHealth

```go
func NewHealth(timeoutMinutes int64) *Health
```
Create a new health instance with the timeout set to the given minutes.

The argument here is minutes, and is converted to a duration in minutes. This is
possibly not the method you want to use.

#### func  NewHealthWithDuration

```go
func NewHealthWithDuration(duration time.Duration) *Health
```
Create a new health instance with the timeout set to the given duration.

#### func (*Health) Healthy

```go
func (h *Health) Healthy() bool
```
Are we healthy?

We are considered healthy if the amount of time since the last heartbeat is
within the timeout.

#### func (*Health) LastHeartbeat

```go
func (h *Health) LastHeartbeat() time.Time
```
Return the timestamp of the last heartbeat.

This is wall-clock time suitable for biologicals. Do not use it for logic.

#### func (*Health) MarshalJSON

```go
func (h *Health) MarshalJSON() ([]byte, error)
```
Encode the health object as JSON.

#### func (*Health) Tick

```go
func (h *Health) Tick()
```
Store current timestamp as the heartbeat value.

The stored time includes a monotonic component. This method is atomic.

#### func (*Health) UserGet

```go
func (h *Health) UserGet(key string) (any, bool)
```
Get the value for a given key from the user data.

#### func (*Health) UserSet

```go
func (h *Health) UserSet(key string, value any)
```
Set the value for the given key in the user data.

#### type Reporter

```go
type Reporter interface {
	Healthy() bool
	LastHeartbeat() time.Time
	Tick()
	UserGet(string) (any, bool)
	UserSet(string, any)
}
```

Health reporter interface type.

Any object that conforms to this interface may be used to report health.

#### func  FromHealth

```go
func FromHealth(ctx context.Context) Reporter
```
FromHealth returns the instance or the fallback.

#### func  GetHealth

```go
func GetHealth(ctx context.Context) (Reporter, error)
```
Get the instance from the given context.

Will return ErrValueNotHealth if the value in the context is not of type
Reporter.

#### func  MustGetHealth

```go
func MustGetHealth(ctx context.Context) Reporter
```
Attempt to get the instance from the given context. Panics if the operation
fails.

#### func  TryGetHealth

```go
func TryGetHealth(ctx context.Context) (Reporter, bool)
```
TryGetHealth returns the instance and true if present and typed.

#### type Ticker

```go
type Ticker interface {
	Channel() <-chan time.Time
	Stop()
}
```


#### func  NewTicker

```go
func NewTicker(duration time.Duration) Ticker
```
