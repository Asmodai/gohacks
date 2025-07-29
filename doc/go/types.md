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

```go
var (
	// Error condition that signals an invalid RFC3339 timestamp  of some
	// kind.
	//
	// This error is usually wrapped around a descriptive message string.
	ErrInvalidRFC3339 error = errors.Base("invalid RFC3339 timestamp")

	// Error condition that signals that an RFC3339 timestamp is not a
	// string format.
	//
	// This error is used by `Set` as well as JSON and YAML methods.
	ErrRFC3339NotString error = errors.Base("RFC3339 timestamp  must be a string")
)
```

#### func  CurrentZone

```go
func CurrentZone() (string, int)
```
Return the current timezone for the host.

#### func  PrettyFormat

```go
func PrettyFormat(dur Duration) string
```
Format a time duration in pretty format.

Example, a duration of 72 minutes becomes "1 hour(s), 12 minute(s)".

#### func  TimeToMySQL

```go
func TimeToMySQL(val time.Time) string
```
Convert a `time.Time` value to a MySQL timestamp for queries.

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

#### type Datum

```go
type Datum = any
```


#### type Duration

```go
type Duration time.Duration
```

Enhanced time duration type.

#### func  NewFromDuration

```go
func NewFromDuration(duration string) (Duration, error)
```

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

Note: this acquires the `readAvailable` semaphore.

#### func (*Mailbox) Empty

```go
func (m *Mailbox) Empty() bool
```
Is the mailbox empty like my heart?

#### func (*Mailbox) Full

```go
func (m *Mailbox) Full() bool
```
Does the mailbox contain a value?

#### func (*Mailbox) Get

```go
func (m *Mailbox) Get() (Datum, bool)
```
Get an element from the mailbox. Defaults to using a context with a deadline of
5 seconds.

#### func (*Mailbox) GetWithContext

```go
func (m *Mailbox) GetWithContext(ctx context.Context) (Datum, bool)
```
Get an element from the mailbox using the provided context.

It is recommended to use a context that has a timeout deadline.

#### func (*Mailbox) Put

```go
func (m *Mailbox) Put(elem Datum) bool
```
Put an element into the mailbox.

#### func (*Mailbox) PutWithContext

```go
func (m *Mailbox) PutWithContext(ctx context.Context, elem Datum) bool
```
Put an element into the mailbox using a context.

#### func (*Mailbox) Reset

```go
func (m *Mailbox) Reset()
```
Reset the mailbox.

#### func (*Mailbox) TryGet

```go
func (m *Mailbox) TryGet() (Datum, bool)
```
Try to get an element from the mailbox.

#### func (*Mailbox) TryPut

```go
func (m *Mailbox) TryPut(item Datum) bool
```
Try to put an element into the mailbox.

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

This is a cheap implementation of a FIFO queue.

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

#### func (*Queue) Empty

```go
func (q *Queue) Empty() bool
```
Is the queue empty?

#### func (*Queue) Full

```go
func (q *Queue) Full() bool
```
Is the queue full?

#### func (*Queue) Get

```go
func (q *Queue) Get() Datum
```
Remove an element from the start of the queue and return it.

This blocks.

#### func (*Queue) GetWithContext

```go
func (q *Queue) GetWithContext(ctx context.Context) (Datum, error)
```
Remove an element from the start of the queue and return it.

Will exit should the context time out or be cancelled.

This blocks.

#### func (*Queue) GetWithoutBlock

```go
func (q *Queue) GetWithoutBlock() (Datum, bool)
```
Remove an element from the start of the queue and return it.

This does not block.

#### func (*Queue) Len

```go
func (q *Queue) Len() int
```
Return the number of elements in the queue.

#### func (*Queue) Put

```go
func (q *Queue) Put(elem Datum)
```
Put an element on to the queue.

This blocks.

#### func (*Queue) PutWithContext

```go
func (q *Queue) PutWithContext(ctx context.Context, elem Datum) error
```
Put an element on to the queue.

Will exit should the context time out or be cancelled.

