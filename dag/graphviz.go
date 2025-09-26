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
	//
	// XXX This is ugly and should be re-written to not be ugly.
	//
	fmt.Fprintln(writer, "digraph DAG {")
	fmt.Fprintln(writer, "  graph [bgcolor=transparent];")
	fmt.Fprintln(writer, "  rankdir=TB;")
	fmt.Fprintln(writer, `  node [shape=record, style="rounded,filled", fontname="monospace", fontsize=10];`)
	fmt.Fprintln(writer, `  edge [color="#666666"];`)

	visited := make(map[*node]string)
	counter := 0

	exportNodeDOT(writer, root, visited, &counter)

	fmt.Fprintln(writer, "}")
}

//nolint:funlen
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

	color := "#6699cc"
	fillColor := "#d8eaf8"
	fontColor := "#003366"
	arrowColor := "#4CAF50"

	if kind == noopIsn {
		color = "#999999"
		fillColor = "#f0f0f0"
		fontColor = "#000000"
	}

	if len(node.ActionName) > 0 {
		parts = append(
			parts,
			"<f2> Action: "+node.ActionName,
		)
		color = "#66cc66"
		fillColor = "#e8f6e0"
		fontColor = "#003300"
	}

	if len(node.FailureName) > 0 {
		parts = append(
			parts,
			"<f3> Failure: "+node.FailureName,
		)
		color = "#66cc66"
		fillColor = "#e8f6e0"
		fontColor = "#003300"
	}

	// Generate the node.
	fmt.Fprintf(
		writer,
		`  %s [label="{%s}" fillcolor="%s", fontcolor="%s", color="%s"];`+"\n",
		nodeID,
		strings.Join(parts, " | "),
		fillColor,
		fontColor,
		color,
	)

	// Link the node to its children.
	for _, child := range node.Children {
		childID := exportNodeDOT(writer, child, visited, counter)
		fmt.Fprintf(
			writer,
			"  %s -> %s [color=%q, style=bold];\n",
			nodeID,
			childID,
			arrowColor,
		)
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
