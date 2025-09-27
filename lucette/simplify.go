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
//
//mock:yes

// * Comments:

// * Package:

package lucette

// * Constants:

const (
	BooleAnd BooleOp = iota
	BooleOr
)

// * Code:

// ** Types:

// Boolean operation type.
type BooleOp int

type buildFn func([]IRNode) IRNode

// ** Interface:

type Simplifier interface {
	Simplify(IRNode) IRNode
}

// ** Structure:

type simplifier struct {
}

// ** Methods:

func (s *simplifier) dedupeFlat(boolop BooleOp, flat []IRNode, build buildFn) IRNode {
	uniq := make([]IRNode, 0, len(flat))
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
		if boolop == BooleAnd {
			return &IRTrue{}
		}

		return &IRFalse{}

	case 1:
		return uniq[0]
	}

	return build(uniq)
}

func (s *simplifier) simplifyFlat(boolop BooleOp, flat []IRNode, build buildFn) IRNode {
	pos := make(map[string]bool, len(flat))
	neg := make(map[string]bool, len(flat))

	for _, kid := range flat {
		if node, isNot := kid.(*IRNot); isNot {
			key := node.Key()

			if pos[key] {
				switch boolop {
				case BooleOr:
					return &IRTrue{}
				case BooleAnd:
					return &IRFalse{}
				}
			}

			neg[key] = true
		} else {
			key := kid.Key()

			if neg[key] {
				switch boolop {
				case BooleOr:
					return &IRTrue{}
				case BooleAnd:
					return &IRFalse{}
				}
			}

			pos[key] = true
		}
	}

	return s.dedupeFlat(boolop, flat, build)
}

func (s *simplifier) simplifyAnd(node *IRAnd) IRNode {
	flat := make([]IRNode, 0, len(node.Kids))

	for _, kid := range node.Kids {
		simp := s.Simplify(kid)

		switch val := simp.(type) {
		case *IRTrue:
		// Drop.

		case *IRAnd:
			flat = append(flat, val.Kids...)

		case *IRFalse:
			return val

		default:
			flat = append(flat, val)
		}
	}

	return s.simplifyFlat(BooleAnd, flat, func(uniq []IRNode) IRNode {
		return &IRAnd{Kids: uniq}
	})
}

func (s *simplifier) simplifyOr(node *IROr) IRNode {
	flat := make([]IRNode, 0, len(node.Kids))

	for _, kid := range node.Kids {
		simp := s.Simplify(kid)

		switch val := simp.(type) {
		case *IRFalse:
		// Drop.

		case *IROr:
			flat = append(flat, val.Kids...)

		case *IRTrue:
			return val

		default:
			flat = append(flat, val)
		}
	}

	return s.simplifyFlat(BooleOr, flat, func(uniq []IRNode) IRNode {
		return &IROr{Kids: uniq}
	})
}

func (s *simplifier) simplifyNot(node *IRNot) IRNode {
	kid := s.Simplify(node.Kid)

	switch kid.(type) {
	case *IRTrue:
		return &IRFalse{}

	case *IRFalse:
		return &IRTrue{}
	}

	return &IRNot{Kid: kid}
}

//nolint:unused
func (s *simplifier) simplifyPhrase(_ *IRPhrase) IRNode {
	// if node.Proximity == 0 || node.HasWildcard() {
	// TODO: This needs to be investigated.
	// }
	return nil
}

func (s *simplifier) Simplify(node IRNode) IRNode {
	switch val := node.(type) {
	case *IRAnd:
		return s.simplifyAnd(val)

	case *IROr:
		return s.simplifyOr(val)

	case *IRNot:
		return s.simplifyNot(val)

	default:
		return node
	}
}

// ** Functions:

func NewSimplifier() Simplifier {
	return &simplifier{}
}

// * simplify.go ends here.
