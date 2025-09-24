// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// validator.go --- Go structure validator.
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

package validator

// * Imports:

import (
	"context"
	"io"

	"github.com/Asmodai/gohacks/dag"
)

// * Code:

// ** Type:

// Validator structure.
type Validator struct {
	cmplr dag.Compiler // The DAG compiler.
	act   *actions     // Actions.
}

// ** Methods:

// Compile an action from an action specification.
func (v *Validator) CompileAction(spec dag.ActionSpec) (dag.ActionFn, error) {
	return v.cmplr.CompileAction(spec)
}

// Compile a failure action from a failure specification.
func (v *Validator) CompileFailure(spec dag.FailureSpec) (dag.ActionFn, error) {
	return v.cmplr.CompileFailure(spec)
}

// Compile a slice of rule specs into a DAG graph.
func (v *Validator) Compile(specs []dag.RuleSpec) []error {
	return v.cmplr.Compile(specs)
}

// Evaluate an input against the validator.
func (v *Validator) Evaluate(input dag.Filterable) {
	v.cmplr.Evaluate(input)
}

// Export the compiler's rulesets to GraphViz DOT format.
func (v *Validator) Export(writer io.Writer) {
	v.cmplr.Export(writer)
}

// Return a list of failure messages (if any) generated during validation.
func (v *Validator) Failures() []error {
	return v.act.errors
}

// Clear the list of failure messages.
func (v *Validator) ClearFailures() {
	v.act.errors = v.act.errors[:0]
}

// ** Functions:

// Create a new validator with the default action set and predicate list.
func NewValidator(ctx context.Context) *Validator {
	inst := &Validator{
		act: &actions{
			errors: []error{},
		},
	}

	dag := dag.NewCompilerWithPredicates(
		ctx,
		inst.act,
		BuildPredicateDict(),
	)

	inst.cmplr = dag

	return inst
}

// * validator.go ends here.
