package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"
	"unsafe"
)

// DirectOutput SDK simulation
type DirectOutputSDK struct {
	useRealSDK bool
	initialized bool
}

// NewDirectOutputSDK creates a new SDK instance
func NewDirectOutputSDK() *DirectOutputSDK {
	sdk := &DirectOutputSDK{}
	
	// Try to detect real SDK (Windows only)
	if runtime.GOOS == "windows" {
		// In a real implementation, this would try to load DirectOutput.dll
		sdk.useRealSDK = false // For demo purposes
	} else {
		sdk.useRealSDK = false
	}
	
	return sdk
}

// Initialize initializes the SDK
func (sdk *DirectOutputSDK) Initialize(pluginName string) error {
	if sdk.initialized {
		return fmt.Errorf("SDK already initialized")
	}
	
	if sdk.useRealSDK {
		log.Printf("Initializing REAL DirectOutput SDK with plugin: %s", pluginName)
	} else {
		log.Printf("Initializing cross-platform DirectOutput SDK with plugin: %s", pluginName)
	}
	
	sdk.initialized = true
	return nil
}

// AddPage adds a page to the device
func (sdk *DirectOutputSDK) AddPage(deviceHandle unsafe.Pointer, page uint32, name string, flags uint32) error {
	if sdk.useRealSDK {
		log.Printf("DirectOutput_AddPage (REAL): device=%p, page=%d, name=%s, flags=0x%08X", deviceHandle, page, name, flags)
	} else {
		log.Printf("DirectOutput_AddPage (simulated): device=%p, page=%d, name=%s, flags=0x%08X", deviceHandle, page, name, flags)
	}
	return nil
}

// SetImage sets an image on the device
func (sdk *DirectOutputSDK) SetImage(deviceHandle unsafe.Pointer, page uint32, index uint32, data []byte) error {
	if sdk.useRealSDK {
		log.Printf("DirectOutput_SetImage (REAL): device=%p, page=%d, index=%d, size=%d", deviceHandle, page, index, len(data))
	} else {
		log.Printf("DirectOutput_SetImage (simulated): device=%p, page=%d, index=%d, size=%d", deviceHandle, page, index, len(data))
	}
	return nil
}

// SetLed sets an LED on the device
func (sdk *DirectOutputSDK) SetLed(deviceHandle unsafe.Pointer, page uint32, index uint32, value uint32) error {
	if sdk.useRealSDK {
		log.Printf("DirectOutput_SetLed (REAL): device=%p, page=%d, index=%d, value=%d", deviceHandle, page, index, value)
	} else {
		log.Printf("DirectOutput_SetLed (simulated): device=%p, page=%d, index=%d, value=%d", deviceHandle, page, index, value)
	}
	return nil
}

