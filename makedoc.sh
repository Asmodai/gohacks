#!/usr/bin/env bash
# -*- Mode: Shell-script -*-
#
# makedoc.sh --- Generate documentation.
#
# SPDX-License-Identifier: MIT
#
# Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
#
# Author:     Paul Ward <paul@lisphacker.uk>
# Maintainer: Paul Ward <paul@lisphacker.uk>
# Created:    26 Jul 2025 07:26:09
#
# {{{ License:
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
# }}}
# {{{ Commentary:
#
# }}}

MODULES="${MODULES:-none}"
path="doc/go"

if [ "${MODULES}" == "none" ]
then
    echo "No modules specified."
    echo
    echo Usage: 'MODULES="..."' $0

    exit 1
fi

echo "Generating docs in '${path}'."

# Make doc directory if needed.
test -d "${path}"                             \
    || mkdir -p "${path}"                     \
    && $(rm -rf "${path}"; mkdir -p "${path}")

for module in ${MODULES}
do
    echo "  Generating ${module} documentation."

    subdir=$(dirname "${module}")
    outfile="${path}/${module}.md"
    outdir=$(dirname "${outfile}")

    mkdir -p "${outdir}"

    godocdown \
        -template doc/templates/gohacks.template \
        "./${module}/" \
        >"${outfile}"
done

# makedoc.sh ends here.
