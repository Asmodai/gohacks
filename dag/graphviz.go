// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// graphviz.go --- Graphviz exporter.
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

package dag

// * Imports:

import (
	"fmt"
	"html"
	"io"
	"strings"
)

// * Constants:

const (
	RankDirectionTB RankDirection = 0
	RankDirectionLR RankDirection = 1

	DefaultMaxInlineItems int = 6

	rankDirTB string = "TB"
	rankDirLR string = "LR"
)

// * Code:

// ** Rank direction:

type RankDirection int

func (rd RankDirection) Direction() string {
	if rd == RankDirectionLR {
		return rankDirLR
	}

	return rankDirTB
}

// ** Configuration:

type DOTOptions struct {
	// "TB" or "LR"
	RankDir RankDirection

	// Emit `Action` nodes?
	EmitActionNodes bool

	// Emit `Failure` nodes?
	EmitFailureNodes bool

	// Maximum number of items to show inline before truncating.
	MaxInlineItems int
}

// ** DOT Emitter:

type dotEmitter struct {
	writer  io.Writer
	options DOTOptions
	predIDs map[*node]string
	predN   int
	actN    int
	failN   int
}

func (d *dotEmitter) summariseNames(list []*nodeAction, limit int) string {
	if len(list) == 0 {
		return ""
	}

	names := make([]string, 0, len(list))
	for _, a := range list {
		names = append(names, actionName(a))
	}

	if len(names) > limit {
		names = append(names[:limit], "...")
	}

	return strings.Join(names, ", ")
}

func (d *dotEmitter) emitFailures(current *node, label string) {
	if d.options.EmitFailureNodes {
		for _, elt := range current.Failures {
			fid := d.emitFailureNode(elt)
			fmt.Fprintf(
				d.writer,
				"  %s -> %s [color=%q, style=dashed, label=%q];\n",
				label,
				fid,
				"#c62828",
				"",
			)
		}
	} else {
		tid := d.emitTerminalNode("FAIL", "#C62828", "#FFEBEE")
		fmt.Fprintf(
			d.writer,
			"  %s -> %s [color=%q, style=dashed, label=%q];\n",
			label,
			tid,
			"#c62828",
			"",
		)
	}
}

func (d *dotEmitter) emitFailureNode(node *nodeAction) string {
	fid := fmt.Sprintf("f%d", d.failN)
	d.failN++

	name := html.EscapeString(actionName(node))

	fmt.Fprintf(
		d.writer,
		`  %s [shape=box, style="rounded,filled", fillcolor="%s", color="%s", fontcolor="%s", label=%q];`+"\n",
		fid,
		"#ffebee",
		"#c62828",
		"#7f0000",
		"FAILURE: "+name,
	)

	return fid
}

func (d *dotEmitter) emitTerminalNode(label, border, fill string) string {
	tid := fmt.Sprintf("t%d", d.failN)
	d.failN++

	fmt.Fprintf(d.writer,
		`  %s [shape=oval, style="filled", fillcolor="%s", color="%s", fontcolor="%s", label=%q];`+"\n",
		tid,
		fill,
		border,
		"#000000",
		label,
	)

	return tid
}

func (d *dotEmitter) emitActions(current *node, label string) {
	for _, elt := range current.Actions {
		aid := d.emitActionNode(elt)
		fmt.Fprintf(
			d.writer,
			"  %s -> %s [color=%q, style=dashed, label=%q];\n",
			label,
			aid,
			"#2e7d32",
			"",
		)
	}
}

func (d *dotEmitter) emitActionNode(node *nodeAction) string {
	aid := fmt.Sprintf("a%d", d.actN)
	d.actN++

	name := html.EscapeString(actionName(node))

	fmt.Fprintf(
		d.writer,
		`  %s [shape=box, style="rounded,filled", fillcolor="%s", color="%s", fontcolor="%s", label=%q];`+"\n",
		aid,
		"#e8f5e9",
		"#2e7d32",
		"#1b5e20",
		"ACTION: "+name,
	)

	return aid
}

func (d *dotEmitter) emitEdges(current *node, label string) {
	for _, elt := range current.Children {
		cid := d.emitPredicate(elt)
		fmt.Fprintf(
			d.writer,
			"  %s -> %s [color=%q, penwidth=2, label=%q];\n",
			label,
			cid,
			"#2e7d32",
			"",
		)
	}
}

