package fip

import (
	"log"
	"strings"
	"time"

	"saitek-controller/internal/usb"
)

// MultiPanel represents a Saitek Flight Multi Panel
type MultiPanel struct {
	device    usb.USBDevice
	connected bool
	vendorID  uint16
	productID uint16
}

// MultiDisplay represents the two 5-digit displays on the multi panel
type MultiDisplay struct {
	TopRow     string // Top row display (5 digits)
	BottomRow  string // Bottom row display (5 digits)
	ButtonLEDs uint8  // Button LED states
}

// Digit encoding map based on fpanels library
var multiDigitMap = map[rune]byte{
	'0': 0x00, '1': 0x01, '2': 0x02, '3': 0x03,
	'4': 0x04, '5': 0x05, '6': 0x06, '7': 0x07,
	'8': 0x08, '9': 0x09, ' ': 0x0F, '-': 0xDE,
}

// Button LED constants
const (
	ButtonAP  = 0x01 // AP button LED
	ButtonHDG = 0x02 // HDG button LED
	ButtonNAV = 0x04 // NAV button LED
	ButtonIAS = 0x08 // IAS button LED
	ButtonALT = 0x10 // ALT button LED
	ButtonVS  = 0x20 // VS button LED
	ButtonAPR = 0x40 // APR button LED
	ButtonREV = 0x80 // REV button LED
)

// NewMultiPanel creates a new multi panel
func NewMultiPanel() *MultiPanel {
	return &MultiPanel{
		vendorID:  0x06A3, // Logitech/Saitek vendor ID
		productID: 0x0D06, // Multi Panel product ID
	}
}

// NewMultiPanelWithUSB creates a new multi panel with custom vendor/product IDs
func NewMultiPanelWithUSB(vendorID, productID uint16) *MultiPanel {
	return &MultiPanel{
		vendorID:  vendorID,
		productID: productID,
	}
}

// Connect connects to the physical multi panel device
func (m *MultiPanel) Connect() error {
	// Try the USB core approach first (like the Python code)
	log.Printf("Trying USB core approach...")
	if usbDev, err := usb.NewUSBCoreDevice(m.vendorID, m.productID); err != nil {
		log.Printf("USB core approach failed: %v", err)

		// Try the standard HID approach as fallback
		log.Printf("Trying HID approach...")
		device, err := usb.OpenDevice(m.vendorID, m.productID)
		if err != nil {
			log.Printf("HID approach also failed: %v", err)
			m.connected = false
			return err
		}
		m.device = device
	} else {
		log.Printf("USB core approach succeeded!")
		m.device = usbDev
	}

	m.connected = true
	return nil
}

// Disconnect disconnects from the multi panel device
func (m *MultiPanel) Disconnect() error {
	if m.device != nil {
		m.device.Close()
		m.device = nil
	}
	m.connected = false
	return nil
}

// IsConnected returns whether the panel is connected
func (m *MultiPanel) IsConnected() bool {
	return m.connected
}

// GetType returns the panel type
func (m *MultiPanel) GetType() usb.PanelType {
	return usb.PanelTypeMulti
}

// GetName returns the panel name
func (m *MultiPanel) GetName() string {
	return "Saitek Flight Multi Panel"
}

// encodeDisplay encodes a string into 5 bytes for the display
// Supports digits 0-9, space, and dash
func encodeMultiDisplay(text string) []byte {
	result := make([]byte, 5)

	// Initialize with spaces
	for i := 0; i < 5; i++ {
		result[i] = multiDigitMap[' ']
	}

	// Process the text from left to right
	textLen := len(text)
	if textLen > 5 {
		text = text[:5] // Take first 5 characters
		textLen = 5
	}

	// Fill from left to right (the multi panel displays correctly from left to right)
	pos := 0
	for i := 0; i < textLen; i++ {
		ch := rune(text[i])
		if val, ok := multiDigitMap[ch]; ok {
			result[pos] = val
			pos++
		} else {
			result[pos] = multiDigitMap[' '] // Default to space for unknown characters
			pos++
		}
	}

	return result
}

