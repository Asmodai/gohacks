// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// parser.go --- Shunting-yard FSM parser.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * Comments:

// * Package:

package lucette

// * Imports:

import (
	"fmt"
	"regexp"
	"strings"
)

// * Constants:

const (
	sExpectPrimary parserState = iota
	sAfterPrimary

	irOR irKind = iota
	irAND
	irImplicitAND
	irNOT
	irREQ
	irPROH
	irFIELDAPPLY
	irLPAREN

	precOR  bindPrec = 10
	precAND bindPrec = 20
	precNOT bindPrec = 30
	precFAP bindPrec = 85
)

// * Variables:

var (
	//nolint:gochecknoglobals
	primaryStarters = map[TokenType]struct{}{
		ttPhrase:   struct{}{},
		ttNumber:   struct{}{},
		ttRegexp:   struct{}{},
		ttField:    struct{}{},
		ttNot:      struct{}{},
		ttLParen:   struct{}{},
		ttLBracket: struct{}{},
		ttLCurly:   struct{}{},
		ttLT:       struct{}{},
		ttLTE:      struct{}{},
		ttGT:       struct{}{},
		ttGTE:      struct{}{},
	}

	//nolint:gochecknoglobals,unused
	primaryEnders = map[TokenType]struct{}{
		ttPhrase:   struct{}{},
		ttNumber:   struct{}{},
		ttRegexp:   struct{}{},
		ttRParen:   struct{}{},
		ttRBracket: struct{}{},
		ttRCurly:   struct{}{},
	}

	//nolint:gochecknoglobals
	precedences = map[irKind]bindPrec{
		irOR:          precOR,
		irAND:         precAND,
		irImplicitAND: precAND,
		irNOT:         precNOT,
		irREQ:         precNOT,
		irPROH:        precNOT,
		irFIELDAPPLY:  precFAP,
	}

	//nolint:gochecknoglobals
	comparators = map[TokenType]CmpKind{
		ttLT:  CmpLT,
		ttLTE: CmpLTE,
		ttGT:  CmpGT,
		ttGTE: CmpGTE,
	}
)

// * Code:

// ** Types:

// Parser state.
type parserState int

// Operand kind.
type irKind int

// Binding precedence.
type bindPrec int

// Operand entry.
type opEntry struct {
	kind irKind
	aux  any
	span Span
}

// ** Diagnostics:

// Diagnostic message.
type Diagnostic struct {
	Msg  string // Diagnostic message.
	At   Span   // Location within token stream.
	Hint string // Hint message, if applicable.
}

// Pretty-print a diagnostic to a string.
func (d *Diagnostic) String() string {
	var sbld strings.Builder

	sbld.WriteString(d.At.String())
	sbld.WriteString(": ")
	sbld.WriteString(d.Msg)

	if len(d.Hint) > 0 {
		sbld.WriteString(" [")
		sbld.WriteString(d.Hint)
		sbld.WriteRune(']')
	}

	return sbld.String()
}

// ** Parser structure:

// Parser structure.
type Parser struct {
	toks  []Token
	index int
	out   []Node
	ops   []opEntry
	state parserState
	diags []Diagnostic
}

// ** Reader Methods:

// Generate a diagnostic error.
func (p *Parser) errorf(span Span, format string, args ...any) {
	p.diags = append(p.diags,
		Diagnostic{
			Msg: fmt.Sprintf(format, args...),
			At:  span})
}

// Return the current token without advancing the token reader.
func (p *Parser) peek() Token {
	return p.toks[p.index]
}

// Get the next token.
func (p *Parser) next() Token {
	tok := p.toks[p.index]
	p.index++

	return tok
}

// Unread a token back to the reader.
func (p *Parser) unread() {
	if p.index > 0 {
		p.index--
	}
}

// Advance the token reader if the current token is of the given type.
func (p *Parser) accept(tType TokenType) bool {
	if p.peek().tokenType == tType {
		p.index++

		return true
	}

	return false
}

