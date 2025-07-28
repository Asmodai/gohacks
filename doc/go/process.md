<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# process -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/process"
```

## Usage

```go
const (
	ContextKeyProcManager = "_DI_PROC_MGR"
)
```

```go
var (
	ErrValueNotProcessManager = errors.Base("value is not process.Manager")
)
```

#### func  SetProcessManager

```go
func SetProcessManager(ctx context.Context, inst Manager) (context.Context, error)
```
Set the process manager value to the context map.

#### type CallbackFn

```go
type CallbackFn func(*State)
```

Callback function.

#### type Config

```go
type Config struct {
	Name      string                // Pretty name.
	Interval  types.Duration        // `RunEvery` time interval.
	Function  CallbackFn            // `Action` callback.
	OnStart   CallbackFn            // `Start` callback.
	OnStop    CallbackFn            // `Stop` callback.
	OnQuery   QueryFn               // `Query` callback.
	Responder responder.Respondable // Responder object.
}
```

Process configuration structure.

#### func  NewConfig

```go
func NewConfig() *Config
```
Create a default process configuration.

#### type Manager

```go
type Manager interface {
	Logger() logger.Logger
	SetContext(context.Context)
	SetLogger(logger.Logger)
	Context() context.Context
	Create(*Config) *Process
	Add(*Process)
	Find(string) (*Process, bool)
	Run(string) bool
	Stop(string) bool
	StopAll() StopAllResults
	Processes() []*Process
	Count() int
}
```

Process manager structure.

To use,

1) Create a new process manager:

```go

    procmgr := process.NewManager()

```

2) Create your process configuration:

```go

    conf := &process.Config{
      Name:     "Windows 95",
      Interval: 10, // seconds
      Function: func(state *State) {
        // Crash or something.
      }
    }

```

3) Create the process itself.

```go

    proc := procmgr.Create(conf)

```

4) Run the process.

```go

    procmgr.Run("Windows 95")

```

/or/

```go

    proc.Run()

```

Manager is optional, as you can create processes directly.

#### func  GetProcessManager

```go
func GetProcessManager(ctx context.Context) (Manager, error)
```
Get the process manager from the given context.

Will return `ErrValueNotProcessManager` if the value in the context is not of
type `process.Manager`.

#### func  MustGetProcessManager

```go
func MustGetProcessManager(ctx context.Context) Manager
```
Attempt to get the process manager from the given context. Panics if the
operation fails.

#### func  NewManager

```go
func NewManager() Manager
```
Create a new process manager.

#### func  NewManagerWithContext

```go
func NewManagerWithContext(parent context.Context) Manager
```
Create a new process manager with a given parent context.

#### type Process

```go
type Process struct {
}
```

Process structure.

To use:

1) Create a config:

```go

    conf := &process.Config{
      Name:     "Windows 95",
      Interval: 10,        // 10 seconds.
      Function: func(state *State) {
        // Crash or something.
      },
    }

```

2) Create a process:

```go

    proc := process.NewProcess(conf)

```

3) Run the process:

```go

    go proc.Run()

```

4) Send data to the process:

```go

    proc.Send("Blue Screen of Death")

```

5) Read data from the process:

```go

    data := proc.Receive()

```

6) Stop the process

```go

    proc.Stop()

```

    will stop the process.

#### func  NewProcess

```go
func NewProcess(config *Config) *Process
```
Create a new process with the given configuration.

#### func  NewProcessWithContext

```go
func NewProcessWithContext(parent context.Context, config *Config) *Process
```
Create a new process with the given configuration and parent context.

#### func (*Process) Context

```go
func (p *Process) Context() context.Context
```
Return the context for the process.

#### func (*Process) Invoke

```go
func (p *Process) Invoke(event events.Event) events.Event
```

#### func (*Process) Name

```go
func (p *Process) Name() string
```

#### func (*Process) Query

```go
func (p *Process) Query(arg interface{}) interface{}
```
Query the running process.

This allows interaction with the process's base object without using `Action`.

#### func (*Process) RespondsTo

```go
func (p *Process) RespondsTo(event events.Event) bool
```

#### func (*Process) Run

```go
func (p *Process) Run() bool
```
Run the process with its action taking place on a continuous loop.

Returns 'true' if the process has been started, or 'false' if it is already
running.

#### func (*Process) Running

```go
func (p *Process) Running() bool
```
Is the process running?

#### func (*Process) Stop

```go
func (p *Process) Stop() bool
```
Stop the process.

Returns 'true' if the process was successfully stopped, or 'false' if it was not
running.

#### func (*Process) Type

```go
func (p *Process) Type() string
```

#### type QueryFn

```go
type QueryFn func(interface{}) interface{}
```


#### type State

```go
type State struct {
}
```

Internal state for processes.

#### func (*State) Context

```go
func (ps *State) Context() context.Context
```
Return the context for the parent process.

#### func (*State) Invoke

```go
func (ps *State) Invoke(event events.Event) (events.Event, bool)
```

#### func (*State) Logger

```go
func (ps *State) Logger() logger.Logger
```

#### func (*State) RespondsTo

```go
func (ps *State) RespondsTo(event events.Event) bool
```

#### type StopAllResults

```go
type StopAllResults map[string]bool
```
