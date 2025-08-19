// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// codegen.go --- Code generation.
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
	"net/netip"
	"regexp"

	"github.com/Asmodai/gohacks/utils"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

//nolint:revive,stylecheck
const (
	opNOP OpCode = iota // No operation.
	//
	// Accumulator operations.
	opRET // Return ACC.
	opLDA // LDA imm: `ACC <- imm`.
	//
	// Flow control.
	opLABEL // LABEL:  Fake instruction for label binding.
	opJMP   // JMP lbl: Jump.
	opJZ    // JZ  lbl: Jump if `ACC == false`.
	opJNZ   // JNZ lbl: Jump if `ACC == true`.
	//
	// Boolean operations.
	opAND // AND: `ACC <- (pop) && ACC`.
	opOR  // OR:  `ACC <- (pop) || ACC`.
	opNOT // NOT: `ACC <- !ACC`.
	//
	// Field context.
	opSET_F // SET.F fid: `CF <- fid`.
	//
	// String/keyword Predicates.
	opEQ_S   // EQ.S   sidx:      `ACC <- (field == S[sidx])`.
	opNEQ_S  // NEQ.S  sidx:      `ACC <- (field != S[sidx])`.
	opPREFIX // PREFIX sidx:      `ACC <- (hasPrefix(field, S[sidx]))`
	opGLOB   // GLOB   sidx:      `ACC <- (matchesGlob(field, S[sidx]))`
	opREGEX  // REGEX  ridx:      `ACC <- (matchesRegex(field, R[ridx]))`.
	opPHRASE // PHRASE sidx,prox: `ACC <- (phrase match; prox == 0 exact phrase, else window-k)`.
	opEXISTS // EXISTS:           `ACC <- (hasAnyValue(field))`.
	//
	// Numeric predicates.
	opEQ_N    // EQ.N    nidx: `ACC <- (field == N[nidx])`.
	opNEQ_N   // NEQ.N   nidx: `ACC <- (field != N[nidx])`.
	opLT_N    // LT.N    nidx: `ACC <- (field < N[nidx])`.
	opLTE_N   // LTE.N   nidx: `ACC <- (field <= N[nidx])`.
	opGT_N    // GT.N    nidx: `ACC <- (field > N[nidx])`.
	opGTE_N   // GTE.N   nidx: `ACC <- (field >= N[nidx])`.
	opRANGE_N // RANGE.N loIdx,hiIdx,incL,incH
	//
	// Datetime predicates.
	opEQ_T    // EQ.T    tidx: `ACC <- (field == T[tidx])`.
	opNEQ_T   // NEQ.T   tidx: `ACC <- (field != T[tidx])`.
	opLT_T    // LT.T    tidx: `ACC <- (field < T[tidx])`.
	opLTE_T   // LTE.T   tidx: `ACC <- (field <= T[tidx])`.
	opGT_T    // GT.T    tidx: `ACC <- (field > T[tidx])`.
	opGTE_T   // GTE.T   tidx: `ACC <- (field >= T[tidx])`.
	opRANGE_T // RANGE.T loIdx,hiIdx,incL,incH
	//
	// IP predicates.
	opEQ_IP    // EQ.IP    ipidx; `ACC <- (field == IP[ipidx])`.
	opNEQ_IP   // NEQ.IP   ipidx; `ACC <- (field != IP[ipidx])`.
	opLT_IP    // LT.IP    ipidx: `ACC <- (field < IP[tidx])`.
	opLTE_IP   // LTE.IP   ipidx: `ACC <- (field <= IP[tidx])`.
	opGT_IP    // GT.IP    ipidx: `ACC <- (field > IP[tidx])`.
	opGTE_IP   // GTE.IP   ipidx: `ACC <- (field >= IP[tidx])`.
	opRANGE_IP // RANGE.IP loIdx,hiIdx,incL,incH
	opIN_CIDR  // IN.CIDR  ipidx,prefixLen
	//
	// Group/metadata.
	opBOOST // BOOST f64: annotate last leaf with boost.
	opFUZZ  // FUZZ  f64: annotate last leaf with fuzz param.

	opMaximum //nolint:unused

	defaultIsnPadding     = 10
	initialLabelArraySize = 16
)

