# Saitek Flight Panel Controller

A comprehensive control software for Saitek/Logitech flight panels designed for filmmaking applications.

## Supported Panels

- **Flight Instrument Panels (FIP)** - Display custom images and instrument data
- **Radio Panel** - Control 7-segment displays and button backlights
- **Switch Panel** - Control landing gear LEDs and read switch states
- **Multi Panel** - Control 7-segment displays and button backlights

## Features

- **Image Display**: Display custom images on Flight Instrument Panels
- **7-Segment Control**: Send data to any 7-segment displays
- **LED Control**: Control button backlights and indicator LEDs
- **Switch Reading**: Read switch states from all panels
- **Modular Design**: Separate modules for each panel type
- **Cross-Platform**: Works on Windows, macOS, and Linux

## Multi Panel Features

The Multi Panel implementation supports:

- **Dual 7-Segment Displays**: Control two 5-digit displays (top and bottom rows)
- **Button LED Control**: Individual control of 8 button backlights (AP, HDG, NAV, IAS, ALT, VS, APR, REV)
- **Switch Reading**: Monitor all switches, buttons, and encoders
- **USB Communication**: Direct USB control using the same protocol as the fpanels library
- **Mock Mode**: Test functionality without physical hardware

### Multi Panel Button LEDs

The button LEDs are controlled with a single byte where each bit represents a button:

- Bit 0 (0x01): AP button
- Bit 1 (0x02): HDG button  
- Bit 2 (0x04): NAV button
- Bit 3 (0x08): IAS button
- Bit 4 (0x10): ALT button
- Bit 5 (0x20): VS button
- Bit 6 (0x40): APR button
- Bit 7 (0x80): REV button

Example: `0x0F` would light up AP, HDG, NAV, and IAS buttons.

## Project Structure

```
saitek-controller/
├── cmd/
│   └── main.go              # Main application entry point
├── internal/
│   ├── fip/                 # Flight Instrument Panel control
│   ├── radio/               # Radio panel control
│   ├── switch/              # Switch panel control
│   ├── multi/               # Multi panel control
│   └── usb/                 # USB communication utilities
├── assets/                  # Images and resources
└── examples/                # Example usage and demos
```

## USB Product IDs

- Flight Instrument Panel: `06a3:0d06`
- Radio Panel: `06a3:0d05`
- Switch Panel: `06a3:0d67`
- Multi Panel: `06a3:0d06`

## Usage

### Command Line Interface

```bash
# Build the application
go build -o saitek-controller cmd/main.go

# Run the application
./saitek-controller

# Multi Panel Examples
./saitek-controller -panel multi -top "250" -bottom "3000" -leds 0x01
./saitek-controller -panel multi -top "120" -bottom "5000" -leds 0x0F

# Radio Panel Examples
./saitek-controller -panel radio -com1a "118.00" -com1s "118.50" -com2a "121.30" -com2s "121.90"

# FIP Examples
./saitek-controller -panel fip -instrument artificial_horizon
./saitek-controller -panel fip -image path/to/image.png
```

### Web-Based GUI

A comprehensive web-based GUI is available for easy control of all panels:

```bash
# Build the GUI application
make build-gui

# Run the GUI application
make run-gui

# Or run directly
go run cmd/saitek-controller-gui/main.go

# Open your browser to http://localhost:8080
```

The GUI provides:
- **Radio Panel Control**: Set COM1 and COM2 frequencies
- **Multi Panel Control**: Set displays and button LEDs
- **Switch Panel Control**: Control landing gear lights
- **Real-time Status**: Monitor panel connections
- **Modern Interface**: Responsive web design

For more details, see [GUI Documentation](cmd/saitek-controller-gui/README.md).

## Dependencies

- Go 1.21+
- USB HID libraries
- Image processing libraries

## License

MIT License 