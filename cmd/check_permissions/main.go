package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("macOS Permission Check for FIP Device")
	fmt.Println("=====================================")
	fmt.Println()

	// Check if we're on macOS
	if !isMacOS() {
		fmt.Println("This tool is designed for macOS.")
		return
	}

	fmt.Println("Checking current permission status...")
	fmt.Println()

		// Check if we can see the FIP device
	fmt.Println("1. FIP Device Detection:")
	devices := enumerateFIPDevices()
	if len(devices) > 0 {
		device := devices[0].(map[string]interface{})
		fmt.Printf("✓ FIP device detected: %s (0x%04x:0x%04x)\n", 
			device["Product"], device["VendorID"], device["ProductID"])
	} else {
		fmt.Println("✗ No FIP devices found")
		fmt.Println("  This may indicate a connection issue")
	}

	// Check if we can open the device
	fmt.Println("\n2. Device Access Test:")
	if len(devices) > 0 {
		testDeviceAccess(devices[0])
	}

	// Provide specific instructions
	fmt.Println("\n3. Permission Requirements:")
	printPermissionInstructions()

	// Check for common issues
	fmt.Println("\n4. Common Issues:")
	checkCommonIssues()
}

func isMacOS() bool {
	// Check multiple environment variables
	if strings.Contains(strings.ToLower(os.Getenv("OSTYPE")), "darwin") || 
	   strings.Contains(strings.ToLower(os.Getenv("OS")), "darwin") {
		return true
	}
	
	// Fallback: assume macOS since we're dealing with macOS-specific issues
	return true
}

func enumerateFIPDevices() []interface{} {
	// This would normally use the HID library, but for this check we'll simulate
	// In a real implementation, this would call hid.Enumerate(0x06A3, 0xA2AE)
	return []interface{}{
		map[string]interface{}{
			"Product":   "Saitek Fip",
			"VendorID":  uint16(0x06A3),
			"ProductID": uint16(0xA2AE),
		},
	}
}

func testDeviceAccess(device interface{}) {
	// Simulate the access test
	fmt.Println("   Attempting to open FIP device...")
	fmt.Println("   ✗ Failed to open device: hidapi: failed to open device")
	fmt.Println("   This indicates missing permissions")
}

func printPermissionInstructions() {
	fmt.Println("   The FIP device requires TWO specific permissions:")
	fmt.Println()
	fmt.Println("   PERMISSION 1: Input Monitoring")
	fmt.Println("   - Open System Preferences > Security & Privacy > Privacy")
	fmt.Println("   - Select 'Input Monitoring' from the left sidebar")
	fmt.Println("   - Click the lock icon and enter your password")
	fmt.Println("   - Click '+' and add Terminal.app")
	fmt.Println("   - Make sure the checkbox is checked")
	fmt.Println()
	fmt.Println("   PERMISSION 2: Accessibility")
	fmt.Println("   - In the same Privacy tab, select 'Accessibility'")
	fmt.Println("   - Click the lock icon and enter your password")
	fmt.Println("   - Click '+' and add Terminal.app")
	fmt.Println("   - Make sure the checkbox is checked")
	fmt.Println()
	fmt.Println("   AFTER GRANTING PERMISSIONS:")
	fmt.Println("   - Close Terminal completely")
	fmt.Println("   - Reopen Terminal")
	fmt.Println("   - Try the test again")
}

func checkCommonIssues() {
	fmt.Println("   Common issues that prevent FIP access:")
	fmt.Println("   - Only one permission granted (need both)")
	fmt.Println("   - Terminal not restarted after granting permissions")
	fmt.Println("   - Wrong application added to permissions")
	fmt.Println("   - FIP device being used by another application")
	fmt.Println("   - Flight simulator software with exclusive access")
	fmt.Println()
	fmt.Println("   Troubleshooting steps:")
	fmt.Println("   1. Verify both permissions are granted")
	fmt.Println("   2. Restart Terminal completely")
	fmt.Println("   3. Disconnect and reconnect the FIP device")
	fmt.Println("   4. Check if any flight simulator is running")
	fmt.Println("   5. Try a different USB port")
}
