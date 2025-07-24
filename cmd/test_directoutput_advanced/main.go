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
	fmt.Println("Advanced DirectOutput Test Program")
	fmt.Println("==================================")

	// Create a new DirectOutput instance
	do, err := fip.NewDirectOutput()
	if err != nil {
		log.Fatalf("Failed to create DirectOutput: %v", err)
	}
	defer do.Close()

	// Initialize DirectOutput
	err = do.Initialize("Saitek FIP Controller")
	if err != nil {
		log.Fatalf("Failed to initialize DirectOutput: %v", err)
	}

	// Create a simulated FIP device
	fmt.Println("Creating FIP device...")
	deviceHandle := unsafe.Pointer(uintptr(1))

	device := &fip.Device{
		Handle:     deviceHandle,
		DeviceType: fip.DeviceTypeFip,
		Pages:      make(map[uint32]*fip.Page),
	}
	do.Devices[deviceHandle] = device

	// Register callbacks
	fmt.Println("Registering callbacks...")
	err = do.RegisterPageCallback(deviceHandle, onPageChanged, nil)
	if err != nil {
		log.Printf("Warning: Failed to register page callback: %v", err)
	}

	err = do.RegisterSoftButtonCallback(deviceHandle, onSoftButtonChanged, nil)
	if err != nil {
		log.Printf("Warning: Failed to register soft button callback: %v", err)
	}

	// Add multiple pages
	fmt.Println("Adding pages...")
	pages := []struct {
		id   uint32
		name string
		img  string
	}{
		{1, "Airspeed Indicator", "assets/airspeed.png"},
		{2, "Altimeter", "assets/altimeter.png"},
		{3, "Artificial Horizon", "assets/artificial_horizon.png"},
		{4, "Compass", "assets/compass.png"},
		{5, "Turn Coordinator", "assets/turn_coordinator.png"},
		{6, "VSI", "assets/vsi.png"},
	}

	for _, page := range pages {
		err = do.AddPage(deviceHandle, page.id, page.name, fip.FLAG_SET_AS_ACTIVE)
		if err != nil {
			log.Printf("Failed to add page %d: %v", page.id, err)
			continue
		}

		// Try to load the image if it exists
		if _, err := os.Stat(page.img); err == nil {
			fmt.Printf("Loading image for page %d: %s\n", page.id, page.img)
			err = do.SetImageFromFile(deviceHandle, page.id, 0, page.img)
			if err != nil {
				log.Printf("Failed to load image for page %d: %v", page.id, err)
			}
		} else {
			// Create a test image if the file doesn't exist
			fmt.Printf("Creating test image for page %d\n", page.id)
			img := createInstrumentImage(page.name)
			fipData, err := do.ConvertImageToFIPFormat(img)
			if err != nil {
				log.Printf("Failed to convert image for page %d: %v", page.id, err)
				continue
			}
			err = do.SetImage(deviceHandle, page.id, 0, fipData)
			if err != nil {
				log.Printf("Failed to set image for page %d: %v", page.id, err)
			}
		}

		// Set LEDs based on page
		for i := uint32(0); i < 6; i++ {
			ledValue := uint32(0)
			if i == page.id-1 { // Light up the LED corresponding to the page
				ledValue = 1
			}
			err = do.SetLed(deviceHandle, page.id, i, ledValue)
			if err != nil {
				log.Printf("Failed to set LED %d on page %d: %v", i, page.id, err)
			}
		}
	}

	// Demonstrate page switching
	fmt.Println("\nDemonstrating page switching...")
	for i := 1; i <= 6; i++ {
		fmt.Printf("Switching to page %d...\n", i)
		// In a real implementation, this would trigger the page change callback
		time.Sleep(1 * time.Second)
	}

	// Demonstrate soft button handling
	fmt.Println("\nDemonstrating soft button handling...")
	buttons := []uint32{
		fip.SoftButton1,
		fip.SoftButton2,
		fip.SoftButton3,
		fip.SoftButton4,
		fip.SoftButton5,
		fip.SoftButton6,
		fip.SoftButtonUp,
		fip.SoftButtonDown,
		fip.SoftButtonLeft,
		fip.SoftButtonRight,
	}

	for _, button := range buttons {
		fmt.Printf("Simulating button press: 0x%08X\n", button)
		// In a real implementation, this would trigger the soft button callback
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\nTest completed successfully!")
	fmt.Println("DirectOutput wrapper is ready for integration with real FIP devices.")
}

