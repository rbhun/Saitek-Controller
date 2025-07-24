package fip

import (
	"log"
	"strings"
	"time"

	"saitek-controller/internal/usb"
)

// RadioPanel represents a Saitek Flight Radio Panel
type RadioPanel struct {
	device    usb.USBDevice
	connected bool
	vendorID  uint16
	productID uint16
}

// RadioDisplay represents the four 5-digit displays on the radio panel
type RadioDisplay struct {
	COM1Active  string // Top Left
	COM1Standby string // Top Right
	COM2Active  string // Bottom Left
	COM2Standby string // Bottom Right
}

// Digit encoding map based on fpanels library
var digitMap = map[rune]byte{
	'0': 0x00, '1': 0x01, '2': 0x02, '3': 0x03,
	'4': 0x04, '5': 0x05, '6': 0x06, '7': 0x07,
	'8': 0x08, '9': 0x09, ' ': 0x0F, '-': 0x0E,
}

// NewRadioPanel creates a new radio panel
func NewRadioPanel() *RadioPanel {
	return &RadioPanel{
		vendorID:  0x06A3, // Logitech/Saitek vendor ID
		productID: 0x0D05, // Radio Panel product ID
	}
}

// NewRadioPanelWithUSB creates a new radio panel with custom vendor/product IDs
func NewRadioPanelWithUSB(vendorID, productID uint16) *RadioPanel {
	return &RadioPanel{
		vendorID:  vendorID,
		productID: productID,
	}
}

// Connect connects to the physical radio panel device
func (r *RadioPanel) Connect() error {
	// Try the USB core approach first (like the Python code)
	log.Printf("Trying USB core approach...")
	if usbDev, err := usb.NewUSBCoreDevice(r.vendorID, r.productID); err != nil {
		log.Printf("USB core approach failed: %v", err)

		// Try the standard HID approach as fallback
		log.Printf("Trying HID approach...")
		device, err := usb.OpenDevice(r.vendorID, r.productID)
		if err != nil {
			log.Printf("HID approach also failed: %v", err)
			r.connected = false
			return err
		}
		r.device = device
	} else {
		log.Printf("USB core approach succeeded!")
		r.device = usbDev
	}

	r.connected = true
	return nil
}

// Disconnect disconnects from the radio panel device
func (r *RadioPanel) Disconnect() error {
	if r.device != nil {
		r.device.Close()
		r.device = nil
	}
	r.connected = false
	return nil
}

// IsConnected returns whether the panel is connected
func (r *RadioPanel) IsConnected() bool {
	return r.connected
}

// GetType returns the panel type
func (r *RadioPanel) GetType() usb.PanelType {
	return usb.PanelTypeRadio
}

// GetName returns the panel name
func (r *RadioPanel) GetName() string {
	return "Saitek Flight Radio Panel"
}

// encodeDisplay encodes a string into 5 bytes for the display
// Supports digits 0-9, space, dash, and decimal points
func encodeDisplay(text string) []byte {
	result := make([]byte, 5)

	// Initialize with spaces
	for i := 0; i < 5; i++ {
		result[i] = digitMap[' ']
	}

	// Process the text from left to right
	textLen := len(text)

	// Fill from left to right
	pos := 0
	for i := 0; i < textLen; i++ {
		ch := rune(text[i])

		if ch == '.' {
			// Decimal point: set high nibble to 0xD (1101) and keep digit value in low nibble
			// This encodes the previous digit with a decimal point
			if pos > 0 {
				// Set the high nibble to 0xD while keeping the digit value
				result[pos-1] = 0xD0 | (result[pos-1] & 0x0F)
			}
			// Don't increment pos for decimal point
		} else if val, ok := digitMap[ch]; ok {
			if pos < 5 {
				result[pos] = val
				pos++
			}
		} else {
			if pos < 5 {
				result[pos] = digitMap[' '] // Default to space for unknown characters
				pos++
			}
		}
	}

	// Ensure we use all 5 positions by padding with spaces if needed
	for pos < 5 {
		result[pos] = digitMap[' ']
		pos++
	}

	return result
}

