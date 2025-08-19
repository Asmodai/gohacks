// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// disassembler.go --- Cute disassembler.
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
	unknownField  = "F[?]"
	unknownString = "S[?]"
	unknownNumber = "N[?]"
	unknownTime   = "T[?]"
	unknownIP     = "IP[?]"
	unknownRegex  = "R[?]"

	PPRangeN PPRangeType = iota
	PPRangeT
	PPRangeIP

	instructionIndent = 8
	minAddressWidth   = 3
	operandWidthPad   = 2
)

// * Variables:

// * Code:

// ** Types:

type PPRangeType int

// ** Pretty printer options:

type PrettyPrinterOpts struct {
	WithComments bool // Include decoded comments.
	AddrWidth    int  // Width of address digits.  0 = auto.
	OpcodeWidth  int  // Pad opcode column. 0 = auto.
	OperandWidth int  // Pad operand column. 0 = auto.
}

// ** Program Methods:

func (p *Program) fieldName(fid int) string {
	if fid >= 0 && fid < len(p.Fields) {
		return p.Fields[fid]
	}

	return unknownField
}

func (p *Program) stringLiteral(idx int) string {
	if idx >= 0 && idx < len(p.Strings) {
		str := p.Strings[idx]

		str = strings.ReplaceAll(str, "\n", "\\n")
		str = strings.ReplaceAll(str, "\t", "\\t")

		return `"` + str + `"`
	}

	return unknownString
}

func (p *Program) numberLiteral(idx int) string {
	if idx >= 0 && idx < len(p.Numbers) {
		return fmt.Sprintf("%g", p.Numbers[idx])
	}

	return unknownNumber
}

func (p *Program) timeLiteral(idx int) string {
	if idx >= 0 && idx < len(p.Times) {
		return strconv.FormatInt(p.Times[idx], 10)
	}

	return unknownTime
}

func (p *Program) ipLiteral(idx int) string {
	if idx >= 0 && idx < len(p.IPs) {
		return p.IPs[idx].String()
	}

	return unknownIP
}

func (p *Program) regexLiteral(idx int) string {
	if idx >= 0 && idx < len(p.Patterns) {
		str := p.Patterns[idx].String()

		return `/` + str + `/`
	}

	return unknownRegex
}

// ** Pretty Printer:

func (p *Program) opcodeWidth(opts PrettyPrinterOpts) int {
	opw := opts.OpcodeWidth

	if opw == 0 {
		for _, isn := range p.Code {
			if oplen := len(opNames[isn.Op]); oplen > opw {
				opw = oplen
			}
		}
	}

	return opw + operandWidthPad
}

func (p *Program) operandWidth(opts PrettyPrinterOpts) int {
	argw := opts.OperandWidth

	if argw == 0 {
		for _, isn := range p.Code {
			if arglen := len(renderArgsBare(p, isn)); arglen > argw {
				argw = arglen
			}
		}
	}

	return argw
}

func (p *Program) addressWidth(opts PrettyPrinterOpts) int {
	addrw := opts.AddrWidth

	if addrw == 0 {
		addrw = len(strconv.Itoa(len(p.Code)))

		if addrw < minAddressWidth {
			addrw = minAddressWidth
		}
	}

	return addrw
}