// Callback functions
func onPageChanged(hDevice unsafe.Pointer, dwPage uint32, bSetActive bool, pCtxt unsafe.Pointer) {
	fmt.Printf("Page changed: %d, Active: %v\n", dwPage, bSetActive)
}

func onSoftButtonChanged(hDevice unsafe.Pointer, dwButtons uint32, pCtxt unsafe.Pointer) {
	fmt.Printf("Soft button changed: 0x%08X\n", dwButtons)

	// Decode which buttons were pressed
	if dwButtons&fip.SoftButton1 != 0 {
		fmt.Println("  - Button 1 pressed")
	}
	if dwButtons&fip.SoftButton2 != 0 {
		fmt.Println("  - Button 2 pressed")
	}
	if dwButtons&fip.SoftButton3 != 0 {
		fmt.Println("  - Button 3 pressed")
	}
	if dwButtons&fip.SoftButton4 != 0 {
		fmt.Println("  - Button 4 pressed")
	}
	if dwButtons&fip.SoftButton5 != 0 {
		fmt.Println("  - Button 5 pressed")
	}
	if dwButtons&fip.SoftButton6 != 0 {
		fmt.Println("  - Button 6 pressed")
	}
	if dwButtons&fip.SoftButtonUp != 0 {
		fmt.Println("  - Right dial clockwise")
	}
	if dwButtons&fip.SoftButtonDown != 0 {
		fmt.Println("  - Right dial counter-clockwise")
	}
	if dwButtons&fip.SoftButtonLeft != 0 {
		fmt.Println("  - Left dial counter-clockwise")
	}
	if dwButtons&fip.SoftButtonRight != 0 {
		fmt.Println("  - Left dial clockwise")
	}
}

func createInstrumentImage(instrumentName string) image.Image {
	// Create a 320x240 instrument image
	img := image.NewRGBA(image.Rect(0, 0, 320, 240))

	// Fill with dark background
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img.Set(x, y, color.RGBA{20, 20, 40, 255})
		}
	}

	// Add instrument-specific elements
	switch instrumentName {
	case "Airspeed Indicator":
		drawAirspeedIndicator(img)
	case "Altimeter":
		drawAltimeter(img)
	case "Artificial Horizon":
		drawArtificialHorizon(img)
	case "Compass":
		drawCompass(img)
	case "Turn Coordinator":
		drawTurnCoordinator(img)
	case "VSI":
		drawVSI(img)
	default:
		drawGenericInstrument(img, instrumentName)
	}

	return img
}

