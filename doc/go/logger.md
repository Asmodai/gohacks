<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# logger -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/logger"
```

## Usage

```go
const ContextKeyLogger = "gohacks/logger@v1"
```
Key used to store the instance in the context's user value.

```go
var ErrValueNotLogger = errors.Base("value is not Logger")
```
Signalled if the instance associated with the context key is not of type Logger.

#### func  SetLogger

```go
func SetLogger(ctx context.Context, inst Logger) (context.Context, error)
```
Set Logger stores the instance in the context map.

#### func  SetLoggerIfAbsent

```go
func SetLoggerIfAbsent(ctx context.Context, inst Logger) (context.Context, error)
```
SetLoggerIfAbsent sets only if not already present.

#### func  WithLogger

```go
func WithLogger(ctx context.Context, fn func(Logger))
```
WithLogger calls fn with the instance or fallback.

#### type Fields

```go
type Fields map[string]any
```

Log field map type.

#### type Level

```go
type Level int8
```


```go
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)
```

#### type Logger

```go
type Logger interface {
	// Sync a log file to disk, if available.
	Sync() error

	// Set whether the logger prints in human-readable 'debug' output or
	// machine-readable JSON format.
	SetDebug(bool)

	// Set the log file to which the logger will write.
	SetLogFile(string)

	// Set sampling options.
	//
	// This is only used by the Zap logger.
	SetSampling(initial, threshold int)

	// Set logger level.
	SetLevel(Level)

	// Log a Go error.
	GoError(error, ...any)

	// Log a debug message.
	Debug(string, ...any)

	// Log an error message.
	Error(string, ...any)

	// Log a warning message.
	Warn(string, ...any)

	// Log an informational message.
	Info(string, ...any)

	// Log a fatal error condition and exit.
	Fatal(string, ...any)

	// Log a panic condition and exit.
	Panic(string, ...any)

	// Log a debug message with printf formatting.
	Debugf(string, ...any)

	// Log an error message with printf formatting.
	Errorf(string, ...any)

	// Log a warning message with printf formatting.
	Warnf(string, ...any)

	// Log an informational message with printf formatting.
	Infof(string, ...any)

	// Log a fatal message with printf formatting.
	Fatalf(string, ...any)

	// Log a panic message with printf formatting.
	Panicf(string, ...any)

	// Encapsulate user-specified metadata fields.
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

#### func  FromLogger

```go
func FromLogger(ctx context.Context) Logger
```
FromLogger returns the instance or the fallback.

#### func  GetLogger

```go
func GetLogger(ctx context.Context) (Logger, error)
```
Get the logger from the given context.

Will return ErrValueNotLogger if the value in the context is not of type Logger.

#### func  MustGetLogger

```go
func MustGetLogger(ctx context.Context) Logger
```
Attempt to get the instance from the given context. Panics if the operation
fails.

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

#### func  TryGetLogger

```go
func TryGetLogger(ctx context.Context) (Logger, bool)
```
TryGetLogger returns the instance and true if present and typed.
