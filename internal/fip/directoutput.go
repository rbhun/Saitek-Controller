package fip

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"unsafe"
)

// DirectOutput wrapper for Saitek FIP panels
type DirectOutput struct {
	// Cross-platform implementation
	Devices map[unsafe.Pointer]*Device
}

// Device represents a DirectOutput device
type Device struct {
	Handle     unsafe.Pointer
	DeviceType [16]byte
	Pages      map[uint32]*Page
}

// Page represents a DirectOutput page
type Page struct {
	ID        uint32
	Name      string
	Active    bool
	Images    map[uint32][]byte
	Leds      map[uint32]uint32
	Callbacks *PageCallbacks
}

// PageCallbacks holds callback functions for a page
type PageCallbacks struct {
	OnPageChanged       func(page uint32, active bool)
	OnSoftButtonChanged func(buttons uint32)
}

// Device GUIDs
var (
	DeviceTypeX52Pro = [16]byte{0x06, 0xD5, 0xDA, 0x29, 0x3B, 0xF9, 0x20, 0x4F, 0x85, 0xFA, 0x1E, 0x02, 0xC0, 0x4F, 0xAC, 0x17}
	DeviceTypeFip    = [16]byte{0xD8, 0x3C, 0x08, 0x3E, 0x37, 0x6A, 0x58, 0x4A, 0x80, 0xA8, 0x3D, 0x6A, 0x2C, 0x07, 0x51, 0x3E}
)

// Soft button constants
const (
	SoftButtonSelect = 0x00000001
	SoftButtonUp     = 0x00000002
	SoftButtonDown   = 0x00000004
	SoftButtonLeft   = 0x00000008
	SoftButtonRight  = 0x00000010
	SoftButton1      = 0x00000020
	SoftButton2      = 0x00000040
	SoftButton3      = 0x00000080
	SoftButton4      = 0x00000100
	SoftButton5      = 0x00000200
	SoftButton6      = 0x00000400
)

// Flags
const (
	FLAG_SET_AS_ACTIVE = 0x00000001
)

// Error codes
const (
	E_PAGENOTACTIVE  = 0xFF040001
	E_BUFFERTOOSMALL = 0xFF040000 | 0x6F // ERROR_BUFFER_OVERFLOW
)

// Callback types
type (
	EnumerateCallback        func(hDevice unsafe.Pointer, pCtxt unsafe.Pointer)
	DeviceChangeCallback     func(hDevice unsafe.Pointer, bAdded bool, pCtxt unsafe.Pointer)
	PageChangeCallback       func(hDevice unsafe.Pointer, dwPage uint32, bSetActive bool, pCtxt unsafe.Pointer)
	SoftButtonChangeCallback func(hDevice unsafe.Pointer, dwButtons uint32, pCtxt unsafe.Pointer)
)

// NewDirectOutput creates a new DirectOutput instance
func NewDirectOutput() (*DirectOutput, error) {
	do := &DirectOutput{
		Devices: make(map[unsafe.Pointer]*Device),
	}

	return do, nil
}

// Initialize initializes the DirectOutput library
func (do *DirectOutput) Initialize(pluginName string) error {
	// On macOS, we'll use a cross-platform approach
	// For now, we'll simulate the DirectOutput behavior
	fmt.Printf("DirectOutput initialized with plugin: %s\n", pluginName)
	return nil
}

// Deinitialize cleans up the DirectOutput library
func (do *DirectOutput) Deinitialize() error {
	// Clean up devices
	do.Devices = make(map[unsafe.Pointer]*Device)
	return nil
}

// RegisterDeviceCallback registers a callback for device changes
func (do *DirectOutput) RegisterDeviceCallback(callback DeviceChangeCallback, context unsafe.Pointer) error {
	// Store callback for later use
	return nil
}

// Enumerate enumerates all DirectOutput devices
func (do *DirectOutput) Enumerate(callback EnumerateCallback, context unsafe.Pointer) error {
	// For now, we'll simulate device enumeration
	// In a real implementation, this would find actual FIP devices
	return nil
}

// RegisterPageCallback registers a callback for page changes
func (do *DirectOutput) RegisterPageCallback(hDevice unsafe.Pointer, callback PageChangeCallback, context unsafe.Pointer) error {
	_, exists := do.Devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}
	
	// Store callback for later use
	return nil
}

// RegisterSoftButtonCallback registers a callback for soft button changes
func (do *DirectOutput) RegisterSoftButtonCallback(hDevice unsafe.Pointer, callback SoftButtonChangeCallback, context unsafe.Pointer) error {
	_, exists := do.Devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	// Store callback for later use
	return nil
}

// GetDeviceType gets the device type GUID
func (do *DirectOutput) GetDeviceType(hDevice unsafe.Pointer) ([16]byte, error) {
	device, exists := do.Devices[hDevice]
	if !exists {
		return [16]byte{}, fmt.Errorf("device not found")
	}
	return device.DeviceType, nil
}

// AddPage adds a page to the device
func (do *DirectOutput) AddPage(hDevice unsafe.Pointer, page uint32, debugName string, flags uint32) error {
	device, exists := do.Devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	if device.Pages == nil {
		device.Pages = make(map[uint32]*Page)
	}

	device.Pages[page] = &Page{
		ID:        page,
		Name:      debugName,
		Active:    (flags & FLAG_SET_AS_ACTIVE) != 0,
		Images:    make(map[uint32][]byte),
		Leds:      make(map[uint32]uint32),
		Callbacks: &PageCallbacks{},
	}

	return nil
}

// RemovePage removes a page from the device
func (do *DirectOutput) RemovePage(hDevice unsafe.Pointer, page uint32) error {
	device, exists := do.Devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	delete(device.Pages, page)
	return nil
}

// SetLed sets an LED on the device
func (do *DirectOutput) SetLed(hDevice unsafe.Pointer, page uint32, index uint32, value uint32) error {
	device, exists := do.Devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	pageObj, exists := device.Pages[page]
	if !exists {
		return fmt.Errorf("page not found")
	}

	pageObj.Leds[index] = value
	return nil
}

// SetImage sets an image on the device
func (do *DirectOutput) SetImage(hDevice unsafe.Pointer, page uint32, index uint32, data []byte) error {
	device, exists := do.Devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	pageObj, exists := device.Pages[page]
	if !exists {
		return fmt.Errorf("page not found")
	}

	pageObj.Images[index] = data
	return nil
}

// SetImageFromFile sets an image from a file
func (do *DirectOutput) SetImageFromFile(hDevice unsafe.Pointer, page uint32, index uint32, filename string) error {
	// Read the image file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read image file: %v", err)
	}

	return do.SetImage(hDevice, page, index, data)
}

// CreateFIPImage creates a 320x240 image for the FIP display
func (do *DirectOutput) CreateFIPImage(width, height int) *image.RGBA {
	// FIP displays are 320x240
	if width == 0 {
		width = 320
	}
	if height == 0 {
		height = 240
	}

	return image.NewRGBA(image.Rect(0, 0, width, height))
}

// ConvertImageToFIPFormat converts an image to FIP format (320x240, 24bpp RGB)
func (do *DirectOutput) ConvertImageToFIPFormat(img image.Image) ([]byte, error) {
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

// SaveImageAsPNG saves an image as PNG for debugging
func (do *DirectOutput) SaveImageAsPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// Close closes the DirectOutput instance
func (do *DirectOutput) Close() error {
	return do.Deinitialize()
}
