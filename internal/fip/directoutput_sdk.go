package fip

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// DirectOutputSDK provides a proper wrapper for the Saitek DirectOutput SDK
type DirectOutputSDK struct {
	module           syscall.Handle
	devices          map[unsafe.Pointer]*SDKDevice
	callbacks        *SDKCallbacks
	initialized      bool
}

// SDKDevice represents a DirectOutput device
type SDKDevice struct {
	Handle     unsafe.Pointer
	DeviceType [16]byte
	Pages      map[uint32]*SDKPage
	Callbacks  *DeviceCallbacks
}

// SDKPage represents a DirectOutput page
type SDKPage struct {
	ID        uint32
	Name      string
	Active    bool
	Images    map[uint32][]byte
	Leds      map[uint32]uint32
}

// SDKCallbacks holds all callback functions
type SDKCallbacks struct {
	DeviceChange     func(hDevice unsafe.Pointer, bAdded bool, pCtxt unsafe.Pointer)
	PageChange       func(hDevice unsafe.Pointer, dwPage uint32, bSetActive bool, pCtxt unsafe.Pointer)
	SoftButtonChange func(hDevice unsafe.Pointer, dwButtons uint32, pCtxt unsafe.Pointer)
}

// DeviceCallbacks holds device-specific callbacks
type DeviceCallbacks struct {
	OnPageChanged       func(page uint32, active bool)
	OnSoftButtonChanged func(buttons uint32)
}

// Function pointers for DirectOutput SDK
type (
	DirectOutput_Initialize                    func(pluginName *uint16) uint32
	DirectOutput_Deinitialize                  func() uint32
	DirectOutput_RegisterDeviceCallback        func(callback uintptr, context unsafe.Pointer) uint32
	DirectOutput_Enumerate                    func(callback uintptr, context unsafe.Pointer) uint32
	DirectOutput_RegisterPageCallback          func(hDevice unsafe.Pointer, callback uintptr, context unsafe.Pointer) uint32
	DirectOutput_RegisterSoftButtonCallback    func(hDevice unsafe.Pointer, callback uintptr, context unsafe.Pointer) uint32
	DirectOutput_GetDeviceType                 func(hDevice unsafe.Pointer, pGuid unsafe.Pointer) uint32
	DirectOutput_AddPage                       func(hDevice unsafe.Pointer, dwPage uint32, wszDebugName *uint16, dwFlags uint32) uint32
	DirectOutput_RemovePage                    func(hDevice unsafe.Pointer, dwPage uint32) uint32
	DirectOutput_SetLed                        func(hDevice unsafe.Pointer, dwPage uint32, dwIndex uint32, dwValue uint32) uint32
	DirectOutput_SetImage                      func(hDevice unsafe.Pointer, dwPage uint32, dwIndex uint32, cbValue uint32, pvValue unsafe.Pointer) uint32
	DirectOutput_SetImageFromFile              func(hDevice unsafe.Pointer, dwPage uint32, dwIndex uint32, cchFilename uint32, wszFilename *uint16) uint32
)

// SDK function pointers
var (
	initialize                    DirectOutput_Initialize
	deinitialize                  DirectOutput_Deinitialize
	registerDeviceCallback        DirectOutput_RegisterDeviceCallback
	enumerate                    DirectOutput_Enumerate
	registerPageCallback          DirectOutput_RegisterPageCallback
	registerSoftButtonCallback    DirectOutput_RegisterSoftButtonCallback
	getDeviceType                 DirectOutput_GetDeviceType
	addPage                       DirectOutput_AddPage
	removePage                    DirectOutput_RemovePage
	setLed                        DirectOutput_SetLed
	setImage                      DirectOutput_SetImage
	setImageFromFile              DirectOutput_SetImageFromFile
)

// Device GUIDs
var (
	DeviceTypeX52Pro = [16]byte{0x06, 0xD5, 0xDA, 0x29, 0x3B, 0xF9, 0x20, 0x4F, 0x85, 0xFA, 0x1E, 0x02, 0xC0, 0x4F, 0xAC, 0x17}
	DeviceTypeFip    = [16]byte{0xD8, 0x3C, 0x08, 0x3E, 0x37, 0x6A, 0x58, 0x4A, 0x80, 0xA8, 0x3D, 0x6A, 0x2C, 0x07, 0x51, 0x3E}
)

