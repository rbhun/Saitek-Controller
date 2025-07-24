package fip

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

// DirectOutputReal provides a real DirectOutput SDK implementation
type DirectOutputReal struct {
	module           syscall.Handle
	devices          map[unsafe.Pointer]*RealDevice
	callbacks        *RealCallbacks
	initialized      bool
	useRealSDK       bool
}

// RealDevice represents a real DirectOutput device
type RealDevice struct {
	Handle     unsafe.Pointer
	DeviceType [16]byte
	Pages      map[uint32]*RealPage
	Callbacks  *DeviceCallbacks
}

// RealPage represents a real DirectOutput page
type RealPage struct {
	ID        uint32
	Name      string
	Active    bool
	Images    map[uint32][]byte
	Leds      map[uint32]uint32
}

// RealCallbacks holds real callback functions
type RealCallbacks struct {
	DeviceChange     func(hDevice unsafe.Pointer, bAdded bool, pCtxt unsafe.Pointer)
	PageChange       func(hDevice unsafe.Pointer, dwPage uint32, bSetActive bool, pCtxt unsafe.Pointer)
	SoftButtonChange func(hDevice unsafe.Pointer, dwButtons uint32, pCtxt unsafe.Pointer)
}

// Real SDK function pointers
var (
	realInitialize                    *syscall.Proc
	realDeinitialize                  *syscall.Proc
	realRegisterDeviceCallback        *syscall.Proc
	realEnumerate                    *syscall.Proc
	realRegisterPageCallback          *syscall.Proc
	realRegisterSoftButtonCallback    *syscall.Proc
	realGetDeviceType                 *syscall.Proc
	realAddPage                       *syscall.Proc
	realRemovePage                    *syscall.Proc
	realSetLed                        *syscall.Proc
	realSetImage                      *syscall.Proc
	realSetImageFromFile              *syscall.Proc
)

// NewDirectOutputReal creates a new real DirectOutput SDK wrapper
func NewDirectOutputReal() (*DirectOutputReal, error) {
	real := &DirectOutputReal{
		devices: make(map[unsafe.Pointer]*RealDevice),
		callbacks: &RealCallbacks{},
	}

	// Try to load the real DirectOutput SDK
	err := real.loadRealSDK()
	if err != nil {
		log.Printf("Warning: Failed to load real DirectOutput SDK: %v", err)
		log.Printf("Falling back to cross-platform implementation")
		real.useRealSDK = false
	} else {
		real.useRealSDK = true
		log.Printf("Successfully loaded real DirectOutput SDK")
	}

	return real, nil
}

// loadRealSDK loads the real DirectOutput SDK DLL
func (real *DirectOutputReal) loadRealSDK() error {
	// Only try to load on Windows
	if runtime.GOOS != "windows" {
		return fmt.Errorf("real DirectOutput SDK only available on Windows")
	}

	// Try to load the DirectOutput DLL
	dllPath := "DirectOutput.dll"
	module, err := syscall.LoadDLL(dllPath)
	if err != nil {
		// Try alternative paths
		alternativePaths := []string{
			"./DirectOutput.dll",
			"./DirectOutput/DirectOutput.dll",
			"../DirectOutput.dll",
		}

		for _, path := range alternativePaths {
			module, err = syscall.LoadDLL(path)
			if err == nil {
				break
			}
		}

		if err != nil {
			return fmt.Errorf("failed to load DirectOutput.dll: %v", err)
		}
	}

	real.module = module.Handle

	// Resolve function pointers
	realInitialize = module.MustFindProc("DirectOutput_Initialize")
	realDeinitialize = module.MustFindProc("DirectOutput_Deinitialize")
	realRegisterDeviceCallback = module.MustFindProc("DirectOutput_RegisterDeviceCallback")
	realEnumerate = module.MustFindProc("DirectOutput_Enumerate")
	realRegisterPageCallback = module.MustFindProc("DirectOutput_RegisterPageCallback")
	realRegisterSoftButtonCallback = module.MustFindProc("DirectOutput_RegisterSoftButtonCallback")
	realGetDeviceType = module.MustFindProc("DirectOutput_GetDeviceType")
	realAddPage = module.MustFindProc("DirectOutput_AddPage")
	realRemovePage = module.MustFindProc("DirectOutput_RemovePage")
	realSetLed = module.MustFindProc("DirectOutput_SetLed")
	realSetImage = module.MustFindProc("DirectOutput_SetImage")
	realSetImageFromFile = module.MustFindProc("DirectOutput_SetImageFromFile")

	return nil
}