// Expect the current token to be of the given type.
//
// If it is not, then generate a diagnostic message.
func (p *Parser) expect(tType TokenType, msg string) Token {
	tok := p.next()

	if tok.tokenType != tType {
		p.errorf(spanTok(tok), "expected %s", msg)

		return newToken(tType, "", tok.start, tok.end)
	}

	return tok
}

// Return the previous token type.
//
//nolint:unused
func (p *Parser) prevType() TokenType {
	if p.index == 0 {
		return ttIllegal
	}

	return p.toks[p.index-1].tokenType
}

// ** Stack methods:

// Push an operator to the stack.
func (p *Parser) pushOp(kind irKind, aux any, span Span) {
	p.ops = append(
		p.ops,
		opEntry{kind: kind, aux: aux, span: span})
}

// Pop an operator from the stack.
func (p *Parser) popOp() (opEntry, bool) {
	if len(p.ops) == 0 {
		return opEntry{}, false
	}

	elt := p.ops[len(p.ops)-1]
	p.ops = p.ops[:len(p.ops)-1]

	return elt, true
}

// Return the operator on the top of the stack.
func (p *Parser) topOp() (opEntry, bool) {
	if len(p.ops) == 0 {
		return opEntry{}, false
	}

	return p.ops[len(p.ops)-1], true
}

// Push a node to the AST.
func (p *Parser) pushNode(node Node) {
	p.out = append(p.out, node)
}

// Pop a node from the AST.
func (p *Parser) popNode() Node {
	if len(p.out) == 0 {
		p.errorf(Span{}, "stack underflow")

		return &NodePred{kind: PK_EXISTS}
	}

	end := len(p.out) - 1
	node := p.out[end]
	p.out = p.out[:end]

	return node
}

// ** Reduction:

// Reduce an operator from the stack.
//
//nolint:forcetypeassert
func (p *Parser) reduceOne() {
	operand, ok := p.popOp()
	if !ok {
		return
	}

	switch operand.kind {
	case irOR:
		right, left := p.popNode(), p.popNode()
		p.pushNode(mkOr(left, right))

	case irAND, irImplicitAND:
		right, left := p.popNode(), p.popNode()
		p.pushNode(mkAnd(left, right))

	case irNOT:
		kid := p.popNode()
		p.pushNode(&NodeNot{
			kid:  kid,
			span: Span{operand.span.start, kid.Span().end}})

	case irREQ:
		kid := p.popNode()
		p.pushNode(&NodeMod{
			kind: ModRequire,
			kid:  kid,
			span: Span{operand.span.start, kid.Span().end}})

	case irPROH:
		kid := p.popNode()
		p.pushNode(&NodeMod{
			kind: ModProhibit,
			kid:  kid,
			span: Span{operand.span.start, kid.Span().end}})

	case irFIELDAPPLY:
		fld := operand.aux.(string)
		kid := p.popNode()
		p.pushNode(applyField(
			fld,
			kid,
			Span{operand.span.start, kid.Span().end}))

	case irLPAREN: // Barrier.
	}
}

// Reduce operators from the top of the stack until the given minimum
// binding precedence is met.
func (p *Parser) reduceWhile(minPrec bindPrec) {
	for {
		top, topOk := p.topOp()
		if !topOk {
			return
		}

		if top.kind == irLPAREN {
			return
		}

		prec := precedence(top.kind)
		if prec < minPrec {
			return
		}

		p.reduceOne()
	}
}

// ** Primary builders:

func (p *Parser) makeRegex(tok Token) Node {
	span := spanTok(tok)

	pattern, ok := tok.literal.value.(string)
	if !ok {
		p.errorf(span,
			"invalid regular expression value: %v",
			tok.literal.value)

		return &NodePred{kind: PK_EXISTS, span: span}
	}

	compiled, err := regexp.Compile(pattern)
	if err != nil {
		p.errorf(span,
			"regular expression compile failed: %q",
			err.Error())

		return &NodePred{kind: PK_EXISTS, span: span}
	}

	//nolint:forcetypeassert
	return &NodePred{
		kind:  PK_REGEX,
		reval: tok.literal.value.(string),
		repat: compiled,
		span:  span}
}