// * Variables:

var (
	//nolint:gochecknoglobals
	opNames = map[OpCode]string{
		opNOP:      "NOP",
		opLDA:      "LDA",
		opRET:      "RET",
		opLABEL:    "LABEL",
		opJMP:      "JMP",
		opJZ:       "JZ",
		opJNZ:      "JNZ",
		opAND:      "AND",
		opOR:       "OR",
		opNOT:      "NOT",
		opSET_F:    "SET.F",
		opEQ_S:     "EQ.S",
		opNEQ_S:    "NEQ.S",
		opPREFIX:   "PREFIX",
		opGLOB:     "GLOB",
		opREGEX:    "REGEX",
		opPHRASE:   "PHRASE",
		opEXISTS:   "EXISTS",
		opEQ_N:     "EQ.N",
		opNEQ_N:    "NEQ.N",
		opLT_N:     "LT.N",
		opLTE_N:    "LTE.N",
		opGT_N:     "GT.N",
		opGTE_N:    "GTE.N",
		opRANGE_N:  "RANGE.N",
		opEQ_T:     "EQ.T",
		opNEQ_T:    "NEQ.T",
		opLT_T:     "LT.T",
		opLTE_T:    "LTE.T",
		opGT_T:     "GT.T",
		opGTE_T:    "GTE.T",
		opRANGE_T:  "RANGE.T",
		opEQ_IP:    "EQ.IP",
		opNEQ_IP:   "NEQ.IP",
		opLT_IP:    "LT.IP",
		opLTE_IP:   "LTE.IP",
		opGT_IP:    "GT.IP",
		opGTE_IP:   "GTE.IP",
		opRANGE_IP: "RANGE.IP",
		opIN_CIDR:  "IN.CIDR",
		opBOOST:    "BOOST",
		opFUZZ:     "FUZZ",
	}

	//nolint:gochecknoglobals
	cmpToOpN = map[CmpKind]OpCode{
		CmpEQ:  opEQ_N,
		CmpNEQ: opNEQ_N,
		CmpLT:  opLT_N,
		CmpLTE: opLTE_N,
		CmpGT:  opGT_N,
		CmpGTE: opGTE_N,
	}

	//nolint:gochecknoglobals
	cmpToOpT = map[CmpKind]OpCode{
		CmpEQ:  opEQ_T,
		CmpNEQ: opNEQ_T,
		CmpLT:  opLT_T,
		CmpLTE: opLTE_T,
		CmpGT:  opGT_T,
		CmpGTE: opGTE_T,
	}

	//nolint:gochecknoglobals
	cmpToOpIP = map[CmpKind]OpCode{
		CmpEQ:  opEQ_IP,
		CmpNEQ: opNEQ_IP,
		CmpLT:  opLT_IP,
		CmpLTE: opLTE_IP,
		CmpGT:  opGT_IP,
		CmpGTE: opGTE_IP,
	}

	ErrLabelMissingID = errors.Base("LABEL missing id")
	ErrLabelBadIDType = errors.Base("LABEL has bad id type")
	ErrJumpMissingArg = errors.Base("jump missing target arg")
	ErrJumpNotLabelID = errors.Base("jump target arg not LabelID")
	ErrUnboundLabel   = errors.Base("unbound label")
)

// * Code:

// ** Instruction type:

type OpCode int

type Instr struct {
	Op   OpCode
	Args []any
}

func (isn Instr) String() string {
	str, found := opNames[isn.Op]
	if !found {
		return "<invalid>"
	}

	return fmt.Sprintf("%s%v",
		utils.Pad(str, defaultIsnPadding),
		isn.Args)
}

// ** Label type:

