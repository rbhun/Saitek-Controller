package fip

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"saitek-controller/internal/usb"
)

// FIPPanel represents a Flight Instrument Panel
type FIPPanel struct {
	device     *usb.Device
	display    *usb.FIPDisplay
	connected  bool
	width      int
	height     int
	title      string
	instrument Instrument
	vendorID   uint16
	productID  uint16
}

// Instrument represents different types of flight instruments
type Instrument int

const (
	InstrumentArtificialHorizon Instrument = iota
	InstrumentAirspeed
	InstrumentAltimeter
	InstrumentCompass
	InstrumentVerticalSpeed
	InstrumentTurnCoordinator
	InstrumentCustom
)

// InstrumentData contains data for instrument display
type InstrumentData struct {
	// Artificial Horizon
	Pitch float64 // degrees
	Roll  float64 // degrees
	
	// Airspeed
	Airspeed float64 // knots
	
	// Altimeter
	Altitude float64 // feet
	Pressure float64 // inHg
	
	// Compass
	Heading float64 // degrees
	
	// Vertical Speed
	VerticalSpeed float64 // feet per minute
	
	// Turn Coordinator
	TurnRate float64 // degrees per second
	Slip     float64 // degrees
}

// NewFIPPanel creates a new FIP panel
func NewFIPPanel(title string, width, height int) (*FIPPanel, error) {
	display, err := usb.NewFIPDisplay(title, width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to create FIP display: %w", err)
	}

	return &FIPPanel{
		display:   display,
		width:     width,
		height:    height,
		title:     title,
		instrument: InstrumentCustom,
	}, nil
}

// NewFIPPanelWithUSB creates a new FIP panel with custom vendor/product IDs
func NewFIPPanelWithUSB(title string, width, height int, vendorID, productID uint16) (*FIPPanel, error) {
	display, err := usb.NewFIPDisplay(title, width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to create FIP display: %w", err)
	}

	return &FIPPanel{
		display:   display,
		width:     width,
		height:    height,
		title:     title,
		instrument: InstrumentCustom,
		vendorID:  vendorID,
		productID: productID,
	}, nil
}

// Connect connects to the physical FIP device
func (f *FIPPanel) Connect() error {
	device, err := usb.OpenDevice(f.vendorID, f.productID)
	if err != nil {
		f.connected = false
		return err
	}
	f.device = device
	f.connected = true
	return nil
}

// Disconnect disconnects from the FIP device
func (f *FIPPanel) Disconnect() error {
	if f.device != nil {
		// Close USB connection
		f.device = nil
	}
	f.connected = false
	return nil
}

// IsConnected returns whether the panel is connected
func (f *FIPPanel) IsConnected() bool {
	return f.connected
}

// GetType returns the panel type
func (f *FIPPanel) GetType() usb.PanelType {
	return usb.PanelTypeFIP
}

// GetName returns the panel name
func (f *FIPPanel) GetName() string {
	return f.title
}

// GetWindow returns the display window
func (f *FIPPanel) GetWindow() *pixelgl.Window {
	if f.display != nil {
		return f.display.Window
	}
	return nil
}

// SetInstrument sets the type of instrument to display
func (f *FIPPanel) SetInstrument(instrument Instrument) {
	f.instrument = instrument
}

// DisplayImage displays an image on the FIP
func (f *FIPPanel) DisplayImage(img image.Image) error {
	return f.display.DisplayImage(img)
}

// DisplayImageFromFile loads and displays an image from file
func (f *FIPPanel) DisplayImageFromFile(filename string) error {
	return f.display.DisplayImageFromFile(filename)
}

// DisplayInstrument displays an instrument with the given data
func (f *FIPPanel) DisplayInstrument(data InstrumentData) error {
	var img image.Image

	switch f.instrument {
	case InstrumentArtificialHorizon:
		img = f.createArtificialHorizon(data.Pitch, data.Roll)
	case InstrumentAirspeed:
		img = f.createAirspeedIndicator(data.Airspeed)
	case InstrumentAltimeter:
		img = f.createAltimeter(data.Altitude, data.Pressure)
	case InstrumentCompass:
		img = f.createCompass(data.Heading)
	case InstrumentVerticalSpeed:
		img = f.createVerticalSpeedIndicator(data.VerticalSpeed)
	case InstrumentTurnCoordinator:
		img = f.createTurnCoordinator(data.TurnRate, data.Slip)
	default:
		img = f.createTestPattern()
	}

	return f.DisplayImage(img)
}

