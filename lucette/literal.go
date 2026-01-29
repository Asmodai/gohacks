// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// literal.go --- Literal type.
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

package lucette

// * Imports:

import "github.com/Asmodai/gohacks/conversion"

// * Code:

// ** Type:

// Literal value structure.
//
//nolint:errname
type Literal struct {
	Value any   // The literal value.
	Err   error // The error to pass as a literal.
}

// ** Methods:

// Is the literal an error?
func (l Literal) IsError() bool {
	return l.Err != nil
}

// Return the string representation of the literal.
func (l Literal) String() string {
	if l.Err != nil {
		return l.Err.Error()
	}

	if normal, ok := conversion.ToString(l.Value); ok {
		return normal
	}

	return "<invalid>"
}

// Return the error message if one is present.
func (l Literal) Error() string {
	if l.Err == nil {
		return ""
	}

	return l.Err.Error()
}

// ** Functions:

// Create a new value literal.
func NewLiteral(value any) Literal {
	return Literal{
		Value: value,
	}
}

// Create a new error literal.
func NewErrorLiteral(err error) Literal {
	return Literal{
		Err: err,
	}
}

// * literal.go ends here.
