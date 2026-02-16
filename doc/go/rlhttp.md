<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# rlhttp -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/rlhttp"
```

## Usage

```go
const (
	// A sane default client timeout value.
	//
	// This is in line with Go, as well as with services like Kubernetes,
	// AWS, et al.
	DefaultClientTimeout = types.Duration(30 * time.Second)

	// Default burst value.
	DefaultBurst int = 1
)
```

```go
const ContextKeyRLHTTP = "gohacks/rlhttp@v1"
```
Key used to store the instance in the context's user value.

```go
var (
	ErrInvalidSettings = errx.Base("invalid rate limiter settings")
	ErrInvalidLimiter  = errx.Base("invalid rate limit")
)
```

```go
var (
	ErrNilRequest = errx.Base("nil request")
)
```

```go
var ErrValueNotRLHTTP = errors.Base("value is not *Client")
```
Signalled if the instance associated with the context key is not of type
*Client.

#### func  SetRLHTTP

```go
func SetRLHTTP(ctx context.Context, inst *Client) (context.Context, error)
```
Set RLHTTP stores the instance in the context map.

#### func  SetRLHTTPIfAbsent

```go
func SetRLHTTPIfAbsent(ctx context.Context, inst *Client) (context.Context, error)
```
SetRLHTTPIfAbsent sets only if not already present.

#### func  WithRLHTTP

```go
func WithRLHTTP(ctx context.Context, fn func(*Client))
```
WithRLHTTP calls fn with the instance or fallback.

#### type Client

```go
type Client struct {
}
```


#### func  FromRLHTTP

```go
func FromRLHTTP(ctx context.Context) *Client
```
FromRLHTTP returns the instance or the fallback.

#### func  GetRLHTTP

```go
func GetRLHTTP(ctx context.Context) (*Client, error)
```
Get the instance from the given context.

Will return ErrValueNotRLHTTP if the value in the context is not of type
*Client.

#### func  MustGetRLHTTP

```go
func MustGetRLHTTP(ctx context.Context) *Client
```
Attempt to get the instance from the given context. Panics if the operation
fails.

#### func  NewClient

```go
func NewClient(cnf *Config) *Client
```
Create a new rate-limited HTTP client instance.

#### func  NewDefault

```go
func NewDefault() *Client
```

#### func  TryGetRLHTTP

```go
func TryGetRLHTTP(ctx context.Context) (*Client, bool)
```
TryGetRLHTTP returns the instance and true if present and typed.

#### func (*Client) Do

```go
func (c *Client) Do(req *http.Request) (*http.Response, error)
```
Perform a HTTP request.

Please note that as this is essentially meant to be used as a middle man between
`apiclient.Client` and `http.Client`, we do not need to drain and close the body
here. That is firmly the responsibility of whatever consumes `rlhttp.Client`.

#### type Config

```go
type Config struct {
	Enabled bool           `json:"enabled"` // Rate limiting enabled?
	Timeout types.Duration `json:"timeout"` // Request timeout.
	Every   types.Duration `json:"every"`   // Time measure.
	Burst   int            `json:"burst"`   // Number of bursts
	Max     int            `json:"max"`     // Max requests per measure.
}
```

Rate limiter configuration.

The limiter spaces requests at an interval of Every/Max (i.e. `Max` requests per
`Every`).

If rate limiting is enabled, then the following conditions hold true:

1) If there is no `Timeout` then a default of 30 seconds is used, 2) If there is
no `Burst` then a default of 1 is used. 3) If there is no `Every` then a
validation error shall be raised. 4) If there is no `Max` then a validation
error shall be raised. 5) If `Burst` is greater than `Max` then a validation
error shall be raised.

It is advised to call `Validate` after populating `Config` and checking if there
are any raised errors.

If `Validate` does raise errors, those will need to be addressed first.

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Create a new default rate limiter configuration.

#### func (*Config) Validate

```go
func (c *Config) Validate() []error
```
Validate rate limiter settings. TODO: Change `Validate` to `Validate/Normalise`.
Don't forget!
