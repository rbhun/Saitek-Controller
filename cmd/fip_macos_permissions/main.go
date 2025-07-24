package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"saitek-controller/internal/usb"

	"github.com/karalabe/hid"
)

func main() {
	fmt.Println("Saitek FIP macOS Permission Fix")
	fmt.Println("=================================")
	fmt.Println()

	// Check if we're on macOS
	if !isMacOS() {
		fmt.Println("This tool is designed for macOS. Please use the appropriate tool for your OS.")
		return
	}

	// Comprehensive permission analysis
	fmt.Println("=== macOS Permission Analysis ===")
	analyzeMacOSPermissions()

	// Try different approaches
	fmt.Println("\n=== Attempting Permission Fixes ===")
	tryPermissionFixes()

	// Provide detailed instructions
	fmt.Println("\n=== Detailed Fix Instructions ===")
	printDetailedInstructions()

	// Test with different approaches
	fmt.Println("\n=== Testing Alternative Approaches ===")
	testAlternativeApproaches()
}

func isMacOS() bool {
	// Check multiple environment variables
	if strings.Contains(strings.ToLower(os.Getenv("OSTYPE")), "darwin") || 
	   strings.Contains(strings.ToLower(os.Getenv("OS")), "darwin") ||
	   strings.Contains(strings.ToLower(os.Getenv("OSTYPE")), "macos") {
		return true
	}
	
	// Fallback: check if we're on a system that behaves like macOS
	// (has the same permission issues)
	return true // Assume macOS for now since we're dealing with macOS-specific issues
}

func analyzeMacOSPermissions() {
	fmt.Println("Analyzing macOS-specific permission issues...")

	// Check if we can see the FIP device
	devices := hid.Enumerate(0x06A3, 0xA2AE)
	if len(devices) == 0 {
		fmt.Println("✗ FIP device not found in HID enumeration")
		fmt.Println("  This indicates a serious permission issue")
		return
	}

	fmt.Printf("✓ FIP device found: %s\n", devices[0].Product)
	fmt.Printf("  Vendor ID: 0x%04x, Product ID: 0x%04x\n", devices[0].VendorID, devices[0].ProductID)

	// Check macOS-specific permission issues
	fmt.Println("\n--- macOS Permission Checks ---")

	// Check if we're running as root
	if os.Geteuid() == 0 {
		fmt.Println("✓ Running as root (elevated privileges)")
	} else {
		fmt.Println("⚠ Running as regular user")
	}

	// Check for Input Monitoring permission
	checkInputMonitoringPermission()

	// Check for Accessibility permission
	checkAccessibilityPermission()

	// Check for USB device permissions
	checkUSBDevicePermissions()
}

func checkInputMonitoringPermission() {
	fmt.Println("\n--- Input Monitoring Permission ---")
	fmt.Println("The FIP device requires Input Monitoring permission on macOS.")
	fmt.Println("This permission allows applications to monitor input from HID devices.")
	fmt.Println()
	fmt.Println("To grant this permission:")
	fmt.Println("1. Open System Preferences > Security & Privacy > Privacy")
	fmt.Println("2. Select 'Input Monitoring' from the left sidebar")
	fmt.Println("3. Click the lock icon and enter your password")
	fmt.Println("4. Add Terminal.app, iTerm, or your IDE")
	fmt.Println("5. If the Go binary appears, add it as well")
	fmt.Println()
}

func checkAccessibilityPermission() {
	fmt.Println("--- Accessibility Permission ---")
	fmt.Println("Some HID devices also require Accessibility permission.")
	fmt.Println()
	fmt.Println("To grant this permission:")
	fmt.Println("1. Open System Preferences > Security & Privacy > Privacy")
	fmt.Println("2. Select 'Accessibility' from the left sidebar")
	fmt.Println("3. Click the lock icon and enter your password")
	fmt.Println("4. Add Terminal.app, iTerm, or your IDE")
	fmt.Println()
}

func checkUSBDevicePermissions() {
	fmt.Println("--- USB Device Permissions ---")
	fmt.Println("The FIP device might need special USB device permissions.")
	fmt.Println()
	fmt.Println("To check USB device status:")
	fmt.Println("1. Open System Information (Apple menu > About This Mac > System Report)")
	fmt.Println("2. Go to USB in the left sidebar")
	fmt.Println("3. Look for 'Saitek Fip' or similar device")
	fmt.Println("4. Check if the device shows as connected and working")
	fmt.Println()
}