//nolint:cyclop
func (p *Program) ppRange(args []any, mode PPRangeType) string {
	low, _ := args[0].(int)
	high, _ := args[1].(int)
	incL, _ := args[2].(bool)
	incH, _ := args[3].(bool)
	loStr := "-∞"
	hiStr := "+∞"
	openStr := "("
	closeStr := ")"

	switch mode {
	case PPRangeN: // Numeric.
		if low >= 0 {
			loStr = p.numberLiteral(low)
		}

		if high >= 0 {
			hiStr = p.numberLiteral(high)
		}

	case PPRangeT: // Time.
		if low >= 0 {
			loStr = p.timeLiteral(low)
		}

		if high >= 0 {
			hiStr = p.timeLiteral(high)
		}

	case PPRangeIP: // IP.
		if low >= 0 {
			loStr = p.ipLiteral(low)
		}

		if high >= 0 {
			hiStr = p.ipLiteral(high)
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

func (p *Program) PrettyPrint(writer io.Writer, opts PrettyPrinterOpts) {
	code := p.Code
	lbls := computeLabels(code)
	opwidth := p.opcodeWidth(opts)
	argwidth := p.operandWidth(opts)
	addrwidth := p.addressWidth(opts)
	isnIndent := strings.Repeat(" ", instructionIndent)

	for pc, isn := range code {
		if name, has := lbls[pc]; has {
			fmt.Fprintf(writer, "%s:\n", name)
		}

		addr := fmt.Sprintf("%0*d", addrwidth, pc)
		opcode := opNames[isn.Op]
		args := renderArgsBare(p, isn)

		if opts.WithComments {
			cmt := renderComment(p, isn, lbls)

			if len(cmt) > 0 {
				fmt.Fprintf(writer,
					"%s%s\t%-*s  %-*s  ; %s\n",
					isnIndent,
					addr,
					opwidth,
					opcode,
					argwidth,
					args,
					cmt)

				continue
			}
		}

		fmt.Fprintf(writer,
			"%s%s\t%-*s  %s\n",
			isnIndent,
			addr,
			opwidth,
			opcode,
			args)
	}
}

// ** Functions:

//nolint:cyclop,exhaustive,funlen
func renderArgsBare(_ *Program, isn Instr) string {
	args := isn.Args
	length := len(args)

	switch isn.Op {
	case opLDA:
		if length == 1 {
			return fmt.Sprintf("imm=%v", args[0])
		}

	case opSET_F:
		if length == 1 {
			return fmt.Sprintf("fid=%v", args[0])
		}

	case opJMP, opJZ, opJNZ:
		if length == 1 {
			if tgt, ok := args[0].(int); ok {
				return fmt.Sprintf("-> %d", tgt)
			}
		}

	case opEQ_S, opNEQ_S, opPREFIX, opGLOB:
		if length == 1 {
			return fmt.Sprintf("sIdx=%v", args[0])
		}

	case opPHRASE:
		if length == 2 { //nolint:mnd
			return fmt.Sprintf("sIdx=%v prox=%v", args[0], args[1])
		}

	case opREGEX:
		if length == 1 {
			return fmt.Sprintf("rIdx=%v", args[0])
		}

	case opRANGE_N, opRANGE_T, opRANGE_IP:
		if length == 4 { //nolint:mnd
			return fmt.Sprintf(
				"lo=%v hi=%v incL=%v incH=%v",
				args[0],
				args[1],
				args[2],
				args[3])
		}

	case opEQ_N, opNEQ_N, opLT_N, opLTE_N, opGT_N, opGTE_N:
		if length == 1 {
			return fmt.Sprintf("nIdx=%v", args[0])
		}

	case opEQ_T, opNEQ_T, opLT_T, opLTE_T, opGT_T, opGTE_T:
		if length == 1 {
			return fmt.Sprintf("tIdx=%v", args[0])
		}

	case opEQ_IP, opNEQ_IP, opLT_IP, opLTE_IP, opGT_IP, opGTE_IP:
		if length == 1 {
			return fmt.Sprintf("ipIdx=%v", args[0])
		}

	case opIN_CIDR:
		if length == 2 { //nolint:mnd
			return fmt.Sprintf("ipIdx=%v /%v", args[0], args[1])
		}

	case opBOOST, opFUZZ:
		if length == 1 {
			return fmt.Sprintf("%v", args[0])
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

//nolint:cyclop,gocognit,gocyclo,exhaustive,funlen
func renderComment(prog *Program, isn Instr, lbls map[int]string) string {
	args := isn.Args
	length := len(args)

	switch isn.Op {
	case opLDA:
		if length == 1 {
			return fmt.Sprintf(`set accumulator=%v`, args[0])
		}

	case opSET_F:
		if length == 1 {
			if fid, ok := args[0].(int); ok {
				return fmt.Sprintf(
					`set field="%s"`,
					prog.fieldName(fid))
			}
		}

	case opJMP, opJZ, opJNZ:
		if length == 1 {
			if tgt, ok := args[0].(int); ok {
				if name, has := lbls[tgt]; has {
					return "jump " + name
				}

				return fmt.Sprintf(
					"jump %d",
					tgt)
			}
		}

	case opEQ_S, opNEQ_S, opPREFIX, opGLOB:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return prog.stringLiteral(idx)
			}
		}

	case opREGEX:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return prog.regexLiteral(idx)
			}
		}

	case opPHRASE:
		if length == 2 { //nolint:mnd
			idx, _ := args[0].(int)
			prox, _ := args[1].(int)

			if prox == 0 {
				return "phrase=" + prog.stringLiteral(idx)
			}

			return fmt.Sprintf("phrase=%s prox=%d",
				prog.stringLiteral(idx),
				prox)
		}

	case opEQ_N, opNEQ_N, opLT_N, opLTE_N, opGT_N, opGTE_N:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return prog.numberLiteral(idx)
			}
		}

	case opEQ_T, opNEQ_T, opLT_T, opLTE_T, opGT_T, opGTE_T:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return prog.timeLiteral(idx)
			}
		}

	case opEQ_IP, opNEQ_IP, opLT_IP, opLTE_IP, opGT_IP, opGTE_IP:
		if length == 1 {
			if idx, ok := args[0].(int); ok {
				return prog.ipLiteral(idx)
			}
		}

	case opRANGE_N:
		if length == 4 { //nolint:mnd
			return prog.ppRange(args, PPRangeN)
		}

	case opRANGE_T:
		if length == 4 { //nolint:mnd
			return prog.ppRange(args, PPRangeT)
		}

	case opRANGE_IP:
		if length == 4 { //nolint:mnd
			return prog.ppRange(args, PPRangeIP)
		}

	case opIN_CIDR:
		if length == 2 { //nolint:mnd
			ipIdx, _ := args[0].(int)
			pfx, _ := args[1].(int)

			return fmt.Sprintf("%s/%d", prog.ipLiteral(ipIdx), pfx)
		}

	case opBOOST:
		if length == 1 {
			if val, ok := args[0].(float64); ok {
				return fmt.Sprintf("boost=%.3f", val)
			}
		}

	case opFUZZ:
		if length == 1 {
			if val, ok := args[0].(float64); ok {
				return fmt.Sprintf("fuzz=%.3f", val)
			}
		}
	}

	return ""
}

//nolint:exhaustive
func computeLabels(code []Instr) map[int]string {
	lbls := map[int]string{}
	seq := 1

	maybeLabel := func(loc int) {
		if loc < 0 || loc >= len(code) {
			return
		}

		if _, ok := lbls[loc]; ok {
			return
		}

		lbls[loc] = fmt.Sprintf("L%03d", seq)
		seq++
	}

	for _, isn := range code {
		switch isn.Op {
		case opJMP, opJZ, opJNZ:
			if len(isn.Args) == 0 {
				continue
			}

			if tgt, ok := isn.Args[0].(int); ok {
				maybeLabel(tgt)
			}
		}
	}

	return lbls
}

func NewDefaultPrettyPrinterOptions() PrettyPrinterOpts {
	return PrettyPrinterOpts{
		WithComments: true,
		AddrWidth:    0,
		OpcodeWidth:  0,
		OperandWidth: 0,
	}
}

// * disassembler.go ends here.
