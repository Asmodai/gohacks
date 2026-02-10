# -*- Mode: Makefile -*-
#
# Makefile --- gohacks makefile.
#
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2021-2026 Paul Ward <paul@lisphacker.uk>
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
PACKAGE = gohacks

# Directories.
DIR = $(CURDIR)

# Source modules.
MODULES = amqp            \
	  apiclient       \
	  apiserver       \
	  app             \
	  config          \
	  contextdi       \
	  contextext      \
	  conversion      \
	  crypto          \
	  dag             \
	  database        \
	  debug           \
	  dynworker       \
	  errx            \
	  events          \
	  fileio          \
	  generics        \
	  health          \
	  logger          \
	  lucette         \
	  math            \
	  memoise         \
	  metadata        \
	  process         \
	  protocols       \
	  responder       \
	  rfc3339         \
	  rlhttp          \
	  scheduler       \
	  secrets         \
	  selector        \
	  semver          \
	  service         \
	  stringy         \
	  sysinfo         \
	  timedcache      \
	  types           \
	  utils           \
	  validator       \
	  wal

# Options
.SHELLFLAGS   := -eu -o pipefail -c
.DEFAULT_GOAL := help

# Binaries.
SHELL      := /bin/bash
PROTOC     ?= protoc
JUNIT2HTML ?= $$(command -v junit2html)

# Settings
LINT_REPORT ?= "golint.xml"
HTML_REPORT ?= "golint.html"

all: deps

.PHONY: configs doc protobuf mocks tooling listdeps lint critic \
	betteralign staticcheck vet test bench bench-save benchstat prof \
	race build run rundebug clean protobuf mocks doc fmt imports vuln ci

deps:
	@echo "Getting dependencies"
	@go work vendor

tidy:
	@echo "Tidying mod dependencies"
	@go mod tidy

tooling:
	@echo "Installing Go tooling..."
	@go install github.com/google/go-licenses@latest
	@go install go.uber.org/mock/mockgen@latest
	@go install github.com/go-critic/go-critic/cmd/gocritic@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/dkorunic/betteralign/cmd/betteralign@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install golang.org/x/perf/cmd/benchstat@latest
	@pip install --break-system-packages junit2html

listdeps:
	@echo "Listing dependencies:"
	@go list -m all

# Please note that the below is absolutely terrible.  golangci really
# wants structured output, so that I can say "Please to be generating junit
# XML to this file".  Or, some other format and a HTML converter.  Either way,
# the sed invocation is terrible.
lint:
	@echo "Running linter."
	@-golangci-lint run                                              \
		--out-format junit-xml                                 \
		| sed -n '1h;1!H;$$ {g;s|\(</testsuites>\).*|\1|; p;}' \
		> $(LINT_REPORT)
	@-if [ -f "$(JUNIT2HTML)" ]; then                     \
		echo "Generating $(HTML_REPORT)";            \
		$(JUNIT2HTML) $(LINT_REPORT) $(HTML_REPORT); \
	fi
	@echo "Done"

critic:
	@echo "Everyone is a critic..."
	@gocritic check

vet:
	@echo "Running Go vet..."
	@go vet ./...

betteralign:
	@echo "Running betteralign..."
	@betteralign ./...

staticcheck:
	@echo "Running staticcheck..."
	@staticcheck ./...

test: deps
	@echo "Running tests"
	@go test $$(go list ./... | grep -v "mocks/") \
		-coverprofile=tests.out               \
		--tags testing
	@go tool cover -html=tests.out -o coverage.html

cover:
	@echo "Checking code coverage..."
	@go test  $$(go list ./... | grep -v "mocks/") \
		-coverprofile=tests.out                \
		--tags testing
	@go tool cover -func=tests.out | tee coverage.txt
	@awk  '/^total:/ { if ($$3+0 < 80) { print "Coverage < 80%"; exit 1 } }' coverage.txt

bench:
	@echo "Benchmarking..."
	@-cp -f bench.txt bench_before.txt
	@-go test -bench=. -benchmem -count=6 ./... | tee bench.txt

benchcmp:
	@benchstat bench_before.txt bench.txt

prof:
	@echo "Profiling..."
	@go test -run=^$ -bench=. -benchmem \
		-cpuprofile=cpu.out         \
		-memprofile=mem.out         \
		./...
	@echo "pprof:"
	@echo "    go tool pprof -http=:0 cpu.out"
	@echo "    go tool pprof -http=:0 mem.out"

race:
	@echo "Checking for race conditions..."
	@go test -race -shuffle=on ./...

fuzz:
	@echo "Fuzzing..."
	@go test ./... -run=^$ -fuzz=Fuzz -fuzztime=30s

fmt:
	@go fmt ./...

imports:
	@goimports -w .

vuln:
	@govulncheck ./...

vulnreport:
	@govulncheck -show=verbose ./...

ci: tidy vet staticcheck betteralign lint test

build:
	@echo "THIS IS A LIBRARY"

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
	@rm -f coverage.txt
	@rm -f bench.txt
	@rm -f bench_before.txt

clean-lint-cache:
	@golangci-lint cache clean

protobuf:
	@PROTOC="$(PROTOC)" ./makeproto.sh
	@echo "Done."

generate:
	@go generate ./...

mocks:
	@./makemocks.sh
	@echo "Done."

doc:
	@MODULES="$(MODULES)" ./makedoc.sh
	@echo "Done."

war:
	@echo "Make l... wait."

love:
	@echo 'Make wa... HEY!'

fs-snapshot:
	@-scripts/fs-snapshot.sh

help:
	@echo "$(PACKAGE) Makefile Help"
	@echo
	@echo "Here are the interesting targets that you can invoke:"
	@echo
	@echo "bench            - Run benchmarks."
	@echo "benchcmp         - Compare benchmarks."
	@echo "betteralign      - Run betteralign and check structs."
	@echo "ci               - tidy, vet, staticcheck, betteralign, lint, test."
	@echo "cover            - Run coverage gate."
	@echo "clean            - Clean up temporary files."
	@echo "clean-lint-cache - Clear golangci cache."
	@echo "critic           - Run gocritic."
	@echo "deps             - Refresh dependencies."
	@echo "doc              - Generate gomarkdown documentation."
	@echo "fmt              - Run gofmt."
	@echo "fuzz             - Run fuzzers."
	@echo "generate         - Run Go generate."
	@echo "help             - you are reading it."
	@echo "imports          - Fix up imports."
	@echo "lint             - Run golangci-lint."
	@echo "listdeps         - List dependencies."
	@echo "mocks            - Generate mocks."
	@echo "prof             - Profile benchmarks."
	@echo "protobuf         - Run protobuf generation."
	@echo "race             - Check for race conditions."
	@echo "staticcheck      - Run staticcheck."
	@echo "test             - Run Go test."
	@echo "tidy             - Tidy up go.mod and go.sum."
	@echo "tooling          - Install Go tooling used by this Makefile."
	@echo "vet              - Run Go vet."
	@echo "vuln             - Perform vulnerability scan."
	@echo "vulnreport       - Perform verbose vulnerability scan."
	@echo
	@echo "Please be sure to 'make tooling'."

# Makefile ends here.
