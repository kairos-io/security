#!/bin/bash
set -ex

DIR="$1"

# recursively run govulncheck for all binaries in the directory
for file in $(find $DIR -type f -executable); do
    # check if the file is an ELF binary
    if [ "$(file $file | grep -o 'ELF')" != "ELF" ]; then
        continue
    fi
    govulncheck -json -mode=binary $file > /reports/$(basename $file).json
done
