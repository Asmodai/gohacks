#!/bin/sh

# usage: ./update_gohacks.sh v1
# Result: "github.com/Asmodai/gohacks/" → "github.com/Asmodai/gohacks/v1/"
#         "github.com/Asmodai/gohacks/v0/" → "github.com/Asmodai/gohacks/v1/"

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <version> (e.g. v1 or v2)"
  exit 1
fi

VERSION="$1"
MODULE_PREFIX="github.com/Asmodai/gohacks"
NEW_IMPORT="${MODULE_PREFIX}/${VERSION}/"

# grep -rl = recursive list of matching files
# sed -i '' for BSD/macOS compatibility; GNU sed users can just use -i

find . -type f -name "*.go" | while IFS= read -r file; do
  if grep -q "${MODULE_PREFIX}/" "$file"; then
    echo "Updating: $file"
    # Replace with or without trailing version
    sed -i.bak -e "s|${MODULE_PREFIX}\(/\(v[0-9][^/]*\)\)\?/|${NEW_IMPORT}|g" "$file" && rm -f "${file}.bak"
  fi
done

echo "All done. You're now riding version ${VERSION} like a majestic code cowboy."

