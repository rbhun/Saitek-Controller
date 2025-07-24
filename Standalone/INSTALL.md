# Saitek Controller - Installation Guide

## Overview

The Saitek Controller is a comprehensive application for controlling Saitek Flight panels (Radio, Multi, Switch) with both command-line and web-based interfaces.

## Quick Start

### Option 1: Use the Launcher Script (Recommended)

1. **Download the application** to your desired location
2. **Open Terminal** and navigate to the application directory
3. **Run the launcher script:**
   ```bash
   ./launch.sh
   ```
4. **Open your web browser** and go to: `http://localhost:8080` (or `http://localhost:8081` if 8080 is in use)

### Option 2: Manual Launch

1. **Open Terminal** and navigate to the application directory
2. **Run the web interface:**
   ```bash
   ./saitek-controller-gui
   ```
3. **Open your web browser** and go to: `http://localhost:8080`

### Option 3: Command Line Interface

For command-line control of the radio panel:
```bash
./set-radio -com1a 118.25 -com1s 118.50 -com2a 121.30 -com2s 121.90
```

## System Requirements

- **macOS**: 10.15 (Catalina) or later
- **USB Ports**: For connecting Saitek panels
- **Web Browser**: Chrome, Safari, Firefox, or Edge
- **Permissions**: May require USB device access permissions

## Supported Hardware

- **Radio Panel**: Saitek Flight Radio Panel (Product ID: 0x0D05)
- **Multi Panel**: Saitek Flight Multi Panel (Product ID: 0x0D06)
- **Switch Panel**: Saitek Flight Switch Panel (Product ID: 0x0D07)

## Features

### Web Interface
- **Radio Panel Control**: Set COM1/COM2 active/standby frequencies
- **Multi Panel Control**: Set displays and button LED states
- **Switch Panel Control**: Control landing gear indicator lights
- **Real-time Status**: Monitor connection status of all panels
- **Responsive Design**: Works on desktop and mobile devices

### Command Line Interface
- **Radio Panel**: Direct frequency setting via command line
- **FIP Panels**: Flight Instrument Panel support
- **Multi Panel**: Display and LED control
- **Switch Panel**: Landing gear light control

## Troubleshooting

### USB Device Not Found
1. **Check connections**: Ensure panels are properly connected via USB
2. **Check permissions**: macOS may require USB device access permissions
3. **Try reconnecting**: Unplug and reconnect the USB cables
4. **Check device IDs**: Verify the correct vendor/product IDs are being used

### Web Interface Not Loading
1. **Check port**: Ensure port 8080 (or 8081) is not in use by another application
2. **Check firewall**: Ensure your firewall allows local connections
3. **Try different browser**: Test with Chrome, Safari, or Firefox

### Radio Panel Display Issues
1. **Check frequency format**: Use standard aviation format (e.g., "118.25")
2. **Verify encoding**: The application now supports all 5 digits with decimal points
3. **Test with command line**: Try the `set-radio` command to verify hardware communication

## File Structure

```
saitek-controller/
├── saitek-controller-gui    # Main web interface application
├── set-radio               # Command-line radio panel control
├── saitek-controller       # Main application (FIP support)
├── launch.sh               # Launcher script
├── assets/                 # Application assets
├── README.md              # Main documentation
└── INSTALL.md             # This installation guide
```

## Advanced Usage

### Custom Port
```bash
./saitek-controller-gui -port 9090
```

### Radio Panel Command Line
```bash
# Set all frequencies
./set-radio -com1a 118.25 -com1s 118.50 -com2a 121.30 -com2s 121.90

# Set individual frequencies
./set-radio -com1a 118.25
```

### FIP Panel Support
```bash
./saitek-controller -panel fip -instrument artificial_horizon
```

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Review the README.md file for detailed documentation
3. Check the logs in the terminal for error messages

## Version Information

- **Version**: 1.0
- **Last Updated**: July 2025
- **Supported Platforms**: macOS 10.15+ 