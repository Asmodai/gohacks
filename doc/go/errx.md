<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# errx -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/errx"
```

## Usage

```go
var (
	AllDetails   = tozd.AllDetails
	As           = tozd.As
	Base         = tozd.Base
	BaseWrap     = tozd.BaseWrap
	BaseWrapf    = tozd.BaseWrapf
	Basef        = tozd.Basef
	Cause        = tozd.Cause
	Details      = tozd.Details
	Errorf       = tozd.Errorf
	Is           = tozd.Is
	Join         = tozd.Join
	New          = tozd.New
	Prefix       = tozd.Prefix
	Unjoin       = tozd.Unjoin
	Unwrap       = tozd.Unwrap
	WithDetails  = tozd.WithDetails
	WithMessage  = tozd.WithMessage
	WithMessagef = tozd.WithMessagef
	WithStack    = tozd.WithStack
	Wrap         = tozd.Wrap
	WrapWith     = tozd.WrapWith
	Wrapf        = tozd.Wrapf
)
```
We simply import the interesting symbols that we want from our given error
handling package.

In this case, the error handling package we use is gitlab.com/tozd/go/errors

#### type Causer

```go
type Causer interface {
	Cause() error
}
```

Error cause interface.

#### type Detailer

```go
type Detailer interface {
	Details() map[string]any
}
```

Error details interface.

#### type Error

```go
type Error tozd.E
```


#### type StackTracer

```go
type StackTracer interface {
	StackTrace() []uintptr
}
```

Stack trace interface.

#### type Unwrapper

```go
type Unwrapper interface {
	Unwrap() error
}
```

Error unwrap interface.

#### type UnwrapperJoined

```go
type UnwrapperJoined interface {
	Unwrap() []error
}
```

Joined error unwrap interface.
