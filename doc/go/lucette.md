<!-- -*- Mode: gfm; auto-fill: t; fill-column: 78; -*- -->

# lucette -- Go Hacks Library

```go
    import "github.com/Asmodai/gohacks/lucette"
```

## Usage

```go
var (
	// Returned when the typer detects an invalid datetime.
	ErrBadDateTime = errors.Base("bad datetime value")

	// Returned should an attempt be made to `unread` after a rune has
	// already been put back into the reader.
	ErrDoubleUnread = errors.Base("double unread")

	// Returned if the lexer is invoked without a valid reader.
	ErrInvalidReader = errors.Base("reader is not valid")

	// Returned when the code generator detects a label without a target.
	ErrJumpMissingArg = errors.Base("jump missing target arg")

	// Returned when the code generator detects a label with an invalid
	// target.
	ErrJumpNotLabelID = errors.Base("jump target arg not LabelID")

	// Returned when the code generator detects a label with a bad ID.
	ErrLabelBadIDType = errors.Base("LABEL has bad id type")

	// Returned when the code generator detects a label that lacks an ID.
	ErrLabelMissingID = errors.Base("LABEL missing id")

	// Returned when the lexer detects an embedded newline in a field
	// name.
	ErrNewlineInField = errors.Base("embedded newline in field")

	// Returned when the lexer detects an embedded newline in a phrase.
	ErrNewlineInPhrase = errors.Base("embedded newline in phrase")

	// Returned when the lexer detects a newline in a regular expression.
	ErrNewlineInRegex = errors.Base("embedded newline in regular expression")

	// Returned if no tokens were provided.
	ErrNoTokens = errors.Base("no tokens")

	// Returned when the lexer detects unsupported flags in a regular
	// expression.
	ErrRegexFlags = errors.Base("regex flags not supported")

	// Returned when the code generator detects a label that has not been
	// bound to a target.
	ErrUnboundLabel = errors.Base("unbound label")

	// Returned when the lexer detects an unexpected bareword in the
	// source code.
	ErrUnexpectedBareword = errors.Base("unexpected bareword (missing quotes or field?)")

	// Returned when the lexer detects an unexpected character.
	ErrUnexpectedRune = errors.Base("unexpected rune")

	// Returned when the lexer detects an unexpected token in the source
	// code.
	ErrUnexpectedToken = errors.Base("unexpected token")

	// Returned when the typer detects an unknown literal.
	ErrUnknownLiteral = errors.Base("unknown literal")

	// Returned when the lexer detects that a quoted field name is
	// unterminated.
	ErrUnterminatedField = errors.Base("unterminated quoted field")

	// Returned when the lexer detects an unterminated regular expression.
	ErrUnterminatedRegex = errors.Base("unterminated regular expression")

	// Returned when the lexer detects an unterminated quoted string.
	ErrUnterminatedString = errors.Base("unterminated string")
)
```

```go
var (
	// A span with zero values.
	ZeroSpan = &Span{}
)
```

#### func  ComparatorKindToString

```go
func ComparatorKindToString(kind ComparatorKind) string
```
Return the string representation for a comparator kind.

#### func  FieldTypeToString

```go
func FieldTypeToString(fType FieldType) string
```

#### func  LiteralKindToString

```go
func LiteralKindToString(lit LiteralKind) string
```
Return the string representation of a literal type.

#### func  ModifierKindToString

```go
func ModifierKindToString(kind ModifierKind) string
```
Return the string representation of a modifier.

#### func  PredicateKindToString

```go
func PredicateKindToString(kind PredicateKind) string
```
Return the string representation for a predicate kind.

#### type ASTAnd

```go
type ASTAnd struct {
	Kids []ASTNode // Child nodes.
}
```

An AST node for the `AND' logical operator.

#### func (ASTAnd) Debug

```go
func (n ASTAnd) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (ASTAnd) Span

```go
func (n ASTAnd) Span() *Span
```
Return the span for the AST node.

#### type ASTComparator

```go
type ASTComparator struct {
	Atom ASTLiteral     // Atom on which to operate.
	Op   ComparatorKind // Comparator operator.
}
```

Comparator structure.

#### func (ASTComparator) Debug

```go
func (c ASTComparator) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### type ASTLiteral

```go
type ASTLiteral struct {
	String string // String value.

	Kind   LiteralKind // Kind of the literal.
	Number float64     // Numeric value.
}
```

An AST node for a `literal' of some kind.

#### func (ASTLiteral) Debug

