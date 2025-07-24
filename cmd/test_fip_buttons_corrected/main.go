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

// Button names based on typical FIP layout
var buttonNames = []string{
	"Button 1", "Button 2", "Button 3", "Button 4",
	"Button 5", "Button 6", "Button 7", "Button 8",
	"Button 9", "Button 10", "Button 11", "Button 12",
}

func main() {
	fmt.Println("FIP Button Test (Corrected)")
	fmt.Println("===========================")
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
	fmt.Println("\n3. Monitoring button presses...")
	fmt.Println("Press buttons on the FIP device to test...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Monitor button presses
	go monitorButtonPresses(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPDevice(device)
}

func monitorButtonPresses(device C.IOHIDDeviceRef) {
	lastButtonStates := make([]bool, 12)
	changeCount := 0

	ticker := time.NewTicker(50 * time.Millisecond) // 20Hz polling
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Read 2-byte report (based on MaxInputReportSize = 2)
			buffer := make([]C.uchar, 2)
			bytesRead := C.readFIPData(device, &buffer[0], C.int(len(buffer)))

			if bytesRead > 0 {
				data := make([]byte, bytesRead)
				for i := 0; i < int(bytesRead); i++ {
					data[i] = byte(buffer[i])
				}

				// Process button states based on device capabilities
				// Each bit in the 2-byte report represents a button (12 buttons total)
				currentButtonStates := make([]bool, 12)

				// First byte: buttons 1-8
				for i := 0; i < 8; i++ {
					currentButtonStates[i] = (data[0] & (1 << i)) != 0
				}

				// Second byte: buttons 9-12 (only first 4 bits)
				for i := 0; i < 4; i++ {
					currentButtonStates[i+8] = (data[1] & (1 << i)) != 0
				}

				// Check for changes
				for i := 0; i < 12; i++ {
					if currentButtonStates[i] != lastButtonStates[i] {
						changeCount++

						action := "pressed"
						if !currentButtonStates[i] {
							action = "released"
						}

						fmt.Printf("\n🎯 %s %s (Button %d)\n", buttonNames[i], action, i+1)
						fmt.Printf("   Raw Data: [%02X %02X]\n", data[0], data[1])
						fmt.Printf("   Binary: %08b %08b\n", data[0], data[1])

						// Show all button states
						fmt.Printf("   All Buttons: ")
						for j := 0; j < 12; j++ {
							if currentButtonStates[j] {
								fmt.Printf("●")
							} else {
								fmt.Printf("○")
							}
							if j == 7 {
								fmt.Printf(" ") // Space between bytes
							}
						}
						fmt.Println()
					}
				}

				// Update last states
				copy(lastButtonStates, currentButtonStates)

				// Show periodic status
				if changeCount%10 == 0 && changeCount > 0 {
					fmt.Printf("   Changes detected: %d\n", changeCount)
				}
			}
		}
	}
}
