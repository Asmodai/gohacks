<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# responder -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/responder"
```

## Usage

```go
const ContextKeyResponderChain = "gohacks/responder@v1"
```
Key used to store the instance in the context's user value.

```go
var (
	// Error condition signalled when an attempt is made to add a
	// non-unique responder to a responder chain.
	ErrDuplicateResponder error = errors.Base("duplicate responder")

	// Error condition signalled when a responder's name via `Name()` is
	// invalid or zero length.
	ErrResponderNameInvalid error = errors.Base("responder name invalid")
)
```

```go
var ErrValueNotResponderChain = errx.Base("value is not *Chain")
```
Signalled if the instance associated with the context key is not of type *Chain.

#### func  SetResponderChain

```go
func SetResponderChain(ctx context.Context, inst *Chain) (context.Context, error)
```
Set ResponderChain stores the instance in the context map.

#### func  SetResponderChainIfAbsent

```go
func SetResponderChainIfAbsent(ctx context.Context, inst *Chain) (context.Context, error)
```
SetResponderChainIfAbsent sets only if not already present.

#### func  WithResponderChain

```go
func WithResponderChain(ctx context.Context, fn func(*Chain))
```
WithResponderChain calls fn with the instance or fallback.

#### type Chain

```go
type Chain struct {
}
```

Responder chain structure.

This attempts to bring a little bit of Smalltalk and Objective-C to the
wonderful world of Go.

It might also attempt to bring a bit of MIT Flavors, too... but don't expect to
see crazy like `defwrapper` and `defwhopper`.

#### func  FromResponderChain

```go
func FromResponderChain(ctx context.Context) *Chain
```
FromResponderChain returns the instance or the fallback.

#### func  GetResponderChain

```go
func GetResponderChain(ctx context.Context) (*Chain, error)
```
Get the instance from the given context.

Will return ErrValueNotResponderChain if the value in the context is not of type
*Chain.

#### func  MustGetResponderChain

```go
func MustGetResponderChain(ctx context.Context) *Chain
```
Attempt to get the instance from the given context. Panics if the operation
fails.

#### func  NewChain

```go
func NewChain(name string) *Chain
```
Create a new responder chain object.

#### func  TryGetResponderChain

```go
func TryGetResponderChain(ctx context.Context) (*Chain, bool)
```
TryGetResponderChain returns the instance and true if present and typed.

#### func (*Chain) Add

```go
func (chain *Chain) Add(responder Respondable) (Respondable, error)
```
Adds a responder to the responder chain.

The responder will have a default priority of 0.

Returns `ErrDuplicateResponder` if a non-unique responder is added.

#### func (*Chain) AddNamed

```go
func (chain *Chain) AddNamed(name string, responder Respondable) (Respondable, error)
```
Adds the supplied responder to the chain using a user-specified name.

The given name overrides that provided by the `Name()` method in the responder,
thus should only be used in use-cases where you need a specific identifier for a
responder.

The responder will have a default priority of 0.

Returns `ErrDuplicateResponder` if a non-unique responder is added.

#### func (*Chain) AddNamedWithPriority

```go
func (chain *Chain) AddNamedWithPriority(
	name string,
	responder Respondable,
	priority int,
) (Respondable, error)
```
Adds the supplied responder to the chain using a user-specified name and
priority.

The given name overrides that provided by the `Name()` method in the responder,
thus should only be used in use-cases where you need a specific identifier for a
responder.

Returns `ErrDuplicateResponder` if a non-unique responder is added.

#### func (*Chain) AddOrReplace

```go
func (chain *Chain) AddOrReplace(responder Respondable) Respondable
```
Adds a responder to the responder chain.

The responder will have a default priority of 0.

If the responder already exists in the chain then it is replaced with the
provided responder.

#### func (*Chain) AddOrReplaceNamed

```go
func (chain *Chain) AddOrReplaceNamed(name string, responder Respondable) Respondable
```
Adds the supplied responder to the chain using a user-specified name.

The responder will have a default priority of 0.

If the responder already exists in the chain then it is replaced with the
provided responder.

The given name overrides that provided by the `Name()` method in the responder,
thus should only be used in use-cases where you need a specific identifier for a
responder.

#### func (*Chain) AddOrReplaceNamedWithPriority

```go
func (chain *Chain) AddOrReplaceNamedWithPriority(
	name string,
	responder Respondable,
	priority int,
) Respondable
```
Adds the supplied responder to the chain using a user-specified name and
priority.

If the responder already exists in the chain then it is replaced with the
provided responder.

The given name overrides that provided by the `Name()` method in the responder,
thus should only be used in use-cases where you need a specific identifier for a
responder.

#### func (*Chain) AddOrReplaceWithPriority

```go
func (chain *Chain) AddOrReplaceWithPriority(responder Respondable, priority int) Respondable
```
Adds a responder to the responder chain.

If the responder already exists in the chain then it is replaced with the
provided responder.

#### func (*Chain) AddWithPriority

```go
func (chain *Chain) AddWithPriority(responder Respondable, priority int) (Respondable, error)
```
Adds the supplied responder to the chain using the given priority.

Returns `ErrDuplicateResponder` if a non-unique responder is added.

#### func (*Chain) Clear

```go
func (chain *Chain) Clear()
```
Clear all responders from the chain.

#### func (*Chain) Count

```go
func (chain *Chain) Count() int
```
Return the number of responders in the chain.

#### func (*Chain) Invoke

```go
func (chain *Chain) Invoke(event events.Event) events.Event
```
Send an event to the chain.

The first object that can respond to the event will consume it.

Implements `Respondable`.

#### func (*Chain) IsEmpty

```go
func (chain *Chain) IsEmpty() bool
```
Is the chain empty?

#### func (*Chain) MustSendFirst

```go
func (chain *Chain) MustSendFirst(event events.Event) events.Event
```
Send a message to the responder chain. Panics if no responder is able to respond
to the event.

#### func (*Chain) MustSendNamed

```go
func (chain *Chain) MustSendNamed(name string, event events.Event) events.Event
```
Send a message to a specific responder. Panics if the responder does not exist
or does not respond to the event.

#### func (*Chain) Name

```go
func (chain *Chain) Name() string
```
Return the name of the chain.

Implements `Respondable`.

#### func (*Chain) Names

```go
func (chain *Chain) Names() []string
```
Return a list of names for the responders currently in the chain.

#### func (*Chain) Remove

```go
func (chain *Chain) Remove(responder Respondable) bool
```

#### func (*Chain) RemoveNamed

```go
func (chain *Chain) RemoveNamed(name string) bool
```
Remove the named responder from the responder chain.

Returns false if no such responder was found.

#### func (*Chain) RespondsTo

```go
func (chain *Chain) RespondsTo(event events.Event) bool
```
Iterate over responders checking if any implement the given event.

The first responder found that responds to the event will result in `true` being
returned.

Implements `Respondable`.

#### func (*Chain) SendAll

```go
func (chain *Chain) SendAll(event events.Event) []events.Event
```
Send a message to all responders in the chain.

All responders capable of responding to the event will receive the event.

Returns a list of objects of interface `events.Event`, which may be the same
event as was passed. Doing sanity on the return values is up to you

This method is thread-safe.

#### func (*Chain) SendFirst

```go
func (chain *Chain) SendFirst(event events.Event) (events.Event, bool)
```
Send a message to the responder chain.

The first responder capable of responding to the event will consume the event.

Returns an object of interface `events.Event`, which may be the same event that
was passed to it. Doing sanity on the return value is up to you.

Unlike `SendNamed`, there is no indicator that either receivers were not found
or that there were no receivers. So it would be wise to assume that a `false`
means that absolutely nothing in the chain received your event.

Example usage would be:

```go

    result, ok := someChain.Send(evt)
    if !ok {
    	log.Warn("No responders have responded to the event.")
    } else {
    	log.Info("At least one responder has responded to the event.")
    }

