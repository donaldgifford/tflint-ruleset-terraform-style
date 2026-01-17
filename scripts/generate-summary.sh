#!/bin/bash
#
# Generates SUMMARY.md for stitchmd from markdown files in a source directory.
# Extracts the first H1 heading from each file as the link text.
#
# Usage: ./scripts/generate-summary.sh <source_directory> [title]
# Example: ./scripts/generate-summary.sh styleguides/terraform/src "Terraform Style Guide"

set -euo pipefail

SRC_DIR="${1:?Usage: $0 <source_directory> [title]}"
TITLE="${2:-Style Guide}"

# Remove trailing slash if present
SRC_DIR="${SRC_DIR%/}"

SUMMARY_FILE="$SRC_DIR/SUMMARY.md"

# Function to extract title from a markdown file (first H1 heading)
get_title() {
    local file="$1"
    # Get the first line starting with "# " and remove the "# " prefix
    grep -m1 '^# ' "$file" 2>/dev/null | sed 's/^# //' || basename "$file" .md
}

# Start generating SUMMARY.md
{
    echo "# $TITLE"
    echo ""

    # Find all .md files except SUMMARY.md, sorted alphabetically
    # First, check for intro.md and add it first if it exists
    if [[ -f "$SRC_DIR/intro.md" ]]; then
        title=$(get_title "$SRC_DIR/intro.md")
        echo "- [$title](intro.md)"
    fi

    # Then add all other files
    for file in "$SRC_DIR"/*.md; do
        filename=$(basename "$file")

        # Skip SUMMARY.md and intro.md (already handled)
        if [[ "$filename" == "SUMMARY.md" ]] || [[ "$filename" == "intro.md" ]]; then
            continue
        fi

        title=$(get_title "$file")
        echo "- [$title]($filename)"
    done
} > "$SUMMARY_FILE"

echo "Generated $SUMMARY_FILE with $(grep -c '^\- \[' "$SUMMARY_FILE") entries"
