# Saitek Controller - Truly Standalone Application

## 🚀 **100% Standalone - No Dependencies Required!**

This folder contains the complete, **truly standalone** Saitek Controller application that can be moved to any Mac and run without installing any additional packages or libraries.

## ✅ **What Makes This Special:**

- **No Homebrew required** - All dependencies are bundled
- **No manual library installation** - Everything is self-contained
- **Works on any Mac** - macOS 10.15+ compatible
- **Zero configuration** - Just copy and run

## 📦 **Contents**

- **`saitek-controller-gui`** - Main web interface application
- **`set-radio`** - Command-line radio panel control
- **`saitek-controller`** - Main application with FIP support
- **`launch.sh`** - Easy launcher script
- **`libs/`** - Bundled dependencies (libusb)
- **`assets/`** - Application assets and images
- **`INSTALL.md`** - Installation and usage guide
- **`NETWORK_ACCESS.md`** - Network access guide

## 🎯 **Quick Start (Any Mac)**

### Option 1: Use the Launcher (Recommended)
```bash
./launch.sh
```
Then open your browser to `http://localhost:8080`

### Option 2: Direct Web Interface
```bash
./saitek-controller-gui
```

### Option 3: Network Access (from other computers)
```bash
./saitek-controller-gui -host 0.0.0.0 -port 8080
```
Then access from other computers at `http://YOUR_IP:8080`

### Option 4: Command Line Radio Control
```bash
./set-radio -com1a 118.25 -com1s 118.50 -com2a 121.30 -com2s 121.90
```

## ✅ **What's Working**

### Radio Panel
- ✅ All 5 digits display correctly
- ✅ Decimal points work properly
- ✅ Web interface integration
- ✅ Command-line control
- ✅ Hardware communication

### Web Interface
- ✅ Radio Panel Control - Set COM1/COM2 frequencies
- ✅ Multi Panel Control - Displays and LED control
- ✅ Switch Panel Control - Landing gear lights
- ✅ Real-time status monitoring
- ✅ Responsive design
- ✅ Network access from other computers

## 🔧 **System Requirements**

- **macOS**: 10.15 (Catalina) or later
- **USB Ports**: For connecting Saitek panels
- **Web Browser**: Chrome, Safari, Firefox, or Edge
- **Permissions**: May require USB device access permissions

## 📋 **Supported Hardware**

- **Radio Panel**: Saitek Flight Radio Panel (Product ID: 0x0D05)
- **Multi Panel**: Saitek Flight Multi Panel (Product ID: 0x0D06)
- **Switch Panel**: Saitek Flight Switch Panel (Product ID: 0x0D07)

## 🌐 **Network Access**

You can access the web interface from other computers on the same network:

1. **Start with network access:**
   ```bash
   ./saitek-controller-gui -host 0.0.0.0 -port 8080
   ```

2. **Access from other computers:**
   ```
   http://YOUR_COMPUTER_IP:8080
   ```

3. **Mobile devices** can also access the interface!

See `NETWORK_ACCESS.md` for detailed instructions.

## 🚀 **Moving This Application**

You can copy this entire folder to any Mac and run it:

```bash
# Copy to your desired location
cp -r Standalone-Bundled/ /path/to/your/desired/location/

# Navigate to the new location
cd /path/to/your/desired/location/

# Run the application
./launch.sh
```

## 📖 **Documentation**

- **`INSTALL.md`** - Complete installation and usage guide
- **`NETWORK_ACCESS.md`** - Network access instructions
- **`STANDALONE_SUMMARY.md`** - Detailed feature summary

## 🎉 **Ready for Distribution!**

This bundled application is:
- ✅ **Completely self-contained**
- ✅ **No external dependencies**
- ✅ **Ready for distribution**
- ✅ **Production-ready**
- ✅ **Network accessible**
- ✅ **Works on any Mac without setup**

**You can now move this entire folder to any Mac and run it immediately!** 🚀 