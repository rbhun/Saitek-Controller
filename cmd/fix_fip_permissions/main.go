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
	fmt.Println("Saitek FIP Permission Fix Tool")
	fmt.Println("===============================")
	fmt.Println()

	// Check current permissions
	fmt.Println("=== Current Permission Status ===")
	checkCurrentPermissions()

	// Try different approaches to fix permissions
	fmt.Println("\n=== Attempting Permission Fixes ===")

	// Approach 1: Check if running as root
	if os.Geteuid() == 0 {
		fmt.Println("✓ Running as root - trying direct access...")
		testFIPAccess()
	} else {
		fmt.Println("⚠ Not running as root")
		fmt.Println("   Trying to access FIP device...")
		testFIPAccess()
	}

	// Approach 2: Provide detailed instructions
	fmt.Println("\n=== Permission Fix Instructions ===")
	printPermissionInstructions()

	// Approach 3: Test with different permission methods
	fmt.Println("\n=== Testing Alternative Access Methods ===")
	testAlternativeMethods()
}

func checkCurrentPermissions() {
	// Check if we can see the FIP device
	devices := hid.Enumerate(0, 0)
	fipFound := false

	for _, dev := range devices {
		if dev.VendorID == 0x06A3 && dev.ProductID == 0xA2AE {
			fipFound = true
			fmt.Printf("✓ FIP device detected: %s\n", dev.Product)
			fmt.Printf("  Vendor ID: 0x%04x, Product ID: 0x%04x\n", dev.VendorID, dev.ProductID)
			fmt.Printf("  Path: %s\n", dev.Path)
			break
		}
	}

	if !fipFound {
		fmt.Println("✗ FIP device not found in HID enumeration")
		fmt.Println("  This may indicate a permission issue")
	}

	// Check system permissions
	fmt.Println("\n--- System Permission Checks ---")

	// Check if we're running as root
	if os.Geteuid() == 0 {
		fmt.Println("✓ Running as root (elevated privileges)")
	} else {
		fmt.Println("⚠ Running as regular user")
	}

	// Check if we can access /dev/hidraw devices (Linux-style)
	if _, err := os.Stat("/dev/hidraw0"); err == nil {
		fmt.Println("✓ HID raw devices accessible")
	} else {
		fmt.Println("⚠ HID raw devices not accessible (normal on macOS)")
	}

	// Check if we can access IOKit devices
	if _, err := os.Stat("/dev/usb"); err == nil {
		fmt.Println("✓ USB devices accessible")
	} else {
		fmt.Println("⚠ USB devices not directly accessible (normal on macOS)")
	}
}

func testFIPAccess() {
	fmt.Println("\n--- Testing FIP Device Access ---")

	// Try to open the FIP device directly
	devices := hid.Enumerate(0x06A3, 0xA2AE)
	if len(devices) == 0 {
		fmt.Println("✗ No FIP devices found in enumeration")
		return
	}

	device := devices[0]
	fmt.Printf("Attempting to open FIP device: %s\n", device.Product)

	// Try to open the device
	hidDevice, err := device.Open()
	if err != nil {
		fmt.Printf("✗ Failed to open FIP device: %v\n", err)
		fmt.Println("  This indicates a permission issue")
		return
	}
	defer hidDevice.Close()

	fmt.Println("✓ Successfully opened FIP device!")

	// Try to read from the device
	fmt.Println("Attempting to read button states...")

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

	// Monitor for button presses for 10 seconds
	fmt.Println("Monitoring button presses for 10 seconds...")
	fmt.Println("Press buttons on the FIP device to test...")

	timeout := time.After(10 * time.Second)
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

func printPermissionInstructions() {
	fmt.Println("To fix FIP device permissions on macOS, try the following:")
	fmt.Println()

	fmt.Println("1. SYSTEM PREFERENCES METHOD:")
	fmt.Println("   - Open System Preferences > Security & Privacy > Privacy")
	fmt.Println("   - Select 'Input Monitoring' from the left sidebar")
	fmt.Println("   - Click the lock icon and enter your password")
	fmt.Println("   - Add your terminal application (Terminal.app or iTerm)")
	fmt.Println("   - Add the Go binary if it appears in the list")
	fmt.Println()

	fmt.Println("2. TERMINAL PERMISSIONS:")
	fmt.Println("   - Go to System Preferences > Security & Privacy > Privacy")
	fmt.Println("   - Select 'Accessibility' from the left sidebar")
	fmt.Println("   - Add Terminal.app or your terminal application")
	fmt.Println()

	fmt.Println("3. RUN WITH ELEVATED PRIVILEGES:")
	fmt.Println("   - Try running the test with sudo:")
	fmt.Println("     sudo go run cmd/test_fip_permissions/main.go")
	fmt.Println()

	fmt.Println("4. CHECK USB DEVICE PERMISSIONS:")
	fmt.Println("   - The FIP device might need USB device permissions")
	fmt.Println("   - Check if the device appears in System Information")
	fmt.Println("   - Try disconnecting and reconnecting the device")
	fmt.Println()

	fmt.Println("5. DEVELOPER TOOLS:")
	fmt.Println("   - If you're a developer, you might need to sign the application")
	fmt.Println("   - Or run with proper entitlements")
	fmt.Println()

	fmt.Println("6. ALTERNATIVE APPROACH:")
	fmt.Println("   - The FIP might require a different access method")
	fmt.Println("   - Try using IOKit directly instead of HID library")
	fmt.Println()
}

func testAlternativeMethods() {
	fmt.Println("Testing alternative access methods...")

	// Method 1: Try with different HID library options
	fmt.Println("\n--- Method 1: HID Library with Different Options ---")
	testHIDWithOptions()

	// Method 2: Try using our USB core directly
	fmt.Println("\n--- Method 2: Direct USB Core Access ---")
	testUSBCoreAccess()

	// Method 3: Check system commands
	fmt.Println("\n--- Method 3: System USB Information ---")
	testSystemUSBInfo()
}

func testHIDWithOptions() {
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

func testSystemUSBInfo() {
	fmt.Println("Checking system USB information...")

	// Try to get USB device information using system commands
	commands := []struct {
		name string
		cmd  []string
	}{
		{"system_profiler", []string{"system_profiler", "SPUSBDataType"}},
		{"ioreg", []string{"ioreg", "-p", "IOUSBHostDevice", "-l"}},
		{"lsusb", []string{"lsusb"}},
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
		if strings.Contains(outputStr, "Saitek") || strings.Contains(outputStr, "FIP") ||
			strings.Contains(outputStr, "06A3") || strings.Contains(outputStr, "A2AE") {
			fmt.Printf("✓ Found Saitek FIP device in %s output\n", command.name)
			// Print relevant lines
			lines := strings.Split(outputStr, "\n")
			for _, line := range lines {
				if strings.Contains(line, "Saitek") || strings.Contains(line, "FIP") ||
					strings.Contains(line, "06A3") || strings.Contains(line, "A2AE") {
					fmt.Printf("  %s\n", strings.TrimSpace(line))
				}
			}
		} else {
			fmt.Printf("⚠ No Saitek FIP device found in %s output\n", command.name)
		}
	}
}