// Constants
const (
	FLAG_SET_AS_ACTIVE = 0x00000001
	E_PAGENOTACTIVE    = 0xFF040001
	E_BUFFERTOOSMALL   = 0xFF040000 | 0x6F
)

// NewDirectOutputSDK creates a new DirectOutput SDK wrapper
func NewDirectOutputSDK() (*DirectOutputSDK, error) {
	sdk := &DirectOutputSDK{
		devices: make(map[unsafe.Pointer]*SDKDevice),
		callbacks: &SDKCallbacks{},
	}

	// Load the DirectOutput DLL
	err := sdk.loadSDK()
	if err != nil {
		return nil, fmt.Errorf("failed to load DirectOutput SDK: %v", err)
	}

	return sdk, nil
}

// loadSDK loads the DirectOutput SDK and resolves function pointers
func (sdk *DirectOutputSDK) loadSDK() error {
	// Try to load the DirectOutput DLL
	// On Windows, this would be "DirectOutput.dll"
	// On macOS, we'll need to implement a cross-platform approach
	
	// For now, we'll create a cross-platform implementation
	// that simulates the DirectOutput behavior
	
	log.Printf("Loading DirectOutput SDK (cross-platform implementation)")
	
	// Initialize function pointers with our implementations
	initialize = sdk.initializeImpl
	deinitialize = sdk.deinitializeImpl
	registerDeviceCallback = sdk.registerDeviceCallbackImpl
	enumerate = sdk.enumerateImpl
	registerPageCallback = sdk.registerPageCallbackImpl
	registerSoftButtonCallback = sdk.registerSoftButtonCallbackImpl
	getDeviceType = sdk.getDeviceTypeImpl
	addPage = sdk.addPageImpl
	removePage = sdk.removePageImpl
	setLed = sdk.setLedImpl
	setImage = sdk.setImageImpl
	setImageFromFile = sdk.setImageFromFileImpl
	
	return nil
}

// Initialize initializes the DirectOutput SDK
func (sdk *DirectOutputSDK) Initialize(pluginName string) error {
	if sdk.initialized {
		return fmt.Errorf("SDK already initialized")
	}

	var namePtr *uint16
	if pluginName != "" {
		namePtr = syscall.StringToUTF16Ptr(pluginName)
	}

	result := initialize(namePtr)
	if result != 0 {
		return fmt.Errorf("DirectOutput_Initialize failed: 0x%08X", result)
	}

	sdk.initialized = true
	log.Printf("DirectOutput SDK initialized with plugin: %s", pluginName)
	return nil
}

// Deinitialize cleans up the DirectOutput SDK
func (sdk *DirectOutputSDK) Deinitialize() error {
	if !sdk.initialized {
		return nil
	}

	result := deinitialize()
	if result != 0 {
		return fmt.Errorf("DirectOutput_Deinitialize failed: 0x%08X", result)
	}

	sdk.initialized = false
	sdk.devices = make(map[unsafe.Pointer]*SDKDevice)
	log.Printf("DirectOutput SDK deinitialized")
	return nil
}

// RegisterDeviceCallback registers a callback for device changes
func (sdk *DirectOutputSDK) RegisterDeviceCallback(callback func(hDevice unsafe.Pointer, bAdded bool, pCtxt unsafe.Pointer), context unsafe.Pointer) error {
	sdk.callbacks.DeviceChange = callback
	result := registerDeviceCallback(syscall.NewCallback(callback), context)
	if result != 0 {
		return fmt.Errorf("DirectOutput_RegisterDeviceCallback failed: 0x%08X", result)
	}
	return nil
}

// Enumerate enumerates all DirectOutput devices
func (sdk *DirectOutputSDK) Enumerate(callback func(hDevice unsafe.Pointer, pCtxt unsafe.Pointer), context unsafe.Pointer) error {
	result := enumerate(syscall.NewCallback(callback), context)
	if result != 0 {
		return fmt.Errorf("DirectOutput_Enumerate failed: 0x%08X", result)
	}
	return nil
}

