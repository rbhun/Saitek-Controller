package fip

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"golang.org/x/image/colornames"
)

// ImageGenerator creates test images for FIP panels
type ImageGenerator struct {
	width  int
	height int
}

// NewImageGenerator creates a new image generator
func NewImageGenerator(width, height int) *ImageGenerator {
	return &ImageGenerator{
		width:  width,
		height: height,
	}
}

// CreateTestPattern creates a test pattern image
func (g *ImageGenerator) CreateTestPattern() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	
	// Create a color test pattern
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			r := uint8((x * 255) / g.width)
			g := uint8((y * 255) / g.height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	
	return img
}

// CreateColorBars creates a color bars test pattern
func (g *ImageGenerator) CreateColorBars() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	
	colors := []color.Color{
		colornames.White,
		colornames.Yellow,
		colornames.Cyan,
		colornames.Green,
		colornames.Magenta,
		colornames.Red,
		colornames.Blue,
		colornames.Black,
	}
	
	barWidth := g.width / len(colors)
	
	for i, c := range colors {
		x1 := i * barWidth
		x2 := (i + 1) * barWidth
		if i == len(colors)-1 {
			x2 = g.width
		}
		
		for y := 0; y < g.height; y++ {
			for x := x1; x < x2; x++ {
				img.Set(x, y, c)
			}
		}
	}
	
	return img
}

// CreateGradient creates a gradient image
func (g *ImageGenerator) CreateGradient() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			r := uint8((x * 255) / g.width)
			g := uint8((y * 255) / g.height)
			b := uint8(255 - r)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	
	return img
}

// CreateInstrumentBackground creates a background for instruments
func (g *ImageGenerator) CreateInstrumentBackground() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	
	// Fill with dark background
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)
	
	// Draw instrument bezel
	centerX := g.width / 2
	centerY := g.height / 2
	radius := g.width / 3
	
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			dx := x - centerX
			dy := y - centerY
			distance := math.Sqrt(float64(dx*dx + dy*dy))
			
			if distance <= float64(radius) {
				img.Set(x, y, colornames.Darkgray)
			}
		}
	}
	
	return img
}

// SaveImage saves an image to a file
func (g *ImageGenerator) SaveImage(img image.Image, filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode image: %w", err)
	}
	
	return nil
}

// GenerateTestImages generates a set of test images
func (g *ImageGenerator) GenerateTestImages(outputDir string) error {
	images := map[string]image.Image{
		"test_pattern.png":    g.CreateTestPattern(),
		"color_bars.png":      g.CreateColorBars(),
		"gradient.png":        g.CreateGradient(),
		"instrument_bg.png":   g.CreateInstrumentBackground(),
	}
	
	for filename, img := range images {
		fullPath := filepath.Join(outputDir, filename)
		if err := g.SaveImage(img, fullPath); err != nil {
			return fmt.Errorf("failed to save %s: %w", filename, err)
		}
		fmt.Printf("Generated: %s\n", fullPath)
	}
	
	return nil
}

// CreateInstrumentImage creates an instrument image with data
func (g *ImageGenerator) CreateInstrumentImage(instrument Instrument, data InstrumentData) image.Image {
	switch instrument {
	case InstrumentArtificialHorizon:
		return g.createArtificialHorizonImage(data.Pitch, data.Roll)
	case InstrumentAirspeed:
		return g.createAirspeedImage(data.Airspeed)
	case InstrumentAltimeter:
		return g.createAltimeterImage(data.Altitude, data.Pressure)
	case InstrumentCompass:
		return g.createCompassImage(data.Heading)
	case InstrumentVerticalSpeed:
		return g.createVSIImage(data.VerticalSpeed)
	case InstrumentTurnCoordinator:
		return g.createTurnCoordinatorImage(data.TurnRate, data.Slip)
	default:
		return g.CreateTestPattern()
	}
}