```go
func (n ASTLiteral) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (ASTLiteral) Span

```go
func (n ASTLiteral) Span() *Span
```
Return the span for the AST node.

#### type ASTModifier

```go
type ASTModifier struct {
	Kid ASTNode // Node to which the modifier applies.

	Kind ModifierKind // Modifier kind.
}
```

An AST node for a `modifier' to an operation..

#### func (ASTModifier) Debug

```go
func (n ASTModifier) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (ASTModifier) Span

```go
func (n ASTModifier) Span() *Span
```
Return the span for the AST node.

#### type ASTNode

```go
type ASTNode interface {
	// Return the span for this node.
	//
	// Spans can be used in diagnostics to show where in the source file
	// an issue exists.
	Span() *Span

	// Print debugging information for the given node.
	Debug(...any) *debug.Debug
}
```

Abstract Syntax Tree.

#### type ASTNot

```go
type ASTNot struct {
	Kid ASTNode // Child node.
}
```

An AST node for the `NOT' logical operator.

#### func (ASTNot) Debug

```go
func (n ASTNot) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (ASTNot) Span

```go
func (n ASTNot) Span() *Span
```
Return the span for the AST node.

#### type ASTOr

```go
type ASTOr struct {
	Kids []ASTNode // Child nodes.
}
```

An AST node for the `OR' logical operator.

#### func (ASTOr) Debug

```go
func (n ASTOr) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (ASTOr) Span

```go
func (n ASTOr) Span() *Span
```
Return the span for the AST node.

#### type ASTPredicate

```go
type ASTPredicate struct {
	Range *ASTRange // Target range value

	Comparator *ASTComparator // Comparator to use.
	Fuzz       *float64       // Levenshtein Distance.
	Boost      *float64       // Boost value.
	Field      string         // Target field.
	String     string         // Target string value.
	Regex      string         // Target regex pattern.

	Kind      PredicateKind // Predicate kind.
	Number    float64       // Target numeric value.
	Proximity int           // String promity.
}
```

An AST node for predicates.

#### func (ASTPredicate) Debug

```go
func (n ASTPredicate) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (ASTPredicate) Span

```go
func (n ASTPredicate) Span() *Span
```
Return the span for the AST node.

#### type ASTRange

```go
type ASTRange struct {
	Lo   *ASTLiteral // Start of range.
	Hi   *ASTLiteral // End of range.
	IncL bool        // Start is inclusive?
	IncH bool        // End is inclusive?
}
```

Range structure.

#### func (ASTRange) Debug

```go
func (r ASTRange) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### type BooleOp

```go
type BooleOp int
```

Boolean operation type.

```go
const (
	BooleAnd BooleOp = iota
	BooleOr
)
```

#### type ComparatorKind

```go
type ComparatorKind int
```

Comparator kind type.

```go
const (
	ComparatorLT  ComparatorKind = iota // Comparator is `LT'.
	ComparatorLTE                       // Comparator is `LTE'.
	ComparatorGT                        // Comparator is `GT'.
	ComparatorGTE                       // Comparator is `GTE'.
	ComparatorEQ                        // Comparator is `EQ'.
	ComparatorNEQ                       // Comparator is `NEQ'.
)
```

#### func  InvertComparator

```go
func InvertComparator(kind ComparatorKind) ComparatorKind
```

#### type Diagnostic

```go
type Diagnostic struct {
	Msg  string // Diagnostic message.
	At   *Span  // Location within source code.
	Hint string // Hint message, if applicable.
}
```


#### func  NewDiagnostic

```go
func NewDiagnostic(msg string, at *Span) Diagnostic
```

#### func  NewDiagnosticHint

```go
func NewDiagnosticHint(msg, hint string, at *Span) Diagnostic
```

#### func (Diagnostic) String

```go
func (d Diagnostic) String() string
```
Return the string representation of a diagnostic.

#### type Disassembler

```go
type Disassembler struct {
}
```


#### func  NewDefaultDisassembler

```go
func NewDefaultDisassembler() *Disassembler
```

#### func  NewDisassembler

```go
func NewDisassembler(opts DisassemblerOpts) *Disassembler
```

#### func (*Disassembler) Dissassemble

```go
func (d *Disassembler) Dissassemble(writer io.Writer)
```

#### func (*Disassembler) SetProgram

```go
func (d *Disassembler) SetProgram(program *Program)
```

#### type DisassemblerOpts

```go
type DisassemblerOpts struct {
	WithComments bool // Include decoded comments?
	AddrWidth    int  // Width of an address.  0 = auto.
	OpcodeWidth  int  // Pad opcode column. 0 auto.
	OperandWidth int  // Pad operand column. 0 = auto.
}
```


#### func  NewDefaultDisassemblerOpts

```go
func NewDefaultDisassemblerOpts() DisassemblerOpts
```

#### type FieldSpec

```go
type FieldSpec struct {
	Name     string    // Name of the field.
	FType    FieldType // Field type of the field.
	Analyser string    // Unused.
	Layouts  []string  // Layouts used for type parsers.
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

#### type IRAnd

```go
type IRAnd struct {
	Kids []IRNode
}
```


#### func (IRAnd) Debug

```go
func (n IRAnd) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRAnd) Emit

```go
func (n IRAnd) Emit(program *Program, tLabel, fLabel LabelID)
```
Emit opcode.

#### func (IRAnd) Key

```go
func (n IRAnd) Key() string
```
Generate key.

#### type IRAny

```go
type IRAny struct {
	Field string
}
```


#### func (IRAny) Debug

```go
func (n IRAny) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRAny) Emit

