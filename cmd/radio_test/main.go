package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"saitek-controller/internal/fip"
)

func main() {
	// Parse command line flags
	var (
		vendorID    = flag.Uint("vendor", 0x06A3, "USB Vendor ID")
		productID   = flag.Uint("product", 0x0D05, "USB Product ID")
		com1Active  = flag.String("com1a", "118.00", "COM1 Active frequency")
		com1Standby = flag.String("com1s", "118.50", "COM1 Standby frequency")
		com2Active  = flag.String("com2a", "121.30", "COM2 Active frequency")
		com2Standby = flag.String("com2s", "121.90", "COM2 Standby frequency")
	)
	flag.Parse()

	// Create radio panel
	radio := fip.NewRadioPanelWithUSB(uint16(*vendorID), uint16(*productID))

	// Connect to the device
	fmt.Printf("Connecting to Saitek Flight Radio Panel...\n")
	fmt.Printf("Vendor ID: 0x%04x, Product ID: 0x%04x\n", *vendorID, *productID)

	if err := radio.Connect(); err != nil {
		log.Printf("Failed to connect to radio panel: %v", err)
		log.Printf("Running in mock mode for testing")
	} else {
		fmt.Printf("Successfully connected to radio panel\n")
	}

	// Set initial display
	fmt.Printf("Setting display:\n")
	fmt.Printf("  COM1 Active:   %s\n", *com1Active)
	fmt.Printf("  COM1 Standby:  %s\n", *com1Standby)
	fmt.Printf("  COM2 Active:   %s\n", *com2Active)
	fmt.Printf("  COM2 Standby:  %s\n", *com2Standby)

	if err := radio.SetDisplay(*com1Active, *com1Standby, *com2Active, *com2Standby); err != nil {
		log.Printf("Failed to set display: %v", err)
	} else {
		fmt.Printf("Successfully sent display data to radio panel\n")
	}

	// Wait a moment to see if the display updates
	fmt.Printf("Waiting 3 seconds to observe display...\n")
	time.Sleep(3 * time.Second)

	// Try a different frequency to test updates
	fmt.Printf("Updating to test frequencies...\n")
	if err := radio.SetDisplay("122.80", "122.90", "123.40", "123.50"); err != nil {
		log.Printf("Failed to update display: %v", err)
	} else {
		fmt.Printf("Successfully updated display\n")
	}

	// Wait another moment
	fmt.Printf("Waiting 3 more seconds...\n")
	time.Sleep(3 * time.Second)

	fmt.Printf("Test completed!\n")
}
