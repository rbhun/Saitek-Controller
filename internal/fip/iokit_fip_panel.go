package fip

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"math"
	"sync"
	"time"

	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#import <IOKit/IOKitLib.h>
#import <IOKit/hid/IOHIDLib.h>
#import <CoreFoundation/CoreFoundation.h>

IOHIDDeviceRef findFIPDevice() {
    CFMutableDictionaryRef matchingDict = IOServiceMatching(kIOHIDDeviceKey);
    if (!matchingDict) {
        return NULL;
    }

    // Set vendor and product ID for Saitek FIP
    CFNumberRef vendorID = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &(int){0x06A3});
    CFNumberRef productID = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &(int){0xA2AE});

    CFDictionarySetValue(matchingDict, CFSTR(kIOHIDVendorIDKey), vendorID);
    CFDictionarySetValue(matchingDict, CFSTR(kIOHIDProductIDKey), productID);

    io_iterator_t iterator;
    kern_return_t result = IOServiceGetMatchingServices(kIOMasterPortDefault, matchingDict, &iterator);

    if (result != kIOReturnSuccess) {
        return NULL;
    }

    io_service_t service = IOIteratorNext(iterator);
    IOObjectRelease(iterator);

    if (!service) {
        return NULL;
    }

    IOHIDDeviceRef device = IOHIDDeviceCreate(kCFAllocatorDefault, service);
    IOObjectRelease(service);

    return device;
}

int openFIPDevice(IOHIDDeviceRef device) {
    if (!device) {
        return -1;
    }

    IOReturn result = IOHIDDeviceOpen(device, kIOHIDOptionsTypeNone);
    return (int)result;
}

int readFIPData(IOHIDDeviceRef device, unsigned char* buffer, int bufferSize) {
    if (!device) {
        return -1;
    }

    CFIndex length = bufferSize;
    IOReturn result = IOHIDDeviceGetReport(device, kIOHIDReportTypeInput, 0, buffer, &length);

    if (result == kIOReturnSuccess) {
        return (int)length;
    }

    return -1;
}

void closeFIPDevice(IOHIDDeviceRef device) {
    if (device) {
        IOHIDDeviceClose(device, kIOHIDOptionsTypeNone);
        CFRelease(device);
    }
}
*/
import "C"

// ButtonEvent represents a button press/release event
type ButtonEvent struct {
	ButtonID int  // 0-11 for the 12 FIP buttons
	Pressed  bool // true for press, false for release
}

// IOKitFIPPanel represents a FIP panel using IOKit for device access
type IOKitFIPPanel struct {
	device      C.IOHIDDeviceRef
	isConnected bool
	mu          sync.RWMutex

	// Display properties
	width      int
	height     int
	title      string
	instrument Instrument

	// Button state tracking
	buttonStates   [12]bool // FIP has 12 buttons
	lastButtonData []byte

	// Event channels
	buttonEvents chan ButtonEvent
	stopChan     chan struct{}

	// Display window
	window *pixelgl.Window
}



// NewIOKitFIPPanel creates a new FIP panel using IOKit
func NewIOKitFIPPanel(title string, width, height int) (*IOKitFIPPanel, error) {
	panel := &IOKitFIPPanel{
		width:          width,
		height:         height,
		title:          title,
		instrument:     InstrumentCustom,
		buttonEvents:   make(chan ButtonEvent, 10),
		stopChan:       make(chan struct{}),
		lastButtonData: make([]byte, 2),
	}

	return panel, nil
}

// Connect attempts to connect to the FIP device
func (p *IOKitFIPPanel) Connect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isConnected {
		return fmt.Errorf("already connected")
	}

	// Find the FIP device
	device := C.findFIPDevice()
	if device == 0 {
		return fmt.Errorf("FIP device not found")
	}

	// Try to open the device
	result := C.openFIPDevice(device)
	if result != 0 {
		return fmt.Errorf("failed to open FIP device: error code %d", result)
	}

	p.device = device
	p.isConnected = true

	log.Println("Successfully connected to FIP device via IOKit")

	// Start button monitoring
	go p.monitorButtons()

	return nil
}

