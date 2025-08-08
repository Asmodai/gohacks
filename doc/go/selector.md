<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# selector -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/selector"
```

## Usage

```go
const (
	DefaultPriority = int(100)
)
```

```go
var (
	ErrForwardLoop       = errors.Base("selector forward loop detected")
	ErrHasNoSelector     = errors.Base("has no selector")
	ErrNoMethodExists    = errors.Base("no method by this name exists")
	ErrNoMethodSpecified = errors.Base("no method specified")
	ErrNoMethodToWrap    = errors.Base("no method to wrap")
	ErrReferenceParse    = errors.Base("reference parse failure")
	ErrSelectorNotFound  = errors.Base("no method for selector")
	ErrSelectorPanic     = errors.Base("panic during selector method")
	ErrUnresolved        = errors.Base("unresolved selector")
)
```

#### func  DisableTrace

```go
func DisableTrace()
```

#### func  DumpIntrospectableInfo

```go
func DumpIntrospectableInfo(obj Introspectable) string
```

#### func  EnableTrace

```go
func EnableTrace()
```

#### func  SetTraceOutput

```go
func SetTraceOutput(out *os.File)
```

#### func  Trace

```go
func Trace(format string, args ...any)
```

#### func  TraceWithWrapper

```go
func TraceWithWrapper(isWrapped bool, format string, args ...any)
```

#### type AuxiliaryMethod

```go
type AuxiliaryMethod struct {
}
```


#### type Entry

```go
type Entry struct {
}
```

Selector table entry.

#### type Introspectable

```go
type Introspectable interface {
	responder.Respondable

	// Return a list of selectors.
	Selectors() []string

	// Return a list of sorted selectors
	SortedSelectors() []string

	// Return a list of methods that the object can respond to.
	Methods() *Table

	// Does the object conform to the given protocol?
	ConformsTo(protocol string) bool

	// List all protocols for which the object claims conformity.
	ListProtocols() []string

	// Return a map of metadata for a selector
	MetadataForSelector(string) (map[string]string, error)
}
```

Objects that implement these methods are considered `introspectable` and are
deemed capable of being asked to describe themselves in various ways.

#### type Method

```go
type Method func(responder.Respondable, events.Event) events.Event
```

Selector method function signature type.

#### type Namespace

```go
type Namespace struct {
	Uses   []*Package
	Shadow map[string]bool
}
```


#### func (*Namespace) Dispatch

```go
func (ns *Namespace) Dispatch(
	reg *Registry,
	raw string, target responder.Respondable, evt events.Event,
) (events.Event, bool, string)
```

#### func (*Namespace) Resolve

```go
func (ns *Namespace) Resolve(reg *Registry, ref Ref) (ResolveResult, bool)
```
Resolve a reference.

#### type Package

```go
type Package struct {
	Name  string
	Table *Table
}
```

Packages.

A package provides a table of selectors, aliases, and defaults.

#### func  NewPackage

```go
func NewPackage(name string) *Package
```

#### func (*Package) Alias

```go
func (pkg *Package) Alias(alias, target string) bool
```
Create an alias that maps alias to target.

#### func (*Package) Export

```go
func (pkg *Package) Export(name string) bool
```
Export a selector.

#### func (*Package) GetDefault

```go
func (pkg *Package) GetDefault(op string) (string, bool)
```
Returns the default for the given operation.

#### func (*Package) IsExported

```go
func (pkg *Package) IsExported(name string) bool
```
Is the given selector exported?

#### func (*Package) ResolveAlias

```go
func (pkg *Package) ResolveAlias(name string) string
```
Resolve a selector.

If the specified selector is an alias, then its target is returned.

#### func (*Package) SetDefault

```go
func (pkg *Package) SetDefault(operator, name string) bool
```
Sets a default selector.

#### func (*Package) Unexport

```go
func (pkg *Package) Unexport(name string)
```

#### type Ref

```go
type Ref struct {
	Package  string // Package name.
	Name     string // Name.
	Version  string // Version, e.g. "v1", "v1.2" et al
	Internal bool   // If true, then thing is package internal.
}
```


#### func  ParseRef

```go
func ParseRef(ref string) (Ref, bool)
```
Parse a reference.

Supports:

    name
    name@version
    package:name
    package::name
    package:name@version
    package::name@version

#### type Registry

```go
type Registry struct {
	GlobalDefault *Package
}
```


#### func  NewRegistry

```go
func NewRegistry() *Registry
```

#### func (*Registry) AddPackage

```go
func (r *Registry) AddPackage(pkg *Package)
```

#### func (*Registry) GetPackage

```go
func (r *Registry) GetPackage(name string) (*Package, bool)
```

#### type ResolveResult

```go
type ResolveResult struct {
	Pkg   *Package
	Table *Table
	Name  string
	Why   string
}
```


#### func (ResolveResult) String

```go
func (r ResolveResult) String() string
```

#### type Respondable

```go
type Respondable struct {
}
```


#### func  NewRespondable

```go
func NewRespondable(name, typeName string) *Respondable
```

#### func (*Respondable) AddProtocol

```go
func (sr *Respondable) AddProtocol(name string)
```

#### func (*Respondable) ConformsTo

```go
func (sr *Respondable) ConformsTo(name string) bool
```

#### func (*Respondable) Invoke

```go
func (sr *Respondable) Invoke(evt events.Event) events.Event
```

#### func (*Respondable) ListProtocols

```go
func (sr *Respondable) ListProtocols() []string
```

#### func (*Respondable) MetadataForSelector

```go
func (sr *Respondable) MetadataForSelector(selector string) (map[string]string, error)
```

#### func (*Respondable) Methods

```go
func (sr *Respondable) Methods() *Table
```

#### func (*Respondable) Name

```go
func (sr *Respondable) Name() string
```

#### func (*Respondable) RespondsTo

```go
func (sr *Respondable) RespondsTo(evt events.Event) bool
```

#### func (*Respondable) Selectors

```go
func (sr *Respondable) Selectors() []string
```

#### func (*Respondable) SortedSelectors

```go
func (sr *Respondable) SortedSelectors() []string
```

#### func (*Respondable) Type

```go
func (sr *Respondable) Type() string
```

#### type SelectorError

```go
type SelectorError struct {
}
```

Selector error event.

NOTE: `golangci-lint` will want this to be called `Error', and that is not what
we want. This is an explicit event, not to be confused with `events.Error`.

