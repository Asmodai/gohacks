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

//
//
//

// * Package:

package dag

// * Imports:

import (
	"context"
	"fmt"

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
	Compile([]RuleSpec) []error
	Evaluate(DataMap)
}

// ** Types:

type compiler struct {
	nodeCache map[string]*node    // Noe cache.
	actions   map[string]ActionFn // Action cache.
	root      *node               // Root node.
	ctx       context.Context     // Owning context.
	lgr       logger.Logger       // Logger instance from DI.
	debugMode bool                // Are we debugging?
	builder   Actions             // Action builder.
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

// Compile a rule spec into a DAG graph.
func (cmplr *compiler) Compile(rules []RuleSpec) []error {
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
	cmplr.root = &node{Predicate: alwaysTruePredicate()}
	cmplr.nodeCache = make(map[string]*node)
}

func (cmplr *compiler) compileRule(rule RuleSpec) error {
	current := cmplr.root

	for _, cond := range rule.Conditions {
		pred, key, err := cmplr.buildPredicate(cond)
		if err != nil {
			return err
		}

		next := cmplr.getOrCreateNode(pred, key)
		cmplr.linkNodes(current, next)

		current = next
	}

	return cmplr.attachAction(current, rule.Action)
}

func (cmplr *compiler) buildPredicate(cond ConditionSpec) (Predicate, string, error) {
	builder, ok := predicateBuilders[cond.Operator]
	if !ok {
		return Predicate{}, "", errors.WithMessagef(
			ErrUnknownOperator,
			"unknown operator: %q",
			cond.Operator,
		)
	}

	pred := builder(cond.Attribute, cond.Value)
	key := fmt.Sprintf(
		"%s %s %v",
		cond.Attribute,
		cond.Operator,
		cond.Value,
	)

	return pred, key, nil
}

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

func (cmplr *compiler) attachAction(current *node, action ActionSpec) error {
	compiled, err := cmplr.CompileAction(action)
	if err != nil {
		return errors.WithMessagef(
			ErrRuleCompileFailed,
			"rule '%s' compilation failed: %s",
			action.Name,
			err.Error(),
		)
	}

	current.Action = compiled

	return nil
}

func (cmplr *compiler) Evaluate(input DataMap) {
	traverse(cmplr.ctx, cmplr.root, input)
}

// ** Functions:

func alwaysTruePredicate() Predicate {
	return Predicate{
		Eval: func(_ *Predicate, _ DataMap) bool {
			return true
		},
	}
}

func NewCompiler(ctx context.Context, builder Actions) Compiler {
	lgr := logger.MustGetLogger(ctx)

	dbg, err := contextdi.GetDebugMode(ctx)
	if err != nil {
		dbg = false
	}

	return &compiler{
		nodeCache: make(map[string]*node, initialNodeCacheSize),
		actions:   make(map[string]ActionFn, initialActionsSize),
		ctx:       ctx,
		lgr:       lgr,
		debugMode: dbg,
		builder:   builder,
	}
}

// * compiler.go ends here.
