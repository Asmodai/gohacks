// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// nnf.go --- Negation Normal Form processor.
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
//
//mock:yes

// * Comments:

// * Package:

package lucette

// * Code:

// ** Interface:

type NNF interface {
	NNF(IRNode) IRNode
}

// ** Structure:

type nnf struct {
}

// ** Functions:

func (n *nnf) nnfAnd(node *IRAnd, neg bool) IRNode {
	kids := make([]IRNode, 0, len(node.Kids))

	for _, kid := range node.Kids {
		kids = append(kids, n.nnf(kid, neg))
	}

	if neg {
		// De Morgan: NOT (A && B) => (NOT A) || (NOT B)
		return &IROr{Kids: kids}
	}

	return &IRAnd{Kids: kids}
}

func (n *nnf) nnfOr(node *IROr, neg bool) IRNode {
	kids := make([]IRNode, 0, len(node.Kids))

	for _, kid := range node.Kids {
		kids = append(kids, n.nnf(kid, neg))
	}

	if neg {
		// NOT (A || B) => (NOT A) && (NOT B)
		return &IRAnd{Kids: kids}
	}

	return &IROr{Kids: kids}
}

func (n *nnf) nnf(node IRNode, neg bool) IRNode {
	switch val := node.(type) {
	case *IRAnd:
		return n.nnfAnd(val, neg)

	case *IROr:
		return n.nnfOr(val, neg)

	case *IRNot:
		return n.nnf(val.Kid, !neg)

	default:
		if !neg {
			return val
		}

		// Try to invert the leaf.
		if inv := n.invertLeaf(val); inv != nil {
			return inv
		}

		// Fallback: keep a single NOT as the leaf.
		return &IRNot{Kid: val}
	}
}

func (n *nnf) invertRangeN(node *IRNumberRange) IRNode {
	var parts []IRNode

	if node.Lo != nil {
		parts = append(
			parts,
			&IRNumberCmp{
				Field: node.Field,
				Op:    ComparatorLT,
				Value: *node.Lo})

		//nolint:forcetypeassert
		if node.IncL {
			parts[len(parts)-1].(*IRNumberCmp).Op = ComparatorLT
		} else {
			parts[len(parts)-1].(*IRNumberCmp).Op = ComparatorLTE
		}
	}

	if node.Hi != nil {
		opcode := ComparatorGTE

		if node.IncH {
			opcode = ComparatorGT
		}

		parts = append(
			parts,
			&IRNumberCmp{
				Field: node.Field,
				Op:    opcode,
				Value: *node.Hi})
	}

	switch len(parts) {
	case 0:
		return &IRPhrase{
			Field:     node.Field,
			Phrase:    "",
			Proximity: 0}

	case 1:
		return parts[0]
	}

	return &IROr{Kids: parts}
}

func (n *nnf) invertRangeT(node *IRTimeRange) IRNode {
	var parts []IRNode

	if node.Lo != nil {
		opcode := ComparatorLTE

		if node.IncL {
			opcode = ComparatorLT
		}

		parts = append(
			parts,
			&IRTimeCmp{
				Field: node.Field,
				Op:    opcode,
				Value: *node.Lo})
	}

	if node.Hi != nil {
		opcode := ComparatorGTE

		if node.IncH {
			opcode = ComparatorGT
		}

		parts = append(
			parts,
			&IRTimeCmp{
				Field: node.Field,
				Op:    opcode,
				Value: *node.Hi})
	}

	switch len(parts) {
	case 0:
		return &IRFalse{}

	case 1:
		return parts[0]
	}

	return &IROr{Kids: parts}
}

func (n *nnf) invertRangeIP(node *IRIPRange) IRNode {
	var parts []IRNode

	if node.Lo != zeroIP {
		opcode := ComparatorLTE

		if node.IncL {
			opcode = ComparatorLT
		}

		parts = append(
			parts,
			&IRIPCmp{
				Field: node.Field,
				Op:    opcode,
				Value: node.Lo})
	}

	if node.Hi != zeroIP {
		opcode := ComparatorGTE

		if node.IncH {
			opcode = ComparatorGT
		}

		parts = append(
			parts,
			&IRIPCmp{
				Field: node.Field,
				Op:    opcode,
				Value: node.Hi})
	}

	switch len(parts) {
	case 0:
		return &IRFalse{}
	case 1:
		return parts[0]
	}

	return &IROr{Kids: parts}
}

func (n *nnf) invertLeaf(node IRNode) IRNode {
	switch val := node.(type) {
	case *IRStringEQ:
		return &IRStringNEQ{Field: val.Field, Value: val.Value}

	case *IRStringNEQ:
		return &IRStringEQ{Field: val.Field, Value: val.Value}

	case *IRNumberCmp:
		return &IRNumberCmp{
			Field: val.Field,
			Op:    InvertComparator(val.Op),
			Value: val.Value}

	case *IRTimeCmp:
		return &IRTimeCmp{
			Field: val.Field,
			Op:    InvertComparator(val.Op),
			Value: val.Value}

	case *IRIPCmp:
		return &IRIPCmp{
			Field: val.Field,
			Op:    InvertComparator(val.Op),
			Value: val.Value}

	case *IRNumberRange:
		return n.invertRangeN(val)

	case *IRTimeRange:
		return n.invertRangeT(val)

	case *IRIPRange:
		return n.invertRangeIP(val)
	}

	return nil
}

func (n *nnf) NNF(tokens IRNode) IRNode {
	return n.nnf(tokens, false)
}

// ** Functions:

func NewNNF() NNF {
	return &nnf{}
}

// * nnf.go ends here.
