// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// program.go --- Code generator.
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
	"strconv"
	"strings"

	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	// Initial size of the labels array.
	initialLabelArraySize = 16

	// Unknown constant entity.
	unknownEntity = "?"
)

// * Code:

// ** Types:

// Label identifier type.
type LabelID int

// ** Structure:

type Program struct {
	Fields      []string         // Field constants.
	Strings     []string         // String constants.
	Numbers     []float64        // Number constants.
	Times       []int64          // Date/time constants.
	IPs         []netip.Addr     // IP address constants.
	Patterns    []*regexp.Regexp // Regular expression constants.
	Code        []Instr          // Bytecode.
	nextLabelID LabelID          // Label ID counter.
}

// ** Accessors:

func (p *Program) fieldName(fid int) string {
	if fid >= 0 && fid < len(p.Fields) {
		return p.Fields[fid]
	}

	return unknownEntity
}

func (p *Program) stringConstant(idx int) string {
	if idx >= 0 && idx < len(p.Strings) {
		str := p.Strings[idx]

		str = strings.ReplaceAll(str, "\n", "\\n")
		str = strings.ReplaceAll(str, "\t", "\\t")

		return `"` + str + `"`
	}

	return unknownEntity
}

func (p *Program) numberConstant(idx int) string {
	if idx >= 0 && idx < len(p.Numbers) {
		return fmt.Sprintf("%g", p.Numbers[idx])
	}

	return unknownEntity
}

func (p *Program) timeConstant(idx int) string {
	if idx >= 0 && idx < len(p.Times) {
		return strconv.FormatInt(p.Times[idx], 10)
	}

	return unknownEntity
}

func (p *Program) ipConstant(idx int) string {
	if idx >= 0 && idx < len(p.IPs) {
		return p.IPs[idx].String()
	}

	return unknownEntity
}

func (p *Program) regexConstant(idx int) string {
	if idx >= 0 && idx < len(p.Patterns) {
		str := p.Patterns[idx].String()

		return `/` + str + `/`
	}

	return unknownEntity
}

// ** Label methods:

// Generate a new label identifier and return it.
func (p *Program) NewLabel() LabelID {
	id := p.nextLabelID
	p.nextLabelID++

	return id
}

// Append a `LABEL` instruction bound to the given label identifier.
//
// The instruction will be removed by the label resolver.
func (p *Program) BindLabel(id LabelID) {
	p.AppendIsn(OpLabel, id)
}

