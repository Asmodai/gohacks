// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// utils.go --- Utilities.
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
	"encoding/json"

	"gitlab.com/tozd/go/errors"
	"gopkg.in/yaml.v3"
)

// * Constants:

const (
	prefixJSON = ""
	indentJSON = "    "
)

// * Variables:

// * Code:

// ** Methods:

// Dump the rule specification to YAML format.
func (rs *RuleSpec) DumpToYAML() (string, error) {
	raw, err := yaml.Marshal(rs)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(raw), nil
}

// Dump the rule specification to JSON format.
func (rs *RuleSpec) DumpToJSON() (string, error) {
	raw, err := json.MarshalIndent(rs, prefixJSON, indentJSON)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(raw), nil
}

// Dump a slice of rule specifications to YAML format.
func DumpRulesToYAML(rules []RuleSpec) (string, error) {
	raw, err := yaml.Marshal(rules)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(raw), nil
}

// Dump a slice of rule specifications to JSON format.
func ParseFromYAML(data string) ([]RuleSpec, error) {
	result := []RuleSpec{}
	raw := []byte(data)

	if err := yaml.Unmarshal(raw, &result); err != nil {
		return []RuleSpec{}, errors.WithStack(err)
	}

	return result, nil
}

func ParseFromJSON(data string) ([]RuleSpec, error) {
	result := []RuleSpec{}
	raw := []byte(data)

	if err := json.Unmarshal(raw, &result); err != nil {
		return []RuleSpec{}, errors.WithStack(err)
	}

	return result, nil
}

// * utils.go ends here.
