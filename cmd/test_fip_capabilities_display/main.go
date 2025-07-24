package main

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"os/signal"
	"syscall"
	"time"
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

void enumerateDeviceCapabilities(IOHIDDeviceRef device) {
    if (!device) {
        return;
    }

    // Get device properties
    CFTypeRef vendorID = IOHIDDeviceGetProperty(device, CFSTR(kIOHIDVendorIDKey));
    CFTypeRef productID = IOHIDDeviceGetProperty(device, CFSTR(kIOHIDProductIDKey));
    CFTypeRef manufacturer = IOHIDDeviceGetProperty(device, CFSTR(kIOHIDManufacturerKey));
    CFTypeRef product = IOHIDDeviceGetProperty(device, CFSTR(kIOHIDProductKey));
    CFTypeRef serialNumber = IOHIDDeviceGetProperty(device, CFSTR(kIOHIDSerialNumberKey));

    printf("Device Properties:\n");
    if (vendorID) printf("  Vendor ID: %ld\n", CFNumberGetValue(vendorID, kCFNumberIntType, NULL));
    if (productID) printf("  Product ID: %ld\n", CFNumberGetValue(productID, kCFNumberIntType, NULL));
    if (manufacturer) printf("  Manufacturer: %s\n", CFStringGetCStringPtr(manufacturer, kCFStringEncodingUTF8));
    if (product) printf("  Product: %s\n", CFStringGetCStringPtr(product, kCFStringEncodingUTF8));
    if (serialNumber) printf("  Serial: %s\n", CFStringGetCStringPtr(serialNumber, kCFStringEncodingUTF8));

    // Get HID elements
    CFArrayRef elements = IOHIDDeviceCopyMatchingElements(device, NULL, kIOHIDOptionsTypeNone);
    if (elements) {
        CFIndex count = CFArrayGetCount(elements);
        printf("HID Elements (%ld):\n", count);

        for (CFIndex i = 0; i < count; i++) {
            IOHIDElementRef element = (IOHIDElementRef)CFArrayGetValueAtIndex(elements, i);

            CFTypeRef usagePage = IOHIDElementGetProperty(element, CFSTR(kIOHIDElementUsagePageKey));
            CFTypeRef usage = IOHIDElementGetProperty(element, CFSTR(kIOHIDElementUsageKey));
            CFTypeRef reportID = IOHIDElementGetProperty(element, CFSTR(kIOHIDElementReportIDKey));
            CFTypeRef reportSize = IOHIDElementGetProperty(element, CFSTR(kIOHIDElementReportSizeKey));
            CFTypeRef reportCount = IOHIDElementGetProperty(element, CFSTR(kIOHIDElementReportCountKey));

            printf("  Element %ld:\n", i);
            if (usagePage) printf("    Usage Page: %ld\n", CFNumberGetValue(usagePage, kCFNumberIntType, NULL));
            if (usage) printf("    Usage: %ld\n", CFNumberGetValue(usage, kCFNumberIntType, NULL));
            if (reportID) printf("    Report ID: %ld\n", CFNumberGetValue(reportID, kCFNumberIntType, NULL));
            if (reportSize) printf("    Report Size: %ld\n", CFNumberGetValue(reportSize, kCFNumberIntType, NULL));
            if (reportCount) printf("    Report Count: %ld\n", CFNumberGetValue(reportCount, kCFNumberIntType, NULL));
        }

        CFRelease(elements);
    }
}

