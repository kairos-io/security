#!/bin/bash
set -e

# Configuration
IMAGE_BASE="quay.io/kairos/framework"
REPO="kairos-io/kairos-framework"
SUFFIX=""
PREFIX="v"
OUTPUT_FILE=images.json
PAGE=1
PER_PAGE=100

# Ensure jq_script is properly defined as a single-line string
jq_script='[.[] | select(.tag_name) | {tag_name: .tag_name, version: .tag_name[1:], image: ("'"$IMAGE_BASE"':" + .tag_name + "'"$SUFFIX"'")}] | group_by(.version | split(".") | .[:2] | join(".")) | map(sort_by(.version | split(".") | map(tonumber)) | reverse | .[0] | {image: .image, tag_name: .tag_name})'

# Initialize or clear the JSON array in the file
echo "[]" > "$OUTPUT_FILE"

# Function to fetch releases, process, and append to the file
fetch_and_append_releases() {
    local page=$1
    local releases=$(curl -s "https://api.github.com/repos/$REPO/releases?per_page=$PER_PAGE&page=$page")

    # Check if we got any releases back
    if [ "$(echo "$releases" | jq length)" -eq 0 ]; then
        echo "No more releases to process."
        return 1
    fi

    # Process releases and append to the output file
    if [ "$ALL_RELEASES" = "true" ]; then
        local new_data=$(echo "$releases" | jq "[.[] | {image: (\"$IMAGE_BASE:$PREFIX\" + .tag_name + \"$SUFFIX\")}]")
    else
        local new_data=$(echo "$releases" | jq "$jq_script")
    fi
    local existing_data=$(cat "$OUTPUT_FILE")
    echo "$existing_data" "$new_data" | jq -s '.[0] + .[1]' > "$OUTPUT_FILE.tmp" && mv "$OUTPUT_FILE.tmp" "$OUTPUT_FILE"

    return 0
}

# Loop through all pages of releases and process them
while fetch_and_append_releases $PAGE; do
    ((PAGE++))
done

echo "Image versions have been written to $OUTPUT_FILE"