```go
func (n IRAny) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRAny) Key

```go
func (n IRAny) Key() string
```
Generate the key.

#### type IRFalse

```go
type IRFalse struct {
}
```


#### func (IRFalse) Debug

```go
func (n IRFalse) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRFalse) Emit

```go
func (n IRFalse) Emit(_ *Program, _, _ LabelID)
```
Emit opcode.

#### func (IRFalse) Key

```go
func (n IRFalse) Key() string
```
Generate key.

#### type IRGlob

```go
type IRGlob struct {
	Field string
	Glob  string
}
```


#### func (IRGlob) Debug

```go
func (n IRGlob) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRGlob) Emit

```go
func (n IRGlob) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Emit opcode.

#### func (IRGlob) Key

```go
func (n IRGlob) Key() string
```
Generate the key.

#### type IRIPCmp

```go
type IRIPCmp struct {
	Field string
	Op    ComparatorKind
	Value netip.Addr
}
```


#### func (IRIPCmp) Debug

```go
func (n IRIPCmp) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRIPCmp) Emit

```go
func (n IRIPCmp) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRIPCmp) Key

```go
func (n IRIPCmp) Key() string
```
Generate the key.

#### type IRIPRange

```go
type IRIPRange struct {
	Field string
	Lo    netip.Addr
	Hi    netip.Addr
	IncL  bool
	IncH  bool
}
```


#### func (IRIPRange) Debug

```go
func (n IRIPRange) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRIPRange) Emit

```go
func (n IRIPRange) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRIPRange) Key

```go
func (n IRIPRange) Key() string
```
Generate the key.

#### type IRNode

```go
type IRNode interface {
	// Return a unique key for the node.
	//
	// This is used during code generation.
	Key() string

	// Print debugging information.
	Debug(...any) *debug.Debug

	// Emit code.
	Emit(*Program, LabelID, LabelID)
}
```

An intermediate representation node of a syntactic element.

#### type IRNot

```go
type IRNot struct {
	Kid IRNode
}
```


#### func (IRNot) Debug

```go
func (n IRNot) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRNot) Emit

```go
func (n IRNot) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Emit opcode.

#### func (IRNot) Key

```go
func (n IRNot) Key() string
```
Generate key.

#### type IRNumberCmp

```go
type IRNumberCmp struct {
	Field string
	Op    ComparatorKind
	Value float64
}
```


#### func (IRNumberCmp) Debug

```go
func (n IRNumberCmp) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRNumberCmp) Emit

```go
func (n IRNumberCmp) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRNumberCmp) Key

```go
func (n IRNumberCmp) Key() string
```
Generate the key.

#### type IRNumberRange

```go
type IRNumberRange struct {
	Field string
	Lo    *float64
	Hi    *float64
	IncL  bool
	IncH  bool
}
```


#### func (IRNumberRange) Debug

```go
func (n IRNumberRange) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRNumberRange) Emit

```go
func (n IRNumberRange) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRNumberRange) Key

```go
func (n IRNumberRange) Key() string
```
Generate the key.

#### type IROr

```go
type IROr struct {
	Kids []IRNode
}
```


#### func (IROr) Debug

```go
func (n IROr) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IROr) Emit

```go
func (n IROr) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Emit opcode.

#### func (IROr) Key

```go
func (n IROr) Key() string
```
Generate key.

#### type IRPhrase

```go
type IRPhrase struct {
	Field     string
	Phrase    string
	Proximity int
	Fuzz      *float64
	Boost     *float64
}
```


