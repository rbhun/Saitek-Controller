package main

import (
	"fmt"
	"log"

	"github.com/google/gousb"
)

func main() {
	ctx := gousb.NewContext()
	defer ctx.Close()

	// List all devices
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return true // Open all devices
	})
	if err != nil {
		log.Fatalf("Failed to list devices: %v", err)
	}

	fmt.Printf("Found %d USB devices:\n", len(devs))
	for i, dev := range devs {
		defer dev.Close()

		// Try to get manufacturer and product strings
		manufacturer := "Unknown"
		product := "Unknown"

		if man, err := dev.Manufacturer(); err == nil {
			manufacturer = man
		}
		if prod, err := dev.Product(); err == nil {
			product = prod
		}

		fmt.Printf("  Device %d: Vendor=0x%04x Product=0x%04x Manufacturer='%s' Product='%s'\n",
			i, dev.Desc.Vendor, dev.Desc.Product, manufacturer, product)

		// Check if this matches any known Saitek/Logitech devices
		if dev.Desc.Vendor == 0x06A3 {
			fmt.Printf("    *** This is a Logitech/Saitek device! ***\n")
		}
	}
}
