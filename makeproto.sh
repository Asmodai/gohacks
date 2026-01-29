#!/bin/env bash
# -*- Mode: Shell-script -*-
#
# SPDX-License-Identifier: MIT
#
# makeproto.sh --- Make protobuf.
#
# Copyright (c) 2025-2026 Paul Ward <paul@lisphacker.uk>
#
# Author:     Paul Ward <paul@lisphacker.uk>
# Maintainer: Paul Ward <paul@lisphacker.uk>
# Created:    11 Jul 2025 17:31:24
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

set -euo pipefail

BASEDIR=$(pwd)

PROTOC=${PROTOC:-protoc}
GO_OUT_OPTS="paths=source_relative"

find . -type f -path "*/protobuf/*.proto" | while read -r protofile
do
    protofile=${protofile#./}
    genfile=${protofile/protobuf/gen}
    genfile=${genfile%.proto}.pb.go
    outdir=$(dirname "${genfile}")
    protopath=$(dirname "${protofile}")

    echo "Generating '$genfile' from '$protofile' with '${PROTOC}'"

    mkdir -p "$outdir"

    ${PROTOC}                          \
        --proto_path="${protopath}"    \
        --go_out="${outdir}"           \
        --go-grpc_out="${outdir}"      \
        --go_opt="${GO_OUT_OPTS}"      \
        --go-grpc_opt="${GO_OUT_OPTS}" \
        "${protofile}"
done

# makeproto.sh ends here.