// Make a predicate form from the given token.
//
//nolint:forcetypeassert
func (p *Parser) makePredForm(tok Token) Node {
	span := spanTok(tok)

	//nolint:exhaustive
	switch tok.tokenType {
	case ttPhrase:
		return &NodePred{
			kind:   PK_PHRASE,
			strval: tok.literal.value.(string),
			span:   span}

	case ttNumber:
		return &NodePred{
			kind: PK_CMP,
			cmp: &Comparator{
				Op: CmpEQ,
				Atom: NodeLit{
					kind:   LNumber,
					numval: tok.literal.value.(float64),
					span:   span}}}

	case ttRegexp:
		return p.makeRegex(tok)

	default:
		p.errorf(span, "unexpected primary")

		return &NodePred{kind: PK_EXISTS, span: span}
	}
}

// Expect that the next token in the reader is one that can be used within
// a range.
//
//nolint:forcetypeassert
func (p *Parser) expectRangeAtom() NodeLit {
	tok := p.next()
	span := spanTok(tok)

	//nolint:exhaustive
	switch tok.tokenType {
	case ttPhrase:
		return NodeLit{
			kind:   LString,
			strval: tok.literal.value.(string),
			span:   span}

	case ttNumber:
		return NodeLit{
			kind:   LNumber,
			numval: tok.literal.value.(float64),
			span:   span}

	default:
		p.errorf(span, "expected range atom (string or number)")

		return NodeLit{
			kind:   LString,
			strval: "",
			span:   span}
	}
}

// Read in a range bound from the token stream.
func (p *Parser) readRangeBound() *NodeLit {
	if p.accept(ttStar) {
		return &NodeLit{kind: LUnbounded}
	}

	atom := p.expectRangeAtom()

	return &atom
}

// Parse a RANGE or a comparator.
func (p *Parser) parseRangeOrCmp(first Token) Node {
	span := spanTok(first)

	//nolint:exhaustive
	switch first.tokenType {
	case ttLT, ttLTE, ttGT, ttGTE:
		atom := p.expectRangeAtom()
		op := comparators[first.tokenType]

		return &NodePred{
			kind: PK_CMP,
			cmp: &Comparator{
				Op:   op,
				Atom: atom},
			span: Span{span.start, atom.Span().end}}
	}

	low := p.readRangeBound()
	p.expect(ttTo, "TO")
	high := p.readRangeBound()
	end := p.next()

	if end.tokenType != ttRBracket && end.tokenType != ttRCurly {
		p.errorf(spanTok(end), "unterminated interval")
	}

	incL := (first.tokenType == ttLBracket)
	incH := (end.tokenType == ttRBracket)

	return &NodePred{
		kind: PK_RANGE,
		rnge: &Range{
			Low:  low,
			High: high,
			IncL: incL,
			IncH: incH},
		span: Span{span.start, spanTok(end).end}}
}

// ** Main loop:

// Handle any required postfix translation.
//
//nolint:forcetypeassert,exhaustive
func (p *Parser) handlePostFix() {
	for {
		switch p.peek().tokenType {
		case ttCaret:
			p.next()

			numTok := p.expect(ttNumber, "number after '^'")
			top := p.popNode()

			p.pushNode(attachBoost(top, numTok.literal.value.(float64)))

		case ttTilde:
			p.next()

			var flt *float64

			if p.peek().tokenType == ttNumber {
				val := p.next().literal.value.(float64)
				flt = &val
			}

			top := p.popNode()
			p.pushNode(attachFuzz(top, flt))

		default:
			return
		}
	}
}

