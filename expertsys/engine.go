// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// engine.go --- Expert system engine.
//
// Copyright (c) 2026 Paul Ward <paul@lisphacker.uk>
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

package expertsys

// * Imports:

import (
	"github.com/Asmodai/gohacks/dag"
	"github.com/Asmodai/gohacks/errx"
)

// * Constants:

// * Variables:

var (
	ErrNotStable = errx.Base("engine did not stabilise")
)

// * Code:

// ** Type:

// Expert system engine.
type Engine struct {
	cmplr dag.Compiler
}

func (e *Engine) RunToFixpoint(wmem WorkingMemory, maxIters int) (int, error) {
	for iter := range maxIters {
		before := wmem.Version()

		e.cmplr.Evaluate(wmem)

		if wmem.Version() == before {
			return iter + 1, nil
		}
	}

	return 0, errx.WithMessagef(
		ErrNotStable,
		"version stuck at %d after %d iterations",
		wmem.Version(),
		maxIters)
}

// * engine.go ends here.
