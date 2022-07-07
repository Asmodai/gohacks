-*- Mode: gfm -*-

# logger -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/logger"
```

## Usage

#### type DefaultLogger

```go
type DefaultLogger struct {
}
```

Default logging structure.

This is a simple implementation of the `ILogger` interface that simply redirects
messages to `log.Printf`.

It is used in the same way as the main `Logger` implementation.

#### func  NewDefaultLogger

```go
func NewDefaultLogger() *DefaultLogger
```
Create a new default logger.

#### func (*DefaultLogger) Debug

```go
func (l *DefaultLogger) Debug(msg string, rest ...interface{})
```
Write a debug message to the log.

#### func (*DefaultLogger) Debugf

```go
func (l *DefaultLogger) Debugf(format string, args ...interface{})
```
Write a debug message to the log.

#### func (*DefaultLogger) Error

```go
func (l *DefaultLogger) Error(msg string, rest ...interface{})
```
Write an error message to the log.

#### func (*DefaultLogger) Errorf

```go
func (l *DefaultLogger) Errorf(msg string, args ...interface{})
```
Write an error message to the log and then exit.

#### func (*DefaultLogger) Fatal

```go
func (l *DefaultLogger) Fatal(msg string, rest ...interface{})
```
Write a fatal message to the log and then exit.

#### func (*DefaultLogger) Fatalf

```go
func (l *DefaultLogger) Fatalf(msg string, args ...interface{})
```
Write a fatal message to the log and then exit.

#### func (*DefaultLogger) Info

```go
func (l *DefaultLogger) Info(msg string, rest ...interface{})
```
Write an information message to the log.

#### func (*DefaultLogger) Infof

```go
func (l *DefaultLogger) Infof(msg string, args ...interface{})
```
Write an information message to the log.

#### func (*DefaultLogger) Panicf

```go
func (l *DefaultLogger) Panicf(msg string, args ...interface{})
```
Write a fatal message to the log and then exit.

#### func (*DefaultLogger) SetDebug

```go
func (l *DefaultLogger) SetDebug(junk bool)
```
Set debug mode.

#### func (*DefaultLogger) SetLogFile

```go
func (l *DefaultLogger) SetLogFile(junk string)
```
Set the log file to use.

#### func (*DefaultLogger) Warn

```go
func (l *DefaultLogger) Warn(msg string, rest ...interface{})
```
Write a warning message to the log.

#### func (*DefaultLogger) Warnf

```go
func (l *DefaultLogger) Warnf(msg string, args ...interface{})
```
Write a warning message to the log.

#### func (*DefaultLogger) WithFields

```go
func (l *DefaultLogger) WithFields(_ Fields) ILogger
```

#### type Fields

```go
type Fields map[string]interface{}
```


#### type ILogger

```go
type ILogger interface {
	SetDebug(bool)
	SetLogFile(string)

	Debug(string, ...interface{})
	Error(string, ...interface{})
	Warn(string, ...interface{})
	Info(string, ...interface{})
	Fatal(string, ...interface{})

	Debugf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Fatalf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})

	WithFields(Fields) ILogger
}
```


#### type Logger

```go
type Logger struct {
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

#### func  NewLogger

```go
func NewLogger() *Logger
```
Create a new logger.

#### func  NewLoggerWithFile

```go
func NewLoggerWithFile(logfile string) *Logger
```
Create a new logger with the given log file.

#### func (*Logger) Debug

```go
func (l *Logger) Debug(msg string, rest ...interface{})
```
Write a debug message to the log.

#### func (*Logger) Debugf

```go
func (l *Logger) Debugf(format string, args ...interface{})
```
Compatibility method.

#### func (*Logger) Error

```go
func (l *Logger) Error(msg string, rest ...interface{})
```
Write an error message to the log.

#### func (*Logger) Errorf

```go
func (l *Logger) Errorf(format string, args ...interface{})
```
Compatibility method.

#### func (*Logger) Fatal

```go
func (l *Logger) Fatal(msg string, rest ...interface{})
```
Write a fatal message to the log and then exit.

#### func (*Logger) Fatalf

```go
func (l *Logger) Fatalf(format string, args ...interface{})
```
Compatibility method.

#### func (*Logger) Info

```go
func (l *Logger) Info(msg string, rest ...interface{})
```
Write an information message to the log.

#### func (*Logger) Infof

```go
func (l *Logger) Infof(format string, args ...interface{})
```
Compatibility method.

#### func (*Logger) Panicf

```go
func (l *Logger) Panicf(format string, args ...interface{})
```
Compatibility method.

#### func (*Logger) SetDebug

```go
func (l *Logger) SetDebug(flag bool)
```
Set debug mode.

Debug mode is a production-friendly runtime mode that will print human-readable
messages to standard output instead of the defined log file.

#### func (*Logger) SetLogFile

```go
func (l *Logger) SetLogFile(file string)
```
Set the log file to use.

#### func (*Logger) Warn

```go
func (l *Logger) Warn(msg string, rest ...interface{})
```
Write a warning message to the log.

#### func (*Logger) Warnf

```go
func (l *Logger) Warnf(format string, args ...interface{})
```
Compatibility method.

#### func (*Logger) WithFields

```go
func (l *Logger) WithFields(fields Fields) ILogger
```

#### type MockLogger

```go
type MockLogger struct {
	Test      *testing.T
	LastFatal string
}
```

Mock logger for Go testing framework.

To use, be sure to set `Test` to your test's `testing.T` instance.

#### func  NewMockLogger

```go
func NewMockLogger(_ string) *MockLogger
```
Create a new mock logger.

#### func (*MockLogger) Debug

```go
func (l *MockLogger) Debug(msg string, rest ...interface{})
```
Write a debug message to the log.

#### func (*MockLogger) Debugf

```go
func (l *MockLogger) Debugf(msg string, rest ...interface{})
```
Write a debug message to the log.

#### func (*MockLogger) Error

```go
func (l *MockLogger) Error(msg string, rest ...interface{})
```
Write an error message to the log.

#### func (*MockLogger) Errorf

```go
func (l *MockLogger) Errorf(msg string, rest ...interface{})
```
Write a fatal message to the log and then exit.

#### func (*MockLogger) Fatal

```go
func (l *MockLogger) Fatal(msg string, rest ...interface{})
```
Write a fatal message to the log and then exit.

#### func (*MockLogger) Fatalf

```go
func (l *MockLogger) Fatalf(msg string, rest ...interface{})
```
Write a fatal message to the log and then exit.

#### func (*MockLogger) Info

```go
func (l *MockLogger) Info(msg string, rest ...interface{})
```
Write an information message to the log.

#### func (*MockLogger) Infof

```go
func (l *MockLogger) Infof(msg string, rest ...interface{})
```
Write an information message to the log.

#### func (*MockLogger) Panicf

```go
func (l *MockLogger) Panicf(msg string, rest ...interface{})
```
Write a fatal message to the log and then exit.

#### func (*MockLogger) SetDebug

```go
func (l *MockLogger) SetDebug(junk bool)
```
Set debug mode.

#### func (*MockLogger) SetLogFile

```go
func (l *MockLogger) SetLogFile(junk string)
```
Set the log file to use.

#### func (*MockLogger) Warn

```go
func (l *MockLogger) Warn(msg string, rest ...interface{})
```
Write a warning message to the log.

#### func (*MockLogger) Warnf

```go
func (l *MockLogger) Warnf(msg string, rest ...interface{})
```
Write a warning message to the log.

#### func (*MockLogger) WithFields

```go
func (l *MockLogger) WithFields(_ Fields) ILogger
```
