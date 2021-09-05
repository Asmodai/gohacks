# -*- Mode: Makefile -*-
#
# Makefile --- gohacks makefile.
#
# Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
#
# Author:     Paul Ward <asmodai@gmail.com>
# Maintainer: Paul Ward <asmodai@gmail.com>
# Created:    11 Aug 2021 04:28:08
#
#{{{ License:
#
# This program is free software: you can redistribute it
# and/or modify it under the terms of the GNU General Public
# License as published by the Free Software Foundation,
# either version 3 of the License, or (at your option) any
# later version.
#
# This program is distributed in the hope that it will be
# useful, but WITHOUT ANY  WARRANTY; without even the implied
# warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR
# PURPOSE.  See the GNU General Public License for more
# details.
#
# You should have received a copy of the GNU General Public
# License along with this program.  If not, see
# <http://www.gnu.org/licenses/>.
#
#}}}
#{{{ Commentary:
#
#}}}

PACKAGE = gohack

DIR = $(PWD)

MODULES = app             \
	  config          \
	  database        \
	  di              \
	  process         \
	  math/conversion \
	  sysinfo         \
	  rfc3339         \
	  types

all: deps

.PHONY: configs doc

deps:
	@echo Getting dependencies
	@go mod vendor

tidy:
	@echo Tidying mod dependencies
	@go mod tidy

listdeps:
	@echo Listing dependencies:
	@go list -m all

prunedeps:
	@echo Pruning dependencies
	@go mod tidy

build:
	@echo THIS IS A LIBRARY

test: deps
	@echo Running tests
	@go test $$(go list ./...) -coverprofile=tests.out
	@go tool cover -html=tests.out -o coverage.html

run:
	@echo THIS IS A LIBRARY

rundebug:
	@echo THIS IS A LIBRARY

clean:
	@echo Cleaning
	@rm *.out
	@rm doc/*.md

doc:
	@echo Generating documentation
	@test -d doc || mkdir doc
	@for dir in $(MODULES); do \
		echo "Generating $${dir}.md"                                        ;\
		godocdown -template doc/gohacks.template ./$${dir}/ >doc/$${dir}.md ;\
	done

# Makefile ends here.
