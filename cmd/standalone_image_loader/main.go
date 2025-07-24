package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ImageLoader provides functionality to load and process images for FIP display
type ImageLoader struct {
	// FIP display specifications
	FIPWidth  int
	FIPHeight int
	FIPFormat string // "RGB" or "RGBA"
	
	// Resize options
	ResizeMode ResizeMode
	Quality    int // 1-100 for JPEG quality
}

// ResizeMode defines how images should be resized
type ResizeMode int

const (
	ResizeModeStretch ResizeMode = iota // Stretch to fit (may distort)
	ResizeModeFit                       // Fit within bounds (maintain aspect ratio)
	ResizeModeCrop                      // Crop to fit (maintain aspect ratio)
	ResizeModeCenter                    // Center and pad with background
)

// NewImageLoader creates a new image loader with FIP specifications
func NewImageLoader() *ImageLoader {
	return &ImageLoader{
		FIPWidth:   320,
		FIPHeight:  240,
		FIPFormat:  "RGB",
		ResizeMode: ResizeModeFit,
		Quality:    90,
	}
}

// LoadImageFromFile loads an image from a file and processes it for FIP display
func (loader *ImageLoader) LoadImageFromFile(filename string) (image.Image, error) {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found: %s", filename)
	}

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	log.Printf("Loaded image: %s (format: %s, size: %dx%d)", filename, format, img.Bounds().Dx(), img.Bounds().Dy())

	// Process the image for FIP display
	return loader.ProcessImageForFIP(img)
}

// ProcessImageForFIP processes an image to fit FIP display requirements
func (loader *ImageLoader) ProcessImageForFIP(img image.Image) (image.Image, error) {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	log.Printf("Processing image: %dx%d -> %dx%d", originalWidth, originalHeight, loader.FIPWidth, loader.FIPHeight)

	// Check if resizing is needed
	if originalWidth == loader.FIPWidth && originalHeight == loader.FIPHeight {
		log.Printf("Image is already the correct size")
		return img, nil
	}

	// Resize the image according to the selected mode
	resizedImg, err := loader.resizeImage(img)
	if err != nil {
		return nil, fmt.Errorf("failed to resize image: %v", err)
	}

	return resizedImg, nil
}

// resizeImage resizes an image according to the selected resize mode
func (loader *ImageLoader) resizeImage(img image.Image) (image.Image, error) {
	targetWidth := loader.FIPWidth
	targetHeight := loader.FIPHeight

	switch loader.ResizeMode {
	case ResizeModeStretch:
		return loader.stretchImage(img, targetWidth, targetHeight), nil

	case ResizeModeFit:
		return loader.fitImage(img, targetWidth, targetHeight), nil

	case ResizeModeCrop:
		return loader.cropImage(img, targetWidth, targetHeight), nil

	case ResizeModeCenter:
		return loader.centerImage(img, targetWidth, targetHeight), nil

	default:
		return nil, fmt.Errorf("unknown resize mode: %d", loader.ResizeMode)
	}
}

// stretchImage stretches the image to fit the target dimensions (may distort)
func (loader *ImageLoader) stretchImage(img image.Image, targetWidth, targetHeight int) image.Image {
	resized := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.Draw(resized, resized.Bounds(), img, img.Bounds().Min, draw.Over)
	return resized
}