type LabelID int

// ** Program type:

type Program struct {
	Fields      []string
	Strings     []string
	Numbers     []float64
	Times       []int64
	IPs         []netip.Addr
	Patterns    []*regexp.Regexp
	Code        []Instr
	nextLabelID LabelID
}

// ** Label methods:

func (p *Program) newLabel() LabelID {
	id := p.nextLabelID
	p.nextLabelID++

	return id
}

func (p *Program) bindLabel(id LabelID) {
	p.appendIsn(opLABEL, id)
}

func (p *Program) appendJump(op OpCode, target LabelID) {
	p.appendIsn(op, target)
}

//nolint:cyclop,funlen
func (p *Program) resolveLabels() error {
	// Phase 1: Build `label -> PC` map.
	labelPC := make(map[LabelID]int, initialLabelArraySize)
	pcOut := 0

	for idx := range p.Code {
		isn := p.Code[idx]

		if isn.Op == opLABEL {
			// Bind label to the next real instruction.
			if len(isn.Args) != 1 {
				return errors.WithMessagef(
					ErrLabelMissingID,
					"%d",
					idx)
			}

			lid, ok := isn.Args[0].(LabelID)
			if !ok {
				return errors.WithMessagef(
					ErrLabelBadIDType,
					"%d",
					idx)
			}

			labelPC[lid] = pcOut

			continue
		}

		pcOut++
	}

	// Safety: ensure labels don't point past end -- i.e., ensure there's
	// a real instruction after each label.
	//
	// If you *intentionally* bind a label at end, make sure a RET
	// follows before peephole.

	// Phase 2: Build stripped code and patch jumps.
	stripped := make([]Instr, 0, pcOut)

	for idx := range len(p.Code) {
		isn := p.Code[idx]

		if isn.Op == opLABEL {
			continue
		}

		// Patch jump.
		if isn.Op == opJMP || isn.Op == opJZ || isn.Op == opJNZ {
			if len(isn.Args) != 1 {
				return errors.WithMessagef(
					ErrJumpMissingArg,
					"%d",
					idx)
			}

			lid, found := isn.Args[0].(LabelID)
			if !found {
				return errors.WithMessagef(
					ErrJumpNotLabelID,
					"%d",
					idx)
			}

			tgt, found := labelPC[lid]
			if !found {
				return errors.WithMessagef(
					ErrUnboundLabel,
					"%v at %d",
					lid,
					idx)
			}

			isn.Args[0] = tgt
		}

		stripped = append(stripped, isn)
	}

	p.Code = stripped

	return nil
}

// ** Utility Methods:

func (p *Program) addField(val string) int {
	for idx := range p.Fields {
		if val == p.Fields[idx] {
			return idx
		}
	}

	p.Fields = append(p.Fields, val)

	return len(p.Fields) - 1
}

func (p *Program) addString(val string) int {
	for idx := range p.Strings {
		if val == p.Strings[idx] {
			return idx
		}
	}

	p.Strings = append(p.Strings, val)

	return len(p.Strings) - 1
}

func (p *Program) addNumber(val float64) int {
	for idx := range p.Numbers {
		if val == p.Numbers[idx] {
			return idx
		}
	}

	p.Numbers = append(p.Numbers, val)

	return len(p.Numbers) - 1
}

func (p *Program) addTime(val int64) int {
	for idx := range p.Times {
		if val == p.Times[idx] {
			return idx
		}
	}

	p.Times = append(p.Times, val)

	return len(p.Times) - 1
}

func (p *Program) addIP(val netip.Addr) int {
	for idx := range p.IPs {
		if val == p.IPs[idx] {
			return idx
		}
	}

	p.IPs = append(p.IPs, val)

	return len(p.IPs) - 1
}

func (p *Program) addRegex(val *regexp.Regexp) int {
	for idx := range p.Patterns {
		if val.String() == p.Patterns[idx].String() {
			return idx
		}
	}

	p.Patterns = append(p.Patterns, val)

	return len(p.Patterns) - 1
}

