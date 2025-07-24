package main

import (
	"fmt"
	"time"

	"saitek-controller/internal/usb"

	"github.com/karalabe/hid"
)

func main() {
	fmt.Println("Testing Saitek FIP Communication")
	fmt.Println("================================")

	// First, let's see what HID devices are available
	fmt.Println("\n=== Available HID Devices ===")
	devices := hid.Enumerate(0, 0)
	fmt.Printf("Found %d HID devices:\n", len(devices))

	var fipDevice *hid.DeviceInfo
	for i, dev := range devices {
		fmt.Printf("  Device %d: Vendor=0x%04x Product=0x%04x Name='%s' Manufacturer='%s' Path='%s'\n",
			i, dev.VendorID, dev.ProductID, dev.Product, dev.Manufacturer, dev.Path)

		// Check if this looks like a Saitek device
		if dev.VendorID == 0x06A3 || dev.Manufacturer == "Saitek" || dev.Manufacturer == "Logitech" {
			fmt.Printf("    *** This looks like a Saitek device! ***\n")

			// Check if this is the FIP
			if dev.VendorID == 0x06A3 && dev.ProductID == 0xA2AE {
				fipDevice = &dev
				fmt.Printf("    *** This is the FIP device! ***\n")
			}
		}
	}

	if fipDevice == nil {
		fmt.Println("FIP device not found in HID enumeration!")
		return
	}

	fmt.Printf("\n=== Testing FIP Device Direct Access ===\n")
	fmt.Printf("FIP Device: Vendor=0x%04x Product=0x%04x Name='%s' Manufacturer='%s'\n",
		fipDevice.VendorID, fipDevice.ProductID, fipDevice.Product, fipDevice.Manufacturer)

	// Try to open the FIP device directly using the HID library
	fmt.Println("\n--- Attempting direct HID access ---")
	handle, err := fipDevice.Open()
	if err != nil {
		fmt.Printf("Failed to open FIP device directly: %v\n", err)
		fmt.Println("This might be a permission issue. Let's try alternative approaches...")
	} else {
		fmt.Printf("Successfully opened FIP device directly!\n")
		defer handle.Close()

		// Try to send some data to the FIP
		fmt.Println("Attempting to send data to FIP...")
		testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
		written, err := handle.Write(testData)
		if err != nil {
			fmt.Printf("Failed to write to FIP: %v\n", err)
		} else {
			fmt.Printf("Successfully wrote %d bytes to FIP\n", written)
		}

		// Try to read data from the FIP
		fmt.Println("Attempting to read data from FIP...")
		readData := make([]byte, 64)
		read, err := handle.Read(readData)
		if err != nil {
			fmt.Printf("Failed to read from FIP: %v\n", err)
		} else {
			fmt.Printf("Successfully read %d bytes from FIP: %v\n", read, readData[:read])
		}
	}

	// Try using our USB abstraction layer
	fmt.Println("\n--- Testing USB abstraction layer ---")
	device, err := usb.OpenDevice(0x06A3, 0xA2AE)
	if err != nil {
		fmt.Printf("Failed to open device via USB abstraction: %v\n", err)
	} else {
		fmt.Printf("Successfully opened device via USB abstraction: %s\n", device.Name)
		fmt.Printf("Device connected: %v\n", device.IsConnected())

		// Try to send a test message
		testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
		if err := device.SendControlMessage(0x21, 0x09, 0x0200, 0x0001, testData); err != nil {
			fmt.Printf("Failed to send control message: %v\n", err)
		} else {
			fmt.Printf("Successfully sent control message\n")
		}

		// Try to read some data
		data, err := device.ReadBulkData(0x81, 64)
		if err != nil {
			fmt.Printf("Failed to read bulk data: %v\n", err)
		} else {
			fmt.Printf("Successfully read %d bytes: %v\n", len(data), data)
		}

		device.Close()
		fmt.Printf("Device closed\n")
	}

	// Try different permission approaches
	fmt.Println("\n--- Testing permission alternatives ---")

	// Check if we can access the device with different methods
	fmt.Println("1. Checking if device is accessible via IOKit...")
	// This would require implementing IOKit-specific code

	fmt.Println("2. Checking if device requires special permissions...")
	fmt.Println("   On macOS, some HID devices require special entitlements or permissions.")
	fmt.Println("   You might need to:")
	fmt.Println("   - Add the device to System Preferences > Security & Privacy > Privacy > Input Monitoring")
	fmt.Println("   - Run the application with sudo")
	fmt.Println("   - Add special entitlements to the application")

	fmt.Println("\n=== Test completed ===")
	time.Sleep(2 * time.Second)
}
