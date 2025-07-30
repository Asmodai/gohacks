<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# dag -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/dag"
```

## Usage

```go
var (
	ErrExpectedParams = errors.Base("expected parameters to be given")
	ErrExpectedString = errors.Base("expected a string value")
	ErrMissingParam   = errors.Base("parameter missing")
	ErrUnknownBuiltin = errors.Base("unknown builtin function")
)
```

```go
var (
	ErrUnknownOperator   = errors.Base("unknown operator")
	ErrRuleCompileFailed = errors.Base("rule compilation failed")
)
```

#### func  FormatIsnf

```go
func FormatIsnf(isn, message string, rest ...any) string
```

#### type ActionFn

```go
type ActionFn func(context.Context, DataMap)
```

Action function callback type.

An action callback is a function that takes a single argument containing the
key/value pair map and returns no value.

#### type ActionParams

```go
type ActionParams map[string]any
```

Action parameters type.

A map of key/value pairs that is passed to the action handler.

#### type ActionSpec

```go
type ActionSpec struct {
	Name    string       `json:"name,omitempty"`    // Action name.
	Perform string       `json:"perform,omitempty"` // Function to perform.
	Params  ActionParams `json:"params,omitempty"`  // Parameters.
}
```

Action specification.

#### type Actions

```go
type Actions interface {
	// Build the given builtin functions.
	Builder(string, ActionParams) (ActionFn, error)
}
```

Action builder interface.

The action builder provides a means of compiling JSON or YAML actions into
explicit function objects

The resulting action is a function that takes `context.Context` and `DataMap`
arguments and then performs some sort of user-defined action.

There are two default builtins provided for you:

    `log`:    Log the contents of the parameters to a logger.
    `mutate`: Change value(s) in the parameters.

To use the `log` builtin, you must provide a `logger.Logger` instance in the
context used with the DAG. For this, you can see `logger.SetLogger`.

#### func  NewDefaultActions

```go
func NewDefaultActions() Actions
```

#### type Compiler

```go
type Compiler interface {
	CompileAction(ActionSpec) (ActionFn, error)
	Compile([]RuleSpec) []error
	Evaluate(DataMap)
}
```


#### func  NewCompiler

```go
func NewCompiler(ctx context.Context, builder Actions) Compiler
```

#### type ConditionSpec

```go
type ConditionSpec struct {
	Attribute string `json:"attribute"` // Attribute to check.
	Operator  string `json:"operator"`  // Predicate operator.
	Value     any    `json:"value"`     // Value to check.
}
```

Condition specification.

#### type DataMap

```go
type DataMap map[string]any
```

Base data type.

This is a map of key/value pairs.

#### type EQBuilder

```go
type EQBuilder struct{}
```


#### func (*EQBuilder) Build

```go
func (bld *EQBuilder) Build(key string, val any) Predicate
```

#### func (*EQBuilder) Token

```go
func (bld *EQBuilder) Token() string
```

#### type EQPredicate

```go
type EQPredicate struct {
	MetaPredicate
}
```


#### func (*EQPredicate) Eval

```go
func (pred *EQPredicate) Eval(input DataMap) bool
```

#### func (*EQPredicate) String

```go
func (pred *EQPredicate) String() string
```

#### type GTBuilder

```go
type GTBuilder struct{}
```


#### func (*GTBuilder) Build

```go
func (bld *GTBuilder) Build(key string, val any) Predicate
```

#### func (*GTBuilder) Token

```go
func (bld *GTBuilder) Token() string
```

#### type GTEBuilder

```go
type GTEBuilder struct{}
```


#### func (*GTEBuilder) Build

```go
func (bld *GTEBuilder) Build(key string, val any) Predicate
```

#### func (*GTEBuilder) Token

```go
func (bld *GTEBuilder) Token() string
```

#### type GTEPredicate

```go
type GTEPredicate struct {
	MetaPredicate
}
```


#### func (*GTEPredicate) Eval

```go
func (pred *GTEPredicate) Eval(input DataMap) bool
```

#### func (*GTEPredicate) String

```go
func (pred *GTEPredicate) String() string
```

#### type GTPredicate

```go
type GTPredicate struct {
	MetaPredicate
}
```


#### func (*GTPredicate) Eval

```go
func (pred *GTPredicate) Eval(input DataMap) bool
```

#### func (*GTPredicate) String

```go
func (pred *GTPredicate) String() string
```

#### type LTBuilder

```go
type LTBuilder struct{}
```


#### func (*LTBuilder) Build

```go
func (bld *LTBuilder) Build(key string, val any) Predicate
```

#### func (*LTBuilder) Token

```go
func (bld *LTBuilder) Token() string
```

#### type LTEBuilder

```go
type LTEBuilder struct{}
```


#### func (*LTEBuilder) Build

```go
func (bld *LTEBuilder) Build(key string, val any) Predicate
```

#### func (*LTEBuilder) Token

```go
func (bld *LTEBuilder) Token() string
```

#### type LTEPredicate

```go
type LTEPredicate struct {
	MetaPredicate
}
```


#### func (*LTEPredicate) Eval

```go
func (pred *LTEPredicate) Eval(input DataMap) bool
```

#### func (*LTEPredicate) String

```go
func (pred *LTEPredicate) String() string
```

#### type LTPredicate

```go
type LTPredicate struct {
	MetaPredicate
}
```


#### func (*LTPredicate) Eval

```go
func (pred *LTPredicate) Eval(input DataMap) bool
```

#### func (*LTPredicate) String

```go
func (pred *LTPredicate) String() string
```

#### type MetaPredicate

```go
type MetaPredicate struct {
}
```


#### func (*MetaPredicate) GetFloatValues

```go
func (meta *MetaPredicate) GetFloatValues(input DataMap) (float64, float64, bool)
```

#### func (*MetaPredicate) GetStringValues

```go
func (meta *MetaPredicate) GetStringValues(input DataMap) (string, string, bool)
```

#### type NEQBuilder

```go
type NEQBuilder struct{}
```


#### func (*NEQBuilder) Build

```go
func (bld *NEQBuilder) Build(key string, val any) Predicate
```

#### func (*NEQBuilder) Token

```go
func (bld *NEQBuilder) Token() string
```

#### type NEQPredicate

```go
type NEQPredicate struct {
	MetaPredicate
}
```


#### func (*NEQPredicate) Eval

```go
func (pred *NEQPredicate) Eval(input DataMap) bool
```

#### func (*NEQPredicate) String

```go
func (pred *NEQPredicate) String() string
```

#### type NOOPPredicate

```go
type NOOPPredicate struct{}
```


#### func (*NOOPPredicate) Eval

```go
func (pred *NOOPPredicate) Eval(_ DataMap) bool
```

#### func (*NOOPPredicate) String

```go
func (pred *NOOPPredicate) String() string
```

#### type Predicate

```go
type Predicate interface {
	Eval(DataMap) bool
	String() string
}
```


#### type PredicateBuilder

```go
type PredicateBuilder interface {
	Token() string
	Build(string, any) Predicate
}
```


#### type PredicateDict

```go
type PredicateDict map[string]PredicateBuilder
```


#### func  BuildPredicateDict

```go
func BuildPredicateDict() PredicateDict
```

#### type PredicateFn

```go
type PredicateFn func(string, any) Predicate
```

Predicate function type.

A predicate is a function that answers a yes-or-no question. In other words: any
expression that can boil down to a boolean.

#### type REIMBuilder

```go
type REIMBuilder struct{}
```


#### func (*REIMBuilder) Build

```go
func (bld *REIMBuilder) Build(key string, val any) Predicate
```

#### func (*REIMBuilder) Token

```go
func (bld *REIMBuilder) Token() string
```

#### type REIMPredicate

```go
type REIMPredicate struct {
	MetaPredicate
}
```


#### func (*REIMPredicate) Eval

```go
func (pred *REIMPredicate) Eval(input DataMap) bool
```

#### func (*REIMPredicate) String

```go
func (pred *REIMPredicate) String() string
```

#### type RESMBuilder

```go
type RESMBuilder struct{}
```


#### func (*RESMBuilder) Build

```go
func (bld *RESMBuilder) Build(key string, val any) Predicate
```

#### func (*RESMBuilder) Token

```go
func (bld *RESMBuilder) Token() string
```

#### type RESMPredicate

```go
type RESMPredicate struct {
	MetaPredicate
}
```


#### func (*RESMPredicate) Eval

```go
func (pred *RESMPredicate) Eval(input DataMap) bool
```

#### func (*RESMPredicate) String

```go
func (pred *RESMPredicate) String() string
```

#### type RuleSpec

```go
type RuleSpec struct {
	Name       string          `json:"name"`       // Rule name.
	Conditions []ConditionSpec `json:"conditions"` // List of conditions.
	Action     ActionSpec      `json:"action"`     // Action to evaluate.
}
```

Filter rule specification.

#### type SIEQBuilder

```go
type SIEQBuilder struct{}
```


#### func (*SIEQBuilder) Build

```go
func (bld *SIEQBuilder) Build(key string, val any) Predicate
```

#### func (*SIEQBuilder) Token

```go
func (bld *SIEQBuilder) Token() string
```

#### type SIEQPredicate

```go
type SIEQPredicate struct {
	MetaPredicate
}
```


#### func (*SIEQPredicate) Eval

```go
func (pred *SIEQPredicate) Eval(input DataMap) bool
```

#### func (*SIEQPredicate) String

```go
func (pred *SIEQPredicate) String() string
```

#### type SINEQBuilder

```go
type SINEQBuilder struct{}
```


#### func (*SINEQBuilder) Build

```go
func (bld *SINEQBuilder) Build(key string, val any) Predicate
```

#### func (*SINEQBuilder) Token

```go
func (bld *SINEQBuilder) Token() string
```

#### type SINEQPredicate

```go
type SINEQPredicate struct {
	MetaPredicate
}
```


#### func (*SINEQPredicate) Eval

```go
func (pred *SINEQPredicate) Eval(input DataMap) bool
```

#### func (*SINEQPredicate) String

```go
func (pred *SINEQPredicate) String() string
```

#### type SSEQBuilder

```go
type SSEQBuilder struct{}
```


#### func (*SSEQBuilder) Build

```go
func (bld *SSEQBuilder) Build(key string, val any) Predicate
```

#### func (*SSEQBuilder) Token

```go
func (bld *SSEQBuilder) Token() string
```

#### type SSEQPredicate

```go
type SSEQPredicate struct {
	MetaPredicate
}
```


#### func (*SSEQPredicate) Eval

```go
func (pred *SSEQPredicate) Eval(input DataMap) bool
```

#### func (*SSEQPredicate) String

```go
func (pred *SSEQPredicate) String() string
```

#### type SSNEQBuilder

```go
type SSNEQBuilder struct{}
```


#### func (*SSNEQBuilder) Build

```go
func (bld *SSNEQBuilder) Build(key string, val any) Predicate
```

#### func (*SSNEQBuilder) Token

```go
func (bld *SSNEQBuilder) Token() string
```

#### type SSNEQPredicate

```go
type SSNEQPredicate struct {
	MetaPredicate
}
```


#### func (*SSNEQPredicate) Eval

```go
func (pred *SSNEQPredicate) Eval(input DataMap) bool
```

#### func (*SSNEQPredicate) String

```go
func (pred *SSNEQPredicate) String() string
```
