package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"unsafe"

	"saitek-controller/internal/fip"
	"saitek-controller/internal/usb"
)

func main() {
	fmt.Println("Real FIP Device Test")
	fmt.Println("====================")

	// First, let's try to detect the FIP using our existing USB methods
	fmt.Println("1. Detecting FIP device via USB...")

		// Look for Saitek FIP devices
	devices, err := usb.FindDevices()
	if err != nil {
		log.Printf("Warning: Failed to enumerate USB devices: %v", err)
	} else {
		fmt.Printf("Found %d USB devices\n", len(devices))
		
		for i, device := range devices {
			fmt.Printf("Device %d: VID=0x%04X, PID=0x%04X, Name=%s, Path=%s\n", 
				i, device.VendorID, device.ProductID, device.Name, device.Path)
			
			// Check if this is a Saitek FIP
			if device.VendorID == 0x06A3 && device.ProductID == 0xA2AE {
				fmt.Printf("✓ Found Saitek FIP device!\n")
				fmt.Printf("  Name: %s\n", device.Name)
				fmt.Printf("  Path: %s\n", device.Path)
				
				// Try to open the device
				fmt.Println("2. Attempting to open FIP device...")
				handle, err := usb.OpenDevice(0x06A3, 0xA2AE)
				if err != nil {
					log.Printf("Failed to open FIP device: %v", err)
				} else {
					fmt.Println("✓ Successfully opened FIP device!")
					defer handle.Close()
					
					// Test basic USB communication
					testUSBCommunication(handle)
				}
			}
		}
	}

	// Now test DirectOutput integration
	fmt.Println("\n3. Testing DirectOutput integration...")
	testDirectOutput()
}

func testUSBCommunication(handle *usb.Device) {
	fmt.Println("   Testing USB communication...")
	
	// Show device info
	fmt.Printf("   Device: VID=0x%04X, PID=0x%04X, Name=%s\n", 
		handle.VendorID, handle.ProductID, handle.Name)
	
	// Test sending a control message (FIP display update)
	fmt.Println("   Testing FIP display update...")
	
	// Create a test image data (22 bytes for radio panel, but FIP might be different)
	testData := make([]byte, 22)
	for i := range testData {
		testData[i] = byte(i)
	}
	
	err := handle.SendControlMessage(0x21, 0x09, 0x0200, 0, testData)
	if err != nil {
		log.Printf("Failed to send control message: %v", err)
	} else {
		fmt.Println("   ✓ Successfully sent control message")
	}
	
	// Try to read some data from the device
	fmt.Println("   Attempting to read data from device...")
	
	data, err := handle.ReadBulkData(0x81, 64) // Common HID endpoint
	if err != nil {
		log.Printf("Failed to read bulk data: %v", err)
	} else {
		fmt.Printf("   ✓ Received %d bytes: %x\n", len(data), data[:min(len(data), 16)])
	}
	
	// Check if device is still connected
	if handle.IsConnected() {
		fmt.Println("   ✓ Device is still connected")
	} else {
		fmt.Println("   ⚠ Device connection lost")
	}
}

func testDirectOutput() {
	// Create DirectOutput instance
	do, err := fip.NewDirectOutput()
	if err != nil {
		log.Printf("Failed to create DirectOutput: %v", err)
		return
	}
	defer do.Close()

	// Initialize DirectOutput
	err = do.Initialize("Real FIP Test")
	if err != nil {
		log.Printf("Failed to initialize DirectOutput: %v", err)
		return
	}

	fmt.Println("   ✓ DirectOutput initialized")

	// Create a test image
	fmt.Println("   Creating test image...")
	img := createTestImage()

	// Convert to FIP format
	fipData, err := do.ConvertImageToFIPFormat(img)
	if err != nil {
		log.Printf("Failed to convert image: %v", err)
		return
	}

	fmt.Printf("   ✓ Image converted to FIP format: %d bytes\n", len(fipData))

	// Save test image for inspection
	err = do.SaveImageAsPNG(img, "real_fip_test_image.png")
	if err != nil {
		log.Printf("Failed to save test image: %v", err)
	} else {
		fmt.Println("   ✓ Test image saved as 'real_fip_test_image.png'")
	}

	// Try to create a simulated device (since we don't have real DirectOutput DLL on macOS)
	fmt.Println("   Creating simulated DirectOutput device...")
	deviceHandle := unsafe.Pointer(uintptr(1))

	device := &fip.Device{
		Handle:     deviceHandle,
		DeviceType: fip.DeviceTypeFip,
		Pages:      make(map[uint32]*fip.Page),
	}
	do.Devices[deviceHandle] = device

	// Add a test page
	err = do.AddPage(deviceHandle, 1, "Real FIP Test", fip.FLAG_SET_AS_ACTIVE)
	if err != nil {
		log.Printf("Failed to add page: %v", err)
		return
	}

	// Set the test image
	err = do.SetImage(deviceHandle, 1, 0, fipData)
	if err != nil {
		log.Printf("Failed to set image: %v", err)
		return
	}

	// Set some LEDs
	for i := uint32(0); i < 6; i++ {
		err = do.SetLed(deviceHandle, 1, i, 1)
		if err != nil {
			log.Printf("Failed to set LED %d: %v", i, err)
		}
	}

	fmt.Println("   ✓ DirectOutput test completed successfully")
	fmt.Println("\n4. Next Steps:")
	fmt.Println("   - On Windows: Install DirectOutput SDK and test with real DLL")
	fmt.Println("   - On macOS: Use our wrapper for development/testing")
	fmt.Println("   - Integrate with your GUI application")
}

func createTestImage() image.Image {
	// Create a 320x240 test image with clear patterns
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
	drawText(img, "FIP TEST", 160, 60, color.RGBA{255, 255, 255, 255})
	drawText(img, "320x240", 160, 80, color.RGBA{255, 255, 0, 255})
	drawText(img, "READY", 160, 180, color.RGBA{0, 255, 0, 255})

	return img
}

func drawText(img *image.RGBA, text string, x, y int, c color.Color) {
	// Simple text drawing
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
