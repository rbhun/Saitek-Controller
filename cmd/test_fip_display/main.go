package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/image/colornames"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#import <IOKit/IOKitLib.h>
#import <IOKit/hid/IOHIDLib.h>
#import <CoreFoundation/CoreFoundation.h>

IOHIDDeviceRef findFIPDevice() {
    CFMutableDictionaryRef matchingDict = IOServiceMatching(kIOHIDDeviceKey);
    if (!matchingDict) {
        return NULL;
    }

    // Set vendor and product ID for Saitek FIP
    CFNumberRef vendorID = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &(int){0x06A3});
    CFNumberRef productID = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &(int){0xA2AE});

    CFDictionarySetValue(matchingDict, CFSTR(kIOHIDVendorIDKey), vendorID);
    CFDictionarySetValue(matchingDict, CFSTR(kIOHIDProductIDKey), productID);

    io_iterator_t iterator;
    kern_return_t result = IOServiceGetMatchingServices(kIOMasterPortDefault, matchingDict, &iterator);

    if (result != kIOReturnSuccess) {
        return NULL;
    }

    io_service_t service = IOIteratorNext(iterator);
    IOObjectRelease(iterator);

    if (!service) {
        return NULL;
    }

    IOHIDDeviceRef device = IOHIDDeviceCreate(kCFAllocatorDefault, service);
    IOObjectRelease(service);

    return device;
}

int openFIPDevice(IOHIDDeviceRef device) {
    if (!device) {
        return -1;
    }

    IOReturn result = IOHIDDeviceOpen(device, kIOHIDOptionsTypeNone);
    return (int)result;
}

int sendImageToFIP(IOHIDDeviceRef device, unsigned char* imageData, int dataSize) {
    if (!device) {
        return -1;
    }

    // Try to send image data using output reports
    // This is a simplified approach - the real DirectOutput protocol is more complex
    IOReturn result = IOHIDDeviceSetReport(device, kIOHIDReportTypeOutput, 0, imageData, dataSize);
    return (int)result;
}

void closeFIPDevice(IOHIDDeviceRef device) {
    if (device) {
        IOHIDDeviceClose(device, kIOHIDOptionsTypeNone);
        CFRelease(device);
    }
}
*/
import "C"

const (
	FIP_WIDTH  = 320
	FIP_HEIGHT = 240
)

func main() {
	fmt.Println("FIP Display Test Tool")
	fmt.Println("=====================")
	fmt.Println()

	// Find the FIP device
	fmt.Println("1. Searching for FIP device...")
	device := C.findFIPDevice()
	if device == 0 {
		fmt.Println("✗ FIP device not found")
		return
	}

	fmt.Println("✓ FIP device found!")

	// Try to open the device
	fmt.Println("\n2. Attempting to open FIP device...")
	result := C.openFIPDevice(device)
	if result != 0 {
		fmt.Printf("✗ Failed to open FIP device: error code %d\n", result)
		return
	}

	fmt.Println("✓ Successfully opened FIP device!")
	fmt.Println("\n3. Testing display functionality...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Test different display patterns
	go testDisplayPatterns(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}

func testDisplayPatterns(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing display patterns...")

	// Test 1: Solid color patterns
	testSolidColors(device)

	// Test 2: Test pattern
	testTestPattern(device)

	// Test 3: Gradient pattern
	testGradientPattern(device)

	// Test 4: Checkerboard pattern
	testCheckerboardPattern(device)
}

func testSolidColors(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing solid colors...")

	colors := []color.Color{
		colornames.Red,
		colornames.Green,
		colornames.Blue,
		colornames.White,
		colornames.Black,
		colornames.Yellow,
		colornames.Cyan,
		colornames.Magenta,
	}

	for i, c := range colors {
		fmt.Printf("   Sending color %d: %v\n", i+1, c)
		sendColorToFIP(device, c)
		time.Sleep(2 * time.Second)
	}
}

func testTestPattern(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing test pattern...")

	// Create a test pattern image
	img := createTestPattern()
	sendImageToFIP(device, img)
	time.Sleep(3 * time.Second)
}

func testGradientPattern(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing gradient pattern...")

	// Create a gradient image
	img := createGradientPattern()
	sendImageToFIP(device, img)
	time.Sleep(3 * time.Second)
}

func testCheckerboardPattern(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing checkerboard pattern...")

	// Create a checkerboard image
	img := createCheckerboardPattern()
	sendImageToFIP(device, img)
	time.Sleep(3 * time.Second)
}

func sendColorToFIP(device C.IOHIDDeviceRef, c color.Color) {
	// Create a solid color image
	img := image.NewRGBA(image.Rect(0, 0, FIP_WIDTH, FIP_HEIGHT))
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)

	sendImageToFIP(device, img)
}

func sendImageToFIP(device C.IOHIDDeviceRef, img image.Image) {
	// Convert image to RGB format (320x240)
	resized := resizeImage(img, FIP_WIDTH, FIP_HEIGHT)

	// Convert to RGB byte array
	rgbData := imageToRGB(resized)

	// Send to FIP device
	result := C.sendImageToFIP(device, (*C.uchar)(&rgbData[0]), C.int(len(rgbData)))

	if result == 0 {
		fmt.Printf("   ✓ Image sent successfully (%d bytes)\n", len(rgbData))
	} else {
		fmt.Printf("   ✗ Failed to send image: error code %d\n", result)
	}
}

func resizeImage(img image.Image, width, height int) image.Image {
	// Simple resize by scaling
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(resized, resized.Bounds(), img, image.Point{}, draw.Src)
	return resized
}

func imageToRGB(img image.Image) []byte {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create RGB buffer (3 bytes per pixel)
	rgbData := make([]byte, width*height*3)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			// Convert from 16-bit to 8-bit
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)

			// Store in RGB format
			index := (y*width + x) * 3
			rgbData[index] = r8
			rgbData[index+1] = g8
			rgbData[index+2] = b8
		}
	}

	return rgbData
}

func createTestPattern() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, FIP_WIDTH, FIP_HEIGHT))

	// Create a test pattern
	for y := 0; y < FIP_HEIGHT; y++ {
		for x := 0; x < FIP_WIDTH; x++ {
			r := uint8((x * 255) / FIP_WIDTH)
			g := uint8((y * 255) / FIP_HEIGHT)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

func createGradientPattern() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, FIP_WIDTH, FIP_HEIGHT))

	// Create a gradient pattern
	for y := 0; y < FIP_HEIGHT; y++ {
		for x := 0; x < FIP_WIDTH; x++ {
			// Red gradient horizontally
			r := uint8((x * 255) / FIP_WIDTH)
			// Green gradient vertically
			g := uint8((y * 255) / FIP_HEIGHT)
			// Blue constant
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return img
}

func createCheckerboardPattern() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, FIP_WIDTH, FIP_HEIGHT))

	// Create a checkerboard pattern
	checkerSize := 20
	for y := 0; y < FIP_HEIGHT; y++ {
		for x := 0; x < FIP_WIDTH; x++ {
			checkerX := (x / checkerSize) % 2
			checkerY := (y / checkerSize) % 2

			if (checkerX+checkerY)%2 == 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}

	return img
}
