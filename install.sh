#!/bin/bash

set -e

LOGFILE="$HOME/.support/install.log"
SUPPORT_DIR="$HOME/.support"
PLUGINS_DIR="$SUPPORT_DIR/plugins"
CONFIG_FILE="$SUPPORT_DIR/config.yaml"

# Create the support directory
echo "Creating the support directory..."
mkdir -p "$SUPPORT_DIR"

echo "Starting installation..." | tee -a "$LOGFILE"

# Check for required dependencies
echo "Checking for required dependencies..." | tee -a "$LOGFILE"
for cmd in git go; do
    if ! command -v $cmd &> /dev/null; then
        echo "Error: $cmd is not installed. Please install it and try again." | tee -a "$LOGFILE"
        exit 1
    fi
done

# Clone the repository
echo "Cloning the repository..." | tee -a "$LOGFILE"
if [ -d "support" ]; then
    echo "Removing existing support directory..." | tee -a "$LOGFILE"
    rm -rf support
fi

git clone https://github.com/Codycody31/support.git | tee -a "$LOGFILE"
cd support

# Build all plugins
echo "Building plugins..." | tee -a "$LOGFILE"
chmod +x build_plugins.sh
./build_plugins.sh | tee -a "$LOGFILE"

# Build the core application
echo "Building the core application..." | tee -a "$LOGFILE"
go build -o dist/support main.go | tee -a "$LOGFILE"

# Create necessary directories
echo "Creating necessary directories..." | tee -a "$LOGFILE"
mkdir -p "$PLUGINS_DIR"

# Create necessary directories
echo "Creating necessary directories..." | tee -a "$LOGFILE"
mkdir -p "$PLUGINS_DIR"

# Move plugins
echo "Moving plugins to $PLUGINS_DIR..." | tee -a "$LOGFILE"

# Check if files already exist and ask for overwrite permission
if [ "$(ls -A $PLUGINS_DIR)" ]; then
    read -p "Plugins directory already contains files. Do you want to overwrite them? (y/n): " choice
    case "$choice" in 
      y|Y ) echo "Overwriting existing files..." | tee -a "$LOGFILE";;
      n|N ) echo "Skipping plugin move to avoid overwrite." | tee -a "$LOGFILE"; exit 0;;
      * ) echo "Invalid choice. Exiting." | tee -a "$LOGFILE"; exit 1;;
    esac
fi

mv dist/plugins/* "$PLUGINS_DIR/"

# Move the core application
echo "Moving the core application to $HOME/.local/bin..." | tee -a "$LOGFILE"
mkdir -p "$HOME/.local/bin"
mv dist/support "$HOME/.local/bin/"

# Update YAML configuration file
echo "Updating the YAML configuration file..." | tee -a "$LOGFILE"
if [ ! -f "$CONFIG_FILE" ]; then
	echo "No existing configuration file found. Creating a new one..." | tee -a "$LOGFILE"
    echo "plugins_dir: \"/etc/support/plugins\"" | sudo tee "$CONFIG_FILE"
fi

if grep -q "plugins_dir" "$CONFIG_FILE"; then
    echo "Updating plugins_dir in configuration file..." | tee -a "$LOGFILE"
    sed -i 's|plugins_dir:.*|plugins_dir: "'$PLUGINS_DIR'"|' "$CONFIG_FILE"
else
    echo "Adding plugins_dir to configuration file..." | tee -a "$LOGFILE"
    echo "plugins_dir: \"$PLUGINS_DIR\"" >> "$CONFIG_FILE"
fi

echo "Installation completed successfully!" | tee -a "$LOGFILE"
