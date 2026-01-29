// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// node.go --- Direct Acyclic Graph node type.
//
// Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
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

package dag

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/logger"
)

// * Code:

// Graph node type.
type node struct {
	Predicate   Predicate // Predicate.
	Action      ActionFn  // Action to execute upon successful predicate.
	ActionName  string    // Action name.
	Failure     ActionFn  // Action to execute upon predicate failure.
	FailureName string    // Failure action name.
	Children    []*node   // Child nodes.
}

// Traverse each child node in the given root node.
//
// If the node has an associated predicate then that is evaluated against
// the given input.
func traverse(ctx context.Context, root *node, input Filterable, debug bool, logger logger.Logger) {
	if !root.Predicate.Eval(ctx, input) {
		if debug {
			logger.Debug(
				"Eval failure",
				"predicate", root.Predicate.Debug(),
				"input", input.String(),
			)
		}

		if root.Failure != nil {
			root.Failure(ctx, input)
		}

		return
	}

	if debug {
		logger.Debug(
			"Eval success",
			"predicate", root.Predicate.Debug(),
			"input", input.String(),
		)
	}

	if root.Action != nil {
		root.Action(ctx, input)
	}

	for _, child := range root.Children {
		traverse(ctx, child, input, debug, logger)
	}
}

// * node.go ends here.
