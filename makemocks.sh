#!/bin/bash
# -*- Mode: Shell-script -*-
#
# makemocks.sh --- Make mocks from defined interfaces.
#
# Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
#
# Author:     Paul Ward <asmodai@gmail.com>
# Maintainer: Paul Ward <asmodai@gmail.com>
# Created:    01 Sep 2021 13:13:28
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

test -d "mocks" && rm -rf "mocks"

ROOT=$(pwd)
FILES=$(find . -iname "i*.go" | grep -v "/vendor/")

MOCK_PATH="mocks"

for file in ${FILES}
do
    # Do we define any interfaces?
    egrep -e "type\s\w.+\sinterface\s{" $file >/dev/null 2>&1
    if [ $? -eq 1 ]
    then
        # No, move on to next file.
        continue
    fi

    fname=$(basename ${file})
    pname=$(basename $(dirname ${file}))
    output="${ROOT}/${MOCK_PATH}/${pname}/${fname}"

    # Interfaces defined, let's run mockgen
    echo "Processing ${fname} (${pname}) => ${output}"

    mockgen                             \
        -package="mock_${pname}" \
        -source=${file}                 \
        -destination="${output}"
    case $? in
        0)
            sed                                           \
                -i ''                                     \
                -e "2s/^//p; 2s/^.*/\/\/ +build testing/" \
                "${output}"
            ;;
        127)
            echo "mockgen not installed."
            exit 1
            ;;
        *)
            echo "mockgen failed, exit code $?"
            exit $?
            ;;
    esac
done

# makemocks.sh ends here.
