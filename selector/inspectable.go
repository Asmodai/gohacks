// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// inspectable.go --- Inspectable responders.
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
//
//mock:yes

// * Comments:

// * Package:

package selector

// * Imports:

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Asmodai/gohacks/responder"
	"github.com/Asmodai/gohacks/utils"
)

// * Code:

// ** Interfaces:

// Objects that implement these methods are considered `introspectable` and
// are deemed capable of being asked to describe themselves in various ways.
type Introspectable interface {
	responder.Respondable

	// Return a list of selectors.
	Selectors() []string

	// Return a list of sorted selectors
	SortedSelectors() []string

	// Return a list of methods that the object can respond to.
	Methods() *Table

	// Does the object conform to the given protocol?
	ConformsTo(protocol string) bool

	// List all protocols for which the object claims conformity.
	ListProtocols() []string

	// Return a map of metadata for a selector
	MetadataForSelector(string) (map[string]string, error)
}

// ** Functions:

func dumpDocString(doc, indent string) string {
	const width = 80

	var sbld strings.Builder

	words := strings.FieldsFunc(doc, unicode.IsSpace)
	lineLen := 0

	sbld.WriteString(indent)

	for _, word := range words {
		if lineLen+len(word)+1 > width {
			sbld.WriteRune('\n')
			sbld.WriteString(indent)
			sbld.WriteString(word)

			lineLen = len(word)
		} else {
			if lineLen > 0 {
				sbld.WriteRune(' ')

				lineLen++
			}

			sbld.WriteString(word)

			lineLen += len(word)
		}
	}

	sbld.WriteRune('\n')

	return sbld.String()
}

//nolint:unparam
func dumpField(field, content, indent string) string {
	const fieldTitleLen = 10

	var sbld strings.Builder

	sbld.WriteString(indent)
	sbld.WriteString(utils.Pad(field+":", fieldTitleLen))
	sbld.WriteRune(' ')
	sbld.WriteString(content)
	sbld.WriteRune('\n')

	return sbld.String()
}

//nolint:cyclop,funlen
func DumpIntrospectableInfo(obj Introspectable) string {
	const indent = "    "

	var sbld strings.Builder

	fmt.Fprintf(&sbld, "Object: %s (%s)\n", obj.Name(), obj.Type())

	//
	// Print protocols.
	sbld.WriteString("Protocols:\n")

	for _, proto := range obj.ListProtocols() {
		sbld.WriteString("  @ ")
		sbld.WriteString(proto)
		sbld.WriteRune('\n')
	}

	//
	// Print selectors.
	sbld.WriteString("Selectors:\n")

	for _, sel := range obj.SortedSelectors() {
		meta, err := obj.MetadataForSelector(sel)
		if err != nil {
			meta = map[string]string{}
		}

		sbld.WriteString("  - ")
		sbld.WriteString(sel)
		sbld.WriteString("\n\n")

		if strval, found := meta["version"]; found {
			sbld.WriteString(dumpField("Version", strval, indent))
		}

		if strval, found := meta["since"]; found {
			sbld.WriteString(dumpField("Since", strval, indent))
		}

		if strval, found := meta["protocol"]; found {
			sbld.WriteString(dumpField("Protocol", strval, indent))
		}

		if strval, found := meta["visibility"]; found {
			sbld.WriteString(dumpField("Visibility", strval, indent))
		}

		if strval, found := meta["author"]; found {
			sbld.WriteString(dumpField("Author", strval, indent))
		}

		if strval, found := meta["tags"]; found {
			sbld.WriteString(dumpField("Tags", strval, indent))
		}

		if strval, found := meta["deprecated"]; found {
			sbld.WriteString(dumpField("Deprecated", strval, indent))
		}

		if docStr, found := meta["doc"]; found {
			sbld.WriteRune('\n')
			sbld.WriteString(dumpDocString(docStr, indent))
		}
	}

	return sbld.String()
}

// * inspectable.go ends here.
