package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"saitek-controller/internal/fip"
)

func main() {
	// Define command line flags
	var (
		imageFile   = flag.String("image", "", "Image file to load (required)")
		resizeMode  = flag.String("resize", "fit", "Resize mode: stretch, fit, crop, center")
		quality     = flag.Int("quality", 90, "JPEG quality (1-100)")
		outputFile  = flag.String("output", "", "Output file for processed image")
		showInfo    = flag.Bool("info", false, "Show image information only")
		validate    = flag.Bool("validate", false, "Validate image size only")
		listFormats = flag.Bool("formats", false, "List supported formats")
	)
	flag.Parse()

	// Create image loader
	loader := fip.NewImageLoader()

	// Handle special flags
	if *listFormats {
		fmt.Println("Supported image formats:")
		for _, format := range loader.GetSupportedFormats() {
			fmt.Printf("  %s\n", format)
		}
		return
	}

	// Check if image file is provided
	if *imageFile == "" {
		fmt.Println("Error: Image file is required")
		fmt.Println("Usage: fip_image_loader -image <filename> [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		return
	}

	// Check if file exists
	if _, err := os.Stat(*imageFile); os.IsNotExist(err) {
		log.Fatalf("Error: Image file not found: %s", *imageFile)
	}

	// Check if format is supported
	if !loader.IsSupportedFormat(*imageFile) {
		log.Fatalf("Error: Unsupported image format: %s", filepath.Ext(*imageFile))
	}

	// Set resize mode
	switch strings.ToLower(*resizeMode) {
	case "stretch":
		loader.SetResizeMode(fip.ResizeModeStretch)
	case "fit":
		loader.SetResizeMode(fip.ResizeModeFit)
	case "crop":
		loader.SetResizeMode(fip.ResizeModeCrop)
	case "center":
		loader.SetResizeMode(fip.ResizeModeCenter)
	default:
		log.Fatalf("Error: Invalid resize mode: %s", *resizeMode)
	}

	// Set quality
	loader.SetQuality(*quality)

	// Load the image
	fmt.Printf("Loading image: %s\n", *imageFile)
	img, err := loader.LoadImageFromFile(*imageFile)
	if err != nil {
		log.Fatalf("Error loading image: %v", err)
	}

	// Get image information
	info := loader.GetImageInfo(img)
	fmt.Printf("Image information:\n")
	fmt.Printf("  Original size: %dx%d\n", info.Width, info.Height)
	fmt.Printf("  FIP size: %dx%d\n", info.FIPWidth, info.FIPHeight)
	fmt.Printf("  Needs resize: %v\n", info.NeedsResize)
	fmt.Printf("  Resize mode: %s\n", *resizeMode)

	// Validate image size if requested
	if *validate {
		err := loader.ValidateImageSize(img)
		if err != nil {
			fmt.Printf("Validation failed: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Printf("Validation passed: Image is correct size for FIP\n")
		}
		return
	}

	// Show info only if requested
	if *showInfo {
		return
	}

	// Convert to FIP format
	fmt.Printf("Converting to FIP format...\n")
	fipData, err := loader.ConvertImageToFIPFormat(img)
	if err != nil {
		log.Fatalf("Error converting to FIP format: %v", err)
	}

	fmt.Printf("FIP data size: %d bytes\n", len(fipData))

	// Save processed image if output file is specified
	if *outputFile != "" {
		fmt.Printf("Saving processed image to: %s\n", *outputFile)
		
		ext := strings.ToLower(filepath.Ext(*outputFile))
		switch ext {
		case ".png":
			err = loader.SaveImageAsPNG(img, *outputFile)
		case ".jpg", ".jpeg":
			err = loader.SaveImageAsJPEG(img, *outputFile)
		default:
			log.Fatalf("Error: Unsupported output format: %s", ext)
		}
		
		if err != nil {
			log.Fatalf("Error saving image: %v", err)
		}
		fmt.Printf("✓ Image saved successfully\n")
	}

	// Simulate sending to FIP
	fmt.Printf("Simulating FIP display...\n")
	
	// Create SDK for simulation
	sdk := NewDirectOutputSDK()
	defer sdk.Close()

	// Initialize SDK
	err = sdk.Initialize("FIP Image Loader")
	if err != nil {
		log.Printf("Warning: Failed to initialize SDK: %v", err)
	} else {
		// Create device handle
		deviceHandle := unsafe.Pointer(uintptr(0x12345678))

		// Add page
		err = sdk.AddPage(deviceHandle, 1, "Image Loader Page", 0x00000001)
		if err != nil {
			log.Printf("Warning: Failed to add page: %v", err)
		} else {
			// Send image to FIP
			err = sdk.SetImage(deviceHandle, 1, 0, fipData)
			if err != nil {
				log.Printf("Warning: Failed to send image to FIP: %v", err)
			} else {
				fmt.Printf("✓ Image sent to FIP successfully\n")
			}

			// Clean up
			sdk.RemovePage(deviceHandle, 1)
		}
	}

	fmt.Printf("\n✅ Image processing completed successfully!\n")
}

// DirectOutput SDK simulation
type DirectOutputSDK struct {
	useRealSDK bool
	initialized bool
}

func NewDirectOutputSDK() *DirectOutputSDK {
	sdk := &DirectOutputSDK{}
	sdk.useRealSDK = false // For demo purposes
	return sdk
}

func (sdk *DirectOutputSDK) Initialize(pluginName string) error {
	if sdk.initialized {
		return fmt.Errorf("SDK already initialized")
	}
	
	log.Printf("Initializing DirectOutput SDK with plugin: %s", pluginName)
	sdk.initialized = true
	return nil
}

func (sdk *DirectOutputSDK) AddPage(deviceHandle unsafe.Pointer, page uint32, name string, flags uint32) error {
	log.Printf("DirectOutput_AddPage: device=%p, page=%d, name=%s, flags=0x%08X", deviceHandle, page, name, flags)
	return nil
}

func (sdk *DirectOutputSDK) SetImage(deviceHandle unsafe.Pointer, page uint32, index uint32, data []byte) error {
	log.Printf("DirectOutput_SetImage: device=%p, page=%d, index=%d, size=%d", deviceHandle, page, index, len(data))
	return nil
}

func (sdk *DirectOutputSDK) RemovePage(deviceHandle unsafe.Pointer, page uint32) error {
	log.Printf("DirectOutput_RemovePage: device=%p, page=%d", deviceHandle, page)
	return nil
}

func (sdk *DirectOutputSDK) Close() error {
	if sdk.initialized {
		log.Printf("DirectOutput SDK closed")
		sdk.initialized = false
	}
	return nil
}