// fitImage fits the image within the target dimensions (maintains aspect ratio)
func (loader *ImageLoader) fitImage(img image.Image, targetWidth, targetHeight int) image.Image {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// Calculate scaling factors
	scaleX := float64(targetWidth) / float64(originalWidth)
	scaleY := float64(targetHeight) / float64(originalHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// Calculate new dimensions
	newWidth := int(float64(originalWidth) * scale)
	newHeight := int(float64(originalHeight) * scale)

	// Resize the image
	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.Draw(resized, resized.Bounds(), img, img.Bounds().Min, draw.Over)

	// Create final image with padding
	final := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	
	// Calculate centering offsets
	offsetX := (targetWidth - newWidth) / 2
	offsetY := (targetHeight - newHeight) / 2

	// Draw the resized image centered
	draw.Draw(final, image.Rect(offsetX, offsetY, offsetX+newWidth, offsetY+newHeight), resized, image.Point{}, draw.Over)

	return final
}

// cropImage crops the image to fit the target dimensions (maintains aspect ratio)
func (loader *ImageLoader) cropImage(img image.Image, targetWidth, targetHeight int) image.Image {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// Calculate scaling factors
	scaleX := float64(targetWidth) / float64(originalWidth)
	scaleY := float64(targetHeight) / float64(originalHeight)
	scale := scaleX
	if scaleY > scaleX {
		scale = scaleY
	}

	// Calculate new dimensions
	newWidth := int(float64(originalWidth) * scale)
	newHeight := int(float64(originalHeight) * scale)

	// Resize the image
	resized := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.Draw(resized, resized.Bounds(), img, img.Bounds().Min, draw.Over)

	// Crop to target dimensions
	final := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	
	// Calculate crop offsets
	offsetX := (newWidth - targetWidth) / 2
	offsetY := (newHeight - targetHeight) / 2

	// Draw the cropped portion
	draw.Draw(final, final.Bounds(), resized, image.Point{offsetX, offsetY}, draw.Over)

	return final
}

// centerImage centers the image and pads with background
func (loader *ImageLoader) centerImage(img image.Image, targetWidth, targetHeight int) image.Image {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()

	// Create final image with background
	final := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	
	// Fill with dark background
	for y := 0; y < targetHeight; y++ {
		for x := 0; x < targetWidth; x++ {
			final.Set(x, y, image.Black)
		}
	}

	// Calculate centering offsets
	offsetX := (targetWidth - originalWidth) / 2
	offsetY := (targetHeight - originalHeight) / 2

	// Ensure image fits within bounds
	if offsetX < 0 {
		offsetX = 0
	}
	if offsetY < 0 {
		offsetY = 0
	}

	// Draw the image centered
	draw.Draw(final, image.Rect(offsetX, offsetY, offsetX+originalWidth, offsetY+originalHeight), img, image.Point{}, draw.Over)

	return final
}

// ValidateImageSize checks if an image meets FIP size requirements
func (loader *ImageLoader) ValidateImageSize(img image.Image) error {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	if width != loader.FIPWidth || height != loader.FIPHeight {
		return fmt.Errorf("image size %dx%d does not match FIP requirements %dx%d", width, height, loader.FIPWidth, loader.FIPHeight)
	}

	return nil
}

// GetImageInfo returns information about an image
func (loader *ImageLoader) GetImageInfo(img image.Image) ImageInfo {
	return ImageInfo{
		Width:       img.Bounds().Dx(),
		Height:      img.Bounds().Dy(),
		FIPWidth:    loader.FIPWidth,
		FIPHeight:   loader.FIPHeight,
		NeedsResize: img.Bounds().Dx() != loader.FIPWidth || img.Bounds().Dy() != loader.FIPHeight,
	}
}

// ImageInfo contains information about an image
type ImageInfo struct {
	Width       int
	Height      int
	FIPWidth    int
	FIPHeight   int
	NeedsResize bool
}

// SetResizeMode sets the resize mode for image processing
func (loader *ImageLoader) SetResizeMode(mode ResizeMode) {
	loader.ResizeMode = mode
}

// SetQuality sets the JPEG quality for saving images
func (loader *ImageLoader) SetQuality(quality int) {
	if quality < 1 {
		quality = 1
	}
	if quality > 100 {
		quality = 100
	}
	loader.Quality = quality
}

// SaveImageAsPNG saves an image as PNG
func (loader *ImageLoader) SaveImageAsPNG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

// SaveImageAsJPEG saves an image as JPEG
func (loader *ImageLoader) SaveImageAsJPEG(img image.Image, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return jpeg.Encode(file, img, &jpeg.Options{Quality: loader.Quality})
}

// GetSupportedFormats returns a list of supported image formats
func (loader *ImageLoader) GetSupportedFormats() []string {
	return []string{".png", ".jpg", ".jpeg", ".gif"}
}

// IsSupportedFormat checks if a file format is supported
func (loader *ImageLoader) IsSupportedFormat(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	supported := loader.GetSupportedFormats()
	
	for _, format := range supported {
		if ext == format {
			return true
		}
	}
	return false
}

// LoadAndConvertToFIP loads an image and converts it to FIP format
func (loader *ImageLoader) LoadAndConvertToFIP(filename string) ([]byte, error) {
	// Load and process the image
	img, err := loader.LoadImageFromFile(filename)
	if err != nil {
		return nil, err
	}

	// Convert to FIP format
	return loader.ConvertImageToFIPFormat(img)
}

// ConvertImageToFIPFormat converts an image to FIP format (320x240, 24bpp RGB)
func (loader *ImageLoader) ConvertImageToFIPFormat(img image.Image) ([]byte, error) {
	// Create a 320x240 RGBA image
	fipImg := image.NewRGBA(image.Rect(0, 0, loader.FIPWidth, loader.FIPHeight))

	// Draw the source image onto the FIP image
	draw.Draw(fipImg, fipImg.Bounds(), img, image.Point{}, draw.Src)

	// Convert to 24bpp RGB format (FIP requirement)
	data := make([]byte, loader.FIPWidth*loader.FIPHeight*3)
	for y := 0; y < loader.FIPHeight; y++ {
		for x := 0; x < loader.FIPWidth; x++ {
			idx := (y*loader.FIPWidth + x) * 3
			c := fipImg.RGBAAt(x, y)
			data[idx] = c.R   // Red
			data[idx+1] = c.G // Green
			data[idx+2] = c.B // Blue
		}
	}

	return data, nil
}

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
	loader := NewImageLoader()

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
		fmt.Println("Usage: standalone_image_loader -image <filename> [options]")
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
		loader.SetResizeMode(ResizeModeStretch)
	case "fit":
		loader.SetResizeMode(ResizeModeFit)
	case "crop":
		loader.SetResizeMode(ResizeModeCrop)
	case "center":
		loader.SetResizeMode(ResizeModeCenter)
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

	fmt.Printf("\n✅ Image processing completed successfully!\n")
	fmt.Printf("The image has been processed for FIP display (320x240, 24bpp RGB)\n")
}