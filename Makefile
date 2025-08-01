# -*- Mode: Makefile -*-
#
# Makefile --- gohacks makefile.
#
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2021-2025 Paul Ward <paul@lisphacker.uk>
#
# Author:     Paul Ward <paul@lisphacker.uk>
# Maintainer: Paul Ward <paul@lisphacker.uk>
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

# Go package.
PACKAGE = gohack

# Directories.
DIR = $(PWD)

# Source modules.
MODULES = amqp            \
	  apiclient       \
	  apiserver       \
	  app             \
	  config          \
	  contextdi       \
	  contextext      \
	  crypto          \
	  dag             \
	  database        \
	  debug           \
	  dynworker       \
	  events          \
	  fileloader      \
	  generics        \
	  logger          \
	  math            \
	  math/conversion \
	  memoise         \
	  process         \
	  responder       \
	  rfc3339         \
	  rlhttp          \
	  secrets         \
	  semver          \
	  service         \
	  sysinfo         \
	  timedcache      \
	  types           \
	  utils           \
	  validator

# Binaries.
PROTOC          ?= protoc

# Settings
LINT_REPORT ?= "golint.xml"
HTML_REPORT ?= "golint.html"

all: deps

.PHONY: configs doc protobuf mocks

deps:
	@echo "Getting dependencies"
	@go work vendor

tidy:
	@echo "Tidying mod dependencies"
	@go mod tidy

tooling:
	@go install github.com/google/go-licenses@latest
	@go install go.uber.org/mock/mockgen@latest
	@go install github.com/go-critic/go-critic/cmd/gocritic@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/google/go-licenses@latest
	@pip install --break-system-packages junit2html

listdeps:
	@echo "Listing dependencies:"
	@go list -m all

prunedeps:
	@echo "Pruning dependencies"
	@go mod tidy

# Please note that the below is absolutely terrible.  golangci really
# wants structured output, so that I can say "Please to be generating junit
# XML to this file".  Or, some other format and a HTML converter.  Either way,
# the sed invocation is terrible.
lint:
	@echo "Running linter."
	@if [ -f "$$HOME/.local/bin/junit2html" ]; then                        \
		echo "Generating $(LINT_REPORT) and $(HTML_REPORT)";           \
		golangci-lint run                                              \
			--out-format junit-xml                                 \
			| sed -n '1h;1!H;$$ {g;s|\(</testsuites>\).*|\1|; p;}' \
			> $(LINT_REPORT);                                      \
		junit2html $(LINT_REPORT) $(HTML_REPORT);                      \
	else                                                                   \
		golangci-lint run;                                             \
	fi
	@echo "Done"

critic:
	@echo "Everyone is a critic..."
	@gocritic check

build:
	@echo "THIS IS A LIBRARY"

test: deps
	@echo "Running tests"
	@go test $$(go list ./... | grep -v "mocks/") \
		-coverprofile=tests.out               \
		--tags testing
	@go tool cover -html=tests.out -o coverage.html

run:
	@echo "THIS IS A LIBRARY"

rundebug:
	@echo "THIS IS A LIBRARY"

clean:
	@echo "Cleaning"
	@rm -f *.out
	@rm -f golint.*
	@rm -f report.xml
	@rm -f coverage.html
	@rm -f doc/go/*.md

protobuf:
	@PROTOC="$(PROTOC)" ./makeproto.sh
	@echo "Done."

mocks:
	@./makemocks.sh
	@echo "Done."

doc:
	@MODULES="$(MODULES)" ./makedoc.sh
	@echo "Done."

# Makefile ends here.
