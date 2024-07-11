-*- Mode: gfm -*-

# logger -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/logger"
```

## Usage

#### type Fields

```go
type Fields map[string]any
```


#### type Logger

```go
type Logger interface {
	SetDebug(bool)
	SetLogFile(string)

	Debug(string, ...any)
	Error(string, ...any)
	Warn(string, ...any)
	Info(string, ...any)
	Fatal(string, ...any)
	Panic(string, ...any)

	Debugf(string, ...any)
	Errorf(string, ...any)
	Warnf(string, ...any)
	Infof(string, ...any)
	Fatalf(string, ...any)
	Panicf(string, ...any)

	WithFields(Fields) Logger
}
```

Logging structure.

To use,

1) Create a logger:

```go

    lgr := logger.NewLogger()

```

2) Do things with it:

```go

    lgr.Warn("Not enough coffee!")
    lgr.Info("Water is heating up.")
    // and so on.

```

If an empty string is passed to `NewLogger`, then the log facility will display
messages on standard output.

#### func  NewDefaultLogger

```go
func NewDefaultLogger() Logger
```
Create a new default logger.

#### func  NewZapLogger

```go
func NewZapLogger() Logger
```
Create a new logger.

#### func  NewZapLoggerWithFile

```go
func NewZapLoggerWithFile(logfile string) Logger
```
Create a new logger with the given log file.

#### type MockLogger

```go
type MockLogger struct {
}
```

MockLogger is a mock of Logger interface.

#### func  NewMockLogger

```go
func NewMockLogger(ctrl *gomock.Controller) *MockLogger
```
NewMockLogger creates a new mock instance.

#### func (*MockLogger) Debug

```go
func (m *MockLogger) Debug(arg0 string, arg1 ...any)
```
Debug mocks base method.

#### func (*MockLogger) Debugf

```go
func (m *MockLogger) Debugf(arg0 string, arg1 ...any)
```
Debugf mocks base method.

#### func (*MockLogger) EXPECT

```go
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use.

#### func (*MockLogger) Error

```go
func (m *MockLogger) Error(arg0 string, arg1 ...any)
```
Error mocks base method.

#### func (*MockLogger) Errorf

```go
func (m *MockLogger) Errorf(arg0 string, arg1 ...any)
```
Errorf mocks base method.

#### func (*MockLogger) Fatal

```go
func (m *MockLogger) Fatal(arg0 string, arg1 ...any)
```
Fatal mocks base method.

#### func (*MockLogger) Fatalf

```go
func (m *MockLogger) Fatalf(arg0 string, arg1 ...any)
```
Fatalf mocks base method.

#### func (*MockLogger) Info

```go
func (m *MockLogger) Info(arg0 string, arg1 ...any)
```
Info mocks base method.

#### func (*MockLogger) Infof

```go
func (m *MockLogger) Infof(arg0 string, arg1 ...any)
```
Infof mocks base method.

#### func (*MockLogger) Panic

```go
func (m *MockLogger) Panic(arg0 string, arg1 ...any)
```
Panic mocks base method.

#### func (*MockLogger) Panicf

```go
func (m *MockLogger) Panicf(arg0 string, arg1 ...any)
```
Panicf mocks base method.

#### func (*MockLogger) SetDebug

```go
func (m *MockLogger) SetDebug(arg0 bool)
```
SetDebug mocks base method.

#### func (*MockLogger) SetLogFile

```go
func (m *MockLogger) SetLogFile(arg0 string)
```
SetLogFile mocks base method.

#### func (*MockLogger) Warn

```go
func (m *MockLogger) Warn(arg0 string, arg1 ...any)
```
Warn mocks base method.

#### func (*MockLogger) Warnf

```go
func (m *MockLogger) Warnf(arg0 string, arg1 ...any)
```
Warnf mocks base method.

#### func (*MockLogger) WithFields

```go
func (m *MockLogger) WithFields(arg0 Fields) Logger
```
WithFields mocks base method.

#### type MockLoggerMockRecorder

```go
type MockLoggerMockRecorder struct {
}
```

MockLoggerMockRecorder is the mock recorder for MockLogger.

#### func (*MockLoggerMockRecorder) Debug

```go
func (mr *MockLoggerMockRecorder) Debug(arg0 any, arg1 ...any) *gomock.Call
```
Debug indicates an expected call of Debug.

#### func (*MockLoggerMockRecorder) Debugf

```go
func (mr *MockLoggerMockRecorder) Debugf(arg0 any, arg1 ...any) *gomock.Call
```
Debugf indicates an expected call of Debugf.

#### func (*MockLoggerMockRecorder) Error

```go
func (mr *MockLoggerMockRecorder) Error(arg0 any, arg1 ...any) *gomock.Call
```
Error indicates an expected call of Error.

#### func (*MockLoggerMockRecorder) Errorf

```go
func (mr *MockLoggerMockRecorder) Errorf(arg0 any, arg1 ...any) *gomock.Call
```
Errorf indicates an expected call of Errorf.

#### func (*MockLoggerMockRecorder) Fatal

```go
func (mr *MockLoggerMockRecorder) Fatal(arg0 any, arg1 ...any) *gomock.Call
```
Fatal indicates an expected call of Fatal.

#### func (*MockLoggerMockRecorder) Fatalf

```go
func (mr *MockLoggerMockRecorder) Fatalf(arg0 any, arg1 ...any) *gomock.Call
```
Fatalf indicates an expected call of Fatalf.

#### func (*MockLoggerMockRecorder) Info

```go
func (mr *MockLoggerMockRecorder) Info(arg0 any, arg1 ...any) *gomock.Call
```
Info indicates an expected call of Info.

#### func (*MockLoggerMockRecorder) Infof

```go
func (mr *MockLoggerMockRecorder) Infof(arg0 any, arg1 ...any) *gomock.Call
```
Infof indicates an expected call of Infof.

#### func (*MockLoggerMockRecorder) Panic

```go
func (mr *MockLoggerMockRecorder) Panic(arg0 any, arg1 ...any) *gomock.Call
```
Panic indicates an expected call of Panic.

#### func (*MockLoggerMockRecorder) Panicf

```go
func (mr *MockLoggerMockRecorder) Panicf(arg0 any, arg1 ...any) *gomock.Call
```
Panicf indicates an expected call of Panicf.

#### func (*MockLoggerMockRecorder) SetDebug

```go
func (mr *MockLoggerMockRecorder) SetDebug(arg0 any) *gomock.Call
```
SetDebug indicates an expected call of SetDebug.

#### func (*MockLoggerMockRecorder) SetLogFile

```go
func (mr *MockLoggerMockRecorder) SetLogFile(arg0 any) *gomock.Call
```
SetLogFile indicates an expected call of SetLogFile.

#### func (*MockLoggerMockRecorder) Warn

```go
func (mr *MockLoggerMockRecorder) Warn(arg0 any, arg1 ...any) *gomock.Call
```
Warn indicates an expected call of Warn.

#### func (*MockLoggerMockRecorder) Warnf

```go
func (mr *MockLoggerMockRecorder) Warnf(arg0 any, arg1 ...any) *gomock.Call
```
Warnf indicates an expected call of Warnf.

#### func (*MockLoggerMockRecorder) WithFields

```go
func (mr *MockLoggerMockRecorder) WithFields(arg0 any) *gomock.Call
```
WithFields indicates an expected call of WithFields.
