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

int readFIPDataWithReportID(IOHIDDeviceRef device, unsigned char* buffer, int bufferSize, int reportID) {
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
	fmt.Println("FIP Report ID Debug Tool")
	fmt.Println("========================")
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
	fmt.Println("\n3. Testing different report IDs...")
	fmt.Println("Press buttons on the FIP device while testing...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Test different report IDs
	go testReportIDs(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}

func testReportIDs(device C.IOHIDDeviceRef) {
	// Test report IDs 0-15 (common range for HID devices)
	for reportID := 0; reportID <= 15; reportID++ {
		fmt.Printf("\n🔍 Testing Report ID %d:\n", reportID)

		// Try to read from this report ID
		buffer := make([]C.uchar, 8)
		bytesRead := C.readFIPDataWithReportID(device, &buffer[0], C.int(len(buffer)), C.int(reportID))

		if bytesRead > 0 {
			data := make([]byte, bytesRead)
			for i := 0; i < int(bytesRead); i++ {
				data[i] = byte(buffer[i])
			}

			fmt.Printf("   ✓ Success! Data (len=%d): %v\n", len(data), data)
			fmt.Printf("   Hex: %02X\n", data)

			// Check if any non-zero data
			hasNonZero := false
			for _, b := range data {
				if b != 0 {
					hasNonZero = true
					break
				}
			}

			if hasNonZero {
				fmt.Printf("   🎯 NON-ZERO DATA FOUND!\n")
			}
		} else {
			fmt.Printf("   ✗ No data or error\n")
		}

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n✅ Report ID testing complete!")
	fmt.Println("Now monitoring for changes...")

	// Monitor the most promising report IDs
	monitorPromisingReports(device)
}

func monitorPromisingReports(device C.IOHIDDeviceRef) {
	// Monitor report IDs that might contain button data
	promisingReports := []int{0, 1, 2, 3} // Common report IDs for buttons

	for _, reportID := range promisingReports {
		go monitorReportID(device, reportID)
	}

	// Keep monitoring
	select {}
}

func monitorReportID(device C.IOHIDDeviceRef, reportID int) {
	lastData := make([]byte, 8)
	lastDataLen := 0
	changeCount := 0

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			buffer := make([]C.uchar, 8)
			bytesRead := C.readFIPDataWithReportID(device, &buffer[0], C.int(len(buffer)), C.int(reportID))

			if bytesRead > 0 {
				data := make([]byte, bytesRead)
				for i := 0; i < int(bytesRead); i++ {
					data[i] = byte(buffer[i])
				}

				// Check if data has changed
				if !dataEqual(data, lastData[:lastDataLen]) {
					changeCount++
					fmt.Printf("\n🔍 Report %d CHANGE #%d (len=%d): %v\n", reportID, changeCount, len(data), data)
					fmt.Printf("   Hex: %02X\n", data)

					// Store for comparison
					copy(lastData, data)
					lastDataLen = len(data)
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