// Disconnect closes the connection to the FIP device
func (p *IOKitFIPPanel) Disconnect() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isConnected {
		return nil
	}

	// Signal stop
	close(p.stopChan)

	// Close device
	if p.device != 0 {
		C.closeFIPDevice(p.device)
		p.device = 0
	}

	p.isConnected = false
	log.Println("Disconnected from FIP device")

	return nil
}

// IsConnected returns whether the panel is connected
func (p *IOKitFIPPanel) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.isConnected
}

// GetButtonEvents returns the channel for button events
func (p *IOKitFIPPanel) GetButtonEvents() <-chan ButtonEvent {
	return p.buttonEvents
}

// GetButtonState returns the current state of a button
func (p *IOKitFIPPanel) GetButtonState(buttonID int) bool {
	if buttonID < 0 || buttonID >= 12 {
		return false
	}

	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.buttonStates[buttonID]
}

// SetInstrument sets the type of instrument to display
func (p *IOKitFIPPanel) SetInstrument(instrument Instrument) {
	p.instrument = instrument
}

// GetName returns the panel name
func (p *IOKitFIPPanel) GetName() string {
	return p.title
}

// GetWindow returns the display window
func (p *IOKitFIPPanel) GetWindow() *pixelgl.Window {
	return p.window
}

// DisplayImage displays an image on the FIP
func (p *IOKitFIPPanel) DisplayImage(img image.Image) error {
	// For now, we'll just log that we would display the image
	// In a full implementation, this would send the image to the FIP display
	log.Printf("Would display image of size %dx%d", img.Bounds().Dx(), img.Bounds().Dy())
	return nil
}

// DisplayImageFromFile loads and displays an image from file
func (p *IOKitFIPPanel) DisplayImageFromFile(filename string) error {
	log.Printf("Would display image from file: %s", filename)
	return nil
}

// DisplayInstrument displays an instrument with the given data
func (p *IOKitFIPPanel) DisplayInstrument(data InstrumentData) error {
	var img image.Image

	switch p.instrument {
	case InstrumentArtificialHorizon:
		img = p.createArtificialHorizon(data.Pitch, data.Roll)
	case InstrumentAirspeed:
		img = p.createAirspeedIndicator(data.Airspeed)
	case InstrumentAltimeter:
		img = p.createAltimeter(data.Altitude, data.Pressure)
	case InstrumentCompass:
		img = p.createCompass(data.Heading)
	case InstrumentVerticalSpeed:
		img = p.createVerticalSpeedIndicator(data.VerticalSpeed)
	case InstrumentTurnCoordinator:
		img = p.createTurnCoordinator(data.TurnRate, data.Slip)
	default:
		img = p.createTestPattern()
	}

	return p.DisplayImage(img)
}

// monitorButtons continuously monitors for button presses
func (p *IOKitFIPPanel) monitorButtons() {
	ticker := time.NewTicker(10 * time.Millisecond) // 100Hz polling
	defer ticker.Stop()

	for {
		select {
		case <-p.stopChan:
			return
		case <-ticker.C:
			p.readButtonData()
		}
	}
}

// readButtonData reads button data and processes button events
func (p *IOKitFIPPanel) readButtonData() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isConnected || p.device == 0 {
		return
	}

	// Read button data
	buffer := make([]C.uchar, 2)
	bytesRead := C.readFIPData(p.device, &buffer[0], C.int(len(buffer)))

	if bytesRead > 0 {
		data := make([]byte, bytesRead)
		for i := 0; i < int(bytesRead); i++ {
			data[i] = byte(buffer[i])
		}

		// Process button data
		p.processButtonData(data)
	}
}

