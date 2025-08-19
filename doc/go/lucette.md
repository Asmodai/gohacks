<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# lucette -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/lucette"
```

## Usage

```go
const (
	PK_EQ_S PredKind = iota
	PK_NEQ_S
	PK_PREFIX
	PK_GLOB
	PK_REGEX
	PK_PHRASE
	PK_EXISTS
	PK_CMP
	PK_RANGE

	CmpLT CmpKind = iota
	CmpLTE
	CmpGT
	CmpGTE
	CmpEQ
	CmpNEQ
)
```

```go
var (
	ErrLabelMissingID = errors.Base("LABEL missing id")
	ErrLabelBadIDType = errors.Base("LABEL has bad id type")
	ErrJumpMissingArg = errors.Base("jump missing target arg")
	ErrJumpNotLabelID = errors.Base("jump target arg not LabelID")
	ErrUnboundLabel   = errors.Base("unbound label")
)
```

```go
var (
	ErrUnterminatedRegex  = errors.Base("unterminated regular expression")
	ErrNewlineInRegex     = errors.Base("embedded newline in regular expression")
	ErrRegexFlags         = errors.Base("regex flags not supported")
	ErrUnexpectedToken    = errors.Base("unexpected token")
	ErrNewlineInPhrase    = errors.Base("embedded newline in phrase")
	ErrUnterminatedField  = errors.Base("unterminated quoted field")
	ErrNewlineInField     = errors.Base("embedded newline in field")
	ErrUnterminatedString = errors.Base("unterminated string")
	ErrUnexpectedRune     = errors.Base("unexpected rune")
	ErrDoubleUnread       = errors.Base("double unread")
	ErrUnexpectedBareword = errors.Base("unexpected bareword (missing quotes or field?)")
)
```

```go
var (
	ErrBadDateTime    = errors.Base("bad datetime")
	ErrUnknownLiteral = errors.Base("unknown literal")
)
```

#### type BoolOp

```go
type BoolOp int
```


```go
const (
	OpAnd BoolOp = iota
	OpOr
)
```

#### type CmpKind

```go
type CmpKind int
```

Comparator kind type.

#### type Comparator

```go
type Comparator struct {
	Op   CmpKind // Comparison operator.
	Atom NodeLit // Atom to compare.
}
```

Comparator structure.

#### type Diagnostic

```go
type Diagnostic struct {
	Msg  string // Diagnostic message.
	At   Span   // Location within token stream.
	Hint string // Hint message, if applicable.
}
```

Diagnostic message.

#### func (*Diagnostic) String

```go
func (d *Diagnostic) String() string
```
Pretty-print a diagnostic to a string.

#### type FieldSpec

```go
type FieldSpec struct {
	Name     string
	FType    FieldType
	Analyser string
	Layouts  []string
}
```


#### type FieldType

```go
type FieldType int
```


```go
const (
	FTKeyword FieldType = iota
	FTText
	FTNumeric
	FTDateTime
	FTIP
)
```

#### type Instr

```go
type Instr struct {
	Op   OpCode
	Args []any
}
```


#### func (Instr) String

```go
func (isn Instr) String() string
```

#### type LabelID

```go
type LabelID int
```


#### type Lexer

```go
type Lexer struct {
}
```

The lexer.

You are not meant to use this directly.

#### func  NewLexer

```go
func NewLexer(reader io.Reader) *Lexer
```

#### type LitKind

```go
type LitKind int
```

Literal kind type.

```go
const (
	LString LitKind = iota
	LNumber
	LUnbounded
)
```

#### type ModKind

```go
type ModKind int
```

Modifier kind type.

```go
const (
	ModRequire ModKind = iota
	ModProhibit
)
```

#### type Node

```go
type Node interface {
	// Return the span for this node.
	//
	// Spans can be used in diagnostics to show where in the source file
	// an issue exists.
	Span() Span

	// Print debugging information for the given node.
	Debug(...any) *debug.Debug
}
```

AST node.

#### type NodeAnd

```go
type NodeAnd struct {
}
```

An AST node for the `AND' logical operator.

#### func (*NodeAnd) Debug

```go
func (n *NodeAnd) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*NodeAnd) Span

```go
func (n *NodeAnd) Span() Span
```
Return the span for the AND within the source code.

#### type NodeLit

```go
type NodeLit struct {
}
```

An AST node representing a literal.

#### func (*NodeLit) Debug

```go
func (n *NodeLit) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*NodeLit) Span

```go
func (n *NodeLit) Span() Span
```
Return the span for the node.

#### type NodeMod

