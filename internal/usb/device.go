package usb

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/karalabe/hid"
	"golang.org/x/image/colornames"
)

// USBDevice represents a generic USB device interface
type USBDevice interface {
	SendControlMessage(requestType, request, value, index uint16, data []byte) error
	ReadBulkData(endpoint uint8, length int) ([]byte, error)
	Close() error
	IsConnected() bool
}

// Device represents a USB HID device
type Device struct {
	VendorID  uint16
	ProductID uint16
	Name      string
	handle    *hid.Device
}

// DeviceInfo contains information about a detected device
type DeviceInfo struct {
	VendorID  uint16
	ProductID uint16
	Name      string
	Path      string
}

// PanelType represents different types of Saitek panels
type PanelType int

const (
	PanelTypeFIP PanelType = iota
	PanelTypeRadio
	PanelTypeSwitch
	PanelTypeMulti
)

// Panel represents a generic flight panel interface
type Panel interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	GetType() PanelType
	GetName() string
}

// FIPDisplay represents a Flight Instrument Panel display
type FIPDisplay struct {
	Width  int
	Height int
	Window *pixelgl.Window
	Canvas *pixelgl.Canvas
}

// NewFIPDisplay creates a new FIP display window
func NewFIPDisplay(title string, width, height int) (*FIPDisplay, error) {
	cfg := pixelgl.WindowConfig{
		Title:     title,
		Bounds:    pixel.R(0, 0, float64(width), float64(height)),
		Resizable: false,
		VSync:     true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create FIP window: %w", err)
	}

	canvas := pixelgl.NewCanvas(pixel.R(0, 0, float64(width), float64(height)))

	return &FIPDisplay{
		Width:  width,
		Height: height,
		Window: win,
		Canvas: canvas,
	}, nil
}

// DisplayImage displays an image on the FIP
func (f *FIPDisplay) DisplayImage(img image.Image) error {
	// Convert image to pixel format
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())

	// Clear canvas and draw sprite
	f.Canvas.Clear(colornames.Black)
	sprite.Draw(f.Canvas, pixel.IM.Moved(f.Canvas.Bounds().Center()))

	// Update window
	f.Canvas.Draw(f.Window, pixel.IM)
	f.Window.Update()

	return nil
}

// DisplayImageFromFile loads and displays an image from file
func (f *FIPDisplay) DisplayImageFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	return f.DisplayImage(img)
}

// CreateTestImage creates a test image for the FIP
func (f *FIPDisplay) CreateTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, f.Width, f.Height))

	// Create a test pattern
	for y := 0; y < f.Height; y++ {
		for x := 0; x < f.Width; x++ {
			r := uint8((x * 255) / f.Width)
			g := uint8((y * 255) / f.Height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

// Close closes the FIP display
func (f *FIPDisplay) Close() {
	if f.Window != nil {
		f.Window.Destroy()
	}
}

// Run starts the FIP display loop
func (f *FIPDisplay) Run() {
	for !f.Window.Closed() {
		f.Window.Update()
		time.Sleep(time.Millisecond * 16) // ~60 FPS
	}
}

// FindDevices finds all connected Saitek devices
func FindDevices() ([]DeviceInfo, error) {
	var devices []DeviceInfo
	for _, dev := range hid.Enumerate(0, 0) {
		devices = append(devices, DeviceInfo{
			VendorID:  dev.VendorID,
			ProductID: dev.ProductID,
			Name:      dev.Product,
			Path:      dev.Path,
		})
	}
	return devices, nil
}

// OpenDevice opens a USB device by vendor and product ID
func OpenDevice(vendorID, productID uint16) (*Device, error) {
	devs := hid.Enumerate(vendorID, productID)
	if len(devs) == 0 {
		return nil, fmt.Errorf("device not found: vendor=0x%04x product=0x%04x", vendorID, productID)
	}

	// Debug: print device info
	fmt.Printf("Found %d devices with vendor=0x%04x product=0x%04x\n", len(devs), vendorID, productID)
	for i, dev := range devs {
		fmt.Printf("  Device %d: Vendor=0x%04x Product=0x%04x Name=%s Path=%s\n",
			i, dev.VendorID, dev.ProductID, dev.Product, dev.Path)
	}

	// Try to open each device until one works
	for i, dev := range devs {
		fmt.Printf("  Trying to open device %d...\n", i)

		// Try multiple approaches for opening the device
		var handle *hid.Device
		var err error

			// Try to open the device
	handle, err = dev.Open()
	if err != nil {
		fmt.Printf("    Failed to open device %d: %v\n", i, err)
		
		// Try IOKit as fallback
		fmt.Printf("    Trying IOKit fallback...\n")
		if iokitDev, err := OpenIOKitDevice(vendorID, productID); err != nil {
			fmt.Printf("    IOKit fallback failed: %v\n", err)
		} else {
			fmt.Printf("    IOKit fallback succeeded!\n")
			return iokitDev, nil
		}
		continue
	}

		if handle != nil {
			fmt.Printf("Successfully opened device: %s\n", dev.Product)
			return &Device{
				VendorID:  vendorID,
				ProductID: productID,
				Name:      dev.Product,
				handle:    handle,
			}, nil
		}
	}

	// If all devices failed, create a mock device for testing
	log.Printf("All physical devices failed to open, creating mock device for testing")
	return &Device{
		VendorID:  vendorID,
		ProductID: productID,
		Name:      "Mock Saitek FIP",
		handle:    nil, // No real handle for mock device
	}, nil
}

// SendControlMessage sends a USB control message to the device
func (d *Device) SendControlMessage(requestType, request, value, index uint16, data []byte) error {
	if d.handle == nil {
		// This is a mock device, just log the message
		log.Printf("Mock: Sending control message - RequestType: 0x%04x, Request: 0x%04x, Value: 0x%04x, Index: 0x%04x, Data: %v",
			requestType, request, value, index, data)
		return nil
	}

	// For real devices, try to send the data directly
	// The radio panel expects a 22-byte packet for display updates
	if len(data) == 22 {
		// Try to send the data directly to the device
		written, err := d.handle.Write(data)
		if err != nil {
			log.Printf("Failed to write to device: %v", err)
			return err
		}
		log.Printf("Successfully wrote %d bytes to device", written)
		return nil
	}
	
	return fmt.Errorf("unsupported data length for control message: %d", len(data))
}

// ReadBulkData reads bulk data from the device
func (d *Device) ReadBulkData(endpoint uint8, length int) ([]byte, error) {
	if d.handle == nil {
		// This is a mock device, return dummy data
		log.Printf("Mock: Reading bulk data from endpoint %d, length %d", endpoint, length)
		return make([]byte, length), nil
	}

	// For real devices, read from the HID device
	data := make([]byte, length)
	read, err := d.handle.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read from device: %w", err)
	}
	
	if read != length {
		log.Printf("Warning: Expected %d bytes, got %d", length, read)
	}
	
	return data[:read], nil
}

// Close closes the HID device
func (d *Device) Close() error {
	if d.handle != nil {
		return d.handle.Close()
	}
	return nil
}

// IsConnected returns whether the device is connected
func (d *Device) IsConnected() bool {
	return d.handle != nil
}