// createArtificialHorizonImage creates an artificial horizon image
func (g *ImageGenerator) createArtificialHorizonImage(pitch, roll float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	
	// Clear to blue (sky)
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Skyblue}, image.Point{}, draw.Src)
	
	// Calculate center
	centerX := g.width / 2
	centerY := g.height / 2
	
	// Draw brown earth (bottom half)
	earthRect := image.Rect(0, centerY, g.width, g.height)
	draw.Draw(img, earthRect, &image.Uniform{colornames.Saddlebrown}, image.Point{}, draw.Src)
	
	// Apply roll rotation
	rollRad := roll * math.Pi / 180
	cosRoll := math.Cos(rollRad)
	sinRoll := math.Sin(rollRad)
	
	// Draw horizon line with roll
	for x := 0; x < g.width; x++ {
		// Calculate rotated position
		relX := float64(x - centerX)
		relY := float64(centerY) + pitch*2 // Simple pitch representation
		
		rotX := relX*cosRoll - relY*sinRoll + float64(centerX)
		rotY := relX*sinRoll + relY*cosRoll + float64(centerY)
		
		if rotY >= 0 && rotY < float64(g.height) {
			img.Set(int(rotX), int(rotY), colornames.White)
		}
	}
	
	return img
}

// createAirspeedImage creates an airspeed indicator image
func (g *ImageGenerator) createAirspeedImage(airspeed float64) image.Image {
	img := g.CreateInstrumentBackground()
	rgba := img.(*image.RGBA)
	
	// Draw airspeed needle
	centerX := g.width / 2
	centerY := g.height / 2
	radius := g.width / 3
	
	// Calculate needle angle (0-360 degrees)
	angle := (airspeed / 200.0) * 360.0 // Scale to 0-200 knots
	angleRad := angle * math.Pi / 180
	
	// Draw needle
	needleLength := float64(radius) * 0.8
	
	// Draw needle line
	for i := 0; i < int(needleLength); i++ {
		x := centerX + int(float64(i)*math.Sin(angleRad))
		y := centerY - int(float64(i)*math.Cos(angleRad))
		if x >= 0 && x < g.width && y >= 0 && y < g.height {
			rgba.Set(x, y, colornames.White)
		}
	}
	
	return img
}

// createAltimeterImage creates an altimeter image
func (g *ImageGenerator) createAltimeterImage(altitude, pressure float64) image.Image {
	img := g.CreateInstrumentBackground()
	rgba := img.(*image.RGBA)
	
	// Draw altimeter face
	centerX := g.width / 2
	centerY := g.height / 2
	radius := g.width / 3
	
	// Draw altitude needle
	altitudeAngle := (altitude / 10000.0) * 360.0 // Scale to 0-10,000 feet
	altitudeRad := altitudeAngle * math.Pi / 180
	
	needleLength := float64(radius) * 0.8
	
	// Draw needle
	for i := 0; i < int(needleLength); i++ {
		x := centerX + int(float64(i)*math.Sin(altitudeRad))
		y := centerY - int(float64(i)*math.Cos(altitudeRad))
		if x >= 0 && x < g.width && y >= 0 && y < g.height {
			rgba.Set(x, y, colornames.White)
		}
	}
	
	return img
}

// createCompassImage creates a compass image
func (g *ImageGenerator) createCompassImage(heading float64) image.Image {
	img := g.CreateInstrumentBackground()
	rgba := img.(*image.RGBA)
	
	// Draw compass face
	centerX := g.width / 2
	centerY := g.height / 2
	radius := g.width / 3
	
	// Draw heading indicator
	headingRad := heading * math.Pi / 180
	indicatorLength := float64(radius) * 0.9
	
	// Draw heading line
	for i := 0; i < int(indicatorLength); i++ {
		x := centerX + int(float64(i)*math.Sin(headingRad))
		y := centerY - int(float64(i)*math.Cos(headingRad))
		if x >= 0 && x < g.width && y >= 0 && y < g.height {
			rgba.Set(x, y, colornames.White)
		}
	}
	
	return img
}

// createVSIImage creates a vertical speed indicator image
func (g *ImageGenerator) createVSIImage(vs float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, g.width, g.height))
	
	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)
	
	// Draw vertical speed gauge
	centerX := g.width / 2
	centerY := g.height / 2
	
	// Simple vertical speed representation
	vsHeight := int(vs / 1000 * float64(g.height/2)) // Scale to reasonable range
	vsY := centerY + vsHeight
	
	if vsY >= 0 && vsY < g.height {
		for x := centerX - 20; x <= centerX + 20; x++ {
			if x >= 0 && x < g.width {
				img.Set(x, vsY, colornames.White)
			}
		}
	}
	
	return img
}

// createTurnCoordinatorImage creates a turn coordinator image
func (g *ImageGenerator) createTurnCoordinatorImage(turnRate, slip float64) image.Image {
	img := g.CreateInstrumentBackground()
	
	// This is a simplified turn coordinator
	// In a real implementation, this would be more complex
	
	return img
} 