```go
type NodeMod struct {
}
```

An AST node representing a modifier.

#### func (*NodeMod) Debug

```go
func (n *NodeMod) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*NodeMod) Span

```go
func (n *NodeMod) Span() Span
```
Return the span for the modifier.

#### type NodeNot

```go
type NodeNot struct {
}
```

An AST node for the `NOT' logical operator.

#### func (*NodeNot) Debug

```go
func (n *NodeNot) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*NodeNot) Span

```go
func (n *NodeNot) Span() Span
```
Return the span for the NOT within the source code.

#### type NodeOr

```go
type NodeOr struct {
}
```

An AST node for the `OR' logical operator.

#### func (*NodeOr) Debug

```go
func (n *NodeOr) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*NodeOr) Span

```go
func (n *NodeOr) Span() Span
```
Return the span for the OR within the source code.

#### type NodePred

```go
type NodePred struct {
}
```

An AST node representing a predicate.

#### func (*NodePred) Debug

```go
func (n *NodePred) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*NodePred) Span

```go
func (n *NodePred) Span() Span
```
Return the span.

#### type OpCode

```go
type OpCode int
```


#### type PPRangeType

```go
type PPRangeType int
```


```go
const (
	PPRangeN PPRangeType = iota
	PPRangeT
	PPRangeIP
)
```

#### type Parser

```go
type Parser struct {
}
```

Parser structure.

#### func  NewParser

```go
func NewParser(toks []Token) *Parser
```
Create a new parser instance.

#### func (*Parser) Parse

```go
func (p *Parser) Parse() (Node, []Diagnostic)
```
Parse the token stream.

#### type Position

```go
type Position struct {
}
```

Position within source code.

#### func (*Position) String

```go
func (p *Position) String() string
```
Pretty-print the position as a string.

#### type PredKind

```go
type PredKind int
```

Predicate kind type.

#### type PrettyPrinterOpts

```go
type PrettyPrinterOpts struct {
	WithComments bool // Include decoded comments.
	AddrWidth    int  // Width of address digits.  0 = auto.
	OpcodeWidth  int  // Pad opcode column. 0 = auto.
	OperandWidth int  // Pad operand column. 0 = auto.
}
```


#### func  NewDefaultPrettyPrinterOptions

```go
func NewDefaultPrettyPrinterOptions() PrettyPrinterOpts
```

#### type Program

```go
type Program struct {
	Fields   []string
	Strings  []string
	Numbers  []float64
	Times    []int64
	IPs      []netip.Addr
	Patterns []*regexp.Regexp
	Code     []Instr
}
```


#### func  NewProgram

```go
func NewProgram() *Program
```

#### func (*Program) Emit

```go
func (p *Program) Emit(node TypedNode)
```

#### func (*Program) Peephole

```go
func (p *Program) Peephole()
```

#### func (*Program) PrettyPrint

```go
func (p *Program) PrettyPrint(writer io.Writer, opts PrettyPrinterOpts)
```

#### type Range

```go
type Range struct {
	Low  *NodeLit // Low literal.
	High *NodeLit // High literal.
	IncL bool     // Low is inclusive?
	IncH bool     // High is inclusive?
}
```

Range structure.

#### type Schema

```go
type Schema map[string]FieldSpec
```


#### type Span

```go
type Span struct {
}
```

Source code span.

#### func (*Span) String

```go
func (s *Span) String() string
```
Return the string representation of a span.

#### type Token

```go
type Token struct {
}
```


#### func (*Token) String

```go
func (t *Token) String() string
```

#### type TokenType

```go
type TokenType uint
```

Lexer token.

#### func (TokenType) Literal

```go
func (t TokenType) Literal() string
```

#### func (TokenType) String

```go
func (t TokenType) String() string
```
Stringer method for lexer tokens.

#### type TypedNode

```go
type TypedNode interface {
	Key() string
	Debug(...any) *debug.Debug
}
```


#### func  Simplify

```go
func Simplify(node TypedNode) TypedNode
```

#### func  ToNNF

```go
func ToNNF(node TypedNode) TypedNode
```

#### type TypedNodeAnd

```go
type TypedNodeAnd struct {
}
```


#### func (*TypedNodeAnd) Debug

```go
func (n *TypedNodeAnd) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeAnd) Key

```go
func (n *TypedNodeAnd) Key() string
```

#### type TypedNodeCmpIP

```go
type TypedNodeCmpIP struct {
}
```


#### func (*TypedNodeCmpIP) Debug

```go
func (n *TypedNodeCmpIP) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeCmpIP) Key

```go
func (n *TypedNodeCmpIP) Key() string
```

