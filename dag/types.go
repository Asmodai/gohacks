// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// types.go --- Base data types.
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

import "context"

// * Code:

// ** Concrete types:

// Base data type.
//
// This is a map of key/value pairs.
type DataMap map[string]any

// Action parameters type.
//
// A map of key/value pairs that is passed to the action handler.
type ActionParams map[string]any

// Predicate function type.
//
// A predicate is a function that answers a yes-or-no question.  In other
// words: any expression that can boil down to a boolean.
type PredicateFn func(string, any) Predicate

// Action function callback type.
//
// An action callback is a function that takes a single argument containing
// the key/value pair map and returns no value.
type ActionFn func(context.Context, DataMap)

// ** Abstract data types:

// Filter rule specification.
type RuleSpec struct {
	// Rule name.
	Name string `json:"name" yaml:"name"`

	// List of conditions.
	Conditions []ConditionSpec `json:"conditions" yaml:"conditions"`

	// Action to evaluate.
	Action ActionSpec `json:"action" yaml:"action"`
}

// Condition specification.
type ConditionSpec struct {
	// Attribute to check.
	Attribute string `json:"attribute" yaml:"attribute"`

	// Predicate operator.
	Operator string `json:"operator" yaml:"operator"`

	// Value to check.
	Value any `json:"value" yaml:"value"`
}

// Action specification.
type ActionSpec struct {
	// Action name.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Function to perform.
	Perform string `json:"perform,omitempty" yaml:"perform,omitempty"`

	// Parameters.
	Params ActionParams `json:"params,omitempty" yaml:"params,omitempty"`
}

// * types.go ends here.
