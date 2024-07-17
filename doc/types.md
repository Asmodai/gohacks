-*- Mode: gfm -*-

# types -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/types"
```

## Usage

```go
const (
	// CSI Pm m -- Character Attributes (SGR).
	NORMAL        = 0 // Normal (default), VT100.
	BOLD          = 1 // Bold, VT100.
	FAINT         = 2 // Faint, decreased intensity, ECMA-48 2nd.
	ITALICS       = 3 // Italicized, ECMA-48 2nd.
	UNDERLINE     = 4 // Underlined, VT100.
	BLINK         = 5 // Blink, VT100.
	INVERSE       = 7 // Inverse, VT100.
	STRIKETHROUGH = 9 // Crossed-out characters, ECMA-48 3rd.

	BLACK   = 0
	RED     = 1
	GREEN   = 2
	YELLOW  = 3
	BLUE    = 4
	MAGENTA = 5
	CYAN    = 6
	WHITE   = 7
	DEFAULT = 9

	FGOFFSET = 30
	BGOFFSET = 40
)
```

```go
const (
	// Amount of time to delay semaphore acquisition loops.
	MailboxDelaySleep  time.Duration = 50 * time.Millisecond
	DefaultCtxDeadline time.Duration = 5 * time.Second
)
```

#### type ColorString

```go
type ColorString struct {
}
```


#### func  MakeColorString

```go
func MakeColorString() *ColorString
```

#### func  MakeColorStringWithAttrs

```go
func MakeColorStringWithAttrs(data string, attr, foreg, backg int) *ColorString
```

#### func (*ColorString) SetAttr

```go
func (cs *ColorString) SetAttr(attr int)
```

#### func (*ColorString) SetBG

```go
func (cs *ColorString) SetBG(col int)
```

#### func (*ColorString) SetFG

```go
func (cs *ColorString) SetFG(col int)
```

#### func (*ColorString) SetString

```go
func (cs *ColorString) SetString(val string)
```

#### func (*ColorString) String

```go
func (cs *ColorString) String() string
```

#### type Error

```go
type Error struct {
	Module  string
	Message string
}
```

Custom error structure.

This is compatible with the `error` interface and provides `Unwrap` support.

#### func  NewError

```go
func NewError(module string, format string, args ...interface{}) *Error
```
Create a new error object.

#### func  NewErrorAndLog

```go
func NewErrorAndLog(module string, format string, args ...interface{}) *Error
```
Create a new error object and immediately log it.

#### func (*Error) Error

```go
func (e *Error) Error() string
```
Return a human-readable string representation of the error.

#### func (*Error) Log

```go
func (e *Error) Log()
```
Log the error.

#### func (*Error) MarshalJSON

```go
func (e *Error) MarshalJSON() ([]byte, error)
```
Convert the error to a JSON string.

#### func (*Error) String

```go
func (e *Error) String() string
```
Return a human-readable string representation of the error.

#### func (*Error) Unwrap

```go
func (e *Error) Unwrap() error
```
Unwrap the error.

#### type Mailbox

```go
type Mailbox struct {
}
```

Mailbox structure.

This is a cheap implementation of a mailbox.

It uses two semaphores to control read and write access, and contains a single
datum.

This is *not* a queue!

#### func  NewMailbox

```go
func NewMailbox() *Mailbox
```
Create and return a new empty mailbox.

Note: this acquires the `preventRead` semaphore.

#### func (*Mailbox) Full

```go
func (m *Mailbox) Full() bool
```
Does the mailbox contain a value?

#### func (*Mailbox) Get

```go
func (m *Mailbox) Get() (interface{}, bool)
```
Get an element from the mailbox. Defaults to using a context with a deadline of
5 seconds.

#### func (*Mailbox) GetWithContext

```go
func (m *Mailbox) GetWithContext(ctx context.Context) (interface{}, bool)
```

#### func (*Mailbox) Put

```go
func (m *Mailbox) Put(elem interface{})
```
Put an element into the mailbox.

#### type Pair

```go
type Pair struct {
	First  any
	Second any
}
```

Pair structure.

This is a cheap implementation of a pair (aka two-value tuple).

#### func  NewEmptyPair

```go
func NewEmptyPair() *Pair
```
Create a new empty pair.

#### func  NewPair

```go
func NewPair(first any, second any) *Pair
```
Create a new pair.

#### func (*Pair) String

```go
func (p *Pair) String() string
```
Return a string representation of the pair.

#### type Queue

```go
type Queue struct {
	sync.Mutex
}
```

Queue structure.

This is a cheap implementation of a LIFO queue.

#### func  NewBoundedQueue

```go
func NewBoundedQueue(bounds int) *Queue
```

#### func  NewQueue

```go
func NewQueue() *Queue
```
Create a new empty queue.

#### func (*Queue) Full

```go
func (q *Queue) Full() bool
```
Is the queue full?

#### func (*Queue) Get

```go
func (q *Queue) Get() (interface{}, bool)
```
Remove an element from the end of the queue and return it.

#### func (*Queue) Len

```go
func (q *Queue) Len() int
```
Return the number of elements in the queue.

#### func (*Queue) Put

```go
func (q *Queue) Put(elem interface{}) bool
```
Append an element to the queue. Returns `false` if there is no more room in the
queue.
