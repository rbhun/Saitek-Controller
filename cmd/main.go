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
	"saitek-controller/internal/usb"

	"github.com/faiface/pixel/pixelgl"
	"github.com/karalabe/hid"
)

func main() {
	// Parse command line flags
	var (
		width       = flag.Int("width", 320, "FIP display width")
		height      = flag.Int("height", 240, "FIP display height")
		title       = flag.String("title", "Saitek FIP Controller", "Window title")
		imageFile   = flag.String("image", "", "Image file to display")
		instrument  = flag.String("instrument", "test", "Instrument type (artificial_horizon, airspeed, altimeter, compass, vsi, turn_coordinator, test)")
		vendorID    = flag.String("vendor", "06a3", "USB vendor ID (hex, e.g. 06a3)")
		productID   = flag.String("product", "0a2ae", "USB product ID (hex, e.g. 0a2ae)")
		listDevices = flag.Bool("list-devices", false, "List all connected HID devices and exit")
		panelType   = flag.String("panel", "fip", "Panel type (fip, radio, multi)")
		com1Active  = flag.String("com1a", "118.00", "COM1 Active frequency (radio panel)")
		com1Standby = flag.String("com1s", "118.50", "COM1 Standby frequency (radio panel)")
		com2Active  = flag.String("com2a", "121.30", "COM2 Active frequency (radio panel)")
		com2Standby = flag.String("com2s", "121.90", "COM2 Standby frequency (radio panel)")
		topRow      = flag.String("top", "250", "Top row display (multi panel)")
		bottomRow   = flag.String("bottom", "3000", "Bottom row display (multi panel)")
		buttonLEDs  = flag.Uint("leds", 0x01, "Button LED states (multi panel)")
	)
	flag.Parse()

	// Parse vendor/product IDs
	var vID, pID uint16
	fmt.Sscanf(*vendorID, "%x", &vID)
	fmt.Sscanf(*productID, "%x", &pID)

	// Add a test function for device open/close
	if *listDevices {
		fmt.Println("Listing all connected HID devices:")
		devices, err := usb.FindDevices()
		if err != nil {
			fmt.Printf("Error enumerating devices: %v\n", err)
			os.Exit(1)
		}
		if len(devices) == 0 {
			fmt.Println("No HID devices found.")
			os.Exit(0)
		}
		for _, d := range devices {
			fmt.Printf("Vendor: 0x%04x Product: 0x%04x Name: %s\n", d.VendorID, d.ProductID, d.Name)
		}
		os.Exit(0)
	}

	if flag.Arg(0) == "test-device" {
		fmt.Printf("Testing open/close for vendor=0x%04x product=0x%04x...\n", vID, pID)
		dev, err := usb.OpenDevice(vID, pID)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("SUCCESS: Device opened.")
		dev.Close()
		fmt.Println("Device closed.")
		os.Exit(0)
	}

	if flag.Arg(0) == "test-communication" {
		fmt.Printf("Testing FIP communication for vendor=0x%04x product=0x%04x...\n", vID, pID)
		dev, err := usb.OpenDevice(vID, pID)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
			os.Exit(1)
		}
		defer dev.Close()

		fmt.Printf("Successfully opened device: %s\n", dev.Name)

		// Test sending control messages
		fmt.Println("Testing control message sending...")
		testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

		for i := 0; i < 3; i++ {
			fmt.Printf("Sending control message %d...\n", i+1)
			err := dev.SendControlMessage(0x21, 0x09, 0x0200, 0, testData)
			if err != nil {
				fmt.Printf("Failed to send control message: %v\n", err)
			} else {
				fmt.Printf("Successfully sent control message %d\n", i+1)
			}
			time.Sleep(500 * time.Millisecond)
		}

		// Test reading data
		fmt.Println("Testing data reading...")
		for i := 0; i < 2; i++ {
			fmt.Printf("Reading bulk data %d...\n", i+1)
			data, err := dev.ReadBulkData(0x81, 64)
			if err != nil {
				fmt.Printf("Failed to read bulk data: %v\n", err)
			} else {
				fmt.Printf("Successfully read %d bytes\n", len(data))
			}
		}

		fmt.Println("FIP communication test completed!")
		os.Exit(0)
	}

	if flag.Arg(0) == "test-direct" {
		fmt.Printf("Testing direct FIP communication for vendor=0x%04x product=0x%04x...\n", vID, pID)

		// Try to find the Saitek FIP directly
		devs := hid.Enumerate(vID, pID)
		if len(devs) == 0 {
			fmt.Println("No Saitek FIP devices found!")
			os.Exit(1)
		}

		fmt.Printf("Found %d Saitek FIP devices\n", len(devs))

		for i, dev := range devs {
			fmt.Printf("  Device %d: Vendor=0x%04x Product=0x%04x Name=%s Path=%s\n",
				i, dev.VendorID, dev.ProductID, dev.Product, dev.Path)

			// Try to open the device
			handle, err := dev.Open()
			if err != nil {
				fmt.Printf("    FAILED to open: %v\n", err)
				continue
			}

			fmt.Printf("    SUCCESS: Device opened\n")

			// Try to send some test data to the FIP
			fmt.Printf("    Sending test data to FIP...\n")

			// Create a simple test pattern
			testData := []byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
			}

			// Try to send the data
			written, err := handle.Write(testData)
			if err != nil {
				fmt.Printf("    Failed to write data: %v\n", err)
			} else {
				fmt.Printf("    Successfully wrote %d bytes to FIP\n", written)
			}

			// Wait a moment
			time.Sleep(1 * time.Second)

			// Try to read some data
			fmt.Printf("    Reading data from FIP...\n")
			readData := make([]byte, 64)
			read, err := handle.Read(readData)
			if err != nil {
				fmt.Printf("    Failed to read data: %v\n", err)
			} else {
				fmt.Printf("    Successfully read %d bytes from FIP: %v\n", read, readData[:read])
			}

			handle.Close()
			fmt.Printf("    Device closed\n")
			break
		}

		fmt.Println("Direct FIP communication test completed!")
		os.Exit(0)
	}

	// Handle radio panel
	if *panelType == "radio" {
		// Parse radio panel vendor/product IDs (default to radio panel)
		if *vendorID == "06a3" && *productID == "0a2ae" {
			*vendorID = "06a3"
			*productID = "0d05" // Radio panel product ID
		}
		fmt.Sscanf(*vendorID, "%x", &vID)
		fmt.Sscanf(*productID, "%x", &pID)

		// Create radio panel
		radio := fip.NewRadioPanelWithUSB(vID, pID)

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

		// Run monitoring loop
		fmt.Printf("\nMonitoring radio panel state... Press Ctrl+C to exit\n")
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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
				radio.Close()
				return
			}
		}
	}

	// Handle multi panel
	if *panelType == "multi" {
		// Parse multi panel vendor/product IDs (default to multi panel)
		if *vendorID == "06a3" && *productID == "0a2ae" {
			*vendorID = "06a3"
			*productID = "0d06" // Multi panel product ID
		}
		fmt.Sscanf(*vendorID, "%x", &vID)
		fmt.Sscanf(*productID, "%x", &pID)

		// Create multi panel
		multi := fip.NewMultiPanelWithUSB(vID, pID)

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

		// Run monitoring loop
		fmt.Printf("\nMonitoring multi panel state... Press Ctrl+C to exit\n")
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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
				multi.Close()
				return
			}
		}
	}

	// Initialize pixelgl for FIP panel
	pixelgl.Run(func() {
		// Create FIP panel with vendor/product IDs
		panel, err := fip.NewFIPPanelWithUSB(*title, *width, *height, vID, pID)
		if err != nil {
			log.Fatalf("Failed to create FIP panel: %v", err)
		}
		defer panel.Close()

		// Try to connect to physical device
		if err := panel.Connect(); err != nil {
			log.Printf("Warning: Could not connect to physical FIP device: %v", err)
			log.Println("Running in virtual mode only")
		}

		// Set instrument type
		switch *instrument {
		case "artificial_horizon":
			panel.SetInstrument(fip.InstrumentArtificialHorizon)
		case "airspeed":
			panel.SetInstrument(fip.InstrumentAirspeed)
		case "altimeter":
			panel.SetInstrument(fip.InstrumentAltimeter)
		case "compass":
			panel.SetInstrument(fip.InstrumentCompass)
		case "vsi":
			panel.SetInstrument(fip.InstrumentVerticalSpeed)
		case "turn_coordinator":
			panel.SetInstrument(fip.InstrumentTurnCoordinator)
		default:
			panel.SetInstrument(fip.InstrumentCustom)
		}

		// Display image if specified
		if *imageFile != "" {
			if err := panel.DisplayImageFromFile(*imageFile); err != nil {
				log.Printf("Failed to display image: %v", err)
			}
		} else {
			// Display test pattern or instrument
			if *instrument == "test" {
				// Create and display test pattern using image generator
				generator := fip.NewImageGenerator(*width, *height)
				testImg := generator.CreateTestPattern()
				if err := panel.DisplayImage(testImg); err != nil {
					log.Printf("Failed to display test pattern: %v", err)
				}
			} else {
				// Display instrument with sample data
				data := fip.InstrumentData{
					Pitch:         5.0,
					Roll:          10.0,
					Airspeed:      120.0,
					Altitude:      5000.0,
					Pressure:      29.92,
					Heading:       180.0,
					VerticalSpeed: 500.0,
					TurnRate:      3.0,
					Slip:          0.0,
				}
				if err := panel.DisplayInstrument(data); err != nil {
					log.Printf("Failed to display instrument: %v", err)
				}
			}
		}

		// Run the display loop
		panel.Run()
	})
}

// Example usage functions
func runFIPExample() {
	fmt.Println("Saitek FIP Controller Example")
	fmt.Println("==============================")

	// This would be called from a separate example program
	// For now, we'll just show the usage
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/main.go -title 'My FIP' -width 320 -height 240")
	fmt.Println("  go run cmd/main.go -instrument artificial_horizon")
	fmt.Println("  go run cmd/main.go -image path/to/image.png")
}

func createSampleImages() {
	fmt.Println("Creating sample images for FIP...")

	// This would create sample instrument images
	// For now, we'll just create a placeholder
	fmt.Println("Sample images would be created here")
}
