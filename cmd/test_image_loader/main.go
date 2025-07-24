package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"saitek-controller/internal/fip"
)

func main() {
	fmt.Println("FIP Image Loader Test")
	fmt.Println("=====================")

	// Create image loader
	loader := fip.NewImageLoader()
	
	// Create SDK for testing
	sdk := NewDirectOutputSDK()
	defer sdk.Close()

	// Initialize SDK
	err := sdk.Initialize("Image Loader Test")
	if err != nil {
		log.Fatalf("Failed to initialize SDK: %v", err)
	}

	// Create device handle
	deviceHandle := unsafe.Pointer(uintptr(0x12345678))

	// Add page
	err = sdk.AddPage(deviceHandle, 1, "Image Loader Test Page", 0x00000001)
	if err != nil {
		log.Fatalf("Failed to add page: %v", err)
	}

	// Test different resize modes
	resizeModes := []struct {
		name fip.ResizeMode
		desc string
	}{
		{fip.ResizeModeStretch, "Stretch (may distort)"},
		{fip.ResizeModeFit, "Fit (maintain aspect ratio)"},
		{fip.ResizeModeCrop, "Crop (maintain aspect ratio)"},
		{fip.ResizeModeCenter, "Center (pad with background)"},
	}

	// Create test images of different sizes
	fmt.Println("1. Creating test images of different sizes...")
	testImages := createTestImages()
	
	for i, img := range testImages {
		filename := fmt.Sprintf("test_image_%d.png", i+1)
		err := loader.SaveImageAsPNG(img, filename)
		if err != nil {
			log.Printf("Warning: Failed to save test image %d: %v", i+1, err)
		} else {
			fmt.Printf("   ✓ Created %s (%dx%d)\n", filename, img.Bounds().Dx(), img.Bounds().Dy())
		}
	}

	// Test loading and processing each image with different resize modes
	fmt.Println("\n2. Testing image loading with different resize modes...")
	
	for _, mode := range resizeModes {
		fmt.Printf("\n   Testing %s mode:\n", mode.desc)
		loader.SetResizeMode(mode.name)
		
		for i := range testImages {
			filename := fmt.Sprintf("test_image_%d.png", i+1)
			
			// Load and process the image
			img, err := loader.LoadImageFromFile(filename)
			if err != nil {
				log.Printf("Warning: Failed to load %s: %v", filename, err)
				continue
			}

			// Get image info
			info := loader.GetImageInfo(img)
			fmt.Printf("     %s: %dx%d -> %dx%d (resize: %v)\n", 
				filename, info.Width, info.Height, info.FIPWidth, info.FIPHeight, info.NeedsResize)

			// Convert to FIP format
			fipData, err := loader.ConvertImageToFIPFormat(img)
			if err != nil {
				log.Printf("Warning: Failed to convert %s: %v", filename, err)
				continue
			}

			// Send to FIP (simulated)
			err = sdk.SetImage(deviceHandle, 1, 0, fipData)
			if err != nil {
				log.Printf("Warning: Failed to send %s: %v", filename, err)
			} else {
				fmt.Printf("     ✓ Sent %s to FIP\n", filename)
			}

			// Save processed image
			outputFilename := fmt.Sprintf("processed_%s_%s", mode.desc, filename)
			err = loader.SaveImageAsPNG(img, outputFilename)
			if err != nil {
				log.Printf("Warning: Failed to save processed %s: %v", outputFilename, err)
			} else {
				fmt.Printf("     ✓ Saved %s\n", outputFilename)
			}
		}
	}

	// Test format validation
	fmt.Println("\n3. Testing format validation...")
	testFormats := []string{
		"test.png",
		"test.jpg", 
		"test.jpeg",
		"test.gif",
		"test.bmp",
		"test.txt",
	}

	for _, format := range testFormats {
		supported := loader.IsSupportedFormat(format)
		fmt.Printf("   %s: %v\n", format, supported)
	}

	// Test image validation
	fmt.Println("\n4. Testing image size validation...")
	
	// Create images of different sizes
	sizes := []struct {
		width, height int
		valid         bool
	}{
		{320, 240, true},   // Correct size
		{640, 480, false},  // Too large
		{160, 120, false},  // Too small
		{320, 120, false},  // Wrong aspect ratio
		{120, 240, false},  // Wrong aspect ratio
	}

	for i, size := range sizes {
		img := image.NewRGBA(image.Rect(0, 0, size.width, size.height))
		
		// Fill with test pattern
		for y := 0; y < size.height; y++ {
			for x := 0; x < size.width; x++ {
				img.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
			}
		}

		err := loader.ValidateImageSize(img)
		valid := err == nil
		
		fmt.Printf("   %dx%d: %v (expected: %v)\n", size.width, size.height, valid, size.valid)
		
		if valid != size.valid {
			fmt.Printf("     Warning: Validation mismatch!\n")
		}
	}

	// Test direct loading and conversion
	fmt.Println("\n5. Testing direct loading and conversion...")
	
	for i := range testImages {
		filename := fmt.Sprintf("test_image_%d.png", i+1)
		
		// Load and convert directly to FIP format
		fipData, err := loader.LoadAndConvertToFIP(filename)
		if err != nil {
			log.Printf("Warning: Failed to load and convert %s: %v", filename, err)
			continue
		}

		fmt.Printf("   ✓ Loaded and converted %s (%d bytes)\n", filename, len(fipData))
		
		// Send to FIP
		err = sdk.SetImage(deviceHandle, 1, 0, fipData)
		if err != nil {
			log.Printf("Warning: Failed to send %s: %v", filename, err)
		} else {
			fmt.Printf("     ✓ Sent to FIP\n")
		}
	}

	// Test JPEG saving
	fmt.Println("\n6. Testing JPEG saving...")
	
	for i, img := range testImages {
		// Process image for FIP
		processedImg, err := loader.ProcessImageForFIP(img)
		if err != nil {
			log.Printf("Warning: Failed to process image %d: %v", i+1, err)
			continue
		}

		// Save as JPEG with different qualities
		qualities := []int{50, 75, 90}
		for _, quality := range qualities {
			loader.SetQuality(quality)
			filename := fmt.Sprintf("test_image_%d_q%d.jpg", i+1, quality)
			
			err := loader.SaveImageAsJPEG(processedImg, filename)
			if err != nil {
				log.Printf("Warning: Failed to save %s: %v", filename, err)
			} else {
				fmt.Printf("   ✓ Saved %s (quality: %d)\n", filename, quality)
			}
		}
	}

	// Clean up
	fmt.Println("\n7. Cleaning up...")
	err = sdk.RemovePage(deviceHandle, 1)
	if err != nil {
		log.Printf("Warning: Failed to remove page: %v", err)
	}

	fmt.Println("\n✅ FIP Image Loader Test Completed!")
	fmt.Println("\nGenerated files:")
	
	// List generated files
	files, err := filepath.Glob("*.png")
	if err == nil {
		for _, file := range files {
			if strings.HasPrefix(file, "test_image_") || strings.HasPrefix(file, "processed_") {
				fmt.Printf("  - %s\n", file)
			}
		}
	}
	
	files, err = filepath.Glob("*.jpg")
	if err == nil {
		for _, file := range files {
			fmt.Printf("  - %s\n", file)
		}
	}
	
	fmt.Println("\nThis demonstrates comprehensive image loading capabilities!")
}