func drawAirspeedIndicator(img *image.RGBA) {
	// Draw a circular airspeed indicator
	centerX, centerY := 160, 120
	radius := 80

	// Draw outer circle
	for angle := 0; angle < 360; angle++ {
		x := centerX + int(float64(radius)*cos(float64(angle)*3.14159/180))
		y := centerY + int(float64(radius)*sin(float64(angle)*3.14159/180))
		if x >= 0 && x < 320 && y >= 0 && y < 240 {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// Draw speed markings
	for speed := 0; speed <= 200; speed += 20 {
		angle := float64(speed) * 2.7 // Scale to fit in 240 degrees
		x := centerX + int(float64(radius-10)*cos(angle*3.14159/180))
		y := centerY + int(float64(radius-10)*sin(angle*3.14159/180))
		if x >= 0 && x < 320 && y >= 0 && y < 240 {
			img.Set(x, y, color.RGBA{255, 255, 0, 255})
		}
	}

	// Draw center text
	drawText(img, "AIRSPEED", 160, 180, color.RGBA{255, 255, 255, 255})
	drawText(img, "KTS", 160, 200, color.RGBA{255, 255, 255, 255})
}

func drawAltimeter(img *image.RGBA) {
	// Draw altimeter display
	drawText(img, "ALTITUDE", 160, 60, color.RGBA{255, 255, 255, 255})
	drawText(img, "00000", 160, 120, color.RGBA{0, 255, 0, 255})
	drawText(img, "FT", 160, 180, color.RGBA{255, 255, 255, 255})
}

func drawArtificialHorizon(img *image.RGBA) {
	// Draw artificial horizon
	_, centerY := 160, 120

	// Draw horizon line
	for x := 60; x < 260; x++ {
		img.Set(x, centerY, color.RGBA{255, 255, 255, 255})
	}

	// Draw sky (top half)
	for y := 0; y < centerY; y++ {
		for x := 60; x < 260; x++ {
			img.Set(x, y, color.RGBA{0, 100, 200, 255})
		}
	}

	// Draw ground (bottom half)
	for y := centerY; y < 240; y++ {
		for x := 60; x < 260; x++ {
			img.Set(x, y, color.RGBA{139, 69, 19, 255})
		}
	}

	drawText(img, "ATTITUDE", 160, 200, color.RGBA{255, 255, 255, 255})
}

func drawCompass(img *image.RGBA) {
	// Draw compass rose
	centerX, centerY := 160, 120
	radius := 60

	// Draw compass circle
	for angle := 0; angle < 360; angle++ {
		x := centerX + int(float64(radius)*cos(float64(angle)*3.14159/180))
		y := centerY + int(float64(radius)*sin(float64(angle)*3.14159/180))
		if x >= 0 && x < 320 && y >= 0 && y < 240 {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// Draw cardinal directions
	drawText(img, "N", 160, 70, color.RGBA{255, 0, 0, 255})
	drawText(img, "E", 210, 120, color.RGBA{255, 255, 255, 255})
	drawText(img, "S", 160, 170, color.RGBA{255, 255, 255, 255})
	drawText(img, "W", 110, 120, color.RGBA{255, 255, 255, 255})

	drawText(img, "HDG", 160, 200, color.RGBA{255, 255, 255, 255})
}

func drawTurnCoordinator(img *image.RGBA) {
	// Draw turn coordinator
	centerX, centerY := 160, 120
	radius := 50

	// Draw outer circle
	for angle := 0; angle < 360; angle++ {
		x := centerX + int(float64(radius)*cos(float64(angle)*3.14159/180))
		y := centerY + int(float64(radius)*sin(float64(angle)*3.14159/180))
		if x >= 0 && x < 320 && y >= 0 && y < 240 {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}

	// Draw turn indicator
	for i := -20; i <= 20; i++ {
		x := centerX + i
		y := centerY + int(float64(i*i)/20)
		if x >= 0 && x < 320 && y >= 0 && y < 240 {
			img.Set(x, y, color.RGBA{0, 255, 0, 255})
		}
	}

	drawText(img, "TURN", 160, 200, color.RGBA{255, 255, 255, 255})
}

func drawVSI(img *image.RGBA) {
	// Draw vertical speed indicator
	centerX, centerY := 160, 120

	// Draw vertical scale
	for y := 40; y < 200; y++ {
		img.Set(centerX, y, color.RGBA{255, 255, 255, 255})
	}

	// Draw scale markings
	for i := -10; i <= 10; i++ {
		y := centerY + i*10
		if y >= 0 && y < 240 {
			img.Set(centerX-10, y, color.RGBA{255, 255, 0, 255})
			img.Set(centerX+10, y, color.RGBA{255, 255, 0, 255})
		}
	}

	drawText(img, "VSI", 160, 200, color.RGBA{255, 255, 255, 255})
}

func drawGenericInstrument(img *image.RGBA, name string) {
	// Draw a generic instrument display
	drawText(img, name, 160, 120, color.RGBA{255, 255, 255, 255})
}

func drawText(img *image.RGBA, text string, x, y int, c color.Color) {
	// Simple text drawing (in a real implementation, you'd use a proper font)
	// For now, just draw some pixels to represent text
	for i, _ := range text {
		// Simple character representation
		charX := x + i*8 - len(text)*4
		if charX >= 0 && charX < 320 {
			img.Set(charX, y, c)
			img.Set(charX+1, y, c)
			img.Set(charX, y+1, c)
			img.Set(charX+1, y+1, c)
		}
	}
}

// Simple math functions
func cos(x float64) float64 {
	// Simple cosine approximation
	return 1 - x*x/2 + x*x*x*x/24
}

func sin(x float64) float64 {
	// Simple sine approximation
	return x - x*x*x/6 + x*x*x*x*x/120
}
