// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// disassembler.go --- Disassembler.
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
	"io"
	"strconv"
	"strings"
)

// * Constants:

const (
	constantField  = "F"
	constantString = "S"
	constantNumber = "N"
	constantTime   = "T"
	constantIP     = "IP"
	constantRegex  = "R"

	indentIsn      = 8
	widthAddress   = 3
	paddingOperand = 2

	ppRangeNumber ppRangeType = iota
	ppRangeTime
	ppRangeIP
)

// * Variables:

// * Code:

// ** Types:

type ppRangeType int

// ** Options:

// *** Structure:

type DisassemblerOpts struct {
	WithComments bool // Include decoded comments?
	AddrWidth    int  // Width of an address.  0 = auto.
	OpcodeWidth  int  // Pad opcode column. 0 auto.
	OperandWidth int  // Pad operand column. 0 = auto.
}

// *** Functions:

func NewDefaultDisassemblerOpts() DisassemblerOpts {
	return DisassemblerOpts{
		WithComments: true,
		AddrWidth:    0,
		OpcodeWidth:  0,
		OperandWidth: 0}
}

// ** Disassembler:

// *** Structure:

type Disassembler struct {
	program      *Program
	opts         DisassemblerOpts
	labels       map[int]string
	opcodeWidth  int
	operandWidth int
	addrWidth    int
}

// *** Methods:

func (d *Disassembler) reset() {
	d.program = nil
	d.labels = make(map[int]string, 0)
}

func (d *Disassembler) computeOpcodeWidth() {
	opw := d.opts.OpcodeWidth

	if opw == 0 {
		for _, isn := range d.program.Code {
			if oplen := len(opNames[isn.Op]); oplen > opw {
				opw = oplen
			}
		}
	}

	d.opcodeWidth = opw
}

func (d *Disassembler) computeOperandWidth() {
	argw := d.opts.OperandWidth

	if argw == 0 {
		for _, isn := range d.program.Code {
			if arglen := len(renderArgsBare(isn)); arglen > argw {
				argw = arglen
			}
		}
	}

	d.operandWidth = argw + paddingOperand
}

func (d *Disassembler) computeAddressWidth() {
	addrw := d.opts.AddrWidth

	if addrw == 0 {
		addrw = len(strconv.Itoa(len(d.program.Code)))

		if addrw < widthAddress {
			addrw = widthAddress
		}
	}

	d.addrWidth = addrw
}

func (d *Disassembler) computeLabels() {
	seq := 1

	d.labels = make(map[int]string, 0)

	maybeLabel := func(loc int) {
		if loc < 0 || loc >= len(d.program.Code) {
			return
		}

		if _, ok := d.labels[loc]; ok {
			return
		}

		d.labels[loc] = fmt.Sprintf("L%03d", seq)
		seq++
	}

	for _, isn := range d.program.Code {
		if isn.IsJump() {
			if len(isn.Args) == 0 {
				continue
			}

			if tgt, ok := isn.Args[0].(int); ok {
				maybeLabel(tgt)
			}
		}
	}
}

func (d *Disassembler) SetProgram(program *Program) {
	d.reset()

	d.program = program

	d.computeOpcodeWidth()
	d.computeOperandWidth()
	d.computeAddressWidth()
	d.computeLabels()
}

func (d *Disassembler) Dissassemble(writer io.Writer) {
	code := d.program.Code
	isnIndent := strings.Repeat(" ", indentIsn)

	for pcounter, isn := range code {
		if name, found := d.labels[pcounter]; found {
			fmt.Fprintf(writer, "%s\n", name)
		}

		addr := fmt.Sprintf("%0*d", d.addrWidth, pcounter)
		opcode := opNames[isn.Op]
		args := renderArgsBare(isn)

		if d.opts.WithComments {
			cmt := renderComment(d.program, isn, d.labels)

			if len(cmt) > 0 {
				fmt.Fprintf(writer,
					"%s%s\t%-*s  %-*s  ; %s\n",
					isnIndent,
					addr,
					d.opcodeWidth,
					opcode,
					d.operandWidth,
					args,
					cmt)

				continue
			}
		}

		fmt.Fprintf(writer,
			"%s%s\t%-*s  %-*s\n",
			isnIndent,
			addr,
			d.opcodeWidth,
			opcode,
			d.operandWidth,
			args)
	}
}

