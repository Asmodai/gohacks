<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# app -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/app"
```

## Usage

#### type Application

```go
type Application interface {
	// Cannot be invoked once the application has been initialised.
	ParseConfig()

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

	// Set the parent context for the application.
	//
	// Danger.  Setting this while the application is running can cause
	// unintended side effects due to the old context's cancel function
	// being executed.
	//
	// It is advisable to run this prior to initialisation.
	SetContext(context.Context)

	// Return the application's process manager instance.
	ProcessManager() process.Manager

	// Return the application's logger instance.
	Logger() logger.Logger

	// Return the application's configuration.
	Configuration() config.Config

	// Set the callback that will be invoked when the application starts.
	//
	// If not set, then the default startup handler will be invoked.
	//
	// This cannot be set once the application has been initialised.
	SetOnStart(CallbackFn)

	// Set the callback that will be invoked when the application exits.
	//
	// If not set, then the default exit handler will be invoked.
	//
	/// This cannot be set once the application has been initialised.
	SetOnExit(CallbackFn)

	// Set the callback that will be invoked whenever the event loop
	// fires.
	//
	// If not set, then the default main loop callback will be invoked.
	//
	// This cannot be set once the application has been initialised.
	SetMainLoop(MainLoopFn)

	// Is the application running?
	IsRunning() bool

	// Is the application in 'debug' mode.
	IsDebug() bool

	// Add a responder to the application's responder chain.
	AddResponder(responder.Respondable) (responder.Respondable, error)

	// Remove a responder from the application's responder chain.
	RemoveResponder(responder.Respondable) bool

	// Send an event to the application's responder.
	//
	// Event will be consumed by the first responder that handles it.
	SendFirstResponder(events.Event) (events.Event, bool)

	// Send an event to all the application's responders.
	SendAllResponders(events.Event) []events.Event

	// Return the name of the application's responder chain.
	//
	// Implements `Respondable`.
	Type() string

	// Ascertain if any of the application's responders will respond to
	// an event.
	//
	// The first responder found that responds to the event will result
	// in `true` being returned.
	//
	// Implements `Respondable`.
	RespondsTo(events.Event) bool

	// Send an event to the application's responders.
	//
	// The first object that can respond to the event will consume it.
	//
	// Implements `Respondable`.
	Invoke(events.Event) events.Event
}
```

Application.

#### func  NewApplication

```go
func NewApplication(cnf *Config) Application
```
Create a new application.

#### type CallbackFn

```go
type CallbackFn func(Application)
```

Signal callback function type.

#### type Config

```go
type Config struct {
	// The application's pretty name.
	Name string

	// The application's version number.
	Version *semver.SemVer

	// The application's configuration.
	AppConfig any

	// Validators used to validate the application's configuration.
	Validators config.ValidatorsMap

	// Require CLI flags like '-config' to be provided?
	RequireCLI bool
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
type OnSignalFn func(os.Signal)
```

Callback for the signal responder.

#### type SignalResponder

```go
type SignalResponder struct {
}
```

Signal responder.

#### func  NewSignalResponder

```go
func NewSignalResponder() *SignalResponder
```

#### func (*SignalResponder) Invoke

```go
func (sr *SignalResponder) Invoke(evt events.Event) events.Event
```
Invokes the given event.

#### func (*SignalResponder) Name

```go
func (sr *SignalResponder) Name() string
```
Returns the name of the responder.

#### func (*SignalResponder) RespondsTo

```go
func (sr *SignalResponder) RespondsTo(evt events.Event) bool
```
Returns whether the responder can respond to a given event.

#### func (*SignalResponder) SetOnSignal

```go
func (sr *SignalResponder) SetOnSignal(callback OnSignalFn)
```
Sets the callback function.

#### func (*SignalResponder) Type

```go
func (sr *SignalResponder) Type() string
```
Returns the type of the responder.
