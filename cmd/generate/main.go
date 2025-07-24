package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"saitek-controller/internal/fip"
)

func main() {
	var (
		outputDir = flag.String("output", "assets", "Output directory for images")
		width     = flag.Int("width", 320, "Image width")
		height    = flag.Int("height", 240, "Image height")
	)
	flag.Parse()

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Create image generator
	generator := fip.NewImageGenerator(*width, *height)

	fmt.Printf("Generating test images in %s (%dx%d)\n", *outputDir, *width, *height)

	// Generate test images
	if err := generator.GenerateTestImages(*outputDir); err != nil {
		log.Fatalf("Failed to generate test images: %v", err)
	}

	// Generate instrument images with sample data
	instruments := []struct {
		name   string
		inst   fip.Instrument
		data   fip.InstrumentData
	}{
		{
			name: "artificial_horizon",
			inst: fip.InstrumentArtificialHorizon,
			data: fip.InstrumentData{Pitch: 5.0, Roll: 10.0},
		},
		{
			name: "airspeed",
			inst: fip.InstrumentAirspeed,
			data: fip.InstrumentData{Airspeed: 120.0},
		},
		{
			name: "altimeter",
			inst: fip.InstrumentAltimeter,
			data: fip.InstrumentData{Altitude: 5000.0, Pressure: 29.92},
		},
		{
			name: "compass",
			inst: fip.InstrumentCompass,
			data: fip.InstrumentData{Heading: 180.0},
		},
		{
			name: "vsi",
			inst: fip.InstrumentVerticalSpeed,
			data: fip.InstrumentData{VerticalSpeed: 500.0},
		},
		{
			name: "turn_coordinator",
			inst: fip.InstrumentTurnCoordinator,
			data: fip.InstrumentData{TurnRate: 3.0, Slip: 0.0},
		},
	}

	for _, inst := range instruments {
		filename := filepath.Join(*outputDir, inst.name+".png")
		img := generator.CreateInstrumentImage(inst.inst, inst.data)
		
		if err := generator.SaveImage(img, filename); err != nil {
			log.Printf("Failed to save %s: %v", filename, err)
		} else {
			fmt.Printf("Generated: %s\n", filename)
		}
	}

	fmt.Println("Image generation complete!")
	fmt.Printf("Generated images are available in: %s\n", *outputDir)
} 