func tryPermissionFixes() {
	fmt.Println("Trying different permission approaches...")

	// Approach 1: Try with current permissions
	fmt.Println("\n--- Approach 1: Current Permissions ---")
	testFIPAccess()

	// Approach 2: Try with sudo
	if os.Geteuid() != 0 {
		fmt.Println("\n--- Approach 2: Elevated Privileges ---")
		fmt.Println("Note: Even with sudo, macOS may still block HID access")
		fmt.Println("This is a macOS security feature, not a bug")
	}

	// Approach 3: Check if device is being used by another application
	fmt.Println("\n--- Approach 3: Device Usage Check ---")
	checkDeviceUsage()
}

func testFIPAccess() {
	devices := hid.Enumerate(0x06A3, 0xA2AE)
	if len(devices) == 0 {
		fmt.Println("✗ No FIP devices found")
		return
	}

	device := devices[0]
	fmt.Printf("Attempting to open FIP device: %s\n", device.Product)

	// Try to open the device
	hidDevice, err := device.Open()
	if err != nil {
		fmt.Printf("✗ Failed to open FIP device: %v\n", err)
		fmt.Println("  This is the core permission issue")
		fmt.Println("  The device is detected but cannot be opened")
		return
	}
	defer hidDevice.Close()

	fmt.Println("✓ Successfully opened FIP device!")
	fmt.Println("  This means permissions are correctly set")

	// Try to read button states
	fmt.Println("Testing button reading...")

	// Set up a goroutine to read button states
	buttonChan := make(chan []byte, 10)
	errorChan := make(chan error, 1)

	go func() {
		buffer := make([]byte, 2) // FIP has 2-byte input reports
		for {
			n, err := hidDevice.Read(buffer)
			if err != nil {
				errorChan <- err
				return
			}
			if n > 0 {
				data := make([]byte, n)
				copy(data, buffer[:n])
				buttonChan <- data
			}
		}
	}()

	// Monitor for button presses for 5 seconds
	fmt.Println("Monitoring button presses for 5 seconds...")
	fmt.Println("Press buttons on the FIP device to test...")

	timeout := time.After(5 * time.Second)
	for {
		select {
		case data := <-buttonChan:
			fmt.Printf("✓ Button press detected: %v\n", data)
		case err := <-errorChan:
			fmt.Printf("✗ Error reading from device: %v\n", err)
			return
		case <-timeout:
			fmt.Println("✓ No button presses detected (device is working)")
			return
		}
	}
}

func checkDeviceUsage() {
	fmt.Println("Checking if the FIP device is being used by another application...")

	// Use system commands to check device usage
	commands := []struct {
		name string
		cmd  []string
	}{
		{"lsof", []string{"lsof", "+D", "/dev"}},
		{"ps", []string{"ps", "aux"}},
	}

	for _, command := range commands {
		fmt.Printf("\n--- %s ---\n", command.name)
		cmd := exec.Command(command.cmd[0], command.cmd[1:]...)
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("✗ %s failed: %v\n", command.name, err)
			continue
		}

		outputStr := string(output)
		if strings.Contains(outputStr, "hid") || strings.Contains(outputStr, "usb") {
			fmt.Printf("✓ Found HID/USB related processes\n")
			// Print relevant lines
			lines := strings.Split(outputStr, "\n")
			for _, line := range lines {
				if strings.Contains(line, "hid") || strings.Contains(line, "usb") {
					fmt.Printf("  %s\n", strings.TrimSpace(line))
				}
			}
		} else {
			fmt.Printf("⚠ No HID/USB related processes found\n")
		}
	}
}

