package main

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"saitek-controller/internal/usb"
)

func main() {
	fmt.Println("FIP Protocol Reverse Engineering Test")
	fmt.Println("=====================================")

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
	fipData := convertImageToFIPFormat(testImage)

	// Test different protocol approaches
	fmt.Println("\n4. Testing different FIP protocol approaches...")

	// Test 1: DirectOutput-style initialization
	testDirectOutputInit(device)

	// Test 2: HID-style communication
	testHIDStyleCommunication(device, fipData)

	// Test 3: Custom protocol based on DirectOutput analysis
	testCustomDirectOutputProtocol(device, fipData)

	// Test 4: Raw USB bulk transfer
	testRawUSBBulkTransfer(device, fipData)

	fmt.Println("\n✅ Protocol reverse engineering test completed!")
	fmt.Println("\nIf you don't see images on your FIP:")
	fmt.Println("1. The FIP may need specific initialization")
	fmt.Println("2. We may need to analyze the DirectOutput service more deeply")
	fmt.Println("3. The FIP might require a specific driver or service")
}

func testDirectOutputInit(device *usb.Device) {
	fmt.Println("\n   Testing DirectOutput-style initialization...")

	// Try different initialization sequences
	initSequences := [][]byte{
		// DirectOutput initialization
		{0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x4F, 0x55, 0x54, 0x50, 0x55, 0x54}, // "DIRECTOUTPUT"
		{0x49, 0x4E, 0x49, 0x54, 0x49, 0x41, 0x4C, 0x49, 0x5A, 0x45},             // "INITIALIZE"
		{0x46, 0x49, 0x50, 0x20, 0x49, 0x4E, 0x49, 0x54},                         // "FIP INIT"
		{0x01, 0x00, 0x00, 0x00},                                                 // Simple init
		{0x02, 0x00, 0x00, 0x00},                                                 // Alternative init
	}

	for i, seq := range initSequences {
		fmt.Printf("   Testing init sequence %d: %v\n", i+1, seq)
		err := device.SendControlMessage(0x21, 0x09, 0x0200, 0, seq)
		if err != nil {
			fmt.Printf("   ✗ Init sequence %d failed: %v\n", i+1, err)
		} else {
			fmt.Printf("   ✓ Init sequence %d sent successfully\n", i+1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testHIDStyleCommunication(device *usb.Device, imageData []byte) {
	fmt.Println("\n   Testing HID-style communication...")

	// Try different HID report formats
	hidReports := [][]byte{
		// HID report with image data
		append([]byte{0x01, 0x00, 0x00, 0x00}, imageData[:64]...),
		// HID report with command
		{0x02, 0x53, 0x45, 0x54, 0x49, 0x4D, 0x41, 0x47, 0x45}, // "SETIMAGE"
		// HID report with page info
		{0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Page 0
	}

	for i, report := range hidReports {
		fmt.Printf("   Testing HID report %d: %d bytes\n", i+1, len(report))
		err := device.SendControlMessage(0x21, 0x09, 0x0200, 0, report)
		if err != nil {
			fmt.Printf("   ✗ HID report %d failed: %v\n", i+1, err)
		} else {
			fmt.Printf("   ✓ HID report %d sent successfully\n", i+1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testCustomDirectOutputProtocol(device *usb.Device, imageData []byte) {
	fmt.Println("\n   Testing custom DirectOutput protocol...")

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
		0x00, 0x00, 0x00, 0x00, // Index 0

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
	binary.LittleEndian.PutUint32(header[40:44], size)

	packet = append(packet, header...)
	packet = append(packet, imageData...)

	fmt.Printf("   Testing custom DirectOutput protocol: %d bytes\n", len(packet))
	err := device.SendControlMessage(0x21, 0x09, 0x0200, 0, packet)
	if err != nil {
		fmt.Printf("   ✗ Custom DirectOutput protocol failed: %v\n", err)
	} else {
		fmt.Printf("   ✓ Custom DirectOutput protocol sent successfully\n")
	}
}

func testRawUSBBulkTransfer(device *usb.Device, imageData []byte) {
	fmt.Println("\n   Testing raw USB bulk transfer...")

	// Try sending image data directly via bulk transfer
	fmt.Printf("   Testing raw USB bulk transfer: %d bytes\n", len(imageData))

	// Try different bulk endpoints
	endpoints := []int{0x01, 0x02, 0x03}

	for _, endpoint := range endpoints {
		fmt.Printf("   Testing bulk endpoint 0x%02X\n", endpoint)
		// Note: We can't directly call bulk transfer, but we can try control messages
		// that might trigger bulk transfers
		err := device.SendControlMessage(0x21, 0x09, uint16(endpoint), 0, imageData[:1024])
		if err != nil {
			fmt.Printf("   ✗ Bulk endpoint 0x%02X failed: %v\n", endpoint, err)
		} else {
			fmt.Printf("   ✓ Bulk endpoint 0x%02X sent successfully\n", endpoint)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func createSimpleTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with a simple pattern
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			// Create a simple gradient
			r := uint8((x * 255) / 320)
			g := uint8((y * 255) / 240)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
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
