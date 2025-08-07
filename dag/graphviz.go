// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// graphviz.go --- Graphviz exporter.
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
	"fmt"
	"io"
	"strings"
)

// * Code:

// Generate a Graphviz visualisation of the DAG starting from the given node.
func ExportToDOT(writer io.Writer, root *node) {
	fmt.Fprintln(writer, "digraph DAG {")
	fmt.Fprintln(writer, "  rankdir=TB;")
	fmt.Fprintln(writer, `  node [shape=record, fontname="monospace", fontsize=10];`)

	visited := make(map[*node]string)
	counter := 0

	exportNodeDOT(writer, root, visited, &counter)

	fmt.Fprintln(writer, "}")
}

func exportNodeDOT(writer io.Writer, node *node, visited map[*node]string, counter *int) string {
	if id, ok := visited[node]; ok {
		return id
	}

	// Node ID.
	nodeID := fmt.Sprintf("n%d", *counter)
	*counter++
	visited[node] = nodeID

	// Label parts.
	kind := node.Predicate.Instruction()
	pred := node.Predicate.String()

	parts := []string{
		"<f0> " + kind,
		"<f1> " + escapeDOT(pred),
	}

	color := "#ddeeff"
	if kind == noopIsn {
		color = "#eeeeee"
	}

	if len(node.ActionName) > 0 {
		parts = append(
			parts,
			"<f2> Action: "+node.ActionName,
		)
		color = "#ccffcc"
	}

	// Generate the node.
	fmt.Fprintf(
		writer,
		`  %s [label="{%s}" style=filled fillcolor="%s"];`+"\n",
		nodeID,
		strings.Join(parts, " | "),
		color,
	)

	// Link the node to its children.
	for _, child := range node.Children {
		childID := exportNodeDOT(writer, child, visited, counter)
		fmt.Fprintf(writer, "  %s -> %s;\n", nodeID, childID)
	}

	return nodeID
}

//nolint:cyclop
func escapeDOT(str string) string {
	escaped := ""

	for _, elt := range str {
		switch elt {
		case '"':
			escaped += `\"`

		case '\n':
			escaped += `\n`

		case '<':
			escaped += `\<`

		case '>':
			escaped += `\>`

		case '{':
			escaped += `\{`

		case '}':
			escaped += `\}`

		case '[':
			escaped += `\[`

		case ']':
			escaped += `\]`

		default:
			escaped += string(elt)
		}
	}

	return escaped
}

// * graphviz.go ends here.
