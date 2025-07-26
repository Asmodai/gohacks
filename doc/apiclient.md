<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# apiclient -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/apiclient"
```

## Usage

```go
const (
	ReasonNone int = iota
	ReasonRetryAfter
	ReasonRateLimitExceeded
	ReasonHTTPStatusCode
	ReasonMaximum

	ReasonInvalid string = "invalid reason"
)
```

```go
var (
	// Triggered when an invalid authentication method is passed via the API
	// parameters.  Will also be triggered if both basic auth and auth token
	// methods are specified in the same parameters.
	ErrInvalidAuthMethod = errors.Base("invalid authentication method")

	// Triggered if a required authentication method argument is not provided in
	// the API parameters.
	ErrMissingArgument = errors.Base("missing argument")

	// Triggered if the result of an API call via the client does not have a
	// `2xx` HTTP status code or fails the user-defined success check.
	ErrNotOk = errors.Base("not ok")
)
```

#### type AuthBasic

```go
type AuthBasic struct {
	Username string
	Password string
}
```

Basic authentication configuration.

#### type AuthToken

```go
type AuthToken struct {
	// The HTTP header added to the request that contains the token data.
	Header string

	// The authentication token.
	Data string
}
```

Authentication token configuration.

This type does not care about the token type. If you intend to make use of token
types such as JWTs then you must implement that code yourself.

#### type Client

```go
type Client interface {
	// Perform a HTTP GET using the given API parameters.
	//
	// Returns the response body as an array of bytes, the HTTP status
	// code, and an error if one is triggered.
	//
	// You will need to remember to check both the error and status code.
	Get(*Params) Response

	// Perform a HTTP POST using the given API parameters.
	//
	// Returns the response body as an array of bytes, the HTTP status
	// code, and an error if one is triggered.
	//
	// You will need to remember to check both the error and status code.
	Post(*Params) Response

	// Perform a HTTP get using the given API parameters and context.
	//
	// Returns the response body as an array of bytes, the HTTP status
	// code, and an error if one is triggered.
	//
	// You will need to remember to check both the error and status code.
	GetWithContext(context.Context, *Params) Response

	// Perform a HTTP POST using the given API parameters and context.
	//
	// Returns the response body as an array of bytes, the HTTP status
	// code, and an error if one is triggered.
	//
	// You will need to remember to check both the error and status code.
	PostWithContext(context.Context, *Params) Response
}
```

API Client

Like a finely-crafted sword, it can be wielded with skill if instructions are
followed.

1) Create your config:

    ```go
    conf := &apiclient.Config{
    	RequestsPerSecond: 5,    // 5 requests per second.
    	Timeout:           types.Duration(5*time.Second),
    }
    ```

2) Create your client

    ```go
    api := apiclient.NewClient(conf)
    ```

3) ???

    ```go
    params := &Params{
    	URL: "http://www.example.com/underpants",
    }
    ```

4) Profit

    ```go
    data, code, err := api.Get(params)
    // check `err` and `code` here.
    // `data` will need to be converted from `[]byte`.
    ```

The client supports both the "Auth Basic" schema and authentication tokens
passed via HTTP headers. You need to ensure you pick either one or the other,
not both. Attempting to use both will generate an `invalid authentication
method` error.

For full information about the API client parameters, please see the
documentation for the `Params` type.

#### func  NewClient

```go
func NewClient(config *Config, logger logger.Logger) Client
```
Create a new API client with the given configuration.

#### type Config

```go
type Config struct {
	// The number of requests per second should rate limiting be required.
	RequestsPerSecond int `json:"requests_per_second"`

	// HTTP connection timeout value.
	Timeout types.Duration `json:"timeout"`

	// Callback to call in order to check request success.
	SuccessCheck SuccessCheckFn `json:"-"`
}
```

API client configuration.

#### func  NewConfig

```go
func NewConfig(reqsPerSec int, timeout types.Duration) *Config
```
Create a new API client configuration.

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Return a new default API client configuration.

#### type ContentType

```go
type ContentType struct {
	// The MIME type that we wish to accept.  Sent in the request via the `Accept`
	// header.
	Accept string

	// The MIME type of the data we are sending.  Sent in the request via the
	// `Content-Type` header.
	Type string
}
```

Content types configuration.

#### type HTTPClient

```go
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}
```


#### type Params

```go
type Params struct {
	// Request URL.
	URL string

	// Use the Basic Auth schema to authenticate to the remote server.
	UseBasic bool

	// Use an 'authentication token' to authenticate to the remote server.
	UseToken bool

	// MIME content type of the data we are sending and wish to accept from the
	// remote server.
	Content ContentType

	// Authentication token configuration.
	Token AuthToken

	// Basic Auth configuration.
	Basic AuthBasic

	// Queries that are sent via the HTTP request.
	Queries []*QueryParam
}
```

API client request parameters.

Although we support both Basic Auth and authentication tokens, only one should
be used. If both are specified, an error will be triggered.

#### func  NewParams

```go
func NewParams() *Params
```
Create a new API parameters object.

#### func (*Params) AddQueryParam

```go
func (p *Params) AddQueryParam(name, content string) *QueryParam
```
Add a new query parameter.

#### func (*Params) ClearQueryParams

```go
func (p *Params) ClearQueryParams()
```
Clear all query parameters.

#### func (*Params) SetUseBasic

```go
func (p *Params) SetUseBasic(val bool)
```
Enable/disable basic authentication.

#### func (*Params) SetUseToken

```go
func (p *Params) SetUseToken(val bool)
```
Enable/disable authentication token.

#### type QueryParam

```go
type QueryParam struct {
	// Name of the query parameter.
	Name string

	// Content of the query parameter
	Content string
}
```

API URL query parameter.

#### func  NewQueryParam

```go
func NewQueryParam(name, content string) *QueryParam
```
Create a new query parameter.

#### type RateLimitInfo

```go
type RateLimitInfo struct {
	Limit      int           // Total limit count.
	Remaining  int           // Number of calls remaining.
	ResetTime  time.Time     // Limit reset time.
	RetryAfter time.Duration // Retry delay.
	ReasonCode int           // Inferred throttle reason.
}
```

Rate limit information.

#### func  NewRateLimitInfo

```go
func NewRateLimitInfo(code int, headers http.Header) RateLimitInfo
```
Return a new rate limit info object.

#### func (RateLimitInfo) IsRateLimited

```go
func (obj RateLimitInfo) IsRateLimited() bool
```
Has the rate limit been reached?

#### func (RateLimitInfo) Reason

```go
func (obj RateLimitInfo) Reason() string
```
Return the throttling reason (if any) as a string.

#### func (RateLimitInfo) String

```go
func (obj RateLimitInfo) String() string
```
Return the rate limit information in string format.

#### type Response

```go
type Response struct {
	Body       []byte
	StatusCode int
	Headers    http.Header
	Limits     RateLimitInfo
	Error      error
}
```


#### func  NewResponse

```go
func NewResponse(code int, body []byte, headers http.Header, err error) Response
```

#### func  NewResponseFromError

```go
func NewResponseFromError(err error) Response
```

#### func  NewResponseWithCodeFromError

```go
func NewResponseWithCodeFromError(code int, err error) Response
```

#### type SuccessCheckFn

```go
type SuccessCheckFn func(int) bool
```
