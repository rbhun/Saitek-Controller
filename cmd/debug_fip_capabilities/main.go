package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework IOKit -framework CoreFoundation
#import <IOKit/IOKitLib.h>
#import <IOKit/hid/IOHIDLib.h>
#import <CoreFoundation/CoreFoundation.h>

void enumerateFIPCapabilities() {
    CFMutableDictionaryRef matchingDict = IOServiceMatching(kIOHIDDeviceKey);
    if (!matchingDict) {
        printf("Failed to create matching dictionary\n");
        return;
    }

    // Set vendor and product ID for Saitek FIP
    CFNumberRef vendorID = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &(int){0x06A3});
    CFNumberRef productID = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &(int){0xA2AE});

    CFDictionarySetValue(matchingDict, CFSTR(kIOHIDVendorIDKey), vendorID);
    CFDictionarySetValue(matchingDict, CFSTR(kIOHIDProductIDKey), productID);

    io_iterator_t iterator;
    kern_return_t result = IOServiceGetMatchingServices(kIOMasterPortDefault, matchingDict, &iterator);

    if (result != kIOReturnSuccess) {
        printf("Failed to get matching services\n");
        return;
    }

    io_service_t service = IOIteratorNext(iterator);
    IOObjectRelease(iterator);

    if (!service) {
        printf("No FIP device found\n");
        return;
    }

    printf("=== FIP Device Capabilities ===\n");

    // Get device properties
    CFMutableDictionaryRef properties;
    result = IORegistryEntryCreateCFProperties(service, &properties, kCFAllocatorDefault, 0);
    if (result == kIOReturnSuccess) {
        printf("Device Properties:\n");
        CFShow(properties);
        CFRelease(properties);
    }

    // Get HID device
    IOHIDDeviceRef device = IOHIDDeviceCreate(kCFAllocatorDefault, service);
    if (device) {
        printf("\n=== HID Device Information ===\n");

        // Get input elements
        CFArrayRef elements = IOHIDDeviceCopyMatchingElements(device, NULL, kIOHIDOptionsTypeNone);
        if (elements) {
            CFIndex count = CFArrayGetCount(elements);
            printf("Input Elements: %ld\n", count);

            for (CFIndex i = 0; i < count; i++) {
                IOHIDElementRef element = (IOHIDElementRef)CFArrayGetValueAtIndex(elements, i);

                CFNumberRef usagePage = (CFNumberRef)IOHIDElementGetProperty(element, CFSTR(kIOHIDElementUsagePageKey));
                CFNumberRef usage = (CFNumberRef)IOHIDElementGetProperty(element, CFSTR(kIOHIDElementUsageKey));
                CFNumberRef reportID = (CFNumberRef)IOHIDElementGetProperty(element, CFSTR(kIOHIDElementReportIDKey));
                CFNumberRef reportSize = (CFNumberRef)IOHIDElementGetProperty(element, CFSTR(kIOHIDElementReportSizeKey));
                CFNumberRef reportCount = (CFNumberRef)IOHIDElementGetProperty(element, CFSTR(kIOHIDElementReportCountKey));

                int usagePageVal = 0, usageVal = 0, reportIDVal = 0, reportSizeVal = 0, reportCountVal = 0;

                if (usagePage) CFNumberGetValue(usagePage, kCFNumberIntType, &usagePageVal);
                if (usage) CFNumberGetValue(usage, kCFNumberIntType, &usageVal);
                if (reportID) CFNumberGetValue(reportID, kCFNumberIntType, &reportIDVal);
                if (reportSize) CFNumberGetValue(reportSize, kCFNumberIntType, &reportSizeVal);
                if (reportCount) CFNumberGetValue(reportCount, kCFNumberIntType, &reportCountVal);

                printf("  Element %ld: UsagePage=%d, Usage=%d, ReportID=%d, Size=%d, Count=%d\n",
                       i, usagePageVal, usageVal, reportIDVal, reportSizeVal, reportCountVal);
            }
            CFRelease(elements);
        }

        CFRelease(device);
    }

    IOObjectRelease(service);
}

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

void closeFIPDevice(IOHIDDeviceRef device) {
    if (device) {
        IOHIDDeviceClose(device, kIOHIDOptionsTypeNone);
        CFRelease(device);
    }
}
*/
import "C"

func main() {
	fmt.Println("FIP Device Capabilities Debug Tool")
	fmt.Println("==================================")
	fmt.Println()

	// Enumerate device capabilities
	fmt.Println("1. Enumerating FIP device capabilities...")
	C.enumerateFIPCapabilities()

	fmt.Println("\n2. Testing device connection...")

	// Find and test the device
	device := C.findFIPDevice()
	if device == 0 {
		fmt.Println("✗ FIP device not found")
		return
	}

	fmt.Println("✓ FIP device found!")

	// Try to open the device
	result := C.openFIPDevice(device)
	if result != 0 {
		fmt.Printf("✗ Failed to open FIP device: error code %d\n", result)
		return
	}

	fmt.Println("✓ Successfully opened FIP device!")
	fmt.Println("\n3. Device is ready for testing...")
	fmt.Println("Press buttons on the FIP device...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}