// SendDisplay sends the display data to the multi panel
func (m *MultiPanel) SendDisplay(display MultiDisplay) error {
	if m.device == nil {
		log.Printf("Mock: Sending multi display - Top: '%s', Bottom: '%s', LEDs: 0x%02x",
			display.TopRow, display.BottomRow, display.ButtonLEDs)
		return nil
	}

	// Encode displays and create packet
	packet := make([]byte, 12) // 10 bytes for displays + 1 byte for LEDs + 1 byte for compatibility

	// Top row (first 5 bytes)
	topEncoded := encodeMultiDisplay(display.TopRow)
	copy(packet[0:5], topEncoded)

	// Bottom row (next 5 bytes)
	bottomEncoded := encodeMultiDisplay(display.BottomRow)
	copy(packet[5:10], bottomEncoded)

	// Button LEDs (11th byte)
	packet[10] = display.ButtonLEDs

	// 12th byte is always 0xFF for compatibility
	packet[11] = 0xFF

	// Debug logging
	log.Printf("Sending packet - Top encoded: %v, Bottom encoded: %v, LEDs: 0x%02x", 
		topEncoded, bottomEncoded, display.ButtonLEDs)

	// Send control message
	// bmRequestType=0x21, bRequest=0x09, wValue=0x0300, wIndex=0
	return m.device.SendControlMessage(0x21, 0x09, 0x0300, 0, packet)
}

// ReadSwitchState reads the current state of switches and encoders
func (m *MultiPanel) ReadSwitchState() ([]byte, error) {
	if m.device == nil {
		log.Printf("Mock: Reading multi panel switch state")
		return make([]byte, 3), nil
	}

	// Read 3 bytes from endpoint 1
	return m.device.ReadBulkData(1, 3)
}

// ParseSwitchState parses the switch state bytes into readable format
func (m *MultiPanel) ParseSwitchState(data []byte) map[string]bool {
	if len(data) < 3 {
		return nil
	}

	state := make(map[string]bool)

	// Byte 1 - Selection switches and encoders
	state["ALT"] = (data[0] & 0x01) != 0
	state["VS"] = (data[0] & 0x02) != 0
	state["IAS"] = (data[0] & 0x04) != 0
	state["HDG"] = (data[0] & 0x08) != 0
	state["CRS"] = (data[0] & 0x10) != 0
	state["ENCODER_CW"] = (data[0] & 0x20) != 0
	state["ENCODER_CCW"] = (data[0] & 0x40) != 0
	state["AP"] = (data[0] & 0x80) != 0

	// Byte 2 - Push buttons
	state["HDG_BTN"] = (data[1] & 0x01) != 0
	state["NAV_BTN"] = (data[1] & 0x02) != 0
	state["IAS_BTN"] = (data[1] & 0x04) != 0
	state["ALT_BTN"] = (data[1] & 0x08) != 0
	state["VS_BTN"] = (data[1] & 0x10) != 0
	state["APR_BTN"] = (data[1] & 0x20) != 0
	state["REV_BTN"] = (data[1] & 0x40) != 0
	state["THROTTLE_ARM"] = (data[1] & 0x80) != 0

	// Byte 3 - Flaps and pitch trim
	state["FLAPS_UP"] = (data[2] & 0x01) != 0
	state["FLAPS_DOWN"] = (data[2] & 0x02) != 0
	state["PITCH_DOWN"] = (data[2] & 0x04) != 0
	state["PITCH_UP"] = (data[2] & 0x08) != 0

	return state
}

// FormatValue formats a value string for display
// Handles common aviation value formats
func FormatMultiValue(value string) string {
	// Remove any non-digit characters except decimal point and dash
	value = strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' {
			return r
		}
		return -1
	}, value)

	// Ensure it's not longer than 5 characters
	if len(value) > 5 {
		value = value[:5]
	}

	return value
}

// SetDisplay sets the multi panel display with formatted values
func (m *MultiPanel) SetDisplay(topRow, bottomRow string, buttonLEDs uint8) error {
	display := MultiDisplay{
		TopRow:     FormatMultiValue(topRow),
		BottomRow:  FormatMultiValue(bottomRow),
		ButtonLEDs: buttonLEDs,
	}

	return m.SendDisplay(display)
}

// SetButtonLEDs sets the button LED states
func (m *MultiPanel) SetButtonLEDs(leds uint8) error {
	// Get current display state and update only the LEDs
	if m.device == nil {
		log.Printf("Mock: Setting button LEDs to 0x%02x", leds)
		return nil
	}

	// For now, we'll need to maintain the current display state
	// This is a simplified version - in a real implementation you'd want to cache the current display
	display := MultiDisplay{
		TopRow:     "     ", // Default empty display
		BottomRow:  "     ", // Default empty display
		ButtonLEDs: leds,
	}

	return m.SendDisplay(display)
}

// Run starts a monitoring loop for the multi panel
func (m *MultiPanel) Run() {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 Hz polling
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if m.IsConnected() {
				data, err := m.ReadSwitchState()
				if err != nil {
					log.Printf("Error reading multi panel state: %v", err)
					continue
				}

				state := m.ParseSwitchState(data)
				if state != nil {
					// Log any active switches/encoders
					for name, active := range state {
						if active {
							log.Printf("Multi Panel: %s activated", name)
						}
					}
				}
			}
		}
	}
}

// Close closes the multi panel
func (m *MultiPanel) Close() {
	m.Disconnect()
}
