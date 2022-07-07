-*- Mode: gfm -*-

# apiclient -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/apiclient"
```

## Usage

#### type AuthBasic

```go
type AuthBasic struct {
	Username string
	Password string
}
```


#### type AuthToken

```go
type AuthToken struct {
	Header string
	Data   string
}
```


#### type Client

```go
type Client struct {
	Client  IHTTPClient
	Limiter *rate.Limiter
	Trace   *httptrace.ClientTrace
}
```

API Client

Like a finely-crafted sword, it can be wielded with skill if instructions are
followed.

1) Create your config:

    	conf := &apiclient.Config{
        RequestsPerSecond: 5,    // 5 requests per second.
    		Timeout:           5,    // 5 seconds.
    	}

2) Create your client

    api := apiclient.NewClient(conf)

3) ???

    params := &Params{
    	Url: "http://www.example.com/underpants",
    }

4) Profit

    data, code, err := api.Get(params)
    // check `err` and `code` here.
    // `data` will need to be converted from `[]byte`.

#### func  NewClient

```go
func NewClient(config *Config, logger logger.ILogger) *Client
```
Create a new API client with the given configuration.

#### func (*Client) Get

```go
func (c *Client) Get(data *Params) ([]byte, int, error)
```
Perform a HTTP GET using the given API parameters.

Returns the response body as an array of bytes, the HTTP status code, and an
error if one is triggered.

You will need to remember to check the error *and* the status code.

#### func (*Client) Post

```go
func (c *Client) Post(data *Params) ([]byte, int, error)
```
Perform a HTTP POST using the given API parameters.

Returns the response body as an array of bytes, the HTTP status code, and an
error if one is triggered.

You will need to remember to check the error *and* the status code.

#### type Config

```go
type Config struct {
	RequestsPerSecond int `json:"requests_per_second"`
	Timeout           int `json:"timeout"`
}
```

API client configuration.

`RequstsPerSecond` is the number of requests per second rate limiting. `Timeout`
is obvious.

#### func  NewConfig

```go
func NewConfig(ReqsPerSec, Timeout int) *Config
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
	Accept string
	Type   string
}
```


#### type IApiClient

```go
type IApiClient interface {
	Get(data *Params) ([]byte, int, error)
	Post(data *Params) ([]byte, int, error)
}
```

API client interface.

#### type IHTTPClient

```go
type IHTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
```

HTTP client interface.

This is to allow mocking of `net/http`'s `http.Client` in unit tests. You could
also, if drunk enough, provide your own HTTP client, as long as it conforms to
the interface. But you wouldn't want to do that, would you.

#### type MockCallbackFn

```go
type MockCallbackFn func(data *Params) ([]byte, int, error)
```


#### type MockClient

```go
type MockClient struct {
	GetFn  MockCallbackFn
	PostFn MockCallbackFn
}
```


#### func  NewMockClient

```go
func NewMockClient(config *Config, lgr logger.ILogger) *MockClient
```

#### func (*MockClient) Get

```go
func (c *MockClient) Get(data *Params) ([]byte, int, error)
```

#### func (*MockClient) Post

```go
func (c *MockClient) Post(data *Params) ([]byte, int, error)
```

#### type Params

```go
type Params struct {
	Url string // API URL.

	UseBasic bool
	UseToken bool

	Content ContentType
	Token   AuthToken
	Basic   AuthBasic

	Queries []*QueryParam
}
```

API client request parameters

#### func  NewParams

```go
func NewParams() *Params
```
Create a new API parameters object,

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
	Name    string
	Content string
}
```

API URL query parameter.

#### func  NewQueryParam

```go
func NewQueryParam(name, content string) *QueryParam
```
Create a new query parameter.
