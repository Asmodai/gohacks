// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// peephole.go --- Peephole optimiser.
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

// * Constants:

// * Variables:

// * Code:

//nolint:cyclop,exhaustive,forcetypeassert
func (p *Program) Peephole() {
	oldCode := p.Code
	newCode := make([]Instr, 0, len(oldCode))
	lastfid := -1
	idx := 0

	for idx < len(oldCode) {
		isn := oldCode[idx]

		switch isn.Op {
		case opLABEL:
			// Keep labels as-is.
			newCode = append(newCode, isn)
			idx++

			continue

		case opSET_F:
			// Remove unchanged field sets.
			fid, _ := isn.Args[0].(int)

			if fid == lastfid {
				idx++

				continue
			}

			lastfid = fid

		case opNOT:
			// Collapse double NOT.
			if idx+1 < len(oldCode) && oldCode[idx+1].Op == opNOT {
				idx++

				continue
			}

		case opRANGE_N:
			// Replace RANGE(lo==hi,inclusives) -> EQ
			low, high := isn.Args[0].(int), isn.Args[1].(int)
			incl, inch := isn.Args[2].(bool), isn.Args[3].(bool)

			if low != -1 && high != -1 && low == high && incl && inch {
				isn = Instr{Op: opEQ_N, Args: []any{low}}
			}

		case opJMP, opJZ, opJNZ:
			// Drop jumps to immediate next label.
			if idx+1 < len(oldCode) && oldCode[idx+1].Op == opLABEL && len(isn.Args) == 1 {
				if tgt, ok := isn.Args[0].(LabelID); ok {
					if lbl, _ := oldCode[idx+1].Args[0].(LabelID); lbl == tgt {
						idx++

						continue
					}
				}
			}
		}

		// Keep re-written or unchanged instructions.
		newCode = append(newCode, isn)

		idx++
	}

	p.Code = newCode
}

// * peephole.go ends here.