func (d *dotEmitter) predicateHTMLLabel(current *node) string {
	kind := html.EscapeString(current.Predicate.Instruction())
	pred := html.EscapeString(current.Predicate.String())

	// Build action/failure summaries (inline mode)
	act := d.summariseNames(current.Actions, d.options.MaxInlineItems)
	fail := d.summariseNames(current.Failures, d.options.MaxInlineItems)

	var sbld strings.Builder

	sbld.WriteString(`<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4" BGCOLOR="#d8eaf8">`)

	// Header row: instruction
	sbld.WriteString(`<TR><TD BGCOLOR="#6699cc"><FONT COLOR="#ffffff"><B>`)
	sbld.WriteString(kind)
	sbld.WriteString(`</B></FONT></TD></TR>`)

	// Predicate row
	sbld.WriteString(`<TR><TD ALIGN="LEFT"><FONT FACE="monospace">`)
	sbld.WriteString(pred)
	sbld.WriteString(`</FONT></TD></TR>`)

	// Actions row (green)
	if !d.options.EmitActionNodes && act != "" {
		sbld.WriteString(`<TR><TD ALIGN="LEFT" BGCOLOR="#e8f5e9"><FONT COLOR="#1b5e20">`)
		sbld.WriteString(`Actions: `)
		sbld.WriteString(html.EscapeString(act))
		sbld.WriteString(`</FONT></TD></TR>`)
	}

	// Failures row (red)
	if !d.options.EmitFailureNodes && fail != "" {
		sbld.WriteString(`<TR><TD ALIGN="LEFT" BGCOLOR="#ffebee"><FONT COLOR="#7f0000">`)
		sbld.WriteString(`Failures: `)
		sbld.WriteString(html.EscapeString(fail))
		sbld.WriteString(`</FONT></TD></TR>`)
	}

	sbld.WriteString(`</TABLE>`)

	return sbld.String()
}

func (d *dotEmitter) emitPredicate(root *node) string {
	if id, ok := d.predIDs[root]; ok {
		return id
	}

	pid := fmt.Sprintf("p%d", d.predN)
	d.predN++
	d.predIDs[root] = pid

	label := d.predicateHTMLLabel(root)
	fmt.Fprintf(d.writer, "  %s [shape=plain, label=<%s>];\n", pid, label)

	if len(root.Failures) > 0 {
		d.emitFailures(root, pid)
	}

	d.emitEdges(root, pid)

	if len(root.Actions) > 0 && d.options.EmitActionNodes {
		d.emitActions(root, pid)
	}

	return pid
}

// ** Functions:

func ExportToDOT(writer io.Writer, root *node) {
	ExportToDOTWithOptions(
		writer,
		root,
		DOTOptions{
			RankDir:          RankDirectionTB,
			EmitActionNodes:  true,
			EmitFailureNodes: true,
			MaxInlineItems:   DefaultMaxInlineItems,
		},
	)
}

func ExportToDOTWithOptions(writer io.Writer, root *node, opt DOTOptions) {
	if opt.RankDir < 0 || opt.RankDir > RankDirectionLR {
		opt.RankDir = RankDirectionTB
	}

	if opt.MaxInlineItems <= 0 {
		opt.MaxInlineItems = DefaultMaxInlineItems
	}

	inst := &dotEmitter{
		writer:  writer,
		options: opt,
		predIDs: make(map[*node]string),
	}

	emitPreamble(writer, opt)
	inst.emitPredicate(root)
	emitPostamble(writer)
}

func emitPreamble(writer io.Writer, opt DOTOptions) {
	fmt.Fprintln(writer, "digraph DAG {")
	fmt.Fprintln(writer, "  graph [bgcolor=transparent];")
	fmt.Fprintf(writer, "  rankdir=%s;\n", opt.RankDir.Direction())
	fmt.Fprintln(writer, `  node [fontname="monospace", fontsize=10];`)
	fmt.Fprintln(writer, `  edge [fontname="monospace", fontsize=9, color="#666666"];`)
}

func emitPostamble(writer io.Writer) {
	fmt.Fprintln(writer, "}")
}

func actionName(node *nodeAction) string {
	if node == nil {
		return ""
	}

	if len(node.Name) > 0 {
		return node.Name
	}

	return "<unnamed>"
}

// * graphviz.go ends here.
