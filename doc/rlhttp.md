-*- Mode: gfm -*-

# rlhttp -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/rlhttp"
```

## Usage

#### type Client

```go
type Client struct {
}
```


#### func  NewClient

```go
func NewClient(rlimiter *rate.Limiter, timeout time.Duration) *Client
```

#### func (*Client) Do

```go
func (c *Client) Do(req *http.Request) (*http.Response, error)
```
