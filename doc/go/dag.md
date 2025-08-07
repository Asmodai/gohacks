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

```go
var (
	ErrInvalidRegexp  = errors.Base("invalid regexp")
	ErrRegexpParse    = errors.Base("error parsing regexp")
	ErrValueNotString = errors.Base("value is not a string")
)
```

#### func  DumpRulesToYAML

```go
func DumpRulesToYAML(rules []RuleSpec) (string, error)
```
Dump a slice of rule specifications to YAML format.

#### func  ExportToDOT

```go
func ExportToDOT(writer io.Writer, root *node)
```
Generate a Graphviz visualisation of the DAG starting from the given node.

#### func  FormatDebugIsnf

```go
func FormatDebugIsnf(isn, message string, rest ...any) string
```
Pretty-print predicate information for debugging.

This includes the token.

#### func  FormatIsnf

```go
func FormatIsnf(message string, rest ...any) string
```
Pretty-print a predicate.

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
Create a new empty `Actions` object.

#### type Compiler

```go
type Compiler interface {
	CompileAction(ActionSpec) (ActionFn, error)
	Compile([]RuleSpec) []error
	Evaluate(Filterable)
	Export(io.Writer)
}
```


#### func  NewCompiler

```go
func NewCompiler(ctx context.Context, build Actions) Compiler
```

#### func  NewCompilerWithPredicates

