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
	fmt.Println("FIP Gentle Test Tool")
	fmt.Println("====================")
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
	fmt.Println("\n3. Monitoring device (gentle polling)...")
	fmt.Println("Press buttons on the FIP device to test...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Monitor device gently
	go monitorDeviceGently(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}

func monitorDeviceGently(device C.IOHIDDeviceRef) {
	lastData := []byte{0, 0}
	changeCount := 0
	readCount := 0

	// Use a slower, gentler polling rate
	ticker := time.NewTicker(200 * time.Millisecond) // 5Hz polling
	defer ticker.Stop()

	fmt.Println("   Starting gentle monitoring...")
	fmt.Println("   (Polling at 5Hz to avoid device resets)")

	for {
		select {
		case <-ticker.C:
			readCount++

			// Read data with gentle approach
			buffer := make([]C.uchar, 2)
			bytesRead := C.readFIPData(device, &buffer[0], C.int(len(buffer)))

			if bytesRead > 0 {
				data := make([]byte, bytesRead)
				for i := 0; i < int(bytesRead); i++ {
					data[i] = byte(buffer[i])
				}

				// Check if data has changed
				if data[0] != lastData[0] || data[1] != lastData[1] {
					changeCount++
					fmt.Printf("\n🔍 DATA CHANGE #%d (read #%d): [%02X %02X]\n",
						changeCount, readCount, data[0], data[1])
					fmt.Printf("   Binary: %08b %08b\n", data[0], data[1])

					// Analyze button states
					buttons := []bool{}
					for i := 0; i < 8; i++ {
						buttons = append(buttons, (data[0]&(1<<i)) != 0)
					}
					for i := 0; i < 4; i++ {
						buttons = append(buttons, (data[1]&(1<<i)) != 0)
					}

					fmt.Printf("   Buttons: ")
					for i, pressed := range buttons {
						if pressed {
							fmt.Printf("●")
						} else {
							fmt.Printf("○")
						}
						if i == 7 {
							fmt.Printf(" ") // Space between bytes
						}
					}
					fmt.Println()

					// Show which buttons changed
					for i, pressed := range buttons {
						if pressed {
							fmt.Printf("   Button %d is PRESSED\n", i+1)
						}
					}

					// Store for comparison
					copy(lastData, data)
				}
			}

			// Show periodic status (every 25 reads = 5 seconds)
			if readCount%25 == 0 {
				fmt.Printf("   Status: %d reads, %d changes detected\n", readCount, changeCount)
			}
		}
	}
}