// SendDisplay sends the display data to the radio panel
func (r *RadioPanel) SendDisplay(display RadioDisplay) error {
	if r.device == nil {
		log.Printf("Mock: Sending radio display - COM1A: '%s', COM1S: '%s', COM2A: '%s', COM2S: '%s'",
			display.COM1Active, display.COM1Standby, display.COM2Active, display.COM2Standby)
		return nil
	}

	// Encode all four displays
	packet := make([]byte, 22) // 20 bytes for displays + 2 bytes for Windows compatibility

	// Top Left (COM1 Active)
	copy(packet[0:5], encodeDisplay(display.COM1Active))

	// Top Right (COM1 Standby)
	copy(packet[5:10], encodeDisplay(display.COM1Standby))

	// Bottom Left (COM2 Active)
	copy(packet[10:15], encodeDisplay(display.COM2Active))

	// Bottom Right (COM2 Standby)
	copy(packet[15:20], encodeDisplay(display.COM2Standby))

	// Last two bytes are zero for Windows compatibility
	packet[20] = 0x00
	packet[21] = 0x00

	// Send control message
	// bmRequestType=0x21, bRequest=0x09, wValue=0x0300, wIndex=0
	return r.device.SendControlMessage(0x21, 0x09, 0x0300, 0, packet)
}

// ReadSwitchState reads the current state of switches and encoders
func (r *RadioPanel) ReadSwitchState() ([]byte, error) {
	if r.device == nil {
		log.Printf("Mock: Reading radio panel switch state")
		return make([]byte, 3), nil
	}

	// Read 3 bytes from endpoint 1
	return r.device.ReadBulkData(1, 3)
}

// ParseSwitchState parses the switch state bytes into readable format
func (r *RadioPanel) ParseSwitchState(data []byte) map[string]bool {
	if len(data) < 3 {
		return nil
	}

	state := make(map[string]bool)

	// Byte 1
	state["COM1_1"] = (data[0] & 0x01) != 0
	state["COM1_2"] = (data[0] & 0x02) != 0
	state["NAV1_1"] = (data[0] & 0x04) != 0
	state["NAV1_2"] = (data[0] & 0x08) != 0
	state["ADF_1"] = (data[0] & 0x10) != 0
	state["DME_1"] = (data[0] & 0x20) != 0
	state["XPDR_1"] = (data[0] & 0x40) != 0
	state["COM2_1"] = (data[0] & 0x80) != 0

	// Byte 2
	state["COM2_2"] = (data[1] & 0x01) != 0
	state["NAV2_1"] = (data[1] & 0x02) != 0
	state["NAV2_2"] = (data[1] & 0x04) != 0
	state["ADF_2"] = (data[1] & 0x08) != 0
	state["DME_2"] = (data[1] & 0x10) != 0
	state["XPDR_2"] = (data[1] & 0x20) != 0
	state["ACT_STBY_1"] = (data[1] & 0x40) != 0
	state["ACT_STBY_2"] = (data[1] & 0x80) != 0

	// Byte 3 - Encoders
	state["ENC1_INNER_CW"] = (data[2] & 0x01) != 0
	state["ENC1_INNER_CCW"] = (data[2] & 0x02) != 0
	state["ENC1_OUTER_CW"] = (data[2] & 0x04) != 0
	state["ENC1_OUTER_CCW"] = (data[2] & 0x08) != 0
	state["ENC2_INNER_CW"] = (data[2] & 0x10) != 0
	state["ENC2_INNER_CCW"] = (data[2] & 0x20) != 0
	state["ENC2_OUTER_CW"] = (data[2] & 0x40) != 0
	state["ENC2_OUTER_CCW"] = (data[2] & 0x80) != 0

	return state
}

// FormatFrequency formats a frequency string for display
// Handles common aviation frequency formats
func FormatFrequency(freq string) string {
	// Remove any non-digit characters except decimal point
	freq = strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' {
			return r
		}
		return -1
	}, freq)

	// Don't truncate - let the encodeDisplay function handle the full frequency
	// Aviation frequencies can be up to 5 digits (e.g., 118.25, 121.90)
	return freq
}

// SetDisplay sets the radio panel display with formatted frequencies
func (r *RadioPanel) SetDisplay(com1Active, com1Standby, com2Active, com2Standby string) error {
	display := RadioDisplay{
		COM1Active:  FormatFrequency(com1Active),
		COM1Standby: FormatFrequency(com1Standby),
		COM2Active:  FormatFrequency(com2Active),
		COM2Standby: FormatFrequency(com2Standby),
	}

	return r.SendDisplay(display)
}

// Run starts a monitoring loop for the radio panel
func (r *RadioPanel) Run() {
	ticker := time.NewTicker(100 * time.Millisecond) // 10 Hz polling
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if r.IsConnected() {
				data, err := r.ReadSwitchState()
				if err != nil {
					log.Printf("Error reading radio panel state: %v", err)
					continue
				}

				state := r.ParseSwitchState(data)
				if state != nil {
					// Log any active switches/encoders
					for name, active := range state {
						if active {
							log.Printf("Radio Panel: %s activated", name)
						}
					}
				}
			}
		}
	}
}

// Close closes the radio panel
func (r *RadioPanel) Close() {
	r.Disconnect()
}
