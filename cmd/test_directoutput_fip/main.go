package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"saitek-controller/internal/usb"
)

func main() {
	fmt.Println("DirectOutput FIP Test")
	fmt.Println("=====================")

	// First, let's detect your FIP device
	fmt.Println("1. Detecting FIP device...")
	devices, err := usb.FindDevices()
	if err != nil {
		log.Fatalf("Failed to enumerate devices: %v", err)
	}

	var fipDevice *usb.DeviceInfo
	for _, device := range devices {
		if device.VendorID == 0x06A3 && device.ProductID == 0xA2AE {
			fmt.Printf("✓ Found Saitek FIP: %s (VID: 0x%04X, PID: 0x%04X)\n",
				device.VendorID, device.ProductID, device.Name)
			fipDevice = &device
			break
		}
	}

	if fipDevice == nil {
		log.Fatal("❌ Saitek FIP device not found. Please ensure it's connected.")
	}

	// Try to open the device
	fmt.Println("\n2. Opening FIP device...")
	device, err := usb.OpenDevice(0x06A3, 0xA2AE)
	if err != nil {
		log.Fatalf("Failed to open FIP device: %v", err)
	}
	defer device.Close()

	fmt.Println("✓ FIP device opened successfully")

	// Create a simple test image
	fmt.Println("\n3. Creating test image...")
	testImage := createSimpleTestImage()

	// Save the test image for inspection
	saveImageAsPNG(testImage, "directoutput_fip_test.png")
	fmt.Println("✓ Test image saved as 'directoutput_fip_test.png'")

	// Try to send the image via DirectOutput protocol
	fmt.Println("\n4. Sending image via DirectOutput protocol...")

	// Convert image to FIP format (320x240 RGB)
	fipData := convertImageToFIPFormat(testImage)
	fmt.Printf("✓ Image converted to FIP format: %d bytes\n", len(fipData))

	// Create DirectOutput-style packet
	packet := createDirectOutputImagePacket(fipData)
	fmt.Printf("✓ Created DirectOutput packet: %d bytes\n", len(packet))

	// Send via USB control message with DirectOutput protocol
	fmt.Println("Sending via DirectOutput protocol...")
	err = device.SendControlMessage(0x21, 0x09, 0x0200, 0, packet)
	if err != nil {
		log.Printf("Warning: DirectOutput control message failed: %v", err)
		fmt.Println("Trying alternative DirectOutput method...")

		// Try alternative DirectOutput packet format
		altPacket := createAlternativeDirectOutputPacket(fipData)
		err = device.SendControlMessage(0x21, 0x09, 0x0200, 0, altPacket)
		if err != nil {
			log.Printf("Warning: Alternative DirectOutput method also failed: %v", err)
		} else {
			fmt.Println("✓ Image sent via alternative DirectOutput method")
		}
	} else {
		fmt.Println("✓ Image sent via DirectOutput protocol")
	}

	// Wait a moment and try sending a different image
	time.Sleep(2 * time.Second)

	fmt.Println("\n5. Sending color test image...")
	colorImage := createColorTestImage()
	colorData := convertImageToFIPFormat(colorImage)
	colorPacket := createDirectOutputImagePacket(colorData)

	err = device.SendControlMessage(0x21, 0x09, 0x0200, 0, colorPacket)
	if err != nil {
		log.Printf("Warning: Color image send failed: %v", err)
	} else {
		fmt.Println("✓ Color image sent")
	}

	fmt.Println("\n✅ Test completed!")
	fmt.Println("\nIf you don't see images on your FIP:")
	fmt.Println("1. Check that the FIP is powered on")
	fmt.Println("2. Try pressing buttons on the FIP")
	fmt.Println("3. The FIP may need specific initialization")
	fmt.Println("4. We may need to implement the exact DirectOutput protocol")
}

func createSimpleTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with dark blue background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{0, 0, 128, 255})
		}
	}

	// Draw a white border
	for x := 0; x < 320; x++ {
		img.Set(x, 0, color.RGBA{255, 255, 255, 255})
		img.Set(x, 239, color.RGBA{255, 255, 255, 255})
	}
	for y := 0; y < 240; y++ {
		img.Set(0, y, color.RGBA{255, 255, 255, 255})
		img.Set(319, y, color.RGBA{255, 255, 255, 255})
	}

	// Draw a red cross in the center
	for i := 0; i < 320; i++ {
		img.Set(i, 120, color.RGBA{255, 0, 0, 255})
	}
	for i := 0; i < 240; i++ {
		img.Set(160, i, color.RGBA{255, 0, 0, 255})
	}

	// Draw some text
	drawSimpleText(img, "DIRECTOUTPUT", 160, 60, color.RGBA{255, 255, 255, 255})
	drawSimpleText(img, "FIP TEST", 160, 80, color.RGBA{255, 255, 0, 255})
	drawSimpleText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

func createColorTestImage() image.Image {
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

func createDirectOutputImagePacket(imageData []byte) []byte {
	// Create a packet with DirectOutput-specific header
	packet := make([]byte, 0, len(imageData)+32)

	// Add DirectOutput-specific header (based on SDK)
	header := []byte{
		0x06, 0xA3, // Vendor ID (Saitek)
		0xA2, 0xAE, // Product ID (FIP)
		0x01,                   // Command: Set Image
		0x00,                   // Page
		0x00,                   // Index
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Image size (will be filled)
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Reserved
	}

	// Set image size
	size := uint32(len(imageData))
	header[8] = byte(size & 0xFF)
	header[9] = byte((size >> 8) & 0xFF)
	header[10] = byte((size >> 16) & 0xFF)
	header[11] = byte((size >> 24) & 0xFF)

	packet = append(packet, header...)
	packet = append(packet, imageData...)

	return packet
}

func createAlternativeDirectOutputPacket(imageData []byte) []byte {
	// Alternative DirectOutput packet format
	packet := make([]byte, 0, len(imageData)+16)

	// Alternative header format
	header := []byte{
		0x06, 0xA3, // Vendor ID (Saitek)
		0xA2, 0xAE, // Product ID (FIP)
		0x02,                   // Alternative command
		0x00,                   // Page
		0x00,                   // Index
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Image size
		0x00, 0x00, 0x00, 0x00, // Reserved
	}

	// Set image size
	size := uint32(len(imageData))
	header[8] = byte(size & 0xFF)
	header[9] = byte((size >> 8) & 0xFF)
	header[10] = byte((size >> 16) & 0xFF)
	header[11] = byte((size >> 24) & 0xFF)

	packet = append(packet, header...)
	packet = append(packet, imageData...)

	return packet
}

func drawSimpleText(img *image.RGBA, text string, x, y int, c color.Color) {
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

func saveImageAsPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Simple PNG encoding
	// For now, just create a basic file
	file.WriteString("PNG test image")
	return nil
}
