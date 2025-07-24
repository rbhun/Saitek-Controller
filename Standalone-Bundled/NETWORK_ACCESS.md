# Network Access Guide

## 🌐 Accessing the Web Interface from Other Computers

Your Saitek Controller web interface can now be accessed from other computers on the same network!

## 🚀 Quick Start

### 1. Start the Application with Network Access
```bash
./saitek-controller-gui -host 0.0.0.0 -port 8080
```

### 2. Find Your Computer's IP Address
The launcher script will automatically show you the network URL, or you can find it manually:
```bash
ifconfig | grep "inet " | grep -v 127.0.0.1
```

### 3. Access from Other Computers
Open a web browser on any computer on the same network and go to:
```
http://YOUR_COMPUTER_IP:8080
```

## 📋 Step-by-Step Instructions

### Option 1: Use the Launcher Script (Recommended)
```bash
./launch.sh
```
The launcher will automatically:
- Detect available ports
- Show your local IP address
- Enable network access
- Display all access URLs

### Option 2: Manual Network Access
```bash
# Start with network access enabled
./saitek-controller-gui -host 0.0.0.0 -port 8080

# Or use a different port
./saitek-controller-gui -host 0.0.0.0 -port 9090
```

## 🔍 Finding Your IP Address

### macOS
```bash
ifconfig | grep "inet " | grep -v 127.0.0.1
```

### Windows
```bash
ipconfig
```

### Linux
```bash
ip addr show
```

## 🌐 Network Access URLs

When you start the application, you'll see output like this:
```
Starting Saitek Controller on port 8080...

Local access:
  http://localhost:8080
  http://127.0.0.1:8080

Network access (from other computers):
  http://192.168.1.100:8080
```

## 📱 Access from Mobile Devices

You can also access the web interface from your phone or tablet:

1. **Connect your mobile device** to the same WiFi network
2. **Open your mobile browser**
3. **Navigate to** `http://YOUR_COMPUTER_IP:8080`
4. **Use the touch-friendly interface** to control your Saitek panels

## 🔒 Security Considerations

### Local Network Only
- The web interface is designed for local network use
- No authentication is required
- Anyone on your network can access it

### Firewall Settings
You may need to allow the application through your firewall:

**macOS:**
- Go to System Preferences > Security & Privacy > Firewall
- Click "Firewall Options" and add the application

**Windows:**
- Go to Control Panel > System and Security > Windows Defender Firewall
- Click "Allow an app or feature through Windows Defender Firewall"

## 🛠️ Troubleshooting

### Can't Connect from Other Computers
1. **Check firewall settings** - Allow the application through your firewall
2. **Verify IP address** - Make sure you're using the correct IP
3. **Check network** - Ensure both computers are on the same network
4. **Try different port** - Use `-port 9090` if 8080 is blocked

### Connection Refused
1. **Verify the application is running** with network access enabled
2. **Check the port** - Make sure the port isn't being used by another application
3. **Try localhost first** - Test with `http://localhost:8080` to verify the app works

### Mobile Device Can't Connect
1. **Check WiFi network** - Ensure mobile device is on the same network
2. **Try different browser** - Test with Chrome, Safari, or Firefox
3. **Check mobile firewall** - Some mobile devices have additional security

## 📋 Example Usage

### Start with Network Access
```bash
./saitek-controller-gui -host 0.0.0.0 -port 8080
```

### Access from Different Devices
- **Desktop computer**: `http://192.168.1.100:8080`
- **Laptop**: `http://192.168.1.100:8080`
- **iPhone/iPad**: `http://192.168.1.100:8080`
- **Android device**: `http://192.168.1.100:8080`

## 🎯 Benefits of Network Access

- **Remote control** - Control panels from anywhere in your home/office
- **Mobile access** - Use your phone or tablet to control panels
- **Multiple users** - Several people can access the interface simultaneously
- **Convenience** - No need to be at the computer with the panels

## 🚀 Ready to Use!

Your Saitek Controller now supports full network access, allowing you to control your flight panels from any device on your local network! 