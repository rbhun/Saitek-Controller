package usb

import (
	"fmt"
	"log"

	"github.com/google/gousb"
)

// USBCoreDevice represents a USB device using direct USB access
type USBCoreDevice struct {
	VendorID  uint16
	ProductID uint16
	Name      string
	device    *gousb.Device
	ctx       *gousb.Context
}

// NewUSBCoreDevice creates a new USB device using direct USB access
func NewUSBCoreDevice(vendorID, productID uint16) (*USBCoreDevice, error) {
	ctx := gousb.NewContext()

	// Find the device
	dev, err := ctx.OpenDeviceWithVIDPID(gousb.ID(vendorID), gousb.ID(productID))
	if err != nil {
		ctx.Close()
		return nil, fmt.Errorf("failed to find device: %w", err)
	}

	if dev == nil {
		ctx.Close()
		return nil, fmt.Errorf("device not found: vendor=0x%04x product=0x%04x", vendorID, productID)
	}

	// Set auto detach to prevent kernel driver issues
	if err := dev.SetAutoDetach(true); err != nil {
		log.Printf("Warning: failed to set auto detach: %v", err)
	}

	return &USBCoreDevice{
		VendorID:  vendorID,
		ProductID: productID,
		Name:      "Saitek Radio Panel (USB Core)",
		device:    dev,
		ctx:       ctx,
	}, nil
}

// SendControlMessage sends a USB control message to the device
func (d *USBCoreDevice) SendControlMessage(requestType, request, value, index uint16, data []byte) error {
	if d.device == nil {
		return fmt.Errorf("device not initialized")
	}

	// Send control transfer exactly like the Python code
	// bmRequestType=0x21, bRequest=0x09, wValue=0x0300, wIndex=0
	_, err := d.device.Control(
		uint8(requestType),
		uint8(request),
		uint16(value),
		uint16(index),
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to send control message: %w", err)
	}

	log.Printf("Successfully sent control message to device")
	return nil
}

// ReadBulkData reads bulk data from the device
func (d *USBCoreDevice) ReadBulkData(endpoint uint8, length int) ([]byte, error) {
	if d.device == nil {
		return nil, fmt.Errorf("device not initialized")
	}

	// For now, return empty data since we're focusing on display updates
	return make([]byte, length), nil
}

// Close closes the USB device
func (d *USBCoreDevice) Close() error {
	if d.device != nil {
		d.device.Close()
	}
	if d.ctx != nil {
		d.ctx.Close()
	}
	return nil
}

// IsConnected returns whether the device is connected
func (d *USBCoreDevice) IsConnected() bool {
	return d.device != nil
}