// ** Functions:

func fmtConstant(pool, val string) string {
	return pool + "[" + val + "]"
}

//nolint:cyclop
func ppRange(prog *Program, args []any, mode ppRangeType) string {
	low, _ := args[0].(int)
	high, _ := args[1].(int)
	incL, _ := args[2].(bool)
	incH, _ := args[3].(bool)
	loStr := "-∞"
	hiStr := "+∞"
	openStr := "("
	closeStr := ")"

	switch mode {
	case ppRangeNumber: // Numeric.
		if low >= 0 {
			loStr = fmtConstant(constantNumber,
				prog.numberConstant(low))
		}

		if high >= 0 {
			hiStr = fmtConstant(constantNumber,
				prog.numberConstant(high))
		}

	case ppRangeTime: // Time.
		if low >= 0 {
			loStr = fmtConstant(constantTime,
				prog.timeConstant(low))
		}

		if high >= 0 {
			hiStr = fmtConstant(constantTime,
				prog.timeConstant(high))
		}

	case ppRangeIP: // IP.
		if low >= 0 {
			loStr = fmtConstant(constantIP,
				prog.ipConstant(low))
		}

		if high >= 0 {
			hiStr = fmtConstant(constantIP,
				prog.ipConstant(high))
		}
	}

	if incL {
		openStr = "["
	}

	if incH {
		closeStr = "]"
	}

	return fmt.Sprintf("%s%s..%s%s",
		openStr,
		loStr,
		hiStr,
		closeStr)
}

//nolint:cyclop,exhaustive,funlen,mnd
func renderArgsBare(isn Instr) string {
	args := isn.Args
	length := len(args)

	switch isn.Op {
	case OpLoadA:
		if length == 1 {
			return fmt.Sprintf("imm=%v", args[0])
		}

	case OpLoadField:
		if length == 1 {
			return fmt.Sprintf("fid=%v", args[0])
		}

	case OpLoadBoost:
		if length == 1 {
			return fmt.Sprintf("boost=%v", args[0])
		}

	case OpLoadFuzzy:
		if length == 1 {
			return fmt.Sprintf("fuzz=%v", args[0])
		}

	case OpJump, OpJumpNZ, OpJumpZ:
		if length == 1 {
			if tgt, ok := args[0].(int); ok {
				return fmt.Sprintf("%03d", tgt)
			}
		}

	case OpStringEQ, OpStringNEQ, OpPrefix, OpGlob:
		if length == 1 {
			return fmt.Sprintf("sIdx=%v", args[0])
		}

	case OpPhrase:
		if length == 2 {
			return fmt.Sprintf("sIdx=%v prox=%v", args[0], args[1])
		}

	case OpRegex:
		if length == 1 {
			return fmt.Sprintf("rIdx=%v", args[0])
		}

	case OpNumberRange, OpTimeRange, OpIPRange:
		if length == 4 {
			return fmt.Sprintf(
				"lo=%v hi=%v incL=%v incH=%v",
				args[0], args[1], args[2], args[3])
		}

	case OpNumberEQ, OpNumberNEQ, OpNumberGT, OpNumberLT, OpNumberGTE, OpNumberLTE:
		if length == 1 {
			return fmt.Sprintf("nIdx=%v", args[0])
		}

	case OpTimeEQ, OpTimeNEQ, OpTimeGT, OpTimeLT, OpTimeGTE, OpTimeLTE:
		if length == 1 {
			return fmt.Sprintf("tIdx=%v", args[0])
		}

	case OpIPEQ, OpIPNEQ, OpIPGT, OpIPLT, OpIPGTE, OpIPLTE:
		if length == 1 {
			return fmt.Sprintf("ipIdx=%v", args[0])
		}
	}

	if length == 0 {
		return ""
	}

	parts := make([]string, length)
	for idx, val := range args {
		parts[idx] = fmt.Sprintf("%v", val)
	}

	return strings.Join(parts, " ")
}

