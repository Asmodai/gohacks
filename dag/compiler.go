// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// compiler.go --- DAG compiler.
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

package dag

// * Imports:

import (
	"context"
	"fmt"
	"io"

	"github.com/Asmodai/gohacks/contextdi"
	"github.com/Asmodai/gohacks/logger"
	"gitlab.com/tozd/go/errors"
)

// * Constants:

const (
	initialNodeCacheSize = 2
	initialActionsSize   = 2
)

// * Variables:

var (
	ErrUnknownOperator   = errors.Base("unknown operator")
	ErrRuleCompileFailed = errors.Base("rule compilation failed")
)

// * Code:

// ** Interface:

type Compiler interface {
	CompileAction(ActionSpec) (ActionFn, error)
	CompileFailure(FailureSpec) (ActionFn, error)
	Compile([]RuleSpec) []error
	Evaluate(Filterable)
	Export(io.Writer)
}

// ** Types:

type compiler struct {
	ctx        context.Context     // Owning context.
	lgr        logger.Logger       // Logger instance from DI.
	builder    Actions             // Action builder.
	nodeCache  map[string]*node    // Node cache.
	actions    map[string]ActionFn // Action cache.
	predicates PredicateDict       // Predicates.
	root       *node               // Root node.
	debugMode  bool                // Are we debugging?
}

// ** Methods:

// Check to see if a key exists in our node cache.
//
// If it does, return the associated node.
//
// If it does not, then create a new node and return that.
//
// If anything else happens, please call 1-800-OMGWTFLUL.
func (cmplr *compiler) getOrCreateNode(pred Predicate, key string) *node {
	if oldNode, ok := cmplr.nodeCache[key]; ok {
		return oldNode
	}

	newNode := &node{Predicate: pred}
	cmplr.nodeCache[key] = newNode

	return newNode
}

// Compile an action from an action specification.
func (cmplr *compiler) CompileAction(spec ActionSpec) (ActionFn, error) {
	built, err := cmplr.builder.Builder(spec.Perform, spec.Params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return built, nil
}

// Compile a failure action from a failure specification.
func (cmplr *compiler) CompileFailure(spec FailureSpec) (ActionFn, error) {
	built, err := cmplr.builder.Builder(spec.Perform, spec.Params)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return built, nil
}

// Initialise predicate dictionary.
func (cmplr *compiler) initPredicates() {
	if cmplr.predicates == nil {
		cmplr.predicates = BuildPredicateDict()
	}
}

// Compile a slice of rule specs into a DAG graph.
func (cmplr *compiler) Compile(rules []RuleSpec) []error {
	if cmplr.predicates == nil {
		cmplr.initPredicates()
	}

	cmplr.initRoot()

	issues := make([]error, 0)

	for _, rule := range rules {
		if err := cmplr.compileRule(rule); err != nil {
			issues = append(issues, err)
		}
	}

	return issues
}

// Build the root node with a dummy predicate.
func (cmplr *compiler) initRoot() {
	cmplr.root = &node{Predicate: &NOOPPredicate{}}
	cmplr.nodeCache = make(map[string]*node)
}

// Compile a rule specification.
func (cmplr *compiler) compileRule(rule RuleSpec) error {
	current := cmplr.root

	for _, cond := range rule.Conditions {
		pred, key, err := cmplr.buildPredicate(cond)
		if err != nil {
			return errors.WithMessagef(
				err,
				"Rule %q",
				rule.Name)
		}

		next := cmplr.getOrCreateNode(pred, key)
		cmplr.linkNodes(current, next)

		current = next
	}

	if err := cmplr.attachAction(current, rule.Action); err != nil {
		return errors.WithMessagef(
			err,
			"Rule %q",
			rule.Name)
	}

	if err := cmplr.attachFailure(current, rule.Failure); err != nil {
		return errors.WithMessagef(
			err,
			"Rule %q",
			rule.Name)
	}

	return nil
}

// Build a predicate from a condition specification.
func (cmplr *compiler) buildPredicate(cond ConditionSpec) (Predicate, string, error) {
	builder, ok := cmplr.predicates[cond.Operator]
	if !ok {
		return nil, "", errors.WithMessagef(
			ErrUnknownOperator,
			"Condition operator %q",
			cond.Operator)
	}

	pred, err := builder.Build(
		cond.Attribute,
		cond.Value,
		cmplr.lgr, cmplr.debugMode,
	)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	key := fmt.Sprintf(
		"%s %s %v",
		cond.Attribute,
		cond.Operator,
		cond.Value,
	)

	return pred, key, nil
}

// Link the given nodes.
func (cmplr *compiler) linkNodes(current, next *node) {
	linked := false

	for _, child := range current.Children {
		if child == next {
			linked = true

			break
		}
	}

	if !linked {
		current.Children = append(current.Children, next)
	}
}

// Attach a success action to the given node.
func (cmplr *compiler) attachAction(current *node, action ActionSpec) error {
	if len(action.Perform) == 0 {
		return nil
	}

	compiled, err := cmplr.CompileAction(action)
	if err != nil {
		return errors.WithMessagef(
			err,
			"Action %q",
			action.Name)
	}

	if len(action.Name) > 0 {
		current.ActionName = action.Name
	} else {
		current.ActionName = action.Perform
	}

	current.Action = compiled

	return nil
}

// Attach a failure action to the given node.
func (cmplr *compiler) attachFailure(current *node, failure FailureSpec) error {
	if len(failure.Perform) == 0 {
		return nil
	}

	compiled, err := cmplr.CompileFailure(failure)
	if err != nil {
		return errors.WithMessagef(
			err,
			"Failure %q",
			failure.Name)
	}

	if len(failure.Name) > 0 {
		current.FailureName = failure.Name
	} else {
		current.FailureName = failure.Perform
	}

	current.Failure = compiled

	return nil
}

// Evaluate an input against the DAG.
func (cmplr *compiler) Evaluate(input Filterable) {
	traverse(cmplr.ctx,
		cmplr.root,
		input,
		cmplr.debugMode,
		cmplr.lgr)
}

// Export the compiler's rulesets to GraphViz DOT format.
func (cmplr *compiler) Export(writer io.Writer) {
	ExportToDOT(writer, cmplr.root)
}

// ** Functions:

// Return a new DAG compiler.
func NewCompiler(ctx context.Context, build Actions) Compiler {
	return NewCompilerWithPredicates(ctx, build, BuildPredicateDict())
}

// Return a new DAG compiler with custom predicates.
func NewCompilerWithPredicates(
	ctx context.Context,
	builder Actions,
	predicates PredicateDict,
) Compiler {
	lgr := logger.MustGetLogger(ctx)

	dbg, err := contextdi.GetDebugMode(ctx)
	if err != nil {
		dbg = false
	}

	return &compiler{
		nodeCache:  make(map[string]*node, initialNodeCacheSize),
		actions:    make(map[string]ActionFn, initialActionsSize),
		ctx:        ctx,
		lgr:        lgr,
		debugMode:  dbg,
		builder:    builder,
		predicates: predicates,
	}
}

// * compiler.go ends here.
