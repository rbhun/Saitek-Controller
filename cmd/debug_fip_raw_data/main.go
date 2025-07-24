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
	fmt.Println("FIP Raw Data Debug Tool")
	fmt.Println("=======================")
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
	fmt.Println("\n3. Monitoring raw data...")
	fmt.Println("Press buttons on the FIP device to see data changes...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Monitor raw data
	go monitorRawData(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}

func monitorRawData(device C.IOHIDDeviceRef) {
	lastData := make([]byte, 8) // Try larger buffer
	lastDataLen := 0
	changeCount := 0

	ticker := time.NewTicker(50 * time.Millisecond) // 20Hz polling
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Read data with larger buffer
			buffer := make([]C.uchar, 8)
			bytesRead := C.readFIPData(device, &buffer[0], C.int(len(buffer)))

			if bytesRead > 0 {
				data := make([]byte, bytesRead)
				for i := 0; i < int(bytesRead); i++ {
					data[i] = byte(buffer[i])
				}

				// Check if data has changed
				if !dataEqual(data, lastData[:lastDataLen]) {
					changeCount++
					fmt.Printf("\n🔍 DATA CHANGE #%d (len=%d): %v\n", changeCount, len(data), data)
					fmt.Printf("   Hex: %02X\n", data)
					fmt.Printf("   Binary: ")
					for i, b := range data {
						if i > 0 {
							fmt.Print(" ")
						}
						fmt.Printf("%08b", b)
					}
					fmt.Println()

					// Analyze each byte
					fmt.Println("   Byte Analysis:")
					for i, b := range data {
						fmt.Printf("   Byte %d: %02X (%08b) ", i, b, b)
						// Check which bits are set
						for bit := 0; bit < 8; bit++ {
							if (b>>bit)&1 == 1 {
								fmt.Printf("Bit%d ", bit)
							}
						}
						fmt.Println()
					}

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
