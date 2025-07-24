package fip

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/karalabe/hid"
)

// FIPDirect provides direct communication with Saitek FIP devices
type FIPDirect struct {
	device    *hid.Device
	connected bool
}

// FIPImage represents a 320x240 image for the FIP display
type FIPImage struct {
	Width  int
	Height int
	Data   []byte // 320x240x3 = 230,400 bytes in RGB format
}

// NewFIPDirect creates a new FIP direct communication instance
func NewFIPDirect() *FIPDirect {
	return &FIPDirect{}
}

// Connect connects to a Saitek FIP device
func (f *FIPDirect) Connect() error {
	// Look for Saitek FIP devices
	devices := hid.Enumerate(0x06A3, 0xA2AE) // Saitek FIP vendor/product IDs

	if len(devices) == 0 {
		return fmt.Errorf("no Saitek FIP devices found")
	}

	// Try to open the first available device
	for i, device := range devices {
		log.Printf("Trying to connect to FIP device %d: %s", i, device.Product)

		hidDevice, err := device.Open()
		if err != nil {
			log.Printf("Failed to open device %d: %v", i, err)
			continue
		}

		f.device = hidDevice
		f.connected = true
		log.Printf("Successfully connected to FIP device: %s", device.Product)
		return nil
	}

	return fmt.Errorf("failed to connect to any FIP device")
}

// Disconnect disconnects from the FIP device
func (f *FIPDirect) Disconnect() error {
	if f.device != nil {
		f.device.Close()
		f.device = nil
		f.connected = false
	}
	return nil
}

// IsConnected returns true if connected to a FIP device
func (f *FIPDirect) IsConnected() bool {
	return f.connected && f.device != nil
}

// SendImage sends a 320x240 image to the FIP display
func (f *FIPDirect) SendImage(img image.Image) error {
	if !f.IsConnected() {
		return fmt.Errorf("not connected to FIP device")
	}

	// Convert image to FIP format
	fipData, err := f.convertImageToFIPFormat(img)
	if err != nil {
		return fmt.Errorf("failed to convert image: %v", err)
	}

	// Send the image data to the FIP
	return f.sendImageData(fipData)
}

// SendImageFromFile sends an image from a file to the FIP display
func (f *FIPDirect) SendImageFromFile(filename string) error {
	// Read the image file
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	return f.SendImage(img)
}

// convertImageToFIPFormat converts an image to FIP format (320x240, 24bpp RGB)
func (f *FIPDirect) convertImageToFIPFormat(img image.Image) ([]byte, error) {
	// Create a 320x240 RGBA image
	fipImg := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Draw the source image onto the FIP image
	draw.Draw(fipImg, fipImg.Bounds(), img, image.Point{}, draw.Src)

	// Convert to 24bpp RGB format (FIP requirement)
	data := make([]byte, 320*240*3)
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			idx := (y*320 + x) * 3
			c := fipImg.RGBAAt(x, y)
			data[idx] = c.R   // Red
			data[idx+1] = c.G // Green
			data[idx+2] = c.B // Blue
		}
	}

	return data, nil
}

// sendImageData sends image data to the FIP device
func (f *FIPDirect) sendImageData(data []byte) error {
	if len(data) != 320*240*3 {
		return fmt.Errorf("invalid image data size: expected %d, got %d", 320*240*3, len(data))
	}

	// The FIP expects the image data in a specific format
	// Based on the DirectOutput SDK, we need to send it as a control message

	// Create the image packet
	packet := f.createImagePacket(data)

	// Send the packet to the device
	written, err := f.device.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to write image data: %v", err)
	}

	log.Printf("Sent %d bytes to FIP device", written)
	return nil
}

// createImagePacket creates a packet for sending image data to the FIP
func (f *FIPDirect) createImagePacket(data []byte) []byte {
	// Based on the DirectOutput SDK and FIP protocol
	// The FIP expects a specific packet format

	// Create a packet with header and image data
	packet := make([]byte, 0, len(data)+64) // Add some header space

	// Add FIP-specific header
	header := []byte{
		0x06, 0xA3, // Vendor ID (Saitek)
		0xA2, 0xAE, // Product ID (FIP)
		0x01, // Command: Set Image
		0x00, // Page
		0x00, // Index
	}

	packet = append(packet, header...)

	// Add image data
	packet = append(packet, data...)

	return packet
}