This blocks.

#### func (*Queue) PutWithoutBlock

```go
func (q *Queue) PutWithoutBlock(elem Datum) bool
```
Append an element to the queue.

Returns `false` if there is no more room in the queue.

This does not block.

#### type RFC3339

```go
type RFC3339 time.Time
```

RFC 3339 time type.

#### func  ParseRFC3339

```go
func ParseRFC3339(data string) (RFC3339, error)
```
Parse the given string for an RFC3339 timestamp.

If the timestamp is not a valid RFC3339 timestamp, then `ErrInvalidRFC3339` is
returned.

#### func  RFC3339FromUnix

```go
func RFC3339FromUnix(unix int64) RFC3339
```
Convert a Unix timestamp to an RFC3339 timestamp.

#### func (RFC3339) Add

```go
func (obj RFC3339) Add(d time.Duration) RFC3339
```
Add a `time.Duration` value to the timestamp, returning a new timestamp.

#### func (RFC3339) After

```go
func (obj RFC3339) After(t time.Time) bool
```
Is the given time after the time in the timestamp?

#### func (RFC3339) Before

```go
func (obj RFC3339) Before(t time.Time) bool
```
Is the given time before the time in the timestamp?

#### func (RFC3339) Equal

```go
func (obj RFC3339) Equal(t time.Time) bool
```
Is the given time equal to the time in the timestamp?

#### func (RFC3339) Format

```go
func (obj RFC3339) Format(format string) string
```
Format the timestamp with the given format.

#### func (RFC3339) IsDST

```go
func (obj RFC3339) IsDST() bool
```
Does the timestamp correspond to a time where DST is in effect?

#### func (RFC3339) IsZero

```go
func (obj RFC3339) IsZero() bool
```
Is the timestamp a zero value?

#### func (RFC3339) MarshalJSON

```go
func (obj RFC3339) MarshalJSON() ([]byte, error)
```
JSON marshalling method.

#### func (RFC3339) MarshalYAML

```go
func (obj RFC3339) MarshalYAML() (any, error)
```
YAML marshalling method.

#### func (RFC3339) MySQL

```go
func (obj RFC3339) MySQL() string
```
Return a string that can be used in MySQL queries.

#### func (*RFC3339) Set

```go
func (obj *RFC3339) Set(str string) error
```
Set the RFC3339 timestamp to that of the given string.

If the string value fails to parse, then `ErrInvalidRFC3339` is returned.

#### func (*RFC3339) SetYAML

```go
func (obj *RFC3339) SetYAML(value any) error
```
Set the RFC3339 value from a YAML value.

If the passed YAML value is not a string, then `ErrRFC3339NotString` is
returned.

Will also return any error condition from the `Set` method.

#### func (RFC3339) String

```go
func (obj RFC3339) String() string
```
Coerce an RFC3339 time value to a string.

#### func (RFC3339) Sub

```go
func (obj RFC3339) Sub(t time.Time) time.Duration
```
Subtract a `time.Time` value from the timestamp, returning a `time.Duration`.

#### func (RFC3339) Time

```go
func (obj RFC3339) Time() time.Time
```
Coerce an RFC3339 time value to a `time.Time` value.

#### func (RFC3339) Type

```go
func (obj RFC3339) Type() string
```
Return the data type name for CLI flag parsing purposes.

#### func (RFC3339) UTC

```go
func (obj RFC3339) UTC() time.Time
```
Return the UTC time for the timestamp.

RFC3339 timestamps are always UTC internally, so `UTC` is provided as a
courtesy.

#### func (RFC3339) Unix

```go
func (obj RFC3339) Unix() int64
```
Return the Unix time for the timestamp.

#### func (*RFC3339) UnmarshalJSON

```go
func (obj *RFC3339) UnmarshalJSON(data []byte) error
```
JSON unmarshalling method.

#### func (*RFC3339) UnmarshalYAML

```go
func (obj *RFC3339) UnmarshalYAML(value *yaml.Node) error
```
YAML unmarshalling method.
