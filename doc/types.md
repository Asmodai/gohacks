<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# types -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/types"
```

## Usage

```go
const (
	Black   Colour = '0' // Black colour -- ANSI.
	Red     Colour = '1' // Red colour -- ANSI.
	Green   Colour = '2' // Green colour -- ANSI.
	Yellow  Colour = '3' // Yellow colour -- ANSI.
	Blue    Colour = '4' // Blue colour -- ANSI.
	Magenta Colour = '5' // Magenta colour -- ANSI.
	Cyan    Colour = '6' // Cyan colour -- ANSI.
	White   Colour = '7' // White colour -- ANSI.
	Default Colour = '9' // Default colour -- ANSI.

	Normal     int = 0 // Reset all attributes.
	Bold       int = 1 // Bold -- VT100.
	Faint      int = 2 // Faint, decreased intensity -- ECMA-48 2e.
	Italic     int = 3 // Italicizsed -- ECMA-48 2e.
	Underline  int = 4 // Underlined -- VT100.
	Blink      int = 5 // Blinking -- VT100.
	Inverse    int = 6 // Inverse video -- VT100.
	Strikethru int = 7 // Crossed-out characters -- ECMA-48 3e.

)
```
CSI Pm [; Pm ...] m -- Character Attributes (SGR).

```go
const (
	// Amount of time to delay semaphore acquisition loops.
	MailboxDelaySleep time.Duration = 50 * time.Millisecond

	// Default deadline for context timeouts.
	DefaultCtxDeadline time.Duration = 5 * time.Second
)
```

```go
var (
	// Error condition that signals an invalid time duration of some kind.
	//
	// This error is usually wrapped around a descriptive message string.
	ErrInvalidDuration error = errors.Base("invalid time duration")

	// Error condition that signals that a duration is not a string value.
	//
	// This error is used by `Set` as well as JSON and YAML methods.
	ErrDurationNotString error = errors.Base("duration must be a string")

	// Error condition that signals that a duration is out of bounds.
	//
	// This is used by `Validate`.
	ErrOutOfBounds error = errors.Base("duration out of bounds")
)
```

#### type ColorString

```go
type ColorString struct {
}
```

Colour string.

Generates a string that, with the right terminal type, display text using
various character attributes.

To find out more, consult your nearest DEC VT340 programmer's manual or the
latest ECMA-48 standard.

#### func  NewColorString

```go
func NewColorString() *ColorString
```
Make a new coloured string.

#### func  NewColorStringWithColors

```go
func NewColorStringWithColors(
	data string,
	foreg, backg Colour,
) *ColorString
```
Make a new coloured string with the given attributes.

#### func (*ColorString) AddAttr

```go
func (cs *ColorString) AddAttr(index int) *ColorString
```
Add a specific attribute.

This ignores the "Normal" attribute. To clear attributes, use `SetNormal` or
`Clear`.

#### func (*ColorString) Clear

```go
func (cs *ColorString) Clear() *ColorString
```
Reset all attributes.

#### func (*ColorString) RemoveAttr

```go
func (cs *ColorString) RemoveAttr(index int) *ColorString
```
Remove a specific attribute.

This ignores the "Normal" attribute. To clear attributes, use `SetNormal` or
`Clear`.

#### func (*ColorString) SetBG

```go
func (cs *ColorString) SetBG(col Colour) *ColorString
```
Set the background colour.

#### func (*ColorString) SetBlink

```go
func (cs *ColorString) SetBlink() *ColorString
```
Set the blink attribute.

#### func (*ColorString) SetBold

```go
func (cs *ColorString) SetBold() *ColorString
```
Set the bold attribute.

#### func (*ColorString) SetFG

```go
func (cs *ColorString) SetFG(col Colour) *ColorString
```
Set the foreground colour.

#### func (*ColorString) SetFaint

```go
func (cs *ColorString) SetFaint() *ColorString
```
Set the faint attribute.

#### func (*ColorString) SetInverse

```go
func (cs *ColorString) SetInverse() *ColorString
```
Set the inverse video attribute.

#### func (*ColorString) SetItalic

```go
func (cs *ColorString) SetItalic() *ColorString
```
Set the italic attribute.

#### func (*ColorString) SetNormal

```go
func (cs *ColorString) SetNormal() *ColorString
```
Set the attribute to normal.

This removes all other attributes.

#### func (*ColorString) SetStrikethru

```go
func (cs *ColorString) SetStrikethru() *ColorString
```
Set the strikethrough attribute.

#### func (*ColorString) SetString

```go
func (cs *ColorString) SetString(val string) *ColorString
```
Set the string to display.

#### func (*ColorString) SetUnderline

```go
func (cs *ColorString) SetUnderline() *ColorString
```
Set the underline attribute.

#### func (*ColorString) String

```go
func (cs *ColorString) String() string
```
Convert to a string.

Warning: the resulting string will contain escape sequences for use with a
compliant terminal or terminal emulator.

#### type Colour

```go
type Colour rune
```

ECMA-48 colour descriptor type.

#### type Duration

```go
type Duration time.Duration
```

Enhanced time duration type.

#### func (Duration) Duration

```go
func (obj Duration) Duration() time.Duration
```
Coerce a duration to a `time.Duration` value.

#### func (Duration) MarshalJSON

```go
func (obj Duration) MarshalJSON() ([]byte, error)
```
JSON marshalling method.

#### func (Duration) MarshalYAML

```go
func (obj Duration) MarshalYAML() (any, error)
```
YAML marshalling method.

#### func (*Duration) Set

```go
func (obj *Duration) Set(str string) error
```
Set the duration to that of the given string.

This method uses `time.ParseDuration`, so any string that `time` understands may
be used.

If the string value fails parsing, then `ErrInvalidDuration` is returned.

#### func (*Duration) SetYAML

```go
func (obj *Duration) SetYAML(value any) error
```
Set the duration value from a YAML value.

If the passed YAML value is not a string, then `ErrDurationNotString` is
returned.

Will also return any error condition from the `Set` method.

#### func (Duration) String

```go
func (obj Duration) String() string
```
Coerce a duration to a string value.

#### func (Duration) Type

```go
func (obj Duration) Type() string
```
Return the data type name for CLI flag parsing purposes.

#### func (*Duration) UnmarshalJSON

```go
func (obj *Duration) UnmarshalJSON(data []byte) error
```
JSON unmarshalling method.

#### func (*Duration) UnmarshalYAML

```go
func (obj *Duration) UnmarshalYAML(value *yaml.Node) error
```
YAML unmarshalling method.

#### func (Duration) Validate

```go
func (obj Duration) Validate(minDuration, maxDuration time.Duration) error
```
Validate a duration.

This ensures a duration is within a given range.

If validation fails, then `ErrOutOfBounds` is returned.

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
func (m *Mailbox) Get() (any, bool)
```
Get an element from the mailbox. Defaults to using a context with a deadline of
5 seconds.

#### func (*Mailbox) GetWithContext

```go
func (m *Mailbox) GetWithContext(ctx context.Context) (any, bool)
```
Get an element from the mailbox using the provided context.

It is recommended to use a context that has a timeout deadline.

#### func (*Mailbox) Put

```go
func (m *Mailbox) Put(elem any)
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
Create a queue that is bounded to a specific size.

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
func (q *Queue) Get() (any, bool)
```
Remove an element from the end of the queue and return it.

#### func (*Queue) Len

```go
func (q *Queue) Len() int
```
Return the number of elements in the queue.

#### func (*Queue) Put

```go
func (q *Queue) Put(elem any) bool
```
Append an element to the queue. Returns `false` if there is no more room in the
queue.
