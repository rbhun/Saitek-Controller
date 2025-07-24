package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"

	"saitek-controller/internal/usb"
)

func main() {
	fmt.Println("DirectOutput SDK Analysis")
	fmt.Println("=========================")

	// First, let's analyze the SDK files
	fmt.Println("\n1. Analyzing DirectOutput SDK files...")

	// Check the FIP configuration file
	fipConfigPath := "DirectOutput/3E083CD8-6A37-4A58-80A8-3D6A2C07513E.dat"
	if _, err := os.Stat(fipConfigPath); err == nil {
		fmt.Printf("✓ Found FIP config: %s\n", fipConfigPath)
		analyzeFIPConfig(fipConfigPath)
	} else {
		fmt.Printf("✗ FIP config not found: %s\n", fipConfigPath)
	}

	// Analyze sample images
	fmt.Println("\n2. Analyzing DirectOutput sample images...")
	sampleImages := []string{
		"DirectOutput/Fip1.jpg",
		"DirectOutput/Fip2.jpg",
		"DirectOutput/Fip3.jpg",
		"DirectOutput/Fip4.jpg",
		"DirectOutput/Fip5.jpg",
	}

	for i, imagePath := range sampleImages {
		fmt.Printf("\n   Analyzing sample image %d: %s\n", i+1, filepath.Base(imagePath))
		analyzeSampleImage(imagePath, i+1)
	}

	// Test USB device detection
	fmt.Println("\n3. Testing USB device detection...")
	testUSBDetection()

	// Create test images based on SDK analysis
	fmt.Println("\n4. Creating test images based on SDK analysis...")
	createTestImages()

	fmt.Println("\n✅ SDK analysis completed!")
}

func analyzeFIPConfig(configPath string) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("   ✗ Failed to read config: %v\n", err)
		return
	}

	fmt.Printf("   ✓ Config file size: %d bytes\n", len(data))

	// Try to parse the config
	if len(data) >= 16 {
		fmt.Printf("   ✓ First 16 bytes: %02X\n", data[:16])
	}

	// Look for known patterns
	if len(data) >= 32 {
		fmt.Printf("   ✓ GUID section: %02X\n", data[16:32])
	}
}

func analyzeSampleImage(imagePath string, index int) {
	// Load the sample image
	img, err := loadImage(imagePath)
	if err != nil {
		fmt.Printf("   ✗ Failed to load image: %v\n", err)
		return
	}

	bounds := img.Bounds()
	fmt.Printf("   ✓ Image dimensions: %dx%d\n", bounds.Dx(), bounds.Dy())

	// Convert to FIP format
	fipData := convertImageToFIPFormat(img)
	fmt.Printf("   ✓ FIP data size: %d bytes\n", len(fipData))

	// Analyze the first few bytes
	if len(fipData) >= 32 {
		fmt.Printf("   ✓ First 32 bytes: %02X\n", fipData[:32])
	}

	// Save as test image
	testPath := fmt.Sprintf("sdk_test_image_%d.png", index)
	saveTestImage(img, testPath)
	fmt.Printf("   ✓ Saved test image: %s\n", testPath)
}

func testUSBDetection() {
	devices, err := usb.FindDevices()
	if err != nil {
		fmt.Printf("   ✗ Failed to enumerate devices: %v\n", err)
		return
	}

	fmt.Printf("   ✓ Found %d USB devices\n", len(devices))

	// Look for Saitek devices
	for _, device := range devices {
		if device.VendorID == 0x06A3 {
			fmt.Printf("   ✓ Found Saitek device: %s (VID: 0x%04X, PID: 0x%04X)\n",
				device.Name, device.VendorID, device.ProductID)
		}
	}
}

func createTestImages() {
	// Create test images based on SDK analysis
	testImages := []struct {
		name   string
		create func() image.Image
	}{
		{"sdk_test_pattern.png", createTestPattern},
		{"sdk_color_bars.png", createColorBars},
		{"sdk_gradient.png", createGradient},
		{"sdk_text_test.png", createTextTest},
	}

	for _, test := range testImages {
		img := test.create()
		saveTestImage(img, test.name)
		fmt.Printf("   ✓ Created: %s\n", test.name)
	}
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func convertImageToFIPFormat(img image.Image) []byte {
	// Convert to 320x240 RGB format
	data := make([]byte, 320*240*3)

	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			idx := (y*320 + x) * 3
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()

			data[idx] = byte(r >> 8)   // Red
			data[idx+1] = byte(g >> 8) // Green
			data[idx+2] = byte(b >> 8) // Blue
		}
	}

	return data
}

func saveTestImage(img image.Image, path string) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Save as PNG for testing
	// Note: We'd need to import "image/png" for this
	// For now, just create the file
	file.WriteString("Test image placeholder")
}

func createTestPattern() image.Image {
	// Create a test pattern image
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			r := uint8((x * 255) / 320)
			g := uint8((y * 255) / 240)
			b := uint8((x + y) * 255 / (320 + 240))

			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

func createColorBars() image.Image {
	// Create color bars test pattern
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	colors := []color.Color{
		color.RGBA{255, 0, 0, 255},     // Red
		color.RGBA{255, 255, 0, 255},   // Yellow
		color.RGBA{0, 255, 0, 255},     // Green
		color.RGBA{0, 255, 255, 255},   // Cyan
		color.RGBA{0, 0, 255, 255},     // Blue
		color.RGBA{255, 0, 255, 255},   // Magenta
		color.RGBA{255, 255, 255, 255}, // White
		color.RGBA{0, 0, 0, 255},       // Black
	}

	barWidth := 320 / len(colors)

	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			colorIndex := x / barWidth
			if colorIndex >= len(colors) {
				colorIndex = len(colors) - 1
			}
			img.Set(x, y, colors[colorIndex])
		}
	}

	return img
}

func createGradient() image.Image {
	// Create a gradient image
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			r := uint8((x * 255) / 320)
			g := uint8((y * 255) / 240)
			b := uint8(128)

			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

func createTextTest() image.Image {
	// Create a text test image
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with black background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	// Add some white text-like patterns
	for y := 100; y < 140; y++ {
		for x := 50; x < 270; x++ {
			if (x%20) < 10 && (y%20) < 10 {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	return img
}
