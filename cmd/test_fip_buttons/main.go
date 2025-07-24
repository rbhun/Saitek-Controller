package main

import (
	"fmt"
	"time"

	"saitek-controller/internal/usb"

	"github.com/karalabe/hid"
)

func main() {
	fmt.Println("Testing Saitek FIP Button Reading")
	fmt.Println("==================================")

	// First, let's see what HID devices are available
	fmt.Println("\n=== Available HID Devices ===")
	devices := hid.Enumerate(0, 0)
	fmt.Printf("Found %d HID devices:\n", len(devices))

	var fipDevice *hid.DeviceInfo
	for i, dev := range devices {
		fmt.Printf("  Device %d: Vendor=0x%04x Product=0x%04x Name='%s' Manufacturer='%s' Path='%s'\n",
			i, dev.VendorID, dev.ProductID, dev.Product, dev.Manufacturer, dev.Path)

		// Check if this is the FIP
		if dev.VendorID == 0x06A3 && dev.ProductID == 0xA2AE {
			fipDevice = &dev
			fmt.Printf("    *** This is the FIP device! ***\n")
		}
	}

	if fipDevice == nil {
		fmt.Println("FIP device not found in HID enumeration!")
		return
	}

	fmt.Printf("\n=== Testing FIP Button Reading ===\n")
	fmt.Printf("FIP Device: Vendor=0x%04x Product=0x%04x Name='%s' Manufacturer='%s'\n",
		fipDevice.VendorID, fipDevice.ProductID, fipDevice.Product, fipDevice.Manufacturer)

	// Try to open the FIP device directly using the HID library
	fmt.Println("\n--- Attempting direct HID access for button reading ---")
	handle, err := fipDevice.Open()
	if err != nil {
		fmt.Printf("Failed to open FIP device directly: %v\n", err)
		fmt.Println("This might be a permission issue. Let's try alternative approaches...")
	} else {
		fmt.Printf("Successfully opened FIP device directly!\n")
		defer handle.Close()

		// Try to read button data from the FIP
		fmt.Println("Attempting to read button data from FIP...")
		fmt.Println("Press buttons on the FIP device to see if we can detect them...")
		fmt.Println("(Press Ctrl+C to stop)")

		// Read loop for button presses
		for {
			readData := make([]byte, 64)
			read, err := handle.Read(readData)
			if err != nil {
				fmt.Printf("Failed to read from FIP: %v\n", err)
				break
			} else if read > 0 {
				fmt.Printf("Received %d bytes from FIP: %v\n", read, readData[:read])

				// Try to interpret the data as button presses
				interpretButtonData(readData[:read])
			}

			time.Sleep(10 * time.Millisecond) // Small delay to avoid busy waiting
		}
	}

	// Try using our USB abstraction layer
	fmt.Println("\n--- Testing USB abstraction layer for button reading ---")
	device, err := usb.OpenDevice(0x06A3, 0xA2AE)
	if err != nil {
		fmt.Printf("Failed to open device via USB abstraction: %v\n", err)
	} else {
		fmt.Printf("Successfully opened device via USB abstraction: %s\n", device.Name)
		fmt.Printf("Device connected: %v\n", device.IsConnected())

		// Try to read button data
		fmt.Println("Attempting to read button data via USB abstraction...")
		fmt.Println("Press buttons on the FIP device to see if we can detect them...")
		fmt.Println("(Press Ctrl+C to stop)")

		for i := 0; i < 10; i++ { // Try reading 10 times
			data, err := device.ReadBulkData(0x81, 64)
			if err != nil {
				fmt.Printf("Failed to read bulk data: %v\n", err)
			} else if len(data) > 0 {
				fmt.Printf("Received %d bytes via USB abstraction: %v\n", len(data), data)
				interpretButtonData(data)
			}

			time.Sleep(100 * time.Millisecond)
		}

		device.Close()
		fmt.Printf("Device closed\n")
	}

	fmt.Println("\n=== Test completed ===")
}

// interpretButtonData tries to interpret the received data as button presses
func interpretButtonData(data []byte) {
	fmt.Printf("Button data interpretation:\n")

	// Print each byte with its position
	for i, b := range data {
		if b != 0 {
			fmt.Printf("  Byte %d: 0x%02x (%d) - ", i, b, b)

			// Try to interpret as button states
			if i == 0 {
				fmt.Printf("Report ID or button state\n")
			} else if i == 1 {
				fmt.Printf("Button state byte\n")
				// Check individual bits for buttons
				for bit := 0; bit < 8; bit++ {
					if b&(1<<bit) != 0 {
						fmt.Printf("    Button %d pressed\n", bit+1)
					}
				}
			} else {
				fmt.Printf("Additional data\n")
			}
		}
	}

	// Try to interpret as a complete button report
	if len(data) >= 2 {
		fmt.Printf("Complete button report: %v\n", data[:2])
	}
}
