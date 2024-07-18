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

#### type MockApplication

```go
type MockApplication struct {
}
```

MockApplication is a mock of Application interface.

#### func  NewMockApplication

```go
func NewMockApplication(ctrl *gomock.Controller) *MockApplication
```
NewMockApplication creates a new mock instance.

#### func (*MockApplication) Commit

```go
func (m *MockApplication) Commit() string
```
Commit mocks base method.

#### func (*MockApplication) Configuration

```go
func (m *MockApplication) Configuration() config.Config
```
Configuration mocks base method.

#### func (*MockApplication) Context

```go
func (m *MockApplication) Context() context.Context
```
Context mocks base method.

#### func (*MockApplication) EXPECT

```go
func (m *MockApplication) EXPECT() *MockApplicationMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use.

#### func (*MockApplication) Init

```go
func (m *MockApplication) Init()
```
Init mocks base method.

#### func (*MockApplication) IsDebug

```go
func (m *MockApplication) IsDebug() bool
```
IsDebug mocks base method.

#### func (*MockApplication) IsRunning

```go
func (m *MockApplication) IsRunning() bool
```
IsRunning mocks base method.

#### func (*MockApplication) Logger

```go
func (m *MockApplication) Logger() logger.Logger
```
Logger mocks base method.

#### func (*MockApplication) Name

```go
func (m *MockApplication) Name() string
```
Name mocks base method.

#### func (*MockApplication) ProcessManager

```go
func (m *MockApplication) ProcessManager() process.Manager
```
ProcessManager mocks base method.

#### func (*MockApplication) Run

```go
func (m *MockApplication) Run()
```
Run mocks base method.

#### func (*MockApplication) SetMainLoop

```go
func (m *MockApplication) SetMainLoop(arg0 MainLoopFn)
```
SetMainLoop mocks base method.

#### func (*MockApplication) SetOnCHLD

```go
func (m *MockApplication) SetOnCHLD(arg0 OnSignalFn)
```
SetOnCHLD mocks base method.

#### func (*MockApplication) SetOnExit

```go
func (m *MockApplication) SetOnExit(arg0 OnSignalFn)
```
SetOnExit mocks base method.

#### func (*MockApplication) SetOnHUP

```go
func (m *MockApplication) SetOnHUP(arg0 OnSignalFn)
```
SetOnHUP mocks base method.

#### func (*MockApplication) SetOnStart

```go
func (m *MockApplication) SetOnStart(arg0 OnSignalFn)
```
SetOnStart mocks base method.

#### func (*MockApplication) SetOnUSR1

```go
func (m *MockApplication) SetOnUSR1(arg0 OnSignalFn)
```
SetOnUSR1 mocks base method.

#### func (*MockApplication) SetOnUSR2

```go
func (m *MockApplication) SetOnUSR2(arg0 OnSignalFn)
```
SetOnUSR2 mocks base method.

#### func (*MockApplication) SetOnWINCH

```go
func (m *MockApplication) SetOnWINCH(arg0 OnSignalFn)
```
SetOnWINCH mocks base method.

#### func (*MockApplication) Terminate

```go
func (m *MockApplication) Terminate()
```
Terminate mocks base method.

#### func (*MockApplication) Version

```go
func (m *MockApplication) Version() *semver.SemVer
```
Version mocks base method.

#### type MockApplicationMockRecorder

```go
type MockApplicationMockRecorder struct {
}
```

MockApplicationMockRecorder is the mock recorder for MockApplication.

#### func (*MockApplicationMockRecorder) Commit

```go
func (mr *MockApplicationMockRecorder) Commit() *gomock.Call
```
Commit indicates an expected call of Commit.

#### func (*MockApplicationMockRecorder) Configuration

```go
func (mr *MockApplicationMockRecorder) Configuration() *gomock.Call
```
Configuration indicates an expected call of Configuration.

#### func (*MockApplicationMockRecorder) Context

```go
func (mr *MockApplicationMockRecorder) Context() *gomock.Call
```
Context indicates an expected call of Context.

#### func (*MockApplicationMockRecorder) Init

```go
func (mr *MockApplicationMockRecorder) Init() *gomock.Call
```
Init indicates an expected call of Init.

#### func (*MockApplicationMockRecorder) IsDebug

```go
func (mr *MockApplicationMockRecorder) IsDebug() *gomock.Call
```
IsDebug indicates an expected call of IsDebug.

#### func (*MockApplicationMockRecorder) IsRunning

```go
func (mr *MockApplicationMockRecorder) IsRunning() *gomock.Call
```
IsRunning indicates an expected call of IsRunning.

#### func (*MockApplicationMockRecorder) Logger

```go
func (mr *MockApplicationMockRecorder) Logger() *gomock.Call
```
Logger indicates an expected call of Logger.

#### func (*MockApplicationMockRecorder) Name

```go
func (mr *MockApplicationMockRecorder) Name() *gomock.Call
```
Name indicates an expected call of Name.

#### func (*MockApplicationMockRecorder) ProcessManager

```go
func (mr *MockApplicationMockRecorder) ProcessManager() *gomock.Call
```
ProcessManager indicates an expected call of ProcessManager.

#### func (*MockApplicationMockRecorder) Run

```go
func (mr *MockApplicationMockRecorder) Run() *gomock.Call
```
Run indicates an expected call of Run.

#### func (*MockApplicationMockRecorder) SetMainLoop

```go
func (mr *MockApplicationMockRecorder) SetMainLoop(arg0 any) *gomock.Call
```
SetMainLoop indicates an expected call of SetMainLoop.

#### func (*MockApplicationMockRecorder) SetOnCHLD

```go
func (mr *MockApplicationMockRecorder) SetOnCHLD(arg0 any) *gomock.Call
```
SetOnCHLD indicates an expected call of SetOnCHLD.

#### func (*MockApplicationMockRecorder) SetOnExit

```go
func (mr *MockApplicationMockRecorder) SetOnExit(arg0 any) *gomock.Call
```
SetOnExit indicates an expected call of SetOnExit.

#### func (*MockApplicationMockRecorder) SetOnHUP

```go
func (mr *MockApplicationMockRecorder) SetOnHUP(arg0 any) *gomock.Call
```
SetOnHUP indicates an expected call of SetOnHUP.

#### func (*MockApplicationMockRecorder) SetOnStart

```go
func (mr *MockApplicationMockRecorder) SetOnStart(arg0 any) *gomock.Call
```
SetOnStart indicates an expected call of SetOnStart.

#### func (*MockApplicationMockRecorder) SetOnUSR1

```go
func (mr *MockApplicationMockRecorder) SetOnUSR1(arg0 any) *gomock.Call
```
SetOnUSR1 indicates an expected call of SetOnUSR1.

#### func (*MockApplicationMockRecorder) SetOnUSR2

```go
func (mr *MockApplicationMockRecorder) SetOnUSR2(arg0 any) *gomock.Call
```
SetOnUSR2 indicates an expected call of SetOnUSR2.

#### func (*MockApplicationMockRecorder) SetOnWINCH

```go
func (mr *MockApplicationMockRecorder) SetOnWINCH(arg0 any) *gomock.Call
```
SetOnWINCH indicates an expected call of SetOnWINCH.

#### func (*MockApplicationMockRecorder) Terminate

```go
func (mr *MockApplicationMockRecorder) Terminate() *gomock.Call
```
Terminate indicates an expected call of Terminate.

#### func (*MockApplicationMockRecorder) Version

```go
func (mr *MockApplicationMockRecorder) Version() *gomock.Call
```
Version indicates an expected call of Version.

#### type OnSignalFn

```go
type OnSignalFn func(Application)
```

Signal callback function type.
