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

// * Comments:

// * Package:

package lucette

// * Imports:

// * Constants:

// * Variables:

// * Code:

// ** Functions:

func nnfAnd(node *TypedNodeAnd, neg bool) TypedNode {
	kids := make([]TypedNode, 0, len(node.kids))

	for _, kid := range node.kids {
		kids = append(kids, nnf(kid, neg))
	}

	if neg {
		// De Morgan: NOT (A && B) => (NOT A) || (NOT B)
		return &TypedNodeOr{kids: kids}
	}

	return &TypedNodeAnd{kids: kids}
}

func nnfOr(node *TypedNodeOr, neg bool) TypedNode {
	kids := make([]TypedNode, 0, len(node.kids))

	for _, kid := range node.kids {
		kids = append(kids, nnf(kid, neg))
	}

	if neg {
		// NOT (A || B) => (NOT A) && (NOT B)
		return &TypedNodeAnd{kids: kids}
	}

	return &TypedNodeOr{kids: kids}
}

func nnf(node TypedNode, neg bool) TypedNode {
	switch val := node.(type) {
	case *TypedNodeAnd:
		return nnfAnd(val, neg)

	case *TypedNodeOr:
		return nnfOr(val, neg)

	case *TypedNodeNot:
		return nnf(val.kid, !neg)

	default:
		if !neg {
			return val
		}

		// Try to invert the leaf.
		if inv := invertLeaf(val); inv != nil {
			return inv
		}

		// Fallback: keep a single NOT as the leaf.
		return &TypedNodeNot{kid: val}
	}
}

func invertRangeN(node *TypedNodeRangeN) TypedNode {
	var parts []TypedNode

	if node.low != nil {
		parts = append(
			parts,
			&TypedNodeCmpN{
				field: node.field,
				op:    CmpLT,
				value: *node.low})

		//nolint:forcetypeassert
		if node.incl {
			parts[len(parts)-1].(*TypedNodeCmpN).op = CmpLT
		} else {
			parts[len(parts)-1].(*TypedNodeCmpN).op = CmpLTE
		}
	}

	if node.high != nil {
		opcode := CmpGTE

		if node.inch {
			opcode = CmpGT
		}

		parts = append(
			parts,
			&TypedNodeCmpN{
				field: node.field,
				op:    opcode,
				value: *node.high})
	}

	switch len(parts) {
	case 0:
		return &TypedNodePhrase{
			field:  node.field,
			phrase: "",
			prox:   0}

	case 1:
		return parts[0]
	}

	return &TypedNodeOr{kids: parts}
}

func invertRangeT(node *TypedNodeRangeT) TypedNode {
	var parts []TypedNode

	if node.low != nil {
		opcode := CmpLTE

		if node.incl {
			opcode = CmpLT
		}

		parts = append(
			parts,
			&TypedNodeCmpT{
				field: node.field,
				op:    opcode,
				value: *node.low})
	}

	if node.high != nil {
		opcode := CmpGTE

		if node.inch {
			opcode = CmpGT
		}

		parts = append(
			parts,
			&TypedNodeCmpT{
				field: node.field,
				op:    opcode,
				value: *node.high})
	}

	switch len(parts) {
	case 0:
		return &TypedNodeFalse{}

	case 1:
		return parts[0]
	}

	return &TypedNodeOr{kids: parts}
}

func invertRangeIP(node *TypedNodeRangeIP) TypedNode {
	var parts []TypedNode

	if node.low != zeroIP {
		opcode := CmpLTE

		if node.incl {
			opcode = CmpLT
		}

		parts = append(
			parts,
			&TypedNodeCmpIP{
				field: node.field,
				op:    opcode,
				value: node.low})
	}

	if node.high != zeroIP {
		opcode := CmpGTE

		if node.inch {
			opcode = CmpGT
		}

		parts = append(
			parts,
			&TypedNodeCmpIP{
				field: node.field,
				op:    opcode,
				value: node.high})
	}

	switch len(parts) {
	case 0:
		return &TypedNodeFalse{}
	case 1:
		return parts[0]
	}

	return &TypedNodeOr{kids: parts}
}

func invertLeaf(node TypedNode) TypedNode {
	switch val := node.(type) {
	case *TypedNodeEqS:
		return &TypedNodeNeqS{field: val.field, value: val.value}

	case *TypedNodeNeqS:
		return &TypedNodeEqS{field: val.field, value: val.value}

	case *TypedNodeCmpN:
		return &TypedNodeCmpN{
			field: val.field,
			op:    invertCmp(val.op),
			value: val.value}

	case *TypedNodeCmpT:
		return &TypedNodeCmpT{
			field: val.field,
			op:    invertCmp(val.op),
			value: val.value}

	case *TypedNodeCmpIP:
		return &TypedNodeCmpIP{
			field: val.field,
			op:    invertCmp(val.op),
			value: val.value}

	case *TypedNodeRangeN:
		return invertRangeN(val)

	case *TypedNodeRangeT:
		return invertRangeT(val)

	case *TypedNodeRangeIP:
		return invertRangeIP(val)
	}

	return nil
}

func invertCmp(opcode CmpKind) CmpKind {
	switch opcode {
	case CmpLT:
		return CmpGTE

	case CmpLTE:
		return CmpGT

	case CmpGT:
		return CmpLTE

	case CmpGTE:
		return CmpLT

	case CmpEQ:
		return CmpNEQ

	case CmpNEQ:
		return CmpEQ

	default:
		return opcode
	}
}

func ToNNF(node TypedNode) TypedNode {
	return nnf(node, false)
}

// * nnf.go ends here.
