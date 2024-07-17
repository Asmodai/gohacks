-*- Mode: gfm -*-

# apiserver -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/apiserver"
```

## Usage

```go
const (
	DefaultFileMode = 0644
)
```

```go
const (
	MinimumTimeout = 5
)
```

```go
var (
	ErrNoProcessManager = errors.Base("no process manager")
	ErrNoDispatcherProc = errors.Base("no dispatcher process")
	ErrWrongReturnType  = errors.Base("wrong return type")
)
```

#### func  CORSMiddleware

```go
func CORSMiddleware() gin.HandlerFunc
```

#### func  GetRouter

```go
func GetRouter(mgr process.Manager) (*gin.Engine, error)
```

#### func  SetDebugMode

```go
func SetDebugMode(debug bool)
```

#### func  Spawn

```go
func Spawn(mgr process.Manager, lgr logger.Logger, config *Config) (*process.Process, error)
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

#### func (*Config) Host

```go
func (c *Config) Host() (string, error)
```

#### func (*Config) Port

```go
func (c *Config) Port() (int, error)
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
func NewDispatcher(lgr logger.Logger, config *Config) *Dispatcher
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
func NewDispatcherProc(lgr logger.Logger, config *Config) *DispatcherProc
```

#### func (*DispatcherProc) Action

```go
func (p *DispatcherProc) Action(state **process.State)
```

#### type Document

```go
type Document struct {
	Data    interface{}    `json:"data,omitempty"`
	Count   int64          `json:"count"`
	Error   *ErrorDocument `json:"error,omitempty"`
	Elapsed string         `json:"elapsed_time,omitempty"`
}
```


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

#### type MockServer

```go
type MockServer struct {
}
```

MockServer is a mock of Server interface.

#### func  NewMockServer

```go
func NewMockServer(ctrl *gomock.Controller) *MockServer
```
NewMockServer creates a new mock instance.

#### func (*MockServer) EXPECT

```go
func (m *MockServer) EXPECT() *MockServerMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use.

#### func (*MockServer) ListenAndServe

```go
func (m *MockServer) ListenAndServe() error
```
ListenAndServe mocks base method.

#### func (*MockServer) ListenAndServeTLS

```go
func (m *MockServer) ListenAndServeTLS(arg0, arg1 string) error
```
ListenAndServeTLS mocks base method.

#### func (*MockServer) SetTLSConfig

```go
func (m *MockServer) SetTLSConfig(arg0 *tls.Config)
```
SetTLSConfig mocks base method.

#### func (*MockServer) Shutdown

```go
func (m *MockServer) Shutdown(arg0 context.Context) error
```
Shutdown mocks base method.

#### type MockServerMockRecorder

```go
type MockServerMockRecorder struct {
}
```

MockServerMockRecorder is the mock recorder for MockServer.

#### func (*MockServerMockRecorder) ListenAndServe

```go
func (mr *MockServerMockRecorder) ListenAndServe() *gomock.Call
```
ListenAndServe indicates an expected call of ListenAndServe.

#### func (*MockServerMockRecorder) ListenAndServeTLS

```go
func (mr *MockServerMockRecorder) ListenAndServeTLS(arg0, arg1 any) *gomock.Call
```
ListenAndServeTLS indicates an expected call of ListenAndServeTLS.

#### func (*MockServerMockRecorder) SetTLSConfig

```go
func (mr *MockServerMockRecorder) SetTLSConfig(arg0 any) *gomock.Call
```
SetTLSConfig indicates an expected call of SetTLSConfig.

#### func (*MockServerMockRecorder) Shutdown

```go
func (mr *MockServerMockRecorder) Shutdown(arg0 any) *gomock.Call
```
Shutdown indicates an expected call of Shutdown.

#### type Server

```go
type Server interface {
	ListenAndServeTLS(string, string) error
	ListenAndServe() error
	Shutdown(context.Context) error
	SetTLSConfig(*tls.Config)
}
```


#### func  NewDefaultServer

```go
func NewDefaultServer() Server
```

#### func  NewServer

```go
func NewServer(addr string, router *gin.Engine) Server
```
