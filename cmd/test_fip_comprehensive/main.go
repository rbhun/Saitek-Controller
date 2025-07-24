package main

import (
	"fmt"
	"os"
	"os/exec"

	"saitek-controller/internal/usb"

	"github.com/karalabe/hid"
)

func main() {
	fmt.Println("Comprehensive Saitek FIP Testing")
	fmt.Println("=================================")

	// Check system information
	fmt.Println("\n=== System Information ===")
	fmt.Printf("OS: %s\n", os.Getenv("OSTYPE"))

	// Check if we're running as root
	if os.Geteuid() == 0 {
		fmt.Println("Running as root (sudo)")
	} else {
		fmt.Println("Running as regular user")
	}

	// Check USB devices using system commands
	fmt.Println("\n=== System USB Device Information ===")

	// Try to get USB device information using system commands
	cmd := exec.Command("system_profiler", "SPUSBDataType")
	if output, err := cmd.Output(); err == nil {
		fmt.Println("USB devices found by system:")
		fmt.Println(string(output))
	} else {
		fmt.Printf("Failed to get USB info: %v\n", err)
	}

	// Check HID devices using system commands
	fmt.Println("\n=== System HID Device Information ===")
	cmd = exec.Command("ioreg", "-l", "-w", "0", "-r", "-c", "IOHIDDevice")
	if output, err := cmd.Output(); err == nil {
		fmt.Println("HID devices found by system:")
		fmt.Println(string(output))
	} else {
		fmt.Printf("Failed to get HID info: %v\n", err)
	}

	// Check what HID devices our library can see
	fmt.Println("\n=== HID Library Device Information ===")
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

	fmt.Printf("\n=== Detailed FIP Device Analysis ===\n")
	fmt.Printf("FIP Device: Vendor=0x%04x Product=0x%04x Name='%s' Manufacturer='%s'\n",
		fipDevice.VendorID, fipDevice.ProductID, fipDevice.Product, fipDevice.Manufacturer)

	// Try multiple approaches to open the device
	fmt.Println("\n--- Testing Multiple Access Methods ---")

	// Method 1: Direct HID library access
	fmt.Println("\n1. Direct HID library access:")
	handle, err := fipDevice.Open()
	if err != nil {
		fmt.Printf("   Failed: %v\n", err)
	} else {
		fmt.Printf("   Success! Device opened directly.\n")
		defer handle.Close()

		// Try to read some data
		readData := make([]byte, 64)
		read, err := handle.Read(readData)
		if err != nil {
			fmt.Printf("   Failed to read: %v\n", err)
		} else {
			fmt.Printf("   Read %d bytes: %v\n", read, readData[:read])
		}
	}

	// Method 2: USB abstraction layer
	fmt.Println("\n2. USB abstraction layer:")
	device, err := usb.OpenDevice(0x06A3, 0xA2AE)
	if err != nil {
		fmt.Printf("   Failed: %v\n", err)
	} else {
		fmt.Printf("   Success! Device opened via USB abstraction.\n")
		fmt.Printf("   Device name: %s\n", device.Name)
		fmt.Printf("   Device connected: %v\n", device.IsConnected())

		// Try to read some data
		data, err := device.ReadBulkData(0x81, 64)
		if err != nil {
			fmt.Printf("   Failed to read: %v\n", err)
		} else {
			fmt.Printf("   Read %d bytes: %v\n", len(data), data)
		}

		device.Close()
	}

	// Method 3: Try different endpoints
	fmt.Println("\n3. Testing different endpoints:")
	if device, err := usb.OpenDevice(0x06A3, 0xA2AE); err == nil {
		defer device.Close()

		endpoints := []uint8{0x81, 0x82, 0x83, 0x84, 0x85}
		for _, endpoint := range endpoints {
			data, err := device.ReadBulkData(endpoint, 64)
			if err != nil {
				fmt.Printf("   Endpoint 0x%02x: Failed - %v\n", endpoint, err)
			} else {
				fmt.Printf("   Endpoint 0x%02x: Read %d bytes\n", endpoint, len(data))
			}
		}
	}

	// Method 4: Try sending control messages
	fmt.Println("\n4. Testing control messages:")
	if device, err := usb.OpenDevice(0x06A3, 0xA2AE); err == nil {
		defer device.Close()

		// Try different control message types
		controlMessages := []struct {
			requestType uint16
			request     uint16
			value       uint16
			index       uint16
			description string
		}{
			{0x21, 0x09, 0x0200, 0x0001, "Standard HID output"},
			{0x21, 0x09, 0x0000, 0x0000, "HID output with zero values"},
			{0x21, 0x01, 0x0000, 0x0000, "HID get report"},
			{0xA1, 0x01, 0x0000, 0x0000, "HID get report (input)"},
		}

		for _, msg := range controlMessages {
			testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
			err := device.SendControlMessage(msg.requestType, msg.request, msg.value, msg.index, testData)
			if err != nil {
				fmt.Printf("   %s: Failed - %v\n", msg.description, err)
			} else {
				fmt.Printf("   %s: Success\n", msg.description)
			}
		}
	}

	// Check for permission issues
	fmt.Println("\n=== Permission Analysis ===")
	fmt.Println("The FIP device cannot be opened, which could be due to:")
	fmt.Println("1. macOS security restrictions")
	fmt.Println("2. Missing entitlements")
	fmt.Println("3. Device requires special drivers")
	fmt.Println("4. Device is not a standard HID device")

	fmt.Println("\nPossible solutions:")
	fmt.Println("1. Add the application to System Preferences > Security & Privacy > Privacy > Input Monitoring")
	fmt.Println("2. Run the application with sudo (already tried)")
	fmt.Println("3. Add special entitlements to the application")
	fmt.Println("4. Use a different USB access method (not HID)")
	fmt.Println("5. Check if the device needs specific drivers or firmware")

	fmt.Println("\n=== Test completed ===")
}
