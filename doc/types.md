-*- Mode: gfm -*-

# types -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/types"
```

## Usage

```go
const (
	// Normal (default) attribtues -- VT100.
	NORMAL = 0

	// Bold -- VT100.
	BOLD = 1

	// Faint, decreased intensity -- ECMA-48 2e.
	FAINT = 2

	// Italicizsed -- ECMA-48 2e.
	ITALICS = 3

	// Underlined -- VT100.
	UNDERLINE = 4

	// Blinking -- VT100.
	BLINK = 5

	// Inverse video -- VT100.
	INVERSE = 7

	// Crossed-out characters -- ECMA-48 3e.
	STRIKETHROUGH = 9

	// Black colour -- ANSI.
	BLACK = 0

	// Red colour -- ANSI.
	RED = 1

	// Green colour -- ANSI.
	GREEN = 2

	// Yellow colour -- ANSI.
	YELLOW = 3

	// Blue colour -- ANSI.
	BLUE = 4

	// Magenta colour -- ANSI.
	MAGENTA = 5

	// Cyan colour -- ANSI.
	CYAN = 6

	// White colour -- ANSI.
	WHITE = 7

	// Default colour -- ANSI.
	DEFAULT = 9

	// Offset for foreground colours.
	FGOFFSET = 30

	// Offset for background colours.
	BGOFFSET = 40
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

#### func  MakeColorString

```go
func MakeColorString() *ColorString
```
Make a new coloured string.

#### func  MakeColorStringWithAttrs

```go
func MakeColorStringWithAttrs(data string, attr, foreg, backg int) *ColorString
```
Make a new coloured string with the given attributes.

#### func (*ColorString) SetAttr

```go
func (cs *ColorString) SetAttr(attr int)
```
Set the character attribute.

#### func (*ColorString) SetBG

```go
func (cs *ColorString) SetBG(col int)
```
Set the background colour.

#### func (*ColorString) SetFG

```go
func (cs *ColorString) SetFG(col int)
```
Set the foreground colour.

#### func (*ColorString) SetString

```go
func (cs *ColorString) SetString(val string)
```
Set the string to display.

#### func (*ColorString) String

```go
func (cs *ColorString) String() string
```
Convert to a string.

Warning: the resulting string will contain escape sequences for use with a
compliant terminal or terminal emulator.

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