```go
func NewCompilerWithPredicates(
	ctx context.Context,
	builder Actions,
	predicates PredicateDict,
) Compiler
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

Data Input.

This structure holds the data against which we wish to filter.

It is passed to functions such as `Compiler.Evaluate` as the input.

#### func  NewDataInput

```go
func NewDataInput() *DataInput
```
Create a new empty `DataInput` object.

#### func  NewDataInputFromMap

```go
func NewDataInputFromMap(input map[string]any) *DataInput
```
Create a new `DataInput` object with a copy of the provided input map.

#### func (*DataInput) Get

```go
func (input *DataInput) Get(key string) (any, bool)
```
Get the value of a field.

#### func (*DataInput) Keys

```go
func (input *DataInput) Keys() []string
```
Return a list of field names as keys.

#### func (*DataInput) Set

```go
func (input *DataInput) Set(key string, value any) bool
```
Set the value of the given field to the given value.

#### func (*DataInput) String

```go
func (input *DataInput) String() string
```
Returns the string representation.

#### type EIRBuilder

```go
type EIRBuilder struct{}
```


#### func (*EIRBuilder) Build

```go
func (bld *EIRBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

EIR - Exclusive In Range predicate.

Returne true if the input value is in the filter range inclusive.

#### func (*EIRPredicate) Debug

```go
func (pred *EIRPredicate) Debug() string
```

#### func (*EIRPredicate) Eval

```go
func (pred *EIRPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*EIRPredicate) Instruction

```go
func (pred *EIRPredicate) Instruction() string
```

#### func (*EIRPredicate) String

```go
func (pred *EIRPredicate) String() string
```

#### func (*EIRPredicate) Token

```go
func (pred *EIRPredicate) Token() string
```

#### type EQBuilder

```go
type EQBuilder struct{}
```


#### func (*EQBuilder) Build

```go
func (bld *EQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

EQ - Numeric equality predicate.

Returns true if the input value matches the filter value.

#### func (*EQPredicate) Debug

```go
func (pred *EQPredicate) Debug() string
```

#### func (*EQPredicate) Eval

```go
func (pred *EQPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*EQPredicate) Instruction

```go
func (pred *EQPredicate) Instruction() string
```

#### func (*EQPredicate) String

```go
func (pred *EQPredicate) String() string
```

#### func (*EQPredicate) Token

```go
func (pred *EQPredicate) Token() string
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

	// Return a string representation.
	String() string
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
func (bld *GTBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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
func (bld *GTEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

GTE - Numeric Greater-Than-or-Equal-To predicate.

Returns true if the input value is greater than or equal to the filter value.

#### func (*GTEPredicate) Debug

```go
func (pred *GTEPredicate) Debug() string
```

#### func (*GTEPredicate) Eval

```go
func (pred *GTEPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*GTEPredicate) Instruction

```go
func (pred *GTEPredicate) Instruction() string
```

#### func (*GTEPredicate) String

```go
func (pred *GTEPredicate) String() string
```

#### func (*GTEPredicate) Token

```go
func (pred *GTEPredicate) Token() string
```

#### type GTPredicate

```go
type GTPredicate struct {
	MetaPredicate
}
```

GT - Numeric Greater-Than predicate.

Returns true of the input value is greater than the filter value.

#### func (*GTPredicate) Debug

```go
func (pred *GTPredicate) Debug() string
```

#### func (*GTPredicate) Eval

```go
func (pred *GTPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*GTPredicate) Instruction

```go
func (pred *GTPredicate) Instruction() string
```

#### func (*GTPredicate) String

```go
func (pred *GTPredicate) String() string
```

#### func (*GTPredicate) Token

```go
func (pred *GTPredicate) Token() string
```

#### type IIRBuilder

```go
type IIRBuilder struct{}
```


#### func (*IIRBuilder) Build

```go
func (bld *IIRBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

IIR - Inclusive In Range predicate.

Returns true if the input value is in the range defined in the filter inclusive.

#### func (*IIRPredicate) Debug

```go
func (pred *IIRPredicate) Debug() string
```

#### func (*IIRPredicate) Eval

```go
func (pred *IIRPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*IIRPredicate) Instruction

```go
func (pred *IIRPredicate) Instruction() string
```

#### func (*IIRPredicate) String

```go
func (pred *IIRPredicate) String() string
```

#### func (*IIRPredicate) Token

```go
func (pred *IIRPredicate) Token() string
```

#### type LTBuilder

```go
type LTBuilder struct{}
```


#### func (*LTBuilder) Build

```go
func (bld *LTBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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
func (bld *LTEBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

LTE - Numeric Less-Than-or-Equal-To predicate.

Returns true if the input value is lesser than or equal to the filter value.

#### func (*LTEPredicate) Debug

```go
func (pred *LTEPredicate) Debug() string
```

#### func (*LTEPredicate) Eval

```go
func (pred *LTEPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*LTEPredicate) Instruction

```go
func (pred *LTEPredicate) Instruction() string
```

#### func (*LTEPredicate) String

```go
func (pred *LTEPredicate) String() string
```

#### func (*LTEPredicate) Token

```go
func (pred *LTEPredicate) Token() string
```

#### type LTPredicate

```go
type LTPredicate struct {
	MetaPredicate
}
```

LT - Numeric Less-Than predicate.

Returns true if the input value is lesser than the filter value.

#### func (*LTPredicate) Debug

```go
func (pred *LTPredicate) Debug() string
```

#### func (*LTPredicate) Eval

```go
func (pred *LTPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*LTPredicate) Instruction

```go
func (pred *LTPredicate) Instruction() string
```

#### func (*LTPredicate) String

```go
func (pred *LTPredicate) String() string
```

#### func (*LTPredicate) Token

```go
func (pred *LTPredicate) Token() string
```

#### type MetaPredicate

```go
type MetaPredicate struct {
}
```

A `meta` predicate used by all predicates.

The meta preducate presents common fields and methods so as to avoid duplicate
code.

#### func (*MetaPredicate) Debug

```go
func (meta *MetaPredicate) Debug(isn, token string) string
```

#### func (*MetaPredicate) EvalExclusiveRange

```go
func (meta *MetaPredicate) EvalExclusiveRange(input Filterable) bool
```
Does the predicate's input value fall within the exclusive range defined in the
predicate's filter value?

#### func (*MetaPredicate) EvalInclusiveRange

```go
func (meta *MetaPredicate) EvalInclusiveRange(input Filterable) bool
```
Does the predicate's input value fall within the inclusive range defined in the
predicate's filter value?

#### func (*MetaPredicate) EvalStringMember

```go
func (meta *MetaPredicate) EvalStringMember(input Filterable, insens bool) bool
```
Is the predicate's input value a member of the array of strings in the
predicate's filter value?

#### func (*MetaPredicate) GetFloatValueFromInput

```go
func (meta *MetaPredicate) GetFloatValueFromInput(input Filterable) (float64, bool)
```
Return the predicate's input value as a 64-bit float.

This will return the value for the key on which the predicate operates.

#### func (*MetaPredicate) GetFloatValues

```go
func (meta *MetaPredicate) GetFloatValues(input Filterable) (float64, float64, bool)
```
Return both the predicate's input value and filter value as a 64-bit float.

#### func (*MetaPredicate) GetPredicateFloatArray

```go
func (meta *MetaPredicate) GetPredicateFloatArray() ([]float64, bool)
```
Return the predicate's filter value as an array of 64-bit floats.

#### func (*MetaPredicate) GetPredicateStringArray

```go
func (meta *MetaPredicate) GetPredicateStringArray() ([]string, bool)
```
Return the predicate's filter value as an array of strings.

#### func (*MetaPredicate) GetStringValues

```go
func (meta *MetaPredicate) GetStringValues(input Filterable) (string, string, bool)
```
Return both the predicate's input value and filter value as a string.

#### func (*MetaPredicate) String

```go
func (meta *MetaPredicate) String(token string) string
```

#### type NEQBuilder

```go
type NEQBuilder struct{}
```


#### func (*NEQBuilder) Build

```go
func (bld *NEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

NEQ - Numeric Inequality predicate.

Returns true if the input value is not equal to the filter value.

#### func (*NEQPredicate) Debug

```go
func (pred *NEQPredicate) Debug() string
```

#### func (*NEQPredicate) Eval

```go
func (pred *NEQPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*NEQPredicate) Instruction

```go
func (pred *NEQPredicate) Instruction() string
```

#### func (*NEQPredicate) String

```go
func (pred *NEQPredicate) String() string
```

#### func (*NEQPredicate) Token

```go
func (pred *NEQPredicate) Token() string
```

#### type NOOPPredicate

```go
type NOOPPredicate struct{}
```

NOOP - No operation.

This predicate is used internally to represent the root node.

It always returns true.

#### func (*NOOPPredicate) Debug

```go
func (pred *NOOPPredicate) Debug() string
```

#### func (*NOOPPredicate) Eval

```go
func (pred *NOOPPredicate) Eval(_ context.Context, _ Filterable) bool
```

#### func (*NOOPPredicate) Instruction

```go
func (pred *NOOPPredicate) Instruction() string
```

#### func (*NOOPPredicate) String

```go
func (pred *NOOPPredicate) String() string
```

#### func (*NOOPPredicate) Token

```go
func (pred *NOOPPredicate) Token() string
```

#### type Predicate

```go
type Predicate interface {
	// Evaluate the predicate against the given `Filterable` object.
	//
	// Returns the result of the predicate.
	Eval(context.Context, Filterable) bool

	// Return the string representation of the predicate.
	String() string

	// Return the instruction name for the predicate.
	//
	// This isn't used in the current version of the directed acyclic
	// graph, but the theory is that this could be used in a tokeniser
	// or as opcode.
	//
	// The value this returns must be unique.
	Instruction() string

	// Return the token name for the predicate.
	//
	// This is the string value used in the action specification.
	Token() string

	// Return a string representation for debugging.
	Debug() string
}
```

Predicate interface.

All predicates must adhere to this interface.

#### type PredicateBuilder

```go
type PredicateBuilder interface {
	// Return the token name for the predicate.
	//
	// This isn't used in the current version of the directed acyclic
	// graph, but the theory is that this could be used in a tokeniser
	// or as opcode.
	//
	// The value this returns must be unique.
	Token() string

	// Build a new predicate.
	//
	// This will create a predicate that operates on the given field
	// and data.
	Build(field string, data any, lgr logger.Logger, dbg bool) (Predicate, error)
}
```

Predicate builder interface.

All predicate builders must adhere to this interface.

#### type PredicateDict

```go
type PredicateDict map[string]PredicateBuilder
```

Dictionary of available predicate builders.

#### func  BuildPredicateDict

```go
func BuildPredicateDict() PredicateDict
```
Build the predicate dictionary for the directed acyclic graph filter.

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
func (bld *REIMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

REIM - Regular Expression (Insensitive) Match predicate.

Returns true if the regular expression in the filter matches against the input
value.

The regular expression will be compiled with a prefix denoting that it does not
care about case.

#### func (*REIMPredicate) Debug

```go
func (pred *REIMPredicate) Debug() string
```

#### func (*REIMPredicate) Eval

```go
func (pred *REIMPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*REIMPredicate) Instruction

```go
func (pred *REIMPredicate) Instruction() string
```

#### func (*REIMPredicate) String

```go
func (pred *REIMPredicate) String() string
```

#### func (*REIMPredicate) Token

```go
func (pred *REIMPredicate) Token() string
```

#### type RESMBuilder

```go
type RESMBuilder struct{}
```


#### func (*RESMBuilder) Build

```go
func (bld *RESMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

RESM - Regular Expression (Sensitive) Match predicate.

Returns true if the regular expression in the filter value matches against the
input value.

The regular expression will not be forced into being case-insensitive.

#### func (*RESMPredicate) Debug

```go
func (pred *RESMPredicate) Debug() string
```

#### func (*RESMPredicate) Eval

```go
func (pred *RESMPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*RESMPredicate) Instruction

```go
func (pred *RESMPredicate) Instruction() string
```

#### func (*RESMPredicate) String

```go
func (pred *RESMPredicate) String() string
```

#### func (*RESMPredicate) Token

```go
func (pred *RESMPredicate) Token() string
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
Parse a rule specification from a string containing JSON.

#### func  ParseFromYAML

```go
func ParseFromYAML(data string) ([]RuleSpec, error)
```
Parse a rule specification from a string containing YAML.

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
func (bld *SIEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

SIEG - String (Insensitive) Equality predicate.

Returns true if the filter value matches the input value.

This predicate does not care about case.

#### func (*SIEQPredicate) Debug

```go
func (pred *SIEQPredicate) Debug() string
```

#### func (*SIEQPredicate) Eval

```go
func (pred *SIEQPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*SIEQPredicate) Instruction

```go
func (pred *SIEQPredicate) Instruction() string
```

#### func (*SIEQPredicate) String

```go
func (pred *SIEQPredicate) String() string
```

#### func (*SIEQPredicate) Token

```go
func (pred *SIEQPredicate) Token() string
```

#### type SIMBuilder

```go
type SIMBuilder struct{}
```


#### func (*SIMBuilder) Build

```go
func (bld *SIMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

SIM - String (Insensitive) Member predicate.

Returns true if the input value is a member of the string array in the filter
value.

#### func (*SIMPredicate) Debug

```go
func (pred *SIMPredicate) Debug() string
```

#### func (*SIMPredicate) Eval

```go
func (pred *SIMPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*SIMPredicate) Instruction

```go
func (pred *SIMPredicate) Instruction() string
```

#### func (*SIMPredicate) String

```go
func (pred *SIMPredicate) String() string
```

#### func (*SIMPredicate) Token

```go
func (pred *SIMPredicate) Token() string
```

#### type SINEQBuilder

```go
type SINEQBuilder struct{}
```


#### func (*SINEQBuilder) Build

```go
func (bld *SINEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

SINEQ - String (Insensitive) Inequality predicate.

Returns true if the input string is not the same as the filter string.

Case is not taken into account.

#### func (*SINEQPredicate) Debug

```go
func (pred *SINEQPredicate) Debug() string
```

#### func (*SINEQPredicate) Eval

```go
func (pred *SINEQPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*SINEQPredicate) Instruction

```go
func (pred *SINEQPredicate) Instruction() string
```

#### func (*SINEQPredicate) String

```go
func (pred *SINEQPredicate) String() string
```

#### func (*SINEQPredicate) Token

```go
func (pred *SINEQPredicate) Token() string
```

#### type SSEQBuilder

```go
type SSEQBuilder struct{}
```


#### func (*SSEQBuilder) Build

```go
func (bld *SSEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

SSEQ - String (Sensitive) Equality predicate.

Returns true if the input value is the same as the filter value.

#### func (*SSEQPredicate) Debug

```go
func (pred *SSEQPredicate) Debug() string
```

#### func (*SSEQPredicate) Eval

```go
func (pred *SSEQPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*SSEQPredicate) Instruction

```go
func (pred *SSEQPredicate) Instruction() string
```

#### func (*SSEQPredicate) String

```go
func (pred *SSEQPredicate) String() string
```

#### func (*SSEQPredicate) Token

```go
func (pred *SSEQPredicate) Token() string
```

#### type SSMBuilder

```go
type SSMBuilder struct{}
```


#### func (*SSMBuilder) Build

```go
func (bld *SSMBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

SSM - String (Sensitive) Member predicate.

Returns true if the input value is a member of the string array in the filter
value.

#### func (*SSMPredicate) Debug

```go
func (pred *SSMPredicate) Debug() string
```

#### func (*SSMPredicate) Eval

```go
func (pred *SSMPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*SSMPredicate) Instruction

```go
func (pred *SSMPredicate) Instruction() string
```

#### func (*SSMPredicate) String

```go
func (pred *SSMPredicate) String() string
```

#### func (*SSMPredicate) Token

```go
func (pred *SSMPredicate) Token() string
```

#### type SSNEQBuilder

```go
type SSNEQBuilder struct{}
```


#### func (*SSNEQBuilder) Build

```go
func (bld *SSNEQBuilder) Build(key string, val any, lgr logger.Logger, dbg bool) (Predicate, error)
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

SSNEQ - String (Sensitive) Inequality predicate.

Returns true if the input value is different to the filter value.

#### func (*SSNEQPredicate) Debug

```go
func (pred *SSNEQPredicate) Debug() string
```

#### func (*SSNEQPredicate) Eval

```go
func (pred *SSNEQPredicate) Eval(_ context.Context, input Filterable) bool
```

#### func (*SSNEQPredicate) Instruction

```go
func (pred *SSNEQPredicate) Instruction() string
```

#### func (*SSNEQPredicate) String

```go
func (pred *SSNEQPredicate) String() string
```

#### func (*SSNEQPredicate) Token

```go
func (pred *SSNEQPredicate) Token() string
```
