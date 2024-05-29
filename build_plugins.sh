#!/bin/sh

set -e

PLUGIN_DIR="./plugins"
OUTPUT_DIR="./plugins_dir"

# Create the output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Iterate over each plugin directory
for dir in "$PLUGIN_DIR"/*; do
    if [ -d "$dir" ]; then
        plugin_name=$(basename "$dir")
        go build -o "$OUTPUT_DIR/${plugin_name}_plugin.so" -buildmode=plugin "$dir/${plugin_name}_plugin.go"
        echo "Built ${plugin_name}_plugin.so"
    fi
done