#### func (IRPhrase) Debug

```go
func (n IRPhrase) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRPhrase) Emit

```go
func (n IRPhrase) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRPhrase) HasWildcard

```go
func (n IRPhrase) HasWildcard() bool
```

#### func (IRPhrase) Key

```go
func (n IRPhrase) Key() string
```
Generate the key.

#### type IRPrefix

```go
type IRPrefix struct {
	Field  string
	Prefix string
}
```


#### func (IRPrefix) Debug

```go
func (n IRPrefix) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRPrefix) Emit

```go
func (n IRPrefix) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRPrefix) Key

```go
func (n IRPrefix) Key() string
```
Generate a key.

#### type IRRegex

```go
type IRRegex struct {
	Field    string
	Pattern  string
	Compiled *regexp.Regexp
}
```


#### func (IRRegex) Debug

```go
func (n IRRegex) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRRegex) Emit

```go
func (n IRRegex) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRRegex) Key

```go
func (n IRRegex) Key() string
```
Generate the key.

#### type IRStringEQ

```go
type IRStringEQ struct {
	Field string
	Value string
}
```


#### func (IRStringEQ) Debug

```go
func (n IRStringEQ) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRStringEQ) Emit

```go
func (n IRStringEQ) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRStringEQ) Key

```go
func (n IRStringEQ) Key() string
```
Generate the key.

#### type IRStringNEQ

```go
type IRStringNEQ struct {
	Field string
	Value string
}
```


#### func (IRStringNEQ) Debug

```go
func (n IRStringNEQ) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRStringNEQ) Emit

```go
func (n IRStringNEQ) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRStringNEQ) Key

```go
func (n IRStringNEQ) Key() string
```
Generate the key.

#### type IRTimeCmp

```go
type IRTimeCmp struct {
	Field string
	Op    ComparatorKind
	Value int64
}
```


#### func (IRTimeCmp) Debug

```go
func (n IRTimeCmp) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRTimeCmp) Emit

```go
func (n IRTimeCmp) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRTimeCmp) Key

```go
func (n IRTimeCmp) Key() string
```
Generate the key.

#### type IRTimeRange

```go
type IRTimeRange struct {
	Field string
	Lo    *int64
	Hi    *int64
	IncL  bool
	IncH  bool
}
```


#### func (IRTimeRange) Debug

```go
func (n IRTimeRange) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRTimeRange) Emit

```go
func (n IRTimeRange) Emit(program *Program, trueLabel, falseLabel LabelID)
```
Generate opcode.

#### func (IRTimeRange) Key

```go
func (n IRTimeRange) Key() string
```
Generate the key.

#### type IRTrue

```go
type IRTrue struct {
}
```


#### func (IRTrue) Debug

```go
func (n IRTrue) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (IRTrue) Emit

```go
func (n IRTrue) Emit(_ *Program, _, _ LabelID)
```
Emit opcode.

#### func (IRTrue) Key

```go
func (n IRTrue) Key() string
```
Generate key.

#### type Instr

```go
type Instr struct {
	Op   OpCode // Instruction's opcode.
	Args []any  // Instruction's operands.
}
```


#### func (Instr) IsJump

```go
func (isn Instr) IsJump() bool
```
Is the instruction a jump of some kind?

#### func (Instr) String

```go
func (isn Instr) String() string
```
Return the string representation of an instruction.

#### type LabelID

```go
type LabelID int
```

Label identifier type.

#### type LexedToken

```go
type LexedToken struct {
	Literal Literal  // Literal value for the token.
	Lexeme  string   // Lexeme for the token.
	Start   Position // Start position within source code.
	End     Position // End position within source code.
	Token   Token    // The token.
}
```

A lexed token.

#### func  NewLexedToken

```go
func NewLexedToken(token Token, lexeme string, start, end Position) LexedToken
```
Return a new lexed token with the given lexeme.

#### func  NewLexedTokenWithError

```go
func NewLexedTokenWithError(token Token, lexeme string, err error, start, end Position) LexedToken
```
Return a new lexed token with the given lexeme and error message.

#### func  NewLexedTokenWithLiteral

```go
func NewLexedTokenWithLiteral(token Token, lexeme string, lit any, start, end Position) LexedToken
```
Return a new lexed token with the given lexeme and literal.

#### func (*LexedToken) Debug

```go
func (lt *LexedToken) Debug(params ...any) *debug.Debug
```
Display debugging information.

#### func (*LexedToken) String

```go
func (lt *LexedToken) String() string
```

#### type Lexer

