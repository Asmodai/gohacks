-*- Mode: gfm -*-

# debug -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/debug"
```

## Usage

```go
const (
	SpacesPerIndent int = 4  // Number of spaces to use for indentation.
	LineFullLength  int = 69 // Length of debug info line.
	LineTitleLength int = 64 // Length of title segment.
)
```

#### func  DebugPrint

```go
func DebugPrint(thing any, params ...any)
```
Print out a `Debugable` thing to standard output.

#### func  DebugString

```go
func DebugString(thing any, params ...any) (string, bool)
```
Return the string value of a `Debugable` thing.

#### type Debug

```go
type Debug struct {
}
```

Debug information.

#### func  NewDebug

```go
func NewDebug(title string) *Debug
```
Create a new debug object.

#### func (*Debug) End

```go
func (obj *Debug) End()
```
Finalise the debug object.

#### func (*Debug) Init

```go
func (obj *Debug) Init(params ...any)
```
Initialise the debug object.

#### func (*Debug) Print

```go
func (obj *Debug) Print()
```
Print out the debug information to standard output.

#### func (*Debug) Printf

```go
func (obj *Debug) Printf(format string, args ...any)
```
Print to the debug information.

#### func (*Debug) String

```go
func (obj *Debug) String() string
```
Return the string representation of the debug object.

#### type Debugable

```go
type Debugable interface {
	Debug(...any) *Debug
}
```

Abstract interface that defines whether an object is debuggable.
