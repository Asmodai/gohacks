-*- Mode: gfm -*-

# app -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/app"
```

## Usage

```go
const (
	// Time to sleep during main loop so we're a nice neighbour.
	EventLoopSleep time.Duration = 250 * time.Millisecond
)
```

#### type Application

```go
type Application struct {
	OnStart  OnSignalFn // Function called on app startup.
	OnExit   OnSignalFn // Function called on app exit.
	OnHUP    OnSignalFn // Function called when SIGHUP received.
	OnUSR1   OnSignalFn // Function called when SIGUSR1 received.
	OnUSR2   OnSignalFn // Function called when SIGUSR2 received.
	OnWINCH  OnSignalFn // Function used when SIGWINCH received.
	OnCHLD   OnSignalFn // Function used when SIGCHLD received.
	MainLoop MainLoopFn // Application main loop function.
}
```


#### func  NewApplication

```go
func NewApplication(cnf *Config) *Application
```
Create a new application.

#### func (*Application) Commit

```go
func (app *Application) Commit() string
```

#### func (*Application) Configuration

```go
func (app *Application) Configuration() config.Config
```

#### func (*Application) Context

```go
func (app *Application) Context() context.Context
```
Return the application's context.

#### func (*Application) Init

```go
func (app *Application) Init()
```

#### func (*Application) IsDebug

```go
func (app *Application) IsDebug() bool
```
Is the application using debug mode?

#### func (*Application) IsRunning

```go
func (app *Application) IsRunning() bool
```
Is the application running?

#### func (*Application) Logger

```go
func (app *Application) Logger() logger.Logger
```

#### func (*Application) Name

```go
func (app *Application) Name() string
```

#### func (*Application) ProcessManager

```go
func (app *Application) ProcessManager() process.IManager
```

#### func (*Application) Run

```go
func (app *Application) Run()
```
Start the application.

#### func (*Application) SetMainLoop

```go
func (app *Application) SetMainLoop(fn MainLoopFn)
```
Set the main loop callback.

#### func (*Application) SetOnCHLD

```go
func (app *Application) SetOnCHLD(fn OnSignalFn)
```
Set the `OnCHLD` callback.

#### func (*Application) SetOnExit

```go
func (app *Application) SetOnExit(fn OnSignalFn)
```
Set the `OnExit` callback.

#### func (*Application) SetOnHUP

```go
func (app *Application) SetOnHUP(fn OnSignalFn)
```
Set the `OnHUP` callback.

#### func (*Application) SetOnStart

```go
func (app *Application) SetOnStart(fn OnSignalFn)
```
Set the `OnStart` callback.

#### func (*Application) SetOnUSR1

```go
func (app *Application) SetOnUSR1(fn OnSignalFn)
```
Set the `OnUSR1` callback.

#### func (*Application) SetOnUSR2

```go
func (app *Application) SetOnUSR2(fn OnSignalFn)
```
Set the `OnUSR2` callback.

#### func (*Application) SetOnWINCH

```go
func (app *Application) SetOnWINCH(fn OnSignalFn)
```
Set the `OnWINCH` callback.

#### func (*Application) Terminate

```go
func (app *Application) Terminate()
```
Stop the application.

#### func (*Application) Version

```go
func (app *Application) Version() *semver.SemVer
```

#### type Config

```go
type Config struct {
	Name           string
	Version        *semver.SemVer
	Logger         logger.Logger
	ProcessManager process.IManager
	AppConfig      any
	Validators     config.ValidatorsMap
}
```


#### func  NewConfig

```go
func NewConfig() *Config
```

#### type MainLoopFn

```go
type MainLoopFn func(*Application) // Main loop callback function.

```


#### type OnSignalFn

```go
type OnSignalFn func(*Application) // Signal callback function.

```
