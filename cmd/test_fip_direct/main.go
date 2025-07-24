package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"saitek-controller/internal/fip"
)

func main() {
	fmt.Println("Direct FIP Image Sender")
	fmt.Println("========================")

	// Create FIP direct instance
	fipDirect := fip.NewFIPDirect()

	// Connect to FIP device
	fmt.Println("1. Connecting to FIP device...")
	err := fipDirect.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to FIP: %v", err)
	}
	defer fipDirect.Disconnect()

	fmt.Println("✓ Connected to FIP device")

	// Get device info
	deviceInfo, err := fipDirect.GetDeviceInfo()
	if err != nil {
		log.Printf("Warning: Could not get device info: %v", err)
	} else {
		fmt.Printf("Device: %s\n", deviceInfo)
	}

	// Create a test image
	fmt.Println("\n2. Creating test image...")
	testImage := createTestImage()

	// Save test image for inspection
	err = fipDirect.SaveImageAsPNG(testImage, "fip_direct_test_image.png")
	if err != nil {
		log.Printf("Warning: Could not save test image: %v", err)
	} else {
		fmt.Println("✓ Test image saved as 'fip_direct_test_image.png'")
	}

	// Send image to FIP
	fmt.Println("\n3. Sending image to FIP...")
	err = fipDirect.SendImage(testImage)
	if err != nil {
		log.Fatalf("Failed to send image to FIP: %v", err)
	}

	fmt.Println("✓ Image sent to FIP successfully!")

	// Test LED control
	fmt.Println("\n4. Testing LED control...")
	for i := 0; i < 6; i++ {
		fmt.Printf("Setting LED %d ON\n", i)
		err = fipDirect.SetLED(i, true)
		if err != nil {
			log.Printf("Failed to set LED %d: %v", i, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)

	for i := 0; i < 6; i++ {
		fmt.Printf("Setting LED %d OFF\n", i)
		err = fipDirect.SetLED(i, false)
		if err != nil {
			log.Printf("Failed to set LED %d: %v", i, err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Test button events
	fmt.Println("\n5. Testing button events (press buttons on FIP)...")
	eventChan, err := fipDirect.ReadButtonEvents()
	if err != nil {
		log.Printf("Warning: Could not read button events: %v", err)
	} else {
		// Listen for button events for 10 seconds
		timeout := time.After(10 * time.Second)
		fmt.Println("Listening for button events for 10 seconds...")

		for {
			select {
			case event := <-eventChan:
				fmt.Printf("Button event: Button=%d, Pressed=%v, Time=%v\n",
					event.Button, event.Pressed, event.Timestamp)
			case <-timeout:
				fmt.Println("Button event listening timeout")
				goto done
			}
		}
	done:
	}

	// Create and send different test images
	fmt.Println("\n6. Creating and sending different test images...")

	// Test image 1: Color bars
	colorBarsImage := createColorBarsImage()
	err = fipDirect.SendImage(colorBarsImage)
	if err != nil {
		log.Printf("Failed to send color bars image: %v", err)
	} else {
		fmt.Println("✓ Color bars image sent")
	}
	time.Sleep(2 * time.Second)

	// Test image 2: Gradient
	gradientImage := createGradientImage()
	err = fipDirect.SendImage(gradientImage)
	if err != nil {
		log.Printf("Failed to send gradient image: %v", err)
	} else {
		fmt.Println("✓ Gradient image sent")
	}
	time.Sleep(2 * time.Second)

	// Test image 3: Text
	textImage := createTextImage("FIP DIRECT TEST")
	err = fipDirect.SendImage(textImage)
	if err != nil {
		log.Printf("Failed to send text image: %v", err)
	} else {
		fmt.Println("✓ Text image sent")
	}
	time.Sleep(2 * time.Second)

	fmt.Println("\n✓ Direct FIP test completed successfully!")
	fmt.Println("\nNext Steps:")
	fmt.Println("- Integrate with your GUI application")
	fmt.Println("- Add flight simulator data integration")
	fmt.Println("- Create instrument-specific displays")
}

func createTestImage() image.Image {
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
	drawText(img, "FIP DIRECT", 160, 60, color.RGBA{255, 255, 255, 255})
	drawText(img, "320x240", 160, 80, color.RGBA{255, 255, 0, 255})
	drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

func createColorBarsImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Create color bars
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

func createGradientImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Create a gradient from top to bottom
	for y := 0; y < 240; y++ {
		ratio := float64(y) / 240.0
		r := uint8(255 * ratio)
		g := uint8(128 + 127*ratio)
		b := uint8(255 * (1.0 - ratio))

		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

func createTextImage(text string) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with black background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
		}
	}

	// Draw text in center
	drawText(img, text, 160, 120, color.RGBA{255, 255, 255, 255})

	return img
}

func drawText(img *image.RGBA, text string, x, y int, c color.Color) {
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