#### type TypedNodeCmpN

```go
type TypedNodeCmpN struct {
}
```


#### func (*TypedNodeCmpN) Debug

```go
func (n *TypedNodeCmpN) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeCmpN) Key

```go
func (n *TypedNodeCmpN) Key() string
```

#### type TypedNodeCmpT

```go
type TypedNodeCmpT struct {
}
```


#### func (*TypedNodeCmpT) Debug

```go
func (n *TypedNodeCmpT) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeCmpT) Key

```go
func (n *TypedNodeCmpT) Key() string
```

#### type TypedNodeEqS

```go
type TypedNodeEqS struct {
}
```


#### func (*TypedNodeEqS) Debug

```go
func (n *TypedNodeEqS) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeEqS) Key

```go
func (n *TypedNodeEqS) Key() string
```

#### type TypedNodeExists

```go
type TypedNodeExists struct {
}
```


#### func (*TypedNodeExists) Debug

```go
func (n *TypedNodeExists) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeExists) Key

```go
func (n *TypedNodeExists) Key() string
```

#### type TypedNodeFalse

```go
type TypedNodeFalse struct {
}
```


#### func (*TypedNodeFalse) Debug

```go
func (n *TypedNodeFalse) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeFalse) Key

```go
func (n *TypedNodeFalse) Key() string
```

#### type TypedNodeGlob

```go
type TypedNodeGlob struct {
}
```


#### func (*TypedNodeGlob) Debug

```go
func (n *TypedNodeGlob) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeGlob) Key

```go
func (n *TypedNodeGlob) Key() string
```

#### type TypedNodeNeqS

```go
type TypedNodeNeqS struct {
}
```


#### func (*TypedNodeNeqS) Debug

```go
func (n *TypedNodeNeqS) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeNeqS) Key

```go
func (n *TypedNodeNeqS) Key() string
```

#### type TypedNodeNot

```go
type TypedNodeNot struct {
}
```


#### func (*TypedNodeNot) Debug

```go
func (n *TypedNodeNot) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeNot) Key

```go
func (n *TypedNodeNot) Key() string
```

#### type TypedNodeOr

```go
type TypedNodeOr struct {
}
```


#### func (*TypedNodeOr) Debug

```go
func (n *TypedNodeOr) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeOr) Key

```go
func (n *TypedNodeOr) Key() string
```

#### type TypedNodePhrase

```go
type TypedNodePhrase struct {
}
```


#### func (*TypedNodePhrase) Debug

```go
func (n *TypedNodePhrase) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodePhrase) Key

```go
func (n *TypedNodePhrase) Key() string
```

#### type TypedNodePrefix

```go
type TypedNodePrefix struct {
}
```


#### func (*TypedNodePrefix) Debug

```go
func (n *TypedNodePrefix) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodePrefix) Key

```go
func (n *TypedNodePrefix) Key() string
```

#### type TypedNodeRangeIP

```go
type TypedNodeRangeIP struct {
}
```


#### func (*TypedNodeRangeIP) Debug

```go
func (n *TypedNodeRangeIP) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeRangeIP) Key

```go
func (n *TypedNodeRangeIP) Key() string
```

#### type TypedNodeRangeN

```go
type TypedNodeRangeN struct {
}
```


#### func (*TypedNodeRangeN) Debug

```go
func (n *TypedNodeRangeN) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeRangeN) Key

```go
func (n *TypedNodeRangeN) Key() string
```

#### type TypedNodeRangeT

```go
type TypedNodeRangeT struct {
}
```


#### func (*TypedNodeRangeT) Debug

```go
func (n *TypedNodeRangeT) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeRangeT) Key

```go
func (n *TypedNodeRangeT) Key() string
```

#### type TypedNodeRegex

```go
type TypedNodeRegex struct {
}
```


#### func (*TypedNodeRegex) Debug

```go
func (n *TypedNodeRegex) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeRegex) Key

```go
func (n *TypedNodeRegex) Key() string
```

#### type TypedNodeTrue

```go
type TypedNodeTrue struct {
}
```


#### func (*TypedNodeTrue) Debug

```go
func (n *TypedNodeTrue) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*TypedNodeTrue) Key

```go
func (n *TypedNodeTrue) Key() string
```

#### type Typer

```go
type Typer struct {
	Sch   Schema
	Diags []Diagnostic
}
```


#### func  NewTyper

```go
func NewTyper(sch Schema) *Typer
```

#### func (*Typer) Type

```go
func (t *Typer) Type(node Node) (TypedNode, []Diagnostic)
```
