package usb

import (
	"fmt"
	"log"

	"github.com/google/gousb"
)

// GoUSBDevice represents a USB device using the gousb library
type GoUSBDevice struct {
	VendorID  uint16
	ProductID uint16
	Name      string
	device    *gousb.Device
	ctx       *gousb.Context
}

// NewGoUSBDevice creates a new USB device using gousb
func NewGoUSBDevice(vendorID, productID uint16) (*GoUSBDevice, error) {
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

	// Set the active configuration
	if err := dev.SetAutoDetach(true); err != nil {
		log.Printf("Warning: failed to set auto detach: %v", err)
	}

	return &GoUSBDevice{
		VendorID:  vendorID,
		ProductID: productID,
		Name:      "Saitek Radio Panel (gousb)",
		device:    dev,
		ctx:       ctx,
	}, nil
}

// SendControlMessage sends a USB control message to the device
func (d *GoUSBDevice) SendControlMessage(requestType, request, value, index uint16, data []byte) error {
	if d.device == nil {
		return fmt.Errorf("device not initialized")
	}

	// Send control transfer
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
func (d *GoUSBDevice) ReadBulkData(endpoint uint8, length int) ([]byte, error) {
	if d.device == nil {
		return nil, fmt.Errorf("device not initialized")
	}

	// Find the interface and endpoint
	config, err := d.device.Config(1)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	iface, err := config.Interface(0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get interface: %w", err)
	}

	ep, err := iface.InEndpoint(int(endpoint))
	if err != nil {
		return nil, fmt.Errorf("failed to get endpoint: %w", err)
	}

	// Read data
	data := make([]byte, length)
	read, err := ep.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read from endpoint: %w", err)
	}

	return data[:read], nil
}

// Close closes the USB device
func (d *GoUSBDevice) Close() error {
	if d.device != nil {
		if err := d.device.Close(); err != nil {
			return err
		}
	}
	if d.ctx != nil {
		return d.ctx.Close()
	}
	return nil
}

// IsConnected returns whether the device is connected
func (d *GoUSBDevice) IsConnected() bool {
	return d.device != nil
}