func (p *Program) appendIsn(op OpCode, args ...any) {
	p.Code = append(p.Code, Instr{Op: op, Args: args})
}

// ** Generation Methods:

func (p *Program) emitAnd(node *TypedNodeAnd, tlbl, flbl LabelID) {
	kidcount := len(node.kids)

	if kidcount == 0 {
		p.appendJump(opJMP, tlbl)

		return
	}

	for idx := range kidcount - 1 {
		cont := p.newLabel()

		p.emit(node.kids[idx], cont, flbl)
		p.bindLabel(cont)
	}

	p.emit(node.kids[kidcount-1], tlbl, flbl)
}

func (p *Program) emitOr(node *TypedNodeOr, tlbl, flbl LabelID) {
	kidcount := len(node.kids)

	if kidcount == 0 {
		p.appendJump(opJMP, flbl)

		return
	}

	for idx := range kidcount - 1 {
		cont := p.newLabel()

		p.emit(node.kids[idx], tlbl, cont)
		p.bindLabel(cont)
	}

	p.emit(node.kids[kidcount-1], tlbl, flbl)
}

func (p *Program) emitNot(node *TypedNodeNot, tlbl, flbl LabelID) {
	if node.kid == nil {
		p.appendJump(opJMP, flbl)

		return
	}

	cont := p.newLabel()
	p.emit(node.kid, tlbl, cont)
	p.appendIsn(opNOT)
	p.bindLabel(cont)
}