#### func  NewSelectorError

```go
func NewSelectorError(err error) *SelectorError
```

#### func (*SelectorError) Error

```go
func (e *SelectorError) Error() error
```

#### func (*SelectorError) Selector

```go
func (e *SelectorError) Selector() string
```

#### func (*SelectorError) String

```go
func (e *SelectorError) String() string
```

#### func (*SelectorError) When

```go
func (e *SelectorError) When() time.Time
```

#### type SelectorEvent

```go
type SelectorEvent interface {
	events.Event

	Selector() string
}
```

This interface represents a Selector-specific event.

NOTE: `golangci-lint` will want this to be called `Event`. this is a bad idea
because this type is explicitly for selector-specific events, and should not be
confused with `events.Event`.

#### type Table

```go
type Table struct {
}
```

Map selector names to method implementations.

#### func  NewTable

```go
func NewTable() *Table
```

#### func (*Table) AddAfter

```go
func (st *Table) AddAfter(selector string, method Method) error
```

#### func (*Table) AddAfterWithPriority

```go
func (st *Table) AddAfterWithPriority(priority int, selector string, method Method) error
```

#### func (*Table) AddBefore

```go
func (st *Table) AddBefore(selector string, method Method) error
```

#### func (*Table) AddBeforeWithPriority

```go
func (st *Table) AddBeforeWithPriority(priority int, selector string, method Method) error
```

#### func (*Table) AllMetadata

```go
func (st *Table) AllMetadata() map[string]map[string]string
```

#### func (*Table) Get

```go
func (st *Table) Get(name string) (*Entry, bool)
```
Return an entry for a selector.

#### func (*Table) HasSelector

```go
func (st *Table) HasSelector(selector string) bool
```
Check whether a selector is defined.

#### func (*Table) InvokeSelector

```go
func (st *Table) InvokeSelector(sel string, tgt responder.Respondable, evt events.Event) (events.Event, bool)
```

#### func (*Table) InvokeSelectorAsync

```go
func (st *Table) InvokeSelectorAsync(
	ctx context.Context,
	sel string,
	tgt responder.Respondable, evt events.Event,
) <-chan events.Event
```

#### func (*Table) ListMetadata

```go
func (st *Table) ListMetadata(selector string) (map[string]string, error)
```

#### func (*Table) Metadata

```go
func (st *Table) Metadata(selector string) (metadata.Metadata, error)
```

#### func (*Table) MustMetadata

```go
func (st *Table) MustMetadata(selector string) metadata.Metadata
```

#### func (*Table) Register

```go
func (st *Table) Register(selector string, method Method)
```
Register a method for a selector.

#### func (*Table) SetDefault

```go
func (st *Table) SetDefault(selector string)
```

#### func (*Table) SetMaxForwardDepth

```go
func (st *Table) SetMaxForwardDepth(val int)
```

#### func (*Table) SetPrimary

```go
func (st *Table) SetPrimary(selector string, method Method) (Method, error)
```

#### func (*Table) Unregister

```go
func (st *Table) Unregister(selector string)
```