func renderJump(isn Instr, loc string) string {
	//nolint:exhaustive
	switch isn.Op {
	case OpJumpNZ:
		return "jump " + loc + " if not zero"

	case OpJumpZ:
		return "jump " + loc + " if zero"

	default:
		return "jump " + loc
	}
}

//nolint:cyclop,gocognit,gocyclo,funlen
func renderComment(prog *Program, isn Instr, lbls map[int]string) string {
	args := isn.Args
	length := len(args)

	//nolint:exhaustive
	switch isn.Op {
	case OpLoadA:
		if length == 1 {
			return fmt.Sprintf("set accumulator=%v", args[0])
		}

	case OpLoadField:
		if length == 1 {
			if fid, ok := args[0].(int); ok {
				return "set field=" +
					fmtConstant(constantField,
						prog.fieldName(fid))
			}
		}

	case OpLoadBoost:
		if length == 1 {
			if val, ok := args[0].(float64); ok {
				return fmt.Sprintf("set boost=%.3f", val)
			}
		}

	case OpLoadFuzzy:
		if length == 1 {
			if val, ok := args[0].(float64); ok {
				return fmt.Sprintf("set fuzz=%.3f", val)
			}
		}

	case OpJump, OpJumpNZ, OpJumpZ:
		if length == 1 {
			if tgt, ok := args[0].(int); ok {
				if name, has := lbls[tgt]; has {
					return renderJump(isn, name)
				}

				return fmt.Sprintf(
					"jump L%03d",
					tgt)
			}
		}

	case OpStringEQ, OpStringNEQ, OpPrefix, OpGlob:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return fmtConstant(constantString,
					prog.stringConstant(idx))
			}
		}

	case OpRegex:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return fmtConstant(constantRegex,
					prog.regexConstant(idx))
			}
		}

	case OpPhrase:
		if length == 2 { //nolint:mnd
			idx, _ := args[0].(int)
			prox, _ := args[1].(int)

			if prox == 0 {
				return fmtConstant(constantString,
					prog.stringConstant(idx))
			}

			return fmt.Sprintf(
				"%s, prox=%d",
				fmtConstant(constantString,
					prog.stringConstant(idx)),
				prox)
		}

	case OpNumberEQ, OpNumberNEQ, OpNumberGT, OpNumberLT, OpNumberGTE, OpNumberLTE:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return prog.numberConstant(idx)
			}
		}

	case OpTimeEQ, OpTimeNEQ, OpTimeGT, OpTimeLT, OpTimeGTE, OpTimeLTE:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return prog.timeConstant(idx)
			}
		}

	case OpIPEQ, OpIPNEQ, OpIPGT, OpIPLT, OpIPGTE, OpIPLTE:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return prog.ipConstant(idx)
			}
		}

	case OpNumberRange:
		if length == 4 { //nolint:mnd
			return ppRange(prog, args, ppRangeNumber)
		}

	case OpTimeRange:
		if length == 4 { //nolint:mnd
			return ppRange(prog, args, ppRangeTime)
		}

	case OpIPRange:
		if length == 4 { //nolint:mnd
			return ppRange(prog, args, ppRangeIP)
		}

	case OpInCIDR:
		if length == 2 { //nolint:mnd
			ipIdx, _ := args[0].(int)
			mask, _ := args[1].(int)

			return fmt.Sprintf("%s, mask=%d",
				fmtConstant(constantIP,
					prog.ipConstant(ipIdx)),
				mask)
		}
	}

	return ""
}

func NewDisassembler(opts DisassemblerOpts) *Disassembler {
	return &Disassembler{opts: opts}
}

func NewDefaultDisassembler() *Disassembler {
	defaults := NewDefaultDisassemblerOpts()

	return NewDisassembler(defaults)
}

// * disassembler.go ends here.