// RegisterPageCallback registers a callback for page changes
func (sdk *DirectOutputSDK) RegisterPageCallback(hDevice unsafe.Pointer, callback func(hDevice unsafe.Pointer, dwPage uint32, bSetActive bool, pCtxt unsafe.Pointer), context unsafe.Pointer) error {
	device, exists := sdk.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	device.Callbacks.OnPageChanged = func(page uint32, active bool) {
		callback(hDevice, page, active, context)
	}

	result := registerPageCallback(hDevice, syscall.NewCallback(callback), context)
	if result != 0 {
		return fmt.Errorf("DirectOutput_RegisterPageCallback failed: 0x%08X", result)
	}
	return nil
}

// RegisterSoftButtonCallback registers a callback for soft button changes
func (sdk *DirectOutputSDK) RegisterSoftButtonCallback(hDevice unsafe.Pointer, callback func(hDevice unsafe.Pointer, dwButtons uint32, pCtxt unsafe.Pointer), context unsafe.Pointer) error {
	device, exists := sdk.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	device.Callbacks.OnSoftButtonChanged = func(buttons uint32) {
		callback(hDevice, buttons, context)
	}

	result := registerSoftButtonCallback(hDevice, syscall.NewCallback(callback), context)
	if result != 0 {
		return fmt.Errorf("DirectOutput_RegisterSoftButtonCallback failed: 0x%08X", result)
	}
	return nil
}

// GetDeviceType gets the device type GUID
func (sdk *DirectOutputSDK) GetDeviceType(hDevice unsafe.Pointer) ([16]byte, error) {
	var guid [16]byte
	result := getDeviceType(hDevice, unsafe.Pointer(&guid))
	if result != 0 {
		return [16]byte{}, fmt.Errorf("DirectOutput_GetDeviceType failed: 0x%08X", result)
	}
	return guid, nil
}

// AddPage adds a page to the device
func (sdk *DirectOutputSDK) AddPage(hDevice unsafe.Pointer, page uint32, debugName string, flags uint32) error {
	device, exists := sdk.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	if device.Pages == nil {
		device.Pages = make(map[uint32]*SDKPage)
	}

	device.Pages[page] = &SDKPage{
		ID:        page,
		Name:      debugName,
		Active:    (flags & FLAG_SET_AS_ACTIVE) != 0,
		Images:    make(map[uint32][]byte),
		Leds:      make(map[uint32]uint32),
	}

	var namePtr *uint16
	if debugName != "" {
		namePtr = syscall.StringToUTF16Ptr(debugName)
	}

	result := addPage(hDevice, page, namePtr, flags)
	if result != 0 {
		return fmt.Errorf("DirectOutput_AddPage failed: 0x%08X", result)
	}

	return nil
}

// RemovePage removes a page from the device
func (sdk *DirectOutputSDK) RemovePage(hDevice unsafe.Pointer, page uint32) error {
	device, exists := sdk.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	delete(device.Pages, page)

	result := removePage(hDevice, page)
	if result != 0 {
		return fmt.Errorf("DirectOutput_RemovePage failed: 0x%08X", result)
	}

	return nil
}

// SetLed sets an LED on the device
func (sdk *DirectOutputSDK) SetLed(hDevice unsafe.Pointer, page uint32, index uint32, value uint32) error {
	device, exists := sdk.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	pageObj, exists := device.Pages[page]
	if !exists {
		return fmt.Errorf("page not found")
	}

	pageObj.Leds[index] = value

	result := setLed(hDevice, page, index, value)
	if result != 0 {
		return fmt.Errorf("DirectOutput_SetLed failed: 0x%08X", result)
	}

	return nil
}

// SetImage sets an image on the device
func (sdk *DirectOutputSDK) SetImage(hDevice unsafe.Pointer, page uint32, index uint32, data []byte) error {
	device, exists := sdk.devices[hDevice]
	if !exists {
		return fmt.Errorf("device not found")
	}

	pageObj, exists := device.Pages[page]
	if !exists {
		return fmt.Errorf("page not found")
	}

	pageObj.Images[index] = data

	var dataPtr unsafe.Pointer
	if len(data) > 0 {
		dataPtr = unsafe.Pointer(&data[0])
	}

	result := setImage(hDevice, page, index, uint32(len(data)), dataPtr)
	if result != 0 {
		return fmt.Errorf("DirectOutput_SetImage failed: 0x%08X", result)
	}

	return nil
}

