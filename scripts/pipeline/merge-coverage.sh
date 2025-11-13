#!/bin/bash
set -euo pipefail

# Merges multiple Go coverage files into one
# Usage: ./merge-coverage.sh <input-dir> <output-file>

INPUT_DIR="${1:-}"
OUTPUT_FILE="${2:-}"

if [[ -z "$INPUT_DIR" || -z "$OUTPUT_FILE" ]]; then
    echo "Usage: $0 <input-dir> <output-file>"
    exit 1
fi

if [[ ! -d "$INPUT_DIR" ]]; then
    echo "Error: input directory not found: $INPUT_DIR"
    exit 1
fi

echo "Searching for coverage files in ${INPUT_DIR}..."

# Find all .txt files recursively
coverage_files=$(find "$INPUT_DIR" -type f -name "*.txt")

if [[ -z "$coverage_files" ]]; then
    echo "⚠️  No coverage files found, creating empty output"
    echo "mode: atomic" > "$OUTPUT_FILE"
    exit 0
fi

file_count=$(echo "$coverage_files" | wc -l | tr -d ' ')
echo "Found ${file_count} coverage files to merge"

# Write the mode line
echo "mode: atomic" > "$OUTPUT_FILE"

total_lines=0

# Merge all coverage files
while IFS= read -r file; do
    if [[ ! -f "$file" ]]; then
        continue
    fi

    # Skip the first line (mode: atomic) and empty lines
    lines=$(tail -n +2 "$file" | grep -v '^[[:space:]]*$' || true)

    if [[ -n "$lines" ]]; then
        echo "$lines" >> "$OUTPUT_FILE"
        line_count=$(echo "$lines" | wc -l | tr -d ' ')
        total_lines=$((total_lines + line_count))
        echo "  Merged $(basename "$file") (${line_count} lines)"
    fi
done <<< "$coverage_files"

echo "Total coverage lines: ${total_lines}"
echo "✅ Successfully merged coverage files to ${OUTPUT_FILE}"
