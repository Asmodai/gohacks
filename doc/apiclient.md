-*- Mode: gfm -*-

# apiclient -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/apiclient"
```

## Usage

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
	// `200` HTTP status code.
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
	// Returns the response body as an array of bytes, the HTTP status code, and
	// an error if one is triggered.
	//
	// You will need to remember to check both the error and status code.
	Get(*Params) ([]byte, int, error)

	// Perform a HTTP POST using the given API parameters.
	//
	// Returns the response body as an array of bytes, the HTTP status code, and
	// an error if one is triggered.
	//
	// You will need to remember to check both the error and status code.
	Post(*Params) ([]byte, int, error)
}
```

API Client

Like a finely-crafted sword, it can be wielded with skill if instructions are
followed.

1) Create your config:

    ```go
    conf := &apiclient.Config{
    	RequestsPerSecond: 5,    // 5 requests per second.
    	Timeout:           5,    // 5 seconds.
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
	Timeout int `json:"timeout"`
}
```

API client configuration.

#### func  NewConfig

```go
func NewConfig(reqsPerSec, timeout int) *Config
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


#### type MockClient

```go
type MockClient struct {
}
```

MockClient is a mock of Client interface.

#### func  NewMockClient

```go
func NewMockClient(ctrl *gomock.Controller) *MockClient
```
NewMockClient creates a new mock instance.

#### func (*MockClient) EXPECT

```go
func (m *MockClient) EXPECT() *MockClientMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use.

#### func (*MockClient) Get

```go
func (m *MockClient) Get(arg0 *Params) ([]byte, int, error)
```
Get mocks base method.

#### func (*MockClient) Post

```go
func (m *MockClient) Post(arg0 *Params) ([]byte, int, error)
```
Post mocks base method.

#### type MockClientMockRecorder

```go
type MockClientMockRecorder struct {
}
```

MockClientMockRecorder is the mock recorder for MockClient.

#### func (*MockClientMockRecorder) Get

```go
func (mr *MockClientMockRecorder) Get(arg0 any) *gomock.Call
```
Get indicates an expected call of Get.

#### func (*MockClientMockRecorder) Post

```go
func (mr *MockClientMockRecorder) Post(arg0 any) *gomock.Call
```
Post indicates an expected call of Post.

#### type MockHTTPClient

```go
type MockHTTPClient struct {
}
```

MockHTTPClient is a mock of HTTPClient interface.

#### func  NewMockHTTPClient

```go
func NewMockHTTPClient(ctrl *gomock.Controller) *MockHTTPClient
```
NewMockHTTPClient creates a new mock instance.

#### func (*MockHTTPClient) Do

```go
func (m *MockHTTPClient) Do(arg0 *http.Request) (*http.Response, error)
```
Do mocks base method.

#### func (*MockHTTPClient) EXPECT

```go
func (m *MockHTTPClient) EXPECT() *MockHTTPClientMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use.

#### type MockHTTPClientMockRecorder

```go
type MockHTTPClientMockRecorder struct {
}
```

MockHTTPClientMockRecorder is the mock recorder for MockHTTPClient.

#### func (*MockHTTPClientMockRecorder) Do

```go
func (mr *MockHTTPClientMockRecorder) Do(arg0 any) *gomock.Call
```
Do indicates an expected call of Do.

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
