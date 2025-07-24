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
#import <IOKit/usb/IOUSBLib.h>
#import <CoreFoundation/CoreFoundation.h>

IOUSBDeviceInterface** findFIPUSBDevice() {
    CFMutableDictionaryRef matchingDict = IOServiceMatching(kIOUSBDeviceClassName);
    if (!matchingDict) {
        return NULL;
    }

    // Set vendor and product ID for Saitek FIP
    CFNumberRef vendorID = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &(int){0x06A3});
    CFNumberRef productID = CFNumberCreate(kCFAllocatorDefault, kCFNumberIntType, &(int){0xA2AE});

    CFDictionarySetValue(matchingDict, CFSTR(kUSBVendorID), vendorID);
    CFDictionarySetValue(matchingDict, CFSTR(kUSBProductID), productID);

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

    // Get the USB device interface
    IOCFPlugInInterface** plugInInterface = NULL;
    SInt32 score;
    result = IOCreatePlugInInterfaceForService(service, kIOUSBDeviceUserClientTypeID, kIOCFPlugInInterfaceID, &plugInInterface, &score);

    if (result != kIOReturnSuccess) {
        IOObjectRelease(service);
        return NULL;
    }

    IOUSBDeviceInterface** deviceInterface = NULL;
    HRESULT res = (*plugInInterface)->QueryInterface(plugInInterface, CFUUIDGetUUIDBytes(kIOUSBDeviceInterfaceID), (LPVOID*)&deviceInterface);

    (*plugInInterface)->Release(plugInInterface);
    IOObjectRelease(service);

    if (res || !deviceInterface) {
        return NULL;
    }

    return deviceInterface;
}

int openFIPUSBDevice(IOUSBDeviceInterface** device) {
    if (!device) {
        return -1;
    }

    IOReturn result = (*device)->USBDeviceOpen(device);
    return (int)result;
}

int sendDirectOutputUSB(IOUSBDeviceInterface** device, unsigned char* data, int dataSize, int endpoint) {
    if (!device) {
        return -1;
    }

    // Try to send data via USB control transfer
    IOUSBDevRequest request;
    request.bmRequestType = 0x40; // Host to device, vendor request
    request.bRequest = 0x01;      // Custom request
    request.wValue = 0x0000;      // Value
    request.wIndex = 0x0000;      // Index
    request.wLength = dataSize;    // Data length
    request.pData = data;          // Data pointer

    IOReturn result = (*device)->ControlRequest(device, 0, &request);
    return (int)result;
}

int sendDirectOutputBulk(IOUSBDeviceInterface** device, unsigned char* data, int dataSize, int endpoint) {
    if (!device) {
        return -1;
    }

    // Try to send data via bulk transfer
    UInt32 numBytes = dataSize;
    IOReturn result = (*device)->WritePipe(device, endpoint, data, &numBytes);
    return (int)result;
}

void closeFIPUSBDevice(IOUSBDeviceInterface** device) {
    if (device) {
        (*device)->USBDeviceClose(device);
        (*device)->Release(device);
    }
}
*/
import "C"

func main() {
	fmt.Println("FIP USB Direct Communication Test")
	fmt.Println("=================================")
	fmt.Println()

	// Find the FIP USB device
	fmt.Println("1. Searching for FIP USB device...")
	device := C.findFIPUSBDevice()
	if device == nil {
		fmt.Println("✗ FIP USB device not found")
		return
	}

	fmt.Println("✓ FIP USB device found!")

	// Try to open the device
	fmt.Println("\n2. Attempting to open FIP USB device...")
	result := C.openFIPUSBDevice(device)
	if result != 0 {
		fmt.Printf("✗ Failed to open FIP USB device: error code %d\n", result)
		return
	}

	fmt.Println("✓ Successfully opened FIP USB device!")
	fmt.Println("\n3. Testing direct USB communication...")
	fmt.Println("Press Ctrl+C to exit")
	fmt.Println()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Test different USB communication approaches
	go testDirectUSBCommunication(device)

	// Wait for interrupt
	<-sigChan
	fmt.Println("\nShutting down...")
	C.closeFIPUSBDevice(device)
}

func testDirectUSBCommunication(device C.IOUSBDeviceInterface) {
	fmt.Println("   Testing direct USB communication...")

	// Test 1: USB control transfers
	testUSBControlTransfers(device)

	// Test 2: USB bulk transfers
	testUSBBulkTransfers(device)

	// Test 3: DirectOutput-like commands via USB
	testDirectOutputUSBCommands(device)
}

func testUSBControlTransfers(device C.IOUSBDeviceInterface) {
	fmt.Println("   Testing USB control transfers...")

	// Test different control transfer commands
	testCommands := [][]byte{
		{0x01, 0x00, 0x00, 0x00}, // Init command
		{0x02, 0x00, 0x00, 0x00}, // Display command
		{0x03, 0x00, 0x00, 0x00}, // Page command
		{0x04, 0x00, 0x00, 0x00}, // Image command
	}

	for i, cmd := range testCommands {
		fmt.Printf("   Testing control transfer %d: %v\n", i+1, cmd)
		result := C.sendDirectOutputUSB(device, (*C.uchar)(&cmd[0]), C.int(len(cmd)), 0)
		if result == 0 {
			fmt.Printf("   ✓ Control transfer %d successful!\n", i+1)
		} else {
			fmt.Printf("   ✗ Control transfer %d failed: %d\n", i+1, result)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testUSBBulkTransfers(device C.IOUSBDeviceInterface) {
	fmt.Println("   Testing USB bulk transfers...")

	// Test different bulk transfer endpoints
	testData := []byte{0x01, 0x02, 0x03, 0x04, 0x05}

	for endpoint := 1; endpoint <= 4; endpoint++ {
		fmt.Printf("   Testing bulk transfer to endpoint %d\n", endpoint)
		result := C.sendDirectOutputBulk(device, (*C.uchar)(&testData[0]), C.int(len(testData)), C.int(endpoint))
		if result == 0 {
			fmt.Printf("   ✓ Bulk transfer to endpoint %d successful!\n", endpoint)
		} else {
			fmt.Printf("   ✗ Bulk transfer to endpoint %d failed: %d\n", endpoint, result)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func testDirectOutputUSBCommands(device C.IOUSBDeviceInterface) {
	fmt.Println("   Testing DirectOutput USB commands...")

	// Based on research, try some DirectOutput-like USB commands
	directOutputCommands := [][]byte{
		{0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x4F, 0x55, 0x54, 0x50, 0x55, 0x54}, // "DIRECTOUTPUT"
		{0x53, 0x45, 0x54, 0x49, 0x4D, 0x41, 0x47, 0x45},                         // "SETIMAGE"
		{0x41, 0x44, 0x44, 0x50, 0x41, 0x47, 0x45},                               // "ADDPAGE"
		{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},                         // Page 1
		{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},                         // Page 2
	}

	for i, cmd := range directOutputCommands {
		fmt.Printf("   Testing DirectOutput USB command %d: %v\n", i+1, cmd)
		result := C.sendDirectOutputUSB(device, (*C.uchar)(&cmd[0]), C.int(len(cmd)), 0)
		if result == 0 {
			fmt.Printf("   ✓ DirectOutput USB command %d successful!\n", i+1)
		} else {
			fmt.Printf("   ✗ DirectOutput USB command %d failed: %d\n", i+1, result)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
