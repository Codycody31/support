#!/bin/sh

set -e

PLUGIN_DIR="./plugins"
OUTPUT_DIR="./dist/plugins"

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Iterate over each plugin directory
for dir in "$PLUGIN_DIR"/*; do
    if [ -d "$dir" ]; then
        plugin_name=$(basename "$dir")
        go build -o "$OUTPUT_DIR/${plugin_name}.so" -buildmode=plugin "$dir/${plugin_name}.go"
        echo "Built ${plugin_name}.so"
    fi
done