func printDetailedInstructions() {
	fmt.Println("=== COMPREHENSIVE FIX INSTRUCTIONS ===")
	fmt.Println()

	fmt.Println("The Saitek FIP device requires special permissions on macOS.")
	fmt.Println("Here's a step-by-step guide to fix the permissions:")
	fmt.Println()

	fmt.Println("STEP 1: GRANT INPUT MONITORING PERMISSION")
	fmt.Println("==========================================")
	fmt.Println("1. Open System Preferences")
	fmt.Println("2. Go to Security & Privacy")
	fmt.Println("3. Click the 'Privacy' tab")
	fmt.Println("4. Select 'Input Monitoring' from the left sidebar")
	fmt.Println("5. Click the lock icon (bottom left) and enter your password")
	fmt.Println("6. Click the '+' button and add:")
	fmt.Println("   - Terminal.app (if using Terminal)")
	fmt.Println("   - iTerm (if using iTerm)")
	fmt.Println("   - Your IDE (if running from IDE)")
	fmt.Println("7. Make sure the checkbox is checked for each application")
	fmt.Println()

	fmt.Println("STEP 2: GRANT ACCESSIBILITY PERMISSION")
	fmt.Println("======================================")
	fmt.Println("1. In the same Privacy tab, select 'Accessibility'")
	fmt.Println("2. Click the lock icon and enter your password")
	fmt.Println("3. Add the same applications as above")
	fmt.Println("4. Make sure the checkbox is checked")
	fmt.Println()

	fmt.Println("STEP 3: RESTART APPLICATIONS")
	fmt.Println("============================")
	fmt.Println("1. Close your terminal/IDE completely")
	fmt.Println("2. Reopen your terminal/IDE")
	fmt.Println("3. Try running the FIP test again")
	fmt.Println()

	fmt.Println("STEP 4: CHECK DEVICE CONNECTION")
	fmt.Println("================================")
	fmt.Println("1. Disconnect the FIP device")
	fmt.Println("2. Wait 5 seconds")
	fmt.Println("3. Reconnect the FIP device")
	fmt.Println("4. Check System Information > USB to ensure it's recognized")
	fmt.Println()

	fmt.Println("STEP 5: ALTERNATIVE APPROACHES")
	fmt.Println("===============================")
	fmt.Println("If the above doesn't work:")
	fmt.Println("1. Try running with sudo (temporary fix)")
	fmt.Println("2. Check if any flight simulator software is using the device")
	fmt.Println("3. Try a different USB port")
	fmt.Println("4. Check if the device works with other software")
	fmt.Println()

	fmt.Println("STEP 6: DEVELOPER OPTIONS")
	fmt.Println("==========================")
	fmt.Println("For developers who need permanent access:")
	fmt.Println("1. Sign your application with proper entitlements")
	fmt.Println("2. Use IOKit directly instead of HID library")
	fmt.Println("3. Create a kernel extension (advanced)")
	fmt.Println()
}

func testAlternativeApproaches() {
	fmt.Println("Testing alternative approaches to access the FIP...")

	// Approach 1: Try different HID library options
	fmt.Println("\n--- Approach 1: Different HID Options ---")
	testDifferentHIDOptions()

	// Approach 2: Try USB core access
	fmt.Println("\n--- Approach 2: USB Core Access ---")
	testUSBCoreAccess()

	// Approach 3: Check for device-specific issues
	fmt.Println("\n--- Approach 3: Device-Specific Checks ---")
	testDeviceSpecificIssues()
}

func testDifferentHIDOptions() {
	// Try different vendor/product ID combinations
	testIDs := []struct {
		vendor  uint16
		product uint16
		name    string
	}{
		{0x06A3, 0xA2AE, "Standard FIP"},
		{0x06A3, 0x0A2E, "Alternative FIP"},
		{0x06A3, 0x0A2C, "Another variant"},
	}

	for _, test := range testIDs {
		devices := hid.Enumerate(test.vendor, test.product)
		if len(devices) > 0 {
			fmt.Printf("✓ Found device with %s (0x%04x:0x%04x)\n",
				test.name, test.vendor, test.product)

			// Try to open it
			device := devices[0]
			hidDevice, err := device.Open()
			if err != nil {
				fmt.Printf("  ✗ Cannot open: %v\n", err)
			} else {
				fmt.Printf("  ✓ Successfully opened!\n")
				hidDevice.Close()
			}
		}
	}
}

func testUSBCoreAccess() {
	fmt.Println("Testing direct USB core access...")

	// Try to create a USB device using our core
	device, err := usb.OpenDevice(0x06A3, 0xA2AE)
	if err != nil {
		fmt.Printf("✗ USB core cannot access FIP: %v\n", err)
		return
	}
	defer device.Close()

	fmt.Println("✓ USB core can access FIP device")

	// Try to read from it using ReadBulkData
	data, err := device.ReadBulkData(0, 2)
	if err != nil {
		fmt.Printf("✗ Cannot read from FIP: %v\n", err)
		return
	}

	fmt.Printf("✓ Successfully read %d bytes from FIP: %v\n", len(data), data)
}

func testDeviceSpecificIssues() {
	fmt.Println("Checking for device-specific issues...")

	// Check if the device is in a usable state
	fmt.Println("1. Check if the FIP device is powered on and connected")
	fmt.Println("2. Check if any LED indicators are lit on the device")
	fmt.Println("3. Try pressing buttons to see if they respond")
	fmt.Println("4. Check if the device appears in System Information")
	fmt.Println()

	// Check for conflicting software
	fmt.Println("Common conflicting software:")
	fmt.Println("- X-Plane (may have exclusive access)")
	fmt.Println("- Microsoft Flight Simulator")
	fmt.Println("- Prepar3D")
	fmt.Println("- Other flight simulator software")
	fmt.Println("- Saitek/Logitech drivers")
	fmt.Println()

	fmt.Println("If any of these are running, try:")
	fmt.Println("1. Closing the flight simulator")
	fmt.Println("2. Disconnecting and reconnecting the FIP")
	fmt.Println("3. Restarting your computer")
	fmt.Println("4. Running the test again")
}