// DirectOutput SDK simulation (same as before)
type DirectOutputSDK struct {
	useRealSDK bool
	initialized bool
}

func NewDirectOutputSDK() *DirectOutputSDK {
	sdk := &DirectOutputSDK{}
	
	if runtime.GOOS == "windows" {
		sdk.useRealSDK = false // For demo purposes
	} else {
		sdk.useRealSDK = false
	}
	
	return sdk
}

func (sdk *DirectOutputSDK) Initialize(pluginName string) error {
	if sdk.initialized {
		return fmt.Errorf("SDK already initialized")
	}
	
	if sdk.useRealSDK {
		log.Printf("Initializing REAL DirectOutput SDK with plugin: %s", pluginName)
	} else {
		log.Printf("Initializing cross-platform DirectOutput SDK with plugin: %s", pluginName)
	}
	
	sdk.initialized = true
	return nil
}

func (sdk *DirectOutputSDK) AddPage(deviceHandle unsafe.Pointer, page uint32, name string, flags uint32) error {
	if sdk.useRealSDK {
		log.Printf("DirectOutput_AddPage (REAL): device=%p, page=%d, name=%s, flags=0x%08X", deviceHandle, page, name, flags)
	} else {
		log.Printf("DirectOutput_AddPage (simulated): device=%p, page=%d, name=%s, flags=0x%08X", deviceHandle, page, name, flags)
	}
	return nil
}

