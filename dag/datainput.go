// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// datainput.go --- DAG input structure.
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

import "fmt"

// * Constants:

// * Variables:

// * Code:

// ** Types:

type DataInput struct {
	fields map[string]any
}

// ** Methods:

func (input *DataInput) Get(key string) (any, bool) {
	val, ok := input.fields[key]
	if !ok {
		return nil, false
	}

	return val, true
}

func (input *DataInput) Keys() []string {
	result := make([]string, 0, len(input.fields))

	for key := range input.fields {
		result = append(result, key)
	}

	return result
}

func (input *DataInput) Set(key string, value any) bool {
	_, found := input.fields[key]
	if !found {
		return false
	}

	input.fields[key] = value

	return true
}

func (input *DataInput) String() string {
	return fmt.Sprintf("%v", input.fields)
}

// ** Functions

func NewDataInput() *DataInput {
	return &DataInput{fields: map[string]any{}}
}

// Create a new `DataInput` object with a copy of the provided input map.
func NewDataInputFromMap(input map[string]any) *DataInput {
	copyMap := make(map[string]any, len(input))

	for key, val := range input {
		copyMap[key] = val
	}

	return &DataInput{fields: copyMap}
}

// * datainput.go ends here.
