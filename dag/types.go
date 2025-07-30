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

// Condition specification.
type ConditionSpec struct {
	Attribute string `json:"attribute"` // Attribute to check.
	Operator  string `json:"operator"`  // Predicate operator.
	Value     any    `json:"value"`     // Value to check.
}

// Filter rule specification.
type RuleSpec struct {
	Name       string          `json:"name"`       // Rule name.
	Conditions []ConditionSpec `json:"conditions"` // List of conditions.
	Action     ActionSpec      `json:"action"`     // Action to evaluate.
}

// Action specification.
type ActionSpec struct {
	Name    string       `json:"name,omitempty"`    // Action name.
	Perform string       `json:"perform,omitempty"` // Function to perform.
	Params  ActionParams `json:"params,omitempty"`  // Parameters.
}

// * types.go ends here.
