package main

import (
	"flag"
	"log"
	"time"

	"saitek-controller/internal/fip"
)

func main() {
	// Parse command line flags
	vendorID := flag.Uint("vendor", 0x06A3, "USB vendor ID")
	productID := flag.Uint("product", 0x0D67, "USB product ID")
	flag.Parse()

	log.Printf("Testing Saitek Switch Panel")
	log.Printf("Vendor ID: 0x%04X, Product ID: 0x%04X", *vendorID, *productID)

	// Create switch panel
	panel := fip.NewSwitchPanelWithUSB(uint16(*vendorID), uint16(*productID))

	// Connect to the panel
	log.Printf("Connecting to switch panel...")
	if err := panel.Connect(); err != nil {
		log.Printf("Failed to connect to switch panel: %v", err)
		log.Printf("Running in mock mode for testing")
	} else {
		log.Printf("Successfully connected to switch panel")
	}

	// Test landing gear light control
	log.Printf("Testing landing gear light control...")

	// Test 1: All lights off
	log.Printf("Test 1: Setting all lights off")
	if err := panel.SetAllLightsOff(); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(1 * time.Second)

	// Test 2: All green lights (gear down)
	log.Printf("Test 2: Setting all green lights (gear down)")
	if err := panel.SetAllLightsGreen(); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(1 * time.Second)

	// Test 3: All red lights (gear up)
	log.Printf("Test 3: Setting all red lights (gear up)")
	if err := panel.SetAllLightsRed(); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(1 * time.Second)

	// Test 4: All yellow lights (gear transition)
	log.Printf("Test 4: Setting all yellow lights (gear transition)")
	if err := panel.SetAllLightsYellow(); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(1 * time.Second)

	// Test 5: Custom light pattern
	log.Printf("Test 5: Setting custom light pattern")
	customLights := fip.LandingGearLights{
		GreenN: true,  // Green N on
		GreenL: false, // Green L off
		GreenR: true,  // Green R on
		RedN:   false, // Red N off
		RedL:   true,  // Red L on
		RedR:   false, // Red R off
	}
	if err := panel.SetLandingGearLights(customLights); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(1 * time.Second)

	// Test 6: Gear state functions
	log.Printf("Test 6: Testing gear state functions")

	log.Printf("  Setting gear down (green)")
	if err := panel.SetGearDown(); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(1 * time.Second)

	log.Printf("  Setting gear transition (yellow)")
	if err := panel.SetGearTransition(); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(1 * time.Second)

	log.Printf("  Setting gear up (red)")
	if err := panel.SetGearUp(); err != nil {
		log.Printf("Error: %v", err)
	}
	time.Sleep(1 * time.Second)

	// Test 7: Read switch state
	log.Printf("Test 7: Reading switch state")
	if state, err := panel.GetSwitchState(); err != nil {
		log.Printf("Error reading switch state: %v", err)
	} else {
		log.Printf("Switch state: %+v", state)
	}

	// Turn off all lights
	log.Printf("Turning off all lights")
	panel.SetAllLightsOff()

	log.Printf("Switch panel test completed")
	panel.Close()
}
