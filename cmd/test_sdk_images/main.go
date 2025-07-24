package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"

	"saitek-controller/internal/usb"
)

func main() {
	fmt.Println("DirectOutput SDK Images Test")
	fmt.Println("============================")

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

	// Test with actual SDK sample images
	fmt.Println("\n3. Testing with DirectOutput SDK sample images...")

	sampleImages := []string{
		"DirectOutput/Fip1.jpg",
		"DirectOutput/Fip2.jpg",
		"DirectOutput/Fip3.jpg",
		"DirectOutput/Fip4.jpg",
		"DirectOutput/Fip5.jpg",
	}

	for i, imagePath := range sampleImages {
		fmt.Printf("\n   Testing sample image %d: %s\n", i+1, filepath.Base(imagePath))

		// Load the sample image
		img, err := loadImage(imagePath)
		if err != nil {
			fmt.Printf("   ✗ Failed to load image: %v\n", err)
			continue
		}

		// Convert to FIP format
		fipData := convertImageToFIPFormat(img)
		fmt.Printf("   ✓ Image loaded: %dx%d -> %d bytes FIP data\n",
			img.Bounds().Dx(), img.Bounds().Dy(), len(fipData))

		// Try to send the image using different protocols
		err = sendImageToFIP(device, fipData, i+1)
		if err != nil {
			fmt.Printf("   ✗ Failed to send image: %v\n", err)
		} else {
			fmt.Printf("   ✓ Image sent successfully\n")
		}
	}

	fmt.Println("\n✅ SDK images test completed!")
	fmt.Println("\nIf you don't see images on your FIP:")
	fmt.Println("1. The FIP may need the DirectOutput service running")
	fmt.Println("2. We may need to analyze the DirectOutput protocol more deeply")
	fmt.Println("3. The FIP might require specific initialization")
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

func sendImageToFIP(device *usb.Device, imageData []byte, imageIndex int) error {
	// Try different protocols to send the image

	// Protocol 1: DirectOutput-style packet
	packet1 := createDirectOutputPacket(imageData, imageIndex)
	err := device.SendControlMessage(0x21, 0x09, 0x0200, 0, packet1)
	if err == nil {
		return nil
	}

	// Protocol 2: Simple image data
	err = device.SendControlMessage(0x21, 0x09, 0x0200, 0, imageData[:1024])
	if err == nil {
		return nil
	}

	// Protocol 3: HID-style packet
	packet3 := createHIDPacket(imageData, imageIndex)
	err = device.SendControlMessage(0x21, 0x09, 0x0200, 0, packet3)
	if err == nil {
		return nil
	}

	return fmt.Errorf("all protocols failed")
}

func createDirectOutputPacket(imageData []byte, imageIndex int) []byte {
	// Create a packet that mimics the DirectOutput protocol
	packet := make([]byte, 0, len(imageData)+64)

	// DirectOutput header
	header := []byte{
		// Magic bytes (based on DirectOutput analysis)
		0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x4F, 0x55, // "DIRECTOU"
		0x54, 0x50, 0x55, 0x54, 0x20, 0x20, 0x20, 0x20, // "TPUT    "

		// Device type (FIP GUID)
		0x3E, 0x08, 0x3C, 0xD8, 0x6A, 0x37, 0x4A, 0x58, // FIP GUID
		0x80, 0xA8, 0x3D, 0x6A, 0x2C, 0x07, 0x51, 0x3E, // (continued)

		// Command: SetImage
		0x53, 0x45, 0x54, 0x49, 0x4D, 0x41, 0x47, 0x45, // "SETIMAGE"

		// Page and Index
		0x00, 0x00, 0x00, 0x00, // Page 0
		byte(imageIndex), 0x00, 0x00, 0x00, // Index

		// Image size
		0x00, 0x00, 0x00, 0x00, // Size (will be filled)

		// Reserved
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Reserved
		0x00, 0x00, 0x00, 0x00, // Reserved
	}

	// Set image size
	size := uint32(len(imageData))
	header[40] = byte(size & 0xFF)
	header[41] = byte((size >> 8) & 0xFF)
	header[42] = byte((size >> 16) & 0xFF)
	header[43] = byte((size >> 24) & 0xFF)

	packet = append(packet, header...)
	packet = append(packet, imageData...)

	return packet
}

func createHIDPacket(imageData []byte, imageIndex int) []byte {
	// Create a HID-style packet
	packet := make([]byte, 0, len(imageData)+16)

	// HID header
	header := []byte{
		0x01,             // Report ID
		byte(imageIndex), // Image index
		0x00, 0x00,       // Reserved

		// Image size
		byte(len(imageData) & 0xFF),
		byte((len(imageData) >> 8) & 0xFF),
		byte((len(imageData) >> 16) & 0xFF),
		byte((len(imageData) >> 24) & 0xFF),

		// Reserved
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	packet = append(packet, header...)
	packet = append(packet, imageData...)

	return packet
}
