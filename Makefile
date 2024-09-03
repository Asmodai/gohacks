# -*- Mode: Makefile -*-
#
# Makefile --- gohacks makefile.
#
# Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
#
# Author:     Paul Ward <asmodai@gmail.com>
# Maintainer: Paul Ward <asmodai@gmail.com>
# Created:    11 Aug 2021 04:28:08
#
#{{{ License:
#
# Permission is hereby granted, free of charge, to any person
# obtaining a copy of this software and associated documentation files
# (the "Software"), to deal in the Software without restriction,
# including without limitation the rights to use, copy, modify, merge,
# publish, distribute, sublicense, and/or sell copies of the Software,
# and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions: 
#
# The above copyright notice and this permission notice shall be
# included in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
# EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
# MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
# NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
# BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
# ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
#
#}}}
#{{{ Commentary:
#
#}}}

PACKAGE = gohack

DIR = $(PWD)

MODULES = apiclient       \
	  apiserver       \
	  app             \
	  config          \
	  context         \
	  crypto          \
	  database        \
	  events          \
	  generics        \
	  logger          \
	  math/conversion \
	  memoise         \
	  process         \
	  rfc3339         \
	  rlhttp          \
	  secrets         \
	  semver          \
	  service         \
	  sysinfo         \
	  types           \
	  utils

all: deps

.PHONY: configs doc

deps: tidy
	@echo "Getting dependencies"
	@go work vendor

tidy:
	@echo "Tidying mod dependencies"
	@go mod tidy

tooling:
	@go install github.com/google/go-licenses@latest
	@go install go.uber.org/mock/mockgen@latest
	@go install github.com/go-critic/go-critic/cmd/gocritic@latest
	@go install github.com/google/go-licenses@latest
	@pip install junit2html

listdeps:
	@echo "Listing dependencies:"
	@go list -m all

prunedeps:
	@echo "Pruning dependencies"
	@go mod tidy

lint:
	@echo "Running linter"
	@if [ -f "$$HOME/.local/bin/junit2html" ]; then               \
		echo "Generating golint.html";                        \
		golangci-lint run --out-format junit-xml >golint.xml; \
		junit2html golint.xml golint.html;                    \
	else                                                          \
		golangci-lint run;                                    \
	fi
	@echo "Done"

build:
	@echo "THIS IS A LIBRARY"

test: deps
	@echo "Running tests"
	@go test $$(go list ./...) -coverprofile=tests.out --tags testing
	@go tool cover -html=tests.out -o coverage.html

run:
	@echo "THIS IS A LIBRARY"

rundebug:
	@echo "THIS IS A LIBRARY"

clean:
	@echo "Cleaning"
	@rm *.out
	@rm golint.*
	@rm coverage.html
	@rm doc/*.md
	@rm -rf vendor

doc:
	@echo "Generating documentation"
	@test -d doc || mkdir doc
	@for dir in $(MODULES); do \
		echo "... Generating $${dir}.md"                                    ;\
		godocdown -template doc/gohacks.template ./$${dir}/ >doc/$${dir}.md ;\
	done
	@echo "Generating license information"
	@test -d vendor || go work vendor -o vendor
	@./checklic.sh >doc/dependencies.md
	@test -d vendor && rm -rf vendor
	@echo "Done."

# Makefile ends here.
