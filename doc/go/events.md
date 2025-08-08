<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# events -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/events"
```

## Usage

#### func  EventType

```go
func EventType(e Event) reflect.Type
```

#### type Error

```go
type Error struct {
	Time

	Err error
}
```


#### func  NewError

```go
func NewError(err error) *Error
```

#### func (*Error) Error

```go
func (e *Error) Error() error
```

#### func (*Error) String

```go
func (e *Error) String() string
```

#### type Event

```go
type Event interface {
	When() time.Time
	String() string
}
```


#### type EventList

```go
type EventList []Event
```


#### type Forward

```go
type Forward struct {
	Time
}
```


#### func  NewForward

```go
func NewForward(to string, event Event) *Forward
```

#### func (*Forward) Event

```go
func (f *Forward) Event() Event
```

#### func (*Forward) String

```go
func (f *Forward) String() string
```

#### func (*Forward) To

```go
func (f *Forward) To() string
```

#### type Interrupt

```go
type Interrupt struct {
	Time
}
```


#### func  NewInterrupt

```go
func NewInterrupt(data any) *Interrupt
```

#### func (*Interrupt) Data

```go
func (i *Interrupt) Data() any
```

#### func (*Interrupt) String

```go
func (i *Interrupt) String() string
```

#### type Message

```go
type Message struct {
	Time
}
```


#### func  NewMessage

```go
func NewMessage(cmd string, data any) *Message
```

#### func (*Message) Command

```go
func (e *Message) Command() string
```

#### func (*Message) Data

```go
func (e *Message) Data() any
```

#### func (*Message) Index

```go
func (e *Message) Index() uint64
```

#### func (*Message) String

```go
func (e *Message) String() string
```

#### type Queue

```go
type Queue struct {
	sync.Mutex
}
```


#### func  NewQueue

```go
func NewQueue() *Queue
```

#### func (*Queue) Capacity

```go
func (e *Queue) Capacity() int
```

#### func (*Queue) Events

```go
func (e *Queue) Events() int
```

#### func (*Queue) Pop

```go
func (e *Queue) Pop() Event
```

#### func (*Queue) Push

```go
func (e *Queue) Push(evt Event)
```

#### type Response

```go
type Response struct {
	Time
}
```


#### func  NewResponse

```go
func NewResponse(msg *Message, rsp any) *Response
```

#### func (*Response) Command

```go
func (e *Response) Command() string
```

#### func (*Response) Index

```go
func (e *Response) Index() uint64
```

#### func (*Response) Received

```go
func (e *Response) Received() time.Time
```

#### func (*Response) Response

```go
func (e *Response) Response() any
```

#### func (*Response) String

```go
func (e *Response) String() string
```

#### type Signal

```go
type Signal struct {
	Time
}
```


#### func  NewSignal

```go
func NewSignal(sig os.Signal) *Signal
```

#### func (*Signal) Signal

```go
func (e *Signal) Signal() os.Signal
```

#### func (*Signal) String

```go
func (e *Signal) String() string
```

#### type Time

```go
type Time struct {
	TStamp time.Time
}
```


#### func  NewTime

```go
func NewTime() *Time
```

#### func (*Time) SetNow

```go
func (e *Time) SetNow()
```

#### func (*Time) SetWhen

```go
func (e *Time) SetWhen(val time.Time)
```

#### func (*Time) String

```go
func (e *Time) String() string
```

#### func (*Time) When

```go
func (e *Time) When() time.Time
```
