// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// instruction.go --- Instructions.
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

	"github.com/Asmodai/gohacks/utils"
)

// * Constants:

//nolint:godot
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

	// Default padding used when pretty-printing instructions
	defaultIsnPadding = 10
)

// * Variables:

var (
	// Map of `opcode -> string` for pretty-printing.
	//
	//nolint:gochecknoglobals
	opNames = map[OpCode]string{
		OpNoOp:        "NOP",
		OpReturn:      "RET",
		OpLabel:       "LABEL",
		OpJump:        "JMP",
		OpJumpZ:       "JZ",
		OpJumpNZ:      "JNZ",
		OpNot:         "NOT",
		OpLoadA:       "LDA",
		OpLoadField:   "LDFLD",
		OpLoadBoost:   "LDBST",
		OpLoadFuzzy:   "LDFZY",
		OpStringEQ:    "EQ.S",
		OpStringNEQ:   "NEQ.S",
		OpPrefix:      "PFX.S",
		OpGlob:        "GLB.S",
		OpRegex:       "REX.S",
		OpPhrase:      "PHR.S",
		OpAny:         "ANY",
		OpNumberEQ:    "EQ.N",
		OpNumberNEQ:   "NEQ.N",
		OpNumberLT:    "LT.N",
		OpNumberLTE:   "LTE.N",
		OpNumberGT:    "GT.N",
		OpNumberGTE:   "GTE.N",
		OpNumberRange: "RNG.N",
		OpTimeEQ:      "EQ.T",
		OpTimeNEQ:     "NEQ.T",
		OpTimeLT:      "LT.T",
		OpTimeLTE:     "LTE.T",
		OpTimeGT:      "GT.T",
		OpTimeGTE:     "GTE.T",
		OpTimeRange:   "RNG.T",
		OpIPEQ:        "EQ.IP",
		OpIPNEQ:       "NEQ.IP",
		OpIPLT:        "LT.IP",
		OpIPLTE:       "LTE.IP",
		OpIPGT:        "GT.IP",
		OpIPGTE:       "GTE.IP",
		OpIPRange:     "RNG.IP",
		OpInCIDR:      "IN.CIDR",
	}

	// Map of `ComparatorKind -> OpCode` for the conversion of
	// comparators to numeric opcode.
	//
	//nolint:gochecknoglobals
	cmpToOpN = map[ComparatorKind]OpCode{
		ComparatorEQ:  OpNumberEQ,
		ComparatorNEQ: OpNumberNEQ,
		ComparatorLT:  OpNumberLT,
		ComparatorLTE: OpNumberLTE,
		ComparatorGT:  OpNumberGT,
		ComparatorGTE: OpNumberGTE,
	}

	// Map of `ComparatorKind -> OpCode` for the conversion of
	// comparators to date/time opcode.
	//
	//nolint:gochecknoglobals
	cmpToOpT = map[ComparatorKind]OpCode{
		ComparatorEQ:  OpTimeEQ,
		ComparatorNEQ: OpTimeNEQ,
		ComparatorLT:  OpTimeLT,
		ComparatorLTE: OpTimeLTE,
		ComparatorGT:  OpTimeGT,
		ComparatorGTE: OpTimeGTE,
	}

	// Map of `ComparatorKind -> OpCode' for the conversion of
	// comparators to IP address opcode.
	//
	//nolint:gochecknoglobals
	cmpToOpIP = map[ComparatorKind]OpCode{
		ComparatorEQ:  OpIPEQ,
		ComparatorNEQ: OpIPNEQ,
		ComparatorLT:  OpIPLT,
		ComparatorLTE: OpIPLTE,
		ComparatorGT:  OpIPGT,
		ComparatorGTE: OpIPGTE,
	}
)

// * Code:

// ** Types:

// Opcode.
type OpCode int

// * Structure:

type Instr struct {
	Op   OpCode // Instruction's opcode.
	Args []any  // Instruction's operands.
}

// ** Methods:

// Return the string representation of an instruction.
func (isn Instr) String() string {
	str, found := opNames[isn.Op]
	if !found {
		return "<invalid>"
	}

	return fmt.Sprintf("%s%v",
		utils.Pad(str, defaultIsnPadding),
		isn.Args)
}

// Is the instruction a jump of some kind?
func (isn Instr) IsJump() bool {
	return isn.Op == OpJump || isn.Op == OpJumpNZ || isn.Op == OpJumpZ
}

// * instruction.go ends here.
