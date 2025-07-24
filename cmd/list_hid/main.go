package main

import (
	"fmt"

	"github.com/karalabe/hid"
)

func main() {
	// List all HID devices
	devices := hid.Enumerate(0, 0)

	fmt.Printf("Found %d HID devices:\n", len(devices))
	for i, dev := range devices {
		fmt.Printf("  Device %d: Vendor=0x%04x Product=0x%04x Manufacturer='%s' Product='%s' Path='%s'\n",
			i, dev.VendorID, dev.ProductID, dev.Manufacturer, dev.Product, dev.Path)

		// Check if this matches any known Saitek/Logitech devices
		if dev.VendorID == 0x06A3 {
			fmt.Printf("    *** This is a Logitech/Saitek device! ***\n")
		}
	}

	// Specifically look for Logitech/Saitek devices
	fmt.Printf("\nLogitech/Saitek devices (Vendor ID 0x06A3):\n")
	for i, dev := range devices {
		if dev.VendorID == 0x06A3 {
			fmt.Printf("  Device %d: Product=0x%04x Manufacturer='%s' Product='%s' Path='%s'\n",
				i, dev.ProductID, dev.Manufacturer, dev.Product, dev.Path)
		}
	}
}
