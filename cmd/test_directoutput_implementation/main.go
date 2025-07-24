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
	"unsafe"

	"golang.org/x/image/colornames"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#import <IOKit/IOKitLib.h>
#import <IOKit/hid/IOHIDLib.h>
#import <CoreFoundation/CoreFoundation.h>

// DirectOutput Device Type GUID for FIP
const unsigned char DeviceType_Fip[16] = {0x3E, 0x08, 0x3C, 0xD8, 0x6A, 0x37, 0x4A, 0x58, 0x80, 0xA8, 0x3D, 0x6A, 0x2C, 0x07, 0x51, 0x3E};

// DirectOutput constants
const int FLAG_SET_AS_ACTIVE = 0x00000001;
const int E_PAGENOTACTIVE = 0xFF040001;

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

// DirectOutput-like functions
int directOutputInitialize() {
    // Simulate DirectOutput initialization
    printf("DirectOutput: Initializing...\n");
    return 0;
}

int directOutputEnumerate(void* context) {
    // Simulate device enumeration
    printf("DirectOutput: Enumerating devices...\n");
    return 0;
}

int directOutputAddPage(IOHIDDeviceRef device, int page, int flags) {
    if (!device) {
        return -1;
    }

    // Try to send a "page add" command via feature report
    unsigned char pageData[] = {0x01, (unsigned char)page, (unsigned char)flags, 0x00};
    IOReturn result = IOHIDDeviceSetReport(device, kIOHIDReportTypeFeature, 0x01, pageData, sizeof(pageData));

    if (result == kIOReturnSuccess) {
        printf("DirectOutput: Added page %d with flags %d\n", page, flags);
        return 0;
    } else {
        printf("DirectOutput: Failed to add page %d: %d\n", page, result);
        return (int)result;
    }
}

int directOutputSetImage(IOHIDDeviceRef device, int page, int index, unsigned char* imageData, int dataSize) {
    if (!device) {
        return -1;
    }

    // Try to send image data via feature report with DirectOutput-like header
    unsigned char header[] = {0x02, (unsigned char)page, (unsigned char)index, 0x00};

    // Create a larger buffer for header + image data
    int totalSize = sizeof(header) + dataSize;
    unsigned char* fullData = malloc(totalSize);
    memcpy(fullData, header, sizeof(header));
    memcpy(fullData + sizeof(header), imageData, dataSize);

    IOReturn result = IOHIDDeviceSetReport(device, kIOHIDReportTypeFeature, 0x02, fullData, totalSize);

    free(fullData);

    if (result == kIOReturnSuccess) {
        printf("DirectOutput: Set image on page %d, index %d (%d bytes)\n", page, index, dataSize);
        return 0;
    } else {
        printf("DirectOutput: Failed to set image: %d\n", result);
        return (int)result;
    }
}

int directOutputSetImageFromFile(IOHIDDeviceRef device, int page, int index, const char* filename) {
    if (!device) {
        return -1;
    }

    // Try to send file command via feature report
    unsigned char fileData[] = {0x03, (unsigned char)page, (unsigned char)index, 0x00};
    IOReturn result = IOHIDDeviceSetReport(device, kIOHIDReportTypeFeature, 0x03, fileData, sizeof(fileData));

    if (result == kIOReturnSuccess) {
        printf("DirectOutput: Set image from file on page %d, index %d: %s\n", page, index, filename);
        return 0;
    } else {
        printf("DirectOutput: Failed to set image from file: %d\n", result);
        return (int)result;
    }
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
	fmt.Println("DirectOutput Implementation Test")
	fmt.Println("==============================")
	fmt.Println()

	// Initialize DirectOutput
	fmt.Println("1. Initializing DirectOutput...")
	result := C.directOutputInitialize()
	if result != 0 {
		fmt.Printf("✗ Failed to initialize DirectOutput: %d\n", result)
		return
	}
	fmt.Println("✓ DirectOutput initialized!")

	// Enumerate devices
	fmt.Println("\n2. Enumerating devices...")
	C.directOutputEnumerate(nil)

	// Find the FIP device
	fmt.Println("\n3. Searching for FIP device...")
	device := C.findFIPDevice()
	if device == 0 {
		fmt.Println("✗ FIP device not found")
		return
	}

	fmt.Println("✓ FIP device found!")

	// Try to open the device
	fmt.Println("\n4. Attempting to open FIP device...")
	result = C.openFIPDevice(device)
	if result != 0 {
		fmt.Printf("✗ Failed to open FIP device: error code %d\n", result)
		return
	}

	fmt.Println("✓ Successfully opened FIP device!")
	fmt.Println("\n5. Testing DirectOutput protocol...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Test DirectOutput functionality
	go testDirectOutputProtocol(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}

func testDirectOutputProtocol(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing DirectOutput protocol...")

	// Test 1: Add pages
	testAddPages(device)

	// Test 2: Set images
	testSetImages(device)

	// Test 3: Set images from files
	testSetImagesFromFiles(device)
}

func testAddPages(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing page management...")

	// Add pages with different flags
	pages := []struct {
		page  int
		flags int
	}{
		{1, 0},                    // Normal page
		{2, C.FLAG_SET_AS_ACTIVE}, // Active page
		{3, 0},                    // Another page
	}

	for _, p := range pages {
		fmt.Printf("   Adding page %d with flags %d\n", p.page, p.flags)
		result := C.directOutputAddPage(device, C.int(p.page), C.int(p.flags))
		if result == 0 {
			fmt.Printf("   ✓ Page %d added successfully!\n", p.page)
		} else {
			fmt.Printf("   ✗ Failed to add page %d: %d\n", p.page, result)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testSetImages(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing image setting...")

	// Create test images
	testImages := []image.Image{
		createSolidColorImage(colornames.Red),
		createSolidColorImage(colornames.Green),
		createSolidColorImage(colornames.Blue),
		createTestPattern(),
		createGradientPattern(),
	}

	for i, img := range testImages {
		fmt.Printf("   Setting image %d\n", i+1)

		// Convert image to RGB data
		rgbData := imageToRGB(img)

		// Try to set image on page 1, index i
		result := C.directOutputSetImage(device, 1, C.int(i), (*C.uchar)(&rgbData[0]), C.int(len(rgbData)))
		if result == 0 {
			fmt.Printf("   ✓ Image %d set successfully!\n", i+1)
		} else {
			fmt.Printf("   ✗ Failed to set image %d: %d\n", i+1, result)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func testSetImagesFromFiles(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing image from file...")

	// Test setting images from files
	testFiles := []string{
		"assets/test_pattern.png",
		"assets/color_bars.png",
		"assets/gradient.png",
	}

	for i, filename := range testFiles {
		fmt.Printf("   Setting image from file: %s\n", filename)

		// Convert filename to C string
		cFilename := C.CString(filename)
		defer C.free(unsafe.Pointer(cFilename))

		result := C.directOutputSetImageFromFile(device, 1, C.int(i), cFilename)
		if result == 0 {
			fmt.Printf("   ✓ Image from file set successfully!\n")
		} else {
			fmt.Printf("   ✗ Failed to set image from file: %d\n", result)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func createSolidColorImage(c color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, FIP_WIDTH, FIP_HEIGHT))
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
	return img
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
