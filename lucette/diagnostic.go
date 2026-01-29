// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// diagnostic.go --- Diagnostic type.
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

import "strings"

// * Constants:

// * Variables:

// * Code:

// ** Structure:

type Diagnostic struct {
	Msg  string // Diagnostic message.
	At   *Span  // Location within source code.
	Hint string // Hint message, if applicable.
}

// ** Methods:

// Return the string representation of a diagnostic.
func (d Diagnostic) String() string {
	var sbld strings.Builder

	sbld.WriteString(d.At.String())
	sbld.WriteString(": ")
	sbld.WriteString(d.Msg)

	if len(d.Hint) > 0 {
		sbld.WriteString(" [")
		sbld.WriteString(d.Hint)
		sbld.WriteRune(']')
	}

	return sbld.String()
}

// ** Functions:

func NewDiagnostic(msg string, at *Span) Diagnostic {
	return Diagnostic{Msg: msg, At: at}
}

func NewDiagnosticHint(msg, hint string, at *Span) Diagnostic {
	return Diagnostic{Msg: msg, Hint: hint, At: at}
}

// * diagnostic.go ends here.