// RemovePage removes a page from the device
func (sdk *DirectOutputSDK) RemovePage(deviceHandle unsafe.Pointer, page uint32) error {
	if sdk.useRealSDK {
		log.Printf("DirectOutput_RemovePage (REAL): device=%p, page=%d", deviceHandle, page)
	} else {
		log.Printf("DirectOutput_RemovePage (simulated): device=%p, page=%d", deviceHandle, page)
	}
	return nil
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
	drawText(img, "DirectOutput SDK", 160, 60, color.RGBA{255, 255, 255, 255})
	drawText(img, "320x240", 160, 80, color.RGBA{255, 255, 0, 255})
	drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
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

// IsUsingRealSDK returns true if using the real DirectOutput SDK
func (sdk *DirectOutputSDK) IsUsingRealSDK() bool {
	return sdk.useRealSDK
}

// Close closes the SDK
func (sdk *DirectOutputSDK) Close() error {
	if sdk.initialized {
		log.Printf("DirectOutput SDK closed")
		sdk.initialized = false
	}
	return nil
}

func main() {
	fmt.Println("Standalone DirectOutput SDK FIP Test")
	fmt.Println("====================================")

	// Create DirectOutput SDK instance
	fmt.Println("1. Creating DirectOutput SDK...")
	sdk := NewDirectOutputSDK()
	defer sdk.Close()

	// Check if we're using the real SDK or fallback
	if sdk.IsUsingRealSDK() {
		fmt.Println("   ✓ Using REAL DirectOutput SDK")
	} else {
		fmt.Println("   ⚠ Using cross-platform fallback")
	}

	// Initialize the SDK
	fmt.Println("2. Initializing SDK...")
	err := sdk.Initialize("Standalone FIP Test")
	if err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
	}

	// Create a simulated device handle
	deviceHandle := unsafe.Pointer(uintptr(0x12345678))

	// Add a page
	fmt.Println("3. Adding FIP page...")
	err = sdk.AddPage(deviceHandle, 1, "Standalone Test Page", 0x00000001) // FLAG_SET_AS_ACTIVE
	if err != nil {
		log.Fatalf("Failed to add page: %v", err)
	}

	// Create and send test images
	fmt.Println("4. Creating and sending test images...")

	// Test 1: Simple test image
	fmt.Println("   Creating simple test image...")
	testImage := sdk.CreateTestImage()
	
	// Convert to FIP format
	fipData, err := sdk.ConvertImageToFIPFormat(testImage)
	if err != nil {
		log.Fatalf("Failed to convert test image: %v", err)
	}

	// Send image
	err = sdk.SetImage(deviceHandle, 1, 0, fipData)
	if err != nil {
		log.Printf("Warning: Failed to send test image: %v", err)
	} else {
		fmt.Println("   ✓ Simple test image sent")
	}

	// Save the test image
	err = sdk.SaveImageAsPNG(testImage, "standalone_sdk_test_image.png")
	if err != nil {
		log.Printf("Warning: Failed to save test image: %v", err)
	} else {
		fmt.Println("   ✓ Test image saved as 'standalone_sdk_test_image.png'")
	}

	// Test 2: Color bars
	fmt.Println("   Creating color bars image...")
	colorImage := createColorBars()
	colorData, err := sdk.ConvertImageToFIPFormat(colorImage)
	if err != nil {
		log.Fatalf("Failed to convert color image: %v", err)
	}

	err = sdk.SetImage(deviceHandle, 1, 0, colorData)
	if err != nil {
		log.Printf("Warning: Failed to send color image: %v", err)
	} else {
		fmt.Println("   ✓ Color bars image sent")
	}

	err = sdk.SaveImageAsPNG(colorImage, "standalone_sdk_color_bars.png")
	if err != nil {
		log.Printf("Warning: Failed to save color image: %v", err)
	} else {
		fmt.Println("   ✓ Color image saved as 'standalone_sdk_color_bars.png'")
	}

	// Test 3: Text pattern
	fmt.Println("   Creating text pattern image...")
	textImage := createTextPattern()
	textData, err := sdk.ConvertImageToFIPFormat(textImage)
	if err != nil {
		log.Fatalf("Failed to convert text image: %v", err)
	}

	err = sdk.SetImage(deviceHandle, 1, 0, textData)
	if err != nil {
		log.Printf("Warning: Failed to send text image: %v", err)
	} else {
		fmt.Println("   ✓ Text pattern image sent")
	}

	err = sdk.SaveImageAsPNG(textImage, "standalone_sdk_text_pattern.png")
	if err != nil {
		log.Printf("Warning: Failed to save text image: %v", err)
	} else {
		fmt.Println("   ✓ Text image saved as 'standalone_sdk_text_pattern.png'")
	}

	// Test LED control
	fmt.Println("5. Testing LED control...")
	for i := 0; i < 6; i++ {
		err = sdk.SetLed(deviceHandle, 1, uint32(i), 1)
		if err != nil {
			log.Printf("Warning: Failed to set LED %d: %v", i, err)
		} else {
			fmt.Printf("   ✓ LED %d turned on\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}

	// Turn off LEDs
	for i := 0; i < 6; i++ {
		err = sdk.SetLed(deviceHandle, 1, uint32(i), 0)
		if err != nil {
			log.Printf("Warning: Failed to turn off LED %d: %v", i, err)
		}
	}

	// Clean up
	fmt.Println("6. Cleaning up...")
	err = sdk.RemovePage(deviceHandle, 1)
	if err != nil {
		log.Printf("Warning: Failed to remove page: %v", err)
	}

	fmt.Println("\n✅ Standalone DirectOutput SDK FIP Test Completed!")
	fmt.Println("\nGenerated test images:")
	fmt.Println("  - standalone_sdk_test_image.png (Simple test)")
	fmt.Println("  - standalone_sdk_color_bars.png (Color bars)")
	fmt.Println("  - standalone_sdk_text_pattern.png (Text pattern)")
	fmt.Println("\nThis demonstrates driver-independent FIP image sending!")
	
	if sdk.IsUsingRealSDK() {
		fmt.Println("\n🎉 SUCCESS: Using the REAL DirectOutput SDK!")
	} else {
		fmt.Println("\n⚠️  NOTE: Using cross-platform fallback")
		fmt.Println("   This is normal for development/testing.")
	}
}

// Image creation functions
func createColorBars() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	colors := []color.Color{
		color.RGBA{255, 0, 0, 255},     // Red
		color.RGBA{255, 255, 0, 255},   // Yellow
		color.RGBA{0, 255, 0, 255},     // Green
		color.RGBA{0, 255, 255, 255},   // Cyan
		color.RGBA{0, 0, 255, 255},     // Blue
		color.RGBA{255, 0, 255, 255},   // Magenta
		color.RGBA{255, 255, 255, 255}, // White
		color.RGBA{128, 128, 128, 255}, // Gray
	}

	barWidth := 320 / len(colors)
	for i, c := range colors {
		x1 := i * barWidth
		x2 := (i + 1) * barWidth
		if i == len(colors)-1 {
			x2 = 320
		}

		for y := 0; y < 240; y++ {
			for x := x1; x < x2; x++ {
				img.Set(x, y, c)
			}
		}
	}

	return img
}

func createTextPattern() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with dark background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{20, 20, 40, 255})
		}
	}

	// Draw text pattern
	drawText(img, "DirectOutput SDK", 160, 40, color.RGBA{255, 255, 255, 255})
	drawText(img, "FIP Image Sender", 160, 60, color.RGBA{255, 255, 0, 255})
	drawText(img, "Driver Independent", 160, 80, color.RGBA{0, 255, 255, 255})
	drawText(img, "320x240 Display", 160, 100, color.RGBA{255, 128, 0, 255})
	drawText(img, "24bpp RGB Format", 160, 120, color.RGBA{128, 255, 0, 255})
	drawText(img, "Test Pattern", 160, 140, color.RGBA{255, 0, 255, 255})
	drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

func drawText(img *image.RGBA, text string, x, y int, c color.Color) {
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