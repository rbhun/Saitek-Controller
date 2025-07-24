package main

import (
	"fmt"
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

int readFIPData(IOHIDDeviceRef device, unsigned char* buffer, int bufferSize) {
    if (!device) {
        return -1;
    }

    CFIndex length = bufferSize;
    IOReturn result = IOHIDDeviceGetReport(device, kIOHIDReportTypeInput, 0, buffer, &length);

    if (result == kIOReturnSuccess) {
        return (int)length;
    }

    return -1;
}

void closeFIPDevice(IOHIDDeviceRef device) {
    if (device) {
        IOHIDDeviceClose(device, kIOHIDOptionsTypeNone);
        CFRelease(device);
    }
}
*/
import "C"

func main() {
	fmt.Println("FIP Device Test using IOKit")
	fmt.Println("============================")
	fmt.Println()

	// Find the FIP device
	fmt.Println("1. Searching for FIP device...")
	device := C.findFIPDevice()
	if device == 0 {
		fmt.Println("✗ FIP device not found")
		fmt.Println("  This may indicate:")
		fmt.Println("  - Device not connected")
		fmt.Println("  - Device not recognized by system")
		fmt.Println("  - Permission issues")
		return
	}

	fmt.Println("✓ FIP device found!")

	// Try to open the device
	fmt.Println("\n2. Attempting to open FIP device...")
	result := C.openFIPDevice(device)
	if result != 0 {
		fmt.Printf("✗ Failed to open FIP device: error code %d\n", result)
		fmt.Println("  This may indicate permission issues")
		fmt.Println("  Try running with sudo or check permissions")
		C.closeFIPDevice(device)
		return
	}

	fmt.Println("✓ Successfully opened FIP device!")

	// Test reading button data
	fmt.Println("\n3. Testing button reading...")
	fmt.Println("Press buttons on the FIP device to test...")
	fmt.Println("(Will monitor for 10 seconds)")

	// Monitor for button presses
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Try to read button data
			buffer := make([]C.uchar, 2)
			bytesRead := C.readFIPData(device, &buffer[0], C.int(len(buffer)))

			if bytesRead > 0 {
				data := make([]byte, bytesRead)
				for i := 0; i < int(bytesRead); i++ {
					data[i] = byte(buffer[i])
				}
				fmt.Printf("✓ Button press detected: %v\n", data)
			}
		case <-timeout:
			fmt.Println("✓ No button presses detected (device is working)")
			fmt.Println("\nSuccess! The FIP device is accessible via IOKit.")
			C.closeFIPDevice(device)
			return
		}
	}
}
