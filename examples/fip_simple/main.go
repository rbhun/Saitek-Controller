package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"saitek-controller/internal/fip"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	fmt.Println("Saitek FIP Simple Example")
	fmt.Println("==========================")

	// Initialize pixelgl
	pixelgl.Run(func() {
		// Create FIP panel with the correct vendor/product IDs
		panel, err := fip.NewFIPPanelWithUSB("FIP Test", 320, 240, 0x06A3, 0xA2AE)
		if err != nil {
			log.Fatalf("Failed to create FIP panel: %v", err)
		}
		defer panel.Close()

		// Try to connect to physical device
		fmt.Println("Attempting to connect to FIP device...")
		if err := panel.Connect(); err != nil {
			log.Printf("Warning: Could not connect to physical FIP device: %v", err)
			log.Println("Running in virtual mode only")
		} else {
			fmt.Println("Successfully connected to FIP device!")
		}

		// Animation loop
		startTime := time.Now()
		window := panel.GetWindow()
		for window != nil && !window.Closed() {
			// Calculate time-based animation
			elapsed := time.Since(startTime).Seconds()

			// Create animated instrument data
			data := fip.InstrumentData{
				Pitch:         10 * math.Sin(elapsed*0.5),       // Oscillating pitch
				Roll:          15 * math.Sin(elapsed*0.3),       // Oscillating roll
				Airspeed:      120 + 20*math.Sin(elapsed*0.2),   // Varying airspeed
				Altitude:      5000 + 500*math.Sin(elapsed*0.1), // Varying altitude
				Pressure:      29.92,
				Heading:       math.Mod(elapsed*10, 360),   // Rotating heading
				VerticalSpeed: 500 * math.Sin(elapsed*0.4), // Oscillating vertical speed
				TurnRate:      3 * math.Sin(elapsed*0.6),   // Varying turn rate
				Slip:          2 * math.Sin(elapsed*0.7),   // Varying slip
			}

			// Display the instrument
			if err := panel.DisplayInstrument(data); err != nil {
				log.Printf("Failed to display instrument: %v", err)
			}

			// Update window
			window.Update()
			time.Sleep(time.Millisecond * 16) // ~60 FPS
		}
	})
}