// Resolve labels to absolute addresses.
//
//nolint:cyclop,funlen
func (p *Program) resolveLabels() error {
	//
	// Phase 1: Build `label -> PC` map.
	labelPC := make(map[LabelID]int, initialLabelArraySize)
	pcOut := 0

	for idx := range p.Code {
		isn := p.Code[idx]

		if isn.Op == OpLabel {
			// Bind the label to the next real instruction.
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

	// NOTE:
	// Safety: ensure labels don't point past end -- i.e., ensure there's
	// a real instruction after each label.
	//
	// If you *intentionally* bind a label at end, make sure a RET
	// follows before peephole.

	//
	// Phase 2: Build stripped code and patch jumps.
	stripped := make([]Instr, 0, pcOut)

	for idx := range len(p.Code) {
		isn := p.Code[idx]

		if isn.Op == OpLabel {
			continue
		}

		// Patch jump.
		if isn.IsJump() {
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

// ** Instruction methods:

// Append a jump instruction to the bytecode.
//
// The instruction's operand will be the target label ID.
func (p *Program) AppendJump(opCode OpCode, target LabelID) {
	p.AppendIsn(opCode, target)
}

// Append an instruction to the bytecode.
//
// Any provided arguments will be added as the instruction's operands.
func (p *Program) AppendIsn(opCode OpCode, args ...any) {
	p.Code = append(p.Code, Instr{Op: opCode, Args: args})
}

// ** Constant methods:

// Add a field name constant.
//
// If the given constant value exists, then an index to its array position
// is returned.
func (p *Program) AddFieldConstant(val string) int {
	for idx := range p.Fields {
		if val == p.Fields[idx] {
			return idx
		}
	}

	p.Fields = append(p.Fields, val)

	return len(p.Fields) - 1
}

// Add a string constant.
//
// If the given constant value exists, then an index to its array position
// is returned.
func (p *Program) AddStringConstant(val string) int {
	for idx := range p.Strings {
		if val == p.Strings[idx] {
			return idx
		}
	}

	p.Strings = append(p.Strings, val)

	return len(p.Strings) - 1
}

// Add a numeric constant.
//
// If the given constant value exists, then an index to its array position
// is returned.
func (p *Program) AddNumberConstant(val float64) int {
	for idx := range p.Numbers {
		if val == p.Numbers[idx] {
			return idx
		}
	}

	p.Numbers = append(p.Numbers, val)

	return len(p.Numbers) - 1
}

// Add a date/time constant.
//
// If the given constant value exists, then an index to its array position
// is returned.
func (p *Program) AddTimeConstant(val int64) int {
	for idx := range p.Times {
		if val == p.Times[idx] {
			return idx
		}
	}

	p.Times = append(p.Times, val)

	return len(p.Times) - 1
}

// Add an IP address constant.
//
// If the given constant value exists, then an index to its array position
// is returned.
func (p *Program) AddIPConstant(val netip.Addr) int {
	for idx := range p.IPs {
		if val == p.IPs[idx] {
			return idx
		}
	}

	p.IPs = append(p.IPs, val)

	return len(p.IPs) - 1
}

// Add a regular expression constant.
//
// If the given constant value exists, then an index to its array position
// is returned.
func (p *Program) AddRegexConstant(val *regexp.Regexp) int {
	for idx := range p.Patterns {
		if val.String() == p.Patterns[idx].String() {
			return idx
		}
	}

	p.Patterns = append(p.Patterns, val)

	return len(p.Patterns) - 1
}

// ** Generation methods:

func (p *Program) Emit(irNode IRNode) {
	tRoot := p.NewLabel() // Root `true` CPS label.
	fRoot := p.NewLabel() // Root `false` CPS label.
	done := p.NewLabel()  // `done` CPS label.

	// Emit the instructions in the IR tree.
	irNode.Emit(p, tRoot, fRoot)

	// The root continuation for a true result.
	p.BindLabel(tRoot)
	p.AppendIsn(OpLoadA, 1)    // Set accumulator to 1.
	p.AppendJump(OpJump, done) // Jump out.

	// The root continuation for a false result.
	p.BindLabel(fRoot)
	p.AppendIsn(OpLoadA, 0) // Set accumulator to 0.

	// The root exit point.
	p.BindLabel(done)
	p.AppendIsn(OpReturn) // Halt the program.

	// Invoke the peephole optimiser.
	p.Peephole()

	// Resolve labels.
	if err := p.resolveLabels(); err != nil {
		// XXX Would be nice to handle this properly.
		panic(err)
	}
}

// ** Functions:

// Return the relevant numeric opcode for the given comparator.
//
// If no opcode is found, then `def` will be used instead.
func GetNumberComparator(cmp ComparatorKind, def OpCode) OpCode {
	if code, found := cmpToOpN[cmp]; found {
		return code
	}

	return def
}

// Return the relevant date/time opcode for the given comparator.
//
// If no opcode is found, then `def` will be used instead.
func GetTimeComparator(cmp ComparatorKind, def OpCode) OpCode {
	if code, found := cmpToOpT[cmp]; found {
		return code
	}

	return def
}

// Return the relevant IP address opcode for the given comparator.
//
// If no opcode is found, then `def` will be used instead.
func GetIPComparator(cmp ComparatorKind, def OpCode) OpCode {
	if code, found := cmpToOpIP[cmp]; found {
		return code
	}

	return def
}

// Create a new program instance.
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

// * program.go ends here.
