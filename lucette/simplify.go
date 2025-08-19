// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// simplify.go --- Simplifier.
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

const (
	OpAnd BoolOp = iota
	OpOr
)

// * Variables:

// * Code:

type BoolOp int

type buildFn func([]TypedNode) TypedNode

func dedupeFlat(opcode BoolOp, flat []TypedNode, build buildFn) TypedNode {
	uniq := make([]TypedNode, 0, len(flat))
	seen := make(map[string]struct{}, len(flat))

	for _, kid := range flat {
		key := kid.Key()

		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}

		uniq = append(uniq, kid)
	}

	switch len(uniq) {
	case 0:
		if opcode == OpAnd {
			return &TypedNodeTrue{}
		}

		return &TypedNodeFalse{}

	case 1:
		return uniq[0]
	}

	return build(uniq)
}

func simplifyFlat(opcode BoolOp, flat []TypedNode, build buildFn) TypedNode {
	pos := make(map[string]bool, len(flat))
	neg := make(map[string]bool, len(flat))

	for _, kid := range flat {
		if node, isNot := kid.(*TypedNodeNot); isNot {
			key := node.Key()

			if pos[key] {
				switch opcode {
				case OpOr:
					return &TypedNodeTrue{}
				case OpAnd:
					return &TypedNodeFalse{}
				}
			}

			neg[key] = true
		} else {
			key := kid.Key()

			if neg[key] {
				switch opcode {
				case OpOr:
					return &TypedNodeTrue{}
				case OpAnd:
					return &TypedNodeFalse{}
				}
			}

			pos[key] = true
		}
	}

	return dedupeFlat(opcode, flat, build)
}

func simplifyAnd(node *TypedNodeAnd) TypedNode {
	flat := make([]TypedNode, 0, len(node.kids))

	for _, kid := range node.kids {
		simp := Simplify(kid)

		switch val := simp.(type) {
		case *TypedNodeTrue:
			// Drop.

		case *TypedNodeAnd:
			flat = append(flat, val.kids...)

		case *TypedNodeFalse:
			return val

		default:
			flat = append(flat, val)
		}
	}

	return simplifyFlat(OpAnd, flat, func(uniq []TypedNode) TypedNode {
		return &TypedNodeAnd{kids: uniq}
	})
}

func simplifyOr(node *TypedNodeOr) TypedNode {
	flat := make([]TypedNode, 0, len(node.kids))

	for _, kid := range node.kids {
		simp := Simplify(kid)

		switch val := simp.(type) {
		case *TypedNodeFalse:
			// Drop.

		case *TypedNodeOr:
			flat = append(flat, val.kids...)

		case *TypedNodeTrue:
			return val

		default:
			flat = append(flat, val)
		}
	}

	return simplifyFlat(OpOr, flat, func(uniq []TypedNode) TypedNode {
		return &TypedNodeOr{kids: uniq}
	})
}

func simplifyNot(node *TypedNodeNot) TypedNode {
	kid := Simplify(node.kid)

	switch kid.(type) {
	case *TypedNodeTrue:
		return &TypedNodeFalse{}
	case *TypedNodeFalse:
		return &TypedNodeTrue{}
	}

	return &TypedNodeNot{kid: kid}
}

func Simplify(node TypedNode) TypedNode {
	switch val := node.(type) {
	case *TypedNodeAnd:
		return simplifyAnd(val)

	case *TypedNodeOr:
		return simplifyOr(val)

	case *TypedNodeNot:
		return simplifyNot(val)

	default:
		return node
	}
}

// * simplify.go ends here.
