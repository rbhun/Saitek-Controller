#!/bin/bash

# Saitek Controller Launcher Script
# This script launches the Saitek Controller application

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$SCRIPT_DIR"

echo "Saitek Controller Launcher"
echo "=========================="

# Check if we're in the right directory
if [ ! -f "$APP_DIR/saitek-controller-gui" ]; then
    echo "Error: saitek-controller-gui not found in $APP_DIR"
    echo "Please run this script from the application directory"
    exit 1
fi

# Check if port 8080 is available, otherwise use 8081
PORT=8080
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "Port 8080 is in use, using port 8081 instead"
    PORT=8081
fi

# Get the local IP address
LOCAL_IP=$(ifconfig | grep "inet " | grep -v 127.0.0.1 | awk '{print $2}' | head -1)

echo "Starting Saitek Controller on port $PORT..."
echo ""
echo "Local access:"
echo "  http://localhost:$PORT"
echo "  http://127.0.0.1:$PORT"
echo ""
echo "Network access (from other computers):"
echo "  http://$LOCAL_IP:$PORT"
echo ""
echo "Press Ctrl+C to stop the application"
echo ""

# Launch the application with network access
"$APP_DIR/saitek-controller-gui" -port "$PORT" -host "0.0.0.0" 