int trySendImageData(IOHIDDeviceRef device, unsigned char* data, int dataSize, int reportID, int reportType) {
    if (!device) {
        return -1;
    }

    IOReturn result;
    if (reportType == 0) {
        // Try output report
        result = IOHIDDeviceSetReport(device, kIOHIDReportTypeOutput, reportID, data, dataSize);
    } else {
        // Try feature report
        result = IOHIDDeviceSetReport(device, kIOHIDReportTypeFeature, reportID, data, dataSize);
    }

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
	fmt.Println("FIP Device Capabilities & Display Test")
	fmt.Println("=====================================")
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

	// Enumerate device capabilities
	fmt.Println("\n3. Device capabilities:")
	C.enumerateDeviceCapabilities(device)

	fmt.Println("\n4. Testing display approaches...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Test different display approaches
	go testDisplayApproaches(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}

func testDisplayApproaches(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing different display approaches...")

	// Test 1: Try different report types and IDs
	testDifferentReportTypes(device)

	// Test 2: Try sending image data in chunks
	testImageDataChunks(device)

	// Test 3: Try different data formats
	testDifferentDataFormats(device)

	// Test 4: Try DirectOutput-like commands
	testDirectOutputCommands(device)
}

func testDifferentReportTypes(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing different report types...")

	testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	// Test output reports with different IDs
	for reportID := 0; reportID <= 5; reportID++ {
		fmt.Printf("   Testing output report ID %d\n", reportID)
		result := C.trySendImageData(device, (*C.uchar)(&testData[0]), C.int(len(testData)), C.int(reportID), 0)
		if result == 0 {
			fmt.Printf("   ✓ Output report ID %d successful!\n", reportID)
		} else {
			fmt.Printf("   ✗ Output report ID %d failed: %d\n", reportID, result)
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Test feature reports with different IDs
	for reportID := 0; reportID <= 5; reportID++ {
		fmt.Printf("   Testing feature report ID %d\n", reportID)
		result := C.trySendImageData(device, (*C.uchar)(&testData[0]), C.int(len(testData)), C.int(reportID), 1)
		if result == 0 {
			fmt.Printf("   ✓ Feature report ID %d successful!\n", reportID)
		} else {
			fmt.Printf("   ✗ Feature report ID %d failed: %d\n", reportID, result)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testImageDataChunks(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing image data chunks...")

	// Create a simple test image
	img := createTestImage()
	rgbData := imageToRGB(img)

	// Try sending in different chunk sizes
	chunkSizes := []int{64, 128, 256, 512}

	for _, chunkSize := range chunkSizes {
		fmt.Printf("   Testing chunk size %d bytes\n", chunkSize)

		for i := 0; i < len(rgbData); i += chunkSize {
			end := i + chunkSize
			if end > len(rgbData) {
				end = len(rgbData)
			}

			chunk := rgbData[i:end]
			result := C.trySendImageData(device, (*C.uchar)(&chunk[0]), C.int(len(chunk)), 0, 0)

			if result == 0 {
				fmt.Printf("   ✓ Chunk %d-%d successful!\n", i, end)
			} else {
				fmt.Printf("   ✗ Chunk %d-%d failed: %d\n", i, end, result)
				break
			}

			time.Sleep(10 * time.Millisecond)
		}
	}
}

func testDifferentDataFormats(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing different data formats...")

	// Test different data formats
	testFormats := [][]byte{
		{0x52, 0x47, 0x42, 0x20},                         // "RGB "
		{0x49, 0x4D, 0x47, 0x20},                         // "IMG "
		{0x44, 0x49, 0x53, 0x50},                         // "DISP"
		{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Page 1
		{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // Page 2
	}

	for i, format := range testFormats {
		fmt.Printf("   Testing data format %d: %v\n", i+1, format)
		result := C.trySendImageData(device, (*C.uchar)(&format[0]), C.int(len(format)), 0, 0)
		if result == 0 {
			fmt.Printf("   ✓ Data format %d successful!\n", i+1)
		} else {
			fmt.Printf("   ✗ Data format %d failed: %d\n", i+1, result)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testDirectOutputCommands(device C.IOHIDDeviceRef) {
	fmt.Println("   Testing DirectOutput commands...")

	// Based on the forum post, try DirectOutput-like commands
	directOutputCommands := [][]byte{
		{0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x4F, 0x55, 0x54, 0x50, 0x55, 0x54}, // "DIRECTOUTPUT"
		{0x53, 0x45, 0x54, 0x49, 0x4D, 0x41, 0x47, 0x45},                         // "SETIMAGE"
		{0x41, 0x44, 0x44, 0x50, 0x41, 0x47, 0x45},                               // "ADDPAGE"
		{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},                         // Page 1
		{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},                         // Page 2
	}

	for i, cmd := range directOutputCommands {
		fmt.Printf("   Testing DirectOutput command %d: %v\n", i+1, cmd)
		result := C.trySendImageData(device, (*C.uchar)(&cmd[0]), C.int(len(cmd)), 0, 0)
		if result == 0 {
			fmt.Printf("   ✓ DirectOutput command %d successful!\n", i+1)
		} else {
			fmt.Printf("   ✗ DirectOutput command %d failed: %d\n", i+1, result)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func createTestImage() image.Image {
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
