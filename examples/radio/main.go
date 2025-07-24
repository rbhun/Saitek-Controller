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
		productID   = flag.Uint("product", 0x0D05, "USB Product ID")
		com1Active  = flag.String("com1a", "118.00", "COM1 Active frequency")
		com1Standby = flag.String("com1s", "118.50", "COM1 Standby frequency")
		com2Active  = flag.String("com2a", "121.30", "COM2 Active frequency")
		com2Standby = flag.String("com2s", "121.90", "COM2 Standby frequency")
		interactive = flag.Bool("interactive", false, "Run in interactive mode")
	)
	flag.Parse()

	// Create radio panel
	radio := fip.NewRadioPanelWithUSB(uint16(*vendorID), uint16(*productID))

	// Connect to the device
	fmt.Printf("Connecting to Saitek Flight Radio Panel...\n")
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
	}

	if *interactive {
		// Run interactive mode
		runInteractive(radio)
	} else {
		// Run monitoring mode
		runMonitoring(radio)
	}
}

func runInteractive(radio *fip.RadioPanel) {
	fmt.Printf("\nInteractive mode - Press Ctrl+C to exit\n")
	fmt.Printf("Commands:\n")
	fmt.Printf("  set <com1a> <com1s> <com2a> <com2s> - Set frequencies\n")
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
				if radio.IsConnected() {
					data, err := radio.ReadSwitchState()
					if err != nil {
						continue
					}

					state := radio.ParseSwitchState(data)
					if state != nil {
						for name, active := range state {
							if active {
								fmt.Printf("Radio Panel: %s activated\n", name)
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

func runMonitoring(radio *fip.RadioPanel) {
	fmt.Printf("\nMonitoring mode - Press Ctrl+C to exit\n")
	fmt.Printf("Monitoring radio panel state...\n")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring loop
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if radio.IsConnected() {
				data, err := radio.ReadSwitchState()
				if err != nil {
					continue
				}

				state := radio.ParseSwitchState(data)
				if state != nil {
					for name, active := range state {
						if active {
							fmt.Printf("Radio Panel: %s activated\n", name)
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
