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

#### func  DumpRulesToYAML

```go
func DumpRulesToYAML(rules []RuleSpec) (string, error)
```
Dump a slice of rule specifications to YAML format.

#### func  FormatIsnf

```go
func FormatIsnf(isn, message string, rest ...any) string
```

#### type ActionFn

```go
type ActionFn func(context.Context, Filterable)
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
	// Action name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Function to perform.
	Perform string `json:"perform,omitempty" yaml:"perform,omitempty"`

	// Parameters.
	Params ActionParams `json:"params,omitempty" yaml:"params,omitempty"`
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

The resulting action is a function that takes `context.Context` and `Filterable`
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
	Evaluate(Filterable)
}
```


#### func  NewCompiler

```go
func NewCompiler(ctx context.Context, builder Actions) Compiler
```

#### type ConditionSpec

```go
type ConditionSpec struct {
	// Attribute to check.
	Attribute string `json:"attribute" yaml:"attribute"`

	// Predicate operator.
	Operator string `json:"operator" yaml:"operator"`

	// Value to check.
	Value any `json:"value" yaml:"value"`
}
```

Condition specification.

#### type DataInput

```go
type DataInput struct {
}
```


#### func  NewDataInput

```go
func NewDataInput() *DataInput
```

#### func  NewDataInputFromMap

```go
func NewDataInputFromMap(input map[string]any) *DataInput
```
Create a new `DataInput` object with a copy of the provided input map.

#### func (*DataInput) Get

```go
func (input *DataInput) Get(key string) (any, bool)
```

#### func (*DataInput) Keys

```go
func (input *DataInput) Keys() []string
```

#### func (*DataInput) Set

```go
func (input *DataInput) Set(key string, value any) bool
```

#### type EIRBuilder

```go
type EIRBuilder struct{}
```


#### func (*EIRBuilder) Build

```go
func (bld *EIRBuilder) Build(key string, val any) Predicate
```

#### func (*EIRBuilder) Token

```go
func (bld *EIRBuilder) Token() string
```

#### type EIRPredicate

```go
type EIRPredicate struct {
	MetaPredicate
}
```


#### func (*EIRPredicate) Eval

```go
func (pred *EIRPredicate) Eval(input Filterable) bool
```

#### func (*EIRPredicate) String

```go
func (pred *EIRPredicate) String() string
```

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
func (pred *EQPredicate) Eval(input Filterable) bool
```

#### func (*EQPredicate) String

```go
func (pred *EQPredicate) String() string
```

#### type Filterable

```go
type Filterable interface {
	// Get the given key from the filterable entity.
	//
	// If the key exists, it is returned along with `true`.
	//
	// If the key does not exist, `false` returned.
	Get(string) (any, bool)

	// Set the given key to the given value.
	//
	// The graph engine should not add new entries, so if an attempt is
	// made to do so, then `false` is returned and nothing happens.
	Set(string, any) bool

	// Get a list of keys from the filterable entity.
	Keys() []string
}
```

Filterable interface.

This interface allows objects to be used with the direct acyclig graph as input.

A 'filterable' entity provides a means of getting at its field contents so the
DAG can look them up.

A decision was made to avoid `reflect` as the DAG might be in a hot path where
reflection adds too big a performance hit.

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
func (pred *GTEPredicate) Eval(input Filterable) bool
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
func (pred *GTPredicate) Eval(input Filterable) bool
```

#### func (*GTPredicate) String

```go
func (pred *GTPredicate) String() string
```

#### type IIRBuilder

```go
type IIRBuilder struct{}
```


#### func (*IIRBuilder) Build

```go
func (bld *IIRBuilder) Build(key string, val any) Predicate
```

#### func (*IIRBuilder) Token

```go
func (bld *IIRBuilder) Token() string
```

#### type IIRPredicate

```go
type IIRPredicate struct {
	MetaPredicate
}
```


#### func (*IIRPredicate) Eval

```go
func (pred *IIRPredicate) Eval(input Filterable) bool
```

#### func (*IIRPredicate) String

