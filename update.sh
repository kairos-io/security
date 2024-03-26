#!/bin/bash
set -e

IMAGE_BASE="quay.io/kairos/framework"
REPO="kairos-io/kairos-framework"
SUFFIX="_generic"
PREFIX="v"
OUTPUT_FILE=images.json
RELEASES=$(curl -s "https://api.github.com/repos/$REPO/releases?per_page=100")

# Parse the releases to create objects with the "image" key, combining the base image path with the version number, and write to the file
if [ "$ALL_RELEASES" = "true" ]; then
    echo "$RELEASES" | jq "[.[] | {image: (\"$IMAGE_BASE:$PREFIX\" + .tag_name + \"$SUFFIX\")}]" > "$OUTPUT_FILE"
fi

# Parse the releases, group by major.minor version, sort by version, and keep only the latest release for each branch
jq_script='
    [ 
        .[] 
        | select(.tag_name ) 
        | { tag_name: .tag_name, version: .tag_name[1:], image: ("'"$IMAGE_BASE"':" + .tag_name + "_generic") }
    ] 
    | group_by(.version | split(".") | .[:2] | join(".")) 
    | map(
        sort_by(.version | split(".") | map(tonumber)) | reverse | .[0] 
        | {image: .image, tag_name: .tag_name}
      )
'

echo "$RELEASES" | jq "$jq_script" > "$OUTPUT_FILE"

echo "Image versions have been written to $OUTPUT_FILE"