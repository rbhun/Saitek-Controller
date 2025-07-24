package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"saitek-controller/internal/fip"
)

func main() {
	// Parse command line flags
	vendorID := flag.Uint("vendor", 0x06A3, "USB vendor ID")
	productID := flag.Uint("product", 0x0D67, "USB product ID")
	flag.Parse()

	log.Printf("Starting Saitek Switch Panel Controller")
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

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the panel monitoring in a goroutine
	go panel.Run()

	// Demo the landing gear lights
	log.Printf("Starting landing gear light demonstration...")

	// Test different light patterns
	demoLights(panel)

	// Wait for interrupt signal
	<-sigChan
	log.Printf("Shutting down...")

	// Turn off all lights before closing
	panel.SetAllLightsOff()
	panel.Close()
}

func demoLights(panel *fip.SwitchPanel) {
	// Demo sequence for landing gear lights
	patterns := []struct {
		name string
		fn   func() error
	}{
		{"All lights off", panel.SetAllLightsOff},
		{"All green (gear down)", panel.SetAllLightsGreen},
		{"All red (gear up)", panel.SetAllLightsRed},
		{"All yellow (gear transition)", panel.SetAllLightsYellow},
		{"Custom pattern", func() error {
			lights := fip.LandingGearLights{
				GreenN: true,  // Green N
				GreenL: false, // Green L off
				GreenR: true,  // Green R
				RedN:   false, // Red N off
				RedL:   true,  // Red L
				RedR:   false, // Red R off
			}
			return panel.SetLandingGearLights(lights)
		}},
	}

	for i, pattern := range patterns {
		log.Printf("Demo %d: %s", i+1, pattern.name)
		if err := pattern.fn(); err != nil {
			log.Printf("Error setting lights: %v", err)
		}
		time.Sleep(2 * time.Second)
	}

	// Continuous gear state simulation
	log.Printf("Starting continuous gear state simulation...")
	go func() {
		states := []struct {
			name string
			fn   func() error
		}{
			{"Gear Down", panel.SetGearDown},
			{"Gear Transition", panel.SetGearTransition},
			{"Gear Up", panel.SetGearUp},
			{"Gear Transition", panel.SetGearTransition},
		}

		for {
			for _, state := range states {
				log.Printf("Gear State: %s", state.name)
				if err := state.fn(); err != nil {
					log.Printf("Error setting gear state: %v", err)
				}
				time.Sleep(3 * time.Second)
			}
		}
	}()
}
