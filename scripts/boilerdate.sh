#!/bin/sh
# -*- Mode: Shell-script -*-
#
# SPDX-License-Identifier: MIT
#
# boilerdate.sh --- Update copyright and email.
#
# Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
#
# Author:     Paul Ward <paul@lisphacker.uk>
# Maintainer: Paul Ward <paul@lisphacker.uk>
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

CURRENT_YEAR=$(date +%Y)
EMAIL="${EMAIL:-paul@lisphacker.uk}"  # default email, override by env

usage() {
  echo "Usage: $0 file1 [file2 ...]" >&2
  echo "Set EMAIL environment variable to override default email address." >&2
  exit 1
}

if [ $# -eq 0 ]; then
  usage
fi

for file in "$@"; do
  if [ ! -f "$file" ]; then
    echo "Warning: File not found: $file" >&2
    continue
  fi

  cp "$file" "${file}.bak" || {
    echo "Error: Could not backup $file" >&2
    continue
  }

  TMPFILE=$(mktemp) || {
    echo "Error: Could not create temp file" >&2
    exit 1
  }

  while IFS= read -r line; do
    case "$line" in
      *Copyright\ \(c\)*)
        YEAR_PART=$(echo "$line" | sed -n 's/.*Copyright (c) \([0-9]\{4\}\)\(-[0-9]\{4\}\)\?.*/\1\2/p')
        START_YEAR=$(echo "$YEAR_PART" | cut -d'-' -f1)
        END_YEAR=$(echo "$YEAR_PART" | cut -s -d'-' -f2)

        if [ -z "$END_YEAR" ]; then
          if [ "$START_YEAR" -lt "$CURRENT_YEAR" ]; then
            NEW_YEAR="$START_YEAR-$CURRENT_YEAR"
          else
            NEW_YEAR="$START_YEAR"
          fi
        else
          if [ "$END_YEAR" -lt "$CURRENT_YEAR" ]; then
            NEW_YEAR="$START_YEAR-$CURRENT_YEAR"
          else
            NEW_YEAR="$START_YEAR-$END_YEAR"
          fi
        fi

        UPDATED_LINE=$(echo "$line" | sed "s/Copyright (c) [0-9]\{4\}\(-[0-9]\{4\}\)\?/Copyright (c) $NEW_YEAR/")
        UPDATED_LINE=$(echo "$UPDATED_LINE" | sed "s/<[^>]*>$/<$EMAIL>/")
        ;;

      *\<*\@*\>*)
        # Replace email inside angle brackets in the line
        UPDATED_LINE=$(echo "$line" | sed "s/<[^>]*>$/<$EMAIL>/")
        ;;

      *)
        UPDATED_LINE="$line"
        ;;
    esac

    printf '%s\n' "$UPDATED_LINE"
  done < "$file" > "$TMPFILE" && mv "$TMPFILE" "$file"
done

