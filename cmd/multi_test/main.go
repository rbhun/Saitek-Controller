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
		vendorID   = flag.Uint("vendor", 0x06A3, "USB Vendor ID")
		productID  = flag.Uint("product", 0x0D06, "USB Product ID")
		topRow     = flag.String("top", "250", "Top row display")
		bottomRow  = flag.String("bottom", "3000", "Bottom row display")
		buttonLEDs = flag.Uint("leds", 0x01, "Button LED states")
	)
	flag.Parse()

	// Create multi panel
	multi := fip.NewMultiPanelWithUSB(uint16(*vendorID), uint16(*productID))

	// Connect to the device
	fmt.Printf("Connecting to Saitek Multi Panel...\n")
	fmt.Printf("Vendor ID: 0x%04x, Product ID: 0x%04x\n", *vendorID, *productID)

	if err := multi.Connect(); err != nil {
		log.Printf("Failed to connect to multi panel: %v", err)
		log.Printf("Running in mock mode for testing")
	} else {
		fmt.Printf("Successfully connected to multi panel\n")
	}

	// Set initial display
	fmt.Printf("Setting display:\n")
	fmt.Printf("  Top Row:      %s\n", *topRow)
	fmt.Printf("  Bottom Row:   %s\n", *bottomRow)
	fmt.Printf("  Button LEDs:  %08b\n", *buttonLEDs)

	if err := multi.SetDisplay(*topRow, *bottomRow, uint8(*buttonLEDs)); err != nil {
		log.Printf("Failed to set display: %v", err)
	} else {
		fmt.Printf("Successfully sent display data to multi panel\n")
	}

	// Wait a moment to see if the display updates
	fmt.Printf("Waiting 3 seconds to observe display...\n")
	time.Sleep(3 * time.Second)

	// Try a different display to test updates
	fmt.Printf("Updating to test values...\n")
	if err := multi.SetDisplay("120", "5000", 0x0F); err != nil {
		log.Printf("Failed to update display: %v", err)
	} else {
		fmt.Printf("Successfully updated display\n")
	}

	// Wait another moment
	time.Sleep(3 * time.Second)

	// Test button LED control
	fmt.Printf("Testing button LED control...\n")
	if err := multi.SetButtonLEDs(0xAA); err != nil {
		log.Printf("Failed to set button LEDs: %v", err)
	} else {
		fmt.Printf("Successfully set button LEDs\n")
	}

	// Wait and test switch reading
	fmt.Printf("Testing switch reading for 5 seconds...\n")
	start := time.Now()
	for time.Since(start) < 5*time.Second {
		if multi.IsConnected() {
			data, err := multi.ReadSwitchState()
			if err != nil {
				log.Printf("Error reading switch state: %v", err)
				continue
			}

			state := multi.ParseSwitchState(data)
			if state != nil {
				// Log any active switches
				for name, active := range state {
					if active {
						fmt.Printf("Multi Panel: %s activated\n", name)
					}
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("Multi panel test completed!\n")
	multi.Close()
}
