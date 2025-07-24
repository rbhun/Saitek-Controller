#!/bin/bash

# Saitek Controller Deployment Script
# This script helps you deploy the standalone application to any location

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_NAME="Saitek-Controller"

echo "Saitek Controller Deployment Script"
echo "=================================="
echo ""

# Check if destination is provided
if [ -z "$1" ]; then
    echo "Usage: ./deploy.sh <destination_path>"
    echo ""
    echo "Examples:"
    echo "  ./deploy.sh ~/Desktop/Saitek-Controller"
    echo "  ./deploy.sh /Applications/Saitek-Controller"
    echo "  ./deploy.sh ~/Documents/Flight-Sim/Saitek-Controller"
    echo ""
    echo "The application will be copied to the specified location."
    exit 1
fi

DESTINATION="$1"

echo "Deploying Saitek Controller to: $DESTINATION"
echo ""

# Create destination directory if it doesn't exist
if [ ! -d "$DESTINATION" ]; then
    echo "Creating destination directory..."
    mkdir -p "$DESTINATION"
fi

# Copy all files
echo "Copying application files..."
cp -r "$SCRIPT_DIR"/* "$DESTINATION/"

# Make executables
echo "Setting executable permissions..."
chmod +x "$DESTINATION/saitek-controller-gui"
chmod +x "$DESTINATION/saitek-controller"
chmod +x "$DESTINATION/set-radio"
chmod +x "$DESTINATION/launch.sh"

echo ""
echo "✅ Deployment complete!"
echo ""
echo "To run the application:"
echo "  cd \"$DESTINATION\""
echo "  ./launch.sh"
echo ""
echo "Or open your browser to: http://localhost:8080"
echo ""
echo "The application is now ready to use on this Mac!"