```

This method is thread-safe.

#### func (*Chain) SendNamed

```go
func (chain *Chain) SendNamed(name string, event events.Event) (events.Event, bool, bool)
```
Send a message to a specific responder.

Returns an object of interface `events.Event`, which may be the same event that
was passed to it. Doing sanity on the return value is up to you. Returns the
resulting event from the responder, a boolean value that states if the responder
was able to respond, and a boolean value that states whether the responder was
found.

Example usage would be:

```go

    result, responds, found := someChain.SendNamed("something", evt)
    if !found {
    	log.Warn("Responder `something` not found!")
    } else if !responds {
    	log.Info("Responder 'something' ignored the event.")
    } else {
    	log.Debug("Event was handled.")
    }

```

This method is thread-safe.

#### func (*Chain) SendType

```go
func (chain *Chain) SendType(typeName string, event events.Event) []events.Event
```
Send a message to all responders of a given type in the chain.

All responders of the given type that are capable of responding to the event
will receive the event.

Returns a list of objects of interface `events.Event`, which may be the same
event as was passed. Doing sanity on the return values is up to you.

This method is thread-safe.

#### func (*Chain) Type

```go
func (chain *Chain) Type() string
```
Return the type name of the chain.

Implements `Respondable`.

#### type Respondable

```go
type Respondable interface {
	// A unique name for the respondable object.
	//
	// As this allows us to send events to a specific thing the value
	// returned here must be unique.
	Name() string

	// The type of the respondable object.
	//
	// This can be the internal Go type, or some arbitrary user-specified
	// value that makes sense to you.
	//
	// This is used to implement a "send to all of type" system.
	Type() string

	// Does the receiver respond to a specific event or event type?
	//
	// There is no definition for what `RespondsTo` should do other than
	// return a boolean that states whether an object responds to an
	// event or not.
	RespondsTo(events.Event) bool

	// Send an event to the object.
	//
	// There is no second return value to indicate success or whether
	// the event was handled or not.  The idea being that the receiver
	// will send an `events.Response` event back.
	Invoke(events.Event) events.Event
}
```

Objects the implement these methods are considered `respondable` and are deemed
capable of being sent messages directly or via responder chain.
