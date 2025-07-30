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

#### type Predicate

```go
type Predicate struct {
	Eval func(*Predicate, DataMap) bool // Function to evaluate.
	Data any                            // Data to test against.
}
```


#### type PredicateFn

```go
type PredicateFn func(string, any) Predicate
```

Predicate function type.

A predicate is a function that answers a yes-or-no question. In other words: any
expression that can boil down to a boolean.

#### type RuleSpec

```go
type RuleSpec struct {
	Name       string          `json:"name"`       // Rule name.
	Conditions []ConditionSpec `json:"conditions"` // List of conditions.
	Action     ActionSpec      `json:"action"`     // Action to evaluate.
}
```

Filter rule specification.
