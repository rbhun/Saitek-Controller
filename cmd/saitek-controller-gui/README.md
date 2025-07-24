# Saitek Controller GUI

A comprehensive web-based GUI for controlling Saitek Flight Radio, Multi, and Switch panels.

## Features

- **Radio Panel Control**: Set COM1 and COM2 active/standby frequencies
- **Multi Panel Control**: Set display values and button LED states
- **Switch Panel Control**: Control landing gear indicator lights
- **Real-time Status**: Monitor connection status of all panels
- **Modern Web Interface**: Responsive design that works on desktop and mobile

## Quick Start

### Prerequisites

- Go 1.21 or later
- Saitek Flight panels (Radio, Multi, Switch) connected via USB
- Web browser

### Building and Running

1. **Build the application:**
   ```bash
   make build-gui
   ```

2. **Run the application:**
   ```bash
   make run-gui
   ```

3. **Or run directly:**
   ```bash
   go run cmd/saitek-controller-gui/main.go
   ```

4. **Open your web browser and navigate to:**
   ```
   http://localhost:8080
   ```

### Command Line Options

- `-port <port>`: Specify the port to listen on (default: 8080)

Example:
```bash
go run cmd/saitek-controller-gui/main.go -port 9090
```

## Usage

### Radio Panel

The Radio Panel section allows you to control the four 5-digit displays:

- **COM1 Active**: Top left display
- **COM1 Standby**: Top right display  
- **COM2 Active**: Bottom left display
- **COM2 Standby**: Bottom right display

Enter frequencies in standard aviation format (e.g., "118.00", "121.30") and click "Set Radio Display" to update the panel.

### Multi Panel

The Multi Panel section controls the two 5-digit displays and button LEDs:

- **Top Row**: Upper 5-digit display
- **Bottom Row**: Lower 5-digit display
- **Button LEDs**: Checkboxes for each button LED (AP, HDG, NAV, IAS, ALT, VS, APR, REV)

Enter values for the displays and select which button LEDs should be illuminated.

### Switch Panel

The Switch Panel section controls the landing gear indicator lights:

- **Individual Control**: Check/uncheck individual lights (Green N/L/R, Red N/L/R)
- **Preset Buttons**:
  - **Gear Down (Green)**: All green lights on, red lights off
  - **Gear Up (Red)**: All red lights on, green lights off  
  - **Gear Transition (Yellow)**: All lights on (creates yellow effect)
  - **All Lights Off**: Turn off all lights

## Panel Status

The application shows real-time connection status for each panel:

- **Green indicator**: Panel is connected and responding
- **Red indicator**: Panel is not connected or not responding

Click "Refresh Status" to update the connection status, or "Reconnect All" to attempt reconnection to all panels.

## Troubleshooting

### Panels Not Connecting

1. **Check USB connections**: Ensure all panels are properly connected via USB
2. **Check permissions**: On macOS, you may need to grant accessibility permissions
3. **Check device IDs**: The application uses standard Saitek vendor/product IDs:
   - Radio Panel: 0x06A3/0x0D05
   - Multi Panel: 0x06A3/0x0D06  
   - Switch Panel: 0x06A3/0x0D67

### Web Interface Not Loading

1. **Check port**: Ensure the port (default 8080) is not in use
2. **Check firewall**: Ensure your firewall allows connections to the port
3. **Try different browser**: Some browsers may have compatibility issues

### Display Not Updating

1. **Check panel connection**: Verify the panel status indicator is green
2. **Check input format**: Ensure frequencies and values are in the correct format
3. **Try reconnecting**: Use the "Reconnect All" button to refresh connections

## Development

### Architecture

The application consists of:

- **PanelManager**: Manages connections to all three panel types
- **Web Server**: Serves the HTML interface and REST API endpoints
- **REST API**: Provides endpoints for setting panel states and getting status

### API Endpoints

- `GET /api/status`: Get current status of all panels
- `POST /api/radio/set`: Set radio panel display
- `POST /api/multi/set`: Set multi panel display and LEDs
- `POST /api/switch/set`: Set switch panel lights
- `POST /api/connect`: Reconnect to all panels

### Adding New Features

To add new functionality:

1. **Backend**: Add methods to `PanelManager` in `main.go`
2. **API**: Add new endpoints to handle the functionality
3. **Frontend**: Add UI elements and JavaScript functions
4. **Testing**: Test with actual hardware

## License

This project is part of the Saitek Controller project. See the main project README for license information. 