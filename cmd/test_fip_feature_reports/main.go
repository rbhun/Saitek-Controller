package main

import (
	"fmt"
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

int readFIPFeatureReport(IOHIDDeviceRef device, unsigned char* buffer, int bufferSize, int reportID) {
    if (!device) {
        return -1;
    }

    CFIndex length = bufferSize;
    IOReturn result = IOHIDDeviceGetReport(device, kIOHIDReportTypeFeature, reportID, buffer, &length);

    if (result == kIOReturnSuccess) {
        return (int)length;
    }

    return -1;
}

int readFIPInputReport(IOHIDDeviceRef device, unsigned char* buffer, int bufferSize, int reportID) {
    if (!device) {
        return -1;
    }

    CFIndex length = bufferSize;
    IOReturn result = IOHIDDeviceGetReport(device, kIOHIDReportTypeInput, reportID, buffer, &length);

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
	fmt.Println("FIP Feature Report Test Tool")
	fmt.Println("============================")
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
	fmt.Println("\n3. Testing different report types...")
	fmt.Println("Press buttons on the FIP device while testing...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Test different report types
	go testDifferentReports(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}

func testDifferentReports(device C.IOHIDDeviceRef) {
	// Test feature reports first
	fmt.Println("   Testing Feature Reports...")
	for reportID := 0; reportID <= 5; reportID++ {
		buffer := make([]C.uchar, 8)
		bytesRead := C.readFIPFeatureReport(device, &buffer[0], C.int(len(buffer)), C.int(reportID))

		if bytesRead > 0 {
			data := make([]byte, bytesRead)
			for i := 0; i < int(bytesRead); i++ {
				data[i] = byte(buffer[i])
			}
			fmt.Printf("   Feature Report %d: %v (len=%d)\n", reportID, data, len(data))
		}
	}

	fmt.Println("\n   Testing Input Reports...")
	// Test input reports with different IDs
	for reportID := 0; reportID <= 5; reportID++ {
		buffer := make([]C.uchar, 8)
		bytesRead := C.readFIPInputReport(device, &buffer[0], C.int(len(buffer)), C.int(reportID))

		if bytesRead > 0 {
			data := make([]byte, bytesRead)
			for i := 0; i < int(bytesRead); i++ {
				data[i] = byte(buffer[i])
			}
			fmt.Printf("   Input Report %d: %v (len=%d)\n", reportID, data, len(data))
		}
	}

	fmt.Println("\n   Now monitoring for changes...")

	// Monitor for changes in both feature and input reports
	go monitorFeatureReports(device)
	go monitorInputReports(device)

	// Keep monitoring
	select {}
}

func monitorFeatureReports(device C.IOHIDDeviceRef) {
	lastData := make(map[int][]byte)
	changeCount := 0

	ticker := time.NewTicker(500 * time.Millisecond) // 2Hz polling
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Test feature reports 0-2
			for reportID := 0; reportID <= 2; reportID++ {
				buffer := make([]C.uchar, 8)
				bytesRead := C.readFIPFeatureReport(device, &buffer[0], C.int(len(buffer)), C.int(reportID))

				if bytesRead > 0 {
					data := make([]byte, bytesRead)
					for i := 0; i < int(bytesRead); i++ {
						data[i] = byte(buffer[i])
					}

					// Check if data changed
					lastDataForReport, exists := lastData[reportID]
					if !exists || !dataEqual(data, lastDataForReport) {
						changeCount++
						fmt.Printf("\n🔍 Feature Report %d CHANGE #%d: %v\n", reportID, changeCount, data)
						lastData[reportID] = data
					}
				}
			}
		}
	}
}

func monitorInputReports(device C.IOHIDDeviceRef) {
	lastData := make(map[int][]byte)
	changeCount := 0

	ticker := time.NewTicker(500 * time.Millisecond) // 2Hz polling
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Test input reports 0-2
			for reportID := 0; reportID <= 2; reportID++ {
				buffer := make([]C.uchar, 8)
				bytesRead := C.readFIPInputReport(device, &buffer[0], C.int(len(buffer)), C.int(reportID))

				if bytesRead > 0 {
					data := make([]byte, bytesRead)
					for i := 0; i < int(bytesRead); i++ {
						data[i] = byte(buffer[i])
					}

					// Check if data changed
					lastDataForReport, exists := lastData[reportID]
					if !exists || !dataEqual(data, lastDataForReport) {
						changeCount++
						fmt.Printf("\n🔍 Input Report %d CHANGE #%d: %v\n", reportID, changeCount, data)
						lastData[reportID] = data
					}
				}
			}
		}
	}
}

func dataEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
