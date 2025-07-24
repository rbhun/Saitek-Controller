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

	"saitek-controller/internal/usb"
)

// FIPUSB provides USB-based communication with Saitek FIP devices
type FIPUSB struct {
	device    *usb.Device
	connected bool
}

// NewFIPUSB creates a new USB-based FIP instance
func NewFIPUSB() *FIPUSB {
	return &FIPUSB{}
}

// Connect connects to a Saitek FIP device using USB
func (f *FIPUSB) Connect() error {
	// Use our existing USB infrastructure
	device, err := usb.OpenDevice(0x06A3, 0xA2AE) // Saitek FIP vendor/product IDs
	if err != nil {
		return fmt.Errorf("failed to open FIP device: %v", err)
	}

	f.device = device
	f.connected = true
	log.Printf("Successfully connected to FIP device via USB")
	return nil
}

// Disconnect disconnects from the FIP device
func (f *FIPUSB) Disconnect() error {
	if f.device != nil {
		f.device.Close()
		f.device = nil
		f.connected = false
	}
	return nil
}

// IsConnected returns true if connected to a FIP device
func (f *FIPUSB) IsConnected() bool {
	return f.connected && f.device != nil
}

// SendImage sends a 320x240 image to the FIP display
func (f *FIPUSB) SendImage(img image.Image) error {
	if !f.IsConnected() {
		return fmt.Errorf("not connected to FIP device")
	}

	// Convert image to FIP format
	fipData, err := f.convertImageToFIPFormat(img)
	if err != nil {
		return fmt.Errorf("failed to convert image: %v", err)
	}

	// Send the image data to the FIP using USB control message
	return f.sendImageData(fipData)
}

// SendImageFromFile sends an image from a file to the FIP display
func (f *FIPUSB) SendImageFromFile(filename string) error {
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
func (f *FIPUSB) convertImageToFIPFormat(img image.Image) ([]byte, error) {
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

// sendImageData sends image data to the FIP device via USB
func (f *FIPUSB) sendImageData(data []byte) error {
	if len(data) != 320*240*3 {
		return fmt.Errorf("invalid image data size: expected %d, got %d", 320*240*3, len(data))
	}

	// Create the image packet for USB transmission
	packet := f.createImagePacket(data)

	// Send via USB control message
	err := f.device.SendControlMessage(0x21, 0x09, 0x0200, 0, packet)
	if err != nil {
		return fmt.Errorf("failed to send image data via USB: %v", err)
	}

	log.Printf("Sent %d bytes to FIP device via USB", len(packet))
	return nil
}

// createImagePacket creates a packet for sending image data to the FIP
func (f *FIPUSB) createImagePacket(data []byte) []byte {
	// Create a packet with FIP-specific header
	packet := make([]byte, 0, len(data)+16)

	// Add FIP-specific header
	header := []byte{
		0x06, 0xA3, // Vendor ID (Saitek)
		0xA2, 0xAE, // Product ID (FIP)
		0x01,                   // Command: Set Image
		0x00,                   // Page
		0x00,                   // Index
		0x00, 0x00, 0x00, 0x00, // Reserved
	}

	packet = append(packet, header...)

	// Add image data
	packet = append(packet, data...)

	return packet
}

// SetLED sets an LED on the FIP device
func (f *FIPUSB) SetLED(index int, value bool) error {
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

	// Send via USB control message
	err := f.device.SendControlMessage(0x21, 0x09, 0x0200, 0, packet)
	if err != nil {
		return fmt.Errorf("failed to set LED: %v", err)
	}

	log.Printf("Set LED %d to %v via USB", index, value)
	return nil
}

// ReadButtonEvents reads button events from the FIP
func (f *FIPUSB) ReadButtonEvents() (chan FIPUSBButtonEvent, error) {
	if !f.IsConnected() {
		return nil, fmt.Errorf("not connected to FIP device")
	}

	eventChan := make(chan FIPUSBButtonEvent, 10)

	go func() {
		defer close(eventChan)

		// Read from USB device
		for f.IsConnected() {
			data, err := f.device.ReadBulkData(0x81, 64) // Common HID endpoint
			if err != nil {
				log.Printf("Error reading from FIP: %v", err)
				break
			}

			if len(data) > 0 {
				event := f.parseButtonEvent(data)
				if event != nil {
					eventChan <- *event
				}
			}
		}
	}()

	return eventChan, nil
}

// FIPUSBButtonEvent represents a button press event from the FIP via USB
type FIPUSBButtonEvent struct {
	Button    int
	Pressed   bool
	Timestamp time.Time
}

// parseButtonEvent parses a button event from the FIP
func (f *FIPUSB) parseButtonEvent(data []byte) *FIPUSBButtonEvent {
	if len(data) < 8 {
		return nil
	}

	// Check if this is a button event packet
	if data[0] == 0x06 && data[1] == 0xA3 && data[2] == 0xA2 && data[3] == 0xAE {
		if data[4] == 0x03 { // Button event command
			button := int(data[5])
			pressed := data[6] != 0

			return &FIPUSBButtonEvent{
				Button:    button,
				Pressed:   pressed,
				Timestamp: time.Now(),
			}
		}
	}

	return nil
}

// CreateTestImage creates a test image for the FIP
func (f *FIPUSB) CreateTestImage() image.Image {
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
	f.drawText(img, "FIP USB", 160, 60, color.RGBA{255, 255, 255, 255})
	f.drawText(img, "320x240", 160, 80, color.RGBA{255, 255, 0, 255})
	f.drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

// drawText draws simple text on the image
func (f *FIPUSB) drawText(img *image.RGBA, text string, x, y int, c color.Color) {
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
func (f *FIPUSB) SaveImageAsPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// GetDeviceInfo returns information about the connected FIP device
func (f *FIPUSB) GetDeviceInfo() (string, error) {
	if !f.IsConnected() {
		return "", fmt.Errorf("not connected to FIP device")
	}

	return fmt.Sprintf("Saitek FIP - USB Device (VID: 0x06A3, PID: 0xA2AE)"), nil
}