```go
type Lexer interface {
	Reset()
	Tokens() []LexedToken
	Lex(io.Reader) ([]LexedToken, error)
}
```


#### func  NewLexer

```go
func NewLexer() Lexer
```
Create a new lexer.

#### type Literal

```go
type Literal struct {
	Value any   // The literal value.
	Err   error // The error to pass as a literal.
}
```

Literal value structure.

#### func  NewErrorLiteral

```go
func NewErrorLiteral(err error) Literal
```
Create a new error literal.

#### func  NewLiteral

```go
func NewLiteral(value any) Literal
```
Create a new value literal.

#### func (Literal) Error

```go
func (l Literal) Error() string
```
Return the error message if one is present.

#### func (Literal) IsError

```go
func (l Literal) IsError() bool
```
Is the literal an error?

#### func (Literal) String

```go
func (l Literal) String() string
```
Return the string representation of the literal.

#### type LiteralKind

```go
type LiteralKind int
```

Literal kind type.

```go
const (
	LString    LiteralKind = iota // Literal is a string.
	LNumber                       // Literal is a number.
	LUnbounded                    // Literal has no value bound to it.
)
```

#### type ModifierKind

```go
type ModifierKind int
```

Modifier kind type.

```go
const (
	// This modifier will mark the predicate to which it is attached as
	// having a required value.  If the value is not present, then the
	// predicate will fail.
	ModRequire ModifierKind = iota

	// This modifier will mark the predicate to which it is attached as
	// having a prohibited value.  If the value is present, then the
	// predicate will fail.
	ModProhibit
)
```

#### type NNF

```go
type NNF interface {
	NNF(IRNode) IRNode
}
```


#### func  NewNNF

```go
func NewNNF() NNF
```

#### type OpCode

```go
type OpCode int
```

Opcode.

