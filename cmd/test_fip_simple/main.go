package main

import (
	"fmt"
	"time"

	"github.com/karalabe/hid"
)

func main() {
	fmt.Println("Simple FIP Access Test")
	fmt.Println("======================")
	fmt.Println()

	// Check if we can enumerate the device
	fmt.Println("1. Enumerating FIP devices...")
	devices := hid.Enumerate(0x06A3, 0xA2AE)

	if len(devices) == 0 {
		fmt.Println("✗ No FIP devices found")
		return
	}

	fmt.Printf("✓ Found %d FIP device(s)\n", len(devices))
	for i, dev := range devices {
		fmt.Printf("  Device %d: %s (0x%04x:0x%04x)\n", i, dev.Product, dev.VendorID, dev.ProductID)
	}

	// Try to open the first device
	fmt.Println("\n2. Attempting to open FIP device...")
	device := devices[0]

	hidDevice, err := device.Open()
	if err != nil {
		fmt.Printf("✗ Failed to open device: %v\n", err)
		fmt.Println("\nThis indicates a permission issue.")
		fmt.Println("Please ensure you have granted:")
		fmt.Println("- Input Monitoring permission to Terminal")
		fmt.Println("- Accessibility permission to Terminal")
		fmt.Println("\nTry restarting Terminal after granting permissions.")
		return
	}
	defer hidDevice.Close()

	fmt.Println("✓ Successfully opened FIP device!")
	fmt.Println("\n3. Testing button reading...")

	// Try to read button states
	fmt.Println("Press buttons on the FIP device to test...")
	fmt.Println("(Will monitor for 10 seconds)")

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
			fmt.Println("\nSuccess! The FIP device is accessible and working.")
			return
		}
	}
}
