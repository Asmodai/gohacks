#!/bin/bash
# -*- Mode: Shell-script -*-
#
# makemocks.sh --- Make mocks from defined interfaces.
#
# Copyright (c) 2021-2024 Paul Ward <asmodai@gmail.com>
#
# Author:     Paul Ward <asmodai@gmail.com>
# Maintainer: Paul Ward <asmodai@gmail.com>
# Created:    01 Sep 2021 13:13:28
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

test -d "mocks" && rm -rf "mocks"

ROOT=$(pwd)
FILES=$(find . -iname "*.go" | grep -v "/vendor/")
MOCK_PATH="mocks"

test -d "${MOCK_PATH}" || mkdir "${MOCK_PATH}"

for file in ${FILES}
do
    # Build mocks from the file?
    grep "mock:yes" $file >/dev/null 2>&1
    if [ $? -eq 1 ]
    then
        # Nope.
        continue
    fi

    fname=$(basename ${file})
    pname=$(basename $(dirname ${file}))
    output="${MOCK_PATH}/${pname}/$(echo ${fname} | cut -d. -f1)_mock.go"

    echo "Processing ${pname}/${fname} => ${output}"

    mockgen                      \
        -package="${pname}"      \
        -source="${file}"        \
        -destination="${output}"
    case $? in
        0)
            echo "Mock generation successful."
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