```go
const (

	// No operation.
	OpNoOp OpCode = iota

	// Return the contents of the accumulator.
	//
	// `RET: <- ACC`
	OpReturn

	// Special instruction used by the code generator to mark a
	// program location as a label for use by the jump instructions.
	//
	// This instruction is not part of the bytecode.  It is removed
	// when the label resolver processes code for jump addresses.
	//
	// Should it be included in a bytecode instruction stream, it will
	// equate to a NOP.
	OpLabel

	// Jump to label.
	//
	// `JMP lbl: (jump)`
	OpJump

	// Jump to label if accumulator is zero.
	//
	// `JZ lbl: (jump if ACC == 0)`
	OpJumpZ

	// Jump to label if accumulator is non-zero.
	//
	// `JNZ lbl: (jump if ACC > 0)`
	OpJumpNZ

	// Negate the value in the accumulator.
	//
	// `NOT: `ACC <- !ACC`
	OpNot

	// Load an immediate into accumulator.
	//
	// `LDA imm: ACC <- imm`
	OpLoadA

	// Load a field ID into the field register.
	//
	// `LDFLD fid: `FIELD <- fid`
	OpLoadField

	// Load a value into the boost register.
	//
	// `LDBST imm: `BOOST <- imm`
	OpLoadBoost

	// Load a value into the fuzzy register.
	//
	// `LDFZY imm: `FUZZY <- imm`
	OpLoadFuzzy

	// Compare the current field to the given string constant for
	// equality.
	//
	// Stores the result in the accumulator.
	//
	// `EQ.S sIdx: ACC <- field[FIELD} == string[sIdx]`
	OpStringEQ

	// Compare the current field to the given string constant for
	// inequality.
	//
	// Stores the result in the accumulator.
	//
	// `EQ.S sIdx: ACC <- field[FIELD} != string[sIdx]`
	OpStringNEQ

	// Test whether the current field has the given string constant as
	// a prefix.
	//
	// Stores the result in the accumulator.
	//
	// `PFX.S sIdx: ACC <- HasPrefix(field[FIELD], string[sIdx])`
	OpPrefix

	// Test whether the current field matches the given glob pattern.
	//
	// Stores the result in the accumulator.
	//
	// `GLB.S sIdx: ACC <- MatchesGlob(field[FIELD], string[sIdx])`
	OpGlob

	// Perform a regular expression match of the current field against
	// the given regular expression constant.
	//
	// Stores the result in the accumulator.
	//
	// `REX.S rIdx: ACC <- MatchesRexeg(field[FIELD], regex[rIdx])`
	OpRegex

	// Test whether the current field contains the given string constant
	// as a phrase.
	//
	// If non-zero, the `prox` argument specifies maximum Levenshtein
	// distance (proximity) allowed for a match.
	//
	// Stores the result in the accumulator.
	//
	// `PHR.S sIdx: ACC <- MatchesPhrase(field[FIELD], string[sIdx])`
	OpPhrase

	// Test whether the current field has any value at all.
	//
	// Stores the result in the accumulator.
	//
	// `ANY: ACC <- HasAnyValue(field[FIELD])`
	OpAny

	// Test whether the current field has equality with the given
	// number constant.
	//
	// Stores the result in the accumulator.
	//
	// `EQ.N nIdx: ACC <- (field[FIELD] == number[nIdx])`
	OpNumberEQ

	// Test whether the current field has inequality with the given
	// number constant.
	//
	// Stores the result in the accumulator.
	//
	// `NEQ.N nIdx: ACC <- (field[FIELD] != number[nIdx])`
	OpNumberNEQ

	// Test whether the current field has a value that is lesser than
	// the given number constant.
	//
	// Stores the result in the accumulator.
	//
	// `LT.N nIdx: ACC <- (field[FIELD] < number[nIdx])
	OpNumberLT

	// Test whether the current field has a value that is lesser than or
	// equal to the given number constant.
	//
	// Stores the result in the accumulator.
	//
	// `LTE.N nIdx: ACC <- (field[FIELD] <= number[nIdx])
	OpNumberLTE

	// Test whether the current field has a value that is greater than
	// the given number constant.
	//
	// Stores the result in the accumulator.
	//
	// `GT.N nIdx: ACC <- (field[FIELD] > number[nIdx])
	OpNumberGT

	// Test whether the current field has a value that is greater than or
	// equal to the given number constant.
	//
	// Stores the result in the accumulator.
	//
	// `GTE.N nIdx: ACC <- (field[FIELD] >= number[nIdx])
	OpNumberGTE

	// Test whether the current field has a value that falls within the
	// given range.
	//
	// `loIdx` is the starting number in the range.
	// `hiIdx` is the ending number in the range.
	// `incL` is non-zero if the range is to be inclusive at the lowest.
	// `incH' is non-zero if the range is to be inclusive at the highest.
	//
	// Stores the results in the accumulator.
	//
	// RNG.N loIdx hiIdx incL incH: ACC <- inRange(field[field]...)`
	OpNumberRange

	// Test whether the current field has equality with the given
	// date/time constant.
	//
	// Stores the result in the accumulator.
	//
	// `EQ.T tIdx: ACC <- (field[FIELD] == time[tIdx])`
	OpTimeEQ

	// Test whether the current field has inequality with the given
	// date/time constant.
	//
	// Stores the result in the accumulator.
	//
	// `NEQ.T tIdx: ACC <- (field[FIELD] != time[tIdx])`
	OpTimeNEQ

	// Test whether the current field has a value that is lesser than
	// the given date/time constant.
	//
	// Stores the result in the accumulator.
	//
	// `LT.T tIdx: ACC <- (field[FIELD] < time[tIdx])
	OpTimeLT

	// Test whether the current field has a value that is lesser than or
	// equal to the given date/time constant.
	//
	// Stores the result in the accumulator.
	//
	// `LTE.T tIdx: ACC <- (field[FIELD] <= time[tIdx])
	OpTimeLTE

	// Test whether the current field has a value that is greater than
	// the given date/time constant.
	//
	// Stores the result in the accumulator.
	//
	// `GT.T nIdx: ACC <- (field[FIELD] > time[tIdx])
	OpTimeGT

	// Test whether the current field has a value that is greater than or
	// equal to the given date/time constant.
	//
	// Stores the result in the accumulator.
	//
	// `GTE.T tIdx: ACC <- (field[FIELD] >= timer[tIdx])
	OpTimeGTE

	// Test whether the current field has a value that falls within the
	// given range.
	//
	// `loIdx` is the starting date/time in the range.
	// `hiIdx` is the ending date/time in the range.
	// `incL` is non-zero if the range is to be inclusive at the lowest.
	// `incH' is non-zero if the range is to be inclusive at the highest.
	//
	// Stores the results in the accumulator.
	//
	// RNG.T loIdx hiIdx incL incH: ACC <- inRange(field[field]...)`
	OpTimeRange

	// Test whether the current field has equality with the given
	// IP address constant.
	//
	// Stores the result in the accumulator.
	//
	// `EQ.IP ipIdx: ACC <- (field[FIELD] == address[ipIdx])`
	OpIPEQ

	// Test whether the current field has inequality with the given
	// IP address constant.
	//
	// Stores the result in the accumulator.
	//
	// `NEQ.IP ipIdx: ACC <- (field[FIELD] != address[ipIdx])`
	OpIPNEQ

	// Test whether the current field has a value that is lesser than
	// the given IP address constant.
	//
	// Stores the result in the accumulator.
	//
	// `LT.IP ipIdx: ACC <- (field[FIELD] < address[ipIdx])
	OpIPLT

	// Test whether the current field has a value that is lesser than or
	// equal to the given IP address constant.
	//
	// Stores the result in the accumulator.
	//
	// `LTE.IP ipIdx: ACC <- (field[FIELD] <= address[ipIdx])
	OpIPLTE

	// Test whether the current field has a value that is greater than
	// the given IP address constant.
	//
	// Stores the result in the accumulator.
	//
	// `GT.IP ipIdx: ACC <- (field[FIELD] > address[ipIdx])
	OpIPGT

	// Test whether the current field has a value that is greater than or
	// equal to the given IP address constant.
	//
	// Stores the result in the accumulator.
	//
	// `GTE.IP ipIdx: ACC <- (field[FIELD] >= address[ipIdx])
	OpIPGTE

	// Test whether the current field has a value that falls within the
	// given range.
	//
	// `loIdx` is the starting IP address in the range.
	// `hiIdx` is the ending IP address in the range.
	// `incL` is non-zero if the range is to be inclusive at the lowest.
	// `incH' is non-zero if the range is to be inclusive at the highest.
	//
	// Stores the results in the accumulator.
	//
	// RNG.IP loIdx hiIdx incL incH: ACC <- inRange(field[field]...)`
	OpIPRange

	// Test whether the current field is within a CIDR range.
	//
	// `IN.CIDR ipIdx, prefix: ACC <- (field[FIELD] = cidr[ipIdx,prefix])`
	OpInCIDR

	// Maximum number of opcode currently supported.
	OpMaximum
)
```

#### func  GetIPComparator

```go
func GetIPComparator(cmp ComparatorKind, def OpCode) OpCode
```
Return the relevant IP address opcode for the given comparator.

If no opcode is found, then `def` will be used instead.

#### func  GetNumberComparator

```go
func GetNumberComparator(cmp ComparatorKind, def OpCode) OpCode
```
Return the relevant numeric opcode for the given comparator.

If no opcode is found, then `def` will be used instead.

#### func  GetTimeComparator

```go
func GetTimeComparator(cmp ComparatorKind, def OpCode) OpCode
```
Return the relevant date/time opcode for the given comparator.

If no opcode is found, then `def` will be used instead.

#### type Parser

```go
type Parser interface {
	// Reset the parser state.
	Reset()

	// Parse the list of lexed tokens and generate an AST.
	Parse([]LexedToken) (ASTNode, []Diagnostic)

	// Return a list of diagnostic messages.
	Diagnostics() []Diagnostic
}
```


#### func  NewParser

```go
func NewParser() Parser
```
Create a new parser instance.

#### type Position

```go
type Position struct {
	Line   int // Line number.
	Column int // Column number.
}
```

Position within source code.

#### func  NewEmptyPosition

```go
func NewEmptyPosition() Position
```
Create a new empty position.

#### func  NewPosition

```go
func NewPosition(line, col int) Position
```
Create a new position with the given line and column numbers.

#### func (Position) String

```go
func (p Position) String() string
```
Return the string representation of a position.

#### type PredicateKind

```go
type PredicateKind int
```

Predicate kind type.

```go
const (
	PredicateCMP    PredicateKind = iota // Predicate is a comparator.
	PredicateEQS                         // Predicate is `EQ.S'.
	PredicateANY                         // Predicate is `ANY'.
	PredicateGLOB                        // Predicate is `GLOB'.
	PredicateNEQS                        // Predicate is `NEQ.S'.
	PredicatePHRASE                      // Predicate is `PHRASE'.
	PredicatePREFIX                      // Predicate is `PREFIX'.
	PredicateRANGE                       // Predicate is `RANGE'.
	PredicateREGEX                       // Predicate is `REGEX'.
)
```

#### type Program

```go
type Program struct {
	Fields   []string         // Field constants.
	Strings  []string         // String constants.
	Numbers  []float64        // Number constants.
	Times    []int64          // Date/time constants.
	IPs      []netip.Addr     // IP address constants.
	Patterns []*regexp.Regexp // Regular expression constants.
	Code     []Instr          // Bytecode.
}
```


#### func  NewProgram

```go
func NewProgram() *Program
```
Create a new program instance.

#### func (*Program) AddFieldConstant

```go
func (p *Program) AddFieldConstant(val string) int
```
Add a field name constant.

If the given constant value exists, then an index to its array position is
returned.

#### func (*Program) AddIPConstant

```go
func (p *Program) AddIPConstant(val netip.Addr) int
```
Add an IP address constant.

If the given constant value exists, then an index to its array position is
returned.

#### func (*Program) AddNumberConstant

```go
func (p *Program) AddNumberConstant(val float64) int
```
Add a numeric constant.

If the given constant value exists, then an index to its array position is
returned.

#### func (*Program) AddRegexConstant

```go
func (p *Program) AddRegexConstant(val *regexp.Regexp) int
```
Add a regular expression constant.

If the given constant value exists, then an index to its array position is
returned.

#### func (*Program) AddStringConstant

```go
func (p *Program) AddStringConstant(val string) int
```
Add a string constant.

If the given constant value exists, then an index to its array position is
returned.

#### func (*Program) AddTimeConstant

```go
func (p *Program) AddTimeConstant(val int64) int
```
Add a date/time constant.

If the given constant value exists, then an index to its array position is
returned.

#### func (*Program) AppendIsn

```go
func (p *Program) AppendIsn(opCode OpCode, args ...any)
```
Append an instruction to the bytecode.

Any provided arguments will be added as the instruction's operands.

#### func (*Program) AppendJump

```go
func (p *Program) AppendJump(opCode OpCode, target LabelID)
```
Append a jump instruction to the bytecode.

The instruction's operand will be the target label ID.

#### func (*Program) BindLabel

```go
func (p *Program) BindLabel(id LabelID)
```
Append a `LABEL` instruction bound to the given label identifier.

The instruction will be removed by the label resolver.

#### func (*Program) Emit

```go
func (p *Program) Emit(irNode IRNode)
```

#### func (*Program) NewLabel

```go
func (p *Program) NewLabel() LabelID
```
Generate a new label identifier and return it.

#### func (*Program) Peephole

```go
func (p *Program) Peephole()
```

#### type Schema

```go
type Schema map[string]FieldSpec
```


#### type Simplifier

```go
type Simplifier interface {
	Simplify(IRNode) IRNode
}
```


#### func  NewSimplifier

```go
func NewSimplifier() Simplifier
```

#### type Span

```go
type Span struct {
}
```

Span within source code.

#### func  NewEmptySpan

```go
func NewEmptySpan() *Span
```
Create a new empty span.

#### func  NewSpan

```go
func NewSpan(start, end Position) *Span
```
Create a new span with the given start and end positions.

#### func (*Span) End

```go
func (s *Span) End() Position
```
Return the span's ending position.

#### func (*Span) Start

```go
func (s *Span) Start() Position
```
Return the span's starting position.

#### func (*Span) String

```go
func (s *Span) String() string
```
Return the string representation of a span.

#### type Token

```go
type Token int
```

Lucette token type.

```go
const (
	TokenEOF Token = iota // End of file.

	TokenNumber // Numeric value.
	TokenPhrase // String phrase value.
	TokenField  // Field name.
	TokenRegex  // Regular expression.

	TokenPlus     // '+'
	TokenMinus    // '-'
	TokenStar     // '*'
	TokenQuestion // '?'
	TokenLParen   // '('
	TokenLBracket // '['
	TokenLCurly   // '{'
	TokenRParen   // ')'
	TokenRBracket // ']'
	TokenRCurly   // '}'
	TokenColon    // ':'
	TokenTilde    // '~'
	TokenCaret    // '^'

	TokenTo  // 'TO'.
	TokenAnd // 'AND'/'&&'.
	TokenOr  // 'OR'/'||'.
	TokenNot // 'NOT'/'!'.
	TokenLT  // '<'.
	TokenLTE // '<='.
	TokenGT  // '>'
	TokenGTE // '>='

	TokenIllegal // Illegal token.
	TokenUnknown // Unknown token.
)
```

#### func (Token) Literal

```go
func (t Token) Literal() string
```
Return the literal string representation of a token if it has one.

#### func (Token) String

```go
func (t Token) String() string
```
Return the string representation of a token.

#### type Typer

```go
type Typer interface {
	// Generate typed IR from the given AST root node.
	Type(ASTNode) (IRNode, []Diagnostic)

	// Return the diagnostic messages generated during IR generation.
	Diagnostics() []Diagnostic
}
```


#### func  NewTyper

```go
func NewTyper(sch Schema) Typer
```
