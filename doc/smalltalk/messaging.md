<!-- -*- mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# Message Passing, Responder Chains, Selectors, and Protocols

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->
**Table of Contents**

- [Message Passing, Responder Chains, Selectors, and Protocols](#message-passing-responder-chains-selectors-and-protocols)
  - [Synopsis](#synopsis)
  - [1. Architecture at a glance](#1-architecture-at-a-glance)
  - [2. The events package](#2-the-events-package)
    - [Included event types](#included-event-types)
    - [Why this exists](#why-this-exists)
    - [Example: creating and responding to a message](#example-creating-and-responding-to-a-message)
    - [Example: using the event queue](#example-using-the-event-queue)
    - [Where events is enough by itself](#where-events-is-enough-by-itself)
  - [3. The responder package](#3-the-responder-package)
    - [What it gives you](#what-it-gives-you)
    - [The responder chain](#the-responder-chain)
    - [Example: a plain respondable object](#example-a-plain-respondable-object)
    - [When to use responder without selector](#when-to-use-responder-without-selector)
  - [4. The selector package](#4-the-selector-package)
    - [Core idea](#core-idea)
    - [Main pieces](#main-pieces)
      - [`selector.Table`](#selectortable)
      - [Forwarding](#forwarding)
      - [Selector responses and errors](#selector-responses-and-errors)
      - [`selector.Respondable`](#selectorrespondable)
    - [Introspection](#introspection)
      - [Example: defining a selector event](#example-defining-a-selector-event)
      - [Example: registering selectors on a selector-aware object](#example-registering-selectors-on-a-selector-aware-object)
      - [Example: metadata and introspection](#example-metadata-and-introspection)
      - [Example: aliasing and package defaults](#example-aliasing-and-package-defaults)
      - [Example: package registry](#example-package-registry)
    - [When selector is the right layer](#when-selector-is-the-right-layer)
  - [5. The protocols package](#5-the-protocols-package)
    - [The registry](#the-registry)
    - [Validate versus Verify](#validate-versus-verify)
      - [Validate](#validate)
      - [Verify](#verify)
    - [Example: defining and validating a protocol](#example-defining-and-validating-a-protocol)
    - [Example: protocol with custom verification](#example-protocol-with-custom-verification)
    - [When to use protocols](#when-to-use-protocols)
  - [6. Typical usage patterns](#6-typical-usage-patterns)
    - [Pattern A: plain event bus-ish objects](#pattern-a-plain-event-bus-ish-objects)
    - [Pattern B: command/operation objects without channel gymnastics](#pattern-b-commandoperation-objects-without-channel-gymnastics)
    - [Pattern C: responder chain + selectors](#pattern-c-responder-chain--selectors)
    - [Pattern D: pluggable object system](#pattern-d-pluggable-object-system)
  - [7. Design notes and trade-offs](#7-design-notes-and-trade-offs)
    - [Why return events rather than (any, error) everywhere?](#why-return-events-rather-than-any-error-everywhere)
    - [Why protocols are separate](#why-protocols-are-separate)
  - [8. A small end-to-end example](#8-a-small-end-to-end-example)
  - [9. Choosing the right package](#9-choosing-the-right-package)

<!-- markdown-toc end -->

## Synopsis

This document explains how the `events`, `responder`, `selector`, and `protocols` packages fit together.

Taken together, they form a small message-passing runtime for Go with a very particular flavour:

- `events` provides the common event vocabulary.
- `responder` provides Objective-C/Cocoa-style responder chains.
- `selector` adds selector-based dispatch, auxiliary methods, forwarding, and introspection.
- `protocols` provides named contracts over selectors.

You can use the layers independently:

- `events` on its own for generic event objects and queues.
- `events` + `responder` for message passing without selector dispatch.
- `events` + `selector` for selector-based objects.
- all four together for a Smalltalk / Flavors / CLOS / Objective-C inspired runtime.

---

## 1. Architecture at a glance

```text
                +-------------------+
                |     protocols     |
                | named contracts   |
                +---------+---------+
                          |
                +---------v---------+
                |      selector     |
                | selectors, before |
                | / primary / after |
                | forwarding, docs  |
                +---------+---------+
                          |
                +---------v---------+
                |     responder     |
                | names, types,     |
                | Invoke, chain     |
                +---------+---------+
                          |
                +---------v---------+
                |       events      |
                | time, message,    |
                | response, error,  |
                | queue, signal     |
                +-------------------+
```

The rough rule of thumb is:

- start with events
- add responder when you want named receivers and chain-based delivery
- add selector when you want symbolic message names and method tables
- add protocols when you want to declare and validate capabilities

## 2. The events package

The root abstraction is `events.Event`:

```go
package events

type Event interface {
    When() time.Time
    String() string
}
```

That is intentionally tiny. Any object with a timestamp and a string form can participate in the system.

### Included event types

The package currently includes these core event types:

- Time: basic timestamp carrier
- Message: command + payload + monotonically increasing index
- Response: response to a Message
- Error: event that wraps an error
- Interrupt: generic interrupt-like event with a payload
- Signal: wraps an os.Signal
- Forward: requests forwarding to another selector target
- Queue: thread-safe FIFO for events

### Why this exists

`events` gives the upper layers a common language without forcing a complicated inheritance hierarchy. It is deliberately boring, which is exactly what you want from a base event package.

### Example: creating and responding to a message

``` go
package main

import (
    "fmt"

    "github.com/Asmodai/gohacks/events"
)

func main() {
    msg := events.NewMessage("cache.flush", map[string]any{
        "scope": "all",
    })

    fmt.Println(msg.Command())
    fmt.Println(msg.Data())

    rsp := events.NewResponse(msg, "ok")
    fmt.Println(rsp.Command())
    fmt.Println(rsp.Response())
}
```

### Example: using the event queue

``` go
package main

import (
    "fmt"

    "github.com/Asmodai/gohacks/events"
)

func main() {
    q := events.NewQueue()

    q.Push(events.NewMessage("reload", nil))
    q.Push(events.NewInterrupt("shutdown requested"))

    for q.Events() > 0 {
        evt := q.Pop()
        fmt.Printf("%T -> %s\n", evt, evt.String())
    }
}

```

### Where events is enough by itself

Use `events` alone when you need:

- a consistent event interface
- timestamped messages
- in-process queues
- simple request/response style event objects
- a foundation for higher-level dispatch later

## 3. The responder package

The `responder` package introduces the idea of a respondable object:

``` go
type Respondable interface {
    ResponderName() string
    ResponderType() string
    RespondsTo(events.Event) bool
    Invoke(events.Event) events.Event
}
```

This is the first big step upward. Now events are not just data; they can be sent to named objects.

### What it gives you

- named targets via `ResponderName()`
- groupable targets via `ResponderType()`
- capability checks via `RespondsTo()`
- dynamic invocation via `Invoke()`
- a thread-safe responder chain via `Chain`

### The responder chain

A `responder.Chain` is a registry of receivers. Responders can be added by name, replaced, removed, and queried by name or type.

Delivery styles include:

- `SendNamed`: send to one named responder
- `MustSendNamed`: same, but panic on failure
- `SendFirst`: send to the first responder in priority order that accepts it
- `MustSendFirst`: same, but panic on failure
- `SendAll`: broadcast to all responders that accept the event

The chain also tracks responders by type and priority, which makes it suitable for small object systems, plugin architectures, actor-ish subsystems, or UI-like propagation.

### Example: a plain respondable object

``` go
package main

import (
    "fmt"

    "github.com/Asmodai/gohacks/events"
    "github.com/Asmodai/gohacks/responder"
)

type Printer struct {
    name string
}

func (p *Printer) ResponderName() string { return p.name }
func (p *Printer) ResponderType() string { return "printer" }

func (p *Printer) RespondsTo(evt events.Event) bool {
    _, ok := evt.(*events.Message)
    return ok
}

func (p *Printer) Invoke(evt events.Event) events.Event {
    msg := evt.(*events.Message)
    text, _ := msg.Data().(string)
    return events.NewResponse(msg, fmt.Sprintf("%s handled %q", p.name, text))
}

func main() {
    chain := responder.NewChain("example")
    _, _ = chain.Add(&Printer{name: "console"})

    result, responds, found := chain.SendNamed(
        "console",
        events.NewMessage("print", "hello"),
    )

    fmt.Println(found, responds)
    fmt.Printf("%T -> %s\n", result, result.String())
}
```

### When to use responder without selector

This layer is enough when you want:

- named receivers
- chain traversal
- dynamic routing by object name or type
- event-driven delivery without a selector table

If you are happy writing your own `RespondsTo` and `Invoke` logic, you can stop here.

## 4. The selector package

The `selector` package takes the responder idea and adds symbolic message names, selector tables, auxiliary methods, forwarding, and introspection.

This is where the runtime starts feeling delightfully Objective-C-ish.

### Core idea

A selector-specific event is any event that also provides a selector string:

``` go
type SelectorEvent interface {
    events.Event
    Selector() string
}
```

A selector-aware object can then dispatch on the selector name instead of doing manual switch logic in `Invoke()`.

### Main pieces

#### `selector.Table`

A `Table` maps selector names to methods. Each selector entry can have:

- one primary method
- zero or more before methods
- zero or more after methods

The primary method has the signature:

``` go
type Method func(responder.Respondable, events.Event) events.Event
```

That gives you a compact object/message model:

``` text
SelectorEvent
   -> :before methods
   -> primary method
   -> :after methods
   -> result event
```


#### Forwarding

If a method returns `*events.Forward`, the table can re-dispatch to another selector target. A maximum forward depth is enforced to avoid loops.

#### Selector responses and errors

The package defines selector-native result types:

- SelectorResponse
- SelectorError

Both are still events, so the whole system stays message-oriented.


#### `selector.Respondable`

This is a ready-made implementation of a selector-aware receiver.

It provides:

- `ResponderName()`
- `ResponderType()`
- a selector Table
- protocol membership tracking
- selector introspection helpers
- an `Invoke()` that dispatches via selector name

### Introspection

Objects implementing `selector.Introspectable` can describe themselves:

- `Selectors()`
- `SortedSelectors()`
- `Methods()`
- `ConformsTo()`
- `ListProtocols()`
- `MetadataForSelector()`

`DumpIntrospectableInfo()` turns that into a human-readable description.

#### Example: defining a selector event

``` go
package main

import (
    "time"

    "github.com/Asmodai/gohacks/events"
)

type Ping struct {
    events.Time
    from string
}

func NewPing(from string) *Ping {
    return &Ping{
        Time: events.Time{TStamp: time.Now()},
        from: from,
    }
}

func (p *Ping) Selector() string { return "ping" }
func (p *Ping) String() string   { return "ping from " + p.from }
```


#### Example: registering selectors on a selector-aware object

``` go
package main

import (
    "fmt"

    "github.com/Asmodai/gohacks/events"
    "github.com/Asmodai/gohacks/selector"
)

func main() {
    obj := selector.NewRespondable("worker-1", "worker")

    obj.Methods().Register("ping", func(_ any, evt events.Event) events.Event {
        selEvt := evt.(selector.SelectorEvent)
        return selector.NewSelectorResponse(selEvt, "pong")
    })

    obj.Methods().AddBefore("ping", func(_ any, evt events.Event) events.Event {
        fmt.Println("before ping")
        return evt
    })

    obj.Methods().AddAfter("ping", func(_ any, evt events.Event) events.Event {
        fmt.Printf("after ping -> %T\n", evt)
        return evt
    })

    rsp := obj.Invoke(NewPing("tester"))

    switch v := rsp.(type) {
    case *selector.SelectorError:
        panic(v.Error())
    case selector.SelectorEvent:
        fmt.Printf("selector=%s\n", v.Selector())
    default:
        fmt.Printf("result=%T\n", rsp)
    }
}
```

#### Example: metadata and introspection

If you attach metadata to selectors, introspection becomes self-documenting.
That makes the runtime suitable for command catalogs, debug dumps, and tooling.

Conceptually, the flow looks like this:

``` go
info := selector.DumpIntrospectableInfo(obj)
fmt.Println(info)
```

This prints the object's name, type, protocols, selectors, and metadata such as version, since, protocol, visibility, tags, author, and doc strings.

#### Example: aliasing and package defaults

`selector.Package` groups a table with exports, aliases, and defaults.

Use it when you want selector namespaces that behave more like a language package:

- export only chosen selectors
- define aliases
- define default selectors for operations

``` go
pkg := selector.NewPackage("core")

pkg.Table.Register("process.start", startMethod)
pkg.Table.Register("process.stop", stopMethod)

pkg.Export("process.start")
pkg.Export("process.stop")

pkg.Alias("start", "process.start")
pkg.SetDefault("lifecycle", "process.start")
```

This is useful when you want stable external names while keeping internal selector names more verbose.

#### Example: package registry

Multiple selector packages can be collected in a `selector.Registry`:

``` go
reg := selector.NewRegistry()
reg.AddPackage(pkg)

core, ok := reg.GetPackage("core")
```

That gives you a natural place to manage subsystems, command namespaces, or pluggable modules.

### When selector is the right layer

Use selector when you want:

- Objective-C style symbolic dispatch
- Smalltalk-ish message passing
- Flavors/CLOS-style before/after behaviour
- explicit forwarding
- method tables instead of giant switch statements
- introspection and metadata

## 5. The protocols package

The `protocols` package adds named contracts over selectors.

A protocol is simply:

``` go
type Protocol struct {
    Name      string
    Selectors []string
}
```

That minimal shape is a feature, not a limitation. It lets you define what an object is expected to implement without entangling the object model itself.

### The registry

A `protocols.Registry` stores:

- registered protocols
- optional verifier functions per protocol

Key operations are:

- Register(proto)
- RegisterWithVerifier(proto, verifier)
- Verify(name, obj)
- Validate(name, rbl)

### Validate versus Verify

They do slightly different jobs.

#### Validate

Checks that an object has all selectors listed by the protocol.

This is structural validation: “do the required methods exist?”

#### Verify

Runs a custom verifier against an introspectable object.

This is semantic validation: “does the object comply with the stronger rules of
this protocol?”

### Example: defining and validating a protocol

``` go
package main

import (
    "fmt"

    "github.com/Asmodai/gohacks/protocols"
)

func main() {
    reg := protocols.NewRegistry()

    proto := &protocols.Protocol{
        Name: "lifecycle",
        Selectors: []string{
            "start",
            "stop",
            "status",
        },
    }

    reg.Register(proto)

    ok := reg.Validate("lifecycle", obj)
    fmt.Println("valid:", ok)
}
```

Where `obj` is something exposing `Methods() *selector.Table`, such as a `selector.Respondable`.


### Example: protocol with custom verification

``` go
reg.RegisterWithVerifier(proto, func(obj selector.Introspectable) error {
    if !obj.ConformsTo("lifecycle") {
        return errors.New("object does not claim lifecycle conformance")
    }

    // Custom policy could inspect metadata, visibility, versioning, etc.
    return nil
})
```

This is useful when “has the selectors” is necessary but not sufficient.

### When to use protocols

Use protocols when you want:

- named capability sets
- conformance checks
- plugin validation
- runtime contracts for selector-aware objects
- richer tooling and documentation

## 6. Typical usage patterns

### Pattern A: plain event bus-ish objects

Use:

- events
- optionally responder

Good when you just want a set of objects that can receive timestamped messages.

### Pattern B: command/operation objects without channel gymnastics

Use:

- events
- selector

This is the sweet spot when you want code that reads more like:

``` go
result := proc.Invoke(evt)
```

instead of building ad-hoc command structs and giant switch statements.

### Pattern C: responder chain + selectors

Use:

- events
- responder
- selector

Good when you want named objects in a chain, but each object also has a proper selector table.

### Pattern D: pluggable object system

Use all four packages.

This is the full-fat model:

- event vocabulary
- named responders
- selector dispatch
- protocol contracts
- introspection and metadata

## 7. Design notes and trade-offs

Why not just use channels?

Go channels are excellent for concurrency and streaming, but they can become ugly when used as a faux RPC layer:

``` go
cmdCh <- command{op: "reload", arg: x, resp: replyCh}
reply := <-replyCh
```

The selector approach gives you a more object/message oriented style:

``` go
reply := obj.Invoke(evt)
```

or, when used with chains:

``` go
reply, ok := chain.SendFirst(evt)
```

### Why return events rather than (any, error) everywhere?

Because the system stays inside one conceptual model:

- input is an event
- output is an event
- errors are also events

That keeps dispatch logic uniform and makes tracing and routing simpler.

### Why protocols are separate

Protocols are kept out of the selector table itself so that capability contracts remain declarative and optional.

That keeps the selector layer lightweight while still allowing richer runtime validation when needed.

## 8. A small end-to-end example

The following sketch shows how the pieces fit together conceptually.

``` go
package main

import (
    "fmt"
    "time"

    "github.com/Asmodai/gohacks/events"
    "github.com/Asmodai/gohacks/protocols"
    "github.com/Asmodai/gohacks/responder"
    "github.com/Asmodai/gohacks/selector"
)

type StartEvent struct {
    events.Time
}

func NewStartEvent() *StartEvent {
    return &StartEvent{Time: events.Time{TStamp: time.Now()}}
}

func (e *StartEvent) Selector() string { return "start" }
func (e *StartEvent) String() string   { return "start" }

func main() {
    proc := selector.NewRespondable("proc-1", "process")
    proc.AddProtocol("lifecycle")

    proc.Methods().Register("start", func(_ responder.Respondable, evt events.Event) events.Event {
        selEvt := evt.(selector.SelectorEvent)
        return selector.NewSelectorResponse(selEvt, "started")
    })

    chain := responder.NewChain("processes")
    _, _ = chain.Add(proc)

    preg := protocols.NewRegistry()
    preg.Register(&protocols.Protocol{
        Name:      "lifecycle",
        Selectors: []string{"start"},
    })

    fmt.Println("protocol valid:", preg.Validate("lifecycle", proc))

    evt := NewStartEvent()
    rsp, ok := chain.SendFirst(evt)
    if !ok {
        panic("no handler")
    }

    switch v := rsp.(type) {
    case *selector.SelectorError:
        panic(v.Error())
    case selector.SelectorEvent:
        fmt.Println("selector:", v.Selector())
        if sr, ok := v.(*selector.SelectorResponse); ok {
            fmt.Println("response:", sr.Response())
        }
    default:
        fmt.Printf("unexpected result: %T\n", rsp)
    }
}
```

The important part is not the exact toy event type. The important part is the flow:

- define a selector event
- register selector methods on a selector-aware object
- optionally put the object in a responder chain
- optionally validate it against a protocol
- invoke it and handle a selector response or selector error

## 9. Choosing the right package

Use `events` when

- you need timestamped event values
- you want a queue
- you want a tiny common interface

Use `responder` when

- you want named targets
- you want chain-based delivery
- you want to route by object identity or type

Use `selector` when

- you want message names instead of a big switch
- you want before/after hooks
- you want forwarding
- you want self-describing objects

Use `protocols` when

- you want formal capability sets
- you want runtime validation
- you want plugins or modules to declare what they implement

10. Summary

These packages form a layered runtime:

- events gives you the atoms
- responder gives you named receivers and chains
- selector gives you symbolic dispatch and method combination
- protocols gives you contracts and validation

That combination is a neat fit when you want message passing in Go without turning everything into command channels, ad-hoc RPC structs, or giant switch-based dispatchers.

11. Do you miss `SELECTQ` and `DEFSELECT`?

Yes.

Yes I do.
