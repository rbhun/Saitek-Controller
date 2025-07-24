package main

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"saitek-controller/internal/fip"
)

func main() {
	fmt.Println("Saitek FIP Controller Example")
	fmt.Println("==============================")

	// Initialize pixelgl
	pixelgl.Run(func() {
		// Create FIP panel
		panel, err := fip.NewFIPPanel("FIP Example - Artificial Horizon", 320, 240)
		if err != nil {
			log.Fatalf("Failed to create FIP panel: %v", err)
		}
		defer panel.Close()

		// Try to connect to physical device
		if err := panel.Connect(); err != nil {
			log.Printf("Warning: Could not connect to physical FIP device: %v", err)
			log.Println("Running in virtual mode only")
		}

		// Animation loop
		startTime := time.Now()
		for !panel.display.Window.Closed() {
			// Calculate time-based animation
			elapsed := time.Since(startTime).Seconds()
			
			// Create animated instrument data
			data := fip.InstrumentData{
				Pitch:        10 * math.Sin(elapsed * 0.5),      // Oscillating pitch
				Roll:         15 * math.Sin(elapsed * 0.3),      // Oscillating roll
				Airspeed:     120 + 20*math.Sin(elapsed*0.2),    // Varying airspeed
				Altitude:     5000 + 500*math.Sin(elapsed*0.1),  // Varying altitude
				Pressure:     29.92,
				Heading:      math.Mod(elapsed*10, 360),         // Rotating heading
				VerticalSpeed: 500 * math.Sin(elapsed * 0.4),    // Oscillating vertical speed
				TurnRate:     3 * math.Sin(elapsed * 0.6),       // Varying turn rate
				Slip:         2 * math.Sin(elapsed * 0.7),       // Varying slip
			}

			// Display the instrument
			if err := panel.DisplayInstrument(data); err != nil {
				log.Printf("Failed to display instrument: %v", err)
			}

			// Update window
			panel.display.Window.Update()
			time.Sleep(time.Millisecond * 16) // ~60 FPS
		}
	})
}

// Example of switching between different instruments
func runInstrumentDemo() {
	pixelgl.Run(func() {
		panel, err := fip.NewFIPPanel("FIP Instrument Demo", 320, 240)
		if err != nil {
			log.Fatalf("Failed to create FIP panel: %v", err)
		}
		defer panel.Close()

		instruments := []fip.Instrument{
			fip.InstrumentArtificialHorizon,
			fip.InstrumentAirspeed,
			fip.InstrumentAltimeter,
			fip.InstrumentCompass,
			fip.InstrumentVerticalSpeed,
			fip.InstrumentTurnCoordinator,
		}

		instrumentNames := []string{
			"Artificial Horizon",
			"Airspeed Indicator",
			"Altimeter",
			"Compass",
			"Vertical Speed Indicator",
			"Turn Coordinator",
		}

		startTime := time.Now()
		instrumentIndex := 0

		for !panel.display.Window.Closed() {
			elapsed := time.Since(startTime).Seconds()

			// Switch instruments every 3 seconds
			if int(elapsed/3) != instrumentIndex {
				instrumentIndex = int(elapsed / 3) % len(instruments)
				panel.SetInstrument(instruments[instrumentIndex])
				fmt.Printf("Switched to: %s\n", instrumentNames[instrumentIndex])
			}

			// Create animated data
			data := fip.InstrumentData{
				Pitch:        10 * math.Sin(elapsed * 0.5),
				Roll:         15 * math.Sin(elapsed * 0.3),
				Airspeed:     120 + 20*math.Sin(elapsed*0.2),
				Altitude:     5000 + 500*math.Sin(elapsed*0.1),
				Pressure:     29.92,
				Heading:      math.Mod(elapsed*10, 360),
				VerticalSpeed: 500 * math.Sin(elapsed * 0.4),
				TurnRate:     3 * math.Sin(elapsed * 0.6),
				Slip:         2 * math.Sin(elapsed * 0.7),
			}

			if err := panel.DisplayInstrument(data); err != nil {
				log.Printf("Failed to display instrument: %v", err)
			}

			panel.display.Window.Update()
			time.Sleep(time.Millisecond * 16)
		}
	})
} 