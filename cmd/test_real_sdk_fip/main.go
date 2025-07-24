package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"
	"unsafe"

	"saitek-controller/internal/fip"
)

func main() {
	fmt.Println("Real DirectOutput SDK FIP Image Sender")
	fmt.Println("=======================================")

	// Create real DirectOutput SDK instance
	fmt.Println("1. Initializing Real DirectOutput SDK...")
	realSDK, err := fip.NewDirectOutputReal()
	if err != nil {
		log.Fatalf("Failed to create Real DirectOutput SDK: %v", err)
	}
	defer realSDK.Close()

	// Check if we're using the real SDK or fallback
	if realSDK.IsUsingRealSDK() {
		fmt.Println("   ✓ Using REAL DirectOutput SDK")
	} else {
		fmt.Println("   ⚠ Using cross-platform fallback (no real SDK available)")
	}

	// Initialize the SDK
	err = realSDK.Initialize("Real FIP Image Sender")
	if err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
	}

	// Register device callback
	fmt.Println("2. Registering device callbacks...")
	err = realSDK.RegisterDeviceCallback(onDeviceChange, nil)
	if err != nil {
		log.Printf("Warning: Failed to register device callback: %v", err)
	}

	// Enumerate devices
	fmt.Println("3. Enumerating DirectOutput devices...")
	err = realSDK.Enumerate(onDeviceEnumerate, nil)
	if err != nil {
		log.Printf("Warning: Failed to enumerate devices: %v", err)
	}

	// Create a simulated FIP device handle
	deviceHandle := unsafe.Pointer(uintptr(0x12345678))

	// Add a page to the device
	fmt.Println("4. Adding FIP page...")
	err = realSDK.AddPage(deviceHandle, 1, "Real FIP Test Page", fip.FLAG_SET_AS_ACTIVE)
	if err != nil {
		log.Fatalf("Failed to add page: %v", err)
	}

	// Register page callback
	err = realSDK.RegisterPageCallback(deviceHandle, onPageChanged, nil)
	if err != nil {
		log.Printf("Warning: Failed to register page callback: %v", err)
	}

	// Register soft button callback
	err = realSDK.RegisterSoftButtonCallback(deviceHandle, onSoftButtonChanged, nil)
	if err != nil {
		log.Printf("Warning: Failed to register soft button callback: %v", err)
	}

	// Create and send test images
	fmt.Println("5. Creating and sending test images...")

	// Test 1: Simple test image
	fmt.Println("   Sending simple test image...")
	testImage := realSDK.CreateTestImage()
	fipData, err := realSDK.ConvertImageToFIPFormat(testImage)
	if err != nil {
		log.Fatalf("Failed to convert test image: %v", err)
	}

	err = realSDK.SetImage(deviceHandle, 1, 0, fipData)
	if err != nil {
		log.Printf("Warning: Failed to send test image: %v", err)
	} else {
		fmt.Println("   ✓ Simple test image sent")
	}

	// Save the test image for inspection
	err = realSDK.SaveImageAsPNG(testImage, "real_sdk_test_image_1.png")
	if err != nil {
		log.Printf("Warning: Failed to save test image: %v", err)
	} else {
		fmt.Println("   ✓ Test image saved as 'real_sdk_test_image_1.png'")
	}

	// Test 2: Color bars
	fmt.Println("   Sending color bars image...")
	colorImage := createColorBars()
	colorData, err := realSDK.ConvertImageToFIPFormat(colorImage)
	if err != nil {
		log.Fatalf("Failed to convert color image: %v", err)
	}

	err = realSDK.SetImage(deviceHandle, 1, 0, colorData)
	if err != nil {
		log.Printf("Warning: Failed to send color image: %v", err)
	} else {
		fmt.Println("   ✓ Color bars image sent")
	}

	err = realSDK.SaveImageAsPNG(colorImage, "real_sdk_test_image_2.png")
	if err != nil {
		log.Printf("Warning: Failed to save color image: %v", err)
	} else {
		fmt.Println("   ✓ Color image saved as 'real_sdk_test_image_2.png'")
	}

	// Test 3: Gradient
	fmt.Println("   Sending gradient image...")
	gradientImage := createGradient()
	gradientData, err := realSDK.ConvertImageToFIPFormat(gradientImage)
	if err != nil {
		log.Fatalf("Failed to convert gradient image: %v", err)
	}

	err = realSDK.SetImage(deviceHandle, 1, 0, gradientData)
	if err != nil {
		log.Printf("Warning: Failed to send gradient image: %v", err)
	} else {
		fmt.Println("   ✓ Gradient image sent")
	}

	err = realSDK.SaveImageAsPNG(gradientImage, "real_sdk_test_image_3.png")
	if err != nil {
		log.Printf("Warning: Failed to save gradient image: %v", err)
	} else {
		fmt.Println("   ✓ Gradient image saved as 'real_sdk_test_image_3.png'")
	}

	// Test 4: Text pattern
	fmt.Println("   Sending text pattern image...")
	textImage := createTextPattern()
	textData, err := realSDK.ConvertImageToFIPFormat(textImage)
	if err != nil {
		log.Fatalf("Failed to convert text image: %v", err)
	}

	err = realSDK.SetImage(deviceHandle, 1, 0, textData)
	if err != nil {
		log.Printf("Warning: Failed to send text image: %v", err)
	} else {
		fmt.Println("   ✓ Text pattern image sent")
	}

	err = realSDK.SaveImageAsPNG(textImage, "real_sdk_test_image_4.png")
	if err != nil {
		log.Printf("Warning: Failed to save text image: %v", err)
	} else {
		fmt.Println("   ✓ Text image saved as 'real_sdk_test_image_4.png'")
	}

	// Test 5: Complex pattern
	fmt.Println("   Sending complex pattern image...")
	patternImage := createComplexPattern()
	patternData, err := realSDK.ConvertImageToFIPFormat(patternImage)
	if err != nil {
		log.Fatalf("Failed to convert pattern image: %v", err)
	}

	err = realSDK.SetImage(deviceHandle, 1, 0, patternData)
	if err != nil {
		log.Printf("Warning: Failed to send pattern image: %v", err)
	} else {
		fmt.Println("   ✓ Complex pattern image sent")
	}

	err = realSDK.SaveImageAsPNG(patternImage, "real_sdk_test_image_5.png")
	if err != nil {
		log.Printf("Warning: Failed to save pattern image: %v", err)
	} else {
		fmt.Println("   ✓ Pattern image saved as 'real_sdk_test_image_5.png'")
	}

	// Test LED control
	fmt.Println("6. Testing LED control...")
	for i := 0; i < 6; i++ {
		err = realSDK.SetLed(deviceHandle, 1, uint32(i), 1)
		if err != nil {
			log.Printf("Warning: Failed to set LED %d: %v", i, err)
		} else {
			fmt.Printf("   ✓ LED %d turned on\n", i)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Turn off LEDs
	for i := 0; i < 6; i++ {
		err = realSDK.SetLed(deviceHandle, 1, uint32(i), 0)
		if err != nil {
			log.Printf("Warning: Failed to turn off LED %d: %v", i, err)
		}
	}

	// Test image from file
	fmt.Println("7. Testing image from file...")
	
	// Create a test image file first
	testFileImage := createTestFileImage()
	err = realSDK.SaveImageAsPNG(testFileImage, "real_test_fip_image.png")
	if err != nil {
		log.Printf("Warning: Failed to create test file: %v", err)
	} else {
		// Try to load and send the file
		err = realSDK.SetImageFromFile(deviceHandle, 1, 0, "real_test_fip_image.png")
		if err != nil {
			log.Printf("Warning: Failed to send image from file: %v", err)
		} else {
			fmt.Println("   ✓ Image from file sent")
		}
	}

	// Test multiple pages
	fmt.Println("8. Testing multiple pages...")
	
	// Add a second page
	err = realSDK.AddPage(deviceHandle, 2, "Second FIP Page", 0)
	if err != nil {
		log.Printf("Warning: Failed to add second page: %v", err)
	} else {
		fmt.Println("   ✓ Second page added")
		
		// Send a different image to the second page
		page2Image := createPage2Image()
		page2Data, err := realSDK.ConvertImageToFIPFormat(page2Image)
		if err != nil {
			log.Printf("Warning: Failed to convert page 2 image: %v", err)
		} else {
			err = realSDK.SetImage(deviceHandle, 2, 0, page2Data)
			if err != nil {
				log.Printf("Warning: Failed to send page 2 image: %v", err)
			} else {
				fmt.Println("   ✓ Page 2 image sent")
			}
		}
		
		err = realSDK.SaveImageAsPNG(page2Image, "real_sdk_page2_image.png")
		if err != nil {
			log.Printf("Warning: Failed to save page 2 image: %v", err)
		} else {
			fmt.Println("   ✓ Page 2 image saved as 'real_sdk_page2_image.png'")
		}
	}

	// Clean up
	fmt.Println("9. Cleaning up...")
	err = realSDK.RemovePage(deviceHandle, 2)
	if err != nil {
		log.Printf("Warning: Failed to remove page 2: %v", err)
	}
	
	err = realSDK.RemovePage(deviceHandle, 1)
	if err != nil {
		log.Printf("Warning: Failed to remove page 1: %v", err)
	}

	fmt.Println("\n✅ Real DirectOutput SDK FIP Image Sender Test Completed!")
	fmt.Println("\nGenerated test images:")
	fmt.Println("  - real_sdk_test_image_1.png (Simple test)")
	fmt.Println("  - real_sdk_test_image_2.png (Color bars)")
	fmt.Println("  - real_sdk_test_image_3.png (Gradient)")
	fmt.Println("  - real_sdk_test_image_4.png (Text pattern)")
	fmt.Println("  - real_sdk_test_image_5.png (Complex pattern)")
	fmt.Println("  - real_test_fip_image.png (File test)")
	fmt.Println("  - real_sdk_page2_image.png (Page 2 test)")
	fmt.Println("\nThis demonstrates driver-independent FIP image sending using the REAL DirectOutput SDK!")
	
	if realSDK.IsUsingRealSDK() {
		fmt.Println("\n🎉 SUCCESS: Using the REAL DirectOutput SDK!")
		fmt.Println("   This means the SDK was found and loaded successfully.")
		fmt.Println("   On Windows, this would communicate with real FIP hardware.")
	} else {
		fmt.Println("\n⚠️  NOTE: Using cross-platform fallback")
		fmt.Println("   The real DirectOutput SDK was not available.")
		fmt.Println("   This is normal on non-Windows systems or when SDK is not installed.")
		fmt.Println("   The functionality is simulated for testing purposes.")
	}
}

// Callback functions
func onDeviceChange(hDevice unsafe.Pointer, bAdded bool, pCtxt unsafe.Pointer) {
	if bAdded {
		fmt.Printf("   ✓ Device added: %p\n", hDevice)
	} else {
		fmt.Printf("   ✓ Device removed: %p\n", hDevice)
	}
}

func onDeviceEnumerate(hDevice unsafe.Pointer, pCtxt unsafe.Pointer) {
	fmt.Printf("   ✓ Found device: %p\n", hDevice)
}

func onPageChanged(hDevice unsafe.Pointer, dwPage uint32, bSetActive bool, pCtxt unsafe.Pointer) {
	if bSetActive {
		fmt.Printf("   ✓ Page %d activated\n", dwPage)
	} else {
		fmt.Printf("   ✓ Page %d deactivated\n", dwPage)
	}
}

func onSoftButtonChanged(hDevice unsafe.Pointer, dwButtons uint32, pCtxt unsafe.Pointer) {
	fmt.Printf("   ✓ Soft buttons changed: 0x%08X\n", dwButtons)
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

func createGradient() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			r := uint8((x * 255) / 320)
			g := uint8((y * 255) / 240)
			b := uint8(((x + y) * 255) / (320 + 240))
			img.Set(x, y, color.RGBA{r, g, b, 255})
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
	drawText(img, "Real DirectOutput SDK", 160, 40, color.RGBA{255, 255, 255, 255})
	drawText(img, "FIP Image Sender", 160, 60, color.RGBA{255, 255, 0, 255})
	drawText(img, "Driver Independent", 160, 80, color.RGBA{0, 255, 255, 255})
	drawText(img, "320x240 Display", 160, 100, color.RGBA{255, 128, 0, 255})
	drawText(img, "24bpp RGB Format", 160, 120, color.RGBA{128, 255, 0, 255})
	drawText(img, "Test Pattern", 160, 140, color.RGBA{255, 0, 255, 255})
	drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

func createComplexPattern() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with dark background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{10, 10, 20, 255})
		}
	}

	// Draw grid
	for x := 0; x < 320; x += 20 {
		for y := 0; y < 240; y++ {
			img.Set(x, y, color.RGBA{50, 50, 100, 255})
		}
	}
	for y := 0; y < 240; y += 20 {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{50, 50, 100, 255})
		}
	}

	// Draw circles
	for i := 0; i < 5; i++ {
		centerX := 80 + i*40
		centerY := 120
		radius := 15
		color := color.RGBA{uint8(50 + i*40), uint8(100 + i*30), uint8(150 + i*20), 255}

		for y := centerY - radius; y <= centerY + radius; y++ {
			for x := centerX - radius; x <= centerX + radius; x++ {
				if x >= 0 && x < 320 && y >= 0 && y < 240 {
					dx := x - centerX
					dy := y - centerY
					if dx*dx + dy*dy <= radius*radius {
						img.Set(x, y, color)
					}
				}
			}
		}
	}

	// Draw diagonal lines
	for i := 0; i < 320; i += 10 {
		if i < 240 {
			img.Set(i, i, color.RGBA{255, 255, 255, 255})
			img.Set(i+1, i, color.RGBA{255, 255, 255, 255})
			img.Set(i, i+1, color.RGBA{255, 255, 255, 255})
		}
	}

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
	drawText(img, "Real DirectOutput SDK", 160, 100, color.RGBA{255, 255, 0, 255})
	drawText(img, "FIP Ready", 160, 140, color.RGBA{0, 255, 0, 255})

	return img
}

func createPage2Image() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with different background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{40, 20, 60, 255})
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
	drawText(img, "Page 2 Test", 160, 80, color.RGBA{255, 255, 255, 255})
	drawText(img, "Multiple Pages", 160, 100, color.RGBA{255, 255, 0, 255})
	drawText(img, "DirectOutput SDK", 160, 120, color.RGBA{0, 255, 255, 255})
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