-*- Mode: gfm -*-

# rpc -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/rpc"
```

## Usage

#### func  Add

```go
func Add(mgr process.IManager, t reflect.Type) (bool, error)
```

#### func  Spawn

```go
func Spawn(mgr process.IManager, lgr logger.ILogger, ctx context.Context, cnf *Config) (*process.Process, error)
```

#### type Client

```go
type Client struct {
}
```


#### func  NewClient

```go
func NewClient(cnf *Config, lgr logger.ILogger) *Client
```

#### func (*Client) Call

```go
func (c *Client) Call(method string, args any, reply any) error
```

#### func (*Client) Close

```go
func (c *Client) Close() error
```

#### func (*Client) Dial

```go
func (c *Client) Dial() error
```

#### type Config

```go
type Config struct {
	Proto string `json:"protocol"`
	Addr  string `json:"address"`
}
```


#### func  NewConfig

```go
func NewConfig(proto, addr string) *Config
```

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```

#### type IManager

```go
type IManager interface {
	SetContext(context.Context)
	SetLogger(logger.ILogger)
	Add(reflect.Type) bool
	Start()
	Shutdown()
}
```


#### type IRPCAble

```go
type IRPCAble interface {
	RpcInfo(NoArgs, *Info) error
}
```


#### func  NewRPCAble

```go
func NewRPCAble(t reflect.Type) IRPCAble
```

#### type Info

```go
type Info struct {
	Name    string
	Version int
}
```


#### func  NewInfo

```go
func NewInfo(name string, version int) *Info
```

#### func (*Info) String

```go
func (i *Info) String() string
```

#### type Manager

```go
type Manager struct {
}
```


#### func  NewManager

```go
func NewManager(cnf *Config, ctx context.Context, lgr logger.ILogger) *Manager
```

#### func (*Manager) Add

```go
func (m *Manager) Add(t reflect.Type) bool
```

#### func (*Manager) SetContext

```go
func (m *Manager) SetContext(parent context.Context)
```

#### func (*Manager) SetLogger

```go
func (m *Manager) SetLogger(lgr logger.ILogger)
```

#### func (*Manager) Shutdown

```go
func (m *Manager) Shutdown()
```

#### func (*Manager) Start

```go
func (m *Manager) Start()
```

#### type ManagerProc

```go
type ManagerProc struct {
	sync.Mutex
}
```


#### func  NewManagerProc

```go
func NewManagerProc(lgr logger.ILogger, ctx context.Context, config *Config) *ManagerProc
```

#### func (*ManagerProc) Action

```go
func (p *ManagerProc) Action(state **process.State)
```

#### type NoArgs

```go
type NoArgs struct{}
```
