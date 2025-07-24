package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"unsafe"

	"saitek-controller/internal/fip"
)

func main() {
	fmt.Println("DirectOutput Test Program")
	fmt.Println("========================")

	// Create a new DirectOutput instance
	do, err := fip.NewDirectOutput()
	if err != nil {
		log.Fatalf("Failed to create DirectOutput: %v", err)
	}
	defer do.Close()

	// Initialize DirectOutput
	err = do.Initialize("Saitek Controller Test")
	if err != nil {
		log.Fatalf("Failed to initialize DirectOutput: %v", err)
	}

	// Create a test image
	fmt.Println("Creating test image...")
	img := createTestImage()

	// Convert to FIP format
	fmt.Println("Converting to FIP format...")
	fipData, err := do.ConvertImageToFIPFormat(img)
	if err != nil {
		log.Fatalf("Failed to convert image: %v", err)
	}

	// Save the FIP image for inspection
	outputPath := "test_fip_image.png"
	fmt.Printf("Saving FIP image to: %s\n", outputPath)
	err = do.SaveImageAsPNG(img, outputPath)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

		// Create a simulated device
	fmt.Println("Creating simulated FIP device...")
	deviceHandle := unsafe.Pointer(uintptr(1)) // Simulated handle
	
	// Create a device in the DirectOutput instance
	device := &fip.Device{
		Handle:     deviceHandle,
		DeviceType: fip.DeviceTypeFip,
		Pages:      make(map[uint32]*fip.Page),
	}
	do.Devices[deviceHandle] = device
	
	// Add a page to the device
	err = do.AddPage(deviceHandle, 1, "Test Page", fip.FLAG_SET_AS_ACTIVE)
	if err != nil {
		log.Fatalf("Failed to add page: %v", err)
	}

	// Set the image on the device
	err = do.SetImage(deviceHandle, 1, 0, fipData)
	if err != nil {
		log.Fatalf("Failed to set image: %v", err)
	}

	// Set some LEDs
	fmt.Println("Setting LEDs...")
	for i := uint32(0); i < 6; i++ {
		err = do.SetLed(deviceHandle, 1, i, 1) // Turn on all LEDs
		if err != nil {
			log.Printf("Failed to set LED %d: %v", i, err)
		}
	}

	fmt.Println("Test completed successfully!")
	fmt.Printf("FIP image data size: %d bytes\n", len(fipData))
	fmt.Printf("Expected size: %d bytes (320x240x3)\n", 320*240*3)
}

func createTestImage() image.Image {
	// Create a 320x240 test image
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with a gradient
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			r := uint8((x * 255) / 320)
			g := uint8((y * 255) / 240)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// Add some text-like patterns
	for y := 50; y < 190; y += 20 {
		for x := 50; x < 270; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	return img
}
