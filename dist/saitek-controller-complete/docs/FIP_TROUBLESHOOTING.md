# FIP Troubleshooting Guide

## macOS HID Device Access Issues

### Problem
The application can detect the Saitek FIP device but cannot open it due to macOS security restrictions. You'll see errors like:
```
Warning: Could not connect to physical FIP device: failed to open device: hidapi: failed to open device
Running in virtual mode only
```

### Current Status
✅ **Virtual Mode Working**: The application can display images in virtual windows
❌ **Physical Device Access**: Cannot connect to the physical FIP hardware

### Solutions to Try

#### 1. Grant Accessibility Permissions
1. Go to **System Preferences** > **Security & Privacy** > **Privacy**
2. Select **Accessibility** from the left sidebar
3. Click the lock icon to make changes
4. Add your terminal application (Terminal.app or iTerm2) to the list
5. Restart your terminal and try again

#### 2. Grant Input Monitoring Permissions
1. Go to **System Preferences** > **Security & Privacy** > **Privacy**
2. Select **Input Monitoring** from the left sidebar
3. Click the lock icon to make changes
4. Add your terminal application to the list
5. Restart your terminal and try again

#### 3. Run with Elevated Privileges
Try running the application with sudo:
```bash
sudo ./bin/saitek-controller -vendor 06a3 -product a2ae -instrument artificial_horizon
```

#### 4. Check Device Permissions
Run this command to see if the device is accessible:
```bash
ls -la /dev/hid*
```

#### 5. Alternative: Use Virtual Mode
If physical device access continues to fail, you can still use the virtual mode for testing:
```bash
./bin/saitek-controller -instrument artificial_horizon
```

### Known Issues
- The HID library uses deprecated macOS APIs (kIOMasterPortDefault)
- macOS requires explicit permission for HID device access
- Some USB hubs may block HID device access

### Workarounds
1. **Virtual Testing**: Use virtual mode for development and testing
2. **Different USB Port**: Try connecting the FIP to a different USB port
3. **Direct Connection**: Connect the FIP directly to the Mac (not through a hub)
4. **Different HID Library**: Consider using a different HID library that supports modern macOS APIs

### Current Functionality
- ✅ Image generation and display
- ✅ Virtual FIP windows
- ✅ Instrument rendering (artificial horizon, airspeed, etc.)
- ✅ USB device detection
- ❌ Physical device communication

The application is fully functional for development and testing in virtual mode. The physical device access is a macOS permission issue that can be resolved with the steps above. 