// createArtificialHorizon creates an artificial horizon instrument
func (f *FIPPanel) createArtificialHorizon(pitch, roll float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	
	// Clear to blue (sky)
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Skyblue}, image.Point{}, draw.Src)
	
	// Calculate center
	centerX := f.width / 2
	centerY := f.height / 2
	
	// Draw brown earth (bottom half)
	earthRect := image.Rect(0, centerY, f.width, f.height)
	draw.Draw(img, earthRect, &image.Uniform{colornames.Saddlebrown}, image.Point{}, draw.Src)
	
	// Apply roll rotation
	rollRad := roll * math.Pi / 180
	cosRoll := math.Cos(rollRad)
	sinRoll := math.Sin(rollRad)
	
	// Draw horizon line with roll
	for x := 0; x < f.width; x++ {
		// Calculate rotated position
		relX := float64(x - centerX)
		relY := float64(centerY) + pitch*2 // Simple pitch representation
		
		rotX := relX*cosRoll - relY*sinRoll + float64(centerX)
		rotY := relX*sinRoll + relY*cosRoll + float64(centerY)
		
		if rotY >= 0 && rotY < float64(f.height) {
			img.Set(int(rotX), int(rotY), colornames.White)
		}
	}
	
	return img
}

// createAirspeedIndicator creates an airspeed indicator
func (f *FIPPanel) createAirspeedIndicator(airspeed float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	
	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)
	
	centerX := f.width / 2
	centerY := f.height / 2
	radius := f.width / 3
	
	// Draw gauge background
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			dx := x - centerX
			dy := y - centerY
			distance := math.Sqrt(float64(dx*dx + dy*dy))
			
			if distance <= float64(radius) {
				img.Set(x, y, colornames.Darkgray)
			}
		}
	}
	
	// Draw airspeed text
	// Note: In a real implementation, you'd use a proper font rendering library
	// For now, we'll just draw a simple representation
	_ = fmt.Sprintf("%.0f", airspeed)
	
	return img
}

// createAltimeter creates an altimeter
func (f *FIPPanel) createAltimeter(altitude, pressure float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	
	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)
	
	// Draw altimeter face
	centerX := f.width / 2
	centerY := f.height / 2
	radius := f.width / 3
	
	// Draw gauge background
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
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

// createCompass creates a compass
func (f *FIPPanel) createCompass(heading float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	
	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)
	
	centerX := f.width / 2
	centerY := f.height / 2
	radius := f.width / 3
	
	// Draw compass face
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
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

// createVerticalSpeedIndicator creates a vertical speed indicator
func (f *FIPPanel) createVerticalSpeedIndicator(vs float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	
	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)
	
	// Draw vertical speed gauge
	centerX := f.width / 2
	centerY := f.height / 2
	
	// Simple vertical speed representation
	vsHeight := int(vs / 1000 * float64(f.height/2)) // Scale to reasonable range
	vsY := centerY + vsHeight
	
	if vsY >= 0 && vsY < f.height {
		for x := centerX - 20; x <= centerX + 20; x++ {
			if x >= 0 && x < f.width {
				img.Set(x, vsY, colornames.White)
			}
		}
	}
	
	return img
}

// createTurnCoordinator creates a turn coordinator
func (f *FIPPanel) createTurnCoordinator(turnRate, slip float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	
	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)
	
	// Draw turn coordinator representation
	// This is a simplified version - a real turn coordinator is more complex
	
	return img
}

// createTestPattern creates a test pattern
func (f *FIPPanel) createTestPattern() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.width, f.height))
	
	// Create a test pattern
	for y := 0; y < f.height; y++ {
		for x := 0; x < f.width; x++ {
			r := uint8((x * 255) / f.width)
			g := uint8((y * 255) / f.height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
	
	return img
}

// Run starts the FIP display loop
func (f *FIPPanel) Run() {
	f.display.Run()
}

// Close closes the FIP panel
func (f *FIPPanel) Close() {
	if f.display != nil {
		f.display.Close()
	}
	f.Disconnect()
} 