```go
func (pred *IIRPredicate) String() string
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
func (pred *LTEPredicate) Eval(input Filterable) bool
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
func (pred *LTPredicate) Eval(input Filterable) bool
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


#### func (*MetaPredicate) EvalExclusiveRange

```go
func (meta *MetaPredicate) EvalExclusiveRange(input Filterable) bool
```

#### func (*MetaPredicate) EvalInclusiveRange

```go
func (meta *MetaPredicate) EvalInclusiveRange(input Filterable) bool
```

#### func (*MetaPredicate) EvalStringMember

```go
func (meta *MetaPredicate) EvalStringMember(input Filterable, insens bool) bool
```

#### func (*MetaPredicate) GetFloatValueFromInput

```go
func (meta *MetaPredicate) GetFloatValueFromInput(input Filterable) (float64, bool)
```

#### func (*MetaPredicate) GetFloatValues

```go
func (meta *MetaPredicate) GetFloatValues(input Filterable) (float64, float64, bool)
```

#### func (*MetaPredicate) GetPredicateFloatArray

```go
func (meta *MetaPredicate) GetPredicateFloatArray() ([]float64, bool)
```

#### func (*MetaPredicate) GetPredicateStringArray

```go
func (meta *MetaPredicate) GetPredicateStringArray() ([]string, bool)
```

#### func (*MetaPredicate) GetStringValues

```go
func (meta *MetaPredicate) GetStringValues(input Filterable) (string, string, bool)
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
func (pred *NEQPredicate) Eval(input Filterable) bool
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
func (pred *NOOPPredicate) Eval(_ Filterable) bool
```

#### func (*NOOPPredicate) String

```go
func (pred *NOOPPredicate) String() string
```

#### type Predicate

```go
type Predicate interface {
	Eval(Filterable) bool
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
func (pred *REIMPredicate) Eval(input Filterable) bool
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
func (pred *RESMPredicate) Eval(input Filterable) bool
```

#### func (*RESMPredicate) String

```go
func (pred *RESMPredicate) String() string
```

#### type RuleSpec

```go
type RuleSpec struct {
	// Rule name.
	Name string `json:"name" yaml:"name"`

	// List of conditions.
	Conditions []ConditionSpec `json:"conditions" yaml:"conditions"`

	// Action to evaluate.
	Action ActionSpec `json:"action" yaml:"action"`
}
```

Filter rule specification.

#### func  ParseFromJSON

```go
func ParseFromJSON(data string) ([]RuleSpec, error)
```

#### func  ParseFromYAML

```go
func ParseFromYAML(data string) ([]RuleSpec, error)
```
Dump a slice of rule specifications to JSON format.

#### func (*RuleSpec) DumpToJSON

```go
func (rs *RuleSpec) DumpToJSON() (string, error)
```
Dump the rule specification to JSON format.

#### func (*RuleSpec) DumpToYAML

```go
func (rs *RuleSpec) DumpToYAML() (string, error)
```
Dump the rule specification to YAML format.

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
func (pred *SIEQPredicate) Eval(input Filterable) bool
```

#### func (*SIEQPredicate) String

```go
func (pred *SIEQPredicate) String() string
```

#### type SIMBuilder

```go
type SIMBuilder struct{}
```


#### func (*SIMBuilder) Build

```go
func (bld *SIMBuilder) Build(key string, val any) Predicate
```

#### func (*SIMBuilder) Token

```go
func (bld *SIMBuilder) Token() string
```

#### type SIMPredicate

```go
type SIMPredicate struct {
	MetaPredicate
}
```


#### func (*SIMPredicate) Eval

```go
func (pred *SIMPredicate) Eval(input Filterable) bool
```

#### func (*SIMPredicate) String

```go
func (pred *SIMPredicate) String() string
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
func (pred *SINEQPredicate) Eval(input Filterable) bool
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
func (pred *SSEQPredicate) Eval(input Filterable) bool
```

#### func (*SSEQPredicate) String

```go
func (pred *SSEQPredicate) String() string
```

#### type SSMBuilder

```go
type SSMBuilder struct{}
```


#### func (*SSMBuilder) Build

```go
func (bld *SSMBuilder) Build(key string, val any) Predicate
```

#### func (*SSMBuilder) Token

```go
func (bld *SSMBuilder) Token() string
```

#### type SSMPredicate

```go
type SSMPredicate struct {
	MetaPredicate
}
```


#### func (*SSMPredicate) Eval

```go
func (pred *SSMPredicate) Eval(input Filterable) bool
```

#### func (*SSMPredicate) String

```go
func (pred *SSMPredicate) String() string
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
func (pred *SSNEQPredicate) Eval(input Filterable) bool
```

#### func (*SSNEQPredicate) String

```go
func (pred *SSNEQPredicate) String() string
```
