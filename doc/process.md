-*- Mode: gfm -*-

# process -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/process"
```

## Usage

```go
const (
	EventLoopSleep time.Duration = 250 * time.Millisecond
)
```

#### type ActionResult

```go
type ActionResult struct {
	Value   any
	Error   error
	Success bool
}
```


#### func  NewActionResult

```go
func NewActionResult(val any, err error) *ActionResult
```

#### func  NewEmptyActionResult

```go
func NewEmptyActionResult() *ActionResult
```

#### func (*ActionResult) IsError

```go
func (ar *ActionResult) IsError() bool
```

#### type CallbackFn

```go
type CallbackFn func(**State)
```

Callback function.

#### type Config

```go
type Config struct {
	Name     string        // Pretty name.
	Interval int           // `RunEvery` time interval.
	Function CallbackFn    // `Action` callback.
	OnStart  CallbackFn    // `Start` callback.
	OnStop   CallbackFn    // `Stop` callback.
	OnQuery  QueryFn       // `Query` callback.
	Logger   logger.Logger // Logger.
}
```

Process configuration structure.

#### func  NewDefaultConfig

```go
func NewDefaultConfig() *Config
```
Create a default process configuration.

#### type Manager

```go
type Manager interface {
	SetLogger(logger.Logger)
	Logger() logger.Logger
	SetContext(context.Context)
	Context() context.Context
	Create(*Config) *Process
	Add(*Process)
	Find(string) (*Process, bool)
	Run(string) bool
	Stop(string) bool
	StopAll() bool
	Processes() *[]*Process
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
      Function: func(state **State) {
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
	sync.Mutex

	Name     string        // Pretty name.
	Function CallbackFn    // `Action` callback.
	OnStart  CallbackFn    // `Start` callback.
	OnStop   CallbackFn    // `Stop` callback.
	OnQuery  QueryFn       // `Query` callback.
	Running  bool          // Is the process running?
	Interval time.Duration // `RunEvery` time interval.
}
```

Process structure.

To use:

1) Create a config:

```go

    conf := &process.Config{
      Name:     "Windows 95",
      Interval: 10,        // 10 seconds.
      Function: func(state **State) {
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
func NewProcessWithContext(config *Config, parent context.Context) *Process
```
Create a new process with the given configuration and parent context.

#### func (*Process) Context

```go
func (p *Process) Context() context.Context
```
Return the context for the process.

#### func (*Process) Query

```go
func (p *Process) Query(arg interface{}) interface{}
```
Query the running process.

This allows interaction with the process's base object without using `Action`.

#### func (*Process) Receive

```go
func (p *Process) Receive() interface{}
```
Receive data from the process with blocking.

#### func (*Process) ReceiveNonBlocking

```go
func (p *Process) ReceiveNonBlocking() (interface{}, bool)
```
Receive data from the process without blocking.

#### func (*Process) Run

```go
func (p *Process) Run() bool
```
Run the process with its action taking place on a continuous loop.

Returns 'true' if the process has been started, or 'false' if it is already
running.

#### func (*Process) Send

```go
func (p *Process) Send(data interface{})
```
Send data to the process with blocking.

#### func (*Process) SendNonBlocking

```go
func (p *Process) SendNonBlocking(data interface{})
```
Send data to the process without blocking.

#### func (*Process) SetContext

```go
func (p *Process) SetContext(parent context.Context)
```
Set the process's context.

#### func (*Process) SetLogger

```go
func (p *Process) SetLogger(lgr logger.Logger)
```
Set the process's logger.

#### func (*Process) SetWaitGroup

```go
func (p *Process) SetWaitGroup(wg *sync.WaitGroup)
```
Set the process's wait group.

#### func (*Process) Stop

```go
func (p *Process) Stop() bool
```
Stop the process.

Returns 'true' if the process was successfully stopped, or 'false' if it was not
running.

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

#### func  NewState

```go
func NewState() *State
```

#### func (*State) Context

```go
func (ps *State) Context() context.Context
```
Return the context for the parent process.

#### func (*State) Logger

```go
func (ps *State) Logger() logger.Logger
```

#### func (*State) Receive

```go
func (ps *State) Receive() (any, bool)
```
Read data from an external entity.

#### func (*State) ReceiveBlocking

```go
func (ps *State) ReceiveBlocking() any
```
Read data from an external entity with blocking.

#### func (*State) Send

```go
func (ps *State) Send(data any) bool
```
Send data from a process to an external entity.

#### func (*State) SendBlocking

```go
func (ps *State) SendBlocking(data any)
```
Send data from a process to an external entity with blocking.
