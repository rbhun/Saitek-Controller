package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"saitek-controller/internal/fip"
)

func main() {
	// Parse command line flags
	var (
		vendorID    = flag.Uint("vendor", 0x06A3, "USB Vendor ID")
		productID   = flag.Uint("product", 0x0D06, "USB Product ID")
		topRow      = flag.String("top", "250", "Top row display")
		bottomRow   = flag.String("bottom", "3000", "Bottom row display")
		buttonLEDs  = flag.Uint("leds", 0x01, "Button LED states")
		interactive = flag.Bool("interactive", false, "Run in interactive mode")
	)
	flag.Parse()

	// Create multi panel
	multi := fip.NewMultiPanelWithUSB(uint16(*vendorID), uint16(*productID))

	// Connect to the device
	fmt.Printf("Connecting to Saitek Multi Panel...\n")
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
	}

	if *interactive {
		// Run interactive mode
		runInteractive(multi)
	} else {
		// Run monitoring mode
		runMonitoring(multi)
	}
}

func runInteractive(multi *fip.MultiPanel) {
	fmt.Printf("\nInteractive mode - Press Ctrl+C to exit\n")
	fmt.Printf("Commands:\n")
	fmt.Printf("  set <top> <bottom> <leds> - Set display and LEDs\n")
	fmt.Printf("  leds <value> - Set button LEDs only\n")
	fmt.Printf("  status - Show current display\n")
	fmt.Printf("  help - Show this help\n")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring in background
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if multi.IsConnected() {
					data, err := multi.ReadSwitchState()
					if err != nil {
						continue
					}

					state := multi.ParseSwitchState(data)
					if state != nil {
						for name, active := range state {
							if active {
								fmt.Printf("Multi Panel: %s activated\n", name)
							}
						}
					}
				}
			}
		}
	}()

	// Wait for exit signal
	<-sigChan
	fmt.Printf("\nShutting down...\n")
}

func runMonitoring(multi *fip.MultiPanel) {
	fmt.Printf("\nMonitoring mode - Press Ctrl+C to exit\n")
	fmt.Printf("Monitoring multi panel state...\n")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring loop
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if multi.IsConnected() {
				data, err := multi.ReadSwitchState()
				if err != nil {
					continue
				}

				state := multi.ParseSwitchState(data)
				if state != nil {
					for name, active := range state {
						if active {
							fmt.Printf("Multi Panel: %s activated\n", name)
						}
					}
				}
			}
		case <-sigChan:
			fmt.Printf("\nShutting down...\n")
			return
		}
	}
}
