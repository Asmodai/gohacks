<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

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
func NewClient(limit int, timeout time.Duration) *Client
```

#### func (*Client) Do

```go
func (c *Client) Do(req *http.Request) (*http.Response, error)
```
