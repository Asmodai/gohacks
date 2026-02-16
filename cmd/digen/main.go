// -*- Mode: Go; auto-fill: t; fill-column: 78; -*-
//
// SPDX-License-Identifier: MIT
//
// main.go --- DI boilerplate generator.
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

package main

// * Imports:

import (
	"flag"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gitlab.com/tozd/go/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

// * Constants:

const (
	defaultFileMode = 0o644
)

// * Code:

type spec struct {
	Basename string
	TypeExpr string
	Key      string
	Fallback string
}

func (s spec) valid() bool {
	return !(len(s.Basename) == 0 || len(s.TypeExpr) == 0 || len(s.Key) == 0 || len(s.Fallback) == 0)
}

type fileJob struct {
	PkgName string
	Out     string
	TestOut string
	Specs   []spec
	Imports map[string]string
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "digen: "+format+"\n", args...)

	os.Exit(1)
}

func makePkgConfig() *packages.Config {
	return &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedModule |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo}
}

func parseCommentGroup(group *ast.CommentGroup) []spec {
	var specs = make([]spec, 0, 1)

	for _, comment := range group.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))

		if !strings.HasPrefix(text, "di:gen") {
			continue
		}

		args := strings.Fields(strings.TrimSpace(strings.TrimPrefix(
			text,
			"di:gen")))

		tmp := map[string]string{}

		for _, kv := range args {
			parts := strings.SplitN(kv, "=", 2) //nolint:mnd

			if len(parts) != 2 { //nolint:mnd
				continue
			}

			tmp[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}

		res := spec{
			Basename: tmp["basename"],
			TypeExpr: tmp["type"],
			Key:      tmp["key"],
			Fallback: tmp["fallback"],
		}

		if !res.valid() {
			continue
		}

		specs = append(specs, res)
	}

	return specs
}

func collectSpecs(pkg *packages.Package) []spec {
	var out []spec

	for _, sfile := range pkg.Syntax {
		if sfile == nil || sfile.Comments == nil {
			continue
		}

		for _, commentGroup := range sfile.Comments {
			out = append(out, parseCommentGroup(commentGroup)...)
		}
	}

	return out
}

//nolint:dupl
func writeTests(pkg *packages.Package, job *fileJob) error {
	sort.Slice(job.Specs, func(i, j int) bool {
		return job.Specs[i].Basename < job.Specs[j].Basename
	})

	var sbld strings.Builder

	if err := testTempl.Execute(&sbld, job); err != nil {
		return errors.WithStack(err)
	}

	formatted, err := imports.Process(
		job.Out,
		[]byte(sbld.String()),
		&imports.Options{Comments: true})
	if err != nil {
		_ = os.WriteFile(job.Out, []byte(sbld.String()), defaultFileMode)

		return errors.WithStack(err)
	}

	dir := "."
	if len(pkg.GoFiles) > 0 {
		dir = filepath.Dir(pkg.GoFiles[0])
	}

	path := filepath.Join(dir, job.TestOut)

	return errors.WithStack(os.WriteFile(path, formatted, defaultFileMode))
}

//nolint:dupl
func writeJob(pkg *packages.Package, job *fileJob) error {
	sort.Slice(job.Specs, func(i, j int) bool {
		return job.Specs[i].Basename < job.Specs[j].Basename
	})

	var sbld strings.Builder

	if err := codeTempl.Execute(&sbld, job); err != nil {
		return errors.WithStack(err)
	}

	formatted, err := imports.Process(
		job.Out,
		[]byte(sbld.String()),
		&imports.Options{Comments: true})
	if err != nil {
		_ = os.WriteFile(job.Out, []byte(sbld.String()), defaultFileMode)

		return errors.WithStack(err)
	}

	dir := "."
	if len(pkg.GoFiles) > 0 {
		dir = filepath.Dir(pkg.GoFiles[0])
	}

	path := filepath.Join(dir, job.Out)

	return errors.WithStack(os.WriteFile(path, formatted, defaultFileMode))
}

//nolint:funlen
func main() {
	var (
		pattern         string
		out             string
		test            string
		contextdiImport string
		errorsImport    string
		failOnEmpty     bool
	)

	flag.StringVar(&pattern, "pattern", ".", "package pattern to load (e.g. ., ./..., ./pkg)")
	flag.StringVar(&out, "out", "di_gen.go", "output file name (per package)")
	flag.StringVar(&test, "test", "di_gen_test.go", "output file name for unit tests")
	flag.StringVar(&contextdiImport, "contextdi", "github.com/Asmodai/gohacks/contextdi", "import path for contextdi")
	flag.StringVar(&errorsImport, "errors", "gitlab.com/tozd/go/errors", "import path for errors package")
	flag.BoolVar(&failOnEmpty, "fail-empty", false, "error if no //di:gen annotations found")
	flag.Parse()

	cfg := makePkgConfig()

	pkgs, err := packages.Load(cfg, pattern)
	if err != nil {
		die("load: %v", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		die("packages contain errors")
	}

	var anyFound bool

	for _, pkg := range pkgs {
		specs := collectSpecs(pkg)

		if len(specs) == 0 {
			if failOnEmpty {
				die("no //di:gen specs found in %s",
					pkg.PkgPath)
			}

			continue
		}

		anyFound = true

		job := fileJob{
			PkgName: pkg.Name,
			Out:     out,
			TestOut: test,
			Specs:   specs,
			Imports: map[string]string{
				"":          "context",
				"contextdi": contextdiImport,
				"errors":    errorsImport,
			},
		}

		if err := writeJob(pkg, &job); err != nil {
			die("write %s: %v", pkg.PkgPath, err)
		}

		if err := writeTests(pkg, &job); err != nil {
			die("write %s: %v", pkg.PkgPath, err)
		}
	}

	if !anyFound && failOnEmpty {
		die("no //di:gen specs found")
	}
}

// * main.go ends here.
