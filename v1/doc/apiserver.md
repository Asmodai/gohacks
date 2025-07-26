<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# apiserver -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/apiserver"
```

## Usage

```go
const (
	MinimumTimeout = 5
)
```

```go
var (
	// List of allowed HTTP headers.
	AllowedHeaders = []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"Accept",
		"Origin",
		"Cache-Control",
		"X-Requested-With",
	}

	// List of allowed HTTP methods.
	AllowedMethods = []string{
		"POST",
		"OPTIONS",
		"GET",
		"PUT",
		"DELETE",
		"PATCH",
	}
)
```

```go
var (
	// Triggered when no process manager instance is provided.
	ErrNoProcessManager = errors.Base("no process manager")

	// Triggered when no API dispatcher instance is provided.
	ErrNoDispatcherProc = errors.Base("no dispatcher process")

	// Triggered if the process responds with an unexpected data type.
	ErrWrongReturnType = errors.Base("wrong return type")
)
```

#### func  CORSMiddleware

```go
func CORSMiddleware() gin.HandlerFunc
```
A gin-gonic handler for handling CORS.

#### func  GetRouter

```go
func GetRouter(mgr process.Manager) (*gin.Engine, error)
```
Get the router currently in use by the registered dispatcher process.

Should no dispatcher process be in the process manager, then
`ErrNoDispatcherProc` will be returned.

Should the dispatcher process return an unexpected value type, then
`ErrWrongReturnType` will be returned.

#### func  SetDebugMode

```go
func SetDebugMode(debug bool)
```
Set debug mode to the given flag.

This will reconfigure Gin-Gonic to either 'release' or 'debug' mode depending on
the boolean value of the flag.

#### func  Spawn

```go
func Spawn(mgr process.Manager, lgr logger.Logger, config *Config) (*process.Process, error)
```
Spawn an API dispatcher process.

#### type Config

```go
type Config struct {
	// Network address that the server will bind to.
	//
	// This is in the format of <address>:<port>.
	// To bind to all available addresses, specify ":<port>" only.
	Addr string `json:"address"`

	// Path to an SSL certificate file if TLS is required.
	Cert string `json:"cert_file"`

	// Path to an SSL key file if TLS is required.
	Key string `json:"key_file"`

	// Should the server use TLS?
	UseTLS bool `json:"use_tls"`

	// Path to the log file for the server.
	LogFile string `json:"log_file"`
}
```

API server configuration.

#### func  NewConfig

```go
func NewConfig(addr, log, cert, key string, tls bool) *Config
```
Create a new configuration.

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Create a new default configuration.

This will create a configuration that has default values.

#### func (*Config) Host

```go
func (c *Config) Host() (string, error)
```
Return the hostname on which the server is bound.

#### func (*Config) Port

```go
func (c *Config) Port() (int, error)
```
Return the port number on which the server is listening.

#### func (*Config) Validate

```go
func (c *Config) Validate() []error
```
Validate the configuration.

#### type Dispatcher

```go
type Dispatcher struct {
}
```

API route dispatcher.

#### func  NewDefaultDispatcher

```go
func NewDefaultDispatcher() *Dispatcher
```
Create a new API route dispatcher with default values.

The dispatcher returned by this function will listen on port 8080 and bind to
all available addresses on the host machine.

#### func  NewDispatcher

```go
func NewDispatcher(lgr logger.Logger, config *Config) *Dispatcher
```
Create a new API route dispatcher.

#### func (*Dispatcher) GetRouter

```go
func (d *Dispatcher) GetRouter() *gin.Engine
```
Return the router used by this dispatcher.

#### func (*Dispatcher) Start

```go
func (d *Dispatcher) Start()
```
Start the API route dispatcher.

#### func (*Dispatcher) Stop

```go
func (d *Dispatcher) Stop()
```
Stop the API route dispatcher.

#### type DispatcherProc

```go
type DispatcherProc struct {
}
```

Dispatcher process.

#### func  NewDispatcherProc

```go
func NewDispatcherProc(lgr logger.Logger, config *Config) *DispatcherProc
```
Create a new dispatcher process.

#### func (*DispatcherProc) Action

```go
func (p *DispatcherProc) Action(state **process.State)
```
Action invoked when the dispatcher process receives a message.

#### type Document

```go
type Document struct {

	// JSON document data.
	Data any `json:"data,omitempty"`

	// Number of elements present should `Data` be an array of some kind.
	Count int64 `json:"count"`

	// Error document.
	Error *ErrorDocument `json:"error,omitempty"`

	// Time taken to generate the JSON document.
	Elapsed string `json:"elapsed_time,omitempty"`
}
```

JSON document.

#### func  NewDocument

```go
func NewDocument(status int, data any) *Document
```
Generate a new JSON document.

#### func  NewErrorDocument

```go
func NewErrorDocument(status int, msg string) *Document
```
Create a new JSON document with an embedded error document.

#### func (*Document) AddHeader

```go
func (d *Document) AddHeader(key, value string)
```
Add a header to the document's HTTP response.

#### func (*Document) SetError

```go
func (d *Document) SetError(err *ErrorDocument)
```
Set the `Error` component of the document.

#### func (*Document) Status

```go
func (d *Document) Status() int
```
Return the document's HTTP status code response.

#### func (*Document) Write

```go
func (d *Document) Write(ctx *gin.Context)
```
Write the document to the given gin-gonic context.

#### type ErrorDocument

```go
type ErrorDocument struct {
	// HTTP status code.
	Status int `json:"status"`

	// Message describing the error.
	Message string `json:"message"`
}
```

JSON error document.

#### func  NewError

```go
func NewError(status int, msg string) *ErrorDocument
```
Create a new error document with the given status and message.

#### type Server

```go
type Server interface {
	// Bind and listen to configured address/port and serve HTTPS requests.
	ListenAndServeTLS(string, string) error

	// Bind and listen to configured address/port and serve HTTP requests.
	ListenAndServe() error

	// Shut down the API server.
	Shutdown(context.Context) error

	// Set the TLS configuration for HTTPS mode.
	SetTLSConfig(*tls.Config)
}
```

API server.

#### func  NewDefaultServer

```go
func NewDefaultServer() Server
```
Create a new API server using default configuration. This will create a new
server configured for HTTP only.

#### func  NewServer

```go
func NewServer(addr string, router *gin.Engine) Server
```
Create a new API server using the given address/port combination and gin-gonic
engine.
