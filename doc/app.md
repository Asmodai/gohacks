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
type Application interface {
	// Initialises the application object.
	//
	// This must be called, as it does several things to set up the
	// various facilities (such as logging) used by the application.
	Init()

	// Run the application.
	//
	// This enters the application's event loop, which will block until
	// the application is subsequently terminated.
	Run()

	// Terminate the application.
	//
	// Breaks out of the event loop, returning control back to the calling
	// function.
	Terminate()

	// Return the application's pretty name.
	Name() string

	// Return the application's version.
	Version() *semver.SemVer

	// Return the application's version control commit identifier.
	Commit() string

	// Return the application's context.
	Context() context.Context

	// Return the application's process manager instance.
	ProcessManager() process.Manager

	// Return the application's logger instance.
	Logger() logger.Logger

	// Return the application's configuration.
	Configuration() config.Config

	// Set the callback that will be invoked when the application starts.
	SetOnStart(OnSignalFn)

	// Set the callback that will be invoked when the application exits.
	//
	// If not set, then the default exit handler will be invoked.
	SetOnExit(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a HUP signal.
	SetOnHUP(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a USR1 signal.
	SetOnUSR1(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a USR2 signal.
	SetOnUSR2(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a WINCH signal.
	//
	// Be careful with this, as it will fire whenever the controlling
	// terminal is resized.
	SetOnWINCH(OnSignalFn)

	// Set the callback that will be invoked when the application
	// receives a CHLD signal.
	SetOnCHLD(OnSignalFn)

	// Set the callback that will be invoked whenever the event loop
	// fires.
	SetMainLoop(MainLoopFn)

	// Is the application running?
	IsRunning() bool

	// Is the application in 'debug' mode.
	IsDebug() bool
}
```

Application.

#### func  NewApplication

```go
func NewApplication(cnf *Config) Application
```
Create a new application.

#### type Config

```go
type Config struct {
	// The application's pretty name.
	Name string

	// The application's version number.
	Version *semver.SemVer

	// The application's logger instance.
	Logger logger.Logger

	// The application's process manager instance.
	ProcessManager process.Manager

	// The application's configuration.
	AppConfig any

	// Validators used to validate the application's configuration.
	Validators config.ValidatorsMap
}
```

Application configuration.

#### func  NewConfig

```go
func NewConfig() *Config
```
Create a new empty configuration.

#### type MainLoopFn

```go
type MainLoopFn func(Application)
```

Main loop callback function type.

#### type OnSignalFn

```go
type OnSignalFn func(Application)
```

Signal callback function type.