// SetLED sets an LED on the FIP device
func (f *FIPDirect) SetLED(index int, value bool) error {
	if !f.IsConnected() {
		return fmt.Errorf("not connected to FIP device")
	}

	if index < 0 || index > 5 {
		return fmt.Errorf("invalid LED index: %d (must be 0-5)", index)
	}

	// Create LED control packet
	packet := []byte{
		0x06, 0xA3, // Vendor ID
		0xA2, 0xAE, // Product ID
		0x02,        // Command: Set LED
		byte(index), // LED index
		byte(0),     // LED value (0 = off, 1 = on)
	}
	if value {
		packet[5] = 1
	}

	written, err := f.device.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to set LED: %v", err)
	}

	log.Printf("Set LED %d to %v (%d bytes written)", index, value, written)
	return nil
}

// ReadButtonEvents reads button events from the FIP
func (f *FIPDirect) ReadButtonEvents() (chan FIPButtonEvent, error) {
	if !f.IsConnected() {
		return nil, fmt.Errorf("not connected to FIP device")
	}

	eventChan := make(chan FIPButtonEvent, 10)

	go func() {
		defer close(eventChan)

		buffer := make([]byte, 64) // Standard HID report size

		for f.IsConnected() {
			read, err := f.device.Read(buffer)
			if err != nil {
				log.Printf("Error reading from FIP: %v", err)
				break
			}

			if read > 0 {
				event := f.parseButtonEvent(buffer[:read])
				if event != nil {
					eventChan <- *event
				}
			}
		}
	}()

	return eventChan, nil
}

// FIPButtonEvent represents a button press event from the FIP
type FIPButtonEvent struct {
	Button    int
	Pressed   bool
	Timestamp time.Time
}

// parseButtonEvent parses a button event from the FIP
func (f *FIPDirect) parseButtonEvent(data []byte) *FIPButtonEvent {
	if len(data) < 8 {
		return nil
	}

	// Check if this is a button event packet
	if data[0] == 0x06 && data[1] == 0xA3 && data[2] == 0xA2 && data[3] == 0xAE {
		if data[4] == 0x03 { // Button event command
			button := int(data[5])
			pressed := data[6] != 0

			return &FIPButtonEvent{
				Button:    button,
				Pressed:   pressed,
				Timestamp: time.Now(),
			}
		}
	}

	return nil
}

// CreateTestImage creates a test image for the FIP
func (f *FIPDirect) CreateTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with dark background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{20, 20, 40, 255})
		}
	}

	// Draw a test pattern
	// Border
	for x := 0; x < 320; x++ {
		img.Set(x, 0, color.RGBA{255, 255, 255, 255})
		img.Set(x, 239, color.RGBA{255, 255, 255, 255})
	}
	for y := 0; y < 240; y++ {
		img.Set(0, y, color.RGBA{255, 255, 255, 255})
		img.Set(319, y, color.RGBA{255, 255, 255, 255})
	}

	// Center cross
	for i := 0; i < 320; i++ {
		img.Set(i, 120, color.RGBA{255, 0, 0, 255})
	}
	for i := 0; i < 240; i++ {
		img.Set(160, i, color.RGBA{0, 255, 0, 255})
	}

	// Test text areas
	f.drawText(img, "FIP TEST", 160, 60, color.RGBA{255, 255, 255, 255})
	f.drawText(img, "320x240", 160, 80, color.RGBA{255, 255, 0, 255})
	f.drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

// drawText draws simple text on the image
func (f *FIPDirect) drawText(img *image.RGBA, text string, x, y int, c color.Color) {
	for i, _ := range text {
		charX := x + i*8 - len(text)*4
		if charX >= 0 && charX < 320 {
			img.Set(charX, y, c)
			img.Set(charX+1, y, c)
			img.Set(charX, y+1, c)
			img.Set(charX+1, y+1, c)
		}
	}
}

// SaveImageAsPNG saves an image as PNG for debugging
func (f *FIPDirect) SaveImageAsPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// GetDeviceInfo returns information about the connected FIP device
func (f *FIPDirect) GetDeviceInfo() (string, error) {
	if !f.IsConnected() {
		return "", fmt.Errorf("not connected to FIP device")
	}

	// Try to get device info
	devices := hid.Enumerate(0x06A3, 0xA2AE)
	for _, device := range devices {
		if device.Product != "" {
			return fmt.Sprintf("Saitek FIP - %s (VID: 0x%04X, PID: 0x%04X)",
				device.Product, device.VendorID, device.ProductID), nil
		}
	}

	return "Saitek FIP Device", nil
}
