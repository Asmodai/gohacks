-*- Mode: gfm -*-

# rlhttp -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/rlhttp"
```

## Usage

#### type RLHTTPClient

```go
type RLHTTPClient struct {
}
```


#### func  NewClient

```go
func NewClient(rl *rate.Limiter, timeout time.Duration) *RLHTTPClient
```

#### func (*RLHTTPClient) Do

```go
func (c *RLHTTPClient) Do(req *http.Request) (*http.Response, error)
```