// Initialize initializes the real DirectOutput SDK
func (real *DirectOutputReal) Initialize(pluginName string) error {
	if real.initialized {
		return fmt.Errorf("SDK already initialized")
	}

	if real.useRealSDK {
		// Use real SDK
		var namePtr *uint16
		if pluginName != "" {
			namePtr, _ = syscall.UTF16PtrFromString(pluginName)
		}

		result, _, _ := realInitialize.Call(uintptr(unsafe.Pointer(namePtr)))
		if result != 0 {
			return fmt.Errorf("DirectOutput_Initialize failed: 0x%08X", result)
		}
	} else {
		// Use cross-platform implementation
		log.Printf("DirectOutput_Initialize (cross-platform): %s", pluginName)
	}

	real.initialized = true
	log.Printf("DirectOutput SDK initialized with plugin: %s", pluginName)
	return nil
}

// Deinitialize cleans up the real DirectOutput SDK
func (real *DirectOutputReal) Deinitialize() error {
	if !real.initialized {
		return nil
	}

	if real.useRealSDK {
		result, _, _ := realDeinitialize.Call()
		if result != 0 {
			return fmt.Errorf("DirectOutput_Deinitialize failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_Deinitialize (cross-platform)")
	}

	real.initialized = false
	real.devices = make(map[unsafe.Pointer]*RealDevice)
	log.Printf("DirectOutput SDK deinitialized")
	return nil
}

// RegisterDeviceCallback registers a callback for device changes
func (real *DirectOutputReal) RegisterDeviceCallback(callback func(hDevice unsafe.Pointer, bAdded bool, pCtxt unsafe.Pointer), context unsafe.Pointer) error {
	real.callbacks.DeviceChange = callback

	if real.useRealSDK {
		callbackPtr := syscall.NewCallback(callback)
		result, _, _ := realRegisterDeviceCallback.Call(callbackPtr, uintptr(context))
		if result != 0 {
			return fmt.Errorf("DirectOutput_RegisterDeviceCallback failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_RegisterDeviceCallback (cross-platform)")
	}

	return nil
}

// Enumerate enumerates all DirectOutput devices
func (real *DirectOutputReal) Enumerate(callback func(hDevice unsafe.Pointer, pCtxt unsafe.Pointer), context unsafe.Pointer) error {
	if real.useRealSDK {
		callbackPtr := syscall.NewCallback(callback)
		result, _, _ := realEnumerate.Call(callbackPtr, uintptr(context))
		if result != 0 {
			return fmt.Errorf("DirectOutput_Enumerate failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_Enumerate (cross-platform)")
		// Simulate device enumeration
		simulatedDevice := unsafe.Pointer(uintptr(0x12345678))
		callback(simulatedDevice, context)
	}

	return nil
}

// RegisterPageCallback registers a callback for page changes
func (real *DirectOutputReal) RegisterPageCallback(hDevice unsafe.Pointer, callback func(hDevice unsafe.Pointer, dwPage uint32, bSetActive bool, pCtxt unsafe.Pointer), context unsafe.Pointer) error {
	device, exists := real.devices[hDevice]
	if !exists {
		// Create device if it doesn't exist (for simulation)
		device = &RealDevice{
			Handle:     hDevice,
			Pages:      make(map[uint32]*RealPage),
			Callbacks:  &DeviceCallbacks{},
		}
		real.devices[hDevice] = device
	}

	device.Callbacks.OnPageChanged = func(page uint32, active bool) {
		callback(hDevice, page, active, context)
	}

	if real.useRealSDK {
		callbackPtr := syscall.NewCallback(callback)
		result, _, _ := realRegisterPageCallback.Call(uintptr(hDevice), callbackPtr, uintptr(context))
		if result != 0 {
			return fmt.Errorf("DirectOutput_RegisterPageCallback failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_RegisterPageCallback (cross-platform)")
	}

	return nil
}

// RegisterSoftButtonCallback registers a callback for soft button changes
func (real *DirectOutputReal) RegisterSoftButtonCallback(hDevice unsafe.Pointer, callback func(hDevice unsafe.Pointer, dwButtons uint32, pCtxt unsafe.Pointer), context unsafe.Pointer) error {
	device, exists := real.devices[hDevice]
	if !exists {
		device = &RealDevice{
			Handle:     hDevice,
			Pages:      make(map[uint32]*RealPage),
			Callbacks:  &DeviceCallbacks{},
		}
		real.devices[hDevice] = device
	}

	device.Callbacks.OnSoftButtonChanged = func(buttons uint32) {
		callback(hDevice, buttons, context)
	}

	if real.useRealSDK {
		callbackPtr := syscall.NewCallback(callback)
		result, _, _ := realRegisterSoftButtonCallback.Call(uintptr(hDevice), callbackPtr, uintptr(context))
		if result != 0 {
			return fmt.Errorf("DirectOutput_RegisterSoftButtonCallback failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_RegisterSoftButtonCallback (cross-platform)")
	}

	return nil
}

// GetDeviceType gets the device type GUID
func (real *DirectOutputReal) GetDeviceType(hDevice unsafe.Pointer) ([16]byte, error) {
	var guid [16]byte

	if real.useRealSDK {
		result, _, _ := realGetDeviceType.Call(uintptr(hDevice), uintptr(unsafe.Pointer(&guid)))
		if result != 0 {
			return [16]byte{}, fmt.Errorf("DirectOutput_GetDeviceType failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_GetDeviceType (cross-platform)")
		// Return FIP device type for simulation
		copy(guid[:], DeviceTypeFip[:])
	}

	return guid, nil
}

// AddPage adds a page to the device
func (real *DirectOutputReal) AddPage(hDevice unsafe.Pointer, page uint32, debugName string, flags uint32) error {
	device, exists := real.devices[hDevice]
	if !exists {
		device = &RealDevice{
			Handle:     hDevice,
			Pages:      make(map[uint32]*RealPage),
			Callbacks:  &DeviceCallbacks{},
		}
		real.devices[hDevice] = device
	}

	if device.Pages == nil {
		device.Pages = make(map[uint32]*RealPage)
	}

	device.Pages[page] = &RealPage{
		ID:        page,
		Name:      debugName,
		Active:    (flags & FLAG_SET_AS_ACTIVE) != 0,
		Images:    make(map[uint32][]byte),
		Leds:      make(map[uint32]uint32),
	}

	if real.useRealSDK {
		var namePtr *uint16
		if debugName != "" {
			namePtr, _ = syscall.UTF16PtrFromString(debugName)
		}

		result, _, _ := realAddPage.Call(uintptr(hDevice), uintptr(page), uintptr(unsafe.Pointer(namePtr)), uintptr(flags))
		if result != 0 {
			return fmt.Errorf("DirectOutput_AddPage failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_AddPage (cross-platform): page=%d, name=%s, flags=0x%08X", page, debugName, flags)
	}

	return nil
}

// RemovePage removes a page from the device
func (real *DirectOutputReal) RemovePage(hDevice unsafe.Pointer, page uint32) error {
	device, exists := real.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	delete(device.Pages, page)

	if real.useRealSDK {
		result, _, _ := realRemovePage.Call(uintptr(hDevice), uintptr(page))
		if result != 0 {
			return fmt.Errorf("DirectOutput_RemovePage failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_RemovePage (cross-platform): page=%d", page)
	}

	return nil
}

// SetLed sets an LED on the device
func (real *DirectOutputReal) SetLed(hDevice unsafe.Pointer, page uint32, index uint32, value uint32) error {
	device, exists := real.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	pageObj, exists := device.Pages[page]
	if !exists {
		return fmt.Errorf("page not found")
	}

	pageObj.Leds[index] = value

	if real.useRealSDK {
		result, _, _ := realSetLed.Call(uintptr(hDevice), uintptr(page), uintptr(index), uintptr(value))
		if result != 0 {
			return fmt.Errorf("DirectOutput_SetLed failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_SetLed (cross-platform): page=%d, index=%d, value=%d", page, index, value)
	}

	return nil
}

// SetImage sets an image on the device
func (real *DirectOutputReal) SetImage(hDevice unsafe.Pointer, page uint32, index uint32, data []byte) error {
	device, exists := real.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	pageObj, exists := device.Pages[page]
	if !exists {
		return fmt.Errorf("page not found")
	}

	pageObj.Images[index] = data

	if real.useRealSDK {
		var dataPtr unsafe.Pointer
		if len(data) > 0 {
			dataPtr = unsafe.Pointer(&data[0])
		}

		result, _, _ := realSetImage.Call(uintptr(hDevice), uintptr(page), uintptr(index), uintptr(len(data)), uintptr(dataPtr))
		if result != 0 {
			return fmt.Errorf("DirectOutput_SetImage failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_SetImage (cross-platform): page=%d, index=%d, size=%d", page, index, len(data))
	}

	return nil
}

// SetImageFromFile sets an image from a file
func (real *DirectOutputReal) SetImageFromFile(hDevice unsafe.Pointer, page uint32, index uint32, filename string) error {
	if real.useRealSDK {
		filenamePtr, _ := syscall.UTF16PtrFromString(filename)
		result, _, _ := realSetImageFromFile.Call(uintptr(hDevice), uintptr(page), uintptr(index), uintptr(len(filename)), uintptr(unsafe.Pointer(filenamePtr)))
		if result != 0 {
			return fmt.Errorf("DirectOutput_SetImageFromFile failed: 0x%08X", result)
		}
	} else {
		log.Printf("DirectOutput_SetImageFromFile (cross-platform): page=%d, index=%d, file=%s", page, index, filename)
		
		// Read the image file and convert it
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open image file: %v", err)
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			return fmt.Errorf("failed to decode image: %v", err)
		}

		fipData, err := real.ConvertImageToFIPFormat(img)
		if err != nil {
			return fmt.Errorf("failed to convert image to FIP format: %v", err)
		}

		return real.SetImage(hDevice, page, index, fipData)
	}

	return nil
}

// ConvertImageToFIPFormat converts an image to FIP format (320x240, 24bpp RGB)
func (real *DirectOutputReal) ConvertImageToFIPFormat(img image.Image) ([]byte, error) {
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

// CreateTestImage creates a test image for the FIP
func (real *DirectOutputReal) CreateTestImage() image.Image {
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
	real.drawText(img, "Real DirectOutput SDK", 160, 60, color.RGBA{255, 255, 255, 255})
	real.drawText(img, "320x240", 160, 80, color.RGBA{255, 255, 0, 255})
	real.drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

// drawText draws simple text on the image
func (real *DirectOutputReal) drawText(img *image.RGBA, text string, x, y int, c color.Color) {
	for i := range text {
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
func (real *DirectOutputReal) SaveImageAsPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// Close closes the real DirectOutput SDK
func (real *DirectOutputReal) Close() error {
	return real.Deinitialize()
}

// IsUsingRealSDK returns true if using the real DirectOutput SDK
func (real *DirectOutputReal) IsUsingRealSDK() bool {
	return real.useRealSDK
}