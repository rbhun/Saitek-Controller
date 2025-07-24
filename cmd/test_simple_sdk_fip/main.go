package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
	"unsafe"

	"saitek-controller/internal/fip"
)

func main() {
	fmt.Println("Simple DirectOutput SDK FIP Test")
	fmt.Println("=================================")

	// Create DirectOutput SDK instance
	fmt.Println("1. Creating DirectOutput SDK...")
	sdk, err := fip.NewDirectOutputReal()
	if err != nil {
		log.Fatalf("Failed to create SDK: %v", err)
	}
	defer sdk.Close()

	// Check if we're using the real SDK or fallback
	if sdk.IsUsingRealSDK() {
		fmt.Println("   ✓ Using REAL DirectOutput SDK")
	} else {
		fmt.Println("   ⚠ Using cross-platform fallback")
	}

	// Initialize the SDK
	fmt.Println("2. Initializing SDK...")
	err = sdk.Initialize("Simple FIP Test")
	if err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
	}

	// Create a simulated device handle
	deviceHandle := unsafe.Pointer(uintptr(0x12345678))

	// Add a page
	fmt.Println("3. Adding FIP page...")
	err = sdk.AddPage(deviceHandle, 1, "Simple Test Page", fip.FLAG_SET_AS_ACTIVE)
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
	err = sdk.SaveImageAsPNG(testImage, "simple_sdk_test_image.png")
	if err != nil {
		log.Printf("Warning: Failed to save test image: %v", err)
	} else {
		fmt.Println("   ✓ Test image saved as 'simple_sdk_test_image.png'")
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

	err = sdk.SaveImageAsPNG(colorImage, "simple_sdk_color_bars.png")
	if err != nil {
		log.Printf("Warning: Failed to save color image: %v", err)
	} else {
		fmt.Println("   ✓ Color image saved as 'simple_sdk_color_bars.png'")
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

	err = sdk.SaveImageAsPNG(textImage, "simple_sdk_text_pattern.png")
	if err != nil {
		log.Printf("Warning: Failed to save text image: %v", err)
	} else {
		fmt.Println("   ✓ Text image saved as 'simple_sdk_text_pattern.png'")
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

	// Test image from file
	fmt.Println("6. Testing image from file...")
	
	// Create a test image file first
	testFileImage := createTestFileImage()
	err = sdk.SaveImageAsPNG(testFileImage, "simple_test_fip_image.png")
	if err != nil {
		log.Printf("Warning: Failed to create test file: %v", err)
	} else {
		// Try to load and send the file
		err = sdk.SetImageFromFile(deviceHandle, 1, 0, "simple_test_fip_image.png")
		if err != nil {
			log.Printf("Warning: Failed to send image from file: %v", err)
		} else {
			fmt.Println("   ✓ Image from file sent")
		}
	}

	// Clean up
	fmt.Println("7. Cleaning up...")
	err = sdk.RemovePage(deviceHandle, 1)
	if err != nil {
		log.Printf("Warning: Failed to remove page: %v", err)
	}

	fmt.Println("\n✅ Simple DirectOutput SDK FIP Test Completed!")
	fmt.Println("\nGenerated test images:")
	fmt.Println("  - simple_sdk_test_image.png (Simple test)")
	fmt.Println("  - simple_sdk_color_bars.png (Color bars)")
	fmt.Println("  - simple_sdk_text_pattern.png (Text pattern)")
	fmt.Println("  - simple_test_fip_image.png (File test)")
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

func createTestFileImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with gradient background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			r := uint8((x * 128) / 320)
			g := uint8((y * 128) / 240)
			b := uint8(64)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Draw border
	for x := 0; x < 320; x++ {
		img.Set(x, 0, color.RGBA{255, 255, 255, 255})
		img.Set(x, 239, color.RGBA{255, 255, 255, 255})
	}
	for y := 0; y < 240; y++ {
		img.Set(0, y, color.RGBA{255, 255, 255, 255})
		img.Set(319, y, color.RGBA{255, 255, 255, 255})
	}

	// Draw text
	drawText(img, "File Test Image", 160, 80, color.RGBA{255, 255, 255, 255})
	drawText(img, "DirectOutput SDK", 160, 100, color.RGBA{255, 255, 0, 255})
	drawText(img, "FIP Ready", 160, 140, color.RGBA{0, 255, 0, 255})

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