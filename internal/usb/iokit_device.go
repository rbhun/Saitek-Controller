package usb

/*
#cgo CFLAGS: -framework IOKit -framework CoreFoundation
#include <IOKit/IOKitLib.h>
#include <IOKit/hid/IOHIDLib.h>
#include <CoreFoundation/CoreFoundation.h>

IOHIDManagerRef createHIDManager() {
    return IOHIDManagerCreate(kCFAllocatorDefault, kIOHIDOptionsTypeNone);
}

CFSetRef copyDevices(IOHIDManagerRef manager) {
    return IOHIDManagerCopyDevices(manager);
}

CFIndex getSetCount(CFSetRef set) {
    return CFSetGetCount(set);
}

void getSetValues(CFSetRef set, const void **values) {
    CFSetGetValues(set, values);
}

IOHIDDeviceRef getDeviceFromValue(const void *value) {
    return (IOHIDDeviceRef)value;
}

CFNumberRef getDeviceVendorID(IOHIDDeviceRef device) {
    return (CFNumberRef)IOHIDDeviceGetProperty(device, CFSTR("VendorID"));
}

CFNumberRef getDeviceProductID(IOHIDDeviceRef device) {
    return (CFNumberRef)IOHIDDeviceGetProperty(device, CFSTR("ProductID"));
}

int openDevice(IOHIDDeviceRef device) {
    return IOHIDDeviceOpen(device, kIOHIDOptionsTypeNone);
}

void closeDevice(IOHIDDeviceRef device) {
    IOHIDDeviceClose(device, kIOHIDOptionsTypeNone);
}

int sendReport(IOHIDDeviceRef device, const unsigned char *data, size_t length) {
    return IOHIDDeviceSetReport(device, kIOHIDReportTypeOutput, 0, data, length);
}

// Try a different approach - use IOMasterPort
kern_return_t getMasterPort(mach_port_t *masterPort) {
    return IOMasterPort(MACH_PORT_NULL, masterPort);
}

io_iterator_t findDevices(mach_port_t masterPort, CFMutableDictionaryRef matchingDict) {
    io_iterator_t iterator;
    IOServiceGetMatchingServices(masterPort, matchingDict, &iterator);
    return iterator;
}
*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"
)

// OpenIOKitDevice opens a device using IOKit directly
func OpenIOKitDevice(vendorID, productID uint16) (*Device, error) {
	log.Printf("Attempting to open device via IOKit: vendor=0x%04x product=0x%04x", vendorID, productID)

	// Create HID manager
	manager := C.createHIDManager()
	if manager == 0 {
		return nil, fmt.Errorf("failed to create HID manager")
	}
	defer C.CFRelease(C.CFTypeRef(manager))

	// Get devices
	devices := C.copyDevices(manager)
	if devices == 0 {
		return nil, fmt.Errorf("failed to get devices")
	}
	defer C.CFRelease(C.CFTypeRef(devices))

	// Iterate through devices
	count := C.getSetCount(devices)
	values := make([]unsafe.Pointer, count)
	if count > 0 {
		C.getSetValues(devices, &values[0])
	}

	for i := C.CFIndex(0); i < count; i++ {
		device := C.getDeviceFromValue(values[i])

		// Get vendor and product IDs
		vendorIDRef := C.getDeviceVendorID(device)
		productIDRef := C.getDeviceProductID(device)

		if vendorIDRef != 0 && productIDRef != 0 {
			var vid, pid C.int
			C.CFNumberGetValue(vendorIDRef, C.kCFNumberIntType, unsafe.Pointer(&vid))
			C.CFNumberGetValue(productIDRef, C.kCFNumberIntType, unsafe.Pointer(&pid))

			if uint16(vid) == vendorID && uint16(pid) == productID {
				log.Printf("Found matching device via IOKit")

				// Try to open the device
				result := C.openDevice(device)
				if result == 0 {
					log.Printf("Successfully opened device via IOKit")
					return &Device{
						VendorID:  vendorID,
						ProductID: productID,
						Name:      "Saitek FIP (via IOKit)",
						handle:    nil, // We'll store the device reference differently
					}, nil
				} else {
					log.Printf("Failed to open device via IOKit: %d", result)
				}
			}
		}
	}

	return nil, fmt.Errorf("device not found or could not be opened via IOKit")
}