// Parse a primary expression.
//
//nolint:forcetypeassert,exhaustive
func (p *Parser) parsePrimary(tok Token) {
	switch tok.tokenType {
	case ttNot:
		p.reduceWhile(precNOT)
		p.pushOp(irNOT, nil, spanTok(tok))

	case ttPlus:
		p.reduceWhile(precNOT)
		p.pushOp(irREQ, nil, spanTok(tok))

	case ttMinus:
		p.reduceWhile(precNOT)
		p.pushOp(irPROH, nil, spanTok(tok))

	case ttField:
		p.expect(ttColon, "':' after field")
		p.reduceWhile(precFAP)

		fld := tok.literal.value.(string)
		p.pushOp(irFIELDAPPLY, fld, spanTok(tok))

	case ttLParen:
		p.pushOp(irLPAREN, nil, spanTok(tok))

	case ttPhrase, ttNumber, ttRegexp:
		p.pushNode(p.makePredForm(tok))
		p.handlePostFix()

		p.state = sAfterPrimary

	case ttLBracket, ttLCurly, ttLT, ttLTE, ttGT, ttGTE:
		p.pushNode(p.parseRangeOrCmp(tok))
		p.handlePostFix()

		p.state = sAfterPrimary

	default:
		p.errorf(spanTok(tok), "expected primary")
	}
}

// Parse an expression that comes after a primary.
//
//nolint:cyclop,forcetypeassert,exhaustive
func (p *Parser) parseAfterPrimary(tok Token) {
	switch tok.tokenType {
	case ttCaret:
		numTok := p.expect(ttNumber, "number after '^'")
		top := p.popNode()

		p.pushNode(attachBoost(top, numTok.literal.value.(float64)))

	case ttTilde:
		var flt *float64

		if p.peek().tokenType == ttNumber {
			val := p.next().literal.value.(float64)
			flt = &val
		}

		top := p.popNode()
		p.pushNode(attachFuzz(top, flt))

	case ttAnd:
		p.reduceWhile(precAND)
		p.pushOp(irAND, nil, spanTok(tok))
		p.state = sExpectPrimary

	case ttOr:
		p.reduceWhile(precOR)
		p.pushOp(irOR, nil, spanTok(tok))
		p.state = sExpectPrimary

	case ttRParen:
		for {
			top, ok := p.topOp()
			if !ok {
				p.errorf(spanTok(tok), "unmatched ')'")

				break
			}

			if top.kind == irLPAREN {
				break
			}

			p.reduceOne()
		}

		if top, ok := p.popOp(); !ok || top.kind != irLPAREN {
			p.errorf(spanTok(tok), "unmatched ')'")
		}

	default:
		if startsPrimary(tok.tokenType) {
			p.unread()
			p.reduceWhile(precAND)
			p.pushOp(irImplicitAND, nil, Span{})
			p.state = sExpectPrimary
		} else {
			p.errorf(spanTok(tok), "expected operator")
		}
	}
}

// Parse the token stream.
func (p *Parser) Parse() (Node, []Diagnostic) {
	for {
		tok := p.next()

		// If we're at EOF, then go drain.
		if tok.tokenType == ttEOF {
			goto DRAIN
		}

		switch p.state {
		case sExpectPrimary:
			p.parsePrimary(tok)

		case sAfterPrimary:
			p.parseAfterPrimary(tok)
		}
	}

DRAIN:
	for {
		if len(p.ops) == 0 {
			break
		}

		top, _ := p.topOp()
		if top.kind == irLPAREN {
			p.errorf(top.span, "unmatched '('")

			_, _ = p.popOp()

			continue
		}

		p.reduceOne()
	}

	if len(p.out) == 0 {
		return &NodePred{kind: PK_EXISTS, span: Span{}}, p.diags
	}

	if len(p.out) != 1 {
		p.errorf(p.out[0].Span(), "internal: multiple roots after reduce")
	}

	return p.out[0], p.diags
}

// ** Functions:

// Return the precedence of the operator kind.
func precedence(kind irKind) bindPrec {
	prec, found := precedences[kind]

	if !found {
		return 0
	}

	return prec
}