func (sdk *DirectOutputSDK) SetImage(deviceHandle unsafe.Pointer, page uint32, index uint32, data []byte) error {
	if sdk.useRealSDK {
		log.Printf("DirectOutput_SetImage (REAL): device=%p, page=%d, index=%d, size=%d", deviceHandle, page, index, len(data))
	} else {
		log.Printf("DirectOutput_SetImage (simulated): device=%p, page=%d, index=%d, size=%d", deviceHandle, page, index, len(data))
	}
	return nil
}

func (sdk *DirectOutputSDK) RemovePage(deviceHandle unsafe.Pointer, page uint32) error {
	if sdk.useRealSDK {
		log.Printf("DirectOutput_RemovePage (REAL): device=%p, page=%d", deviceHandle, page)
	} else {
		log.Printf("DirectOutput_RemovePage (simulated): device=%p, page=%d", deviceHandle, page)
	}
	return nil
}

func (sdk *DirectOutputSDK) Close() error {
	if sdk.initialized {
		log.Printf("DirectOutput SDK closed")
		sdk.initialized = false
	}
	return nil
}

// createTestImages creates test images of different sizes
func createTestImages() []image.Image {
	images := []image.Image{}
	
	// Test image 1: Correct size (320x240)
	img1 := image.NewRGBA(image.Rect(0, 0, 320, 240))
	for y := 0; y < 240; y++ {
		for x := 0; x < 320; x++ {
			img1.Set(x, y, color.RGBA{uint8(x), uint8(y), 128, 255})
		}
	}
	images = append(images, img1)

	// Test image 2: Large image (640x480)
	img2 := image.NewRGBA(image.Rect(0, 0, 640, 480))
	for y := 0; y < 480; y++ {
		for x := 0; x < 640; x++ {
			img2.Set(x, y, color.RGBA{uint8(x/2), uint8(y/2), 255, 255})
		}
	}
	images = append(images, img2)

	// Test image 3: Small image (160x120)
	img3 := image.NewRGBA(image.Rect(0, 0, 160, 120))
	for y := 0; y < 120; y++ {
		for x := 0; x < 160; x++ {
			img3.Set(x, y, color.RGBA{255, uint8(x*2), uint8(y*2), 255})
		}
	}
	images = append(images, img3)

	// Test image 4: Wide image (640x240)
	img4 := image.NewRGBA(image.Rect(0, 0, 640, 240))
	for y := 0; y < 240; y++ {
		for x := 0; x < 640; x++ {
			img4.Set(x, y, color.RGBA{uint8(x/2), 128, uint8(y), 255})
		}
	}
	images = append(images, img4)

	// Test image 5: Tall image (320x480)
	img5 := image.NewRGBA(image.Rect(0, 0, 320, 480))
	for y := 0; y < 480; y++ {
		for x := 0; x < 320; x++ {
			img5.Set(x, y, color.RGBA{128, uint8(x*2), uint8(y/2), 255})
		}
	}
	images = append(images, img5)

	return images
}