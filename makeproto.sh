#!/bin/env bash
# -*- Mode: Shell-script -*-
#
# makeproto.sh --- Make protobuf.
#
# Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
#
# Author:     Paul Ward <paul@lisphacker.uk>
# Maintainer: Paul Ward <paul@lisphacker.uk>
# Created:    11 Jul 2025 17:31:24
#
# {{{ License:
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