// Generate a span for the given token.
func spanTok(tok Token) Span {
	return Span{start: tok.start, end: tok.end}
}

// Does the given token type start a primary expression?
func startsPrimary(tType TokenType) bool {
	_, found := primaryStarters[tType]

	return found
}

// Does the given token type end a primary expression?
//
//nolint:unused
func endsPrimary(tType TokenType) bool {
	_, found := primaryEnders[tType]

	return found
}

// Reduce an AND expression.
func mkAnd(lhs, rhs Node) Node {
	splhs := lhs.Span()
	sprhs := rhs.Span()

	if andlhs, ok := lhs.(*NodeAnd); ok {
		if andrhs, ok := rhs.(*NodeAnd); ok {
			return &NodeAnd{
				kids: append(andlhs.kids, andrhs.kids...),
				span: Span{splhs.start, sprhs.end}}
		}

		return &NodeAnd{
			kids: append(andlhs.kids, rhs),
			span: Span{splhs.start, sprhs.end}}
	}

	if andrhs, ok := rhs.(*NodeAnd); ok {
		return &NodeAnd{
			kids: append([]Node{lhs}, andrhs.kids...),
			span: Span{splhs.start, sprhs.end}}
	}

	return &NodeAnd{
		kids: []Node{lhs, rhs},
		span: Span{splhs.start, sprhs.end}}
}

// Reduce an OR expression.
func mkOr(lhs, rhs Node) Node {
	slhs := lhs.Span()
	srhs := rhs.Span()

	if orlhs, ok := lhs.(*NodeOr); ok {
		if orrhs, ok := rhs.(*NodeOr); ok {
			return &NodeOr{
				kids: append(orlhs.kids, orrhs.kids...),
				span: Span{slhs.start, srhs.end}}
		}

		return &NodeOr{
			kids: append(orlhs.kids, rhs),
			span: Span{slhs.start, srhs.end}}
	}

	if orrhs, ok := rhs.(*NodeOr); ok {
		return &NodeOr{
			kids: append([]Node{lhs}, orrhs.kids...),
			span: Span{slhs.start, srhs.end}}
	}

	return &NodeOr{
		kids: []Node{lhs, rhs},
		span: Span{slhs.start, srhs.end}}
}

// Apply the given field to the given node.
//
// Applies the field to the child nodes where applicable.
func applyField(field string, node Node, span Span) Node {
	switch val := node.(type) {
	case *NodePred:
		res := *val
		res.field = field
		res.span = span

		return &res

	case *NodeAnd:
		kids := make([]Node, len(val.kids))

		for idx, elt := range val.kids {
			kids[idx] = applyField(
				field,
				elt,
				Span{val.span.start, elt.Span().end})
		}

		return &NodeAnd{kids: kids, span: span}

	case *NodeOr:
		kids := make([]Node, len(val.kids))

		for idx, elt := range val.kids {
			kids[idx] = applyField(
				field,
				elt,
				Span{val.span.start, elt.Span().end})
		}

		return &NodeOr{kids: kids, span: span}

	case *NodeNot:
		return &NodeNot{
			kid:  applyField(field, val.kid, val.span),
			span: span}

	case *NodeMod:
		return &NodeMod{
			kind: val.kind,
			kid:  applyField(field, val.kid, val.span),
			span: span}

	default:
		return node
	}
}

// Attach a boost to the node.
func attachBoost(node Node, val float64) Node {
	if pred, ok := node.(*NodePred); ok {
		res := *pred

		res.boost = &val

		return &res
	}

	return node
}

// Attach a fuzz to the node.
func attachFuzz(node Node, val *float64) Node {
	if pred, ok := node.(*NodePred); ok {
		res := *pred

		res.fuzz = val

		return &res
	}

	return node
}

// Create a new parser instance.
func NewParser(toks []Token) *Parser {
	return &Parser{
		toks:  toks,
		state: sExpectPrimary,
	}
}

// * parser.go ends here.
