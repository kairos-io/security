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
CUTOFF_VERSION="v2.10.0"

# Ensure jq_script is properly defined as a single-line string
jq_script='[.[] | select(.tag_name) | {tag_name: .tag_name, version: .tag_name[1:], image: ("'"$IMAGE_BASE"':" + .tag_name + "'"$SUFFIX"'")}] | group_by(.version | split(".") | .[:2] | join(".")) | map(sort_by(.version | split(".") | map(tonumber)) | reverse | .[0] | {image: .image, tag_name: .tag_name})'

# Initialize or clear the JSON array in the file
echo "[]" > "$OUTPUT_FILE"

# Function to compare versions
version_gt() {
    [ "$(printf '%s\n' "$@" | sort -V | head -n 1)" != "$1" ]
}

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
    local new_data=$(echo "$releases" | jq "$jq_script")
    local filtered_data="[]"
    for release in $(echo "$new_data" | jq -r '.[].tag_name'); do
        if version_gt "$release" "$CUTOFF_VERSION"; then
            filtered_data=$(echo "$filtered_data" | jq --argjson release "$(echo "$new_data" | jq --arg tag_name "$release" '.[] | select(.tag_name == $tag_name)')" '. + [$release]')
        fi
    done
    local existing_data=$(cat "$OUTPUT_FILE")
    echo "$existing_data" "$filtered_data" | jq -s '.[0] + .[1] | group_by(.tag_name | split(".")[:2] | join(".")) | map(max_by(.tag_name))' > "$OUTPUT_FILE.tmp" && mv "$OUTPUT_FILE.tmp" "$OUTPUT_FILE"

    return 0
}

# Loop through all pages of releases and process them
while fetch_and_append_releases $PAGE; do
    ((PAGE++))
done

echo "Image versions have been written to $OUTPUT_FILE"