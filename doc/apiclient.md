-*- Mode: gfm -*-

# apiclient -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/apiclient"
```

## Usage

```go
var (
	ErrInvalidAuthMethod = errors.Base("invalid authentication method")
	ErrMissingArgument   = errors.Base("missing argument")
	ErrNotOk             = errors.Base("not ok")
)
```

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
type Client interface {
	Get(*Params) ([]byte, int, error)
	Post(*Params) ([]byte, int, error)
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
    	URL: "http://www.example.com/underpants",
    }

4) Profit

    data, code, err := api.Get(params)
    // check `err` and `code` here.
    // `data` will need to be converted from `[]byte`.

#### func  NewClient

```go
func NewClient(config *Config, logger logger.Logger) Client
```
Create a new API client with the given configuration.

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
	Accept string
	Type   string
}
```


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
	URL string // API URL.

	UseBasic bool
	UseToken bool

	Content ContentType
	Token   AuthToken
	Basic   AuthBasic

	Queries []*QueryParam
}
```

API client request parameters.

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