func (p *Program) emitEqS(node *TypedNodeEqS, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	sidx := p.addString(node.value)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(opEQ_S, sidx)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitNeqS(node *TypedNodeNeqS, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	sidx := p.addString(node.value)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(opNEQ_S, sidx)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitCmpN(node *TypedNodeCmpN, tlbl, flbl LabelID) {
	op := getNComparator(node.op, opEQ_N)
	fidx := p.addField(node.field)
	nidx := p.addNumber(node.value)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(op, nidx)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitCmpT(node *TypedNodeCmpT, tlbl, flbl LabelID) {
	op := getTComparator(node.op, opEQ_N)
	fidx := p.addField(node.field)
	tidx := p.addTime(node.value)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(op, tidx)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitCmpIP(node *TypedNodeCmpIP, tlbl, flbl LabelID) {
	op := getIPComparator(node.op, opEQ_N)
	fidx := p.addField(node.field)
	ipidx := p.addIP(node.value)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(op, ipidx)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitRangeN(node *TypedNodeRangeN, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	p.appendIsn(opSET_F, fidx)

	lowIdx, highIdx := -1, -1

	if node.low != nil {
		lowIdx = p.addNumber(*node.low)
	}

	if node.high != nil {
		highIdx = p.addNumber(*node.high)
	}

	p.appendIsn(opRANGE_N, lowIdx, highIdx, node.incl, node.inch)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitRangeT(node *TypedNodeRangeT, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	p.appendIsn(opSET_F, fidx)

	lowIdx, highIdx := -1, -1

	if node.low != nil {
		lowIdx = p.addTime(*node.low)
	}

	if node.high != nil {
		highIdx = p.addTime(*node.high)
	}

	p.appendIsn(opRANGE_T, lowIdx, highIdx, node.incl, node.inch)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitRangeIP(node *TypedNodeRangeIP, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	p.appendIsn(opSET_F, fidx)

	lowIdx := p.addIP(node.low)
	highIdx := p.addIP(node.high)

	p.appendIsn(opRANGE_IP, lowIdx, highIdx, node.incl, node.inch)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitPrefix(node *TypedNodePrefix, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	sidx := p.addString(node.prefix)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(opPREFIX, sidx)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitGlob(node *TypedNodeGlob, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	gidx := p.addString(node.glob)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(opGLOB, gidx)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitPhrase(node *TypedNodePhrase, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	pidx := p.addString(node.phrase)

	p.appendIsn(opSET_F, fidx)

	if node.fuzz != nil {
		p.appendIsn(opFUZZ, *node.fuzz)
	}

	if node.boost != nil {
		p.appendIsn(opBOOST, *node.boost)
	}

	p.appendIsn(opPHRASE, pidx, node.prox)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitExists(node *TypedNodeExists, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(opEXISTS)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

func (p *Program) emitRegex(node *TypedNodeRegex, tlbl, flbl LabelID) {
	fidx := p.addField(node.field)
	ridx := p.addRegex(node.compiled)

	p.appendIsn(opSET_F, fidx)
	p.appendIsn(opREGEX, ridx)
	p.appendJump(opJNZ, tlbl)
	p.appendJump(opJMP, flbl)
}

//nolint:cyclop
func (p *Program) emit(node TypedNode, tlbl, flbl LabelID) {
	switch val := node.(type) {
	case *TypedNodeTrue:
		p.appendIsn(opLDA, 1)
	case *TypedNodeFalse:
		p.appendIsn(opLDA, 0)

	case *TypedNodeAnd:
		p.emitAnd(val, tlbl, flbl)
	case *TypedNodeOr:
		p.emitOr(val, tlbl, flbl)
	case *TypedNodeNot:
		p.emitNot(val, tlbl, flbl)

	case *TypedNodeEqS:
		p.emitEqS(val, tlbl, flbl)
	case *TypedNodeNeqS:
		p.emitNeqS(val, tlbl, flbl)

	case *TypedNodePrefix:
		p.emitPrefix(val, tlbl, flbl)
	case *TypedNodeGlob:
		p.emitGlob(val, tlbl, flbl)
	case *TypedNodeRegex:
		p.emitRegex(val, tlbl, flbl)
	case *TypedNodePhrase:
		p.emitPhrase(val, tlbl, flbl)
	case *TypedNodeExists:
		p.emitExists(val, tlbl, flbl)

	case *TypedNodeCmpN:
		p.emitCmpN(val, tlbl, flbl)
	case *TypedNodeCmpT:
		p.emitCmpT(val, tlbl, flbl)
	case *TypedNodeCmpIP:
		p.emitCmpIP(val, tlbl, flbl)

	case *TypedNodeRangeN:
		p.emitRangeN(val, tlbl, flbl)
	case *TypedNodeRangeT:
		p.emitRangeT(val, tlbl, flbl)
	case *TypedNodeRangeIP:
		p.emitRangeIP(val, tlbl, flbl)
	}
}

func (p *Program) Emit(node TypedNode) {
	troot := p.newLabel()
	froot := p.newLabel()
	done := p.newLabel()

	p.emit(node, troot, froot)

	p.bindLabel(troot)
	p.appendIsn(opLDA, 1)
	p.appendJump(opJMP, done)

	p.bindLabel(froot)
	p.appendIsn(opLDA, 0)

	p.bindLabel(done)
	p.appendIsn(opRET)

	p.Peephole()

	if err := p.resolveLabels(); err != nil {
		panic(err)
	}
}

// ** Functions:

func getNComparator(cmp CmpKind, def OpCode) OpCode {
	if code, found := cmpToOpN[cmp]; found {
		return code
	}

	return def
}

func getTComparator(cmp CmpKind, def OpCode) OpCode {
	if code, found := cmpToOpT[cmp]; found {
		return code
	}

	return def
}

func getIPComparator(cmp CmpKind, def OpCode) OpCode {
	if code, found := cmpToOpIP[cmp]; found {
		return code
	}

	return def
}

func NewProgram() *Program {
	return &Program{
		Fields:   []string{},
		Strings:  []string{},
		Numbers:  []float64{},
		Times:    []int64{},
		IPs:      []netip.Addr{},
		Patterns: []*regexp.Regexp{},
		Code:     []Instr{}}
}

// * codegen.go ends here.