// processButtonData interprets the button data and generates events
func (p *IOKitFIPPanel) processButtonData(data []byte) {
	if len(data) < 2 {
		return
	}

	// FIP button mapping (based on typical HID report structure)
	// This may need adjustment based on actual FIP button layout
	buttonMap := []struct {
		byteIndex int
		bitMask   byte
	}{
		{0, 0x01}, // Button 1
		{0, 0x02}, // Button 2
		{0, 0x04}, // Button 3
		{0, 0x08}, // Button 4
		{0, 0x10}, // Button 5
		{0, 0x20}, // Button 6
		{0, 0x40}, // Button 7
		{0, 0x80}, // Button 8
		{1, 0x01}, // Button 9
		{1, 0x02}, // Button 10
		{1, 0x04}, // Button 11
		{1, 0x08}, // Button 12
	}

	// Check for button state changes
	for buttonID, mapping := range buttonMap {
		if buttonID >= len(p.buttonStates) {
			continue
		}

		currentState := (data[mapping.byteIndex] & mapping.bitMask) != 0
		previousState := p.buttonStates[buttonID]

		if currentState != previousState {
			p.buttonStates[buttonID] = currentState

			// Send button event
			select {
			case p.buttonEvents <- ButtonEvent{
				ButtonID: buttonID,
				Pressed:  currentState,
			}:
			default:
				// Channel full, skip this event
			}

			// Log button event
			action := "pressed"
			if !currentState {
				action = "released"
			}
			log.Printf("FIP Button %d %s", buttonID+1, action)
		}
	}

	// Store last data for debugging
	copy(p.lastButtonData, data)
}

// GetLastButtonData returns the last button data for debugging
func (p *IOKitFIPPanel) GetLastButtonData() []byte {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]byte, len(p.lastButtonData))
	copy(result, p.lastButtonData)
	return result
}

// Instrument creation methods (same as original FIP panel)
func (p *IOKitFIPPanel) createArtificialHorizon(pitch, roll float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	// Clear to blue (sky)
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Skyblue}, image.Point{}, draw.Src)

	// Calculate center
	centerX := p.width / 2
	centerY := p.height / 2

	// Draw brown earth (bottom half)
	earthRect := image.Rect(0, centerY, p.width, p.height)
	draw.Draw(img, earthRect, &image.Uniform{colornames.Saddlebrown}, image.Point{}, draw.Src)

	// Apply roll rotation
	rollRad := roll * math.Pi / 180
	cosRoll := math.Cos(rollRad)
	sinRoll := math.Sin(rollRad)

	// Draw horizon line with roll
	for x := 0; x < p.width; x++ {
		// Calculate rotated position
		relX := float64(x - centerX)
		relY := float64(centerY) + pitch*2 // Simple pitch representation

		rotX := relX*cosRoll - relY*sinRoll + float64(centerX)
		rotY := relX*sinRoll + relY*cosRoll + float64(centerY)

		if rotY >= 0 && rotY < float64(p.height) {
			img.Set(int(rotX), int(rotY), colornames.White)
		}
	}

	return img
}

func (p *IOKitFIPPanel) createAirspeedIndicator(airspeed float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)

	centerX := p.width / 2
	centerY := p.height / 2
	radius := p.width / 3

	// Draw gauge background
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
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

func (p *IOKitFIPPanel) createAltimeter(altitude, pressure float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)

	centerX := p.width / 2
	centerY := p.height / 2
	radius := p.width / 3

	// Draw gauge background
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
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

func (p *IOKitFIPPanel) createCompass(heading float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)

	centerX := p.width / 2
	centerY := p.height / 2
	radius := p.width / 3

	// Draw compass face
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
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

func (p *IOKitFIPPanel) createVerticalSpeedIndicator(vs float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)

	centerX := p.width / 2
	centerY := p.height / 2

	// Simple vertical speed representation
	vsHeight := int(vs / 1000 * float64(p.height/2)) // Scale to reasonable range
	vsY := centerY + vsHeight

	if vsY >= 0 && vsY < p.height {
		for x := centerX - 20; x <= centerX+20; x++ {
			if x >= 0 && x < p.width {
				img.Set(x, vsY, colornames.White)
			}
		}
	}

	return img
}

func (p *IOKitFIPPanel) createTurnCoordinator(turnRate, slip float64) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	// Clear to black
	draw.Draw(img, img.Bounds(), &image.Uniform{colornames.Black}, image.Point{}, draw.Src)

	return img
}

func (p *IOKitFIPPanel) createTestPattern() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, p.width, p.height))

	// Create a test pattern
	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			r := uint8((x * 255) / p.width)
			g := uint8((y * 255) / p.height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

// Run starts the FIP display loop
func (p *IOKitFIPPanel) Run() {
	// For now, just keep the program running
	// In a full implementation, this would start the display window
	log.Println("IOKit FIP Panel running...")

	// Keep running until stopped
	<-p.stopChan
}

// Close closes the FIP panel
func (p *IOKitFIPPanel) Close() {
	p.Disconnect()
}
