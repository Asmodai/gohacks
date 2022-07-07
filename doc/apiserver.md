-*- Mode: gfm -*-

# apiserver -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/apiserver"
```

## Usage

#### func  CORSMiddleware

```go
func CORSMiddleware() gin.HandlerFunc
```

#### func  GetRouter

```go
func GetRouter(mgr process.IManager) (*gin.Engine, error)
```

#### func  SetDebugMode

```go
func SetDebugMode(debug bool)
```

#### func  Spawn

```go
func Spawn(mgr process.IManager, lgr logger.ILogger, config *Config) (*process.Process, error)
```

#### type Config

```go
type Config struct {
	Addr    string `json:"address"`
	Cert    string `json:"cert_file"`
	Key     string `json:"key_file"`
	UseTLS  bool   `json:"use_tls"`
	LogFile string `json:"log_file"`
}
```


#### func  NewConfig

```go
func NewConfig(addr, log, cert, key string, tls bool) *Config
```

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```

#### type Dispatcher

```go
type Dispatcher struct {
}
```


#### func  NewDefaultDispatcher

```go
func NewDefaultDispatcher() *Dispatcher
```

#### func  NewDispatcher

```go
func NewDispatcher(lgr logger.ILogger, config *Config) *Dispatcher
```

#### func (*Dispatcher) GetRouter

```go
func (d *Dispatcher) GetRouter() *gin.Engine
```

#### func (*Dispatcher) Start

```go
func (d *Dispatcher) Start()
```

#### func (*Dispatcher) Stop

```go
func (d *Dispatcher) Stop()
```

#### type DispatcherProc

```go
type DispatcherProc struct {
	sync.Mutex
}
```


#### func  NewDispatcherProc

```go
func NewDispatcherProc(lgr logger.ILogger, config *Config) *DispatcherProc
```

#### func (*DispatcherProc) Action

```go
func (p *DispatcherProc) Action(state **process.State)
```

#### type Document

```go
type Document struct {
	Data    interface{}    `json:"data"`
	Error   *ErrorDocument `json:"error"`
	Elapsed string         `json:"elapsed_time",omitempty`
}
```

nolint:govet

#### func  NewDocument

```go
func NewDocument(status int, data interface{}) *Document
```

#### func  NewErrorDocument

```go
func NewErrorDocument(status int, msg string) *Document
```

#### func (*Document) AddHeader

```go
func (d *Document) AddHeader(key, value string)
```

#### func (*Document) SetError

```go
func (d *Document) SetError(err *ErrorDocument)
```

#### func (*Document) Status

```go
func (d *Document) Status() int
```

#### func (*Document) Write

```go
func (d *Document) Write(ctx *gin.Context)
```

#### type ErrorDocument

```go
type ErrorDocument struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
```


#### func  NewError

```go
func NewError(status int, msg string) *ErrorDocument
```

#### type IServer

```go
type IServer interface {
	ListenAndServeTLS(string, string) error
	ListenAndServe() error
	Shutdown(context.Context) error
	SetTLSConfig(*tls.Config)
}
```


#### type Server

```go
type Server struct {
}
```


#### func  NewDefaultServer

```go
func NewDefaultServer() *Server
```

#### func  NewServer

```go
func NewServer(addr string, router *gin.Engine) *Server
```

#### func (*Server) ListenAndServe

```go
func (s *Server) ListenAndServe() error
```

#### func (*Server) ListenAndServeTLS

```go
func (s *Server) ListenAndServeTLS(cert, key string) error
```

#### func (*Server) SetTLSConfig

```go
func (s *Server) SetTLSConfig(conf *tls.Config)
```

#### func (*Server) Shutdown

```go
func (s *Server) Shutdown(ctx context.Context) error
```
