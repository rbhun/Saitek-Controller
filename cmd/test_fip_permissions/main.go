package main

import (
	"fmt"
	"os"
	"time"

	"github.com/karalabe/hid"
)

func main() {
	fmt.Println("Saitek FIP Button Reading Test")
	fmt.Println("===============================")
	fmt.Println()
	fmt.Println("This test will help you set up proper permissions for the FIP device.")
	fmt.Println()

	// Check if we're running with proper permissions
	fmt.Println("=== Permission Check ===")
	if os.Geteuid() == 0 {
		fmt.Println("✓ Running as root (sudo)")
	} else {
		fmt.Println("⚠ Running as regular user")
		fmt.Println("   Note: The FIP device may require elevated permissions")
	}

	// Check for FIP device
	fmt.Println("\n=== FIP Device Detection ===")
	devices := hid.Enumerate(0, 0)
	var fipDevice *hid.DeviceInfo

	for _, dev := range devices {
		if dev.VendorID == 0x06A3 && dev.ProductID == 0xA2AE {
			fipDevice = &dev
			fmt.Printf("✓ Found FIP device: %s\n", dev.Product)
			break
		}
	}

	if fipDevice == nil {
		fmt.Println("✗ FIP device not found!")
		fmt.Println("\nTroubleshooting:")
		fmt.Println("1. Make sure the FIP is connected via USB")
		fmt.Println("2. Try running with sudo: sudo go run cmd/test_fip_permissions/main.go")
		fmt.Println("3. Check if the device appears in System Information > USB")
		return
	}

	// Try to open the device
	fmt.Println("\n=== Attempting to Open FIP Device ===")
	handle, err := fipDevice.Open()
	if err != nil {
		fmt.Printf("✗ Failed to open FIP device: %v\n", err)
		fmt.Println("\n=== Permission Setup Required ===")
		fmt.Println("The FIP device requires special permissions on macOS.")
		fmt.Println("Follow these steps to enable access:")
		fmt.Println()
		fmt.Println("1. Go to System Preferences > Security & Privacy > Privacy")
		fmt.Println("2. Select 'Input Monitoring' from the left sidebar")
		fmt.Println("3. Click the lock icon to make changes (enter your password)")
		fmt.Println("4. Click the '+' button and add your terminal application:")
		fmt.Println("   - Terminal.app (if using Terminal)")
		fmt.Println("   - iTerm.app (if using iTerm)")
		fmt.Println("   - Or add the Go executable directly")
		fmt.Println("5. Make sure the checkbox is enabled for the application")
		fmt.Println("6. Restart your terminal application")
		fmt.Println("7. Run this test again")
		fmt.Println()
		fmt.Println("Alternative: Run with sudo (temporary solution):")
		fmt.Println("   sudo go run cmd/test_fip_permissions/main.go")
		return
	}

	fmt.Println("✓ Successfully opened FIP device!")
	defer handle.Close()

	// Test button reading
	fmt.Println("\n=== Testing Button Reading ===")
	fmt.Println("Press buttons on the FIP device to see if we can detect them...")
	fmt.Println("(Press Ctrl+C to stop)")
	fmt.Println()

	// Read loop for button presses
	for {
		readData := make([]byte, 64)
		read, err := handle.Read(readData)
		if err != nil {
			fmt.Printf("✗ Failed to read from FIP: %v\n", err)
			break
		} else if read > 0 {
			fmt.Printf("✓ Received %d bytes from FIP: %v\n", read, readData[:read])

			// Interpret button data
			interpretFIPButtonData(readData[:read])
		}

		time.Sleep(10 * time.Millisecond) // Small delay to avoid busy waiting
	}
}

// interpretFIPButtonData interprets the received data as FIP button presses
func interpretFIPButtonData(data []byte) {
	fmt.Printf("Button interpretation:\n")

	// Based on the HID descriptor, the FIP has 12 buttons
	// Each button is represented by a bit in the first 2 bytes
	if len(data) >= 2 {
		buttonStates := data[:2]
		fmt.Printf("  Button states: %08b %08b\n", buttonStates[0], buttonStates[1])

		// Check each button
		for i := 0; i < 12; i++ {
			byteIndex := i / 8
			bitIndex := i % 8

			if byteIndex < len(buttonStates) {
				pressed := (buttonStates[byteIndex] & (1 << bitIndex)) != 0
				if pressed {
					fmt.Printf("  ✓ Button %d pressed\n", i+1)
				}
			}
		}
	}

	// Show raw data for debugging
	fmt.Printf("  Raw data: %v\n", data)
}