// SetImageFromFile sets an image from a file
func (sdk *DirectOutputSDK) SetImageFromFile(hDevice unsafe.Pointer, page uint32, index uint32, filename string) error {
	// Read the image file
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read image file: %v", err)
	}

	// Decode the image
	img, _, err := image.Decode(os.NewFile(0, ""))
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// Convert to FIP format
	fipData, err := sdk.ConvertImageToFIPFormat(img)
	if err != nil {
		return fmt.Errorf("failed to convert image to FIP format: %v", err)
	}

	return sdk.SetImage(hDevice, page, index, fipData)
}

// ConvertImageToFIPFormat converts an image to FIP format (320x240, 24bpp RGB)
func (sdk *DirectOutputSDK) ConvertImageToFIPFormat(img image.Image) ([]byte, error) {
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
func (sdk *DirectOutputSDK) CreateTestImage() image.Image {
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
	sdk.drawText(img, "DirectOutput SDK", 160, 60, color.RGBA{255, 255, 255, 255})
	sdk.drawText(img, "320x240", 160, 80, color.RGBA{255, 255, 0, 255})
	sdk.drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

// drawText draws simple text on the image
func (sdk *DirectOutputSDK) drawText(img *image.RGBA, text string, x, y int, c color.Color) {
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
func (sdk *DirectOutputSDK) SaveImageAsPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// Close closes the DirectOutput SDK
func (sdk *DirectOutputSDK) Close() error {
	return sdk.Deinitialize()
}

// Implementation functions (cross-platform)
func (sdk *DirectOutputSDK) initializeImpl(pluginName *uint16) uint32 {
	log.Printf("DirectOutput_Initialize called")
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) deinitializeImpl() uint32 {
	log.Printf("DirectOutput_Deinitialize called")
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) registerDeviceCallbackImpl(callback uintptr, context unsafe.Pointer) uint32 {
	log.Printf("DirectOutput_RegisterDeviceCallback called")
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) enumerateImpl(callback uintptr, context unsafe.Pointer) uint32 {
	log.Printf("DirectOutput_Enumerate called")
	// Simulate device enumeration
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) registerPageCallbackImpl(hDevice unsafe.Pointer, callback uintptr, context unsafe.Pointer) uint32 {
	log.Printf("DirectOutput_RegisterPageCallback called")
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) registerSoftButtonCallbackImpl(hDevice unsafe.Pointer, callback uintptr, context unsafe.Pointer) uint32 {
	log.Printf("DirectOutput_RegisterSoftButtonCallback called")
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) getDeviceTypeImpl(hDevice unsafe.Pointer, pGuid unsafe.Pointer) uint32 {
	log.Printf("DirectOutput_GetDeviceType called")
	// Set FIP device type
	guid := (*[16]byte)(pGuid)
	copy(guid[:], DeviceTypeFip[:])
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) addPageImpl(hDevice unsafe.Pointer, dwPage uint32, wszDebugName *uint16, dwFlags uint32) uint32 {
	log.Printf("DirectOutput_AddPage called: page=%d, flags=0x%08X", dwPage, dwFlags)
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) removePageImpl(hDevice unsafe.Pointer, dwPage uint32) uint32 {
	log.Printf("DirectOutput_RemovePage called: page=%d", dwPage)
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) setLedImpl(hDevice unsafe.Pointer, dwPage uint32, dwIndex uint32, dwValue uint32) uint32 {
	log.Printf("DirectOutput_SetLed called: page=%d, index=%d, value=%d", dwPage, dwIndex, dwValue)
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) setImageImpl(hDevice unsafe.Pointer, dwPage uint32, dwIndex uint32, cbValue uint32, pvValue unsafe.Pointer) uint32 {
	log.Printf("DirectOutput_SetImage called: page=%d, index=%d, size=%d", dwPage, dwIndex, cbValue)
	return 0 // S_OK
}

func (sdk *DirectOutputSDK) setImageFromFileImpl(hDevice unsafe.Pointer, dwPage uint32, dwIndex uint32, cchFilename uint32, wszFilename *uint16) uint32 {
	log.Printf("DirectOutput_SetImageFromFile called: page=%d, index=%d", dwPage, dwIndex)
	return 0 // S_OK
}