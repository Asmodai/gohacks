-*- Mode: gfm -*-

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
}
```


#### func  NewError

```go
func NewError(err error) *Error
```

#### func (*Error) Error

```go
func (e *Error) Error() string
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
func NewMessage(cmd int, data any) *Message
```

#### func (*Message) Command

```go
func (e *Message) Command() int
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

#### type MockEvent

```go
type MockEvent struct {
}
```

MockEvent is a mock of Event interface.

#### func  NewMockEvent

```go
func NewMockEvent(ctrl *gomock.Controller) *MockEvent
```
NewMockEvent creates a new mock instance.

#### func (*MockEvent) EXPECT

```go
func (m *MockEvent) EXPECT() *MockEventMockRecorder
```
EXPECT returns an object that allows the caller to indicate expected use.

#### func (*MockEvent) String

```go
func (m *MockEvent) String() string
```
String mocks base method.

#### func (*MockEvent) When

```go
func (m *MockEvent) When() time.Time
```
When mocks base method.

#### type MockEventMockRecorder

```go
type MockEventMockRecorder struct {
}
```

MockEventMockRecorder is the mock recorder for MockEvent.

#### func (*MockEventMockRecorder) String

```go
func (mr *MockEventMockRecorder) String() *gomock.Call
```
String indicates an expected call of String.

#### func (*MockEventMockRecorder) When

```go
func (mr *MockEventMockRecorder) When() *gomock.Call
```
When indicates an expected call of When.

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
func (e *Response) Command() int
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
